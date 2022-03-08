package exec

import (
	"os/exec"
	"strings"
)

type Cmd struct {
	cmd *exec.Cmd
}

func (c *Cmd) SetEnv(env ...string) {
	c.cmd.Env = env
}

func (c *Cmd) SetDir(dir string) {
	c.cmd.Dir = dir
}

func (c Cmd) Run() (string, error) {
	ret, err := c.cmd.CombinedOutput()
	return strings.Trim(string(ret), "\n"), err
}

func NewCmdRunner(name string, args ...string) *Cmd {
	return &Cmd{cmd: exec.Command(name, args...)}
}
