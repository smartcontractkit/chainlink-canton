package openapi

//go:generate go tool oapi-codegen -generate models,client -config ./cfg.yaml -o ./gen/tokenMetadataV1/tokenMetadataV1.gen.go -package tokenMetadataV1 ./src/token-metadata-v1.yaml
//go:generate go tool oapi-codegen -generate models,client -config ./cfg.yaml -o ./gen/transferInstructionV1/transferInstructionV1.gen.go -package transferInstructionV1 ./src/transfer-instruction-v1.yaml
//go:generate go tool oapi-codegen -generate models,client -config ./cfg.yaml -o ./gen/scanProxy/scanProxy.gen.go -package scanProxy ./src/scan-proxy.yaml
