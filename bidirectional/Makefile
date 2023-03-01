.PHONY: all
all: build run

.PHONY: run
run:
	./app put socks 2
	./app put socks 1
	./app get socks

.PHONY: build
build:
	go build -o app
	go build -o counter-go-grpc ./plugin-go-grpc

.PHONY: pbuf
pbuf:
	protoc -I proto/ proto/kv.proto --go_out=plugins=grpc:proto/

.PHONY: clean
clean:
	rm -f ./app
	rm -f ./counter-go-grpc
	rm -f ./kv_store_*
