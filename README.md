
#### TO DO
- [ ] SSL Configure gRPC
- [ ] Writing Test
- [ ] Create Container with options
- [x] And more ...



#### Generate pb.proto file
     protoc -I=./api/proto/v1 --go_out=plugins=grpc:./pkg/api/v1/ ./api/proto/v1/scaler-api.proto
