package main
import (
	"os"
	"github.com/akaspin/bar/command"
	"github.com/tamtam-im/logx"
)

func main() {
	logx.OnFatal(command.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
