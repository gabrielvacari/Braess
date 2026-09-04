<!--
Sync Impact Report
- Version change: (none) → 1.0.0 (initial ratification)
- Modified principles: n/a (first adoption)
- Added sections:
  - Core Principles: I. Engine Independent of UI, II. Simulate Before You Visualize
    (NON-NEGOTIABLE), III. Decentralized, Selfish Agent Behavior, IV. Validate in
    Text Before Validating Visually, V. Explain Reasoning, Not Just Code
  - Technology Constraints
  - Development Workflow
  - Governance
- Removed sections: none
- Templates requiring updates:
  - .specify/templates/plan-template.md ⚠ pending manual check (verify Constitution
    Check gates reference the roadmap-order and engine-independence principles)
  - .specify/templates/spec-template.md ✅ no changes required (principle-agnostic)
  - .specify/templates/tasks-template.md ✅ no changes required (principle-agnostic)
  - .specify/templates/checklist-template.md ✅ no changes required (principle-agnostic)
- Follow-up TODOs: none — all placeholders resolved from AGENTS.md
-->

# Braess Traffic Simulator Constitution

## Core Principles

### I. Engine Independent of UI
The simulation core (graph + agents + congestion model) MUST run and be fully
testable standalone — e.g. as a terminal script or CLI — with zero dependency
on any UI framework. UI code MAY depend on the engine; the engine MUST NOT
depend on UI code, rendering libraries, or browser-only APIs.
**Rationale**: The project's own architecture splits Graph, Agents, and
Frontend into independent blocks specifically so simulation logic can be
validated before any visuals exist. Coupling the engine to a UI layer would
make correctness (and Braess's Paradox itself) impossible to verify in
isolation.

### II. Simulate Before You Visualize (NON-NEGOTIABLE)
Implementation MUST proceed in this order, and a stage MUST NOT begin until
the previous stage demonstrably works:
1. Pure graph in code (nodes and edges), no UI.
2. Simple agents routed A→B by shortest path, no congestion.
3. Congestion: travel time increasing with volume, with Braess's Paradox
   shown to appear via console/log output.
4. Multiple distinct origins (houses) and destinations (companies).
5. Web frontend visualizing what is already validated underneath.
**Rationale**: This is the project's explicit golden rule — jumping to the
frontend before the graph, agents, and congestion model are proven in
text/logs risks building a visualization of behavior that was never actually
verified.

### III. Decentralized, Selfish Agent Behavior
Each agent (car) MUST choose and recalculate its route to minimize its own
travel time given current conditions. Route choice MUST be computed
per-agent from local/shared road-state information, never dictated by a
central optimizer that assigns system-optimal routes. Emergent effects
(congestion, Braess's Paradox) MUST arise from the aggregate of individual
selfish decisions, not be scripted or hard-coded.
**Rationale**: The project exists to experimentally observe whether
individually rational choices worsen collective outcomes. A centrally
optimized router would make that question unanswerable and defeats the
project's purpose as a decentralized multi-agent systems exercise.

### IV. Validate in Text Before Validating Visually
Every new piece of core-simulation behavior (a routing rule, a congestion
function, a new topology feature) MUST be exercised and checked via
console/log output — showing the effect on travel times and, where relevant,
the presence or absence of Braess's Paradox — before or independently of any
graphical verification.
**Rationale**: Text/log validation is cheap, scriptable, and unambiguous; it
is the project's chosen falsification method for whether "adding a road"
helped or hurt, and must not be skipped in favor of eyeballing an animation.

### V. Explain Reasoning, Not Just Code
When proposing simulation or agent-behavior code, the reasoning behind
design decisions MUST be explained alongside the code — trade-offs
considered, why a particular data structure or algorithm was chosen — rather
than handing over a finished solution with no discussion.
**Rationale**: This is explicitly an AI Engineering learning project; the
author's stated goal is conceptual synthesis through dialogue (agent
fundamentals, graph engineering), not just a working artifact.

## Technology Constraints

- Simulation engine: Go.
- Frontend: React with Konva.js for canvas-based graph drawing/animation.
- These choices may be adjusted as development progresses, but any change
  MUST be a deliberate decision recorded in AGENTS.md and, if it affects a
  Core Principle, reflected here via an amendment.

## Development Workflow

- Commits MUST follow Conventional Commits.
- Branching MUST follow Gitflow (feature branches, no direct unreviewed
  commits to the trunk for non-trivial change).
- Any plan or task list produced for this project MUST be checked against
  Principle II's roadmap order before work starts.

## Governance

This constitution supersedes ad hoc practice for this project. Amendments
are made by editing this file together with, where applicable, AGENTS.md,
and MUST include a Sync Impact Report per the amendment procedure defined in
the Spec Kit constitution workflow.

Versioning follows semantic versioning: MAJOR for backward-incompatible
principle removals/redefinitions, MINOR for new principles or materially
expanded guidance, PATCH for wording/clarification fixes.

All specs, plans, and task lists produced via Spec Kit commands MUST verify
compliance with these principles — in particular Principles I, II, and III
— before implementation begins.

**Version**: 1.0.0 | **Ratified**: 2026-09-03 | **Last Amended**: 2026-09-03
