package exec

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Cmd struct {
	cmd *exec.Cmd
}

// SetCommand name is command path, i.e. echo, awk, /path/to/executable, args is command's arguments
func (c *Cmd) SetCommand(name string, args ...string) {
	c.SetCommandContext(context.Background(), name, args...)
}

// SetCommandContext if ctx.Done(), kill the running command
func (c *Cmd) SetCommandContext(ctx context.Context, name string, args ...string) {
	c.cmd = exec.CommandContext(ctx, name, args...)
}

// SetEnv set additional env when running command, overwrite control duplicate key behave
func (c Cmd) SetEnv(env map[string]string, overwrite bool) {
	for k, v := range env {
		c.cmd.Env = append(c.cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	if overwrite {
		c.cmd.Env = append(os.Environ(), c.cmd.Env...)
	} else {
		c.cmd.Env = append(c.cmd.Env, os.Environ()...)

	}
}

// SetWorkingDir set working dir for command when it is running
func (c Cmd) SetWorkingDir(dir string) {
	c.cmd.Dir = dir
}

// SetStdin set stdin
func (c Cmd) SetStdin(r io.Reader) {
	c.cmd.Stdin = r
}

// SetStringAsStdin pass string as stdin
func (c Cmd) SetStringAsStdin(s string) {
	c.cmd.Stdin = bytes.NewBufferString(s)
}

// RunWithCombinedOutput return stdout and stderr together
func (c Cmd) RunWithCombinedOutput() (string, error) {
	out := new(bytes.Buffer)
	c.cmd.Stdout = out
	c.cmd.Stderr = out
	err := c.cmd.Run()
	return strings.Trim(out.String(), "\n"), err
}

// RunWithSeparatedOutput return stdout and stderr separately
func (c Cmd) RunWithSeparatedOutput() (string, string, error) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	c.cmd.Stdout = stdout
	c.cmd.Stderr = stderr
	err := c.cmd.Run()
	return strings.Trim(stdout.String(), "\n"), strings.Trim(stderr.String(), "\n"), err
}

// RunWithStdout return stdout only
func (c Cmd) RunWithStdout() (string, error) {
	stdout := new(bytes.Buffer)
	c.cmd.Stdout = stdout
	err := c.cmd.Run()
	return strings.Trim(stdout.String(), "\n"), err
}

// RunWithStderr return stderr only
func (c Cmd) RunWithStderr() (string, error) {
	stderr := new(bytes.Buffer)
	c.cmd.Stderr = stderr
	err := c.cmd.Run()
	return strings.Trim(stderr.String(), "\n"), err
}

// ExitCode return exit code , return 0 means success
func (c Cmd) ExitCode() int {
	return c.cmd.ProcessState.ExitCode()
}

// pipe  command1  | command2
func (c Cmd) pipe(next *Cmd) *Cmd {
	r, w := io.Pipe()
	c.cmd.Stdout = w
	next.cmd.Stdin = r
	go func() {
		_ = w.CloseWithError(c.cmd.Run())
	}()
	return next
}

// Pipe  command1  | command2 | command3 | ...
func (c Cmd) Pipe(cmd ...*Cmd) *Cmd {
	next := c.pipe(cmd[0])
	for _, i := range cmd[1:] {
		next = next.pipe(i)
	}
	return next
}

func New(options ...Option) *Cmd {
	cmd := new(Cmd)
	for _, option := range options {
		option(cmd)
	}
	return cmd
}

// Pipe  command1  | command2 | command3 | ...
func Pipe(cmd ...*Cmd) *Cmd {
	return cmd[0].Pipe(cmd[1:]...)
}
