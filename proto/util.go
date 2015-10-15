package proto
import (
	"io"
	"bytes"
)


// Peek blob kind from given reader
func PeekManifest(in io.Reader) (r io.Reader, isShadow bool, err error) {
	var n int
	buf := make([]byte, len([]byte(MANIFEST_HEADER)))

	// check header signature
	if n, err = in.Read(buf); err != nil {
		return
	}

	r = io.MultiReader(bytes.NewBuffer(buf[:n]), in)
	if string(buf) == MANIFEST_HEADER {
		isShadow = true
	}
	return

}
