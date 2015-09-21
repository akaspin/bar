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
	f, err := ioutil.TempFile("", "")
	assert.NoError(t, err)
	defer f.Close()
	name = f.Name()

	var i int64
	var j uint8
	buf := bufio.NewWriter(f)

	for i=0; i < size; i++ {
		_, err = buf.Write([]byte{j})
		assert.NoError(t, err)
		j++
		if j > 126 {
			j = 0
		}
	}
	assert.NoError(t, buf.Flush())
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
