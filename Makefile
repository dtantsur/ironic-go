all: build

# Run tests
test: fmt vet
	go test ./... -coverprofile cover.out

build: fmt vet
	go build ./...

fmt:
	go fmt ./...

vet:
	go vet ./...
