# BLOB Shadow 

BLOB shadow replaces BLOB in working tree and git index.

    BAR:SHADOW
    
    version <version>
    id <hash>
    size <size>
    (ending \n)
    [optional]
    id <hash>
    size <size>
    offset <offset>
    (ending \n)

Shadow file always starts with `BAR:BLOB:SHADOW` header. `<version>` is just 
for compatibility check. `<hash>` is regular SHA3-256 hash in hex notation. 
`<size>` is regular BLOB size in bytes. For example:

    BAR:SHADOW
    
    version 0.1.0
    id 3339defdb3e5b3a2a71941b6b2bbdf7bb6525b61ba7eafb2cdb47428b3b65110
    size 52428800
    
    
    id 3339defdb3e5b3a2a71941b6b2bbdf7bb6525b61ba7eafb2cdb47428b3b65110
    size 52428800
    offset 0
    
    ---

If chunks is skipped:

    BAR:SHADOW
        
    version 0.1.0
    id 3339defdb3e5b3a2a71941b6b2bbdf7bb6525b61ba7eafb2cdb47428b3b65110
    size 52428800
    
    ---

    
    
