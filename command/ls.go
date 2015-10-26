package command
import (
	"github.com/spf13/cobra"
	"fmt"
	"github.com/akaspin/bar/client/model"
	"github.com/akaspin/bar/proto"
	"github.com/akaspin/bar/client/lists"
	"github.com/akaspin/bar/client/transport"
	"text/tabwriter"
"strings"
	"sort"
)

type LsCmd struct {
	*Environment
	*CommonOptions

	UseGit bool

	NoHeader bool
	FullID bool

	NoBlobs bool
	NoManifests bool

	NoRemote bool
	NoName bool
	NoID bool
	NoSize bool
}

func (c *LsCmd) Init(cc *cobra.Command) {
	cc.Use = "ls"
	cc.Short = "show info about bar-tracked BLOBs"

	cc.Flags().BoolVarP(&c.UseGit, "git", "", false, "use git infrastructure")

	cc.Flags().BoolVarP(&c.NoHeader, "no-header", "", false,
		"do not print header")
	cc.Flags().BoolVarP(&c.FullID, "full-id", "", false,
		"do not trim blob IDs")

	cc.Flags().BoolVarP(&c.NoBlobs, "no-blobs", "", false,
		"do not collect BLOBs")
	cc.Flags().BoolVarP(&c.NoManifests, "no-manifests", "", false,
		"do not collect manifests")

	cc.Flags().BoolVarP(&c.NoRemote, "no-remote", "", false,
		"do not request bard for BLOBs states")
	cc.Flags().BoolVarP(&c.NoName, "no-name", "", false,
		"do not print BLOB filenames")
	cc.Flags().BoolVarP(&c.NoID, "no-id", "", false, "do not print BLOB IDs")
	cc.Flags().BoolVarP(&c.NoSize, "no-size", "", false,
		"do not print BLOB sizes")
}

func (c *LsCmd) Run(args ...string) (err error) {

	if c.NoBlobs && c.NoManifests {
		err = fmt.Errorf("both -no-blobs and -no-manifests are on")
		return
	}

	mod, err := model.New(c.WD, c.UseGit, proto.CHUNK_SIZE, c.PoolSize)
	if err != nil {
		return
	}

	feed := lists.NewFileList(args...).ListDir(c.WD)

	var dirty map[string]struct{}
	if c.UseGit {
		if feed, err = mod.Git.FilterByAttr("bar", feed...); err != nil {
			return
		}
		var dirtyst []string
		if dirtyst, err = mod.Git.DiffFiles(feed...); err != nil {
			return
		}
		dirty = map[string]struct{}{}
		for _, n := range dirtyst {
			dirty[n] = struct{}{}
		}
	}

	blobs, err := mod.FeedManifests(!c.NoBlobs, !c.NoManifests, true, feed...)
	if err != nil {
		return
	}

	missingOnRemote := map[proto.ID]struct{}{}
	if !c.NoRemote {
		var exists []proto.ID
		trans := transport.NewTransport(mod, "", c.Endpoint, c.PoolSize)
		if exists, err = trans.Check(blobs.IDMap().IDs()); err != nil {
			return
		}
		for _, id := range exists {
			missingOnRemote[id] = struct{}{}
		}
	}

	// print this stuff
	w := new(tabwriter.Writer)
	w.Init(c.Stdout, 0, 8, 2, '\t', 0)

	var line []string
	toLine := func(term string) {
		line = append(line, term)
	}
	flushLine := func() {
		fmt.Fprintln(w, strings.Join(line, "\t"))
		line = []string{}
	}

	if !c.NoHeader {
		if !c.NoName {
			toLine("NAME")
		}
		if !c.NoBlobs && !c.NoManifests {
			toLine("BLOB")
		}
		if !c.NoRemote {
			toLine("SYNC")
		}
		if c.UseGit {
			toLine("GIT")
		}
		if !c.NoID {
			toLine("ID")
		}
		if !c.NoSize {
			toLine("SIZE")
		}
		flushLine()
	}

	var names sort.StringSlice
	for n, _ := range blobs {
		names = append(names, n)
	}
	names.Sort()

	var blobMap map[string]bool
	if !c.NoBlobs && !c.NoManifests {
		if blobMap, err = mod.IsBlobs(names...); err != nil {
			return
		}
	}

	for _, name := range names {
		if !c.NoName {
			toLine(name)
		}
		if !c.NoBlobs && !c.NoManifests {
			if blobMap[name] {
				toLine("yes")
			} else {
				toLine("no")
			}
		}
		if !c.NoRemote {
			if _, missing := missingOnRemote[blobs[name].ID]; missing {
				toLine("no")
			} else {
				toLine("yes")
			}
		}
		if c.UseGit {
			if _, bad := dirty[name]; !bad {
				toLine("ok")
			} else {
				toLine("dirty")
			}
		}
		if !c.NoID {
			if !c.FullID {
				toLine(blobs[name].ID.String()[:12])
			} else {
				toLine(blobs[name].ID.String())
			}
		}
		if !c.NoSize {
			toLine(fmt.Sprintf("%d", blobs[name].Size))
		}

		flushLine()
	}
	w.Flush()

	return
}