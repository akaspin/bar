package cmd_test
import (
	"testing"
	"github.com/akaspin/bar/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/akaspin/bar/cmd"
	"os"
	"bytes"
	"github.com/akaspin/bar/shadow"
	"fmt"
)

func Test_CleanCmd_FromBLOB(t *testing.T) {
	var id []byte
	fmt.Sscanf("0a8808fe65a3d752b175bcc420be536c4f6e4b1328ad52bab591ac6cd6e8b410", "%x", &id)
	expect := (shadow.Shadow{
		Version: "0.1.0",
		ID: id,
		Size: 5242926,
	}).String()

	bn, err := fixtures.MakeBLOB(1024 * 1024 * 5 + 46)
	assert.NoError(t, err)
	defer fixtures.KillBLOB(bn)

	in, err := os.Open(bn)
	assert.NoError(t, err)
	defer in.Close()
	out := bytes.NewBuffer(nil)

	err = cmd.CleanCmd([]string{"clean"}, in, out, os.Stderr)
	assert.NoError(t, err)

	assert.Equal(t, expect, string(out.Bytes()))
}

func Test_CleanCmd_FromManifest(t *testing.T) {
	var id []byte
	fmt.Sscanf("0a8808fe65a3d752b175bcc420be536c4f6e4b1328ad52bab591ac6cd6e8b410", "%x", &id)
	expect := (shadow.Shadow{
		Version: "0.1.0",
		ID: id,
		Size: 5242926,
	}).String()

	in := bytes.NewBuffer([]byte(expect))
	out := bytes.NewBuffer(nil)

	err := cmd.CleanCmd([]string{"clean"}, in, out, os.Stderr)
	assert.NoError(t, err)

	assert.Equal(t, expect, string(out.Bytes()))
}
