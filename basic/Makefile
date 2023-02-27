.PHONY: all
all: build run

build:
	go build -o app .
	go build -o hello_plugin ./hello_plugin_example

run:
	./app

clean:
	rm -f ./app
	rm -f ./hello_plugin
