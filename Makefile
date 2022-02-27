run:
	go run cmd/k8s-context-create/*.go --insecure --token $(TOKEN) --host $(HOST)


build-amd64:
	GOOS=windows GOARCH=amd64 go build -o bin/k8s-context-create-amd64-win *.go

build-macos:
	GOOS=darwin GOARCH=amd64 go build -o bin/k8s-context-create-amd64-darwin *.go

build-all:
	make build-amd64
	make build-macos