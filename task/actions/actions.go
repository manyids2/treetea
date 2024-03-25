package actions

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	tw "github.com/manyids2/tasktea/task"
)

// Paths Get paths to taskdata and taskrc
func Paths() (taskdata, taskrc string) {
	taskdata, _ = os.LookupEnv("TASKDATA")
	taskrc, _ = os.LookupEnv("TASKRC")
	return taskdata, taskrc
}

type Config map[string]string

// ParseRc Parse taskrc and return map of keys and values
func ParseRc(path string) (cfg Config, err error) {
	readFile, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	r := regexp.MustCompile(`(?P<Key>[0-9a-zA-Z_\-\.]+)=(?P<Val>.*)`)

	// First load all into a map
	cfg = map[string]string{}
	for fileScanner.Scan() {
		if len(fileScanner.Text()) == 0 {
			continue
		}
		if fileScanner.Text()[0] == '#' {
			continue
		}
		match := r.FindStringSubmatch(fileScanner.Text())
		if len(match) != 3 {
			continue
		}
		cfg[match[1]] = match[2]
	}
	readFile.Close()

	// Parse into TaskRC struct

	return cfg, nil
}

// RcExtract Extract contexts and reports which are 2 level nested, and rest are kept at 2 levels
// if one level, the second key is ""
// if no dots, just query the map, do not use this function
func RcExtract(key string, cfg Config) (values map[string]map[string]string) {
	values = map[string]map[string]string{}
	for k, v := range cfg {
		parts := strings.Split(k, ".")
		if (parts[0] == key) && (len(parts) > 1) {
			_, ok := values[parts[1]]
			if !ok {
				values[parts[1]] = map[string]string{}
			}
			subkey := ""
			if len(parts) > 2 {
				subkey = parts[2]
			}
			if len(parts) > 3 {
				subkey = strings.Join(parts[2:], ".")
			}
			values[parts[1]][subkey] = v
		}
	}
	return values
}

// List Get list of tasks given filters
//
// `task [filters] export`
func List(filters []string) (tasks_list []tw.Task, err error) {
	filters = append(filters, "export")
	out, err := exec.Command("task", filters...).Output()
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(out, &tasks_list); err != nil {
		return nil, err
	}
	return tasks_list, nil
}

// List Get list of active tasks
//
// `task export active`
func Active() (active []string, err error) {
	var tasks_list []tw.Task
	out, err := exec.Command("task", "export", "active").Output()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(out, &tasks_list); err != nil {
		return nil, err
	}

	active = []string{}
	for _, v := range tasks_list {
		active = append(active, v.UUID)
	}
	return active, nil
}

// Contexts Get all contexts
//
// `task _context`
func Contexts() (contexts []string, err error) {
	out, err := exec.Command("task", "_context").Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\n")
	contexts = lines[:len(lines)-1] // Remove last empty line
	contexts = append(contexts, "none")
	return contexts, nil
}

// Context Parse current context to get read and write filters
//
//	Assumes 1st line is context, 3rd is read filter and 4th is write filter
//	- In case of "none", expected output is `No context is currently applied.`,
//	  so it returns the correct results
//
// `task context show`
func Context() (context string, read_filters string, write_filters string, err error) {
	context = "none"
	out, err := exec.Command("task", "context", "show").Output()
	if err != nil {
		return context, "", "", err
	}
	lines := strings.Split(string(out), "\n")

	r := regexp.MustCompile(`Context '(?P<Context>[a-zA-Z]+)' with`)
	match := r.FindStringSubmatch(lines[0])
	context = match[1]

	r = regexp.MustCompile(`\* read filter: '(?P<ReadFilter>.*)'`)
	match = r.FindStringSubmatch(lines[2])
	read_filters = match[1]

	r = regexp.MustCompile(`\* write filter: '(?P<WriteFilter>.*)'`)
	match = r.FindStringSubmatch(lines[3])
	write_filters = match[1]

	return context, read_filters, write_filters, nil
}

// Projects Get all Projects
//
// `task _projects`
func Projects() (projects []string, err error) {
	out, err := exec.Command("task", "_projects").Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\n")
	projects = lines[:len(lines)-1]
	return projects, nil
}

// Projects Get all Projects
//
// `task context [context]`
func SetContext(context string) (err error) {
	_, err = exec.Command("task", "context", context).Output()
	return err
}

// SetStatus Set task status
//
// `task [uuid] modify status:[status]`
func SetStatus(uuid string, status string) (err error) {
	args := []string{uuid, "modify", fmt.Sprintf("status:%s", status)}
	_, err = exec.Command("task", args...).Output()
	return err
}

// SetActive Start/stop task
//
// `task [uuid] [start/stop]`
func SetActive(uuid string, startstop string) (err error) {
	cmd := exec.Command("task", uuid, startstop)
	_, err = cmd.Output()
	return err
}

// SetParent Make task child of parent
//
// `task [uuid] modify depends:[parent]`
func SetParent(uuid string, parent string) (err error) {
	args := []string{parent, "modify", fmt.Sprintf("depends:%s", uuid)}
	_, err = exec.Command("task", args...).Output()
	return err
}

// UnlinkFromParent In effect, dedent the task
//
// `task [uuid] modify depends:`
// `task [parent] modify depends:[all children other than uuid]`
// `task [grandparent] modify depends:[uuid]`
func SetGrandparentAsParent(uuid string, parent string, grandparent string, depends []string) (err error) {
	// Remove all dependencies of parent
	args := []string{parent, "modify", "depends:"}
	_, err = exec.Command("task", args...).Output()
	if err != nil {
		return err
	}

	// Add back non uuid depends
	for _, v := range depends {
		if v == uuid {
			continue
		}
		args = []string{parent, "modify", fmt.Sprintf("depends:%s", v)}
		_, err = exec.Command("task", args...).Output()
		if err != nil {
			return err
		}
	}

	// Add current to grandparent
	if grandparent != "" {
		args = []string{grandparent, "modify", fmt.Sprintf("depends:%s", uuid)}
		_, err := exec.Command("task", args...).Output()
		if err != nil {
			return err
		}
	}
	return err
}

// Add Add task, get uuid and update parent's children
//
// `task add [tags] -- [description]`
// `task [parent] modify depends:[uuid]`
func Add(filters []string, description string, parent string) (uuid string, err error) {
	// Removing things like `-tag`, `/.../` which dont make sense
	// TODO: Maybe use regex, read spec to find out also whats up
	tags := []string{}
	for _, f := range filters {
		if len(f) == 0 {
			continue
		}
		if (string(f[0]) != "-") && (string(f[0]) != "/") {
			tags = append(tags, f)
		}
	}

	// Add task
	args := []string{"add"}
	args = append(args, tags...)
	args = append(args, "--", description)
	out, err := exec.Command("task", args...).Output()
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(out), "\n")

	// Extract assigned UUID
	r := regexp.MustCompile(`Created task '(?P<UUID>.*)'`)
	for _, v := range lines {
		match := r.FindStringSubmatch(v)
		if len(match) > 0 {
			uuid = match[1]
		}
	}

	// Update parent's dependencies
	if (parent != "") && (uuid != "") {
		args = []string{parent, "modify", fmt.Sprintf("depends:%s", uuid)}
		_, err = exec.Command("task", args...).Output()
	}
	return uuid, err
}

// Modify task description, i.e. edit description
//
// `task [uuid] modify -- [description]`
func ModifyDescription(uuid string, description string) (err error) {
	args := []string{uuid, "modify", "--", description}
	_, err = exec.Command("task", args...).Output()
	return err
}

// Modify task with arguments
//
// `task [uuid] modify [args]`
func Modify(uuid string, args []string) (err error) {
	out := []string{uuid, "modify"}
	out = append(out, args...)
	_, err = exec.Command("task", out...).Output()
	return err
}

// ModifyBatch Modify multiple tasks with same arguments
//
// `task [uuids] modify [args]`
func ModifyBatch(uuids []string, args []string) (err error) {
	out := append(uuids, "modify")
	out = append(out, args...)
	_, err = exec.Command("task", out...).Output()
	return err
}

// Delete task, turn of confirmation
//
// `task [uuid] rc.confirmation=off delete`
func Delete(uuid string) (err error) {
	args := []string{uuid, "rc.confirmation=off", "delete"}
	_, err = exec.Command("task", args...).Output()
	return err
}
