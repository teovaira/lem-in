.PHONY: build test vet race clean run

build:
	go build -o lem-in ./cmd/lem-in

test:
	go test ./...

vet:
	go vet ./...

race:
	go test -race ./...

clean:
	rm -f lem-in

run:
	go run ./cmd/lem-in $(ARGS)
