.PHONY:build dist-wasm test dist watch go-install-deps

default: watch

build: test build-wasm

build-runner:
	mkdir -p dist
	go build -o dist/runner ./runner/cmd/runner

build-wasm:
	mkdir -p dist
	rm -f dist/main.wasm
	@GOOS=js GOARCH=wasm go build -o dist/main.wasm ./ui

build-dummycli-wasm:
	docker build -t rustup_wasm:latest ./context -f ./context/RustWasm.Dockerfile
	docker run --rm -v ${PWD}/runner/rust/dummycli:/usr/src/dummycli -w /usr/src/dummycli rustup_wasm:latest cargo build --release --target=wasm32-unknown-unknown
	cp ${PWD}/runner/rust/dummycli/target/wasm32-unknown-unknown/release/dummycli.wasm dist/

build-hchan:
	mkdir -p dist
	go build -o dist/hchan ./hchan/cmd/hchan

test:
	go test ./runner/...
	go test ./hchan/...

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

run-runner: build-runner build-dummycli-wasm
	./dist/runner -entry app_main ./dist/dummycli.wasm