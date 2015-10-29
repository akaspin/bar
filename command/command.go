package command

import (
	"github.com/spf13/cobra"
	"io"
)

// Command environment
type Environment struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// Common bar options
type CommonOptions struct {
	LoggingLevel string
	WD           string

	Endpoint   string
	ChunkSize  int64
	PoolSize   int
	BufferSize int
}

type Command interface {

	// Init flags
	Init(c *cobra.Command)

	// Run command
	Run(args ...string) (err error)
}

// Attach command to backend
func Attach(c Command, e *Environment, cmds ...*cobra.Command) (res *cobra.Command) {
	res = &cobra.Command{}
	res.SetOutput(e.Stderr)
	c.Init(res)
	if len(cmds) == 0 {
		res.RunE = func(cc *cobra.Command, args []string) (err error) {
			return c.Run(args...)
		}
	}
	res.AddCommand(cmds...)
	return
}
