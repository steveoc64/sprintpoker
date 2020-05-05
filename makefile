all: build

build:
	go generate .
	go build .
	ls -ltra poker

dist: build
	go run dist.go
	doas cp poker /usr/local/bin
	#doas cp config.toml /usr/local/etc/poker.toml
	#doas vi /usr/local/etc/poker.toml
	doas mkdir -p /var/www/poker
	doas cp dist/* /var/www/poker

run: dist
	./poker

dev: build
	POKER_DEV=1 ./poker

# Deploy poker server to Bastille container (FreeBSD jail)
bastille: build
	go run dist.go
	cd dist && doas bastille cp poker . /var/www/poker
	doas bastille cmd poker killall poker
	doas bastille cp poker ./poker /usr/local/bin

bastille-config:
	doas bastille cmd poker vi /usr/local/etc/poker.toml
	doas bastille cmd poker killall poker

bastille-log:
	doas bastille cmd poker tail -f /var/log/poker.log

#only need to do this once
bastille-boot: bastille
	doas bastille cmd poker mkdir -p /var/www/poker
	#doas bastille cp poker ./config.toml /usr/local/etc/poker.toml
	doas bastille cmd poker vi /usr/local/etc/poker.toml
	doas bastille cmd poker daemon -o /var/log/poker.log -p /var/run/poker.pid -R 1 poker

tinydist:
	tinygo build -o dist/tinymain.wasm -target wasm ./main_wasm.go
