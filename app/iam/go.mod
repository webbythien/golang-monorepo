module github.com/monorepo/app/iam

go 1.24.1

replace monorepo/api => ../../api

require monorepo/api v0.0.0-00010101000000-000000000000

require (
	connectrpc.com/connect v1.18.1 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250115164207-1a7da9e5054f // indirect
	google.golang.org/grpc v1.71.1 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
