package cmd
import (
	"flag"
	"io"
	"github.com/akaspin/bar/barctl/git"
	"fmt"
	"strings"
	"github.com/akaspin/bar/barctl/transport"
	"net/url"
	"github.com/nu7hatch/gouuid"
	"sync"
	"github.com/akaspin/bar/shadow"
	"path/filepath"
	"time"
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
	errOut io.Writer
	endpoint string
	trans *transport.TransportPool
	txId string
	gitCli *git.Git
}

func (c *GitPreCommitCmd) Bind(fs *flag.FlagSet, in io.Reader, out, errOut io.Writer) (err error) {
	c.errOut = errOut
	fs.StringVar(&c.endpoint, "endpoint", "http://localhost:3000/v1",
		"bard endpoint")
	return
}

func (c *GitPreCommitCmd) Do() (err error) {
	u, err := url.Parse(c.endpoint)
	if err != nil {
		return
	}
	c.trans = transport.NewTransportPool(u, 16, time.Minute * 5)

	txUUID, err := uuid.NewV4()
	if err != nil {
		return
	}
	c.txId = txUUID.String()

	if c.gitCli, err = git.NewGit(""); err != nil {
		return
	}

	// Check dirty status
	dirty, err := c.gitCli.DirtyFiles()
	if err != nil {
		return
	}
	dirty, err = c.gitCli.FilterByDiff("bar", dirty...)
	if len(dirty) > 0 {
		err = fmt.Errorf("Dirty BLOBs in working tree. Run following command to add BLOBs:\n\n    git -C %s add %s\n",
			c.gitCli.Root, strings.Join(dirty, " "))
		return
	}

	// Collect BLOBs from diff
	diffr, err := c.gitCli.Diff()
	if err != nil {
		return
	}
	fromDiff, err := c.gitCli.ParseDiff(diffr)
	if err != nil {
		return
	}

	toUpload, err := c.declareTx(fromDiff)
	wg := &sync.WaitGroup{}
	stat := map[string]error{}
	for _, b := range toUpload {
		wg.Add(1)
		go c.uploadOne(wg, b, stat)
	}
	wg.Wait()

	if len(stat) > 0 {
		err = fmt.Errorf("errors while upload: %v", stat)
	}
	return
}

func (c *GitPreCommitCmd) declareTx(diff []git.CommitBLOB) (res []git.CommitBLOB, err error) {
	var reqIDs []string
	idmap := map[string]git.CommitBLOB{}
	for _, b := range diff {
		reqIDs = append(reqIDs, b.ID)
		idmap[b.ID] = b
	}
	t, err := c.trans.Take()
	if err != nil {
		return
	}
	defer c.trans.Release(t)
	resIDs, err := t.DeclareCommitTx(c.txId, reqIDs)
	if err != nil {
		return
	}
	for _, id := range resIDs {
		delete(idmap, id)
	}
	for _, cb := range idmap {
		res = append(res, cb)
	}
	return
}

func (c *GitPreCommitCmd) uploadOne(wg *sync.WaitGroup, what git.CommitBLOB, stat map[string]error) (err error) {
	defer wg.Done()
	var s *shadow.Shadow
	var catR io.Reader
	if catR, err = c.gitCli.Cat(what.OID); err != nil {
		stat[what.Filename] = err
		return
	}
	if s, err = shadow.New(catR, 0); err != nil {
		stat[what.Filename] = err
		return
	}
	t, err := c.trans.Take()
	if err != nil {
		stat[what.Filename] = err
		return
	}
	defer c.trans.Release(t)
	fmt.Fprintf(c.errOut, "uploading %s", what.Filename)
	if err = t.Push(filepath.Join(c.gitCli.Root, what.Filename), s); err != nil {
		stat[what.Filename] = err
		return
	}
	return
}
