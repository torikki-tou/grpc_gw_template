proto-vendor:
    mkdir -p protovendor/googleapis/google/api
    curl -L https://raw.githubusercontent.com/googleapis/googleapis/1b5f8632487bce889ce05366647addc6ef5ee36d/google/api/annotations.proto -o protovendor/googleapis/google/api/annotations.proto
    curl -L https://raw.githubusercontent.com/googleapis/googleapis/1b5f8632487bce889ce05366647addc6ef5ee36d/google/api/http.proto -o protovendor/googleapis/google/api/http.proto

gen:
    docker run --rm -v $(pwd):/workspace my-bufgen:latest generate

docker-build:
    docker build -t my-bufgen:latest .