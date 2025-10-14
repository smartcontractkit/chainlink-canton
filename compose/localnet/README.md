# Localnet

This contains a Docker compose setup to quickly start a localnet environment including an SV and 5 participant nodes.

## Start
To start:
```
docker compose up -d
```

## Console
To open the console:
```
docker compose run --rm console
```

## Access APIs

A Nginx server will be started listening on post `8080` to allow for easy access to all node's ledger and admin APIs:

Routing is set up to allow access to the following nodes:
- `sv`
- `participant1`
- `participant2`
- `participant3`
- `participant4`
- `participant5`

To e.g. access the SV instead of Participant 1, replace `participant1` with `sv` in the following URLs:

### Ledger API (HTTP)

http://participant1.json-ledger-api.localhost


### Ledger API (gRPC)

grpc://participant1.grpc-ledger-api.localhost


## Admin API (grpc)

grpc://participant1.admin-api.localhost
