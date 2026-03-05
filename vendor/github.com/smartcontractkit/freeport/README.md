This is a port of inactive Hashicorp freeport library.

It's primary function is to reduce the initialization panics in the original implementation. We observed many flakey CI runs because of
`cannot allocate port block` in the [alloc func](https://github.com/hashicorp/consul/blob/f3c5d71cbf7944fa99df9bf4f33fc213223f170d/sdk/freeport/freeport.go#L274)

The root cause is many concurrent memory spaces (ie concurrent `go test`) that collide because the [hardcode block size is far too large](https://github.com/hashicorp/consul/blob/f3c5d71cbf7944fa99df9bf4f33fc213223f170d/sdk/freeport/freeport.go#L102) for our use case, limiting the [number of port blocks to < 30](https://github.com/hashicorp/consul/blob/f3c5d71cbf7944fa99df9bf4f33fc213223f170d/sdk/freeport/freeport.go#L32).

This migration 
- decreases the [port block size](https://github.com/hashicorp/consul/blob/f3c5d71cbf7944fa99df9bf4f33fc213223f170d/sdk/freeport/freeport.go#L102) to 128, which by inspection of tests is still plenty big enough
- increases the [max number of blocks](https://github.com/hashicorp/consul/blob/f3c5d71cbf7944fa99df9bf4f33fc213223f170d/sdk/freeport/freeport.go#L32) to 512
- changes [the search in alloc](https://github.com/hashicorp/consul/blob/f3c5d71cbf7944fa99df9bf4f33fc213223f170d/sdk/freeport/freeport.go#L264) to try more options before panic'ing

it also drops depreciated functions and window-specific builds, since they are not relevant.