package task

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// List Get list of tasks given filters
func List(filters []string) ([]Task, error) {
	// Setup command
	filters = append(filters, "export")
	cmd := exec.Command("task", filters...)

	// Get json formatted array of tasks, return empty if error
	out, err := cmd.Output()
	if err != nil {
		return []Task{}, err
	}

	// Parse json output, return empty if error
	var tasks_list []Task
	if err := json.Unmarshal(out, &tasks_list); err != nil {
		return []Task{}, err
	}

	return tasks_list, err
}

// Parse current context
func Context() string {
	// Run context command
	context := "none"
	cmd := exec.Command("task", "context", "show")
	out, err := cmd.Output()
	if err != nil {
		return context
	}

	// Find in first line of output
	lines := strings.Split(string(out), "\n")
	r := regexp.MustCompile(`Context '(?P<Context>[a-zA-Z]+)' with`)
	match := r.FindStringSubmatch(lines[0])
	if len(match) != 2 {
		return context
	}
	context = match[1]
	return context
}

// SetStatus Set task status
func SetStatus(uuid string, status string) error {
	args := []string{uuid, "modify", fmt.Sprintf("status:%s", status)}
	cmd := exec.Command("task", args...)
	_, err := cmd.Output()
	return err
}

// LinkToParent Make task child of parent
func LinkToParent(uuid string, parent string) error {
	args := []string{parent, "modify", fmt.Sprintf("depends:%s", uuid)}
	cmd := exec.Command("task", args...)
	_, err := cmd.Output()
	return err
}

// UnlinkFromParent In effect, dedent the task
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

// Add Add task and get uuid
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
