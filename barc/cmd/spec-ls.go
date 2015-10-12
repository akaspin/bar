package cmd
import "github.com/akaspin/bar/proto/manifest"

type SpecLsCmd struct {
	*BaseSubCommand

	endpoint string
	useGit bool
	chunkSize int64
	pool int

}

func NewSpecLsCmd(s *BaseSubCommand) SubCommand  {
	c := &SpecImportCmd{BaseSubCommand: s}
	s.FS.StringVar(&c.endpoint, "endpoint", "http://localhost:3000/v1",
		"bard endpoint")
	s.FS.BoolVar(&c.useGit, "git", false, "use git infrastructure")
	s.FS.Int64Var(&c.chunkSize, "chunk", manifest.CHUNK_SIZE, "preferred chunk size")
	s.FS.IntVar(&c.pool, "pool", 16, "pool sizes")
	return c
}

func (c *SpecLsCmd) Do() (err error) {
	return
}