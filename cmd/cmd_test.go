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

func Test_CleanCmd_FromBLOBShort(t *testing.T) {
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

func Test_CleanCmd_FromBLOBFull(t *testing.T) {
	var id []byte
	fmt.Sscanf("0a8808fe65a3d752b175bcc420be536c4f6e4b1328ad52bab591ac6cd6e8b410", "%x", &id)

	manifest := `BAR:SHADOW

		version 0.1.0
		id 0a8808fe65a3d752b175bcc420be536c4f6e4b1328ad52bab591ac6cd6e8b410
		size 5242926


		id 7dd05dc961a9bfd1d7d3310e890d8be20a7f0667b17b1380a520adc491f202ce
		size 1048576
		offset 0

		id ce3ab016cb42aa12dc1eab3976ad3a49e134966e5cdf24c2c48f93b56239460e
		size 1048576
		offset 1048576

		id 39342fd4b4fe70f6b6640f67f86a8fbc50dc6fdf70a414fa517bd2d57073bddf
		size 1048576
		offset 2097152

		id 824b76f05b554f8ff29f5a0cf67bc173a47d3919050de53c9515e961b25f0060
		size 1048576
		offset 3145728

		id 5f2da039a71d1693b7df7b3eaaad25981cd50e733ea474c0243e47add598d451
		size 1048576
		offset 4194304

		id eac8da2fde911b90daad9b0700102258f98f3ffdfff4a9af8ba352346837c1b5
		size 46
		offset 5242880
		`
	s := &shadow.Shadow{}
	assert.NoError(t, s.FromAny(bytes.NewReader([]byte(manifest)), true))

	bn, err := fixtures.MakeBLOB(1024 * 1024 * 5 + 46)
	assert.NoError(t, err)
	defer fixtures.KillBLOB(bn)

	in, err := os.Open(bn)
	assert.NoError(t, err)
	defer in.Close()
	out := bytes.NewBuffer(nil)

	err = cmd.CleanCmd([]string{"clean", "-full"}, in, out, os.Stderr)
	assert.NoError(t, err)

	assert.Equal(t, (*s).String(), string(out.Bytes()))
}

func Test_CleanCmd_FromManifestShort(t *testing.T) {
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
