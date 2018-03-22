default: cover
cover:
	@go test ./... -coverprofile=coverage.out
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out
	@rm coverage.out
