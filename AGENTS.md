# Repository Guidelines

## Project Structure & Module Organization
Core parsing logic lives in `config.go`, with generated schema helpers under `schema/`. Code generation tooling resides in `cmd/gen_values`, reusable configuration fixtures sit in `examples/`, and parsed sample inputs belong in `testdata/`. Keep new assets beside the code they support to simplify `go generate` output diffs and targeted tests.

## Build, Test & Generation Commands
- `go build ./...` – verifies the module compiles with Go 1.23.2 across every package.
- `go test ./...` – runs table-driven suites in `config_test.go` and generator checks under `cmd/gen_values`.
- `go test ./... -run Read -v` – focus on parser regressions when iterating quickly.
- `go generate ./...` – refreshes `schema/args/args_generated.go` and `schema/values/values_generated.go`; requires `goimports` on PATH.
- `go run cmd/gen_values/main.go -pkg values` – regenerate a single package when experimenting with value handlers.

## Coding Style & Naming Conventions
Follow standard Go formatting: tabs for indentation, max 100-character logical lines, and `gofmt` + `goimports` before commit. Exported APIs should use PascalCase, internal helpers camelCase, and tests mirror the function under test with `TestXxx` names. Prefer short, purposeful package-level comments that state intent rather than implementation trivia.

## Testing Guidelines
All new behavior needs table-driven coverage using the `testing` package and `t.Run` descriptions matching feature intent. Place configuration fixtures in `testdata/` with descriptive `.conf` filenames and assert node counts or contents as shown in existing suites. Run `go test ./... -cover` before submitting and include generator-focused assertions if schema helpers change.

## Commit & Pull Request Guidelines
Use concise, imperative subject lines (e.g., `Add parser guard for empty blocks`) capped near 72 characters, mirroring the current Git history. Squash noisy commits locally, reference related issues in the body, and note follow-up work when applicable. Pull requests should describe the scenario exercised, list commands executed (`go test`, `go generate`), and attach snippets or diffs for complex schema updates so reviewers can trace configuration impacts quickly.

## Configuration & Schema Tips
Treat generated files as read-only; adjust templates in `cmd/gen_values` or schema builders, then rerun generation. Sync example configurations in `examples/` with parser features to keep documentation accurate. When introducing new directives, document expected arguments in code comments and add a matching snippet or macro example to help downstream integrators.
