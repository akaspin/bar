package cmd_test
import (
	"testing"
	"github.com/akaspin/bar/fixtures"
	"github.com/stretchr/testify/assert"
	"os"
)

func Test_Integration_GitInit_Normal(t *testing.T) {
	skip(t)

	tree := fixtures.NewTree("")
	defer tree.Squash()

	assert.NoError(t, tree.Populate())
	assert.NoError(t, tree.InitGit())

	endpoint := fixtures.RunServer(t, "test_srv_data_git-init")
	defer os.RemoveAll("test_srv_data_git-init")

	_, err := tree.Run("barc", "-log-level=DEBUG", "git-init",
		"-endpoint=" + endpoint.String(),
		"-log=DEBUG",
	).Output()
	assert.NoError(t, err)
}



