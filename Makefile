.PHONY:build dist-wasm test dist watch go-install-deps

default: watch

build: test build-wasm

build-wasm:
	mkdir -p dist
	rm -f dist/main.wasm
	@GOOS=js GOARCH=wasm go build -o dist/main.wasm ./ui

build-hchan:
	mkdir -p dist
	go build -o dist/hchan ./hchan/cmd/hchan

test:
	echo "no tests so far"

wasm-dist: test build-wasm
	cp assets/* dist/

go-install-deps:
	GO111MODULE=off go get -u github.com/radovskyb/watcher/cmd/watcher
	GO111MODULE=off go get -u github.com/andrebq/wfsd


# run targets

watch:
	watcher -cmd="make wasm-dist" -startcmd=true -dotfiles=false -keepalive -ignore "./dist"

run-hchan: build-hchan
	./dist/hchan -bind "127.0.0.1:8082"

run-serve:
	cd dist; wfsd -p "127.0.0.1:8081"