# Tamo CLI Application

Tamo is a command-line interface (CLI) application for task and memo management with JSON persistence. It's designed for developers and AI agents who need a simple yet powerful tool to manage tasks, checklists, and associated information.

## Overview

Tamo allows you to:
- Manage tasks with order and completion status
- Create and associate memos with tasks
- Store data in JSON format
- Parse Markdown-like syntax for task and memo creation
- Flatten tasks with memo references for AI prompts

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/zishida/tamo.git
cd tamo

# Build the application
go build -o tamo ./cmd/tamo

# Move the binary to a directory in your PATH (optional)
sudo mv tamo /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/zishida/tamo/cmd/tamo@latest
```

## Getting Started

### Initialize Tamo

Before using Tamo, you need to initialize it in your project directory:

```bash
tamo init
```

This creates a `.tamo` directory in the current directory with an empty `data.json` file.

## Basic Usage

### Task Management

```bash
# Add a task at the end of the list
tamo add task "Complete documentation" -d "Write comprehensive docs for the project"

# Add a task at the beginning of the list
tamo unshift task "High priority task" -d "This needs to be done first"

# List all tasks
tamo list tasks

# Mark a task as done
tamo done <task_id>

# Show task details
tamo show <task_id>
```

### Memo Management

```bash
# Add a memo
tamo add memo "Important Information" -c "This is important content to remember"

# List all memos
tamo list memos

# Show memo details
tamo show <memo_id>
```

### Task Workflow

```bash
# Show the next task (first undone task)
tamo next

# Mark the first task as done
tamo shift task --done

# Remove the last task
tamo pop task --rm
```

## Advanced Features

### Markdown Parsing

Tamo can create tasks and memos from Markdown files:

```bash
# Create a task from a Markdown file
tamo add task -f task_description.md

# Create a task from stdin
cat task_description.md | tamo add task --from-stdin
```

The Markdown file should follow this format:

```markdown
# Task Title

Task description goes here.

```memo
This becomes a separate memo that's linked to the task.
```

More task description.

```memo
Another memo with additional information.
```
```

### Task Flattening for AI Prompts

You can flatten a task with all its memo references expanded:

```bash
tamo flattask <task_id>
```

This is useful for creating comprehensive prompts for AI tools or getting a complete view of a task with all its associated information.

## Project Structure

```
tamo/
├── cmd/
│   └── tamo/
│       └── main.go         # Application entry point
├── internal/
│   ├── cli/
│   │   ├── cli.go          # CLI command handling
│   │   └── markdown_parser.go # Markdown parsing logic
│   ├── model/
│   │   └── model.go        # Data models (Task, Memo, Store)
│   ├── storage/
│   │   └── storage.go      # JSON persistence
│   └── utils/
│       └── utils.go        # Utility functions
├── go.mod                  # Go module file
└── README.md               # This file
```

## Data Models

- **Task**: Represents work to be done with properties like ID, title, description, order, completion status, and memo references
- **Memo**: Stores information related to tasks with properties like ID, title, and content
- **Store**: The main data structure that contains all tasks and memos

## License

[MIT License](LICENSE)
