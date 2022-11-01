package exec

import (
	"testing"
)

func TestNewCmdRunnerCombinedOutput(t *testing.T) {
	r := NewCmdRunner()
	r.SetCmd("/tmp/1.sh")
	r.SetEnv(map[string]string{"AAA": "a1111111111a"})
	t.Log(r.RunWithCombinedOutput())
	t.Log(r.ExitCode())
}

func TestNewCmdRunnerSeperatedOutput(t *testing.T) {
	r := NewCmdRunner()
	r.SetCmd("/tmp/1.sh")
	r.SetEnv(map[string]string{"AAA": "a1111111111a"})
	t.Log(r.RunWithSeparatedOutput())
	t.Log(r.ExitCode())
}

func TestNewCmdRunnerStdout(t *testing.T) {
	r := NewCmdRunner()
	r.SetCmd("/tmp/1.sh")
	r.SetEnv(map[string]string{"AAA": "a1111111111a"})
	t.Log(r.RunWithStdout())
	t.Log(r.ExitCode())
}

func TestNewCmdRunnerStderr(t *testing.T) {
	r := NewCmdRunner()
	r.SetCmd("/tmp/1.sh")
	r.SetEnv(map[string]string{"AAA": "a1111111111a"})
	t.Log(r.RunWithStderr())
	t.Log(r.ExitCode())
}

func TestNewCmdRunnerStdin(t *testing.T) {
	r := NewCmdRunner()
	r.SetCmd("awk", `NR==2{print}`)
	r.SetStdinString("1 2 3\n2 3 4\n 1 3 4")
	t.Log(r.RunWithCombinedOutput())
}

func TestNewCmdRunnerPipe(t *testing.T) {
	n1 := NewCmdRunner("echo", "1 2 3\n2 3 4\n 1 3 4")
	n2 := NewCmdRunner("awk", `{print $1}`)
	n3 := NewCmdRunner("awk", `NR==2{print $1}`)

	// f := Pipe(n1, n2, n3)
	// t.Log(f.RunWithSeparatedOutput())

	f := n1.Pipe(n2, n3)
	t.Log(f.RunWithSeparatedOutput())

}
