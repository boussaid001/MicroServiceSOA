module github.com/boussaid001/go-microservices-project/services/grpc-service

go 1.20

require (
	github.com/boussaid001/go-microservices-project/proto v0.0.0-00010101000000-000000000000
	github.com/lib/pq v1.10.9
	google.golang.org/grpc v1.58.3
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.15.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)

replace github.com/boussaid001/go-microservices-project/proto => ../../proto
