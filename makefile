all: build

build:
	go generate .
	go build .
	ls -ltra poker

dist: build
	go run dist.go
	doas cp poker /usr/local/bin
	doas cp config.toml /usr/local/etc/poker.toml
	doas mkdir -p /var/www/poker
	doas cp dist/* /var/www/poker

run: dist
	./poker

dev: build
	POKER_DEV=1 ./poker

bastille: dist
	doas bastille cp poker ./poker /usr/local/bin
	doas bastille cp poker ./config.toml /usr/local/etc/poker.toml
	cd dist && doas bastille cp poker . /var/www/poker
