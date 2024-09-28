FROM golang:alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /app/app .

RUN grpc_cpp_plugin=/usr/local/bin/grpc_cpp_plugin \
    && go get -u google.golang.org/grpc/cmd/protoc-gen-go \
    && go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

CMD ["./app"]