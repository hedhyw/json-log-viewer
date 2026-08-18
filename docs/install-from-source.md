# Install from source

## Requirements
- `git`;
- `make`;
- `go` 1.25+.

## Go install

```shell
go install github.com/hedhyw/json-log-viewer/cmd/jlv@latest
```

The binary is placed in `$(go env GOPATH)/bin`.

## Build manually

Get the source code:
```shell
git clone git@github.com:hedhyw/json-log-viewer.git && cd json-log-viewer
```

Compile:
```shell
make build
```

Install:
```shell
cp ./bin/jlv /usr/local/bin
chmod +x /usr/local/bin/jlv
```

Use as [jlv](usage.md).
