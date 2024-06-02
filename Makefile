GOLANG_CI_LINT_VER:=v1.59.0
OUT_BIN?=${PWD}/bin/jlv
COVER_PACKAGES=./...
VERSION?=${shell git describe --tags}

all: lint test build

run: build
	./bin/jlv assets/example.log
.PHONY: run

run.version: build
	./bin/jlv --version
.PHONY: run.version

run.stdin: build
	./bin/jlv < assets/example.log
.PHONY: run.stdin

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

lint: bin/golangci-lint-${GOLANG_CI_LINT_VER}
	./bin/golangci-lint-${GOLANG_CI_LINT_VER} run
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

bin/golangci-lint-${GOLANG_CI_LINT_VER}:
	curl \
		-sSfL \
		https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
		| sh -s $(GOLANG_CI_LINT_VER)
	mv ./bin/golangci-lint ./bin/golangci-lint-${GOLANG_CI_LINT_VER}
