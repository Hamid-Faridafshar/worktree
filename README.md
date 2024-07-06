# Repository Manager CLI with Worktree

This Go application helps you manage multiple repositories with ease. It allows you to create worktrees for each repository and manage multiple active branches.

## Features

- Manage multiple repositories with ease
- Worktree support for efficient branch management
- Manage multiple active branches with worktree for each repository
- Interactive CLI for seamless user experience

## Installation

To use this project, clone the repository from GitHub, build it, and run it.

### Prerequisites

This application requires Go installed on your machine. If you don't have Go installed, follow the instructions at the official Go website: [golang.org](https://golang.org/doc/install).

To check if Go is installed, run the following command:

```sh
go version
```

### Build the Project

```sh
go build -o <output-name>
```

### Prepare the Environment
Set the environment variable `WT_ENTRY_POINT` to the path where your repositories are stored.

```sh
export WT_ENTRY_POINT=<path-to-your-repositories>
```

**Important note:** The application requires all your repositories to be in the same directory, and each repository to have either master or main branch.

### Example:
If you set my-repositories as your WT_ENTRY_POINT environment variable, application will show you the list of repository_1, repository_2, repository_3 and will ignore random_dir because it does not have a main/master directory and Random_file.

```
my-repositories/
├── repository_1/
│   └── main/
│       └── repository_1_files
├── repository_2/
│   └── master/
│       └── repository_2_files
├── repository_3/
│   └── master/
│       └── repository_3_files
├── random_dir/
│   └── dir_1/
│       └── dir_1_files
└── Random_file
```

### Run the Project
After building the project, you can run it from the output path.
```sh
./<output-name> 
```

## Usage
Once the application is running, you will be guided through the interactive CLI to manage your repositories with worktree.

## License
This project is licensed under the MIT License