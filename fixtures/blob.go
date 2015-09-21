package fixtures
import (
	"io/ioutil"
	"bufio"
	"os"
	"io"
	"bytes"
	"strings"
	"github.com/akaspin/bar/shadow"
)

// Make temporary BLOB
func MakeBLOB(size int64) (name string, err error)  {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return
	}
	defer f.Close()
	name = f.Name()

	var i int64
	var j uint8
	buf := bufio.NewWriter(f)

	for i=0; i < size; i++ {
		if _, err = buf.Write([]byte{j}); err != nil {
			return
		}
		j++
		if j > 126 {
			j = 0
		}
	}
	err = buf.Flush()
	return
}

func KillBLOB(name string) (err error) {
	return os.Remove(name)
}

func CleanManifest(in string) (out string) {
	r := bufio.NewReader(bytes.NewReader([]byte(in)))
	o := new(bytes.Buffer)
	var buf []byte
	var err error
	for {
		buf, _, err = r.ReadLine()
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return
		}
		_, err = o.Write([]byte(strings.TrimSpace(string(buf)) + "\n"))
	}
	out = string(o.Bytes())
	return
}

func NewShadowFromFile(filename string, full bool, chunkSize int64) (res *shadow.Shadow, err error) {
	r, err := os.Open(filename)
	if err != nil {
		return
	}
	defer r.Close()
	res = &shadow.Shadow{}
	err = res.FromAny(r, full, chunkSize)
	return
}
