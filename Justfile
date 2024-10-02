default: check

## Main Commands

build: fmt clean test build-win

clean: clean-bin tidy
    go clean -i -cache -testcache

## Supporting Commands

tidy:
    go mod tidy

fmt: tidy
    trunk fmt

fmt-all: tidy
    trunk fmt --all

check: fmt
    trunk check

check-all: fmt-all
    trunk check --all

test:
    go test ./7dmt/...

clean-bin:
    rm -f bin/*

update: upgrade
upgrade: tidy
    go get -u
    trunk upgrade

## Build sub-commands

build-win:
    GOOS=windows GOARCH=amd64 go build -o "bin/$(basename ${PWD}).exe" ./main.go

## Git Hooks
pre-commit: clean check test
    git add go.mod go.sum
