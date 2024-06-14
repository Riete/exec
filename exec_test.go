package exec

import (
	"testing"
)

func TestCombinedOutput(t *testing.T) {
	r := New()
	r.SetCommand("/tmp/1.sh")
	r.SetEnv(map[string]string{"AAA": "a1111111111a"}, true)
	t.Log(r.RunWithCombinedOutput())
	t.Log(r.ExitCode())
}

func TestSeperatedOutput(t *testing.T) {
	r := New()
	r.SetCommand("/tmp/1.sh")
	r.SetEnv(map[string]string{"AAA": "a1111111111a"}, true)
	t.Log(r.RunWithSeparatedOutput())
	t.Log(r.ExitCode())
}

func TestStdout(t *testing.T) {
	r := New()
	r.SetCommand("/tmp/1.sh")
	r.SetEnv(map[string]string{"AAA": "a1111111111a"}, true)
	t.Log(r.RunWithStdout())
	t.Log(r.ExitCode())
}

func TestStderr(t *testing.T) {
	r := New()
	r.SetCommand("/tmp/1.sh")
	r.SetEnv(map[string]string{"AAA": "a1111111111a"}, true)
	t.Log(r.RunWithStderr())
	t.Log(r.ExitCode())
}

func TestStdin(t *testing.T) {
	r := New()
	r.SetCommand("awk", `NR==2{print}`)
	r.SetStringAsStdin("1 2 3\n2 3 4\n 1 3 4")
	t.Log(r.RunWithCombinedOutput())
}

func TestPipe(t *testing.T) {
	n1 := New(WithCommand("echo", "1 2 3\n2 3 4\n 1 3 4"))
	n2 := New(WithCommand("awk", `{print $1}`))
	n3 := New(WithCommand("awk", `NR==2{print $1}`))

	f := Pipe(n1, n2, n3)
	t.Log(f.RunWithSeparatedOutput())

	// f := n1.Pipe(n2, n3)
	// t.Log(f.RunWithSeparatedOutput())

}
