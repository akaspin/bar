# Bar

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
 
    $ echo "test" > /my/blobs/test.txt
    $ git add -A
    $ git commit -m "initial commit"  # <-- Here BLOBs will be uploaded
    ...

For this moment all is simple. 
