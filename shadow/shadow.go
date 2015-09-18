package shadow
import (
	"io"
	"bytes"
	"fmt"
	"gopkg.in/bufio.v1"
	"golang.org/x/crypto/sha3"
	"hash"
	"os"
	"strings"
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

func (c Chunk) String() string {
	return fmt.Sprintf("id %x\nsize %d\noffset %d\n\n", c.ID, c.Size, c.Offset)
}

type Shadow struct {
	IsFromShadow bool
	Version string
	ID []byte
	Size int64
	Chunks []Chunk
}

func NewShadowFromFile(filename string, full bool, chunkSize int64) (res *Shadow, err error) {
	r, err := os.Open(filename)
	if err != nil {
		return
	}
	defer r.Close()
	res = &Shadow{}
	err = res.FromAny(r, full, chunkSize)
	return
}

func (s Shadow) String() (res string) {
	data := []byte{}
	buf := bytes.NewBuffer(data)
	(&s).Serialize(buf)
	res = string(buf.Bytes())
	return
}

func (s *Shadow) HasChunks() bool {
	return len(s.Chunks) > 0
}

// write initialized manifest to specific stream
func (s *Shadow) Serialize(out io.Writer) (err error) {
	if _, err = out.Write([]byte(SHADOW_HEADER)); err != nil {
		return
	}
	body := fmt.Sprintf("\n\nversion %s\nid %x\nsize %d\n\n",
		VERSION, s.ID, s.Size)
	if _, err = out.Write([]byte(body)); err != nil {
		return
	}

	if s.HasChunks() {
		if _, err = out.Write([]byte("\n")); err != nil {
			return
		}
		for _, chunk := range s.Chunks {
			if _, err = out.Write([]byte(chunk.String())); err != nil {
				return
			}
		}
	}
	return
}

// Initialise from any source
func (s *Shadow) FromAny(in io.Reader, full bool, chunkSize int64) (err error) {
	var n int
	maybeHeader := make([]byte, len([]byte(SHADOW_HEADER)))

	// check header signature
	if n, err = in.Read(maybeHeader); err != nil {
		return
	}

	r := io.MultiReader(bytes.NewBuffer(maybeHeader[:n]), in)
	if string(maybeHeader) == SHADOW_HEADER {
		err = s.FromManifest(r, full)
	} else {
		err = s.FromBlob(r, full, chunkSize)
	}
	return
}

// Parse manifest
func (s *Shadow) FromManifest(in io.Reader, full bool) (err error) {
	s.IsFromShadow = true

	var buf []byte
	var data string

	br := bufio.NewReader(in)

	// Check header
	if data, err = s.nextLine(br, buf); err != nil {
		return
	}
	if data != SHADOW_HEADER {
		err = fmt.Errorf("bad shadow header %s", data)
		return
	}
	// check delimiter
	if data, err = s.nextLine(br, buf); err != nil {
		return
	}
	if len(buf) > 0 {
		err = fmt.Errorf("bad delimiter %s", data)
	}

	// Read version
	if data, err = s.nextLine(br, buf); err != nil {
		return
	}
	if _, err = fmt.Sscanf(data, "version %s", &s.Version); err != nil {
		return
	}

	// TODO: semver check

	// Read id
	if data, err = s.nextLine(br, buf); err != nil {
		return
	}
	if _, err = fmt.Sscanf(data, "id %x", &s.ID); err != nil {
		return
	}

	// Read size
	if data, err = s.nextLine(br, buf); err != nil {
		return
	}
	if _, err = fmt.Sscanf(data, "size %d", &s.Size); err != nil {
		return
	}

	// check ending
	if data, err = s.nextLine(br, buf); err != nil {
		return
	}
	if len(buf) > 0 {
		err = fmt.Errorf("bad ending %s", data)
	}

	if !full {
		return
	}

	// read chunks
	data, err = s.nextLine(br, buf)
	if err == io.EOF {
		// short
		err = nil
		return
	} else if err != nil {
		return
	}

	// Chunks present
	var aggr int64
	for {
		if aggr >= s.Size {
			if aggr != s.Size {
				err = fmt.Errorf("bad chunks in manifest")
			}
			break
		}
		chunk := Chunk{}

		// chunk id
		if data, err = s.nextLine(br, buf); err != nil {
			return
		}
		if _, err = fmt.Sscanf(data, "id %x", &chunk.ID); err != nil {
			return
		}
		// chunk size
		if data, err = s.nextLine(br, buf); err != nil {
			return
		}
		if _, err = fmt.Sscanf(data, "size %d", &chunk.Size); err != nil {
			return
		}
		// chunk offset
		if data, err = s.nextLine(br, buf); err != nil {
			return
		}
		if _, err = fmt.Sscanf(data, "offset %d", &chunk.Offset); err != nil {
			return
		}
		if data, err = s.nextLine(br, buf); err != nil {
			return
		}
		if len(buf) > 0 {
			err = fmt.Errorf("bad ending %s", data)
		}
		aggr += chunk.Size
		s.Chunks = append(s.Chunks, chunk)
	}

	return
}

// Initialize manifest from BLOB
func (s *Shadow) FromBlob(in io.Reader, full bool, chunkSize int64) (err error) {
	var chunkHasher hash.Hash
	var w io.Writer
	var chunk Chunk
	hasher := sha3.New256()
	var written int64

	for {
		if full {
			chunk = Chunk{
				Offset: s.Size,
			}
			chunkHasher = sha3.New256()
			w = io.MultiWriter(hasher, chunkHasher)
			written, err = io.CopyN(w, in, chunkSize)
			s.Size += written
			chunk.Size = written
			chunk.ID = chunkHasher.Sum(nil)
			s.Chunks = append(s.Chunks, chunk)
		} else {
			written, err = io.CopyN(hasher, in, chunkSize)
			s.Size += written
		}


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

func (s *Shadow) nextLine(in *bufio.Reader, buf []byte) (res string, err error) {
	buf, _, err = in.ReadLine()
	if err != nil {
		return
	}
	res = strings.TrimSpace(string(buf))
	return
}