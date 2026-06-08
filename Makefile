.PHONY: build test run clean

build:
	go build -o bin/relay cmd/relay/main.go

test:
	go test -v -cover ./...

run:
	go run cmd/relay/main.go

clean:
	rm -rf bin/ data/
