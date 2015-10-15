
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


/**
* Wire service
**/
service Bar {

    ServerInfo GetInfo(),

////

    /**
    * Creates new upload on bard and returns missing chunks.
    **/
    list<ID> CreateUpload (

        /** upload id (UUIDv4) */
        1: binary id,

        /** requested manifests */
        2: list<Manifest> manifests,

        /** upload TTL */
        3: i64 ttl,
    ),

    /**
    * Upload BLOB chunk
    **/
    void UploadChunk(
        /** upload id */
        1: binary uploadId,

        /** chunk id */
        2: ID chunkId,

        /** Chunk body */
        3: binary body,
    ),

    /**
    * Mark upload as finished. This action will
    * immediately remove all upload data.
    **/
    void FinishUpload (

        /** Upload id */
        1: binary uploadId,
    ),

////

    /**
    * Tag blobs. Untagged blobs will be removed by GC.
    **/
//    void TagBlobs (
//        1: list<ID> ids,
//        2: list<binary> tags,
//    ),
//
//    void UntagBlobs (
//        1: list<ID> ids,
//        2: list<binary> tags,
//    ),

///

    list<ID> GetMissingBlobIds (
        1: list<ID> ids
    ),

////

    /**
    * Get manifests by their ids
    **/
    list<Manifest> GetManifests (
        1: list<ID> ids,
    )

    /**
    * Fetch chunk from bard
    **/
    binary FetchChunk (
        /** Blob ID */
        1: ID blobID,

        /** Chunk spec */
        2: Chunk chunk,
    )

////

    void UploadSpec (
        1: Spec spec,
    )

    Spec FetchSpec (
        1: ID id,
    )
}