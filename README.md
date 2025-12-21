# Worktree - Git Worktree Manager

A terminal UI application for managing multiple Git repositories using Git's native worktree feature. Create isolated branches with separate working directories and open them in your favorite IDE.

**[View Interactive Demo](https://hamid-faridafshar.github.io/worktree-landing/)**

## Preview

**Main Screen** - Repository list with actions and worktrees:
```
┌─Repositories────────────┐┌─Actions──────────────────────────┐
│                         ││                                  │
│  my-project            ▶││  ..                              │
│  another-repo           ││  Settings                        │
│  api-service            ││  Add New Worktree                │
│                         ││  ─────────────────────────────── │
│                         ││  feature/login                   │
│                         ││  bugfix/header                   │
│                         ││  refactor/api                    │
│                         │├─Logs─────────────────────────────┤
│                         ││ Changed directory to /repos/my..│
│                         ││                                  │
└─────────────────────────┘└──────────────────────────────────┘
```

**Worktree Actions** - Select a worktree to see available IDEs:
```
┌─Actions for feature/login────────────────┐
│                                          │
│  ..                                      │
│  Open Cursor                             │
│  Open VS Code                            │
│  Remove Worktree                         │
│                                          │
└──────────────────────────────────────────┘
```

**Settings Modal** - Configure which IDEs to show:
```
┌─Settings - Toggle IDEs───────────────────┐
│                                          │
│  [x] Cursor                              │
│  [x] VS Code                             │
│  [ ] Zed                                 │
│  [ ] Sublime Text                        │
│  [ ] Neovim                              │
│                                          │
│  [Save]  [Cancel]                        │
│                                          │
└──────────────────────────────────────────┘
```

**Add Worktree Modal** - Create a new worktree:
```
┌─Add New Worktree─────────────────────────┐
│                                          │
│  Enter Branch Name: feature/new-feature  │
│                                          │
│  [Add]  [To Cancel Press Esc]            │
│                                          │
└──────────────────────────────────────────┘
```

## Features

- **Multi-repository management** - Scan and manage all repositories in a directory
- **Worktree support** - Create and remove Git worktrees for parallel branch development
- **Configurable IDE support** - Choose from Cursor, VS Code, Zed, Sublime Text, or Neovim
- **Persistent settings** - IDE preferences saved to `~/.config/worktree/config.json`
- **Interactive terminal UI** - Keyboard-driven navigation

## Installation

### Prerequisites

- Go 1.19 or later
- Git with worktree support

```sh
go version
```

### Build

```sh
git clone https://github.com/Hamid-Faridafshar/worktree.git
cd worktree
go build -o worktree
```

### Run

```sh
# Using environment variable
export WT_ENTRY_POINT=/path/to/your/repositories
./worktree

# Or using command-line flag
./worktree --entry-point /path/to/your/repositories
```

## Repository Structure

The application expects repositories to be organized in a parent directory, each with a `main` or `master` branch:

```
my-repositories/
├── project-1/
│   └── main/           # Required: main or master branch
│       └── ...
├── project-2/
│   └── master/         # Required: main or master branch
│       └── ...
└── not-a-repo/         # Ignored: no main/master directory
    └── random-files/
```

## Configuration

IDE preferences are stored in `~/.config/worktree/config.json`:

```json
{
  "editors": {
    "cursor": {
      "enabled": true,
      "command": "cursor",
      "displayName": "Cursor"
    },
    "vscode": {
      "enabled": true,
      "command": "code",
      "displayName": "VS Code"
    },
    "zed": {
      "enabled": false,
      "command": "zed",
      "displayName": "Zed"
    },
    "sublime": {
      "enabled": false,
      "command": "subl",
      "displayName": "Sublime Text"
    },
    "neovim": {
      "enabled": false,
      "command": "nvim",
      "displayName": "Neovim"
    }
  }
}
```

On first run, a default config is created with Cursor and VS Code enabled.

## Usage

### Workflow

1. **Launch** the app - see your repositories in the left panel
2. **Select a repository** - view its worktrees and actions in the right panel
3. **Open Settings** - toggle which IDEs appear in the menu
4. **Add New Worktree** - create a new branch with its own working directory
5. **Select a worktree** - choose an IDE to open the code

### Keyboard Navigation

| Key | Action |
|-----|--------|
| `↑` `↓` | Navigate lists |
| `Enter` | Select item |
| `Esc` | Cancel / Go back |
| `Tab` | Move between form fields |

## Workspace Support

If a `workspace.code-workspace` file exists in the main branch, it will be:
- Automatically copied to new worktrees
- Used when opening the worktree in VS Code or Cursor

## License

This project is licensed under the MIT License.
