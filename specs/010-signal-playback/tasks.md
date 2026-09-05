---

description: "Task list template for feature implementation"
---

# Tasks: Signal-Controlled Live Playback

**Input**: Design documents from `/specs/010-signal-playback/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/playback-contract.md, quickstart.md. Depends on features 008 (`queuesim`) and 009 (`cmd/server`'s `/api/queue-run`, `web/`'s Signals mode), both extended additively.

**Tests**: Included — Go tests for the new engine recordings (validated in text before UI, per Principle II), a Go contract test extension, and Vitest for the new pure playback functions, per every prior feature's precedent.

**Organization**: Tasks are grouped by user story (spec.md) to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)

## Path Conventions

`queuesim/` and `cmd/server/` (extended, not new files) and `web/src/`
(new model/animation/component files plus `SignalApp.tsx` and
`NetworkCanvas.tsx` extensions). `graph/`, `agent/`, `simulation/`,
`cmd/graphcli/`, and every equilibrium-mode file are not touched.

---

## Phase 1: Foundational (Blocking Prerequisites)

**Purpose**: The engine recordings and their wire types every user story needs to render anything

**⚠️ CRITICAL**: No user story work can begin until this phase is complete — Principle II requires the recordings validated in Go tests before any frontend code reads them

- [X] T001 [P] Add `PositionSample` type, extend `QueueSample` with `Green bool`, extend `AgentReport` with `ID int`, extend `RunResult` with `Positions []PositionSample` in `queuesim/report.go` (data-model.md)
- [X] T002 Record per-tick position samples and signal phase in `queuesim/run.go`: assign each spawned agent a stable `ID`, record one `PositionSample` per in-flight agent per tick (progress while moving, 0 while queued) plus a final `Progress: 1` sample on the tick an agent arrives, and set `QueueSample.Green` from `Signal.IsGreenAt(t)` (research.md decisions #1-#3; depends on T001)
- [X] T003 [P] Go tests in `queuesim/run_test.go` and `queuesim/two_road_test.go` asserting: a position sample exists for every in-flight agent every tick; `Queued` matches the agent's actual state; `Progress` stays in `[0, 1]`; an arriving agent's last sample has `Progress: 1`; `QueueSample.Green` matches `Signal.IsGreenAt(Time)` (depends on T002)
- [X] T004 [P] Extend `cmd/server/queuedto.go`: `PositionSampleDTO`, `QueueSampleDTO.Green`, `AgentReportDTO.ID`, `QueueRunResponse.Positions`, and update `toQueueRunResponse` to map them (contracts/playback-contract.md; depends on T001)
- [X] T005 [P] Extend `cmd/server/queuehandlers_test.go`: a valid `/api/queue-run` response includes non-empty `positions` and at least one `queueSamples[].green == true` and one `== false` (depends on T004)
- [X] T006 [P] Mirror the new fields/types in `web/src/model/queueApi.ts`: `PositionSampleDTO`, `QueueSampleDTO.green`, `AgentReportDTO.id`, `QueueRunResponse.positions` (contracts/playback-contract.md; depends on T004)

**Checkpoint**: The engine records and serves everything a playback needs — validated by Go tests before any UI touches it

---

## Phase 2: User Story 1 - Watch agents move through a signal-controlled network (Priority: P1) 🎯 MVP

**Goal**: Agents visibly move from origin to destination, pausing while queued and resuming when their signal turns green.

**Independent Test**: Run the classic two-road, opposite-signal scenario and confirm agents visibly queue on one road while moving on the other, then swap roles when the signals flip.

### Implementation for User Story 1

- [X] T007 [US1] Implement `groupPositionsByAgent` and `positionAtSimTime` in `web/src/model/queuePlayback.ts` (data-model.md; depends on T006)
- [X] T008 [P] [US1] Implement `usePlaybackClock` in `web/src/animation/usePlaybackClock.ts`: a single simulated-time clock advanced via `requestAnimationFrame`, clamped to `[0, duration]`, with play/pause/restart (research.md decision #4)
- [X] T009 [US1] Wire agent playback into `SignalApp.tsx`: build node/edge lookup maps from `network`, group the run's `positions` by agent, compute one `AgentMarker` per agent each frame via `positionAtSimTime`, pass into the existing `NetworkCanvas` `agents` prop, and add play/restart controls driven by `usePlaybackClock` (depends on T007, T008)

### Tests for User Story 1

- [X] T010 [P] [US1] Vitest tests for `web/src/model/queuePlayback.test.ts`: `positionAtSimTime` interpolates along an edge while moving, holds at the edge's start while queued, and returns `null` before an agent's first recorded sample; `groupPositionsByAgent` groups and preserves time order (depends on T007)

**Checkpoint**: Agents visibly move and queue/resume in the browser — User Story 1 is independently functional and demoable.

---

## Phase 3: User Story 2 - See each signal's live state as the animation plays (Priority: P1)

**Goal**: Each signal-controlled road visually shows green or red, updating as playback advances.

**Independent Test**: Run a scenario with one signal-controlled road, play the animation, and confirm the road's visual state changes from green to red and back on schedule, matching the signal's configured durations.

### Implementation for User Story 2

- [X] T011 [US2] Implement `signalGreenAtSimTime` in `web/src/model/queuePlayback.ts` (depends on T007, same file)
- [X] T012 [P] [US2] Add an optional `roadColors?: Record<string, string>` prop to `web/src/components/NetworkCanvas.tsx`: when a road's id is a key, its line uses that color instead of the default, unless it's the current selection (research.md decision #5)
- [X] T013 [US2] Wire road coloring into `SignalApp.tsx`: for every directed road with a signal, compute its color each frame via `signalGreenAtSimTime` and pass the result as `roadColors` into `NetworkCanvas` (depends on T011, T012, T009)

### Tests for User Story 2

- [X] T014 [P] [US2] Vitest test for `signalGreenAtSimTime` in `queuePlayback.test.ts`: returns the nearest recorded sample at or before the queried time, and falls back to the earliest sample when queried before it (depends on T011)

**Checkpoint**: Roads visibly flip green/red in sync with agents stopping and starting — User Stories 1 and 2 both independently pass.

---

## Phase 4: User Story 3 - Keep the existing summary alongside the animation (Priority: P3)

**Goal**: The queue-length chart and arrival/wait summary from feature 009 remain visible together with the new animation.

**Independent Test**: Run a scenario and confirm both the live animation and the existing chart/summary are visible together, not one in place of the other.

### Implementation for User Story 3

- [X] T015 [US3] Confirm `SignalApp.tsx`'s layout renders the playback (canvas + controls) together with the existing `QueueResultsPanel`, without removing or hiding either (depends on T013)

### Tests for User Story 3

- [X] T016 [P] [US3] Confirm via `git diff` that `components/QueueResultsPanel.tsx`, `components/QueueChart.tsx`, and `model/queueChart.ts` are unchanged by this feature (FR-008)

**Checkpoint**: All three user stories are independently functional.

---

## Phase 5: Polish & Cross-Cutting Concerns

- [X] T017 [P] Run `gofmt -l .` and `go vet ./...`; `npx tsc --noEmit` and `npx oxlint` for `web/`; fix any findings
- [X] T018 Run the full quickstart.md validation end to end: `go test ./...`, `npm test` (in `web/`), then a live curl + manual browser walkthrough (depends on T003, T005, T010, T014, T016)
- [X] T019 [P] Confirm via `git diff --stat` that `graph/`, `agent/`, `simulation/`, `cmd/graphcli/`, `animation/useAgentAnimation.ts`, `animation/AgentTicker.tsx`, and the `EquilibriumApp` portion of `App.tsx` are unchanged (constitution Principle I; FR-009)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Foundational (Phase 1)**: BLOCKS all user stories
- **User Story 1 (Phase 2)**: Depends on Foundational only
- **User Story 2 (Phase 3)**: Depends on Foundational directly, and reuses User Story 1's `SignalApp` per-frame wiring (T009) to attach road coloring to
- **User Story 3 (Phase 4)**: Depends on User Story 2's wiring being in place, since it's a layout/preservation check on the same component
- **Polish (Phase 5)**: Depends on all three user stories

### Parallel Opportunities

- `T001`, `T004` (Foundational, independent files) in parallel; `T003`, `T005`, `T006` each depend on one of those but not each other
- `T008` (playback clock) and `T007` (playback math) in parallel — different files
- `T012` (NetworkCanvas prop) can start as soon as Foundational is done, in parallel with `T007`-`T010`
- `T017` and `T019` (Polish) in parallel

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Foundational
2. Complete Phase 2: User Story 1
3. **STOP and VALIDATE**: agents visibly move and queue in the browser
4. This proves the recorded data and the playback math are correct before signal coloring is layered on

### Incremental Delivery

1. Foundational → engine recordings validated in Go tests, served over the API
2. User Story 1 → validate independently (agents move and queue)
3. User Story 2 → validate independently (roads show live signal state)
4. User Story 3 → validate independently (nothing prior was lost)
5. Polish → lint/build check + confirm zero changes to the engine's other packages or the equilibrium mode

---

## Notes

- [P] tasks touch different files with no dependency between them
- [Story] labels trace every implementation task back to spec.md
- Commit after each task or logical group, using Conventional Commits (constitution Development Workflow)
- No scrubber/seek control, and no change to how many signals a road may
  have, by design (spec.md Assumptions) — do not add them here
