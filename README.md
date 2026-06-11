# lem-in

A Go program that simulates an ant colony. Given a colony file describing rooms, tunnels, and a number of ants, it finds the optimal set of non-overlapping paths from `##start` to `##end` and moves all ants across in the minimum number of turns.

Built at Zone01 — Go · TDD · Conventional Commits.

---

## Team

| Member | Packages |
|---|---|
| KRYSTALLENIA | `graph` · `parser` |
| VASILIKI | `pathfinder` |
| THEO | `simulator` · `output` · `main` |

---

## Build

```bash
git clone <repo-url>
cd lem-in
make build
```

Requires Go 1.22+. No third-party dependencies.

---

## Run

```bash
./lem-in testdata/example00.txt
```

Or without building:

```bash
go run ./cmd/lem-in testdata/example00.txt
```

---

## Output format

The program prints the original colony file verbatim, one blank line, then one line per turn:

```
4
##start
0 0 3
2 2 5
3 4 0
##end
1 8 3
0-2
2-3
3-1

L1-2
L1-3 L2-2
L1-1 L2-3 L3-2
L2-1 L3-3 L4-2
L3-1 L4-3
L4-1
```

Each move is formatted as `L{ant_id}-{room_name}`. Only ants that moved are printed per turn.

---

## Error handling

All errors are written to `stderr` with exit code 1:

```
ERROR: invalid data format, no start room found
```

---

## Test

```bash
make test        # run all tests
make race        # run with race detector
make vet         # run go vet
```

---

## Project reference

See [lem-in-reference.md](lem-in-reference.md) for the full implementation plan, algorithm details, and team workflow.
