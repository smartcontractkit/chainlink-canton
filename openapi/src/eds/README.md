# CCIP Explicit Disclosure APIs

Since users are required to interact with contracts that they're not a stakeholder of, Explicit Disclosure APIs are
provided to users in order to grant them access. Users are required to query these APIs before sending/executing
messages.

While some of these APIs will be run by Chainlink Labs, **in the case of third-party token pools / CCVs / Executors, the
operator of these is required to run their own API and make it available to users.**

## Available APIs

### Global CCIP API

The Global CCIP API provides endpoints for users to query explicit disclosures for the global CCIP contracts, such as:

- FeeQuoter
- OnRamp
- OffRamp
- PerPartyRouterFactory
- RMNRemote
- TokenAdminRegistry

This API is only run by Chainlink Labs, and is made available for all users.

#### PerPartyRouterFactory

> Endpoint: `POST /ccip/v1/global/perPartyRouter/factory`

This endpoint allows users to get an explicit disclosure for the CCIP per-party-router-factory.

It is required for a party to instantiate a per-party-router, which is required for sending and executing messages via
CCIP.

#### Look up a Token Pool on the Token Admin Registry

> Endpoint: `GET /ccip/v1/global/tokenAdminRegistry/token/{instrumentId}`

This allows the user to look up a token pool for a given instrument.

If the provided Instrument has a registered token pool, the API will return the instance address of the token pool.

If the token pool's operator has registered and API endpoint with the global EDS registry, then this will also return
the base URL for the token pool's API, which can be used to query for explicit disclosures from the token pool.

#### Send

> Endpoint: `POST /ccip/v1/global/message/send`

This endpoint allows the user to query the necessary disclosures for sending a message via CCIP.

#### Execute

> Endpoint: `POST /ccip/v1/global/message/execute`

This endpoint allows the user to query the necessary disclosures for executing a message via CCIP.

### Token Pool API

The Token Pool API is provided and ran by the operator of a token pool, and provides endpoints for users to query
explicit disclosures for a specific token pool contract.

#### Send

> Endpoint: `POST /ccip/v1/external/tokenPool/{address}/send`

Send accepts the outgoing message including a token transfer.

It returns the necessary disclosures and inputs for the token pool, as well as a list of required CCVs that this token
transfer requires, which the user will need to query before actually sending the message.

#### Execute

> Endpoint: `POST /ccip/v1/external/tokenPool/{address}/execute`

Execute accepts the encoded message that is to-be-executed on Canton.

It returns the necessary disclosures and input for the token pool.

### CCV API

The CCV API is provided and ran by the operator of a CCV, and provides endpoints for users to query explicit disclosures
for a specific CCV contract.

#### Send

> Endpoint: `POST /ccip/v1/external/ccv/{address}/send`

Send accepts the outgoing message, as well as a list of all CCVs that the sender and token pool require for this
message.

It returns the necessary disclosures and inputs for the CCV.

#### Execute

> Endpoint: `POST /ccip/v1/external/ccv/{address}/execute`

Execute accepts the encoded message that is to-be-executed on Canton.

It returns the necessary disclosures and inputs for the CCV.

### Executor API

The Executor API is provided and ran by the operator of an Executor, and provides endpoints for users to query explicit
disclosures for a specific Executor contract.

#### Send

> Endpoint: `POST /ccip/v1/external/executor/{address}/send`

Send accepts the outgoing message, as well as a list of all CCVs that the sender and token pool require for this
message.

It returns the necessary disclosures and inputs for the Executor.

#### Execute

The Executor provides no API for execution, as it has no contract involved in the execution of a message on the
destination chain.

## Sending a message

### Overview

Sending a message requires the following API calls:

1. (If the message includes as token transfer) The user needs to query the token pool API for the token's token pool to
   get the token pool disclosures along with a list of required CCVs for this token transfer.
    1. If the token pool is unknown, the user needs to query the global CCIP API to look up the token pool for the
       token from the Token Admin Registry. This might also return the base URL for the token pool API if the token pool
       operator has registered it on the global EDS registry (see [GlobalRegistry](#global-eds-registry))
2. The user needs to query the global CCIP API to get the disclosures for sending a message, providing the list of
   required CCVs from the previous step along with any sender-required CCVs that the sender-contract might require.
    1. The global CCIP API will return the disclosures for sending a message
    2. Along with an aggregated list of all CCVs that are required for this message (this includes the sender-required
       CCVs, token-pool-required CCVs, and any default/lane-mandates CCVs that the global CCIP API determines are
       required for this message based on the message content and the provided lists of required CCVs).
    3. If an Executor has been requested/specified for the message, then the global CCIP API will also return the
       address of the Executor contract.
3. The user needs to query the Executor API to get the disclosures for the Executor, providing the list of all CCVs that
   were returned by the previous step.
    1. The Executor API accepts a list of CCVs along with the message, as an Executor can limit the CCVs that can be
       used with it. This allows the Executor to provide disclosures based on the specific CCVs that will be used with
       this message or reject the request if the provided CCVs are not compatible with this Executor.
4. The user needs to query the CCV API for each CCV that was returned by the global CCIP API.
5. The user submits the transaction to the Canton participant, which includes the message and all required disclosures
   from the previous steps, along with the necessary inputs for each of the TP,CCV, and Executor contracts which have
   been returned by their corresponding APIs.

### Diagram

```mermaid
---
title: Sending a message
config:
    sequence:
        messageAlign: center
        noteAlign: left
        noteMargin: 10
    
---
sequenceDiagram
    actor User
    participant TP as Token Pool API
    participant CCIP as CCIP API
    participant Executor as Executor API
    participant CCV as CCV API
    participant Canton as Canton Participant

    opt If token transfer
        opt If the token pool is unknown
            User ->>+ CCIP: GET /ccip/v1/global/tokenAdminRegistry/token/{instrumentId}
            CCIP -->>- User: Return
            note over User, CCIP: Returns: <br/> - instanceAddress: InstanceAddress of the pool <br/> - rawInstanceAddress: RawInstanceAddress of the pool <br/> - baseURL: URL of the token pool API
        end

        User ->>+ TP: POST /ccip/v1/external/tokenPool/{address}/send
        note over User, TP: Provide: <br/> - canton2AnyMessage: the message to be sent
        TP -->>- User: Return requiredCCVs
        note over User, TP: Returns: <br/> - requiredCCVs: List of CCVs that this token transfer requires <br/> - contractId: <br/> - instanceAddress: <br/> - rawInstanceAddress: <br/> - disclosedContracts[]
    end

    User ->>+ CCIP: POST /ccip/v1/global/message/send
    note over User, CCIP: Provide: <br/> - canton2AnyMessage: the message to be sent <br/> - senderRequiredCCVs: List of CCVs that the sender requires <br/> - tokenPoolRequiredCCVs: List of CCVs that the token pool requires
    CCIP -->>- User: Return messageId
    note over User, CCIP: Returns: <br/> - ccvs[]: <br/> ---- instanceAddress: <br/> ---- rawInstanceAddress: <br/> ---- baseURL: URL of the CCV API <br/> - disclosedContracts[]
    User ->>+ Executor: POST /ccip/v1/external/executor/{address}/send
    note over User, Executor: Provide: <br/> - canton2AnyMessage: the message to be sent <br/> - ccvs[]: List of all CCVs that will be used for this message
    Executor -->>- User: Return
    note over User, Executor: Returns: <br/> - contractId: <br/> - instanceAddress: <br/> - rawInstanceAddress: <br/> - baseURL: <br/> - disclosedContracts[]

    loop For each CCV
        User ->>+ CCV: POST /ccip/v1/external/ccv/{address}/send
        note over User, CCV: Provide: <br/> - canton2AnyMessage: the message to be sent
        CCV -->>- User: Return
        note over User, CCV: Returns: <br/> - contractId: <br/> - instanceAddress: <br/> - rawInstanceAddress: <br/> - baseURL: <br/> - disclosedContracts[]
    end

    User ->> Canton: Submit Transaction
    note over User, Canton: The user submits the transaction to the Canton participant, which includes the message and all required disclosures from the previous steps. <br/> Along with the necessary inputs for each of the TP,CCV, and Executor contracts which have been returned by their corresponding APIs.
```

## Executing a message

### Overview

Executing a message requires the following API calls:

1. The user queries the global CCIP API, providing the encoded message itself and a list of addresses of all CCVs that
   have verified this message.
    1. The global CCIP API will return the disclosures for all global CCIP contracts that are relevant for execution
    2. It will also look up all provided CCVs in the global EDS registry and return the `rawInstanceAddress` and
       `baseURL` for each CCV that is found in the registry, which the user can use to query the CCV APIs for these
       CCVs.
    3. If the message includes a token transfer, and the token pool operator has registered the token pool with the
       global EDS registry (see [GlobalRegistry](#global-eds-registry)), then it will also return the
       `rawInstanceAddress` and `baseURL` for the token pool, which
       the user can use to query the token pool API for disclosures related to execution.
2. The user needs to query the token pool API for the token pool of the transferred token, providing the encoded
   message.
3. The user needs to query the CCV API for each CCV that verified this message to get the disclosures for execution,
   providing the encoded message to each CCV API.
4. The user submits the transaction to the Canton participant, which includes the message and all required disclosures
   from the previous steps, along with the necessary inputs for each of the TP and CCV contracts which have been
   returned by their corresponding APIs.

### Diagram

```mermaid
---
title: Executing a message
config:
    sequence:
        messageAlign: center
        noteAlign: left
        noteMargin: 10
---

sequenceDiagram
    actor User
    participant CCIP as CCIP API
    participant TP as Token Pool API
    participant CCV as CCV API
    participant Canton as Canton Participant
    User ->>+ CCIP: POST /ccip/v1/global/message/execute
    note over User, CCIP: Provide: <br/> - encodedMessage: the encoded message to-be-executed <br/> - ccvs[]: List of CCVs that have verified this message
    CCIP -->>- User: Return
    note over User, CCIP: Returns: <br/> - ccvs[]: <br/> ---- instanceAddress: <br/> ---- rawInstanceAddress: <br/> ---- baseURL: URL of the CCV API <br/> - tokenPool: <br/> ---- instanceAddress: <br/> ---- rawInstanceAddress: <br/> ---- baseURL: URL of the Token Pool API <br/> - disclosedContracts[]
    User ->>+ TP: POST /ccip/v1/external/tokenPool/{address}/execute
    note over User, TP: Provide: <br/> - encodedMessage: the encoded message to-be-executed <br/>
    TP -->>- User: Return
    note over User, TP: Returns: <br/> - contractId: <br/> - instanceAddress: <br/> - rawInstanceAddress: <br/> - baseURL: <br/> - disclosedContracts[]

    loop For each CCV
        User ->>+ CCV: POST /ccip/v1/external/ccv/{address}/execute
        note over User, CCV: Provide: <br/> - encodedMessage: the encoded message to-be-executed <br/>
        CCV -->>- User: Return
        note over User, CCV: Returns: <br/> - contractId: <br/> - instanceAddress: <br/> - rawInstanceAddress: <br/> - baseURL: <br/> - disclosedContracts[]
    end

    User ->> Canton: Submit Transaction
    note over User, Canton: The user submits the transaction to the Canton participant, which includes the message and verifier results and all required disclosures from the previous steps. <br/> Along with the necessary inputs for each of the TP,CCV, and Executor contracts which have been returned by their corresponding APIs.
```

## Global EDS Registry

The global EDS registry is an optional component of the CCIP Explicit Disclosure system, which serves as a registry for
CCIP-relevant contracts such as token pools and CCVs. Operators of these contracts can register their contract instances
on the global EDS registry, along with metadata such as the base URL for the contract's API, which allows users to
discover the necessary APIs to query for disclosures when interacting with these contracts.

It consists of a registration record for each contract instance:

```haskell
template ContractRegistration
    with
        instanceId : Text
        owner : Party
        
        ccipOwner : Party
        
        edsBaseUrl : Text
    where
        signatory owner
        observer ccipOwner
        ensure assertValidInstanceId instanceId
```

The operator of a contract, e.g. CCV or Token Pool, can create a registration for their contract instance on the global
EDS registry, providing the `instanceId` of the contract, and the `edsBaseUrl` for their API.

The global CCIP API will index all registration records from the global EDS registry, if multiple registrations exist
for the same contract instance, the global CCIP API will prioritize the registration that was created last.

When a user queries the global CCIP API for disclosures related to a message, the global CCIP API will look up any
relevant contract instances, such as token pools or CCVs, in the global EDS registry to find their corresponding API
base URLs, which are then returned to the user along with the disclosures from the global CCIP API.

> [!NOTE]
> The global EDS registry is entirely optional, and is not required for the functioning of CCIP itself. If a contract
> instance is not registered on the global EDS registry, users can still interact with it, but they would need to know
> the API endpoint for the contract's API through other means in order to query for disclosures related to that contract
> instance.
