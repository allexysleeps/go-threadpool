test:
	go test -race -v ./...
format:
	gofumpt -l -w .