package task

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
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

// View View functional for task
func View(level int, status string, description string) string {
	return fmt.Sprintf("%s%s %s",
		strings.Repeat(" ", level*2), // Indent
		MARKERS["icon"][status],      // Status marker
		description)                  // Description
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
	if err != nil {
		log.Fatal("Failed adding task: ", err)
	}
	return err
}

// UnlinkFromParent In effect, dedent the task
func UnlinkFromParent(uuid string, parent string, grandparent string, depends []string) error {
	// Remove all dependencies of parent
	args := []string{parent, "modify", "depends:"}
	cmd := exec.Command("task", args...)
	_, err := cmd.Output()
	if err != nil {
		log.Fatal("Failed removing dependencies: ", err, uuid, parent, depends)
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
			log.Fatal("Failed adding dep: ", err, uuid, parent, v, depends)
		}
	}

	// Add current to grandparent
	if grandparent != "" {
		args := []string{grandparent, "modify", fmt.Sprintf("depends:%s", uuid)}
		cmd := exec.Command("task", args...)
		_, err := cmd.Output()
		if err != nil {
			log.Fatal("Failed adding back uuid: ", err, uuid, parent, uuid, depends)
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
		_, err := cmd.Output()
		if err != nil {
			return uuid, err
		}
	}
	return uuid, err
}
