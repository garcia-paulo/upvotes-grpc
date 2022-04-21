gen: 
	protoc --proto_path=proto --go_out=proto/gen --go_opt=paths=source_relative --go-grpc_out=proto/gen --go-grpc_opt=paths=source_relative proto/*.proto
clean: 
	rm -f pb/*.go
build: 
	go build main.go