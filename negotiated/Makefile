.PHONY: all
all: build run

.PHONY: run
run:
	./app --plugin=3 put hello "big wide world"
	./app --plugin=2 get hello

.PHONY: build
build:
	go build -o app
	go build -o kv-plugin ./plugin-go

.PHONY: pbuf
pbuf:
	protoc ./proto/kv.proto \
		--go_out=. \
		--go-grpc_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--proto_path=.

.PHONY: clean
clean:
	rm -f ./app
	rm -f ./kv-plugin
	rm -f ./kv_store_*
