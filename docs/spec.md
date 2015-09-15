# BLOB Shadow 

BLOB shadow replaces BLOB in working tree and git index.

    BAR:BLOB:SHADOW
    
    <version>
    <hash>
    <size>
    (ending \n)

Shadow file always starts with `BAR:BLOB:SHADOW` header. `<version>` is just 
for compatibility check. 

`<hash>` must be formed as following:

    sha256:...
    
`<size>` is regular BLOB size in bytes.
