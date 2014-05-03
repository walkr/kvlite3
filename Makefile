build:
	@go build

install:
	@go install

test:
	@go test

bench:
	@go test -bench=".*"

coverage:
	@go test -covermode=count
