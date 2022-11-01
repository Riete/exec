package exec

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"unsafe"
)

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

type Cmd struct {
	cmd *exec.Cmd
}

// SetCmd name is command path, i.e. echo, awk, /path/to/executable, args is command's arguments
func (c *Cmd) SetCmd(name string, args ...string) {
	c.SetCmdWithContext(context.Background(), name, args...)
}

// SetCmdWithContext if ctx.Done(), kill the running command
func (c *Cmd) SetCmdWithContext(ctx context.Context, name string, args ...string) {
	c.cmd = exec.CommandContext(ctx, name, args...)
}

// SetEnv set additional env when running command
func (c Cmd) SetEnv(env map[string]string) {
	for k, v := range env {
		c.cmd.Env = append(c.cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
}

// SetDir set working dir for command when it running
func (c Cmd) SetDir(dir string) {
	c.cmd.Dir = dir
}

// SetStdin set stdin
func (c Cmd) SetStdin(r io.Reader) {
	c.cmd.Stdin = r
}

// SetStdinString pass string as stdin
func (c Cmd) SetStdinString(s string) {
	var b bytes.Buffer
	b.WriteString(s)
	c.cmd.Stdin = &b
}

// RunWithCombinedOutput return stdout and stderr together
func (c Cmd) RunWithCombinedOutput() (string, error) {
	ret, err := c.cmd.CombinedOutput()
	return strings.Trim(BytesToString(ret), "\n"), err
}

// RunWithSeparatedOutput return stdout and stderr separately
func (c Cmd) RunWithSeparatedOutput() (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	c.cmd.Stdout = &stdout
	c.cmd.Stderr = &stderr
	err := c.cmd.Run()
	return strings.Trim(stdout.String(), "\n"), strings.Trim(stderr.String(), "\n"), err
}

// RunWithStdout return stdout only
func (c Cmd) RunWithStdout() (string, error) {
	ret, err := c.cmd.Output()
	return strings.Trim(BytesToString(ret), "\n"), err
}

// RunWithStderr return stderr only
func (c Cmd) RunWithStderr() (string, error) {
	var stderr bytes.Buffer
	c.cmd.Stderr = &stderr
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

func NewCmdRunner(args ...string) *Cmd {
	cmd := new(Cmd)
	if len(args) > 0 {
		cmd.SetCmd(args[0], args[1:]...)
	}
	return cmd
}

// Pipe  command1  | command2 | command3 | ...
func Pipe(cmd ...*Cmd) *Cmd {
	return cmd[0].Pipe(cmd[1:]...)
}
