package openspec

import (
	"errors"
	"regexp"
	"strings"
)

// ErrAmbiguousTask is returned by ToggleTaskByText when the cursor's text
// matches more than one task. Task descriptions are expected to be unique
// within a file (OpenSpec tasks are numbered), so a duplicate is malformed and
// the toggle refuses rather than flipping the wrong line (#115).
var ErrAmbiguousTask = errors.New("ambiguous task: multiple tasks share this text")

type ItemKind int

const (
	KindSection ItemKind = iota
	KindTask
)

// String / MarshalJSON give ItemKind a stable, human-readable, cross-language
// serialization ("section"/"task") for the golden harness, instead of the raw
// iota int. No production code marshals TaskItem today; this only affects the
// shared golden contract (see testdata/corpus/README.md).
func (k ItemKind) String() string {
	if k == KindTask {
		return "task"
	}
	return "section"
}

func (k ItemKind) MarshalJSON() ([]byte, error) {
	return []byte(`"` + k.String() + `"`), nil
}

type TaskItem struct {
	Kind    ItemKind `json:"kind"`
	Text    string   `json:"text"`
	Done    bool     `json:"done"`
	LineNum int      `json:"line_num"`
}

var (
	rxSection = regexp.MustCompile(`^## (.+)$`)
	rxPending = regexp.MustCompile(`^- \[ \] (.+)$`)
	rxDone    = regexp.MustCompile(`^- \[x\] (.+)$`)
)

// splitLines splits s on "\n" and strips a trailing "\r" from each line, so
// CRLF-authored files parse identically to LF ones. Use it for parsing and
// content matching — NOT on a write path that rejoins and persists the file,
// since that would rewrite CRLF line endings to LF.
func splitLines(s string) []string {
	lines := strings.Split(s, "\n")
	for i := range lines {
		lines[i] = strings.TrimSuffix(lines[i], "\r")
	}
	return lines
}

func ParseTasks(content string) []TaskItem {
	lines := splitLines(content)
	items := make([]TaskItem, 0, len(lines))
	for i, line := range lines {
		switch {
		case rxSection.MatchString(line):
			m := rxSection.FindStringSubmatch(line)
			items = append(items, TaskItem{Kind: KindSection, Text: m[1], LineNum: i})
		case rxPending.MatchString(line):
			m := rxPending.FindStringSubmatch(line)
			items = append(items, TaskItem{Kind: KindTask, Text: m[1], Done: false, LineNum: i})
		case rxDone.MatchString(line):
			m := rxDone.FindStringSubmatch(line)
			items = append(items, TaskItem{Kind: KindTask, Text: m[1], Done: true, LineNum: i})
		}
	}
	return items
}

// FindCursorByText returns the index of the first KindTask item with the given
// text, or the index of the first KindTask item if the text is not found.
func FindCursorByText(items []TaskItem, text string) int {
	first := -1
	for i, item := range items {
		if item.Kind != KindTask {
			continue
		}
		if first == -1 {
			first = i
		}
		if item.Text == text {
			return i
		}
	}
	if first == -1 {
		return 0
	}
	return first
}

// countTasksWithText counts KindTask items whose text matches exactly. Used to
// detect an ambiguous (duplicate) toggle target; see ErrAmbiguousTask.
func countTasksWithText(items []TaskItem, text string) int {
	n := 0
	for _, it := range items {
		if it.Kind == KindTask && it.Text == text {
			n++
		}
	}
	return n
}

// ToggleTask flips the done state of items[idx] in memory and on disk.
func (l *Loader) ToggleTask(path string, items []TaskItem, idx int) error {
	data, err := l.fs.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	if idx >= len(items) || items[idx].LineNum >= len(lines) {
		return nil
	}
	if items[idx].Done {
		lines[items[idx].LineNum] = strings.Replace(lines[items[idx].LineNum], "- [x] ", "- [ ] ", 1)
		items[idx].Done = false
	} else {
		lines[items[idx].LineNum] = strings.Replace(lines[items[idx].LineNum], "- [ ] ", "- [x] ", 1)
		items[idx].Done = true
	}
	return l.fs.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644)
}

// ToggleTaskPackage is a package-level wrapper for ToggleTask.
func ToggleTask(path string, items []TaskItem, idx int) error {
	return defaultLoader.ToggleTask(path, items, idx)
}

// ToggleTaskByText re-reads and re-parses the file, locates the task by its
// text, and toggles that fresh line — never trusting a line index captured at
// render time, so an external edit between render and keypress can't flip the
// wrong line. It returns the freshly parsed items (so the caller renders the
// real on-disk state). If no task with the given text is found it makes no
// change and returns the current items. CRLF endings are preserved (the write
// path splits on raw "\n").
func (l *Loader) ToggleTaskByText(path, text string) ([]TaskItem, error) {
	data, err := l.fs.ReadFile(path)
	if err != nil {
		return nil, err
	}
	content := string(data)
	items := ParseTasks(content)
	idx := FindCursorByText(items, text)
	if idx >= len(items) || items[idx].Kind != KindTask || items[idx].Text != text {
		return items, nil
	}
	// Refuse to toggle when the text is ambiguous (duplicate descriptions are a
	// malformed file) rather than flipping the first match (#115).
	if countTasksWithText(items, text) > 1 {
		return items, ErrAmbiguousTask
	}
	lines := strings.Split(content, "\n")
	ln := items[idx].LineNum
	if ln >= len(lines) {
		return items, nil
	}
	if items[idx].Done {
		lines[ln] = strings.Replace(lines[ln], "- [x] ", "- [ ] ", 1)
	} else {
		lines[ln] = strings.Replace(lines[ln], "- [ ] ", "- [x] ", 1)
	}
	newContent := strings.Join(lines, "\n")
	if err := l.fs.WriteFile(path, []byte(newContent), 0644); err != nil {
		return nil, err
	}
	return ParseTasks(newContent), nil
}

// ToggleTaskByText is a package-level wrapper for Loader.ToggleTaskByText.
func ToggleTaskByText(path, text string) ([]TaskItem, error) {
	return defaultLoader.ToggleTaskByText(path, text)
}
