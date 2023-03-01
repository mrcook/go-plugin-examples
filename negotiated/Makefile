.PHONY: all
all: build run

.PHONY: run
run:
	KV_PROTO=grpc ./app put hello "big wide world"
	KV_PROTO=netrpc ./app get hello

.PHONY: build
build:
	go build -o app
	go build -o kv-plugin ./plugin-go

.PHONY: pbuf
pbuf:
	protoc -I proto/ proto/kv.proto --go_out=plugins=grpc:proto/

.PHONY: clean
clean:
	rm -f ./app
	rm -f ./kv-plugin
	rm -f ./kv_store_*
