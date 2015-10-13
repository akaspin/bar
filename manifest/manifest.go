package manifest
import (
	"io"
	"bytes"
	"fmt"
	"bufio"
	"golang.org/x/crypto/sha3"
	"hash"
	"strings"
	"encoding/hex"
	"github.com/akaspin/bar/proto/bar"
)

const (
	MANIFEST_HEADER = "BAR:MANIFEST"
)

type Manifest struct {
	Data
	Chunks []Chunk
}

func (s Manifest) String() (res string) {
	data := []byte{}
	buf := bytes.NewBuffer(data)
	(&s).Serialize(buf)
	res = string(buf.Bytes())
	return
}

// write initialized manifest to specific stream
func (s *Manifest) Serialize(out io.Writer) (err error) {
	if _, err = out.Write([]byte(MANIFEST_HEADER)); err != nil {
		return
	}
	body := fmt.Sprintf("\n\nid %s\nsize %d\n\n", s.ID, s.Size)
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

func NewFromAny(in io.Reader, chunkSize int64) (res *Manifest, err error) {
	r, isShadow, err := Peek(in)
	if err != nil {
		return
	}

	if isShadow {
		res, err = NewFromManifest(r)
	} else {
		res, err = NewFromBLOB(r, chunkSize)
	}
	return
}

// Make shadow from manifest
func NewFromManifest(in io.Reader) (res *Manifest, err error) {
	res = &Manifest{}
	err = res.ParseManifest(in)
	return
}

// Make shadow from BLOB
func NewFromBLOB(in io.Reader, chunkSize int64) (res *Manifest, err error) {
	res = &Manifest{}
	err = res.ParseBlob(in, chunkSize)
	return
}

// Parse manifest
func (s *Manifest) ParseManifest(in io.Reader) (err error) {
	var buf []byte
	var data string

	br := bufio.NewReader(in)

	// Check header
	if data, err = s.nextLine(br, buf); err != nil {
		return
	}
	if data != MANIFEST_HEADER {
		err = fmt.Errorf("bad manifest header %s", data)
		return
	}
	// check delimiter
	if data, err = s.nextLine(br, buf); err != nil {
		return
	}
	if len(buf) > 0 {
		err = fmt.Errorf("bad delimiter %s", data)
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
func (s *Manifest) ParseBlob(in io.Reader, chunkSize int64) (err error) {
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
	return
}

func (s *Manifest) nextLine(in *bufio.Reader, buf []byte) (res string, err error) {
	buf, _, err = in.ReadLine()
	if err != nil {
		return
	}
	res = strings.TrimSpace(string(buf))
	return
}

func (s Manifest) MarshalThrift() (res bar.Manifest, err error) {

	return
}