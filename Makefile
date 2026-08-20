GOLANG_CI_LINT_VER:=v2.13.0
# Pinned by digest: a tag alone is mutable. Docker rejects the reference
# if the tag and the digest disagree, so the two cannot drift apart.
GORELEASER_VERSION:=v2.3.2@sha256:d62b4a18dfe3af7bd4da9e5954b496548ef04e73ae8f98cd75ba63a9ed4d73e5
GOVERSIONINFO_VERSION:=v1.7.0
OUT_BIN?=${PWD}/bin/jlv
COVER_PACKAGES=./...
VERSION?=${shell git describe --tags}
# Windows VERSIONINFO resources only accept numeric fields.
VERSION_NUM=${shell echo "${VERSION}" | sed -E 's/^v//; s/-.*$$//'}

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
		--ldflags "-X main.version=${VERSION}" \
		./cmd/jlv
.PHONY: build

# Generates Windows VERSIONINFO resources (resource_windows_*.syso) that are
# picked up automatically by the Go linker for windows targets.
versioninfo:
	cd cmd/jlv && GOOS= GOARCH= GOARM= GOFLAGS=-mod=mod go run \
		github.com/josephspurrier/goversioninfo/cmd/goversioninfo@${GOVERSIONINFO_VERSION} \
		-platform-specific \
		-propagate-ver-strings \
		-file-version "${VERSION_NUM}" \
		-product-version "${VERSION_NUM}" \
		versioninfo.json
.PHONY: versioninfo

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
	GOBIN=$(PWD)/bin go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANG_CI_LINT_VER)
	mv ./bin/golangci-lint ./bin/golangci-lint-${GOLANG_CI_LINT_VER}
