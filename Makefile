test:
	@go test

bench:
	@go test -bench=".*"

coverage:
	@go test -covermode=count
