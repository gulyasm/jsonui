default: build

installdep:
	@go get github.com/golang/lint/golint
	@go get

build: installdep
	@go fmt
	@go vet
	@golint
	@go test ./...
	@go build

cover:
	@go test ./... -coverprofile=coverage.out
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out
	@rm coverage.out
