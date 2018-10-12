all: build

.PHONY: build
build b:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go generate ./...
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build\
		-o bin/pigeon-http *.go
build-osx bx:
	@CGO_ENABLED=0 go generate ./...
	@CGO_ENABLED=0 GOARCH=amd64 go build\
		-o bin/pigeon-http *.go
deploy d: build
	@docker-compose build
	@docker-compose up
