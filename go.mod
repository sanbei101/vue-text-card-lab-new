module github.com/sanbei101/blue-card-engine

go 1.26.5

require (
	connectrpc.com/connect v1.20.0
	github.com/phuslu/log v1.0.127
	github.com/phuslu/lru v1.0.21
	github.com/rivo/uniseg v0.4.7
	golang.org/x/image v0.43.0
	google.golang.org/protobuf v1.36.11
)

require golang.org/x/text v0.38.0 // indirect

tool (
	connectrpc.com/connect/cmd/protoc-gen-connect-go
	google.golang.org/protobuf/cmd/protoc-gen-go
)
