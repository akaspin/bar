package main
import (
	"os"
	"github.com/akaspin/bar/bar/command"
	"github.com/tamtam-im/logx"
)

func main() {
	logx.OnFatal(command.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))

//	cwd, err := os.Getwd()
//	logx.OnFatal(err)
//
//	base := cmd.NewBaseCmd(os.Args, os.Stdin, os.Stdout, os.Stderr)
//	root, err := cmd.NewRootCmd(base)
//	logx.OnFatal(err)
//
//	err = root.Do(cwd, os.Args)
//	logx.OnFatal(err)
}
