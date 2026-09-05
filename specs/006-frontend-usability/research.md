# Phase 0 Research: Frontend Usability and Visual Design

No `[NEEDS CLARIFICATION]` markers remained in the Technical Context, so
this covers the technology/design decisions made and why — per
constitution Principle V.

## 1. Plain CSS custom properties, not a new dependency

**Decision**: Build the design system (color, spacing, typography) as CSS
custom properties in `index.css`, applied via ordinary class names — no
Tailwind, no CSS-in-JS library, no component library.

**Rationale**: Every prior feature's research.md has applied "don't add a
dependency the problem doesn't need" — feature 002 rejected `testify`,
feature 003 rejected a graph library, feature 005 rejected Redux/Zustand
for the same reason. A consistent, deliberate visual system doesn't
require new tooling; it requires actually defining tokens once and using
them everywhere, which plain CSS does perfectly well at this project's
scale (a handful of components, one page).

**Alternatives considered**: Tailwind — rejected; it's a fine choice in
general, but adding a utility-class build pipeline for one page's worth of
components is more machinery than the problem needs, and this project has
consistently preferred the smallest tool that does the job. A component
library (e.g. MUI, Chakra) — rejected for the same reason, and because it
would fight `react-konva`'s canvas-based node/edge rendering, which is
already bespoke and can't use a component library's controls anyway.

## 2. Copy and its selection logic live in one pure, tested module

**Decision**: Every explanatory string (the intro text, the legend
entries, each mode's description, the empty-canvas and missing-demand
hints) lives in `web/src/content/copy.ts`, including three pure functions
— `modeDescription`, `emptyCanvasHint`, `missingDemandHint` — that decide
*which* message applies given the current state. Components import from
here; no explanatory string is written inline in JSX.

**Rationale**: This is a "content" feature, but FR-003/FR-004/FR-009 each
describe real decision logic (which text to show, given the mode or the
network's current size) — logic that can be wrong independently of how
it's styled. Extracting it as plain functions makes it Vitest-testable
exactly like `web/src/model/network.ts`'s edit operations were in feature
005, rather than leaving it as inline JSX conditionals that could only be
checked by eye. This is the same "keep pure logic separate from
rendering" split feature 005 already established, applied to text
instead of network data.

**Alternatives considered**: Inline conditional JSX in each component
(e.g. `{mode === "place" && <p>...</p>}` scattered across `Toolbar.tsx`)
— rejected; it's exactly the pattern that made the original implementation
hard to review for completeness (is every mode covered? every network
size?), and it can't be unit-tested independently of rendering.

## 3. Always-visible static content, not a dismissible tour

**Decision**: The intro banner and legend are always rendered, with no
dismiss button, "don't show again" state, or step-by-step tour.

**Rationale**: The spec's own Assumptions rule out a "dynamic help
system, guided tour, or tooltip framework" — and feature 005 already
established that this frontend has no cross-session persistence.
A dismissible banner would need to persist "already dismissed" somewhere
(session state at minimum) for the dismissal to feel intentional rather
than reappearing on every reload, which is exactly the kind of state this
project has deliberately avoided adding to the frontend. Always-visible,
static content sidesteps that entirely while still fully satisfying
FR-001's "visible as soon as the app loads."

**Alternatives considered**: A dismissible banner using `localStorage` to
remember it was closed — rejected; it's a reasonable pattern in general,
but introduces the project's first piece of browser-persisted state for a
problem ("the banner takes up space after the first visit") that no
success criterion or user story actually asks to solve.

## 4. Redundant text/icon everywhere color currently carries meaning

**Decision**: Every node already renders a short text label inside its
shape (`H`/`C`, feature 005) — extended to give the intersection marker a
label too (`I`) — and the before/after comparison already states "got
worse by X" / "improved by X" in words, not just red/green. This feature
keeps and makes deliberate what was partly already true by accident, and
closes the one gap (intersection's blank label).

**Rationale**: FR-007/SC-005 require that color never be the *only*
carrier of meaning. Auditing the existing implementation found most of
this already present (feature 005's `NodeShape` text labels, `ResultsPanel`'s
worded delta) — the fix is small and targeted rather than a redesign of
how results are reported.

**Alternatives considered**: Adding icons instead of/alongside letters for
node types — considered, but a single bold letter is simpler to render
inside a small Konva shape than importing or drawing icon glyphs, and
satisfies FR-007 equally well (text, not color, carries the meaning).

## 5. Empty-canvas and missing-demand guidance render inline, not as a dialog

**Decision**: `emptyCanvasHint` renders as an overlay/message inside the
existing canvas area when there are no nodes yet; `missingDemandHint`
renders as a small notice near the "Run" button when there are no demands
yet. Neither opens a modal or a separate screen.

**Rationale**: Both are guidance about the *current* screen, not a
separate task — a modal would interrupt the exact place the user needs to
look (the canvas, or the Run button) to act on the guidance. Inline
placement keeps the fix scoped to "add missing text," not "add a new
interaction pattern," matching the spec's presentation-only framing.

**Alternatives considered**: A single global "what's wrong" banner
listing all current issues (no nodes, no demands, etc.) — rejected as
more indirection than two small, locally-placed hints; a user looking at
an empty canvas shouldn't have to look elsewhere for the explanation.
