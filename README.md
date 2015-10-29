# Bar

**This software is beta**

Bar was developed to solve one simple problem. Large BLOBs can be stored in 
working tree. But our developers do not need these BLOBs. Or they are need 
some of them.

We use git. Bar designed to to minor additions to the workflow. Also bar can 
work without git.

Consider the following scenario:

1. About 2TB of BLOBs in project.
2. All BLOBs should be under version control.
3. All blobs along are not needed locally.
4. Some staff do not have a clue about version control.
5. Some works under windows.
6. They need one-button solution.
7. Cry!
8. Use bar.

## Basic usage

> For now assume what `bard` is deployed and listening at `:3001`.

To set the bar in the repository use `bar git-init` command. Only one 
important flag "endpoint" is just HTTP endpoint of `bard` server. 
    
    $ bar git install --endpoint=localhost:3001
    
Bar-tracked BLOBs are defined by git attributes.

    # .gitattributes
    my/blobs    filter=bar
    
After that you can add, commit and push as usual. Bar will upload new and 
changed BLOBs to `bard` server on commit.
 
    $ echo "test" > my/blobs/test.txt
    $ git add -A
    $ git commit -m "initial commit"  # <-- Here BLOBs will be uploaded
    ...

For this moment all is simple. To transform BLOBs to manifests 
use `git bar squash`:

    $ git bar-squash
    ...
    
    $ git status 
    On branch master
    nothing to commit, working directory clean
    
    $ cat my/blobs/test.txt
    BAR:MANIFEST
    
    version 0.1.0
    id 309a3490190131517180e3827398a665c1eef2b9b2b41108a08a59f4bb15d301
    size 4
    
    
    id 309a3490190131517180e3827398a665c1eef2b9b2b41108a08a59f4bb15d301
    size 3
    offset 0
    
Now BLOB transformed to manifest. Shadow is small text manifest describing BLOB. 
Bar installs own `filter` to git repo. This leverages git to store 
manifests in index.

    $ git ls-files -s --cache my/blobs/test.txt
    100644 c8c0b7267bf989a8d15c885dcdbafb10b60d244e 0	my/blobs/test.txt
    
    $ git cat-file -p c8c0b7267bf989a8d15c885dcdbafb10b60d244e
    BAR:MANIFEST
        
    version 0.1.0
    id 309a3490190131517180e3827398a665c1eef2b9b2b41108a08a59f4bb15d301
    size 4
    
    
    id 309a3490190131517180e3827398a665c1eef2b9b2b41108a08a59f4bb15d301
    size 3
    offset 0
    
To upload BLOBs without squashing them use `git bar-up`.
    
Also on checkout all non-existent BLOBs will be become stored as manifests. To 
get them back use `git bar-down`.

    $ git bar-down my/blobs/test.txt
    ... 
    $ git status
    On branch master
    nothing to commit, working directory clean
        
    $ cat my/blobs/test.txt
    test

To check status of bar-tracked BLOBs use `git bar-ls`
 
    $ git bar-ls
    NAME                BLOB    SYNC    ID                  SIZE
    my/blobs/test.txt   yes     yes     309a349019013151    4
    
    $ git bar-squash
    $ git bar-ls
    NAME                BLOB    SYNC    ID                  SIZE
    my/blobs/test.txt   no      yes     309a349019013151    4

## Covert ops

With `bar divert` you can change specific files from inactive branch without 
touching other files in working tree.

    $ git bar-divert other-branch my/pretty/dir "my pretty/file"
    ...
    
    $ git bar-divert-finish
    # or $ git bar-divert-abort to abort
    
To check divert status use `bar divert status`

    $ git bar-divert-status
    
You can also push inactive branch:

    $ git bar-divert-push other-branch
   
## Installation

To install just grab latest binaries archive from releases and unpack 
somewhere to `PATH`.

## `bar server` BLOB server

    $ bar server run \
        --bind-http=:3000 \
        --bind-rpc=:3001 \
        --storage=block:root=testdata \
        --endpoint=${HOSTNAME}:3001 \
        --endpoint-http=http://${HOSTNAME}:3000/v1 \
        --bindir=dist/bindir
    DEBUG server.go:23: serving at http://0.0.0.0:3000/v1

Use `bar server run -help` to get available options.

## Gitless usage

Bar doesn't require git to work. All `git bar-*` commands is just git aliases:

    git bar-squash  ->  bar up --squash
    git bar-up      ->  bar up
    git bar-down    ->  bar down
    git bar-ls      ->  bar ls
    
To use git infrastructure use `up`, `down` and `ls` with `-git` flag.

## Specs

Specs is just filename->id maps. Specs can be exported and imported by 
`bar spec export` and `bar spec import`. Specs can be generated locally 
or uploaded to bard.

### Windows, brains and two buttons 

Under windows you can use one-button no-brain solution. Just go to HTTP 
host:port and follow instructions in browser.

