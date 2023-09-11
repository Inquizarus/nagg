SHELL=/bin/bash

.PHONY: test
test:
	go test -count=1 -v ./...

test_with_docker:
	podman run -v ./:/app -w /app public.ecr.aws/docker/library/golang:1.21.1 go test -count=1 -v ./...

tidy_with_docker:
	podman run -v ./:/app -w /app public.ecr.aws/docker/library/golang:1.21.1 go mod tidy

.PHONY: build_with_docker
build_with_docker:
	docker run -e CGO_ENABLED=0 -e GOOS=linux -e GPARCH=amd64 -v "$(pwd)":/app -w /app public.ecr.aws/docker/library/golang:1.21.1 go build -ldflags "-extldflags '-static'" -o ./build/nagg ./cmd/nagg

.PHONY: run_local_cmd
run_local_cmd:
	NAGG_CONFIG_PATH="./examples/gateway.json" go run cmd/nagg/main.go

.PHONY: clean
clean:
	rm -rf build/
	rm -rf vendor/
