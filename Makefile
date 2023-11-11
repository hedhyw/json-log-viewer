GOLANG_CI_LINT_VER:=v1.55.2
OUT_BIN?=${PWD}/bin/jlv
COVER_PACKAGES=./...
VERSION?=${shell git describe --tags}

all: lint test build

run:
	@echo "building ${VERSION}"
	go run ./cmd/jlv assets/example.log
.PHONY: build

build:
	@echo "building ${VERSION}"
	go build \
		-o ${OUT_BIN} \
		--ldflags "-s -w -X main.version=${VERSION}" \
		./cmd/jlv
.PHONY: build

install:
	go install ./cmd/jlv
.PHONY: install

lint: bin/golangci-lint
	./bin/golangci-lint run
.PHONY: lint

test:
	go test \
		-coverpkg=${COVER_PACKAGES} \
		-covermode=atomic \
		-race \
		-coverprofile=coverage.out \
		./...
	go tool cover -func=coverage.out
.PHONY: test

vendor:
	go mod tidy
	go mod vendor
.PHONY: vendor

bin/golangci-lint:
	curl \
		-sSfL \
		https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
		| sh -s $(GOLANG_CI_LINT_VER)