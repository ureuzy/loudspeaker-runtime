IMG ?= ureuzy/loudspeaker-runtime:latest

fmt:
	go fmt ./...

vet:
	go vet ./...

build: fmt vet
	go build -o bin/runtime ./cmd/runtime/main.go

run: fmt vet
	go run ./cmd/runtime/main.go

test:
	go test ./... -coverprofile cover.out

docker-build: test
	docker build -t ${IMG} .

docker-push:
	docker push ${IMG}
