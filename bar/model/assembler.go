package model
import (
	"io/ioutil"
	"io"
	"github.com/akaspin/go-contentaddressable"
	"golang.org/x/crypto/sha3"
	"path/filepath"
	"github.com/akaspin/bar/bar/lists"
	"os"
	"github.com/tamtam-im/logx"
	"github.com/akaspin/bar/proto"
	"strings"
	"golang.org/x/net/context"
	"github.com/akaspin/bar/concurrent"
)

type Assembler struct  {
	Where string
	model *Model
}

func NewAssembler(m *Model) (res *Assembler, err error) {
	where, err := ioutil.TempDir("", "")
	if err != nil {
		return
	}
	res = &Assembler{where, m}
	return
}

// Store chunk in assemble
func (a *Assembler) StoreChunk(r io.Reader, id proto.ID) (err error) {
	lock, err := a.model.FdLocks.Take()
	if err != nil {
		return
	}
	defer lock.Close()

	caOpts := contentaddressable.DefaultOptions()
	caOpts.Hasher = sha3.New256()

	w, err := contentaddressable.NewFileWithOptions(
		filepath.Join(a.Where, id.String()), caOpts)
	if os.IsExist(err) {
		err = nil
		return
	}
	if err != nil {
		return
	}
	defer w.Close()

	if _, err = io.Copy(w, r); err != nil {
		return
	}
	err = w.Accept()
	return
}

func (a *Assembler) StoredChunks() (res []string, err error) {
	err = filepath.Walk(a.Where,
		func(path string, info os.FileInfo, errIn error) (err error) {
			if info.IsDir() {
				return
			}
			n := filepath.Base(path)
			if !strings.HasSuffix(n, "-temp") {
				res = append(res, n)
			}
			return
		})
	return
}

// Assemble target files from stored chunks
func (a *Assembler) Done(what lists.BlobMap) (err error) {
	logx.Tracef("assembling %s", what.Names())

	var req, res []interface{}
	for k, v := range what {
		req = append(req, lists.BlobLink{v, k})
	}

	err = a.model.BatchPool.Do(
		func(ctx context.Context, in interface{}) (out interface{}, err error) {
			r := in.(lists.BlobLink)

			lock, err := a.model.FdLocks.Take()
			if err != nil {
				return
			}
			defer lock.Close()

			w, err := os.Create(filepath.Join(a.model.WD,
				r.Name + r.Manifest.ID.String()))
			if err != nil {
				return
			}
			defer w.Close()

			for _, chunk := range r.Manifest.Chunks {
				if err = a.writeChunkTo(w, chunk.ID); err != nil {
					return
				}
			}
			err = a.commitBlob(r.Name, r.Manifest.ID)

			return
		}, &req, &res, concurrent.DefaultBatchOptions().AllowErrors(),
	)
	if err != nil {
		return
	}
	defer a.Close()
	return
}

func (a *Assembler) commitBlob(name string, id proto.ID) (err error) {
	dst := filepath.Join(a.model.WD, name)
	src := dst + id.String()
	bak := dst + ".bak"

	os.Rename(dst, bak)
	if err = os.Rename(src, dst); err != nil {
		os.Remove(dst)
		os.Rename(bak, dst)
		return
	}
	defer os.Remove(src)
	return
}

func (a *Assembler) writeChunkTo(w io.Writer, id proto.ID) (err error) {
	lock, err := a.model.FdLocks.Take()
	if err != nil {
		return
	}
	defer lock.Close()

	name := filepath.Join(a.Where, id.String())
	r, err := os.Open(name)
	if err != nil {
		return
	}
	defer r.Close()

	_, err = io.Copy(w, r)
	return
}

func (a *Assembler) Close() {
	os.RemoveAll(a.Where)
}