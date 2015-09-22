package cmd


// Git pre-commit hook. Used to upload all new/changed blobs
// to bard server:
//
// - Fails on uncommited bar-tracked BLOBs.
// - If working directory is clean - uploads BLOBs to bard.
type GitPreCommitCmd struct {
	endpoint string
}

func (c *GitPreCommitCmd) Bind()
