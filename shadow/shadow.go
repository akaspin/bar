package shadow
import (
	"io"
	"bytes"
	"fmt"
	"gopkg.in/bufio.v1"
	"golang.org/x/crypto/sha3"
	"hash"
	"strings"
	"encoding/hex"
)

const (
	CHUNK_SIZE = 1024 * 1024   // 1M
	SHADOW_HEADER = "BAR:SHADOW"
	VERSION = "0.1.0"
)



type Shadow struct {
	IsFromShadow bool
	Version string
	ID string
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
	if _, err = out.Write([]byte(SHADOW_HEADER)); err != nil {
		return
	}
	body := fmt.Sprintf("\n\nversion %s\nid %s\nsize %d\n\n",
		VERSION, s.ID, s.Size)
	if _, err = out.Write([]byte(body)); err != nil {
		return
	}

	if _, err = out.Write([]byte("\n")); err != nil {
		return
	}
	for _, chunk := range s.Chunks {
		if _, err = out.Write([]byte(chunk.String())); err != nil {
			return
		}
	}
	return
}

func New(in io.Reader, size int64) (res *Shadow, err error) {
	r, isShadow, err := Peek(in)
	if err != nil {
		return
	}

	res = &Shadow{}
	if isShadow {
		err = res.parseManifest(r)
	} else {
		err = res.parseBlob(r, GuessChunkSize(size))
	}
	return
}

// Parse manifest
func (s *Shadow) parseManifest(in io.Reader) (err error) {
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

	// Read id
	if data, err = s.nextLine(br, buf); err != nil {
		return
	}
	if _, err = fmt.Sscanf(data, "id %s", &s.ID); err != nil {
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
		if _, err = fmt.Sscanf(data, "id %s", &chunk.ID); err != nil {
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
func (s *Shadow) parseBlob(in io.Reader, chunkSize int64) (err error) {
	var chunkHasher hash.Hash
	var w io.Writer
	var chunk Chunk
	hasher := sha3.New256()
	var written int64

	for {
		chunk = Chunk{
			Offset: s.Size,
		}
		chunkHasher = sha3.New256()
		w = io.MultiWriter(hasher, chunkHasher)
		written, err = io.CopyN(w, in, chunkSize)
		s.Size += written
		chunk.Size = written
		chunk.ID = hex.EncodeToString(chunkHasher.Sum(nil))
		s.Chunks = append(s.Chunks, chunk)

		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return
		}
	}
	s.ID = hex.EncodeToString(hasher.Sum(nil))
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
