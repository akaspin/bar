package shadow_test
import (
	"testing"
	"github.com/akaspin/bar/shadow"
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"bytes"
	"os"
	"strings"
	"io/ioutil"
	"gopkg.in/bufio.v1"
)

// Make temporary BLOB
func makeBLOB(size int64) (name string, err error)  {
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

func killBLOB(name string) (err error) {
	return os.Remove(name)
}

func Test_Shadow_ToString(t *testing.T) {
	idHex := []byte("ac934d9a88b42aa3b40ef7c81a9dee1aad5a2cddccb00ae6abab9c38095fc15c")
	id := make([]byte, hex.DecodedLen(len(idHex)))
	_, err := hex.Decode(id, idHex)
	assert.NoError(t, err)

	sh := &shadow.Shadow{
		false,
		"0.0.1",
		id,
		1234,
		nil,
	}

	assert.Equal(t, "BAR:SHADOW\n\nversion 0.1.0\nid ac934d9a88b42aa3b40ef7c81a9dee1aad5a2cddccb00ae6abab9c38095fc15c\nsize 1234\n\n", (*sh).String())
}

func Test_Shadow_FromManifest(t *testing.T) {
	in := "BAR:SHADOW\n\nversion 0.1.0\nid ac934d9a88b42aa3b40ef7c81a9dee1aad5a2cddccb00ae6abab9c38095fc15c\nsize 1234\n\n"
	m := &shadow.Shadow{}
	err := m.FromManifest(bytes.NewReader(
		[]byte(in),
	))
	assert.NoError(t, err)
	assert.Equal(t, in, (*m).String())
}

func Test_Shadow_FromBLOB_50M(t *testing.T)  {
	bn, err := makeBLOB(1024 * 1024 * 50)
	assert.NoError(t, err)
	defer killBLOB(bn)

	m := &shadow.Shadow{}
	r, err := os.Open(bn)
	assert.NoError(t, err)
	defer r.Close()

	err = m.FromBlob(r)
	assert.NoError(t, err)
	assert.Equal(t, "BAR:SHADOW\n\nversion 0.1.0\nid 3339defdb3e5b3a2a71941b6b2bbdf7bb6525b61ba7eafb2cdb47428b3b65110\nsize 52428800\n\n", (*m).String())
}

func Test_Shadow_FromBLOB_5M(t *testing.T)  {
	bn, err := makeBLOB(1024 * 1024 * 5)
	assert.NoError(t, err)
	defer killBLOB(bn)

	m := &shadow.Shadow{}
	r, err := os.Open(bn)
	assert.NoError(t, err)
	defer r.Close()

	err = m.FromBlob(r)
	assert.NoError(t, err)
	assert.Equal(t, "BAR:SHADOW\n\nversion 0.1.0\nid 8a4038267545a17fa808226ac40c1d402275eb865e82ddc11ca37e6c326121a4\nsize 5242880\n\n", (*m).String())
}

func Test_Shadow_FromBLOB_2K(t *testing.T)  {
	bn, err := makeBLOB(1024 * 2)
	assert.NoError(t, err)
	defer killBLOB(bn)

	m := &shadow.Shadow{}
	r, err := os.Open(bn)
	assert.NoError(t, err)
	defer r.Close()

	err = m.FromBlob(r)
	assert.NoError(t, err)
	assert.Equal(t, "BAR:SHADOW\n\nversion 0.1.0\nid be4215176932949d887fa82241bbe0b03a44dc16ee2d23eedbc973e511ae8bb2\nsize 2048\n\n", (*m).String())
}

func Test_Shadow_FromBLOB_3b(t *testing.T)  {
	bn, err := makeBLOB(3)
	assert.NoError(t, err)
	defer killBLOB(bn)

	m := &shadow.Shadow{}
	r, err := os.Open(bn)
	assert.NoError(t, err)
	defer r.Close()

	err = m.FromBlob(r)
	assert.NoError(t, err)
	assert.Equal(t, "BAR:SHADOW\n\nversion 0.1.0\nid 1186d49a4ad620618f760f29da2c593b2ec2cc2ced69dc16817390d861e62253\nsize 3\n\n", (*m).String())
}

func Test_Shadow_FromAny_Manifest(t *testing.T) {
	in := "BAR:SHADOW\n\nversion 0.1.0\nid 82783ef12d68fd4c57fd7a8d7e42e7b71fc0fd13e5e30d459f15bc64a298395c\nsize 4\n\n"
	m := &shadow.Shadow{}
	err := m.FromAny(strings.NewReader(in))
	assert.NoError(t, err)
	assert.Equal(t, in, (*m).String())
}

func Test_Shadow_FromAny_50M(t *testing.T)  {
	bn, err := makeBLOB(1024 * 1024 * 50)
	assert.NoError(t, err)
	defer killBLOB(bn)

	m := &shadow.Shadow{}
	r, err := os.Open(bn)
	assert.NoError(t, err)
	defer r.Close()

	err = m.FromAny(r)
	assert.NoError(t, err)
	assert.Equal(t, "BAR:SHADOW\n\nversion 0.1.0\nid 3339defdb3e5b3a2a71941b6b2bbdf7bb6525b61ba7eafb2cdb47428b3b65110\nsize 52428800\n\n", (*m).String())
}

func Test_Shadow_FromAny_5M(t *testing.T)  {
	bn, err := makeBLOB(1024 * 1024 * 5)
	assert.NoError(t, err)
	defer killBLOB(bn)

	m := &shadow.Shadow{}
	r, err := os.Open(bn)
	assert.NoError(t, err)
	defer r.Close()

	err = m.FromAny(r)
	assert.NoError(t, err)
	assert.Equal(t, "BAR:SHADOW\n\nversion 0.1.0\nid 8a4038267545a17fa808226ac40c1d402275eb865e82ddc11ca37e6c326121a4\nsize 5242880\n\n", (*m).String())
}

func Test_Shadow_FromAny_2K(t *testing.T)  {
	bn, err := makeBLOB(1024 * 2)
	assert.NoError(t, err)
	defer killBLOB(bn)

	m := &shadow.Shadow{}
	r, err := os.Open(bn)
	assert.NoError(t, err)
	defer r.Close()

	err = m.FromAny(r)
	assert.NoError(t, err)
	assert.Equal(t, "BAR:SHADOW\n\nversion 0.1.0\nid be4215176932949d887fa82241bbe0b03a44dc16ee2d23eedbc973e511ae8bb2\nsize 2048\n\n", (*m).String())
}

func Test_Shadow_FromAny_3b(t *testing.T)  {
	bn, err := makeBLOB(3)
	assert.NoError(t, err)
	defer killBLOB(bn)

	m := &shadow.Shadow{}
	r, err := os.Open(bn)
	assert.NoError(t, err)
	defer r.Close()

	err = m.FromAny(r)
	assert.NoError(t, err)
	assert.Equal(t, "BAR:SHADOW\n\nversion 0.1.0\nid 1186d49a4ad620618f760f29da2c593b2ec2cc2ced69dc16817390d861e62253\nsize 3\n\n", (*m).String())
}
