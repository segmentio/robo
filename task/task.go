package task

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// Example usage.
type Example struct {
	Description string
	Command     string
}

// Task definition.
type Task struct {
	LookupPath string
	Name       string `yaml:"-"`
	Summary    string
	Command    string
	Script     string
	Exec       string
	Usage      string
	Examples   []*Example
	Trace      bool
}

// Run the task with xtrace shell option
func (t *Task) RunTrace(args []string) error {
	t.Trace = true
	return t.Run(args)
}

// Run the task with `args`.
func (t *Task) Run(args []string) error {
	shell := os.Getenv("ROBO_SHELL")
	if shell == "" {
		shell = "sh"
	}

	if t.Exec != "" {
		return t.RunExec(args)
	}

	if t.Script != "" {
		return t.RunScript(shell, args)
	}

	if t.Command != "" {
		return t.RunCommand(shell, args)
	}

	return fmt.Errorf("nothing to run (add script, command, or exec key)")
}

// RunScript runs the target shell `script` file.
func (t *Task) RunScript(shell string, args []string) error {
	path := filepath.Join(t.LookupPath, t.Script)
	args = append([]string{path}, args...)
	if t.Trace {
		args = append([]string{"-x"}, args...)
	}
	cmd := exec.Command(shell, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunCommand runs the `command` via the shell.
func (t *Task) RunCommand(shell string, args []string) error {
	args = append([]string{"-c", t.Command, "sh"}, args...)
	if t.Trace {
		args = append([]string{"-x"}, args...)
	}
	cmd := exec.Command(shell, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunExec runs the `exec` command.
func (t *Task) RunExec(args []string) error {
	fields := strings.Fields(t.Exec)
	bin := fields[0]

	path, err := exec.LookPath(bin)
	if err != nil {
		return err
	}

	args = append(fields, args...)
	return syscall.Exec(path, args, os.Environ())
}
