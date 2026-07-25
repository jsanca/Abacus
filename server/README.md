# Server

The server will provide the Go HTTP API using Chi. It will own the immutable operation registry, generate the operation manifest, validate requests, execute operations, and shut down gracefully.

No backend implementation is included in this bootstrap task.
