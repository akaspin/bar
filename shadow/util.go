package shadow
import (
	"io"
	"bytes"
)


func Detect(in io.Reader) (r io.Reader, isShadow bool, err error) {
	var n int
	buf := make([]byte, len([]byte(SHADOW_HEADER)))

	// check header signature
	if n, err = in.Read(buf); err != nil {
		return
	}

	r = io.MultiReader(bytes.NewBuffer(buf[:n]), in)
	if string(buf) == SHADOW_HEADER {
		isShadow = true
	}
	return

}
