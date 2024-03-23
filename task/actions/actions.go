package actions

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	tw "github.com/manyids2/tasktea/task"
)

// List Get list of tasks given filters
//
// `task [filters] export`
func List(filters []string) ([]tw.Task, error) {
	// Setup command
	filters = append(filters, "export")
	cmd := exec.Command("task", filters...)

	// Get json formatted array of tasks, return empty if error
	out, err := cmd.Output()
	if err != nil {
		return []tw.Task{}, err
	}

	// Parse json output, return empty if error
	var tasks_list []tw.Task
	if err := json.Unmarshal(out, &tasks_list); err != nil {
		return []tw.Task{}, err
	}

	return tasks_list, err
}

// List Get list of active tasks
//
// `task export active`
func Active() []string {
	active := []string{}
	cmd := exec.Command("task", "export", "active")
	out, err := cmd.Output()
	if err != nil {
		return active
	}

	// Parse json output, return empty if error
	var tasks_list []tw.Task
	if err := json.Unmarshal(out, &tasks_list); err != nil {
		return active
	}

	// Add all active tasks
	for _, v := range tasks_list {
		active = append(active, v.UUID)
	}
	return active
}

// Contexts Get all contexts
//
// `task _context`
func Contexts() []string {
	// Run context command
	contexts := []string{}
	cmd := exec.Command("task", "_context")
	out, err := cmd.Output()
	if err != nil {
		return contexts
	}

	// Split by lines
	lines := strings.Split(string(out), "\n")
	contexts = lines[:len(lines)-1]
	return contexts
}

// Context Parse current context to get read and write filters
//
// `task context show`
func Context() (string, string, string) {
	// Run context command
	context := "none"
	read_filters := ""
	write_filters := ""
	cmd := exec.Command("task", "context", "show")
	out, err := cmd.Output()
	if err != nil {
		return context, read_filters, write_filters
	}

	// Split by lines
	lines := strings.Split(string(out), "\n")

	// Context name
	r := regexp.MustCompile(`Context '(?P<Context>[a-zA-Z]+)' with`)
	match := r.FindStringSubmatch(lines[0])
	if len(match) != 2 {
		return context, read_filters, write_filters
	}
	context = match[1]

	// Read filter
	r = regexp.MustCompile(`\* read filter: '(?P<ReadFilter>.*)'`)
	match = r.FindStringSubmatch(lines[2])
	if len(match) != 2 {
		return context, read_filters, write_filters
	}
	read_filters = match[1]

	// Write filter
	r = regexp.MustCompile(`\* write filter: '(?P<WriteFilter>.*)'`)
	match = r.FindStringSubmatch(lines[3])
	if len(match) != 2 {
		return context, read_filters, write_filters
	}
	write_filters = match[1]

	return context, read_filters, write_filters
}

// Projects Get all Projects
//
// `task _projects`
func Projects() []string {
	// Run context command
	projects := []string{}
	cmd := exec.Command("task", "_projects")
	out, err := cmd.Output()
	if err != nil {
		return projects
	}

	// Split by lines
	lines := strings.Split(string(out), "\n")
	projects = lines[:len(lines)-1]
	return projects
}

// Projects Get all Projects
//
// `task context [context]`
func SetContext(context string) error {
	cmd := exec.Command("task", "context", context)
	_, err := cmd.Output()
	return err
}

// SetStatus Set task status
//
// `task [uuid] modify status:[status]`
func SetStatus(uuid string, status string) error {
	args := []string{uuid, "modify", fmt.Sprintf("status:%s", status)}
	cmd := exec.Command("task", args...)
	_, err := cmd.Output()
	return err
}

// SetActive Start/stop task
//
// `task [uuid] [start/stop]`
func SetActive(uuid string, startstop string) error {
	cmd := exec.Command("task", uuid, startstop)
	_, err := cmd.Output()
	return err
}

// LinkToParent Make task child of parent
//
// `task [uuid] modify depends:[parent]`
func LinkToParent(uuid string, parent string) error {
	args := []string{parent, "modify", fmt.Sprintf("depends:%s", uuid)}
	cmd := exec.Command("task", args...)
	_, err := cmd.Output()
	return err
}

// UnlinkFromParent In effect, dedent the task
//
// `task [uuid] modify depends:`
// `task [parent] modify depends:[all children other than uuid]`
// `task [grandparent] modify depends:[uuid]`
func UnlinkFromParent(uuid string, parent string, grandparent string, depends []string) error {
	// Remove all dependencies of parent
	args := []string{parent, "modify", "depends:"}
	cmd := exec.Command("task", args...)
	_, err := cmd.Output()
	if err != nil {
		return err
	}

	// Add back non uuid depends
	for _, v := range depends {
		if v == uuid {
			continue
		}
		args := []string{parent, "modify", fmt.Sprintf("depends:%s", v)}
		cmd := exec.Command("task", args...)
		_, err := cmd.Output()
		if err != nil {
			return err
		}
	}

	// Add current to grandparent
	if grandparent != "" {
		args := []string{grandparent, "modify", fmt.Sprintf("depends:%s", uuid)}
		cmd := exec.Command("task", args...)
		_, err := cmd.Output()
		if err != nil {
			return err
		}
	}

	return err
}

// Add Add task, get uuid and update parent's children
//
// `task add [tags] '[description]'`
// `task [parent] modify depends:[uuid]`
func Add(filters []string, description string, parent string) (string, error) {
	// Removing things like `-tag`, `/.../` which dont make sense
	// TODO: Maybe use regex, read spec to find out also whats up
	tags := []string{}
	for _, f := range filters {
		if (string(f[0]) != "-") && (string(f[0]) != "/") {
			tags = append(tags, f)
		}
	}

	// Add task
	args := []string{"add"}
	args = append(args, tags...)
	args = append(args, fmt.Sprintf("'%s'", description))
	cmd := exec.Command("task", args...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Extract assigned UUID
	// TODO: use regex or something a bit more solid
	parts := strings.Split(string(out), " ")
	uuid := parts[len(parts)-1]
	uuid = strings.TrimRight(uuid, "\n")
	uuid = strings.TrimRight(uuid, ".")

	// Update parent's dependencies
	if parent != "" {
		args = []string{parent, "modify", fmt.Sprintf("depends:%s", uuid)}
		cmd := exec.Command("task", args...)
		_, err = cmd.Output()
	}
	return uuid, err
}

// Modify task description
//
// `task [uuid] modify '[description]'`
func ModifyDescription(uuid string, description string) error {
	args := []string{uuid, "modify", fmt.Sprintf("'%s'", description)}
	cmd := exec.Command("task", args...)
	_, err := cmd.Output()
	return err
}

// Modify task with arguments
//
// `task [uuid] modify [args]`
func Modify(uuid string, args []string) error {
	out := []string{uuid, "modify"}
	out = append(out, args...)
	cmd := exec.Command("task", out...)
	_, err := cmd.Output()
	return err
}

// ModifyBatch Modify multiple tasks with same arguments
//
// `task [uuids] modify [args]`
func ModifyBatch(uuids []string, args []string) error {
	out := append(uuids, "modify")
	out = append(out, args...)
	cmd := exec.Command("task", out...)
	_, err := cmd.Output()
	return err
}

// Delete task, turn of confirmation
//
// `task [uuid] rc.confirmation=off delete`
func Delete(uuid string) error {
	args := []string{uuid, "rc.confirmation=off", "delete"}
	cmd := exec.Command("task", args...)
	_, err := cmd.Output()
	return err
}
