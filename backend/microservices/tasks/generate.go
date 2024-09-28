package tasks

//go:generate protoc --proto_path=internal/proto/proto --go_out=internal/proto/gen/go/ --go-grpc_out=internal/proto/gen/go/ internal/proto/proto/*.proto

// //go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 vet --file=./internal/database/erd/sqlc.yaml
//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate --file=./internal/database/erd/sqlc.yaml
// //go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 diff ./database/erd --file=./internal/database/erd/sqlc.yaml

// //go:generate protoc --proto_path=proto/proto --go_out=proto/gen/go/ --go-grpc_out=proto/gen/go/ proto/proto/*.proto
