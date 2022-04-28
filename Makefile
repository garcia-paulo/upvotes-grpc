gen: 
	protoc --proto_path=proto --go_out=proto/gen --go_opt=paths=source_relative --go-grpc_out=proto/gen --go-grpc_opt=paths=source_relative proto/*.proto
clean: 
	rm -f proto/gen/*.go
run: 
	cd server && go run .