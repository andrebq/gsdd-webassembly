.PHONY:build dist-wasm test dist watch go-install-deps

watch:
	watcher -cmd="make dist" -startcmd=true -dotfiles=false -keepalive -ignore "./dist"

build: test build-wasm

build-wasm:
	mkdir -p dist
	rm -f dist/main.wasm
	export GOOS=js
	export GOARCH=wasm
	@go build -o dist/main.wasm ./ui

test:
	echo "no tests so far"

dist: build test
	cp assets/* dist/

go-install-deps:
	GO111MODULE=off go get -u github.com/radovskyb/watcher/cmd/watcher