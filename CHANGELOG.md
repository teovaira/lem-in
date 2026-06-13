# Changelog

All notable changes to this project are documented here.

Format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

---

## [1.0.0] — 2026-06-13

### Added
- Go module (`lemin`), stdlib-only, no third-party dependencies
- Package structure: `internal/graph`, `internal/parser`, `internal/pathfinder`, `internal/simulator`, `internal/output`, `cmd/lem-in`
- Makefile with `build`, `test`, `vet`, `race`, `clean`, `run`, `fmt` targets
- GitHub Actions CI: vet → test (race) → build on every push and PR
- `graph`: `Colony`, `Room`, `Path`, `Move` types; `Neighbours` and `HasLink` helpers
- `parser`: state-machine parser with full validation — ant count, room names, coordinates, `##start`/`##end`, duplicate rooms, duplicate/self links
- `pathfinder`: Edmonds-Karp max-flow with node-splitting for vertex-disjoint paths; greedy optimal subset selection via `computeTurns`
- `simulator`: greedy ant-to-path assignment minimising total turns; turn-by-turn move generation
- `output`: verbatim colony echo followed by turn lines via `io.Writer`
- `main`: wires all packages, writes errors to stderr with exit code 1
- Testdata: `example00`–`example07` and `badexample00`–`badexample01` covering all auditor cases

### Fixed
- `.gitignore`: anchored binary pattern to repo root (`/lem-in`) so `cmd/lem-in/` source is never ignored
- `pathfinder`: separated max-flow phase from path-extraction phase to guarantee vertex-disjoint paths
