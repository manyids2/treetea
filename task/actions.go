package task

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"strings"
)

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

func ReadProjects(path string) ([]Project, error) {
	out, _ := os.ReadFile(path)
	var projects []Project
	if err := json.Unmarshal(out, &projects); err != nil {
		return []Project{}, err
	}
	return projects, nil
}

func ReadContexts(path string) (contexts []Context, err error) {
	out, _ := os.ReadFile(path)
	if err = json.Unmarshal(out, &contexts); err != nil {
		return []Context{}, err
	}
	return contexts, nil
}

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
