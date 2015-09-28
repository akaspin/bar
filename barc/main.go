package main
import (
	"github.com/akaspin/bar/barc/cmd"
	"os"
	"fmt"
)

func main() {
	err := cmd.Root(os.Args, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}
