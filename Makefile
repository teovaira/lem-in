.PHONY: build test vet race clean run fmt

build:
	go build -o lem-in ./cmd/lem-in

test:
	go test ./...

vet:
	go vet ./...

race:
	go test -race ./...

fmt:
	gofmt -w .

clean:
	rm -f lem-in

run:
	go run ./cmd/lem-in $(ARGS)
