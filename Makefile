generate-mocks:
	mockgen -source=./ratelimiter/adapter/storage_adapter.go -destination ./ratelimiter/mocks/storage_adapter.go -package mocks
	mockgen -source=./ratelimiter/responsewriter/response_writer.go -destination ./ratelimiter/mocks/response_writer.go -package mocks

test:
	go test ./... -v
