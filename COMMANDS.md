# Tamo Command Reference

This document provides a complete reference of all available commands and their options in the Tamo CLI application.

## Table of Contents

- [Tamo Command Reference](#tamo-command-reference)
  - [Table of Contents](#table-of-contents)
  - [General Commands](#general-commands)
    - [init](#init)
    - [help](#help)
  - [Task Commands](#task-commands)
    - [add task](#add-task)
    - [push task](#push-task)
    - [unshift task](#unshift-task)
    - [list tasks](#list-tasks)
    - [show task](#show-task)
    - [edit task](#edit-task)
    - [done](#done)
    - [undone](#undone)
    - [mv (move)](#mv-move)
    - [rm task](#rm-task)
  - [Memo Commands](#memo-commands)
    - [add memo](#add-memo)
    - [list memos](#list-memos)
    - [show memo](#show-memo)
    - [edit memo](#edit-memo)
    - [rm memo](#rm-memo)
  - [Workflow Commands](#workflow-commands)
    - [pop task](#pop-task)
    - [shift task](#shift-task)
    - [next](#next)
  - [Special Commands](#special-commands)
    - [flattask](#flattask)
  - [Common Patterns](#common-patterns)
    - [ID References](#id-references)
    - [Listing Options](#listing-options)

## General Commands

### init

Initializes Tamo in the current directory.

```
tamo init
```

**Description:**
- Creates a `.tamo` directory in the current directory if it doesn't exist
- Creates an empty `data.json` file in the `.tamo` directory if it doesn't exist
- If Tamo is already initialized, displays a message and does nothing

**Options:** None

### help

Displays help information.

```
tamo help
```

**Description:**
- Shows a list of all available commands and their descriptions
- Displays usage information

**Options:** None

## Task Commands

### add task

Adds a new task.

```
tamo add task "<title>" [-d "<description>"] [-m <memo_id>,...]
tamo add task -f <filepath>
tamo add task --from-stdin
```

**Description:**
- Creates a new task with the specified title and optional description
- Can optionally link the task to existing memos
- Can create a task from a Markdown file or standard input
- When using standard addition (not from Markdown), the task is added at the end of the list (equivalent to `push task`)

**Options:**
- `-d "<description>"`: Task description
- `-m <memo_id>,...`: Comma-separated list of memo IDs to reference
- `-f <filepath>`: Create task from Markdown file
- `--from-stdin`: Create task from Markdown input on stdin

### push task

Adds a new task at the end of the list.

```
tamo push task "<title>" [-d "<description>"] [-m <memo_id>,...]
```

**Description:**
- Creates a new task with the specified title and optional description
- Places the task at the end of the list (highest order value)
- Can optionally link the task to existing memos

**Options:**
- `-d "<description>"`: Task description
- `-m <memo_id>,...`: Comma-separated list of memo IDs to reference

### unshift task

Adds a new task at the beginning of the list.

```
tamo unshift task "<title>" [-d "<description>"] [-m <memo_id>,...]
```

**Description:**
- Creates a new task with the specified title and optional description
- Places the task at the beginning of the list (lowest order value)
- Can optionally link the task to existing memos

**Options:**
- `-d "<description>"`: Task description
- `-m <memo_id>,...`: Comma-separated list of memo IDs to reference

### list tasks

Lists tasks.

```
tamo list [tasks] [--done|--undone] [--refs <memo_id>]
```

**Description:**
- Lists tasks ordered by their `order` value
- Can filter tasks by completion status and memo references
- If no subcommand is specified, defaults to listing tasks

**Options:**
- `--done`: Show only completed tasks
- `--undone`: Show only uncompleted tasks
- `--refs <memo_id>`: Show only tasks referencing the specified memo ID

### show task

Shows details of a specific task.

```
tamo show <task_id>
```

**Description:**
- Displays detailed information about the specified task
- Shows ID, title, order, status, timestamps, description, and referenced memos
- Can use either the full UUID or a prefix of the ID

**Options:** None

### edit task

Edits a task.

```
tamo edit <task_id> [--editor]
```

**Description:**
- Allows editing of a task's title, description, and memo references
- By default, uses a simple prompt-based editor
- With `--editor`, uses the system's default editor (specified by the `EDITOR` environment variable)

**Options:**
- `--editor`: Use the system's default editor

### done

Marks a task as completed.

```
tamo done <task_id>
```

**Description:**
- Sets the `done` flag of the specified task to `true`
- Updates the task's `updated_at` timestamp

**Options:** None

### undone

Marks a task as not completed.

```
tamo undone <task_id>
```

**Description:**
- Sets the `done` flag of the specified task to `false`
- Updates the task's `updated_at` timestamp

**Options:** None

### mv (move)

Moves a task to a specific order or relative to another task.

```
tamo mv <task_id> <target_order>
tamo mv <task_id> before|after <other_task_id>
```

**Description:**
- Changes the `order` value of the specified task
- Can set an absolute order value or position the task relative to another task
- Updates the task's `updated_at` timestamp

**Options:** None

### rm task

Removes a task.

```
tamo rm <task_id> [-f|--force]
```

**Description:**
- Deletes the specified task from the data store
- Can use either the full UUID or a prefix of the ID

**Options:**
- `-f, --force`: Force removal without confirmation

## Memo Commands

### add memo

Adds a new memo.

```
tamo add memo [<title>] [-c "<content>" | --from-stdin | --editor]
```

**Description:**
- Creates a new memo with the specified title (optional) and content
- Content can be provided via command-line argument, standard input, or editor

**Options:**
- `-c "<content>"`: Memo content
- `--from-stdin`: Read content from stdin
- `--editor`: Open editor to input content (not fully implemented yet)

### list memos

Lists memos.

```
tamo list memos
```

**Description:**
- Lists all memos with their ID, title (if any), and a preview of the content

**Options:** None

### show memo

Shows details of a specific memo.

```
tamo show <memo_id>
```

**Description:**
- Displays detailed information about the specified memo
- Shows ID, title (if any), timestamps, and full content
- Can use either the full UUID or a prefix of the ID

**Options:** None

### edit memo

Edits a memo.

```
tamo edit <memo_id> [--editor]
```

**Description:**
- Allows editing of a memo's title and content
- By default, uses a simple prompt-based editor
- With `--editor`, uses the system's default editor (specified by the `EDITOR` environment variable)

**Options:**
- `--editor`: Use the system's default editor

### rm memo

Removes a memo.

```
tamo rm <memo_id> [-f|--force]
```

**Description:**
- Deletes the specified memo from the data store
- If the memo is referenced by any tasks, displays a warning and requires confirmation
- Can use either the full UUID or a prefix of the ID

**Options:**
- `-f, --force`: Force removal without confirmation, even if the memo is referenced by tasks

## Workflow Commands

### pop task

Shows, marks as done, or removes the last task.

```
tamo pop task [--done | --rm [-f]]
```

**Description:**
- Operates on the last task in the list (highest order value)
- Without options, displays the task details
- With `--done`, marks the task as completed
- With `--rm`, removes the task (requires confirmation unless `-f` is specified)

**Options:**
- `--done`: Mark the last task as done
- `--rm`: Remove the last task
- `-f`: Force removal without confirmation (with `--rm`)

### shift task

Shows, marks as done, or removes the first task.

```
tamo shift task [--done | --rm [-f]]
```

**Description:**
- Operates on the first task in the list (lowest order value)
- Without options, displays the task details
- With `--done`, marks the task as completed
- With `--rm`, removes the task (requires confirmation unless `-f` is specified)

**Options:**
- `--done`: Mark the first task as done
- `--rm`: Remove the first task
- `-f`: Force removal without confirmation (with `--rm`)

### next

Shows the first undone task.

```
tamo next
```

**Description:**
- Displays the details of the first undone task (lowest order value among undone tasks)
- Equivalent to `tamo shift task` but only considers undone tasks

**Options:** None

## Special Commands

### flattask

Flattens a task by expanding all memo references.

```
tamo flattask <task_id>
```

**Description:**
- Generates a Markdown document that includes the task title, status, description, and the content of all referenced memos
- Useful for creating comprehensive prompts for AI tools or getting a complete view of a task
- Can use either the full UUID or a prefix of the ID

**Options:** None

## Common Patterns

### ID References

Most commands that operate on a specific task or memo accept either:
- The full UUID (e.g., `123e4567-e89b-12d3-a456-426614174000`)
- A prefix of the UUID (e.g., `123e4567`)

### Listing Options

The `list` command can be used with different subcommands and options:
- `tamo list` or `tamo list tasks`: List all tasks
- `tamo list memos`: List all memos
- `tamo list all`: List all tasks and memos
- `tamo list --done`: List completed tasks
- `tamo list --undone`: List uncompleted tasks
- `tamo list --refs <memo_id>`: List tasks referencing a specific memo
