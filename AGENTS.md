# AGENTS.md

Guide for AI coding agents working on this repository.

## What this project is

`json-log-viewer` (binary name: `jlv`) is an interactive terminal (TUI)
tool for viewing and analyzing structured JSON logs. It shows a compact,
colorized, filterable table of log entries and can expand any entry into a
prettified JSON tree. It is built on
[charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) for
the TUI and [antonmedv/fx](https://github.com/antonmedv/fx) for JSON
viewing. It is an application, not a library: `internal/` packages are not
importable by other projects.

## CLI usage

```sh
jlv assets/example.log            # open a file
jlv < assets/example.log          # read from stdin
kubectl logs pod/NAME -f | jlv    # follow piped logs
jlv -config example.jlv.jsonc app.log
jlv --version
```

Flags: `-config <path>` (path to config) and `-version`. Without `-config`,
the app looks for `.jlv.jsonc` in `$PWD`, then in `$HOME`. The config is
JSONC (JSON with comments); see `example.jlv.jsonc` for all fields
(customizable columns with JSONPath selectors, time formats, etc.).

## Repository layout

- `cmd/jlv/` - entry point: flag parsing, stdin/file input wiring.
- `internal/app/` - TUI application. State-machine pattern: each
  `state*.go` file (`stateinitial`, `stateloaded`, `statefiltering`,
  `statefiltered`, `stateviewrow`, `stateerror`) is a bubbletea model for
  one application state; `logstable.go` / `lazytable.go` render the table.
- `internal/keymap/` - key bindings.
- `internal/pkg/config/` - configuration loading and validation.
- `internal/pkg/source/` - log input: reading, streaming (`steamer.go`),
  parsing entries, log-level detection.
- `internal/pkg/events/` - bubbletea event (message) definitions.
- `internal/pkg/widgets/` - reusable TUI widgets (JSON view, pill input,
  plain-text view).
- `internal/pkg/tests/` - shared test helpers.
- `docs/` - user documentation (features, usage, customization).
- `assets/` - example logs and images used in docs and manual testing.

There is no code generation in this repo; all files are hand-written.

## Build, test, lint

```sh
make build        # builds ./bin/jlv (version injected via ldflags)
make test         # go test -race with coverage; prints per-func coverage
make lint         # golangci-lint (auto-installs pinned version into ./bin)
make all          # lint + test + build
make run          # build and open assets/example.log
make fix          # gofumpt + golangci-lint --fix
make install      # go install ./cmd/jlv
```

- Go 1.25+ is required (see `go` directive in `go.mod`).
- golangci-lint version is pinned by `GOLANG_CI_LINT_VER` in the
  `Makefile`; config is `.golangci.json` (v2 format, `default: all` with an
  explicit disable list). Do not disable additional linters just to
  silence findings in new code.
- The `vendor/` directory is gitignored; do not commit it. Releases run
  `make vendor` themselves.

## Conventions

- Conventional commit messages (`feat:`, `fix:`, `chore:`, `ci:`,
  `docs:`); PR titles are validated by the semantic-pull-request workflow.
- Code style is enforced by gofumpt and golangci-lint; run `make fix`
  before committing.
- Tests live next to the code (`*_test.go`) and use
  `github.com/stretchr/testify`. TUI behavior is tested with
  `charmbracelet/x/exp/teatest`. Keep the high coverage bar (~94%).
- `go.mod` contains `replace` directives pointing
  `github.com/antonmedv/fx` and `github.com/charmbracelet/bubbles` to
  `hedhyw/*` forks. Do not remove these when updating dependencies.
- CI (`.github/workflows/check.yml`) runs `make build`, `make lint`,
  `make test` on every PR; actions are pinned to commit SHAs.
