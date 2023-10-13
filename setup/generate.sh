#!/bin/bash

BASEDIR=$(dirname $(realpath "$0"))

protoc -I $BASEDIR/proto \
--go_out ./pkg/pb \
--go_opt paths=source_relative \
--go-grpc_out ./pkg/pb \
--go-grpc_opt paths=source_relative \
--grpc-gateway_out ./pkg/pb \
--grpc-gateway_opt paths=source_relative \
--openapiv2_out ./docs \
./proto/api/*.proto

