.PHONY: proto mocks build

proto:
	protoc \
		--proto_path=api/proto \
		--go_out=internal/pb \
		--go-grpc_out=internal/pb \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		rates.proto

mocks:
	mockery

build:
	go build -o bin/usdt-rate-service ./cmd/main.go

test:
	go test ./... 

docker-build:
	docker build -t usdt-rate-service:latest .

run:
	docker-compose up --build

lint:
	golangci-lint run 