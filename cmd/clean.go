package cmd
import (
	"io"
	"github.com/akaspin/bar/shadow"
	"flag"
	"github.com/tamtam-im/flags"
)

// Implement clean subcommand
// Use "-full" option to get chunks
func CleanCmd(args []string, in io.Reader, out, errOut io.Writer) (err error) {
	var full bool
	fs := flag.NewFlagSet("clean", flag.ExitOnError)
	fs.BoolVar(&full, "full", false, "include chunks to manifest")
	if err = fs.Parse(args[1:]); err != nil {
		return
	}
	flags.New().SetFlagSet(fs).Boot()

	s := &shadow.Shadow{}
	if err = s.FromAny(in, full); err != nil {
		errOut.Write([]byte(err.Error()))
		return
	}
	err = s.Serialize(out)
	return
}
