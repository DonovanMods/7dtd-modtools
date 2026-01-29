.PHONY: default check test tidy format upgrade update clean build package build-darwin build-linux build-windows

default: check

tidy:
	go mod tidy

format: tidy
	trunk fmt

check: format
	trunk check

test:
	go test -v -cover ./...

upgrade: tidy
	go get -u
	trunk upgrade

update: upgrade

clean:
	rm -f bin/*
	rm -f pkg/*
	go clean -i -cache -testcache

build: check test build-darwin build-linux build-windows

package: build

build-darwin:
	GOOS=darwin GOARCH=arm64 go build -o bin/7dmt-darwin-arm64 ./main.go
	zip -r pkg/7dmt-darwin-arm64.zip bin/7dmt-darwin-arm64

build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/7dmt-linux-amd64 ./main.go
	zip -r pkg/7dmt-linux-amd64.zip bin/7dmt-linux-amd64

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/7dmt-windows-amd64.exe ./main.go
	zip -r pkg/7dmt-windows-amd64.exe.zip bin/7dmt-windows-amd64.exe
