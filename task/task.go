package task

import (
	"fmt"
	"strings"
	"time"
)

// Datetime format of Taskwarrior
type ISO8601 string

// Format for go to parse
const FORMAT_ISO8601 = "20060102T150405Z"

// Text or unicode markers for status
var MARKERS = map[string]map[string]string{
	"icon": {
		"pending":   "󰝦 ",
		"deleted":   "󰩺 ",
		"completed": " ",
		"waiting":   "󰔛 ",
		"recurring": "󰑐 ",
		"default":   "󰝦 ",
	},
	"text": {
		"pending":   "- [ ]",
		"deleted":   "- [ ]",
		"waiting":   "- [ ]",
		"recurring": "- [ ]",
		"default":   "- [ ]",
		// Only this is recognized due to checkbox constraint in markdown syntax
		"completed": "- [x]",
	},
}

// ISO8601_to_DateTime Convert timewarrior format to go datetime with proper
// default in case of empty string
func ISO8601_to_DateTime(d ISO8601) time.Time {
	t, e := time.Parse(FORMAT_ISO8601, string(d))
	if e != nil {
		t = time.Time{}
	}
	return t
}

// Task struct defined at [Taskwarrior JSON format]()
type Task struct {
	ID          int               `json:"id,omitempty"`
	Status      string            `json:"status"`
	UUID        string            `json:"uuid"`
	Entry       ISO8601           `json:"entry"`
	Description string            `json:"description"`
	Start       ISO8601           `json:"start,omitempty"`
	End         ISO8601           `json:"end,omitempty"`
	Due         ISO8601           `json:"due,omitempty"`
	Until       ISO8601           `json:"until,omitempty"`
	Wait        ISO8601           `json:"wait,omitempty"`
	Modified    ISO8601           `json:"modified,omitempty"`
	Scheduled   ISO8601           `json:"scheduled,omitempty"`
	Recur       string            `json:"recur,omitempty"`
	Mask        string            `json:"mask,omitempty"`
	Imask       uint32            `json:"imask,omitempty"`
	Parent      string            `json:"parent,omitempty"`
	Project     string            `json:"project,omitempty"`
	Priority    string            `json:"priority,omitempty"`
	Depends     []string          `json:"depends,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	Annotation  map[string]string `json:"annotation,omitempty"`
}

func (t Task) Key() string {
	return t.UUID
}

func (t Task) Val() string {
	return t.Description
}

func (t Task) Children() []string {
	return t.Depends
}

func (t Task) String() string {
	return fmt.Sprintf("%s %s", MARKERS["icon"][t.Status], t.Description)
}

func (t Task) Extra(name string) string {
	switch name {
	case "due":
		if t.Due == "" {
			return ""
		}
		return fmt.Sprintf("  %s", ISO8601_to_DateTime(t.Due).Format("2 Jan"))
	case "tags":
		if len(t.Tags) == 0 {
			return ""
		}
		return "  " + strings.Join(t.Tags, " ")
	case "uuid":
		return t.UUID[:8]
	case "id":
		return fmt.Sprintf("%d", t.ID)
	case "project":
		return t.Project
	}
	return ""
}

// Context struct
type ContextS string

func (c ContextS) Key() string {
	return string(c)
}

func (c ContextS) Val() string {
	return string(c)
}

func (c ContextS) Children() []string {
	return []string{}
}

func (c ContextS) Extra(string) string {
	return ""
}
