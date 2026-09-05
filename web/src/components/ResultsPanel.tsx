import type { RunResponse } from "../model/api";

interface ResultsPanelProps {
  current: RunResponse | null;
  previous: RunResponse | null;
}

function Summary({ label, result }: { label: string; result: RunResponse }) {
  return (
    <div>
      <h3>{label}</h3>
      <dl style={{ margin: 0, display: "grid", gridTemplateColumns: "auto auto", gap: "var(--space-1) var(--space-3)" }}>
        <dt className="text-muted">Converged</dt>
        <dd style={{ margin: 0 }}>{result.converged ? "yes" : "no"}</dd>
        <dt className="text-muted">Rounds</dt>
        <dd style={{ margin: 0 }}>{result.rounds}</dd>
        <dt className="text-muted">Total travel time</dt>
        <dd style={{ margin: 0 }}>{result.totalTravelTime}</dd>
        <dt className="text-muted">Average travel time</dt>
        <dd style={{ margin: 0 }}>{result.averageTravelTime}</dd>
      </dl>
    </div>
  );
}

/**
 * Shows the current run's result and, once a second run has happened,
 * the previous one alongside it — so a user can tell whether an edit
 * (e.g. adding a road) helped or hurt (FR-007, SC-004).
 */
export function ResultsPanel({ current, previous }: ResultsPanelProps) {
  if (!current) return null;

  const delta = previous ? current.averageTravelTime - previous.averageTravelTime : null;
  // FR-007: the comparison's meaning is stated in words, not carried by color alone.
  const deltaText =
    delta === null
      ? null
      : delta > 0
        ? `▲ worse by ${delta.toFixed(2)}`
        : delta < 0
          ? `▼ better by ${Math.abs(delta).toFixed(2)}`
          : "unchanged";
  const deltaColor = delta === null || delta === 0 ? "var(--color-text-muted)" : delta > 0 ? "var(--color-danger)" : "var(--color-success)";

  return (
    <div className="card">
      <div style={{ display: "flex", gap: "var(--space-6)", flexWrap: "wrap" }}>
        {previous && <Summary label="Previous run" result={previous} />}
        <Summary label="Current run" result={current} />
      </div>
      {deltaText && (
        <p style={{ marginTop: "var(--space-3)", fontWeight: 600, color: deltaColor }}>
          Average travel time: {deltaText}
        </p>
      )}
    </div>
  );
}
