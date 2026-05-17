# Contributing

Read [lem-in-reference.md](lem-in-reference.md) in full before writing a single line of code.

---

## Who does what

| Member | Packages |
|---|---|
| KRYSTALLENIA | `graph` · `parser` |
| VASILIKI | `pathfinder` |
| THEO | `simulator` · `output` · `main` |

Do not touch another person's package without asking first.

---

## Workflow

1. Pull the latest `main` before starting any work.
2. Create a branch: `git checkout -b <scope>/<short-description>`
   Examples: `parser/state-machine`, `pathfinder/bfs`, `fix/duplicate-link-error`
3. Write the test first. It must fail. Commit it.
4. Write the minimum code to make it pass. Commit it.
5. Open a PR. At least one teammate must review before merge.
6. Never merge your own PR without a review.

---

## Commit format

Every commit follows [Conventional Commits](https://www.conventionalcommits.org/):

```
type(scope): short description in lowercase
```

**Types:** `feat` · `fix` · `test` · `refactor` · `docs` · `chore` · `style`

**Scopes:** `graph` · `parser` · `pathfinder` · `simulator` · `output` · `main` · `repo` · `all`

**Rules:**
- Max 72 characters in the subject line
- Lowercase, no period at the end
- Present tense: `add tests` not `added tests`
- Scope is required — no `feat: do something`

**Examples:**

```
feat(graph): define Colony, Room, Path, Move types
test(parser): add all parser tests — all red
fix(pathfinder): handle self-link without infinite loop
chore(repo): add GitHub Actions CI workflow
```

---

## TDD — non-negotiable

1. Write the test. It must **fail**. Commit: `test(scope): add failing tests for X`
2. Write the minimum production code to pass. Commit: `feat(scope): implement X`
3. Refactor if needed. All tests still pass. Commit: `refactor(scope): clean up X`

Never commit production code before the test that covers it exists.

---

## Before every commit

```bash
go vet ./...       # must be clean
gofmt -d .         # must produce no output
go test ./...      # must pass
```

---

## Before opening a PR

- All tests pass: `make test`
- Race detector clean: `make race`
- No vet issues: `make vet`
- No commented-out code
- Every exported symbol has a godoc comment

---

## Allowed packages

**Only the Go standard library.** No `go get`. No third-party imports. The auditor checks this — a single external import fails the project immediately.
