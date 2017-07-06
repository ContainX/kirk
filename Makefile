GO_FMT     = gofmt -s -w -l .
BUILD_TIME = `date +%FT%T%z`

all: deps test compile

deps:
	go get ./...

test:
	go test ./...

compile:
	go build

format:
	$(GO_FMT)
