module github.com/smartcontractkit/chainlink-canton/eds

go 1.25.6

replace github.com/digital-asset/dazl-client/v8 => github.com/noders-team/dazl-client/v8 v8.7.1-2

require (
	github.com/digital-asset/dazl-client/v8 v8.7.1
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/gorilla/mux v1.8.1
	github.com/stretchr/testify v1.11.1
	google.golang.org/grpc v1.78.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260114163908-3f89685c29c3 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/smartcontractkit/chainlink-canton => ..
