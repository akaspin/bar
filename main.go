package main

import (
	"github.com/akaspin/bar/command"
	"github.com/tamtam-im/logx"
	"os"
)

func main() {
	logx.OnFatal(command.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
