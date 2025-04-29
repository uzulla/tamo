# Tamo Usage Guide

This guide provides detailed examples of how to use each feature of the Tamo CLI application.

## Table of Contents

- [Tamo Usage Guide](#tamo-usage-guide)
  - [Table of Contents](#table-of-contents)
  - [Initialization](#initialization)
  - [Task Management](#task-management)
    - [Adding Tasks](#adding-tasks)
      - [Standard Task Addition (at the end)](#standard-task-addition-at-the-end)
      - [Adding a Task with Memo References](#adding-a-task-with-memo-references)
      - [Adding a Task at the End (Push)](#adding-a-task-at-the-end-push)
      - [Adding a Task at the Beginning (Unshift)](#adding-a-task-at-the-beginning-unshift)
    - [Listing Tasks](#listing-tasks)
      - [List All Tasks](#list-all-tasks)
      - [List Only Completed Tasks](#list-only-completed-tasks)
      - [List Only Uncompleted Tasks](#list-only-uncompleted-tasks)
      - [List Tasks Referencing a Specific Memo](#list-tasks-referencing-a-specific-memo)
    - [Showing Task Details](#showing-task-details)
    - [Editing Tasks](#editing-tasks)
    - [Marking Tasks as Done/Undone](#marking-tasks-as-doneundone)
    - [Moving Tasks](#moving-tasks)
      - [Moving a Task to a Specific Order](#moving-a-task-to-a-specific-order)
      - [Moving a Task Relative to Another Task](#moving-a-task-relative-to-another-task)
    - [Removing Tasks](#removing-tasks)
  - [Memo Management](#memo-management)
    - [Adding Memos](#adding-memos)
      - [Adding a Memo with Content](#adding-a-memo-with-content)
      - [Adding a Memo from Standard Input](#adding-a-memo-from-standard-input)
    - [Listing Memos](#listing-memos)
    - [Showing Memo Details](#showing-memo-details)
    - [Editing Memos](#editing-memos)
    - [Removing Memos](#removing-memos)
  - [Task Workflow](#task-workflow)
    - [Pop Command](#pop-command)
    - [Shift Command](#shift-command)
    - [Next Command](#next-command)
  - [Markdown Parsing](#markdown-parsing)
    - [Creating Tasks from Markdown Files](#creating-tasks-from-markdown-files)
    - [Creating Tasks from Standard Input](#creating-tasks-from-standard-input)
    - [Markdown Format](#markdown-format)

## Initialization

Before using Tamo, you need to initialize it in your project directory:

```bash
tamo init
```

This creates a `.tamo` directory in the current directory with an empty `data.json` file. All your tasks and memos will be stored in this file.

If Tamo is already initialized in the directory, you'll see a message indicating that.

## Task Management

### Adding Tasks

There are several ways to add tasks in Tamo:

#### Standard Task Addition (at the end)

```bash
tamo add task "Complete documentation" -d "Write comprehensive docs for the project"
```

This adds a task at the end of the task list. The `-d` flag allows you to specify a description.

#### Adding a Task with Memo References

```bash
tamo add task "Implement feature X" -d "Follow the spec" -m memo1_id,memo2_id
```

The `-m` flag allows you to specify comma-separated memo IDs that this task references.

#### Adding a Task at the End (Push)

```bash
tamo push task "Low priority task" -d "This can be done later"
```

This is equivalent to `add task` and adds the task at the end of the list.

#### Adding a Task at the Beginning (Unshift)

```bash
tamo unshift task "High priority task" -d "This needs to be done first"
```

This adds the task at the beginning of the list, making it the first task to be done.

### Listing Tasks

#### List All Tasks

```bash
tamo list tasks
```

This lists all tasks ordered by their `order` value.

#### List Only Completed Tasks

```bash
tamo list tasks --done
```

#### List Only Uncompleted Tasks

```bash
tamo list tasks --undone
```

#### List Tasks Referencing a Specific Memo

```bash
tamo list tasks --refs <memo_id>
```

### Showing Task Details

To see the full details of a specific task:

```bash
tamo show <task_id>
```

You can use either the full UUID or a prefix of the ID. The output includes:
- Task ID
- Title
- Order
- Status (completed or not)
- Creation and update timestamps
- Description (if any)
- Referenced memos (if any)

### Editing Tasks

To edit a task:

```bash
tamo edit <task_id>
```

This opens a simple prompt-based editor where you can update:
- Title
- Description
- Memo references

If you prefer to use your system's text editor:

```bash
tamo edit <task_id> --editor
```

This opens your default editor (specified by the `EDITOR` environment variable) with the task content.

### Marking Tasks as Done/Undone

To mark a task as completed:

```bash
tamo done <task_id>
```

To mark a completed task as not done:

```bash
tamo undone <task_id>
```

### Moving Tasks

#### Moving a Task to a Specific Order

```bash
tamo mv <task_id> <order_number>
```

This sets the task's order to the specified number.

#### Moving a Task Relative to Another Task

```bash
# Move a task before another task
tamo mv <task_id> before <other_task_id>

# Move a task after another task
tamo mv <task_id> after <other_task_id>
```

### Removing Tasks

To remove a task:

```bash
tamo rm <task_id>
```

To force removal without confirmation:

```bash
tamo rm <task_id> -f
```

## Memo Management

### Adding Memos

#### Adding a Memo with Content

```bash
tamo add memo "Important Information" -c "This is important content to remember"
```

The title is optional. If you don't provide a title, the memo will be created without one.

#### Adding a Memo from Standard Input

```bash
tamo add memo --from-stdin
```

This allows you to type the memo content directly. Press Ctrl+D when finished.

### Listing Memos

To list all memos:

```bash
tamo list memos
```

To list all tasks and memos:

```bash
tamo list all
```

### Showing Memo Details

To see the full details of a specific memo:

```bash
tamo show <memo_id>
```

You can use either the full UUID or a prefix of the ID. The output includes:
- Memo ID
- Title (if any)
- Creation and update timestamps
- Content

### Editing Memos

To edit a memo:

```bash
tamo edit <memo_id>
```

This opens a simple prompt-based editor where you can update:
- Title
- Content

If you prefer to use your system's text editor:

```bash
tamo edit <memo_id> --editor
```

This opens your default editor (specified by the `EDITOR` environment variable) with the memo content.

### Removing Memos

To remove a memo:

```bash
tamo rm <memo_id>
```

If the memo is referenced by any tasks, you'll see a warning. To force removal:

```bash
tamo rm <memo_id> -f
```

## Task Workflow

Tamo provides special commands to help with task workflow management.

### Pop Command

The `pop` command operates on the last task in the list (highest order value):

```bash
# Show the last task
tamo pop task

# Mark the last task as done
tamo pop task --done

# Remove the last task
tamo pop task --rm
```

### Shift Command

The `shift` command operates on the first task in the list (lowest order value):

```bash
# Show the first task
tamo shift task

# Mark the first task as done
tamo shift task --done

# Remove the first task
tamo shift task --rm
```

### Next Command

The `next` command shows the first undone task:

```bash
tamo next
```

This is useful for quickly seeing what task you should work on next.

## Markdown Parsing

Tamo can create tasks and associated memos from Markdown files.

### Creating Tasks from Markdown Files

```bash
tamo add task -f task_description.md
```

### Creating Tasks from Standard Input

```bash
cat task_description.md | tamo add task --from-stdin
```

### Markdown Format

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

The parser:
1. Uses the first H1 heading as the task title
2. Creates a separate memo for each ```memo block
3. Replaces the memo blocks in the original text with references to the created memos
4. Uses the modified text as the task description

## Task Flattening

### Flattening Tasks for AI Prompts

The `flattask` command expands all memo references in a task, creating a single comprehensive document:

```bash
tamo flattask <task_id>
```

This outputs a Markdown document that includes:
- The task title as an H1 heading
- The task status
- The task description
- The content of all referenced memos, properly formatted

This is particularly useful for:
- Creating comprehensive prompts for AI tools
- Getting a complete view of a task with all its associated information
- Sharing task information with others
