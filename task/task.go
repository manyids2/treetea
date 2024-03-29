package task

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Datetime format of Taskwarrior
type ISO8601 time.Time

// Format for go to parse
const FORMAT_ISO8601 = "20060102T150405Z"

// Task struct defined at [Taskwarrior JSON format]()
type Task struct {
	ID          int               `json:"id,omitempty"` // NOTE: ID is not a key!
	UUID        string            `json:"uuid"`         // UUID is primary key
	Status      string            `json:"status"`
	Description string            `json:"description"`
	Entry       ISO8601           `json:"entry"`
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

func (j *ISO8601) UnmarshalJSON(b []byte) error {
	d := strings.Trim(string(b), "\"")
	t, e := time.Parse(FORMAT_ISO8601, string(d))
	if e != nil {
		t = time.Time{}
	}
	*j = ISO8601(t)
	return nil
}

func (j ISO8601) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(j))
}

func (t Task) Key() string {
	return t.UUID
}

func (t Task) Children() []string {
	return t.Depends
}

func (t Task) Desc(name string) string {
	switch name {
	case "due":
		if time.Time(t.Due).IsZero() {
			return ""
		}
		return fmt.Sprintf("  %s", time.Time(t.Due).Format("2 Jan"))
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
	case "description":
		return t.Description
	}
	return ""
}

type StringItem string

func (t StringItem) Key() string {
	return string(t)
}

func (t StringItem) Children() []string {
	return []string{}
}

func (t StringItem) Desc(name string) string {
	return string(t)
}
