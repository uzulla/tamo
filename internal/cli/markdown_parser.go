package cli

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/zishida/tamo/internal/model"
	"github.com/zishida/tamo/internal/storage"
	"github.com/zishida/tamo/internal/utils"
)

// MarkdownParser handles parsing Markdown files to extract tasks and memos
type MarkdownParser struct {
	store *model.Store
}

// NewMarkdownParser creates a new MarkdownParser
func NewMarkdownParser(store *model.Store) *MarkdownParser {
	return &MarkdownParser{
		store: store,
	}
}

// ParseFromFile parses a Markdown file and extracts task and memos
func (p *MarkdownParser) ParseFromFile(filePath string) (*model.Task, []*model.Memo, error) {
	// Read file content
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Get filename for default title
	filename := filePath
	if lastSlash := strings.LastIndex(filePath, "/"); lastSlash >= 0 {
		filename = filePath[lastSlash+1:]
	}
	if lastDot := strings.LastIndex(filename, "."); lastDot >= 0 {
		filename = filename[:lastDot]
	}

	return p.parseMarkdown(string(content), filename)
}

// ParseFromStdin parses Markdown content from stdin
func (p *MarkdownParser) ParseFromStdin() (*model.Task, []*model.Memo, error) {
	// Read from stdin
	scanner := bufio.NewScanner(os.Stdin)
	var contentBuilder strings.Builder
	for scanner.Scan() {
		contentBuilder.WriteString(scanner.Text())
		contentBuilder.WriteString("\n")
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("error reading from stdin: %w", err)
	}
	content := contentBuilder.String()

	return p.parseMarkdown(content, "Task from stdin")
}

// parseMarkdown parses Markdown content and extracts task and memos
func (p *MarkdownParser) parseMarkdown(content, defaultTitle string) (*model.Task, []*model.Memo, error) {
	// Extract title (first H1 heading)
	title := defaultTitle
	titleRegex := regexp.MustCompile(`(?m)^# (.+)$`)
	titleMatch := titleRegex.FindStringSubmatch(content)
	if len(titleMatch) > 1 {
		title = titleMatch[1]
		// Remove the title from the content
		content = titleRegex.ReplaceAllString(content, "")
	}

	// Extract memo blocks
	memoRegex := regexp.MustCompile("(?s)```memo\n(.*?)\n```")
	memoMatches := memoRegex.FindAllStringSubmatch(content, -1)

	// Create memos and replace blocks with references
	var memos []*model.Memo
	for _, match := range memoMatches {
		if len(match) > 1 {
			// Generate UUID for memo
			memoID, err := utils.GenerateUUID()
			if err != nil {
				return nil, nil, fmt.Errorf("failed to generate UUID for memo: %w", err)
			}

			// Create memo
			memo := model.NewMemo(memoID, nil, match[1])
			memos = append(memos, memo)

			// Replace memo block with reference
			memoRef := fmt.Sprintf("[memo](%s)", memoID)
			content = strings.Replace(content, match[0], memoRef, 1)
		}
	}

	// Clean up the content (remove extra newlines, etc.)
	content = strings.TrimSpace(content)

	// Generate UUID for task
	taskID, err := utils.GenerateUUID()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate UUID for task: %w", err)
	}

	// Extract memo references
	var memoRefs []string
	for _, memo := range memos {
		memoRefs = append(memoRefs, memo.ID)
	}

	// Create task
	task := model.NewTask(taskID, title, content, memoRefs)

	// Set task order to max + 1.0
	task.Order = p.store.GetMaxTaskOrder() + 1.0

	return task, memos, nil
}

// SaveTaskAndMemos saves the task and memos to the store
func (p *MarkdownParser) SaveTaskAndMemos(task *model.Task, memos []*model.Memo, s *storage.Storage) error {
	// Add memos to store
	for _, memo := range memos {
		p.store.AddMemo(memo)
	}

	// Add task to store
	p.store.AddTask(task)

	// Save store
	if err := s.Save(p.store); err != nil {
		return fmt.Errorf("failed to save data: %w", err)
	}

	return nil
}
