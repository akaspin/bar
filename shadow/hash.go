package shadow
import (
	"io"
	"bytes"
	"fmt"
	"gopkg.in/bufio.v1"
	"golang.org/x/crypto/sha3"
)

const (
	CHUNK_SIZE = 1024 * 1024   // 1M
	SHADOW_HEADER = "BAR:SHADOW"
	VERSION = "0.1.0"
)

type Chunk struct  {
	ID []byte
	Size int64
	Offset int64
}

type Shadow struct {
	IsFromShadow bool
	Version string
	ID []byte
	Size int64
	Chunks []Chunk
}

func (s Shadow) String() (res string) {
	data := []byte{}
	buf := bytes.NewBuffer(data)
	(&s).Serialize(buf)
	res = string(buf.Bytes())
	return
}

// write initialized manifest to specific stream
func (s *Shadow) Serialize(out io.Writer) (err error) {
	// Header
	if _, err = out.Write([]byte(SHADOW_HEADER)); err != nil {
		return
	}
	body := fmt.Sprintf("\n\nversion %s\nid %x\nsize %d\n\n",
		VERSION, s.ID, s.Size)
	_, err = out.Write([]byte(body))
	return
}

// Initialise from any source
func (s *Shadow) FromAny(in io.Reader) (err error) {
	var n int
	maybeHeader := make([]byte, len([]byte(SHADOW_HEADER)))

	// check header signature
	if n, err = in.Read(maybeHeader); err != nil {
		return
	}

	r := io.MultiReader(bytes.NewBuffer(maybeHeader[:n]), in)
	if string(maybeHeader) == SHADOW_HEADER {
		err = s.FromManifest(r)
	} else {

		err = s.FromBlob(r)
	}
	return
}

// Parse manifest
func (s *Shadow) FromManifest(in io.Reader) (err error) {
	var data []byte

	br := bufio.NewReader(in)

	// Check header
	if data, _, err = br.ReadLine(); err != nil {
		return
	}
	if string(data) != SHADOW_HEADER {
		err = fmt.Errorf("bad shadow header %s", string(data))
		return
	}
	// check delimiter
	if data, _, err = br.ReadLine(); err != nil {
		return
	}
	if len(data) > 0 {
		err = fmt.Errorf("bad delimiter %s", string(data))
	}

	// Read version
	if data, _, err = br.ReadLine(); err != nil {
		return
	}
	if _, err = fmt.Sscanf(string(data), "version %s", &s.Version); err != nil {
		return
	}

	// TODO: semver check

	// Read id
	if data, _, err = br.ReadLine(); err != nil {
		return
	}
	if _, err = fmt.Sscanf(string(data), "id %x", &s.ID); err != nil {
		return
	}

	// Read size
	if data, _, err = br.ReadLine(); err != nil {
		return
	}
	if _, err = fmt.Sscanf(string(data), "size %d", &s.Size); err != nil {
		return
	}

	// check ending
	if data, _, err = br.ReadLine(); err != nil {
		return
	}
	if len(data) > 0 {
		err = fmt.Errorf("bad ending %s", string(data))
	}

	return
}

// Initialize manifest from BLOB
func (s *Shadow) FromBlob(in io.Reader) (err error) {
	hasher := sha3.New256()
	var written int64
	for {
		written, err = io.CopyN(hasher, in, CHUNK_SIZE)
		s.Size += written
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return
		}
	}
	s.ID = hasher.Sum(nil)
	s.Version = VERSION
	return
}
