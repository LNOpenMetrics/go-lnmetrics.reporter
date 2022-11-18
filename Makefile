CC=go
FMT=gofmt
NAME=go-lnmetrics
BASE_DIR=/script
OS=linux
ARCH=386
ARM=
GORPC_COMMIT=58377b3f766604dc383e9cc638e2239bd5bf9c49

default: fmt lint
	$(CC) build -o $(NAME) cmd/go-lnmetrics.reporter/main.go

fmt:
	$(CC) fmt ./...

lint:
	golangci-lint run

check:
	$(CC) test -v ./...

check-dev:
	richgo test ./... -v

build:
	env GOOS=$(OS) GOARCH=$(ARCH) GOARM=$(ARM) $(CC) build -o $(NAME)-$(OS)-$(ARCH) cmd/go-lnmetrics.reporter/main.go

dep:
	$(CC) get -u github.com/LNOpenMetrics/lnmetrics.utils
	$(CC) get -u github.com/vincenzopalazzo/cln4go@$(GORPC_COMMIT)
	$(CC) mod vendor
