# Agents

This file tells AI coding assistants (Claude Code, Copilot, Cursor, etc.) how to work in this repository.

---

## Read this first

Before touching any file, read [lem-in-reference.md](lem-in-reference.md). It is the non-negotiable specification for every type, function signature, algorithm, and error message in the project. Do not deviate from it without updating the document first.

---

## Package ownership

| Package | Owner | File |
|---|---|---|
| `graph` | KRYSTALLENIA | `internal/graph/graph.go`, `internal/graph/graph_test.go` |
| `parser` | KRYSTALLENIA | `internal/parser/parser.go`, `internal/parser/parser_test.go` |
| `pathfinder` | VASILIKI | `internal/pathfinder/pathfinder.go`, `internal/pathfinder/pathfinder_test.go` |
| `simulator` | THEO | `internal/simulator/simulator.go`, `internal/simulator/simulator_test.go` |
| `output` | THEO | `internal/output/output.go`, `internal/output/output_test.go` |
| `main` | THEO | `cmd/lem-in/main.go`, `cmd/lem-in/main_test.go` |

Do not generate code in another owner's package unless explicitly asked.

---

## Hard rules

- **Standard library only.** Never suggest or add a third-party import. No `go get`. If you need something, use `fmt`, `os`, `bufio`, `strings`, `strconv`, `sort`, `errors`, or other stdlib packages.
- **TDD.** Tests come before production code. If asked to implement a function, write the test first, confirm it fails, then write the implementation.
- **Conventional commits.** Every commit message must follow `type(scope): description`. See `CONTRIBUTING.md`.
- **No global mutable state.** No package-level variables that change at runtime.
- **No `panic`.** Return errors. Let `main.go` decide what to print and when to exit.
- **No `os.Exit` outside `main.go`.**
- **`gofmt` before committing.** Generated code must be correctly formatted.

---

## Function signatures are frozen

The signatures in `lem-in-reference.md` Section 3.6 are locked after Phase 1 merges:

```go
func Parse(filename string) (*graph.Colony, error)
func FindPaths(c *graph.Colony) ([]graph.Path, error)
func Simulate(c *graph.Colony, paths []graph.Path) [][]graph.Move
func Print(w io.Writer, c *graph.Colony, turns [][]graph.Move)
```

Do not change these without a team discussion and a document update.

---

## Error format contract

Packages `parser` and `pathfinder` return **plain suffix strings** as errors:

```go
return errors.New("invalid room name")      // correct
return fmt.Errorf("context: %w", err)       // wrong — pollutes the output
```

`main.go` wraps them: `fmt.Fprintf(os.Stderr, "ERROR: invalid data format, %s\n", err)`

---

## Test conventions

- Table-driven tests for multiple similar cases.
- Use `os.CreateTemp` for parser tests — write content to a temp file, pass the path to `Parse`.
- Clean up with `defer os.Remove(f.Name())`.
- Use `t.Fatalf` when continuing after failure makes no sense.
- Never use `time.Sleep` in tests.
- Performance tests use `t.Parallel()` and a generous timeout via `go test -timeout`.

---

## What success looks like

```bash
go build ./cmd/lem-in   # zero errors
go test ./...           # zero failures
go test -race ./...     # zero races
go vet ./...            # zero issues
gofmt -d .              # zero output
./lem-in testdata/example00.txt  # 4 ants, ≤6 turns
```
