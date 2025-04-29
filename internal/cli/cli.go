package cli

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/zishida/tamo/internal/model"
	"github.com/zishida/tamo/internal/storage"
	"github.com/zishida/tamo/internal/utils"
)

// Command represents a CLI command
type Command struct {
	Name        string
	Description string
	Execute     func(args []string) error
}

// CLI represents the command-line interface
type CLI struct {
	commands map[string]Command
}

// NewCLI creates a new CLI
func NewCLI() *CLI {
	cli := &CLI{
		commands: make(map[string]Command),
	}

	// Register commands
	cli.registerCommands()

	return cli
}

// registerCommands registers all available commands
func (c *CLI) registerCommands() {
	// Register init command
	c.commands["init"] = Command{
		Name:        "init",
		Description: "Initialize tamo in the current directory",
		Execute:     c.executeInit,
	}

	// Register help command
	c.commands["help"] = Command{
		Name:        "help",
		Description: "Show help information",
		Execute:     c.executeHelp,
	}

	// Register add commands
	c.commands["add"] = Command{
		Name:        "add",
		Description: "Add a new task or memo",
		Execute:     c.executeAdd,
	}

	// Register push command (alias for add task with order at end)
	c.commands["push"] = Command{
		Name:        "push",
		Description: "Add a new task at the end of the list",
		Execute:     c.executePush,
	}

	// Register unshift command (alias for add task with order at beginning)
	c.commands["unshift"] = Command{
		Name:        "unshift",
		Description: "Add a new task at the beginning of the list",
		Execute:     c.executeUnshift,
	}

	// Register list command
	c.commands["list"] = Command{
		Name:        "list",
		Description: "List tasks and/or memos",
		Execute:     c.executeList,
	}

	// Register show command
	c.commands["show"] = Command{
		Name:        "show",
		Description: "Show details of a task or memo",
		Execute:     c.executeShow,
	}

	// Register remove command
	c.commands["rm"] = Command{
		Name:        "rm",
		Description: "Remove a task or memo",
		Execute:     c.executeRemove,
	}

	// Register edit command
	c.commands["edit"] = Command{
		Name:        "edit",
		Description: "Edit a task or memo",
		Execute:     c.executeEdit,
	}

	// Register done command
	c.commands["done"] = Command{
		Name:        "done",
		Description: "Mark a task as done",
		Execute:     c.executeDone,
	}

	// Register undone command
	c.commands["undone"] = Command{
		Name:        "undone",
		Description: "Mark a task as not done",
		Execute:     c.executeUndone,
	}

	// Register move command
	c.commands["mv"] = Command{
		Name:        "mv",
		Description: "Move a task to a specific order or relative to another task",
		Execute:     c.executeMove,
	}

	// Register pop command
	c.commands["pop"] = Command{
		Name:        "pop",
		Description: "Show, mark as done, or remove the last task",
		Execute:     c.executePop,
	}

	// Register shift command
	c.commands["shift"] = Command{
		Name:        "shift",
		Description: "Show, mark as done, or remove the first task",
		Execute:     c.executeShift,
	}

	// Register next command (alias for shift task)
	c.commands["next"] = Command{
		Name:        "next",
		Description: "Show the first undone task",
		Execute:     c.executeNext,
	}

	// Register flattask command
	c.commands["flattask"] = Command{
		Name:        "flattask",
		Description: "Flatten a task by expanding all memo references",
		Execute:     c.executeFlattask,
	}
}

// Execute executes the CLI with the given arguments
func Execute() error {
	cli := NewCLI()

	// If no arguments, show help
	if len(os.Args) < 2 {
		return cli.executeHelp([]string{})
	}

	// Get command name
	cmdName := os.Args[1]

	// Find command
	cmd, ok := cli.commands[cmdName]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmdName)
		return cli.executeHelp([]string{})
	}

	// Execute command
	return cmd.Execute(os.Args[2:])
}

// executeInit initializes tamo in the current directory
func (c *CLI) executeInit(args []string) error {
	// Parse flags
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	initCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: tamo init\n\n")
		fmt.Fprintf(os.Stderr, "Initialize tamo in the current directory\n\n")
		initCmd.PrintDefaults()
	}

	if err := initCmd.Parse(args); err != nil {
		return err
	}

	// Create storage
	s := storage.NewStorage()

	// Check if already initialized
	if s.Exists() {
		fmt.Println("tamo is already initialized in this directory")
		return nil
	}

	// Initialize storage
	if err := s.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize tamo: %w", err)
	}

	fmt.Println("tamo initialized successfully")
	return nil
}

// executeHelp shows help information
func (c *CLI) executeHelp(args []string) error {
	fmt.Println("tamo - Task and Memo Management CLI")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  tamo <command> [arguments]")
	fmt.Println()
	fmt.Println("Available commands:")

	// Get max command name length for alignment
	maxLen := 0
	for _, cmd := range c.commands {
		if len(cmd.Name) > maxLen {
			maxLen = len(cmd.Name)
		}
	}

	// Print commands
	for _, cmd := range c.commands {
		fmt.Printf("  %-*s  %s\n", maxLen, cmd.Name, cmd.Description)
	}

	return nil
}

// executeAdd handles the 'add' command for both tasks and memos
func (c *CLI) executeAdd(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing subcommand: 'task' or 'memo'")
	}

	subCmd := args[0]
	switch subCmd {
	case "memo":
		return c.executeAddMemo(args[1:])
	case "task":
		return c.executeAddTask(args[1:], "add")
	default:
		return fmt.Errorf("unknown subcommand: %s", subCmd)
	}
}

// executeAddMemo handles the 'add memo' command
func (c *CLI) executeAddMemo(args []string) error {
	// Create flag set
	memoCmd := flag.NewFlagSet("add memo", flag.ExitOnError)

	// Define flags
	contentFlag := memoCmd.String("c", "", "Memo content")
	fromStdinFlag := memoCmd.Bool("from-stdin", false, "Read content from stdin")
	editorFlag := memoCmd.Bool("editor", false, "Open editor to input content")

	// Set usage
	memoCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: tamo add memo [<title>] [-c \"<content>\" | --from-stdin | --editor]\n\n")
		fmt.Fprintf(os.Stderr, "Add a new memo\n\n")
		memoCmd.PrintDefaults()
	}

	// Parse flags
	if err := memoCmd.Parse(args); err != nil {
		return err
	}

	// Get title (optional)
	var title *string
	if memoCmd.NArg() > 0 {
		t := memoCmd.Arg(0)
		title = &t
	}

	// Get content based on flags
	var content string

	// Check if multiple content sources are specified
	contentSources := 0
	if *contentFlag != "" {
		contentSources++
	}
	if *fromStdinFlag {
		contentSources++
	}
	if *editorFlag {
		contentSources++
	}

	if contentSources > 1 {
		return fmt.Errorf("only one of -c, --from-stdin, or --editor can be specified")
	}

	// Get content from the specified source
	if *contentFlag != "" {
		content = *contentFlag
	} else if *fromStdinFlag {
		// Read from stdin
		scanner := bufio.NewScanner(os.Stdin)
		var contentBuilder strings.Builder
		for scanner.Scan() {
			contentBuilder.WriteString(scanner.Text())
			contentBuilder.WriteString("\n")
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading from stdin: %w", err)
		}
		content = contentBuilder.String()
	} else if *editorFlag {
		// TODO: Implement editor support
		return fmt.Errorf("editor support not implemented yet")
	} else {
		// Default to simple input if no flag is specified
		// For now, we'll just use a simple prompt
		fmt.Println("Enter memo content (press Ctrl+D when finished):")
		scanner := bufio.NewScanner(os.Stdin)
		var contentBuilder strings.Builder
		for scanner.Scan() {
			contentBuilder.WriteString(scanner.Text())
			contentBuilder.WriteString("\n")
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading content: %w", err)
		}
		content = contentBuilder.String()
	}

	// Generate UUID
	id, err := utils.GenerateUUID()
	if err != nil {
		return fmt.Errorf("failed to generate UUID: %w", err)
	}

	// Create new memo
	memo := model.NewMemo(id, title, content)

	// Load store
	s := storage.NewStorage()
	store, err := s.Load()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	// Add memo to store
	store.AddMemo(memo)

	// Save store
	if err := s.Save(store); err != nil {
		return fmt.Errorf("failed to save data: %w", err)
	}

	fmt.Printf("Memo added with ID: %s\n", id)
	return nil
}

// executeAddTask handles the 'add task' command
func (c *CLI) executeAddTask(args []string, mode string) error {
	// Check for Markdown parsing options
	if len(args) > 0 && (args[0] == "-f" || args[0] == "--from-stdin") {
		return c.executeAddTaskFromMarkdown(args)
	}

	// Manual argument parsing
	// Set usage
	usage := func() {
		fmt.Fprintf(os.Stderr, "Usage: tamo %s task \"<title>\" [-d \"<description>\"] [-m <memo_id>,...]\n", mode)
		fmt.Fprintf(os.Stderr, "       tamo %s task -f <filepath> | --from-stdin\n\n", mode)
		fmt.Fprintf(os.Stderr, "Add a new task\n\n")
		fmt.Fprintf(os.Stderr, "  -d <description>    Task description\n")
		fmt.Fprintf(os.Stderr, "  -m <memo_id>,...    Comma-separated list of memo IDs\n")
		fmt.Fprintf(os.Stderr, "  -f <filepath>       Create task from Markdown file\n")
		fmt.Fprintf(os.Stderr, "  --from-stdin        Create task from Markdown input on stdin\n")
	}

	// Check if we have at least a title
	if len(args) < 1 {
		usage()
		return fmt.Errorf("missing task title")
	}

	// Get title
	title := args[0]

	// Parse remaining arguments for flags
	var description string
	var memoRefsStr string

	for i := 1; i < len(args); i++ {
		if args[i] == "-d" && i+1 < len(args) {
			description = args[i+1]
			i++ // Skip the next argument
		} else if args[i] == "-m" && i+1 < len(args) {
			memoRefsStr = args[i+1]
			i++ // Skip the next argument
		}
	}

	// Parse memo refs
	var memoRefs []string
	if memoRefsStr != "" {
		inputRefs := strings.Split(memoRefsStr, ",")
		// Trim whitespace from each memo ID
		for _, ref := range inputRefs {
			memoRefs = append(memoRefs, strings.TrimSpace(ref))
		}
	}

	// Load store
	s := storage.NewStorage()
	store, err := s.Load()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	// Convert partial memo IDs to full IDs
	for i, refID := range memoRefs {
		// Find the full memo ID if a partial ID is provided
		if len(refID) < 36 {
			found := false
			for _, memo := range store.Memos {
				if strings.HasPrefix(memo.ID, refID) {
					memoRefs[i] = memo.ID
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("memo with ID %s not found", refID)
			}
		}
	}

	// Generate UUID
	id, err := utils.GenerateUUID()
	if err != nil {
		return fmt.Errorf("failed to generate UUID: %w", err)
	}

	// Validate memo refs (for full IDs)
	for _, memoID := range memoRefs {
		if len(memoID) == 36 && store.FindMemoByID(memoID) == nil {
			return fmt.Errorf("memo with ID %s not found", memoID)
		}
	}

	// Create new task
	task := model.NewTask(id, title, description, memoRefs)

	// Set order based on mode
	switch mode {
	case "add", "push":
		// Add to end (max order + 1.0)
		task.Order = store.GetMaxTaskOrder() + 1.0
	case "unshift":
		// Add to beginning (min order - 1.0)
		task.Order = store.GetMinTaskOrder() - 1.0
	}

	// Add task to store
	store.AddTask(task)

	// Save store
	if err := s.Save(store); err != nil {
		return fmt.Errorf("failed to save data: %w", err)
	}

	fmt.Printf("Task added with ID: %s\n", id)
	return nil
}

// executePush handles the 'push task' command (add to end)
func (c *CLI) executePush(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing subcommand: 'task'")
	}

	subCmd := args[0]
	if subCmd != "task" {
		return fmt.Errorf("unknown subcommand: %s", subCmd)
	}

	return c.executeAddTask(args[1:], "push")
}

// executeUnshift handles the 'unshift task' command (add to beginning)
func (c *CLI) executeUnshift(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing subcommand: 'task'")
	}

	subCmd := args[0]
	if subCmd != "task" {
		return fmt.Errorf("unknown subcommand: %s", subCmd)
	}

	return c.executeAddTask(args[1:], "unshift")
}

// executeList handles the 'list' command
func (c *CLI) executeList(args []string) error {
	// Create flag set
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	// Define flags
	doneFlag := listCmd.Bool("done", false, "Show only completed tasks")
	undoneFlag := listCmd.Bool("undone", false, "Show only uncompleted tasks")
	refsFlag := listCmd.String("refs", "", "Show tasks referencing the specified memo ID")

	// Set usage
	listCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: tamo list [tasks|memos|all] [--done|--undone] [--refs <memo_id>]\n\n")
		fmt.Fprintf(os.Stderr, "List tasks and/or memos\n\n")
		listCmd.PrintDefaults()
	}

	// Parse flags
	if err := listCmd.Parse(args); err != nil {
		return err
	}

	// Get subcommand (default to "tasks")
	subCmd := "tasks"
	if listCmd.NArg() > 0 {
		subCmd = listCmd.Arg(0)
		if subCmd != "tasks" && subCmd != "memos" && subCmd != "all" {
			return fmt.Errorf("unknown subcommand: %s", subCmd)
		}
	}

	// Check for conflicting flags
	if *doneFlag && *undoneFlag {
		return fmt.Errorf("--done and --undone flags cannot be used together")
	}

	// Load store
	s := storage.NewStorage()
	store, err := s.Load()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	// List items based on subcommand
	switch subCmd {
	case "tasks", "all":
		// Filter tasks
		var filteredTasks []*model.Task
		for _, task := range store.Tasks {
			// Filter by done/undone
			if *doneFlag && !task.Done {
				continue
			}
			if *undoneFlag && task.Done {
				continue
			}

			// Filter by memo reference
			if *refsFlag != "" && !containsString(task.MemoRefs, *refsFlag) {
				continue
			}

			filteredTasks = append(filteredTasks, task)
		}

		// Sort tasks by order
		sortTasksByOrder(filteredTasks)

		// Print tasks
		if len(filteredTasks) > 0 {
			fmt.Println("Tasks:")
			for _, task := range filteredTasks {
				doneStr := "[ ]"
				if task.Done {
					doneStr = "[x]"
				}
				fmt.Printf("  %s  %.1f  %s  %s\n", task.ID[:8], task.Order, doneStr, task.Title)
			}
		} else {
			fmt.Println("No tasks found")
		}
	}

	if subCmd == "memos" || subCmd == "all" {
		// Filter memos
		var filteredMemos []*model.Memo
		for _, memo := range store.Memos {
			// Filter by reference
			if *refsFlag != "" {
				// Skip this memo if we're filtering by refs (memos don't reference other memos)
				continue
			}

			filteredMemos = append(filteredMemos, memo)
		}

		// Print memos
		if len(filteredMemos) > 0 {
			if subCmd == "all" {
				fmt.Println() // Add a newline if we're listing both tasks and memos
			}
			fmt.Println("Memos:")
			for _, memo := range filteredMemos {
				titleStr := "<no title>"
				if memo.Title != nil {
					titleStr = *memo.Title
				}

				// Get first line of content
				contentLines := strings.SplitN(memo.Content, "\n", 2)
				contentPreview := contentLines[0]
				if len(contentPreview) > 50 {
					contentPreview = contentPreview[:47] + "..."
				}

				fmt.Printf("  %s  %s  %s\n", memo.ID[:8], titleStr, contentPreview)
			}
		} else {
			fmt.Println("No memos found")
		}
	}

	return nil
}

// executeShow handles the 'show' command
func (c *CLI) executeShow(args []string) error {
	// Create flag set
	showCmd := flag.NewFlagSet("show", flag.ExitOnError)

	// Set usage
	showCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: tamo show <id>\n\n")
		fmt.Fprintf(os.Stderr, "Show details of a task or memo\n\n")
		showCmd.PrintDefaults()
	}

	// Parse flags
	if err := showCmd.Parse(args); err != nil {
		return err
	}

	// Check if ID is provided
	if showCmd.NArg() < 1 {
		return fmt.Errorf("missing ID")
	}

	// Get ID
	id := showCmd.Arg(0)

	// Load store
	s := storage.NewStorage()
	store, err := s.Load()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	// Try to find task by ID or prefix
	var task *model.Task
	if len(id) == 36 { // Full UUID
		task = store.FindTaskByID(id)
	} else {
		// Try to find by prefix
		for _, t := range store.Tasks {
			if strings.HasPrefix(t.ID, id) {
				task = t
				break
			}
		}
	}

	if task != nil {
		// Print task details
		doneStr := "[ ] Not completed"
		if task.Done {
			doneStr = "[x] Completed"
		}

		fmt.Printf("Task ID: %s\n", task.ID)
		fmt.Printf("Title: %s\n", task.Title)
		fmt.Printf("Order: %.1f\n", task.Order)
		fmt.Printf("Status: %s\n", doneStr)
		fmt.Printf("Created: %s\n", task.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated: %s\n", task.UpdatedAt.Format("2006-01-02 15:04:05"))

		if task.Description != "" {
			fmt.Println("\nDescription:")
			fmt.Println(task.Description)
		}

		if len(task.MemoRefs) > 0 {
			fmt.Println("\nReferenced Memos:")
			for _, memoID := range task.MemoRefs {
				memo := store.FindMemoByID(memoID)
				if memo != nil {
					titleStr := "<no title>"
					if memo.Title != nil {
						titleStr = *memo.Title
					}
					fmt.Printf("  %s  %s\n", memoID[:8], titleStr)
				} else {
					fmt.Printf("  %s  <memo not found>\n", memoID[:8])
				}
			}
		}

		return nil
	}

	// Try to find memo by ID or prefix
	var memo *model.Memo
	if len(id) == 36 { // Full UUID
		memo = store.FindMemoByID(id)
	} else {
		// Try to find by prefix
		for _, m := range store.Memos {
			if strings.HasPrefix(m.ID, id) {
				memo = m
				break
			}
		}
	}

	if memo != nil {
		// Print memo details
		fmt.Printf("Memo ID: %s\n", memo.ID)
		if memo.Title != nil {
			fmt.Printf("Title: %s\n", *memo.Title)
		}
		fmt.Printf("Created: %s\n", memo.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated: %s\n", memo.UpdatedAt.Format("2006-01-02 15:04:05"))

		referencingTasks := findTasksReferencingMemo(store, memo.ID)
		if len(referencingTasks) > 0 {
			fmt.Println("\nReference Tasks:")
			for _, task := range referencingTasks {
				fmt.Printf("%s %s\n", task.ID[:8], task.Title)
			}
		}

		fmt.Println("\nContent:")
		fmt.Println(memo.Content)

		return nil
	}

	return fmt.Errorf("no task or memo found with ID: %s", id)
}

// executeRemove handles the 'rm' command
func (c *CLI) executeRemove(args []string) error {
	// Manual argument parsing
	// Set usage
	usage := func() {
		fmt.Fprintf(os.Stderr, "Usage: tamo rm <id> [-f|--force]\n\n")
		fmt.Fprintf(os.Stderr, "Remove a task or memo\n\n")
		fmt.Fprintf(os.Stderr, "  -f, --force    Force removal without confirmation\n")
	}

	// Check if we have at least an ID
	if len(args) < 1 {
		usage()
		return fmt.Errorf("missing ID")
	}

	// Get ID
	id := args[0]

	// Check for force flag
	force := false
	for i := 1; i < len(args); i++ {
		if args[i] == "-f" || args[i] == "--force" {
			force = true
			break
		}
	}

	// Load store
	s := storage.NewStorage()
	store, err := s.Load()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	// Try to find task by ID or prefix
	var task *model.Task
	if len(id) == 36 { // Full UUID
		task = store.FindTaskByID(id)
	} else {
		// Try to find by prefix
		for _, t := range store.Tasks {
			if strings.HasPrefix(t.ID, id) {
				task = t
				break
			}
		}
	}

	if task != nil {
		// Remove task
		removeTask(store, task.ID)

		// Save store
		if err := s.Save(store); err != nil {
			return fmt.Errorf("failed to save data: %w", err)
		}

		fmt.Printf("Task '%s' removed\n", task.Title)
		return nil
	}

	// Try to find memo by ID or prefix
	var memo *model.Memo
	if len(id) == 36 { // Full UUID
		memo = store.FindMemoByID(id)
	} else {
		// Try to find by prefix
		for _, m := range store.Memos {
			if strings.HasPrefix(m.ID, id) {
				memo = m
				break
			}
		}
	}

	if memo != nil {
		// Check if memo is referenced by any tasks
		referencingTasks := findTasksReferencingMemo(store, memo.ID)
		if len(referencingTasks) > 0 {
			if !force {
				fmt.Printf("Memo is referenced by %d tasks. Use -f or --force to remove anyway.\n", len(referencingTasks))
				for _, task := range referencingTasks {
					fmt.Printf("  %s  %s\n", task.ID[:8], task.Title)
				}
				return fmt.Errorf("memo removal aborted")
			} else {
				fmt.Printf("Forcing removal of memo referenced by %d tasks\n", len(referencingTasks))
			}
		}

		// Remove memo
		removeMemo(store, memo.ID)

		// Save store
		if err := s.Save(store); err != nil {
			return fmt.Errorf("failed to save data: %w", err)
		}

		titleStr := "<no title>"
		if memo.Title != nil {
			titleStr = *memo.Title
		}
		fmt.Printf("Memo '%s' removed\n", titleStr)
		return nil
	}

	return fmt.Errorf("no task or memo found with ID: %s", id)
}

// Helper functions

// sortTasksByOrder sorts tasks by their order field
func sortTasksByOrder(tasks []*model.Task) {
	// Simple bubble sort for now
	for i := 0; i < len(tasks); i++ {
		for j := i + 1; j < len(tasks); j++ {
			if tasks[i].Order > tasks[j].Order {
				tasks[i], tasks[j] = tasks[j], tasks[i]
			}
		}
	}
}

// containsString checks if a string slice contains a string
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// removeTask removes a task from the store
func removeTask(store *model.Store, id string) {
	for i, task := range store.Tasks {
		if task.ID == id {
			// Remove task from slice
			store.Tasks = append(store.Tasks[:i], store.Tasks[i+1:]...)
			break
		}
	}
}

// removeMemo removes a memo from the store
func removeMemo(store *model.Store, id string) {
	for i, memo := range store.Memos {
		if memo.ID == id {
			// Remove memo from slice
			store.Memos = append(store.Memos[:i], store.Memos[i+1:]...)
			break
		}
	}

	// Also remove references to this memo from all tasks
	for _, task := range store.Tasks {
		for i, memoID := range task.MemoRefs {
			if memoID == id {
				// Remove reference from slice
				task.MemoRefs = append(task.MemoRefs[:i], task.MemoRefs[i+1:]...)
				break
			}
		}
	}
}

// findTasksReferencingMemo finds all tasks that reference a memo
func findTasksReferencingMemo(store *model.Store, memoID string) []*model.Task {
	var tasks []*model.Task
	for _, task := range store.Tasks {
		if containsString(task.MemoRefs, memoID) {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

// readLine reads a line from stdin
func readLine() string {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

// executeEdit handles the 'edit' command
func (c *CLI) executeEdit(args []string) error {
	// Create flag set
	editCmd := flag.NewFlagSet("edit", flag.ExitOnError)

	// Define flags
	editorFlag := editCmd.Bool("editor", false, "Use editor to edit content")

	// Set usage
	editCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: tamo edit <id> [--editor]\n\n")
		fmt.Fprintf(os.Stderr, "Edit a task or memo\n\n")
		editCmd.PrintDefaults()
	}

	// Parse flags
	if err := editCmd.Parse(args); err != nil {
		return err
	}

	// Check if ID is provided
	if editCmd.NArg() < 1 {
		return fmt.Errorf("missing ID")
	}

	// Get ID
	id := editCmd.Arg(0)

	// Load store
	s := storage.NewStorage()
	store, err := s.Load()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	// Try to find task by ID or prefix
	var task *model.Task
	if len(id) == 36 { // Full UUID
		task = store.FindTaskByID(id)
	} else {
		// Try to find by prefix
		for _, t := range store.Tasks {
			if strings.HasPrefix(t.ID, id) {
				task = t
				break
			}
		}
	}

	if task != nil {
		// Edit task
		return editTask(task, store, s, *editorFlag)
	}

	// Try to find memo by ID or prefix
	var memo *model.Memo
	if len(id) == 36 { // Full UUID
		memo = store.FindMemoByID(id)
	} else {
		// Try to find by prefix
		for _, m := range store.Memos {
			if strings.HasPrefix(m.ID, id) {
				memo = m
				break
			}
		}
	}

	if memo != nil {
		// Edit memo
		return editMemo(memo, store, s, *editorFlag)
	}

	return fmt.Errorf("no task or memo found with ID: %s", id)
}

// editTask edits a task using an editor or simple prompts
func editTask(task *model.Task, store *model.Store, s *storage.Storage, useEditor bool) error {
	if useEditor {
		// Get editor from environment
		editor := os.Getenv("EDITOR")
		if editor == "" {
			// Default to a simple editor if not set
			editor = "nano"
		}

		// Create temporary file
		tmpFile, err := ioutil.TempFile("", "tamo-task-*.md")
		if err != nil {
			return fmt.Errorf("failed to create temporary file: %w", err)
		}
		defer os.Remove(tmpFile.Name())

		// Write task content to temporary file
		content := fmt.Sprintf("# %s\n\n%s\n\n# Memo References (one ID per line):\n%s\n",
			task.Title,
			task.Description,
			strings.Join(task.MemoRefs, "\n"))

		if _, err := tmpFile.Write([]byte(content)); err != nil {
			tmpFile.Close()
			return fmt.Errorf("failed to write to temporary file: %w", err)
		}
		tmpFile.Close()

		// Open editor
		cmd := exec.Command(editor, tmpFile.Name())
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("editor command failed: %w", err)
		}

		// Read edited content
		editedContent, err := ioutil.ReadFile(tmpFile.Name())
		if err != nil {
			return fmt.Errorf("failed to read edited content: %w", err)
		}

		// Parse edited content
		lines := strings.Split(string(editedContent), "\n")

		// Extract title, description, and memo refs
		var title string
		var description strings.Builder
		var memoRefs []string

		mode := "title"
		for _, line := range lines {
			if mode == "title" && strings.HasPrefix(line, "# ") {
				title = strings.TrimPrefix(line, "# ")
				mode = "description"
			} else if mode == "description" && strings.HasPrefix(line, "# Memo References") {
				mode = "refs"
			} else if mode == "description" {
				description.WriteString(line)
				description.WriteString("\n")
			} else if mode == "refs" && line != "" && !strings.HasPrefix(line, "# ") {
				// Add memo ref if it's not empty and not a heading
				memoRefs = append(memoRefs, strings.TrimSpace(line))
			}
		}

		// Update task
		task.Title = title
		task.Description = strings.TrimSpace(description.String())
		task.MemoRefs = memoRefs
		task.UpdatedAt = model.CustomTime{Time: time.Now().UTC()}

		// Save store
		if err := s.Save(store); err != nil {
			return fmt.Errorf("failed to save data: %w", err)
		}

		fmt.Printf("Task '%s' updated\n", task.Title)
		return nil
	} else {
		// Simple prompt-based editing
		fmt.Printf("Editing task: %s\n", task.ID)

		// Edit title
		fmt.Printf("Title [%s]: ", task.Title)
		title := readLine()
		if title != "" {
			task.Title = title
		}

		// Edit description
		fmt.Printf("Description [Press Enter to keep, 'edit' to edit]:\n")
		descAction := readLine()
		if descAction == "edit" {
			fmt.Println("Enter new description (press Ctrl+D when finished):")
			scanner := bufio.NewScanner(os.Stdin)
			var descBuilder strings.Builder
			for scanner.Scan() {
				descBuilder.WriteString(scanner.Text())
				descBuilder.WriteString("\n")
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("error reading description: %w", err)
			}
			task.Description = strings.TrimSpace(descBuilder.String())
		}

		// Edit memo refs
		fmt.Printf("Memo References [%s] (comma-separated): ", strings.Join(task.MemoRefs, ","))
		refsStr := readLine()
		if refsStr != "" {
			task.MemoRefs = strings.Split(refsStr, ",")
			// Trim whitespace from each memo ID
			for i, ref := range task.MemoRefs {
				task.MemoRefs[i] = strings.TrimSpace(ref)
			}
		}

		// Update timestamp
		task.UpdatedAt = model.CustomTime{Time: time.Now().UTC()}

		// Save store
		if err := s.Save(store); err != nil {
			return fmt.Errorf("failed to save data: %w", err)
		}

		fmt.Printf("Task '%s' updated\n", task.Title)
		return nil
	}
}

// editMemo edits a memo using an editor or simple prompts
func editMemo(memo *model.Memo, store *model.Store, s *storage.Storage, useEditor bool) error {
	if useEditor {
		// Get editor from environment
		editor := os.Getenv("EDITOR")
		if editor == "" {
			// Default to a simple editor if not set
			editor = "nano"
		}

		// Create temporary file
		tmpFile, err := ioutil.TempFile("", "tamo-memo-*.md")
		if err != nil {
			return fmt.Errorf("failed to create temporary file: %w", err)
		}
		defer os.Remove(tmpFile.Name())

		// Write memo content to temporary file
		var content string
		if memo.Title != nil {
			content = fmt.Sprintf("# %s\n\n%s\n", *memo.Title, memo.Content)
		} else {
			content = fmt.Sprintf("# \n\n%s\n", memo.Content)
		}

		if _, err := tmpFile.Write([]byte(content)); err != nil {
			tmpFile.Close()
			return fmt.Errorf("failed to write to temporary file: %w", err)
		}
		tmpFile.Close()

		// Open editor
		cmd := exec.Command(editor, tmpFile.Name())
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("editor command failed: %w", err)
		}

		// Read edited content
		editedContent, err := ioutil.ReadFile(tmpFile.Name())
		if err != nil {
			return fmt.Errorf("failed to read edited content: %w", err)
		}

		// Parse edited content
		lines := strings.Split(string(editedContent), "\n")

		// Extract title and content
		var title string
		var contentBuilder strings.Builder

		mode := "title"
		for i, line := range lines {
			if i == 0 && strings.HasPrefix(line, "# ") {
				title = strings.TrimPrefix(line, "# ")
				mode = "content"
			} else if mode == "content" {
				contentBuilder.WriteString(line)
				contentBuilder.WriteString("\n")
			}
		}

		// Update memo
		if title != "" {
			memo.Title = &title
		} else {
			memo.Title = nil
		}
		memo.Content = strings.TrimSpace(contentBuilder.String())
		memo.UpdatedAt = model.CustomTime{Time: time.Now().UTC()}

		// Save store
		if err := s.Save(store); err != nil {
			return fmt.Errorf("failed to save data: %w", err)
		}

		titleStr := "<no title>"
		if memo.Title != nil {
			titleStr = *memo.Title
		}
		fmt.Printf("Memo '%s' updated\n", titleStr)
		return nil
	} else {
		// Simple prompt-based editing
		fmt.Printf("Editing memo: %s\n", memo.ID)

		// Edit title
		titleStr := "<no title>"
		if memo.Title != nil {
			titleStr = *memo.Title
		}
		fmt.Printf("Title [%s]: ", titleStr)
		title := readLine()
		if title != "" {
			memo.Title = &title
		} else if title == "<no title>" {
			memo.Title = nil
		}

		// Edit content
		fmt.Printf("Content [Press Enter to keep, 'edit' to edit]:\n")
		contentAction := readLine()
		if contentAction == "edit" {
			fmt.Println("Enter new content (press Ctrl+D when finished):")
			scanner := bufio.NewScanner(os.Stdin)
			var contentBuilder strings.Builder
			for scanner.Scan() {
				contentBuilder.WriteString(scanner.Text())
				contentBuilder.WriteString("\n")
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("error reading content: %w", err)
			}
			memo.Content = strings.TrimSpace(contentBuilder.String())
		}

		// Update timestamp
		memo.UpdatedAt = model.CustomTime{Time: time.Now().UTC()}

		// Save store
		if err := s.Save(store); err != nil {
			return fmt.Errorf("failed to save data: %w", err)
		}

		titleStr = "<no title>"
		if memo.Title != nil {
			titleStr = *memo.Title
		}
		fmt.Printf("Memo '%s' updated\n", titleStr)
		return nil
	}
}

// executeDone handles the 'done' command
func (c *CLI) executeDone(args []string) error {
	// Create flag set
	doneCmd := flag.NewFlagSet("done", flag.ExitOnError)

	// Set usage
	doneCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: tamo done <task_id>\n\n")
		fmt.Fprintf(os.Stderr, "Mark a task as done\n\n")
		doneCmd.PrintDefaults()
	}

	// Parse flags
	if err := doneCmd.Parse(args); err != nil {
		return err
	}

	// Check if task ID is provided
	if doneCmd.NArg() < 1 {
		return fmt.Errorf("missing task ID")
	}

	// Get task ID
	taskID := doneCmd.Arg(0)

	// Load store
	s := storage.NewStorage()
	store, err := s.Load()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	// Find task by ID or prefix
	var task *model.Task
	if len(taskID) == 36 { // Full UUID
		task = store.FindTaskByID(taskID)
	} else {
		// Try to find by prefix
		for _, t := range store.Tasks {
			if strings.HasPrefix(t.ID, taskID) {
				task = t
				break
			}
		}
	}

	if task == nil {
		return fmt.Errorf("no task found with ID: %s", taskID)
	}

	// Mark task as done
	task.Done = true
	task.UpdatedAt = model.CustomTime{Time: time.Now().UTC()}

	// Save store
	if err := s.Save(store); err != nil {
		return fmt.Errorf("failed to save data: %w", err)
	}

	fmt.Printf("Task '%s' marked as done\n", task.Title)
	return nil
}

// executeUndone handles the 'undone' command
func (c *CLI) executeUndone(args []string) error {
	// Create flag set
	undoneCmd := flag.NewFlagSet("undone", flag.ExitOnError)

	// Set usage
	undoneCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: tamo undone <task_id>\n\n")
		fmt.Fprintf(os.Stderr, "Mark a task as not done\n\n")
		undoneCmd.PrintDefaults()
	}

	// Parse flags
	if err := undoneCmd.Parse(args); err != nil {
		return err
	}

	// Check if task ID is provided
	if undoneCmd.NArg() < 1 {
		return fmt.Errorf("missing task ID")
	}

	// Get task ID
	taskID := undoneCmd.Arg(0)

	// Load store
	s := storage.NewStorage()
	store, err := s.Load()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	// Find task by ID or prefix
	var task *model.Task
	if len(taskID) == 36 { // Full UUID
		task = store.FindTaskByID(taskID)
	} else {
		// Try to find by prefix
		for _, t := range store.Tasks {
			if strings.HasPrefix(t.ID, taskID) {
				task = t
				break
			}
		}
	}

	if task == nil {
		return fmt.Errorf("no task found with ID: %s", taskID)
	}

	// Mark task as not done
	task.Done = false
	task.UpdatedAt = model.CustomTime{Time: time.Now().UTC()}

	// Save store
	if err := s.Save(store); err != nil {
		return fmt.Errorf("failed to save data: %w", err)
	}

	fmt.Printf("Task '%s' marked as not done\n", task.Title)
	return nil
}

// executeMove handles the 'mv' command
func (c *CLI) executeMove(args []string) error {
	// Manual argument parsing
	// Set usage
	usage := func() {
		fmt.Fprintf(os.Stderr, "Usage: tamo mv <task_id> <target_order>\n")
		fmt.Fprintf(os.Stderr, "       tamo mv <task_id> before|after <other_task_id>\n\n")
		fmt.Fprintf(os.Stderr, "Move a task to a specific order or relative to another task\n")
	}

	// Check if we have at least a task ID and a target
	if len(args) < 2 {
		usage()
		return fmt.Errorf("missing arguments")
	}

	// Get task ID
	taskID := args[0]

	// Load store
	s := storage.NewStorage()
	store, err := s.Load()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	// Find task by ID or prefix
	var task *model.Task
	if len(taskID) == 36 { // Full UUID
		task = store.FindTaskByID(taskID)
	} else {
		// Try to find by prefix
		for _, t := range store.Tasks {
			if strings.HasPrefix(t.ID, taskID) {
				task = t
				break
			}
		}
	}

	if task == nil {
		return fmt.Errorf("no task found with ID: %s", taskID)
	}

	// Sort tasks by order
	var tasks []*model.Task
	tasks = append(tasks, store.Tasks...)
	sortTasksByOrder(tasks)

	// Handle different move types
	if args[1] == "before" || args[1] == "after" {
		// Relative move
		if len(args) < 3 {
			usage()
			return fmt.Errorf("missing target task ID")
		}

		// Get target task ID
		targetTaskID := args[2]

		// Find target task
		var targetTask *model.Task
		if len(targetTaskID) == 36 { // Full UUID
			targetTask = store.FindTaskByID(targetTaskID)
		} else {
			// Try to find by prefix
			for _, t := range store.Tasks {
				if strings.HasPrefix(t.ID, targetTaskID) {
					targetTask = t
					break
				}
			}
		}

		if targetTask == nil {
			return fmt.Errorf("no target task found with ID: %s", targetTaskID)
		}

		// Calculate new order
		var newOrder float64

		if args[1] == "before" {
			// Find the task before the target task
			var prevTask *model.Task
			for i, t := range tasks {
				if t.ID == targetTask.ID && i > 0 {
					prevTask = tasks[i-1]
					break
				}
			}

			if prevTask != nil {
				// Place between prev and target
				newOrder = (prevTask.Order + targetTask.Order) / 2.0
			} else {
				// Place before the first task
				newOrder = targetTask.Order - 1.0
			}
		} else { // after
			// Find the task after the target task
			var nextTask *model.Task
			for i, t := range tasks {
				if t.ID == targetTask.ID && i < len(tasks)-1 {
					nextTask = tasks[i+1]
					break
				}
			}

			if nextTask != nil {
				// Place between target and next
				newOrder = (targetTask.Order + nextTask.Order) / 2.0
			} else {
				// Place after the last task
				newOrder = targetTask.Order + 1.0
			}
		}

		// Update task order
		task.Order = newOrder
		task.UpdatedAt = model.CustomTime{Time: time.Now().UTC()}

		// Save store
		if err := s.Save(store); err != nil {
			return fmt.Errorf("failed to save data: %w", err)
		}

		fmt.Printf("Task '%s' moved %s task '%s'\n", task.Title, args[1], targetTask.Title)
		return nil
	} else {
		// Absolute move
		targetOrder, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			usage()
			return fmt.Errorf("invalid target order: %s", args[1])
		}

		// Update task order
		task.Order = targetOrder
		task.UpdatedAt = model.CustomTime{Time: time.Now().UTC()}

		// Save store
		if err := s.Save(store); err != nil {
			return fmt.Errorf("failed to save data: %w", err)
		}

		fmt.Printf("Task '%s' moved to order %.1f\n", task.Title, targetOrder)
		return nil
	}
}

// executePop handles the 'pop task' command
func (c *CLI) executePop(args []string) error {
	// Manual argument parsing
	// Set usage
	usage := func() {
		fmt.Fprintf(os.Stderr, "Usage: tamo pop task [--done | --rm [-f]]\n\n")
		fmt.Fprintf(os.Stderr, "Show, mark as done, or remove the last task\n\n")
		fmt.Fprintf(os.Stderr, "  --done    Mark the last task as done\n")
		fmt.Fprintf(os.Stderr, "  --rm      Remove the last task\n")
		fmt.Fprintf(os.Stderr, "  -f        Force removal without confirmation\n")
	}

	// Check if we have at least the 'task' subcommand
	if len(args) < 1 || args[0] != "task" {
		usage()
		return fmt.Errorf("missing or invalid subcommand: expected 'task'")
	}

	// Parse options
	doneFlag := false
	rmFlag := false
	forceFlag := false

	for i := 1; i < len(args); i++ {
		if args[i] == "--done" {
			doneFlag = true
		} else if args[i] == "--rm" {
			rmFlag = true
		} else if args[i] == "-f" {
			forceFlag = true
		}
	}

	// Check for conflicting flags
	if doneFlag && rmFlag {
		return fmt.Errorf("--done and --rm flags cannot be used together")
	}

	// Load store
	s := storage.NewStorage()
	store, err := s.Load()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	// Find the last task (highest order)
	var lastTask *model.Task
	maxOrder := -1.0

	for _, task := range store.Tasks {
		if task.Order > maxOrder {
			lastTask = task
			maxOrder = task.Order
		}
	}

	if lastTask == nil {
		return fmt.Errorf("no tasks found")
	}

	// Handle different actions
	if doneFlag {
		// Mark as done
		lastTask.Done = true
		lastTask.UpdatedAt = model.CustomTime{Time: time.Now().UTC()}

		// Save store
		if err := s.Save(store); err != nil {
			return fmt.Errorf("failed to save data: %w", err)
		}

		fmt.Printf("Task '%s' marked as done\n", lastTask.Title)
	} else if rmFlag {
		// Remove task
		if !forceFlag {
			// Ask for confirmation
			fmt.Printf("Are you sure you want to remove task '%s'? (y/N): ", lastTask.Title)
			confirmation := readLine()
			if strings.ToLower(confirmation) != "y" {
				fmt.Println("Task removal aborted")
				return nil
			}
		}

		// Remove task
		removeTask(store, lastTask.ID)

		// Save store
		if err := s.Save(store); err != nil {
			return fmt.Errorf("failed to save data: %w", err)
		}

		fmt.Printf("Task '%s' removed\n", lastTask.Title)
	} else {
		// Show task details
		doneStr := "[ ] Not completed"
		if lastTask.Done {
			doneStr = "[x] Completed"
		}

		fmt.Printf("Task ID: %s\n", lastTask.ID)
		fmt.Printf("Title: %s\n", lastTask.Title)
		fmt.Printf("Order: %.1f\n", lastTask.Order)
		fmt.Printf("Status: %s\n", doneStr)
		fmt.Printf("Created: %s\n", lastTask.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated: %s\n", lastTask.UpdatedAt.Format("2006-01-02 15:04:05"))

		if lastTask.Description != "" {
			fmt.Println("\nDescription:")
			fmt.Println(lastTask.Description)
		}

		if len(lastTask.MemoRefs) > 0 {
			fmt.Println("\nReferenced Memos:")
			for _, memoID := range lastTask.MemoRefs {
				memo := store.FindMemoByID(memoID)
				if memo != nil {
					titleStr := "<no title>"
					if memo.Title != nil {
						titleStr = *memo.Title
					}
					fmt.Printf("  %s  %s\n", memoID[:8], titleStr)
				} else {
					fmt.Printf("  %s  <memo not found>\n", memoID[:8])
				}
			}
		}
	}

	return nil
}

// executeShift handles the 'shift task' command
func (c *CLI) executeShift(args []string) error {
	// Manual argument parsing
	// Set usage
	usage := func() {
		fmt.Fprintf(os.Stderr, "Usage: tamo shift task [--done | --rm [-f]]\n\n")
		fmt.Fprintf(os.Stderr, "Show, mark as done, or remove the first task\n\n")
		fmt.Fprintf(os.Stderr, "  --done    Mark the first task as done\n")
		fmt.Fprintf(os.Stderr, "  --rm      Remove the first task\n")
		fmt.Fprintf(os.Stderr, "  -f        Force removal without confirmation\n")
	}

	// Check if we have at least the 'task' subcommand
	if len(args) < 1 || args[0] != "task" {
		usage()
		return fmt.Errorf("missing or invalid subcommand: expected 'task'")
	}

	// Parse options
	doneFlag := false
	rmFlag := false
	forceFlag := false

	for i := 1; i < len(args); i++ {
		if args[i] == "--done" {
			doneFlag = true
		} else if args[i] == "--rm" {
			rmFlag = true
		} else if args[i] == "-f" {
			forceFlag = true
		}
	}

	// Check for conflicting flags
	if doneFlag && rmFlag {
		return fmt.Errorf("--done and --rm flags cannot be used together")
	}

	// Load store
	s := storage.NewStorage()
	store, err := s.Load()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	// Find the first task (lowest order)
	var firstTask *model.Task
	minOrder := math.MaxFloat64

	for _, task := range store.Tasks {
		if task.Order < minOrder {
			firstTask = task
			minOrder = task.Order
		}
	}

	if firstTask == nil {
		return fmt.Errorf("no tasks found")
	}

	// Handle different actions
	if doneFlag {
		// Mark as done
		firstTask.Done = true
		firstTask.UpdatedAt = model.CustomTime{Time: time.Now().UTC()}

		// Save store
		if err := s.Save(store); err != nil {
			return fmt.Errorf("failed to save data: %w", err)
		}

		fmt.Printf("Task '%s' marked as done\n", firstTask.Title)
	} else if rmFlag {
		// Remove task
		if !forceFlag {
			// Ask for confirmation
			fmt.Printf("Are you sure you want to remove task '%s'? (y/N): ", firstTask.Title)
			confirmation := readLine()
			if strings.ToLower(confirmation) != "y" {
				fmt.Println("Task removal aborted")
				return nil
			}
		}

		// Remove task
		removeTask(store, firstTask.ID)

		// Save store
		if err := s.Save(store); err != nil {
			return fmt.Errorf("failed to save data: %w", err)
		}

		fmt.Printf("Task '%s' removed\n", firstTask.Title)
	} else {
		// Show task details
		doneStr := "[ ] Not completed"
		if firstTask.Done {
			doneStr = "[x] Completed"
		}

		fmt.Printf("Task ID: %s\n", firstTask.ID)
		fmt.Printf("Title: %s\n", firstTask.Title)
		fmt.Printf("Order: %.1f\n", firstTask.Order)
		fmt.Printf("Status: %s\n", doneStr)
		fmt.Printf("Created: %s\n", firstTask.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated: %s\n", firstTask.UpdatedAt.Format("2006-01-02 15:04:05"))

		if firstTask.Description != "" {
			fmt.Println("\nDescription:")
			fmt.Println(firstTask.Description)
		}

		if len(firstTask.MemoRefs) > 0 {
			fmt.Println("\nReferenced Memos:")
			for _, memoID := range firstTask.MemoRefs {
				memo := store.FindMemoByID(memoID)
				if memo != nil {
					titleStr := "<no title>"
					if memo.Title != nil {
						titleStr = *memo.Title
					}
					fmt.Printf("  %s  %s\n", memoID[:8], titleStr)
				} else {
					fmt.Printf("  %s  <memo not found>\n", memoID[:8])
				}
			}
		}
	}

	return nil
}

// executeNext handles the 'next' command (alias for shift task with focus on undone tasks)
func (c *CLI) executeNext(args []string) error {
	// Load store
	s := storage.NewStorage()
	store, err := s.Load()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	// Find the first undone task (lowest order)
	var firstUndoneTask *model.Task
	minOrder := math.MaxFloat64

	for _, task := range store.Tasks {
		if !task.Done && task.Order < minOrder {
			firstUndoneTask = task
			minOrder = task.Order
		}
	}

	if firstUndoneTask == nil {
		return fmt.Errorf("no undone tasks found")
	}

	// Show task details
	fmt.Printf("Task ID: %s\n", firstUndoneTask.ID)
	fmt.Printf("Title: %s\n", firstUndoneTask.Title)
	fmt.Printf("Order: %.1f\n", firstUndoneTask.Order)
	fmt.Printf("Status: [ ] Not completed\n")
	fmt.Printf("Created: %s\n", firstUndoneTask.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated: %s\n", firstUndoneTask.UpdatedAt.Format("2006-01-02 15:04:05"))

	if firstUndoneTask.Description != "" {
		fmt.Println("\nDescription:")
		fmt.Println(firstUndoneTask.Description)
	}

	if len(firstUndoneTask.MemoRefs) > 0 {
		fmt.Println("\nReferenced Memos:")
		for _, memoID := range firstUndoneTask.MemoRefs {
			memo := store.FindMemoByID(memoID)
			if memo != nil {
				titleStr := "<no title>"
				if memo.Title != nil {
					titleStr = *memo.Title
				}
				fmt.Printf("  %s  %s\n", memoID[:8], titleStr)
			} else {
				fmt.Printf("  %s  <memo not found>\n", memoID[:8])
			}
		}
	}

	return nil
}

// executeFlattask handles the 'flattask' command
func (c *CLI) executeFlattask(args []string) error {
	// Create flag set
	flattaskCmd := flag.NewFlagSet("flattask", flag.ExitOnError)

	// Set usage
	flattaskCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: tamo flattask <task_id>\n\n")
		fmt.Fprintf(os.Stderr, "Flatten a task by expanding all memo references\n\n")
		flattaskCmd.PrintDefaults()
	}

	// Parse flags
	if err := flattaskCmd.Parse(args); err != nil {
		return err
	}

	// Check if task ID is provided
	if flattaskCmd.NArg() < 1 {
		return fmt.Errorf("missing task ID")
	}

	// Get task ID
	taskID := flattaskCmd.Arg(0)

	// Load store
	s := storage.NewStorage()
	store, err := s.Load()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	// Find task by ID or prefix
	var task *model.Task
	if len(taskID) == 36 { // Full UUID
		task = store.FindTaskByID(taskID)
	} else {
		// Try to find by prefix
		for _, t := range store.Tasks {
			if strings.HasPrefix(t.ID, taskID) {
				task = t
				break
			}
		}
	}

	if task == nil {
		return fmt.Errorf("no task found with ID: %s", taskID)
	}

	// Generate Markdown document
	var doc strings.Builder

	// Add task title and status
	doc.WriteString(fmt.Sprintf("# %s\n\n", task.Title))

	if task.Done {
		doc.WriteString("**Status:** Completed\n\n")
	} else {
		doc.WriteString("**Status:** Not completed\n\n")
	}

	// Add task description if available
	if task.Description != "" {
		doc.WriteString("## Description\n\n")
		doc.WriteString(task.Description)
		doc.WriteString("\n\n")
	}

	// Add referenced memos
	if len(task.MemoRefs) > 0 {
		doc.WriteString("## Referenced Memos\n\n")

		for _, memoID := range task.MemoRefs {
			memo := store.FindMemoByID(memoID)
			if memo != nil {
				// Add memo title
				if memo.Title != nil {
					doc.WriteString(fmt.Sprintf("### %s\n\n", *memo.Title))
				} else {
					doc.WriteString(fmt.Sprintf("### Memo %s\n\n", memoID[:8]))
				}

				// Add memo content
				doc.WriteString(memo.Content)
				doc.WriteString("\n\n")
			} else {
				doc.WriteString(fmt.Sprintf("### Memo %s (not found)\n\n", memoID[:8]))
			}

		}
	}

	// Print the document
	fmt.Println(doc.String())

	return nil
}

// executeAddTaskFromMarkdown handles the 'add task' command with Markdown parsing
func (c *CLI) executeAddTaskFromMarkdown(args []string) error {
	// Check if we have the right arguments
	if len(args) == 0 {
		return fmt.Errorf("missing arguments for Markdown parsing")
	}

	// Parse options
	var filePath string
	fromStdin := false

	if args[0] == "-f" {
		if len(args) < 2 {
			return fmt.Errorf("missing file path after -f")
		}
		filePath = args[1]
	} else if args[0] == "--from-stdin" {
		fromStdin = true
	} else {
		return fmt.Errorf("invalid option: %s", args[0])
	}

	// Load store
	s := storage.NewStorage()
	store, err := s.Load()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	// Create parser
	parser := NewMarkdownParser(store)

	// Parse Markdown
	var task *model.Task
	var memos []*model.Memo

	if fromStdin {
		task, memos, err = parser.ParseFromStdin()
	} else {
		task, memos, err = parser.ParseFromFile(filePath)
	}

	if err != nil {
		return fmt.Errorf("failed to parse Markdown: %w", err)
	}

	// Save task and memos
	if err := parser.SaveTaskAndMemos(task, memos, s); err != nil {
		return fmt.Errorf("failed to save task and memos: %w", err)
	}

	// Print success message
	fmt.Printf("Task added with ID: %s\n", task.ID)
	if len(memos) > 0 {
		fmt.Printf("Created %d memos:\n", len(memos))
		for _, memo := range memos {
			fmt.Printf("  Memo ID: %s\n", memo.ID[:8])
		}
	}

	return nil
}
