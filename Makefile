proto-gen:
	protoc --proto_path=proto/ --go_out=paths=source_relative,plugins=grpc:./pb proto/*/*.proto

server:
	go run cmd/main.go