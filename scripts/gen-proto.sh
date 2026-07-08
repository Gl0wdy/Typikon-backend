#!/bin/bash
set -e

SERVICES=(
    "wiki:api/proto/wiki.proto"
    "auth:api/proto/auth.proto"
    "data:api/proto/data.proto"
    "spaces:api/proto/spaces.proto"
)

GOOGLEAPIS_DIR="./third_party/googleapis"

if [ ! -d "$GOOGLEAPIS_DIR" ]; then
    echo "Downloading dependencies to $GOOGLEAPIS_DIR..."
    mkdir -p "$GOOGLEAPIS_DIR/google/api"
    curl -s -o "$GOOGLEAPIS_DIR/google/api/annotations.proto" \
        https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto
    curl -s -o "$GOOGLEAPIS_DIR/google/api/http.proto" \
        https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto
    echo "Dependencies downloaded."
fi

generate_proto() {
    local proto_path=$1

    if [ ! -f "$proto_path" ]; then
        echo "File not found: $proto_path"
        return 1
    fi

    echo "Generating: $proto_path"

    protoc \
        -I . \
        -I ./api/proto \
        -I "$GOOGLEAPIS_DIR" \
        --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
        "$proto_path"

    echo "Done: $proto_path"
}

if [ -n "$1" ]; then
    for entry in "${SERVICES[@]}"; do
        name="${entry%%:*}"
        path="${entry#*:}"
        if [ "$1" == "$name" ]; then
            generate_proto "$path"
            exit 0
        fi
    done
    echo "Unknown service: $1"
    echo "Available: ${SERVICES[@]%%:*}"
    exit 1
fi

echo "Select a service to regenerate proto:"
select entry in "${SERVICES[@]}" "All services" "Cancel"; do
    case "$entry" in
        "All services")
            for s in "${SERVICES[@]}"; do
                path="${s#*:}"
                generate_proto "$path"
            done
            break
            ;;
        "Cancel")
            echo "Cancelled."
            exit 0
            ;;
        *)
            if [ -n "$entry" ]; then
                path="${entry#*:}"
                generate_proto "$path"
                break
            else
                echo "Invalid choice, try again."
            fi
            ;;
    esac
done