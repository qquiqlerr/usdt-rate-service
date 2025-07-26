.PHONY: proto

proto:
	protoc \
		--proto_path=api/proto \
		--go_out=internal/pb \
		--go-grpc_out=internal/pb \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		rates.proto