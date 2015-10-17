package main
import (
	"github.com/akaspin/bar/barc/cmd"
	"os"
	"github.com/tamtam-im/logx"
)

func main() {
	cwd, err := os.Getwd()
	logx.OnFatal(err)

	base := cmd.NewBaseCmd(os.Args, os.Stdin, os.Stdout, os.Stderr)
	root, err := cmd.NewRootCmd(base)
	logx.OnFatal(err)

	err = root.Do(cwd, os.Args)
	logx.OnFatal(err)
}
