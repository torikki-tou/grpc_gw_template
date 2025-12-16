FROM golang:1.24-alpine AS builder

ENV GOBIN=/tmp/bin/go-install

RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest \
    && go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest \
    && go install github.com/sudorandom/protoc-gen-connect-openapi@latest

FROM bufbuild/buf:latest

RUN apk add --no-cache protobuf

COPY --from=builder /tmp/bin/go-install /usr/local/bin

WORKDIR /workspace
