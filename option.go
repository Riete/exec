package exec

import (
	"context"
	"io"
)

type Option func(*Cmd)

func WithCommand(name string, args ...string) Option {
	return func(c *Cmd) {
		c.SetCommand(name, args...)
	}
}

func WithCommandContext(ctx context.Context, name string, args ...string) Option {
	return func(c *Cmd) {
		c.SetCommandContext(ctx, name, args...)
	}
}

func WithEnv(env map[string]string, overwrite bool) Option {
	return func(c *Cmd) {
		c.SetEnv(env, overwrite)
	}
}

func WithWorkingDir(dir string) Option {
	return func(c *Cmd) {
		c.SetWorkingDir(dir)
	}
}

func WithStdin(r io.Reader) Option {
	return func(c *Cmd) {
		c.SetStdin(r)
	}
}

func WithStringAsStdin(s string) Option {
	return func(c *Cmd) {
		c.SetStringAsStdin(s)
	}
}
