# Change this and commit to create new release
VERSION=0.1.0
REVISION := $(shell git rev-parse --short HEAD;)

.PHONY: bootstrap
bootstrap: export GO111MODULE=on
bootstrap:
	go mod download && go mod vendor

.PHONY: build
build: export CGO_ENABLED=0
build: export GO111MODULE=on
build:
	go build -v --ldflags="-w -X main.Version=$(VERSION) -X main.Revision=$(REVISION)" \
		-o bin/regdoc cmd/regdoc/main.go

.PHONY: clean
clean:
	git status --ignored --short | grep '^!! ' | sed 's/!! //' | xargs rm -rf
