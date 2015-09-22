package fixtures
import (
	"io/ioutil"
	"bufio"
	"os"
	"io"
	"bytes"
	"strings"
	"github.com/akaspin/bar/shadow"
	"testing"
	"github.com/stretchr/testify/assert"
)

// Make temporary BLOB
func MakeBLOB(t *testing.T, size int64) (name string)  {
	name, err := MakeBLOBPure(size)
	assert.NoError(t, err)
	return
}

func MakeBLOBPure(size int64) (name string, err error)  {
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
		_, err = buf.Write([]byte{j})
		if err != nil {
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

// Clean input and return new reader
func CleanInput(in string) (out io.Reader, size int64) {
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
	out = bytes.NewReader(o.Bytes())
	size = int64(len(o.Bytes()))
	return
}

func FixStream(in string) (res string) {
	r, _ := CleanInput(in)
	data, _ := ioutil.ReadAll(r)
	res = string(data)
	return
}

func NewShadowFromFile(filename string) (res *shadow.Shadow, err error) {
	info, err := os.Stat(filename)
	if err != nil {
		return
	}

	r, err := os.Open(filename)
	if err != nil {
		return
	}
	defer r.Close()
	res, err = shadow.New(r, info.Size())
	return
}
