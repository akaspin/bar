package cmd
import (
	"fmt"
	"io"
)

type SneakCmd struct {
	*Base
}

func NewSneakCmd(s *Base) SubCommand {
	c := &SneakCmd{Base: s}

	return c
}

func (c *SneakCmd) Do(args []string) (err error) {
	fmt.Fprintln(c.Stderr, args)
	io.Copy(c.Stdout, c.Stdin)
	return
}