grpc_proto_gen:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api_proto/*.proto
	protoc -I=. --js_out=import_style=commonjs:frontend/src/ --grpc-web_out=import_style=commonjs,mode=grpcwebtext:frontend/src/ api_proto/*.proto
