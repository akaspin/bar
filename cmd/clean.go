package cmd
import (
	"io"
	"github.com/akaspin/bar/shadow"
	"flag"
)

type CleanSubCommand struct {
	endpoint string
	full bool
	strict bool
}

func (c *CleanSubCommand) FS(fs *flag.FlagSet) {
	fs.BoolVar(&c.full, "full", false, "include chunks into manifest")
	fs.BoolVar(&c.strict, "strict", false, "check blobs on server")
	addEndpointToFS(fs, &c.endpoint)
}

func (c *CleanSubCommand) Do(in io.Reader, out, errOut io.Writer) (err error) {

	s := &shadow.Shadow{}
	if err = s.FromAny(in, c.full); err != nil {
		errOut.Write([]byte(err.Error()))
		return
	}
	err = s.Serialize(out)
	return
}
