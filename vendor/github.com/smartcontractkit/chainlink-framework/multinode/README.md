# MultiNode

Enables the use of multiple RPCs in chain integrations. Performs critical health checks,
RPC selection, node metrics, and is used to send transactions to all RPCs and aggregate results.
MultiNode is used by all other components which require reading from or writing to the chain.

## Components

### RPCClient
Interface for wrapping an RPC of any chain type. Required for integrating a new chain with MultiNode.

### Node
Wrapper of an RPCClient with state and lifecycles to handle health of an individual RPC.

### MultiNode
Manages all nodes performing node selection and load balancing, health checks and metrics, and running actions across all nodes.

### Poller
Used to poll for new heads and finalized heads within subscriptions.

### Transaction Sender
Used to send transactions to all healthy RPCs and aggregate the results.