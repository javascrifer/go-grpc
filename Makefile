.PHONY:
pb-generate:
	sh ./scripts/proto-generate.sh

.PHONY:
run-client:
	go run ./cmd/client/main.go

.PHONY:
run-server:
	go run ./cmd/server/main.go