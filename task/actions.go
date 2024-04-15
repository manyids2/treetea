package task

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"strings"
)

/* ----------------------------------------------------------------
---                    General                                  ---
---------------------------------------------------------------- */

func Output(c *exec.Cmd) ([]byte, []byte, error) {
	if c.Stdout != nil {
		return nil, nil, errors.New("exec: Stdout already set")
	}
	var stdout bytes.Buffer
	c.Stdout = &stdout

	if c.Stderr != nil {
		return nil, nil, errors.New("exec: Stderr already set")
	}
	var stderr bytes.Buffer
	c.Stderr = &stderr

	err := c.Run()
	return stdout.Bytes(), stderr.Bytes(), err
}

func GetEnv() (taskdata, taskrc string) {
	taskdata = os.Getenv("TASKDATA")
	taskrc = os.Getenv("TASKRC")
	return taskdata, taskrc
}

func FileExists(taskrc string) bool {
	if _, err := os.Stat(taskrc); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

// https://stackoverflow.com/questions/66643946/how-to-remove-duplicates-strings-or-int-from-slice-in-go
func RemoveDuplicateProjects(sliceList []Project) []Project {
	allKeys := make(map[string]bool)
	list := []Project{}
	for _, item := range sliceList {
		if _, value := allKeys[item.Name]; !value {
			allKeys[item.Name] = true
			list = append(list, item)
		}
	}
	return list
}

/* ----------------------------------------------------------------
---                    Projects                                 ---
---------------------------------------------------------------- */

func ReadProjects(path string) ([]Project, error) {
	out, _ := os.ReadFile(path)
	var projects []Project
	if err := json.Unmarshal(out, &projects); err != nil {
		return []Project{}, err
	}
	return projects, nil
}

func TaskProjects() ([]Project, error) {
	cmd := exec.Command("task", "_projects")
	stdout, _, err := Output(cmd) // Ignore stderr
	if err != nil {
		return []Project{}, err
	}
	names := strings.Split(string(stdout), "\n")
	names = names[:len(names)-1] // Remove last newline
	projects := make([]Project, len(names))
	parents := make([]string, len(names))
	for i, n := range names {
		projects[i] = Project{Name: n, Subprojects: []string{}}
		parts := strings.Split(n, ".")
		if len(parts) > 1 {
			parents[i] = strings.Join(parts[:len(parts)-1], ".")
		}
	}
	for i, p := range parents {
		if p == "" {
			continue
		}
		// Have to loop through list to get index
		for j, n := range projects {
			if n.Name == p {
				projects[j].Subprojects = append(n.Subprojects, names[i])
			}
		}
	}
	return projects, nil
}

func SaveProjects(projects []Project, path string) error {
	bytes, err := json.Marshal(projects)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}
	return err
}

func MergeProjects(projects_a, projects_b []Project) []Project {
	// TODO: Basics work, but need to think about overrides
	projects := append(projects_a, projects_b...)
	projects = RemoveDuplicateProjects(projects)
	return projects
}

func ImportProjects(projects_from, projects_to []Project, parent string) []Project {
	// TODO: Basics work, but need to think about overrides
	for i, p := range projects_to {
		if p.Name != parent {
			continue
		} else {
			for j := range projects_from {
				name := projects_to[i].Name + "." + projects_from[j].Name
				projects_from[j].Name = name
				projects_to[i].Subprojects = append(projects_to[i].Subprojects, name)
			}
		}
	}
	projects := append(projects_from, projects_to...)
	return projects
}

/* ----------------------------------------------------------------
---                    Contexts                                 ---
---------------------------------------------------------------- */

func ReadContexts(path string) (contexts []Context, err error) {
	out, _ := os.ReadFile(path)
	if err = json.Unmarshal(out, &contexts); err != nil {
		return []Context{}, err
	}
	return contexts, nil
}

/* ----------------------------------------------------------------
---                     Tasks                                   ---
---------------------------------------------------------------- */

func List(filters string) (tasks []Task, err error) {
	args := strings.Split(filters, " ")
	args = append(args, "export")
	cmd := exec.Command("task", args...)
	stdout, _, err := Output(cmd) // Ignore stderr
	if err := json.Unmarshal(stdout, &tasks); err != nil {
		return []Task{}, err
	}
	return tasks, nil
}
