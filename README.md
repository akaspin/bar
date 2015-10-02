# Bar

**This software is pre-alpha**

Bar was developed to solve one simple problem. Large BLOBs can be stored in 
working tree. But our developers do not need these BLOBs. Or they are need 
some of them.

We use git. Bar designed to to minor additions to the workflow. Also bar can 
work without git.

## Basic usage

> For now assume what `bard` is deployed and listening at `:3000`.

To set the bar in the repository use `barc git-init` command. Only one 
important flag "endpoint" is just HTTP endpoint of `bard` server. 
    
    $ barc git-init -endpoint=http://localhost:3000/v1
    
Bar-tracked BLOBs are defined by git attributes.

    # .gitattributes
    my/blobs    filter=bar  diff=bar
    
After that you can add, commit and push as usual. Bar will upload new and 
changed BLOBs to `bard` server on commit.
 
    $ echo "test" > my/blobs/test.txt
    $ git add -A
    $ git commit -m "initial commit"  # <-- Here BLOBs will be uploaded
    ...

For this moment all is simple. To transform BLOBs to shadows 
use `git bar-squash`:

    $ git bar-squash
    ...
    
    $ git status 
    On branch master
    nothing to commit, working directory clean
    
    $ cat my/blobs/test.txt
    BAR:SHADOW
    
    version 0.1.0
    id 309a3490190131517180e3827398a665c1eef2b9b2b41108a08a59f4bb15d301
    size 4
    
    
    id 309a3490190131517180e3827398a665c1eef2b9b2b41108a08a59f4bb15d301
    size 3
    offset 0
    
Now BLOB transformed to shadow. Shadow is small text manifest describing BLOB. 
Bar installs own `filter` and `diff` to git repo. This leverages git to store 
shadows in index.

    $ git ls-files -s --cache my/blobs/test.txt
    100644 c8c0b7267bf989a8d15c885dcdbafb10b60d244e 0	my/blobs/test.txt
    
    $ git cat-file -p c8c0b7267bf989a8d15c885dcdbafb10b60d244e
    BAR:SHADOW
        
    version 0.1.0
    id 309a3490190131517180e3827398a665c1eef2b9b2b41108a08a59f4bb15d301
    size 4
    
    
    id 309a3490190131517180e3827398a665c1eef2b9b2b41108a08a59f4bb15d301
    size 3
    offset 0
    
To upload BLOBs without squashing them use `git bar-up`
    
Also on checkout all non-existent BLOBs will be become stored as shadows. To 
get them back use `git bar-down`.

    $ git bar-down my/blobs/test.txt
    ...    

To check status of bar-tracked BLOBs use `git bar-ls`
 
    $ git bar-ls
    NAME                BLOB    SYNC    ID                  SIZE
    my/blobs/test.txt   yes     yes     309a349019013151    4
    
    $ git bar-squash
    $ git bar-ls
    NAME                BLOB    SYNC    ID                  SIZE
    my/blobs/test.txt   no      yes     309a349019013151    4
    
## Installation

To install just grab latest binaries archive from releases and unpack 
somewhere to `PATH`.

## `bard` BLOB server

For `bard` is quick and dirty HTTP file server with block backend:

    $ bard \
        -bind=0.0.0.0:3000 \
        -storage-block-root=bard/blobs \
        -logging-level=DEBUG
    DEBUG server.go:23: serving at http://0.0.0.0:3000/v1

## Git-less usage

Bar doesn't require git to work. All `git bar-*` commands is just git aliases:

    git bar-squash  ->  barc up
    git bar-up      ->  barc up -noop
    git bar-down    ->  barc down
    git bar-ls      ->  barc ls
    
To use git infrastructure with `up`, `down` and `ls` use `-git` flag.
