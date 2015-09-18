# barctl subcommands

git-clean
: git clean. Takes stream for `stdin` and prints shadow manifest to `stdout`. 

git-cat
: cat for git diff. Prints `BAR-SHADOW-BLOB-<is-blob> <filename> <id>`. Used 
for `pre-commit` hook.

upload
: upload blobs to bard. Takes list of filenames with optional hashes from 
`stdin`. 
