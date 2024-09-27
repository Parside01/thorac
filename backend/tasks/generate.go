package tasks

//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 vet --file=./database/erd/sqlc.yaml
//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate ./database/erd --file=./database/erd/sqlc.yaml
//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 diff ./database/erd --file=./database/erd/sqlc.yaml
