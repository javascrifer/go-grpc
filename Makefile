.PHONY:
pb-generate:
	sh ./scripts/proto-generate.sh

.PHONY:
run-client:
	go run ./cmd/client/main.go

.PHONY:
run-server:
	go run ./cmd/server/main.go

.PHONY:
run-server-docker:
	docker container rm -f go-grpc
	docker container run -d -p 50051:50051 --name go-grpc nikasdocker/go-grpc-server

.PHONY:
build-server-image:
	rm -f ./bin/grpc-server/server
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/grpc-server/server ./cmd/server/main.go
	docker image build -t nikasdocker/go-grpc-server .
