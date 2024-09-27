package tasks

//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 vet --file=./internal/database/erd/sqlc.yaml
//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate --file=./internal/database/erd/sqlc.yaml
//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 diff ./database/erd --file=./internal/database/erd/sqlc.yaml
