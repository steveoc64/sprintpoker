all: build

build:
	go generate .
	go build .
	ls -ltra poker

run: build
	./poker
