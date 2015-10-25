# git flow

Git flow is quite usual.

# Setup

Configure git

    $ barctl init -endpoint=http://my.bar/v1
    git config filter.bar.clean "barctl git-clean %f"
    git config filter.bar.smudge "barctl git-smudge -endpoint=http://my.bar/v1 %f"
    git config diff.bar "barctl git-cat"
    
    $ barctl init -endpoint=http://my.bar/v1 | sh
    $ barctl install-hook -endpoint=http://my.bar/v1
    
    $ cat .git/hook/pre-commit
    barctl pre-commit -endpoint=http://my.bar/v1
    
# Usage

Add bar-controlled files to `.gitattributes`:

    /my/blobs/** filter=bar diff=bar

Add, commit and push as usual. Bar will upload new and changed BLOBs to bard 
server on commit. 

    $ echo "test" > /my/blobs/test.txt
    $ git add -A
    $ git commit -m "initial commit" # Here BLOBs will be uploaded
    ...
    
But on checkout all non-existent BLOBs will be become *Shadows*. Shadow file 
has same name. But by design it is just small text manifest. To get normal 
BLOB instead shadows use `barctl blow`

    $ barctl blow -git /my/blobs/*
    
`barctl blow` scans all shadows and replaces them with downloaded BLOBs.

To replace BLOB with shadow use `barctl squash`.

    $ barctl squash -git /my/blobs/*
    
`barctl squash` will replace all BLOBs with them shadows. If BLOB is not 
exists on bard server - it will be uploaded before replace.

To get status of blobs use `barctl status`:

    $ barctl status -git /my/blobs/*
    
    FILE                SHADOW      REMOTE
    my/blobs/test.txt   no          yes
    

