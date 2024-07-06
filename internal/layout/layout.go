package layout

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"worktree/internal/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	primaryColor   = tcell.ColorDarkBlue
	secondaryColor = tcell.ColorGreenYellow
	textColor      = tcell.ColorWhite
	borderColor    = tcell.ColorYellow
)

func CreateLogView() *tview.TextView {
	logView := tview.NewTextView()
	logView.SetWrap(true)
	logView.SetWordWrap(true)
	logView.SetBorder(true)
	logView.SetTitle("Logs")
	logView.SetTitleColor(secondaryColor)
	logView.SetBorderColor(borderColor)
	logView.SetTextColor(textColor)
	logView.SetBackgroundColor(primaryColor)
	logView.SetScrollable(true)
	return logView
}

func CreateActionList() *tview.List {
	actionList := tview.NewList()
	actionList.ShowSecondaryText(false).
		SetTitle("Actions").
		SetBorder(true).
		SetTitleColor(secondaryColor).
		SetBorderColor(borderColor).
		SetBackgroundColor(primaryColor)

	return actionList
}

func CreateDirectoryList() *tview.List {
	leftList := tview.NewList()
	leftList.SetBorder(true).SetTitle("Repositories")
	leftList.SetTitleColor(secondaryColor)
	leftList.SetBorderColor(borderColor)
	leftList.SetMainTextColor(textColor)
	leftList.SetBackgroundColor(primaryColor)

	return leftList
}

func CreateAddWorktreeModal() *tview.Form {
	modal := tview.NewForm()
	modal.
		SetBorder(true).
		SetTitle("Add New Worktree").
		SetTitleColor(textColor).
		SetBorderColor(borderColor).
		SetBackgroundColor(secondaryColor)

	return modal
}

type Layout struct {
	App           *tview.Application
	Layout        *tview.Flex
	LeftList      *tview.List
	ActionList    *tview.List
	LogText       *tview.TextView
	WorktreeModal *tview.Form
}

func NewLayout(app *tview.Application, entryPoint string) *Layout {
	logView := CreateLogView()
	actionList := CreateActionList()

	rightPanel := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(actionList, 0, 2, true).
		AddItem(logView, 0, 1, false)

	leftPanel := CreateDirectoryList()
	leftPanel.SetFocusFunc(func() {
		logView.Clear()
		actionList.Clear()
		actionList.SetTitle("Actions")
		if err := os.Chdir(entryPoint); err != nil {
			fmt.Println("failed to change directory to entry point: %w", err)
		}
	})

	layout := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(leftPanel, 0, 1, true).
		AddItem(rightPanel, 0, 2, false)

	worktreeModal := CreateAddWorktreeModal()

	return &Layout{
		App:           app,
		Layout:        layout,
		LeftList:      leftPanel,
		ActionList:    actionList,
		LogText:       logView,
		WorktreeModal: worktreeModal,
	}
}

func (l *Layout) SetRoot() error {
	if err := l.App.SetRoot(l.Layout, true).Run(); err != nil {
		return fmt.Errorf("failed to set root: %w", err)
	}

	l.App.SetFocus(l.LeftList)

	return nil
}

func (l *Layout) Log(format string, args ...interface{}) {
	l.LogText.SetText(fmt.Sprintln(l.LogText.GetText(false), '\n', fmt.Sprintf(format, args...)))
}

func (l *Layout) SetupLayoutContentMenus(path string) {
	if err := utils.ChangeDirectory(path); err != nil {
		l.Log("Failed to change directory to %v: %v", path, err)
	}

	dirs, err := utils.ListDirectories(path)
	if err != nil {
		l.Log("Error listing directories: %v", err)
	}

	for _, dir := range dirs {
		mainBranchExist, err := utils.CheckMainOrMasterExists(dir)
		if err != nil && mainBranchExist == "" {
			l.Log("Can not find main or master branch for this repository, error: %v", err)
			continue
		}

		l.LeftList.AddItem(dir, "", 0, func() {
			l.setupActions(filepath.Join(path, dir))
			l.App.SetFocus(l.ActionList)
		}).SetSelectedBackgroundColor(secondaryColor)
	}
}

func (l *Layout) setupActions(path string) {
	l.ActionList.Clear()
	l.ActionList.SetTitle("Actions for " + path)

	if err := utils.ChangeDirectory(path); err != nil {
		l.Log("Failed to change directory to %v: %v", path, err)
	}
	l.Log("Changed directory to %v", utils.Pwd())

	l.ActionList.AddItem(" ..", "", 0, func() {
		l.Log("Back to Directory List")

		l.ActionList.Clear()
		l.App.SetFocus(l.LeftList)
	}).SetSelectedBackgroundColor(secondaryColor)

	l.ActionList.AddItem(" Add New Worktree", "", 0, func() {
		l.showWorktreeModal(path)
	}).SetSelectedBackgroundColor(secondaryColor)

	worktreeBranches, err := utils.GetWorktrees(l.Log)
	if err != nil {
		l.Log("Failed to get worktreeBranches: %v", err)
	}

	if len(worktreeBranches) > 0 {
		l.ActionList.AddItem("", " ", 0, nil)
		for _, worktreeBranch := range worktreeBranches {
			l.ActionList.AddItem(worktreeBranch, "", 0, func() {
				l.selectedWorktreeActionList(filepath.Join(path, worktreeBranch), worktreeBranch)
			}).SetSelectedBackgroundColor(secondaryColor).SetBlurFunc(func() {
				l.ActionList.Clear()
				l.setupActions(path)
			})
		}
	} else {
		l.ActionList.AddItem("No worktrees found, add a new worktree!", "", 0, nil).SetSelectedBackgroundColor(secondaryColor)
	}
}

func (l *Layout) selectedWorktreeActionList(path string, selectedWorktree string) {
	l.ActionList.Clear()
	l.ActionList.SetTitle("Actions for " + path)

	if err := utils.ChangeDirectory(path); err != nil {
		l.Log("Failed to change directory to %v: %v", path, err)
	}
	l.Log("Changed directory to %v", utils.Pwd())

	l.ActionList.AddItem(" ..", "", 0, func() {
		l.Log("Back to Action List")
		l.App.SetFocus(l.ActionList)
	}).SetSelectedBackgroundColor(secondaryColor)

	l.ActionList.AddItem("Open VSCode", "", 0, func() {
		if err := utils.OpenVSCode(".", l.Log); err != nil {
			l.Log("Failed to open vscode: %v", err)
		}
	}).SetSelectedBackgroundColor(secondaryColor)

	l.ActionList.AddItem("Remove Worktree", "", 0, func() {
		if err := utils.RemoveWorktree(path, selectedWorktree, l.Log); err != nil {
			l.Log("Failed to remove worktree: %v", err)
		}

		l.App.SetFocus(l.ActionList)
	}).SetSelectedBackgroundColor(secondaryColor)
}

func (l *Layout) showWorktreeModal(path string) {
	l.WorktreeModal.
		AddInputField("Enter Branch Name", "", 30, func(text string, lastChar rune) bool {
			if strings.Contains(text, " ") {
				l.Log("Branch name cannot contain space")
				return false
			}
			return true
		}, nil).
		AddButton("Add", func() {
			l.Log("Adding new worktree")
			newBranchName := l.WorktreeModal.GetFormItemByLabel("Enter Branch Name").(*tview.InputField).GetText()
			if newBranchName == "" {
				l.Log("Branch name cannot be empty")
				return
			}

			if err := utils.AddWorktree(path, newBranchName, l.Log); err != nil {
				l.Log("Failed to add worktree, error: %v", err)
				return
			}

			l.selectedWorktreeActionList(filepath.Join(path, newBranchName), newBranchName)
			_ = l.SetRoot()
		}).
		AddButton("To Cancel Press Esc", func() {
			l.Log("Canceling adding new worktree")
			_ = l.SetRoot()
		}).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyEscape {
				l.Log("Canceling adding new worktree")
				_ = l.SetRoot()
			}

			if event.Key() == tcell.KeyDown {
				l.WorktreeModal.GetButton(0)
			}
			return event
		})

	// Create a primitive that overlays the modal on top of the current root
	modalOverlay := func() tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(l.WorktreeModal, 0, 1, true).
				AddItem(nil, 0, 1, false), 0, 2, true).
			AddItem(nil, 0, 1, false)
	}

	// Set the overlay as the root and focus on the modal
	l.App.SetRoot(modalOverlay(), true).SetFocus(l.WorktreeModal)
}
