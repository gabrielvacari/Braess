# Phase 0 Research: Graph Simulation Engine

No `[NEEDS CLARIFICATION]` markers remained in the Technical Context, so
this covers the technology/design decisions made and why — per constitution
Principle V (explain reasoning, not just code).

## 1. Module layout: library package + separate `cmd/` CLI

**Decision**: Put the engine in a root-level Go package (`graph`) and the
terminal entry point required by FR-006 in `cmd/graphcli`, which imports
`graph` — never the reverse.

**Rationale**: This is the idiomatic Go project layout, and it makes
Principle I (Engine Independent of UI) a structural fact rather than a
promise: `graph` cannot import a UI package that doesn't exist yet, and any
future frontend/agent code will depend on `graph`, not the other way
around. It also keeps `go test ./graph/...` runnable with zero CLI
concerns.

**Alternatives considered**: A generic `src/`/`lib/` tree (as suggested by
the plan template's default Option 1) was rejected — it isn't how Go code
is organized in practice, and following it here would fight the language's
own tooling (`go test`, `go build`, import paths) for no benefit.

## 2. Travel-time function representation

**Decision**: Represent each edge's travel-time behavior as a function
value, `TravelTimeFunc func(volume float64) (float64, error)`, stored on
the edge, rather than a single global formula computed by the graph. Ship
two small constructors as the initial, swappable defaults: `Linear(freeFlow,
slope float64)` (time = freeFlow + slope·volume) and `Constant(t float64)`
(time is fixed regardless of volume).

**Rationale**: The spec's Assumptions section deliberately leaves the exact
formula open — "any function that is non-decreasing in volume and returns
the free-flow time at zero volume satisfies it." That's not laziness: the
classic Braess's Paradox network (the one this whole project exists to
reproduce) needs *different* function shapes on different edges in the same
graph — some roads get slower as more cars use them, others stay constant
regardless of volume. A single hardcoded formula for every edge would make
that network impossible to build later. `Linear` and `Constant` are enough
to construct that classic example once agents exist, without committing the
engine to one "true" congestion model.

**Alternatives considered**: A full BPR curve (`t0·(1 + α·(v/capacity)^β)`)
as the only option was considered — it's the traffic-engineering standard —
but rejected *as the sole option* because it can't express the constant-time
edges the classic paradox example requires. It remains a reasonable
additional constructor to add later if a specific study needs it; nothing
here prevents that.

## 3. Testing approach

**Decision**: Go's standard `testing` package, table-driven tests per file,
run via `go test ./...`. No third-party assertion or mocking library.

**Rationale**: The engine has no external dependencies to mock and no
behavior complex enough to need a matcher DSL — `if got != want` is
sufficient and keeps `graph` dependency-free, reinforcing decision #1.

**Alternatives considered**: `testify` — common in the Go ecosystem, but an
unnecessary dependency for this scope; would only be worth adding if
assertions become unwieldy, which is not the case for a handful of node/edge
invariants.

## 4. CLI implementation

**Decision**: `cmd/graphcli/main.go` uses only `fmt` to build a small sample
network in code and print it — no `flag`-based options are needed yet since
User Story 3 only requires proving the engine is inspectable in text, not a
configurable tool.

**Rationale**: Matches Principle IV (validate in text before visually) with
the smallest possible surface. A configurable CLI (choose graph from a
file, etc.) is out of scope for this spec — it's not required by any FR and
would be speculative for a phase that exists purely to prove the engine
works standalone.

**Alternatives considered**: A CLI framework (e.g. `cobra`) — rejected as
unjustified complexity for a single, argument-free command.
