package cmd_test
import (
	"testing"
	"github.com/akaspin/bar/fixtures"
	"github.com/akaspin/bar/barctl/cmd"
	"flag"
	"github.com/stretchr/testify/assert"
	"os"
	"github.com/tamtam-im/flags"
	"github.com/akaspin/bar/barctl/transport"
)


func TestUpload1(t *testing.T)  {
	root := "test-upload1"
	endpoint := fixtures.RunServer(t, root)
	defer os.RemoveAll(root)

	uploadCmd := &cmd.UploadCommand{}

	subFS := flag.NewFlagSet(root, flag.ExitOnError)
	assert.NoError(t, uploadCmd.Bind(subFS, os.Stdin, os.Stdout, os.Stderr))

	// make some blobs
	bn1 := fixtures.MakeBLOB(t, 1234)
	defer fixtures.KillBLOB(bn1)
	sh1, err := fixtures.NewShadowFromFile(bn1)
	assert.NoError(t, err)

	bn2 := fixtures.MakeBLOB(t, 12)
	defer fixtures.KillBLOB(bn2)
	sh2, err := fixtures.NewShadowFromFile(bn2)
	assert.NoError(t, err)

	bn3 := fixtures.MakeBLOB(t, 4567)
	defer fixtures.KillBLOB(bn3)
	sh3, err := fixtures.NewShadowFromFile(bn3)
	assert.NoError(t, err)

	// upload third blob
	tr := &transport.Transport{endpoint}
	err = tr.Push(bn3, sh3)
	assert.NoError(t, err)

	// Replace third blob with it's shadow
	w, err := os.Create(bn3)
	err = sh3.Serialize(w)
	w.Close()

	flags.New(subFS).NoEnv().Boot([]string{
		"upload",
		"-endpoint=" + endpoint.String(),
		bn1,
		bn2,
		bn3,
	})

	// Upload
	err = uploadCmd.Do()
	assert.NoError(t, err)

	// Check on server
	info, err := os.Stat(fixtures.ServerStoredName(root, sh1.ID))
	assert.NoError(t, err)
	assert.False(t, info.IsDir())
	info, err = os.Stat(fixtures.ServerStoredName(root, sh2.ID))
	assert.NoError(t, err)
	assert.False(t, info.IsDir())
}
