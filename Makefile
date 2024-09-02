compile: 
	protoc -I=api/v1 --go_out=api/v1 --go_opt=paths=source_relative api/v1/*.proto

test:
	go test -race ./...