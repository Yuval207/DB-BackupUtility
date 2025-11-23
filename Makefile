BINARY_NAME=dbbackup

build:
	go build -o $(BINARY_NAME) ./cmd/dbbackup

run: build
	./$(BINARY_NAME)

test:
	go test -v ./...

clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f *.sql *.gz *.archive

deps:
	go mod download
