package cli

import (
        "bytes"
        "io"
        "os"
        "strings"
        "testing"
)

// Helper function to capture stdout for testing
func captureOutput(f func() error) (string, error) {
        old := os.Stdout
        r, w, _ := os.Pipe()
        os.Stdout = w

        err := f()

        w.Close()
        os.Stdout = old

        var buf bytes.Buffer
        io.Copy(&buf, r)
        return buf.String(), err
}

// TestExecuteHelp tests the help command
func TestExecuteHelp(t *testing.T) {
        cli := NewCLI()

        output, err := captureOutput(func() error {
                return cli.executeHelp([]string{})
        })

        if err != nil {
                t.Errorf("Unexpected error: %v", err)
        }

        // Check that the output contains expected help text
        if !strings.Contains(output, "tamo - Task and Memo Management CLI") {
                t.Errorf("Expected help output to contain app description, got: %s", output)
        }

        if !strings.Contains(output, "Available commands:") {
                t.Errorf("Expected help output to list available commands, got: %s", output)
        }
}

// TestExecuteInit tests the init command
func TestExecuteInit(t *testing.T) {
        // Create a temporary directory for testing
        tempDir, err := os.MkdirTemp("", "tamo-test")
        if err != nil {
                t.Fatalf("Failed to create temp dir: %v", err)
        }
        defer os.RemoveAll(tempDir)

        // Change to the temporary directory
        oldWd, err := os.Getwd()
        if err != nil {
                t.Fatalf("Failed to get current directory: %v", err)
        }
        if err := os.Chdir(tempDir); err != nil {
                t.Fatalf("Failed to change to temp dir: %v", err)
        }
        defer os.Chdir(oldWd)

        cli := NewCLI()

        // Test initialization
        output, err := captureOutput(func() error {
                return cli.executeInit([]string{})
        })

        if err != nil {
                t.Errorf("Unexpected error: %v", err)
        }

        if !strings.Contains(output, "tamo initialized successfully") {
                t.Errorf("Expected output to contain initialization success message, got: %s", output)
        }

        // Check that the .tamo directory was created
        tamoDir := ".tamo"
        if _, err := os.Stat(tamoDir); os.IsNotExist(err) {
                t.Errorf("Expected .tamo directory to exist, but it doesn't")
        }

        // Check that the data.json file was created
        dataFile := ".tamo/data.json"
        if _, err := os.Stat(dataFile); os.IsNotExist(err) {
                t.Errorf("Expected data.json file to exist, but it doesn't")
        }

        // Test initialization when already initialized
        output, err = captureOutput(func() error {
                return cli.executeInit([]string{})
        })

        if err != nil {
                t.Errorf("Unexpected error: %v", err)
        }

        if !strings.Contains(output, "already initialized") {
                t.Errorf("Expected output to contain 'already initialized', got: %s", output)
        }
}

// TestExecuteAdd tests the add command
func TestExecuteAdd(t *testing.T) {
        cli := NewCLI()

        // Test missing subcommand
        _, err := captureOutput(func() error {
                return cli.executeAdd([]string{})
        })

        if err == nil || !strings.Contains(err.Error(), "missing subcommand") {
                t.Errorf("Expected error about missing subcommand, got: %v", err)
        }

        // Test unknown subcommand
        _, err = captureOutput(func() error {
                return cli.executeAdd([]string{"unknown"})
        })

        if err == nil || !strings.Contains(err.Error(), "unknown subcommand") {
                t.Errorf("Expected error about unknown subcommand, got: %v", err)
        }
}

// TestExecuteAddTask tests the add task command
func TestExecuteAddTask(t *testing.T) {
        // Create a temporary directory for testing
        tempDir, err := os.MkdirTemp("", "tamo-test")
        if err != nil {
                t.Fatalf("Failed to create temp dir: %v", err)
        }
        defer os.RemoveAll(tempDir)

        // Change to the temporary directory
        oldWd, err := os.Getwd()
        if err != nil {
                t.Fatalf("Failed to get current directory: %v", err)
        }
        if err := os.Chdir(tempDir); err != nil {
                t.Fatalf("Failed to change to temp dir: %v", err)
        }
        defer os.Chdir(oldWd)

        // Initialize tamo
        cli := NewCLI()
        if err := cli.executeInit([]string{}); err != nil {
                t.Fatalf("Failed to initialize tamo: %v", err)
        }

        // Test adding a task
        output, err := captureOutput(func() error {
                return cli.executeAddTask([]string{"Test Task", "-d", "Test Description"}, "add")
        })

        if err != nil {
                t.Errorf("Unexpected error: %v", err)
        }

        if !strings.Contains(output, "Task added with ID") {
                t.Errorf("Expected output to contain task added message, got: %s", output)
        }
}

// TestExecuteList tests the list command
func TestExecuteList(t *testing.T) {
        // Create a temporary directory for testing
        tempDir, err := os.MkdirTemp("", "tamo-test")
        if err != nil {
                t.Fatalf("Failed to create temp dir: %v", err)
        }
        defer os.RemoveAll(tempDir)

        // Change to the temporary directory
        oldWd, err := os.Getwd()
        if err != nil {
                t.Fatalf("Failed to get current directory: %v", err)
        }
        if err := os.Chdir(tempDir); err != nil {
                t.Fatalf("Failed to change to temp dir: %v", err)
        }
        defer os.Chdir(oldWd)

        // Initialize tamo
        cli := NewCLI()
        if err := cli.executeInit([]string{}); err != nil {
                t.Fatalf("Failed to initialize tamo: %v", err)
        }

        // Add a task
        if err := cli.executeAddTask([]string{"Test Task", "-d", "Test Description"}, "add"); err != nil {
                t.Fatalf("Failed to add task: %v", err)
        }

        // Test listing tasks
        output, err := captureOutput(func() error {
                return cli.executeList([]string{"tasks"})
        })

        if err != nil {
                t.Errorf("Unexpected error: %v", err)
        }

        if !strings.Contains(output, "Test Task") {
                t.Errorf("Expected output to contain task title, got: %s", output)
        }
}

// TestExecuteDone tests the done command
func TestExecuteDone(t *testing.T) {
        // Create a temporary directory for testing
        tempDir, err := os.MkdirTemp("", "tamo-test")
        if err != nil {
                t.Fatalf("Failed to create temp dir: %v", err)
        }
        defer os.RemoveAll(tempDir)

        // Change to the temporary directory
        oldWd, err := os.Getwd()
        if err != nil {
                t.Fatalf("Failed to get current directory: %v", err)
        }
        if err := os.Chdir(tempDir); err != nil {
                t.Fatalf("Failed to change to temp dir: %v", err)
        }
        defer os.Chdir(oldWd)

        // Initialize tamo
        cli := NewCLI()
        if err := cli.executeInit([]string{}); err != nil {
                t.Fatalf("Failed to initialize tamo: %v", err)
        }

        // Add a task
        var taskID string
        output, err := captureOutput(func() error {
                return cli.executeAddTask([]string{"Test Task", "-d", "Test Description"}, "add")
        })
        if err != nil {
                t.Fatalf("Failed to add task: %v", err)
        }

        // Extract task ID from output
        idStart := strings.Index(output, "Task added with ID: ") + len("Task added with ID: ")
        if idStart >= len("Task added with ID: ") {
                taskID = strings.TrimSpace(output[idStart:])
        } else {
                t.Fatalf("Failed to extract task ID from output: %s", output)
        }

        // Test marking task as done
        output, err = captureOutput(func() error {
                return cli.executeDone([]string{taskID})
        })

        if err != nil {
                t.Errorf("Unexpected error: %v", err)
        }

        if !strings.Contains(output, "marked as done") {
                t.Errorf("Expected output to contain 'marked as done', got: %s", output)
        }

        // Test marking non-existent task as done
        _, err = captureOutput(func() error {
                return cli.executeDone([]string{"nonexistent"})
        })

        if err == nil || !strings.Contains(err.Error(), "no task found") {
                t.Errorf("Expected error about task not found, got: %v", err)
        }
}

// TestExecuteUndone tests the undone command
func TestExecuteUndone(t *testing.T) {
        // Create a temporary directory for testing
        tempDir, err := os.MkdirTemp("", "tamo-test")
        if err != nil {
                t.Fatalf("Failed to create temp dir: %v", err)
        }
        defer os.RemoveAll(tempDir)

        // Change to the temporary directory
        oldWd, err := os.Getwd()
        if err != nil {
                t.Fatalf("Failed to get current directory: %v", err)
        }
        if err := os.Chdir(tempDir); err != nil {
                t.Fatalf("Failed to change to temp dir: %v", err)
        }
        defer os.Chdir(oldWd)

        // Initialize tamo
        cli := NewCLI()
        if err := cli.executeInit([]string{}); err != nil {
                t.Fatalf("Failed to initialize tamo: %v", err)
        }

        // Add a task
        var taskID string
        output, err := captureOutput(func() error {
                return cli.executeAddTask([]string{"Test Task", "-d", "Test Description"}, "add")
        })
        if err != nil {
                t.Fatalf("Failed to add task: %v", err)
        }

        // Extract task ID from output
        idStart := strings.Index(output, "Task added with ID: ") + len("Task added with ID: ")
        if idStart >= len("Task added with ID: ") {
                taskID = strings.TrimSpace(output[idStart:])
        } else {
                t.Fatalf("Failed to extract task ID from output: %s", output)
        }

        // Mark task as done
        if err := cli.executeDone([]string{taskID}); err != nil {
                t.Fatalf("Failed to mark task as done: %v", err)
        }

        // Test marking task as undone
        output, err = captureOutput(func() error {
                return cli.executeUndone([]string{taskID})
        })

        if err != nil {
                t.Errorf("Unexpected error: %v", err)
        }

        if !strings.Contains(output, "marked as not done") {
                t.Errorf("Expected output to contain 'marked as not done', got: %s", output)
        }
}

// TestExecuteMove tests the move command
func TestExecuteMove(t *testing.T) {
        // Create a temporary directory for testing
        tempDir, err := os.MkdirTemp("", "tamo-test")
        if err != nil {
                t.Fatalf("Failed to create temp dir: %v", err)
        }
        defer os.RemoveAll(tempDir)

        // Change to the temporary directory
        oldWd, err := os.Getwd()
        if err != nil {
                t.Fatalf("Failed to get current directory: %v", err)
        }
        if err := os.Chdir(tempDir); err != nil {
                t.Fatalf("Failed to change to temp dir: %v", err)
        }
        defer os.Chdir(oldWd)

        // Initialize tamo
        cli := NewCLI()
        if err := cli.executeInit([]string{}); err != nil {
                t.Fatalf("Failed to initialize tamo: %v", err)
        }

        // Add two tasks
        var taskID1, taskID2 string
        output, err := captureOutput(func() error {
                return cli.executeAddTask([]string{"Task 1", "-d", "Description 1"}, "add")
        })
        if err != nil {
                t.Fatalf("Failed to add task 1: %v", err)
        }

        // Extract task ID from output
        idStart := strings.Index(output, "Task added with ID: ") + len("Task added with ID: ")
        if idStart >= len("Task added with ID: ") {
                taskID1 = strings.TrimSpace(output[idStart:])
        } else {
                t.Fatalf("Failed to extract task ID from output: %s", output)
        }

        output, err = captureOutput(func() error {
                return cli.executeAddTask([]string{"Task 2", "-d", "Description 2"}, "add")
        })
        if err != nil {
                t.Fatalf("Failed to add task 2: %v", err)
        }

        // Extract task ID from output
        idStart = strings.Index(output, "Task added with ID: ") + len("Task added with ID: ")
        if idStart >= len("Task added with ID: ") {
                taskID2 = strings.TrimSpace(output[idStart:])
        } else {
                t.Fatalf("Failed to extract task ID from output: %s", output)
        }

        // Test moving task to absolute order
        output, err = captureOutput(func() error {
                return cli.executeMove([]string{taskID1, "5.0"})
        })

        if err != nil {
                t.Errorf("Unexpected error: %v", err)
        }

        if !strings.Contains(output, "moved to order 5.0") {
                t.Errorf("Expected output to contain 'moved to order 5.0', got: %s", output)
        }

        // Test moving task relative to another task
        output, err = captureOutput(func() error {
                return cli.executeMove([]string{taskID1, "after", taskID2})
        })

        if err != nil {
                t.Errorf("Unexpected error: %v", err)
        }

        if !strings.Contains(output, "moved after") {
                t.Errorf("Expected output to contain 'moved after', got: %s", output)
        }
}

// TestExecuteFlattask tests the flattask command
func TestExecuteFlattask(t *testing.T) {
        // Create a temporary directory for testing
        tempDir, err := os.MkdirTemp("", "tamo-test")
        if err != nil {
                t.Fatalf("Failed to create temp dir: %v", err)
        }
        defer os.RemoveAll(tempDir)

        // Change to the temporary directory
        oldWd, err := os.Getwd()
        if err != nil {
                t.Fatalf("Failed to get current directory: %v", err)
        }
        if err := os.Chdir(tempDir); err != nil {
                t.Fatalf("Failed to change to temp dir: %v", err)
        }
        defer os.Chdir(oldWd)

        // Initialize tamo
        cli := NewCLI()
        if err := cli.executeInit([]string{}); err != nil {
                t.Fatalf("Failed to initialize tamo: %v", err)
        }

        // Add a memo
        var memoID string
        output, err := captureOutput(func() error {
                return cli.executeAddMemo([]string{"Test Memo", "-c", "Test Memo Content"})
        })
        if err != nil {
                t.Fatalf("Failed to add memo: %v", err)
        }

        // Extract memo ID from output
        idStart := strings.Index(output, "Memo added with ID: ") + len("Memo added with ID: ")
        if idStart >= len("Memo added with ID: ") {
                memoID = strings.TrimSpace(output[idStart:])
        } else {
                t.Fatalf("Failed to extract memo ID from output: %s", output)
        }

        // Add a task with memo reference
        var taskID string
        output, err = captureOutput(func() error {
                return cli.executeAddTask([]string{"Test Task", "-d", "Test Description", "-m", memoID}, "add")
        })
        if err != nil {
                t.Fatalf("Failed to add task: %v", err)
        }

        // Extract task ID from output
        idStart = strings.Index(output, "Task added with ID: ") + len("Task added with ID: ")
        if idStart >= len("Task added with ID: ") {
                taskID = strings.TrimSpace(output[idStart:])
        } else {
                t.Fatalf("Failed to extract task ID from output: %s", output)
        }

        // Test flattask command
        output, err = captureOutput(func() error {
                return cli.executeFlattask([]string{taskID})
        })

        if err != nil {
                t.Errorf("Unexpected error: %v", err)
        }

        if !strings.Contains(output, "Test Task") || !strings.Contains(output, "Test Description") {
                t.Errorf("Expected output to contain task title and description, got: %s", output)
        }

        if !strings.Contains(output, "Test Memo") {
                t.Errorf("Expected output to contain memo title, got: %s", output)
        }
}

// TestExecuteShowMemo tests the show command for memos, including the display of referencing tasks
func TestExecuteShowMemo(t *testing.T) {
        // Create a temporary directory for testing
        tempDir, err := os.MkdirTemp("", "tamo-test")
        if err != nil {
                t.Fatalf("Failed to create temp dir: %v", err)
        }
        defer os.RemoveAll(tempDir)

        // Change to the temporary directory
        oldWd, err := os.Getwd()
        if err != nil {
                t.Fatalf("Failed to get current directory: %v", err)
        }
        if err := os.Chdir(tempDir); err != nil {
                t.Fatalf("Failed to change to temp dir: %v", err)
        }
        defer os.Chdir(oldWd)

        // Initialize tamo
        cli := NewCLI()
        if err := cli.executeInit([]string{}); err != nil {
                t.Fatalf("Failed to initialize tamo: %v", err)
        }

        // Add a memo
        var memoID string
        output, err := captureOutput(func() error {
                return cli.executeAddMemo([]string{"-c", "This is a test memo content"})
        })
        if err != nil {
                t.Fatalf("Failed to add memo: %v", err)
        }

        // Extract memo ID from output
        idStart := strings.Index(output, "Memo added with ID: ") + len("Memo added with ID: ")
        if idStart >= len("Memo added with ID: ") {
                memoID = strings.TrimSpace(output[idStart:])
        } else {
                t.Fatalf("Failed to extract memo ID from output: %s", output)
        }

        // Add a task that references the memo
        output, err = captureOutput(func() error {
                return cli.executeAddTask([]string{"Task referencing memo", "-d", "Task description", "-m", memoID}, "add")
        })
        if err != nil {
                t.Fatalf("Failed to add task: %v", err)
        }

        // Add another task that references the memo
        output, err = captureOutput(func() error {
                return cli.executeAddTask([]string{"Another task referencing memo", "-d", "Another description", "-m", memoID}, "add")
        })
        if err != nil {
                t.Fatalf("Failed to add second task: %v", err)
        }

        // Test showing the memo
        output, err = captureOutput(func() error {
                return cli.executeShow([]string{memoID[:8]}) // Use prefix of memo ID
        })

        if err != nil {
                t.Errorf("Unexpected error: %v", err)
        }

        // Check that the output contains the memo content
        if !strings.Contains(output, "This is a test memo content") {
                t.Errorf("Expected output to contain memo content, got: %s", output)
        }

        // Check that the output contains the "Referenced by Tasks" section
        if !strings.Contains(output, "Referenced by Tasks:") {
                t.Errorf("Expected output to contain 'Referenced by Tasks' section, got: %s", output)
        }

        // Check that the output contains both referencing tasks
        if !strings.Contains(output, "Task referencing memo") {
                t.Errorf("Expected output to contain first referencing task, got: %s", output)
        }

        if !strings.Contains(output, "Another task referencing memo") {
                t.Errorf("Expected output to contain second referencing task, got: %s", output)
        }

        // Check that the "Referenced by Tasks" section appears twice (before and after content)
        firstOccurrence := strings.Index(output, "Referenced by Tasks:")
        if firstOccurrence == -1 {
                t.Errorf("Expected output to contain 'Referenced by Tasks' section, got: %s", output)
        }

        secondOccurrence := strings.Index(output[firstOccurrence+1:], "Referenced by Tasks:")
        if secondOccurrence == -1 {
                t.Errorf("Expected output to contain 'Referenced by Tasks' section twice, got: %s", output)
        }
}

// TestExecuteShowTask tests the show command for tasks
func TestExecuteShowTask(t *testing.T) {
        // Create a temporary directory for testing
        tempDir, err := os.MkdirTemp("", "tamo-test")
        if err != nil {
                t.Fatalf("Failed to create temp dir: %v", err)
        }
        defer os.RemoveAll(tempDir)

        // Change to the temporary directory
        oldWd, err := os.Getwd()
        if err != nil {
                t.Fatalf("Failed to get current directory: %v", err)
        }
        if err := os.Chdir(tempDir); err != nil {
                t.Fatalf("Failed to change to temp dir: %v", err)
        }
        defer os.Chdir(oldWd)

        // Initialize tamo
        cli := NewCLI()
        if err := cli.executeInit([]string{}); err != nil {
                t.Fatalf("Failed to initialize tamo: %v", err)
        }

        // Add a task
        var taskID string
        output, err := captureOutput(func() error {
                return cli.executeAddTask([]string{"Test Task", "-d", "Test Description"}, "add")
        })
        if err != nil {
                t.Fatalf("Failed to add task: %v", err)
        }

        // Extract task ID from output
        idStart := strings.Index(output, "Task added with ID: ") + len("Task added with ID: ")
        if idStart >= len("Task added with ID: ") {
                taskID = strings.TrimSpace(output[idStart:])
        } else {
                t.Fatalf("Failed to extract task ID from output: %s", output)
        }

        // Test showing the task
        output, err = captureOutput(func() error {
                return cli.executeShow([]string{taskID[:8]}) // Use prefix of task ID
        })

        if err != nil {
                t.Errorf("Unexpected error: %v", err)
        }

        // Check that the output contains the task details
        if !strings.Contains(output, "Task ID:") {
                t.Errorf("Expected output to contain 'Task ID:', got: %s", output)
        }

        if !strings.Contains(output, "Test Task") {
                t.Errorf("Expected output to contain task title, got: %s", output)
        }

        if !strings.Contains(output, "Test Description") {
                t.Errorf("Expected output to contain task description, got: %s", output)
        }
}
