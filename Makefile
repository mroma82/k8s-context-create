all:
	make build

run:
	go run cmd/k8s-context-create/*.go --insecure --token $(TOKEN) --host $(HOST)

build-amd64:
	GOOS=linux GOARCH=amd64 go build -o bin/k8s-context-create-amd64 cmd/k8s-context-create/*.go

build-amd64-win:
	GOOS=windows GOARCH=amd64 go build -o bin/k8s-context-create-amd64-win.exe cmd/k8s-context-create/*.go

build-macos:
	GOOS=darwin GOARCH=amd64 go build -o bin/k8s-context-create-amd64-darwin cmd/k8s-context-create/*.go

build:
	make build-amd64
	make build-amd64-win
	make build-macos
