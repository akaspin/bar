
/** Bard server info */
struct ServerInfo {
    /** HTTP endpoint (http://bard.served:3000/v1) */
    1: string httpEndpoint,

    /** Thrift rpc endpoints (tcp://bard.served:3000) */
    2: list<string> rpcEndpoints,

    /** Preferred chunk size. */
    3: i64 chunkSize,

    /** Preferred max connections for client */
    4: i32 maxConn,

    /** Thrift client buffer size */
    5: i32 bufferSize,
}

/** SHA3-256 */
typedef binary ID

/**
* Info about data entity
**/
struct DataInfo {
    1: ID id,
    /** data size */
    2: i64 size,
}

/**
* Blob manifest.
**/
struct Manifest {
    1: DataInfo info,
    2: list<Chunk> chunks,
}

/**
* Chunk info
**/
struct Chunk {
    1: DataInfo info,
    3: i64 offset,
}

struct Spec {
    1: ID id,
    2: i64 timestamp,
    3: map<string, ID> blobs,
    4: list<string> removes,
}


service Bar {

////

    /**
    * Creates new upload on bard and returns missing chunks.
    **/
    list<DataInfo> CreateUpload (

        /** upload id */
        1: binary id,

        /** requested manifests */
        2: list<Manifest> manifests,
    ),

    /**
    * Upload BLOB chunk
    **/
    void UploadChunk(
        /** upload id */
        1: ID uploadId,
        2: DataInfo info,
        3: binary Body,
    ),

    /**
    * Finish upload BLOB
    **/
    void FinishUploadBlob (
        1: binary uploadId,
        2: ID blobId,
        3: list<binary> tags,
    ),

    /**
    * Mark upload as finished. This action will
    * immediately remove all upload data.
    **/
    oneway void FinishUpload (
        1: binary uploadId,
    ),

////

    /**
    * Tag blobs. Untagged blobs will be removed by GC.
    **/
    void TagBlobs (
        1: list<ID> ids,
        2: list<binary> tags,
    ),

    void UntagBlobs (
        1: list<ID> ids,
        2: list<binary> tags,
    ),

///

    list<ID> IsBlobExists (
        1: list<ID> ids
    ),

////

    list<Manifest> GetFetch (
        1: list<ID> ids,
    )

    binary FetchChunk (
        1: ID blobID,
        2: ID chunkID,
    )

////

    void UploadSpec (
        1: Spec spec,
    )

    Spec FetchSpec (
        1: ID id,
    )
}