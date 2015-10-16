package main
import (
	"github.com/akaspin/bar/barc/cmd"
	"os"
	"fmt"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	base := cmd.NewBaseCmd(os.Args, os.Stdin, os.Stdout, os.Stderr)
	root, err := cmd.NewRootCmd(base)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	err = root.Do(cwd, os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}
