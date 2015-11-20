package git

import (
	"encoding/json"
	"fmt"
	"github.com/tamtam-im/logx"
	"os"
	"path/filepath"
	"github.com/akaspin/bar/client/lists"
)

// Diversion spec
type DivertSpec struct {
	// Home ref
	Head string

	// Replaced files
	ToRecover []string

	// Target branch
	TargetBranch string

	// Target files
	TargetFiles []string
}

/*
Divert is ability to temporary switch branch preserving all blobs
in working tree.

	# Reset index and HEAD to otherbranch
	git reset otherbranch

	# make commit for otherbranch
	git add file-to-commit
	git commit "edited file"

	# force recreate otherbranch to here
	git branch -f otherbranch

	# Go back to where we were before
	# (two commits ago, the reset and the commit)
	git reset HEAD@{2}
*/
type Divert struct {
	*Git
}

// Make new divert on existing git repo
func NewDivert(host *Git) *Divert {
	return &Divert{host}
}

// Consistency checks before begin diversion and return spec
func (d *Divert) PrepareBegin(branch string, names ...string) (res DivertSpec, err error) {
	inProgress, err := d.IsInProgress()
	if err != nil {
		return
	}
	if inProgress {
		err = fmt.Errorf("divert already in progress")
		return
	}

	// check branches
	currentBranch, otherBranches, err := d.Git.GetBranches()
	if err != nil {
		return
	}
	if branch == currentBranch {
		err = fmt.Errorf("can not divert to current branch %s", branch)
		return
	}
	var exists bool
	for _, br := range otherBranches {
		if br == branch {
			exists = true
			break
		}
	}
	if !exists {
		err = fmt.Errorf("can not divert to nonexistent branch %s", branch)
		return
	}

	// get HEAD ref
	head, err := d.Git.GetRevParse("HEAD")
	if err != nil {
		return
	}

	res.Head = head
	res.TargetBranch = branch

	// Collect and check recoverable files
	if res.ToRecover, err = d.Git.LsTree(head, names...); err != nil {
		return
	}
	if len(res.ToRecover) > 0 {
		var dirty []string
		if dirty, err = d.Git.DiffFiles(res.ToRecover...); err != nil {
			return
		}
		if len(dirty) > 0 {
			err = fmt.Errorf("working tree is dirty %s. ", dirty)
			return
		}
		if dirty, err = d.Git.DiffIndex(res.ToRecover...); err != nil {
			return
		}
		if len(dirty) > 0 {
			err = fmt.Errorf("uncommited changes in index %s. ", dirty)
			return
		}
	}

	// collect target files
	if res.TargetFiles, err = d.Git.LsTree(branch, names...); err != nil {
		return
	}
	if len(res.TargetFiles) == 0 {
		err = fmt.Errorf("no files found in target branch")
		return
	}

	return
}

// Begin diversion to on branch where "names"
// is git-specific <tree-ish>
func (d *Divert) Begin(spec DivertSpec) (err error) {

	if err = d.writeSpec(spec); err != nil {
		return
	}

	// OOOK!!! Let's play with hammer!
	if err = d.Git.Reset(spec.TargetBranch); err != nil {
		return
	}

	if err = d.Git.Checkout(spec.TargetBranch, spec.TargetFiles...); err != nil {
		return
	}

	logx.Info("DIVERSION IN PROGRESS! DO NOT use any git commands!")
	return
}

// Finish diversion.
func (d *Divert) Commit(spec DivertSpec, message string) (err error) {
	// Add
	if err = d.Git.Add(spec.TargetFiles...); err != nil {
		return
	}
	// Commit

	if err = d.Git.Commit(message); err != nil {
		return
	}
	// force recreate index
	if err = d.Git.BranchRecreate(spec.TargetBranch); err != nil {
		return
	}
	logx.Infof("diversion commited")
	return
}

// Abort diversion
func (d *Divert) Cleanup(spec DivertSpec) (err error) {
	// Remove orphan diverted files
	orphans := map[string]struct{}{}
	for _, f := range spec.TargetFiles {
		orphans[f] = struct{}{}
	}
	for _, f := range spec.ToRecover {
		_, ok := orphans[f]
		if ok {
			delete(orphans, f)
		}
	}
	for f, _ := range orphans {
		os.Remove(lists.OSFromSlash(lists.OSJoin(d.Git.Root, f)))
		logx.Debugf("removed orphan %s", f)
	}

	// Reset to Head
	if err = d.Git.Reset(spec.Head); err != nil {
		return
	}

	if err = d.Git.Checkout(spec.Head, spec.ToRecover...); err != nil {
		return
	}

	logx.Info("cleanup finished")
	return
}

// Is divert in progress
func (d *Divert) IsInProgress() (res bool, err error) {
	_, statErr := os.Stat(d.specFilename())
	if statErr == nil {
		return true, nil
	}
	if os.IsNotExist(statErr) {
		return
	}
	err = statErr
	return
}

func (d *Divert) ReadSpec() (res DivertSpec, err error) {
	f, err := os.Open(filepath.FromSlash(d.specFilename()))
	if err != nil {
		return
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&res)
	return
}

func (d *Divert) CleanSpec() (err error) {
	return os.Remove(d.specFilename())
}

func (d *Divert) writeSpec(spec DivertSpec) (err error) {
	os.MkdirAll(filepath.Dir(d.specFilename()), 0755)
	f, err := os.Create(d.specFilename())
	if err != nil {
		return
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(spec)
	return
}

func (d *Divert) specFilename() string {
	return filepath.Join(d.Git.Root, ".git", "bar", "divert.json")
}
