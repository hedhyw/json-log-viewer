# JSON Log Viewer

![Version](https://img.shields.io/github/v/tag/hedhyw/json-log-viewer)
![Build Status](https://github.com/hedhyw/json-log-viewer/actions/workflows/check.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/hedhyw/json-log-viewer)](https://goreportcard.com/report/github.com/hedhyw/json-log-viewer)
[![Coverage Status](https://coveralls.io/repos/github/hedhyw/json-log-viewer/badge.svg?branch=main)](https://coveralls.io/github/hedhyw/json-log-viewer?branch=main)

It is an interactive tool for viewing and analyzing complex log files with structured JSON logs.

![Animation](./assets/animation.webp)

Main features:
1. It is interactive.
2. Is shows similified log records.
3. It is possible to see the full prettified JSON after clicking.
4. It includes non-JSON logs as they are.
5. It understands different field names.
6. It supports case-insensitive filtering.
7. It is simple.

It uses [antonmedv/fx](https://github.com/antonmedv/fx) for viewing JSON and [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) for terminal UI. The tool is inspired by the project [json-log-viewer](https://github.com/gistia/json-log-viewer) which is unfortunately outdated.

## Table of content

- [JSON Log Viewer](#json-log-viewer)
    - [Table of content](#table-of-content)
    - [Usage](#usage)
    - [Install](#install)
        - [MacOS/Linux HomeBrew](#macoslinux-homebrew)
        - [Go](#go)
        - [Package](#package)
        - [Standalone Binary](#standalone-binary)
        - [Source](#source)
    - [Roadmap](#roadmap)
    - [Resources](#resources)
    - [License](#license)


## Usage

```sh
jlv file.json
```

| Key    | Action         |
| ------ | -------------- |
| Enter  | Open/Close log |
| F      | Filter         |
| Ctrl+C | Exit           |
| Esc    | Back           |
| ↑↓     | Navigation     |

## Install

### MacOS/Linux HomeBrew

```sh
brew install hedhyw/main/json-log-viewer
```

### Go

```bash
go install github.com/hedhyw/json-log-viewer/cmd/jlv@latest
```

### Package

Latest DEB and RPM packages are available on [the releases page](https://github.com/hedhyw/json-log-viewer/releases/latest).

### Standalone Binary

Download latest archive `*.tar.gz` for your target platform from [the releases page](https://github.com/hedhyw/json-log-viewer/releases/latest) and extract it to `/usr/local/bin/jlv`. Add this path to `PATH` environment.

### Source

```
git clone git@github.com:hedhyw/json-log-viewer.git
cd json-log-viewer
make build
cp ./bin/jlv /usr/local/bin
chmod +x /usr/local/bin/jlv
```

## Roadmap

- Accept stream of logs.
- Add colors to log levels.
- Add a configuration file (similar to `.json-log-viewer`).
- Convert number timestamps.

## Resources

Alternatives:
- [mightyguava/jl](https://github.com/mightyguava/jl) - Pretty Viewer for JSON logs.
- [pamburus/hl](https://github.com/pamburus/hl) - A log viewer that translates JSON logs into human-readable representation.
- [json-log-viewer](https://github.com/gistia/json-log-viewer) - Powerful terminal based viewer for JSON logs using ncurses.

## License

[MIT License](LICENSE).
