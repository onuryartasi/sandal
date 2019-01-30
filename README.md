
#### TO DO
- [ ] Writing Test
- [ ] SSL Configure gRPC
- [ ] Create Container with options
- [ ] Docker Metrics data
- [ ] Cobra cli tool
- [x] And more ...



#### Generate pb.go file
     protoc -I=./api/proto/v1 --go_out=plugins=grpc:./pkg/api/v1/ ./api/proto/v1/scaler-api.proto
