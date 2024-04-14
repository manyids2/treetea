package task

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
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

func RunCmd() {
	cmd := exec.Command("task")

	// Get json formatted array of tasks, return empty if error
	stdout, stderr, err := Output(cmd)
	fmt.Println("Stdout:\n", string(stdout))
	fmt.Println("Stderr:\n", string(stderr))
	fmt.Println(err)
}
