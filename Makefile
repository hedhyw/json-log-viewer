GOLANG_CI_LINT_VER:=v2.0.2
GORELEASER_VERSION:=v2.3.2
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

fix: bin/golangci-lint-${GOLANG_CI_LINT_VER}
	 gofumpt -l -w .
	./bin/golangci-lint-${GOLANG_CI_LINT_VER} run --fix
.PHONY: lint-fix

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

goreleaser.check:
	docker run --rm -it \
		-v ${PWD}:/go/src/github.com/hedhyw/json-log-viewer \
		-w /go/src/github.com/hedhyw/json-log-viewer \
		goreleaser/goreleaser:${GORELEASER_VERSION} check
.PHONY: goreleaser.check

bin/golangci-lint-${GOLANG_CI_LINT_VER}:
	curl \
		-sSfL \
		https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
		| sh -s $(GOLANG_CI_LINT_VER)
	mv ./bin/golangci-lint ./bin/golangci-lint-${GOLANG_CI_LINT_VER}
