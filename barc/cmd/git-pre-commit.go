package cmd
import (
	"flag"
	"io"
	"github.com/akaspin/bar/barc/git"
	"fmt"
	"strings"
	"github.com/akaspin/bar/barc/transport"
	"net/url"
	"github.com/nu7hatch/gouuid"
	"sync"
	"github.com/akaspin/bar/shadow"
	"path/filepath"
	"time"
	"github.com/tamtam-im/logx"
)


/*
Git pre-commit hook. Used to upload all new/changed blobs
to bard server:

- Fails on uncommited bar-tracked BLOBs.
- If working directory is clean - uploads BLOBs to bard.

To use with git git-clean MUST be registered in git. Also
git pre-commit hook MUST be registered:

	$ cat > .git/hooks/pre-commit <<EOF
	#!/usr/bin/env sh
	set -e
	barctl git-pre-commit -endpoint=http://my.bar.server/v1
	EOF
	chmod +x .git/hooks/pre-commit
*/

type GitPreCommitCmd struct {
	endpoint string

	// commit transaction ID
	txId string
	transport *transport.TransportPool
	git *git.Git
}

func (c *GitPreCommitCmd) Bind(wd string, fs *flag.FlagSet, in io.Reader, out io.Writer) (err error) {
	fs.StringVar(&c.endpoint, "endpoint", "http://localhost:3000/v1",
		"bard endpoint")
	return
}

func (c *GitPreCommitCmd) Do() (err error) {
	txUUID, err := uuid.NewV4()
	if err != nil {
		return
	}
	c.txId = txUUID.String()

	logx.Debugf("bar commit %s", c.txId)

	u, err := url.Parse(c.endpoint)
	if err != nil {
		return
	}
	c.transport = transport.NewTransportPool(u, 16, time.Minute * 5)

	if c.git, err = git.NewGit(""); err != nil {
		return
	}

	// Check dirty status
	dirty, err := c.git.DiffFiles()
	if err != nil {
		return
	}

	if len(dirty) > 0 {
		dirty, err = c.git.FilterByDiff("bar", dirty...)
		if err != nil {
			return
		}
	}

	if len(dirty) > 0 {
		err = fmt.Errorf(
			"Dirty BLOBs in working tree. Run following command to add BLOBs:\n\n    git -C %s add %s\n",
			c.git.Root, strings.Join(dirty, " "))
		return
	}

	// Collect BLOBs from diff
	diffr, err := c.git.Diff()
	if err != nil {
		return
	}
	fromDiff, err := c.git.ParseDiff(diffr)
	if err != nil {
		return
	}

	toUpload, err := c.declareTx(fromDiff)
	if err != nil {
		return
	}

	wg := &sync.WaitGroup{}
	errs := []error{}
	for _, b := range toUpload {
		wg.Add(1)
		go func(oid string, filename string) {
			defer wg.Done()
			if err1 := c.uploadOne(oid, filename); err1 != nil {
				errs = append(errs, err1)
			}
		}(b.OID, b.Filename)
	}
	wg.Wait()

	if len(errs) > 0 {
		err = fmt.Errorf("errors while upload: %s", errs)
	}
	return
}

// Declare transaction
func (c *GitPreCommitCmd) declareTx(diff []git.DiffEntry) (res []git.DiffEntry, err error) {
	if len(diff) == 0 {
		logx.Debugf("no files to upload")
		return
	}

	// Prepare data for request
	var reqIDs, resIDs []string
	idmap := map[string]git.DiffEntry{}
	for _, b := range diff {
		reqIDs = append(reqIDs, b.ID)
		idmap[b.ID] = b
	}

	t, err := c.transport.Take()
	if err != nil {
		return
	}
	defer c.transport.Release(t)

	if resIDs, err = t.DeclareCommitTx(c.txId, reqIDs); err != nil {
		return
	}
	for _, id := range resIDs {
		res = append(res, idmap[id])
	}
	return
}

func (c *GitPreCommitCmd) uploadOne(oid string, filename string) (err error) {
	var s *shadow.Shadow
	t, err := c.transport.Take()
	if err != nil {
		return
	}
	defer c.transport.Release(t)

	var catR io.Reader
	if catR, err = c.git.Cat(oid); err != nil {
		return
	}
	if s, err = shadow.NewFromManifest(catR); err != nil {
		return
	}

	logx.Infof("uploading %s: %d bytes", filename, s.Size)

	if err = t.Push(filepath.Join(c.git.Root, filename), s); err != nil {
		return
	}
	return
}
