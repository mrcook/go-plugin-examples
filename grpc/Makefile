.PHONY: all
all: build-go

.PHONY: build-go
build-go:
	go build -o app
	go build -o kv-go-grpc ./plugin-go-grpc
	go build -o kv-go-netrpc ./plugin-go-netrpc

.PHONY: pbufs
pbufs: pbufs-go pbufs-py

.PHONY: pbufs-go
pbufs-go:
	protoc -I proto/ proto/kv.proto --go_out=plugins=grpc:proto/

.PHONY: pbufs-py
pbufs-py:
	python -m grpc_tools.protoc -I ./proto/ --python_out=./plugin-python/ --grpc_python_out=./plugin-python/ ./proto/kv.proto

.PHONY: clean
clean:
	rm -f ./app
	rm -f ./kv-go-grpc
	rm -f ./kv-go-netrpc
	rm -f ./kv_grpc_*
	rm -f ./kv_rpc_*
	rm -f ./kv_py_*
