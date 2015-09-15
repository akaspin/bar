package shadow_test
import (
	"testing"
	"github.com/akaspin/bar/shadow"
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"bytes"
	"os"
	"strings"
"github.com/akaspin/bar/fixtures"
)



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

func Test_Shadow_FromBLOB_20M(t *testing.T)  {
	bn, err := fixtures.MakeBLOB(1024 * 1024 * 20 + 5)
	assert.NoError(t, err)
	defer fixtures.KillBLOB(bn)

	m := &shadow.Shadow{}
	r, err := os.Open(bn)
	assert.NoError(t, err)
	defer r.Close()

	err = m.FromBlob(r)
	assert.NoError(t, err)
	assert.Equal(t, "BAR:SHADOW\n\nversion 0.1.0\nid 9e39ad7cf632a038a5a2e0f9144f6ea4aff04ff11803c169cb24f60e56444f08\nsize 20971525\n\n", (*m).String())
}

func Test_Shadow_FromBLOB_2M(t *testing.T)  {
	bn, err := fixtures.MakeBLOB(1024 * 1024 * 2 + 467)
	assert.NoError(t, err)
	defer fixtures.KillBLOB(bn)

	m := &shadow.Shadow{}
	r, err := os.Open(bn)
	assert.NoError(t, err)
	defer r.Close()

	err = m.FromBlob(r)
	assert.NoError(t, err)
	assert.Equal(t, "BAR:SHADOW\n\nversion 0.1.0\nid fd76eb2f9866a12c6c3a2f884d5350b38319bc510106a7ba78789cc5507158b8\nsize 2097619\n\n", (*m).String())
}

func Test_Shadow_FromBLOB_2K(t *testing.T)  {
	bn, err := fixtures.MakeBLOB(1024 * 2)
	assert.NoError(t, err)
	defer fixtures.KillBLOB(bn)

	m := &shadow.Shadow{}
	r, err := os.Open(bn)
	assert.NoError(t, err)
	defer r.Close()

	err = m.FromBlob(r)
	assert.NoError(t, err)
	assert.Equal(t, "BAR:SHADOW\n\nversion 0.1.0\nid be4215176932949d887fa82241bbe0b03a44dc16ee2d23eedbc973e511ae8bb2\nsize 2048\n\n", (*m).String())
}

func Test_Shadow_FromBLOB_3b(t *testing.T)  {
	bn, err := fixtures.MakeBLOB(3)
	assert.NoError(t, err)
	defer fixtures.KillBLOB(bn)

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

func Test_Shadow_FromAny_20M(t *testing.T)  {
	bn, err := fixtures.MakeBLOB(1024 * 1024 * 20 + 5)
	assert.NoError(t, err)
	defer fixtures.KillBLOB(bn)

	m := &shadow.Shadow{}
	r, err := os.Open(bn)
	assert.NoError(t, err)
	defer r.Close()

	err = m.FromAny(r)
	assert.NoError(t, err)
	assert.Equal(t, "BAR:SHADOW\n\nversion 0.1.0\nid 9e39ad7cf632a038a5a2e0f9144f6ea4aff04ff11803c169cb24f60e56444f08\nsize 20971525\n\n", (*m).String())
}

func Test_Shadow_FromAny_2M(t *testing.T)  {
	bn, err := fixtures.MakeBLOB(1024 * 1024 * 2 + 467)
	assert.NoError(t, err)
	defer fixtures.KillBLOB(bn)

	m := &shadow.Shadow{}
	r, err := os.Open(bn)
	assert.NoError(t, err)
	defer r.Close()

	err = m.FromAny(r)
	assert.NoError(t, err)
	assert.Equal(t, "BAR:SHADOW\n\nversion 0.1.0\nid fd76eb2f9866a12c6c3a2f884d5350b38319bc510106a7ba78789cc5507158b8\nsize 2097619\n\n", (*m).String())
}

func Test_Shadow_FromAny_2K(t *testing.T)  {
	bn, err := fixtures.MakeBLOB(1024 * 2)
	assert.NoError(t, err)
	defer fixtures.KillBLOB(bn)

	m := &shadow.Shadow{}
	r, err := os.Open(bn)
	assert.NoError(t, err)
	defer r.Close()

	err = m.FromAny(r)
	assert.NoError(t, err)
	assert.Equal(t, "BAR:SHADOW\n\nversion 0.1.0\nid be4215176932949d887fa82241bbe0b03a44dc16ee2d23eedbc973e511ae8bb2\nsize 2048\n\n", (*m).String())
}

func Test_Shadow_FromAny_3b(t *testing.T)  {
	bn, err := fixtures.MakeBLOB(3)
	assert.NoError(t, err)
	defer fixtures.KillBLOB(bn)

	m := &shadow.Shadow{}
	r, err := os.Open(bn)
	assert.NoError(t, err)
	defer r.Close()

	err = m.FromAny(r)
	assert.NoError(t, err)
	assert.Equal(t, "BAR:SHADOW\n\nversion 0.1.0\nid 1186d49a4ad620618f760f29da2c593b2ec2cc2ced69dc16817390d861e62253\nsize 3\n\n", (*m).String())
}
