package exec

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type Cmd struct {
	cmd *exec.Cmd
}

func (c *Cmd) SetCmd(name string, args ...string) {
	c.SetCmdWithContext(context.Background(), name, args...)
}

func (c *Cmd) SetCmdWithContext(ctx context.Context, name string, args ...string) {
	c.cmd = exec.CommandContext(ctx, name, args...)
}

func (c Cmd) SetEnv(env map[string]string) {
	for k, v := range env {
		c.cmd.Env = append(c.cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
}

func (c Cmd) SetDir(dir string) {
	c.cmd.Dir = dir
}

func (c Cmd) RunWithCombinedOutput() (string, error) {
	ret, err := c.cmd.CombinedOutput()
	return strings.Trim(string(ret), "\n"), err
}

func (c Cmd) RunWithSeparatedOutput() (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	c.cmd.Stdout = &stdout
	c.cmd.Stderr = &stderr
	err := c.cmd.Run()
	return strings.Trim(stdout.String(), "\n"), strings.Trim(stderr.String(), "\n"), err
}

func (c Cmd) RunWithStdout() (string, error) {
	ret, err := c.cmd.Output()
	return strings.Trim(string(ret), "\n"), err
}

func (c Cmd) RunWithStderr() (string, error) {
	var stderr bytes.Buffer
	c.cmd.Stderr = &stderr
	err := c.cmd.Run()
	return strings.Trim(stderr.String(), "\n"), err
}

func (c Cmd) ExitCode() int {
	return c.cmd.ProcessState.ExitCode()
}

func NewCmdRunner(args ...string) *Cmd {
	cmd := new(Cmd)
	if len(args) > 0 {
		cmd.SetCmd(args[0], args[1:]...)
	}
	return cmd
}
