package cmd
import "github.com/akaspin/bar/proto/manifest"


/*
Import spec from bard and populate manifests
*/
type SpecImportCmd struct  {
	*BaseSubCommand

	endpoint string
	useGit bool
	chunkSize int64
	pool int
}

func NewSpecImportCmd(s *BaseSubCommand) SubCommand  {
	c := &SpecExportCmd{BaseSubCommand: s}
	s.FS.StringVar(&c.endpoint, "endpoint", "http://localhost:3000/v1",
		"bard endpoint")
	s.FS.BoolVar(&c.useGit, "git", false, "use git infrastructure")
	s.FS.Int64Var(&c.chunkSize, "chunk", manifest.CHUNK_SIZE, "preferred chunk size")
	s.FS.IntVar(&c.pool, "pool", 16, "pool sizes")
	return c
}

func (c *SpecImportCmd) Do() (err error) {
	return
}

//c := &SpecExportCmd{BaseSubCommand: s}
//s.FS.StringVar(&c.endpoint, "endpoint", "http://localhost:3000/v1",
//"bard endpoint")
//s.FS.BoolVar(&c.useGit, "git", false, "use git infrastructure")
//s.FS.Int64Var(&c.chunkSize, "chunk", manifest.CHUNK_SIZE, "preferred chunk size")
//s.FS.BoolVar(&c.upload, "upload", false, "upload spec to bard and print URL")
//s.FS.IntVar(&c.pool, "pool", 16, "pool sizes")
//return c