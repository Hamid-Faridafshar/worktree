package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func ListDirectories(path string) ([]string, error) {
	var dirs []string

	files, err := os.ReadDir(path)

	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file.Name())
		}
	}

	return dirs, nil
}

func CheckMainOrMasterExists(selectedRepo string) (string, error) {
	mainExists, err := dirExists(filepath.Join(selectedRepo, "main"))
	if err != nil {
		return "", err
	}
	if mainExists {
		return "main", nil
	}
	masterExists, err := dirExists(filepath.Join(selectedRepo, "master"))

	if err != nil {
		return "", err
	}
	if masterExists {
		return "master", nil
	}

	return "", fmt.Errorf("main or master branch not found")
}

func dirExists(path string) (bool, error) {
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

func CreateNewBranch(branchName string, mainBranchName string, logger func(format string, args ...interface{})) error {
	destinationPath := filepath.Join("..", branchName)

	cmd := exec.Command("git", "worktree", "add", "-b", branchName, destinationPath, mainBranchName)
	return cmd.Run()
}

func PrepareDestinationConfig(branchName string, logger func(format string, args ...interface{})) error {
	if _, err := os.Stat("workspace.code-workspace"); os.IsNotExist(err) {
		logger("Workspace file not found")
		return nil
	}

	workspaceFile := "workspace.code-workspace"
	src := filepath.Join(".", workspaceFile)
	dest := filepath.Join("..", branchName, workspaceFile)

	err := copyFile(src, dest)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("error copying file contents: %v", err)
	}

	return nil
}

func OpenVSCode(openingDir string, logger func(format string, args ...interface{})) error {
	_, err := os.Stat(filepath.Join(openingDir, "workspace.code-workspace"))
	if !os.IsNotExist(err) {
		openingDir = filepath.Join(openingDir, "workspace.code-workspace")
	}

	logger("Opening vscode ...")
	cmd := exec.Command("code", openingDir)
	return cmd.Start()
}

func GetWorktrees(logger func(format string, args ...interface{})) ([]string, error) {
	mainBranchName, err := CheckMainOrMasterExists(".")
	if err != nil || mainBranchName == "" {
		return nil, fmt.Errorf("failed checking main or master branch: %w", err)
	}

	if err := ChangeDirectory(mainBranchName); err != nil {
		return nil, fmt.Errorf("failed to change directory to main branch %w", err)
	}

	if err := exec.Command("git", "status").Run(); err != nil {
		return nil, fmt.Errorf("failed to run git status: %w", err)
	}

	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	output, err := cmd.Output()

	if err != nil {
		return nil, fmt.Errorf("failed to get worktrees: %w", err)
	}

	lines := bytes.Split(output, []byte("\n"))
	var branches []string
	for i := 0; i < len(lines); i++ {
		if bytes.HasPrefix(lines[i], []byte("worktree ")) {
			for j := i + 1; j < len(lines) && !bytes.HasPrefix(lines[j], []byte("worktree ")); j++ {
				if bytes.HasPrefix(lines[j], []byte("branch ")) {
					branch := strings.TrimSpace(string(bytes.TrimPrefix(lines[j], []byte("branch refs/heads/"))))
					branches = append(branches, branch)
				}
			}
		}
	}

	if err := ChangeDirectory(".."); err != nil {
		return nil, fmt.Errorf("failed to change directory to parent: %w", err)
	}

	return branches, nil
}

func AddWorktree(path string, newBranchName string, logger func(format string, args ...interface{})) error {
	if err := ChangeDirectory(path); err != nil {
		return fmt.Errorf("failed to change directory to %v: %w", path, err)
	}

	mainBranchName, err := CheckMainOrMasterExists(".")
	if err != nil || mainBranchName == "" {
		logger("Failed checking main or master branch, error: %v", err)
		return err
	}

	if err := ChangeDirectory(mainBranchName); err != nil {
		logger("Failed to go to %s dir, error: %v", mainBranchName, err)

		return err
	}

	logger("Adding worktree %v", newBranchName)
	if err := CreateNewBranch(newBranchName, mainBranchName, logger); err != nil {
		logger("Failed to create a new branch, error: %v", err)

		return err
	}

	if err := PrepareDestinationConfig(newBranchName, logger); err != nil {
		logger("Failed to prepare destination config, error: %v", err)

		return err
	}

	logger("Worktree %v added successfully", newBranchName)

	return nil
}

func RemoveWorktree(path string, branchName string, logger func(format string, args ...interface{})) error {
	dirLevel := strings.Count(branchName, "/") + 1
	if err := ChangeDirectory(strings.Repeat("../", dirLevel)); err != nil {
		logger("Failed to change directory to repository, error: %v", err)
	}

	mainBranchName, err := CheckMainOrMasterExists(".")
	if err != nil || mainBranchName == "" {
		logger("Failed checking main or master branch, error: %v", err)
	}

	if err := ChangeDirectory(mainBranchName); err != nil {
		logger("Failed to change directory to main branch, error: %v", err)
	}

	logger("Removing worktree %v", path)
	if err := execCommand("git", "worktree", "remove", path); err != nil {
		return fmt.Errorf("failed to remove worktree: %w", err)
	}

	logger("Removing git branch %v", branchName)
	if err := execCommand("git", "branch", "-D", branchName); err != nil {
		return fmt.Errorf("failed to delete branch: %w", err)
	}

	return nil
}

func execCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	var outBuffer, errBuffer bytes.Buffer
	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(60 * time.Second):
		if err := cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill process: %w", err)
		}
		return fmt.Errorf("command timed out")
	case err := <-done:
		if err != nil {
			return fmt.Errorf("command failed: %s", errBuffer.String())
		}
	}

	return nil
}

func Pwd() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Failed to get current directory, error: %s", err)
		return ""
	}

	return dir
}

func ChangeDirectory(path string) error {
	if err := os.Chdir(path); err != nil {
		return fmt.Errorf("failed to change directory to %v: %w", path, err)
	}

	return nil
}
