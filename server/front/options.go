package front

import "github.com/akaspin/bar/proto"

type Options struct {
	Info   *proto.ServerInfo
	Bind   string
	BinDir string
}
