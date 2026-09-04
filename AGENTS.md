# AGENTS.md

NeON is a build tool (YAML-driven, `build.yml`) written in Go. The repo dogfoods it: the project's own build file is `build.yml`, driven by an installed `neon` binary.

## Commands

- Daily work: `neon <target>` from the repo root (`neon -targets` lists all). Avoid bare `neon` — the default target is `clean, fmt, lint, test, bugs`, and `clean` wipes `build/` plus the Go build/test cache.
- `neon test` — runs `go test -cover ./...`. This also executes the YAML integration tests (see Testing), so it is the main verification step.
- `neon build` — builds `build/neon` for the current platform, injecting the version via `-ldflags -X github.com/c4s4/neon/neon/build.NeonVersion=$(git rev-parse --short HEAD)`. Without that flag `neon -version` prints `UNKNOWN`.
- `neon bugs` — regression suite (see Testing); depends on `build` and invokes the `neon` binary from PATH.
- `neon fmt` / `neon lint` — `gofmt -s -w .` and `golangci-lint run` (lint needs the tool; `neon tools` installs gox, golangci-lint, gobinsec).
- Without neon installed, fall back to plain `go test ./...`, `gofmt -s`, `golangci-lint run`. CI (`.travis.yml`) runs `go test -race -covermode=atomic ./...` on Go 1.26.
- Single package / single test: `go test -count=1 ./neon/builtin/`, `go test -count=1 -run TestIntegration ./neon/task/`.
- Trust `build.yml` over the "Build" section of `README.md` (the README's ldflags example has a stale variable path).

## Prerequisites / gotchas

- `build.yml` extends `c4s4/build/golang.yml`, resolved from the neon plugin repository `~/.neon`. If neon fails to load the build file, install it with `neon -install c4s4/build`.
- `configuration: '~/.neon/github.yml'` loads the maintainer's out-of-repo config (release token). `release`/`deploy` targets are maintainer-only; never commit secrets.
- `build/` is gitignored local output (binaries, `build/tst` test scratch, PDFs) — never commit it.
- Go modules; `go.mod` requires Go >= 1.26.5.

## Architecture

- Single module `github.com/c4s4/neon`. Entrypoint: `neon/main.go`. Core packages: `neon/build` (engine), `neon/task` (build steps), `neon/builtin` (expression functions), `neon/util`.
- Tasks and builtins self-register: each file has an `init()` calling `build.AddTask`/`build.AddBuiltin`; `main.go` blank-imports both packages. Adding a task/builtin is just adding a file to the package — no wiring needed.
- The `Help` text in each `TaskDesc`/`BuiltinDesc` is the source of truth for the generated reference docs.

## Generated docs

- `doc/tasks.md` and `doc/builtins.md` are generated from the built binary (`neon -tasks-ref` / `neon -builtins-ref`, via the `refs` target). Regenerate them whenever task/builtin help text changes.

## Testing

- Unit tests are colocated `*_test.go` files.
- Integration tests are plain YAML build files: `test/builtin/*.yml` and `test/task/*.yml`, executed by `TestIntegration` inside each package when `go test` runs. They write scratch files under `build/tst`.
- `test/bugs/` holds regression build files for past bugs (some with paired `.tst` mini build files); only run by the `bugs` target, i.e. `neon -file test/bugs/<name>.yml` per file.
- When reproducing a bug, add a build file under `test/bugs/` so it is covered by `neon bugs`.

## Conventions

- Branch flow: work lands on `develop`, merged to `master` ("Merge develop" commits).
- Update `CHANGELOG.md` for user-facing changes: newest entry first, format `## YYYY-MM-DD: X.Y.Z` followed by bullet points.
