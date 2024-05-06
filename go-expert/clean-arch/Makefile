install-cli-tools:
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/99designs/gqlgen@latest
	go install github.com/ktr0731/evans@latest

genwire:
	cd cmd/ordersystem/ && wire

gengraphql:
	gqlgen

genprotobuf:
	cd internal/infra/grpc/protofiles && protoc --go_out=../pb --go_opt=paths=source_relative --go-grpc_out=../pb --go-grpc_opt=paths=source_relative *.proto

