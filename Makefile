GOLANG_TARGET := protogen/golang
OPENAPI_TARGET := protogen/openapi

protoc:
	 rm -rf ${GOLANG_TARGET}/* &&  rm -rf ${OPENAPI_TARGET}/* && docker run --rm -v ${PWD}:${PWD} -w ${PWD} rvolosatovs/protoc \
                    --proto_path=proto \
                    --go_out=${GOLANG_TARGET} \
                    --go_opt=paths=source_relative \
                    --go-grpc_out=${GOLANG_TARGET} \
                    --go-grpc_opt=paths=source_relative \
                    --grpc-gateway_out=${GOLANG_TARGET} \
                    --grpc-gateway_opt paths=source_relative \
                    --grpc-gateway_opt generate_unbound_methods=true \
                    --openapiv2_out ${OPENAPI_TARGET} \
                    --openapiv2_opt allow_merge=true \
                    proto/**/*.proto

google:
	curl -L https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto -o proto/google/api/annotations.proto && \
    curl -L https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto -o proto/google/api/http.proto

run_grpc:
	go run cmd/grpc/main.go

run_rest:
	go run cmd/rest/main.go