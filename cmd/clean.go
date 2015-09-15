package cmd
import (
	"io"
	"github.com/akaspin/bar/shadow"
)

// Implement clean subcommand
func CleanCmd(args []string, in io.Reader, out, errOut io.Writer) (err error) {
	s := &shadow.Shadow{}
	if err = s.FromAny(in); err != nil {
		errOut.Write([]byte(err.Error()))
		return
	}
	err = s.Serialize(out)
	return
}
