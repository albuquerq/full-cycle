gengraphql:
	gqlgen

genprotobuf:
	cd internal/infra/grpc/protofiles && protoc --go_out=../pb --go_opt=paths=source_relative --go-grpc_out=../pb --go-grpc_opt=paths=source_relative *.proto