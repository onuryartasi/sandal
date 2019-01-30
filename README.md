protoc -I=./api/proto/v1 --go_out=plugins=grpc:./pkg/api/v1/ ./api/proto/v1/scaler-api.proto
