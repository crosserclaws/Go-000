.PHONY: all
all: proto wire build

.PHONY: build
build: client server

.PHONY: client server
client server:
	go build -o bin/$@ ./cmd/$@

.PHONY: proto
proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/video/v1/videoinfo.proto

.PHONY: wire
wire:
	wire ./cmd/server

.PHONY: clean
clean:
	rm -r bin/*