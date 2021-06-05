CC=go
FMT=gofmt
NAME=go-lnmetrics
BASE_DIR=/script

default: fmt
	$(CC) build -o $(NAME) cmd/go-metrics-reported/main.go

fmt:
	$(CC) fmt ./...

check:
	echo "Nothings yet"
