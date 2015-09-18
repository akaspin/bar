package main
import (
	"github.com/akaspin/bar/barctl/cmd"
	"os"
)

func main() {
	err := cmd.Root(os.Args, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		os.Exit(2)
	}
}
