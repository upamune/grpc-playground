.PHONY: protoc
protoc:
	protoc --proto_path=$(GOPATH)/src:.  --go_out=plugins=grpc:. ./protobuf/*.proto

