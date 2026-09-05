import { useState } from "react";
import type { ClientNode, NodeType, TimeDemand, TravelTimeSpec } from "../model/queueNetwork";
import { nodeDisplayLabel } from "../model/network";

export type SignalMode = "place" | "connect" | "select";

const MODES: { value: SignalMode; label: string }[] = [
  { value: "place", label: "Place node" },
  { value: "connect", label: "Draw road" },
  { value: "select", label: "Select / remove" },
];

const NODE_TYPES: NodeType[] = ["house", "company", "intersection"];

interface SegmentedProps<T extends string> {
  options: { value: T; label: string }[];
  value: T;
  onChange: (v: T) => void;
}

function Segmented<T extends string>({ options, value, onChange }: SegmentedProps<T>) {
  return (
    <div role="radiogroup" style={{ display: "flex", border: "1px solid var(--color-border)", borderRadius: "var(--radius-md)", overflow: "hidden" }}>
      {options.map((opt) => (
        <button
          key={opt.value}
          type="button"
          role="radio"
          aria-checked={value === opt.value}
          onClick={() => onChange(opt.value)}
          style={{
            flex: 1,
            border: "none",
            borderRadius: 0,
            background: value === opt.value ? "var(--color-accent)" : "var(--color-surface)",
            color: value === opt.value ? "var(--color-accent-contrast)" : "var(--color-text)",
          }}
        >
          {opt.label}
        </button>
      ))}
    </div>
  );
}

interface SignalToolbarProps {
  mode: SignalMode;
  onModeChange: (mode: SignalMode) => void;
  nodeTypeToPlace: NodeType;
  onNodeTypeChange: (type: NodeType) => void;
  travelTime: TravelTimeSpec;
  onTravelTimeChange: (t: TravelTimeSpec) => void;
  hasSignal: boolean;
  onHasSignalChange: (v: boolean) => void;
  greenDuration: number;
  redDuration: number;
  onGreenDurationChange: (v: number) => void;
  onRedDurationChange: (v: number) => void;
  connectingFrom: string | null;
  hasSelection: boolean;
  onRemoveSelected: () => void;
  nodes: ClientNode[];
  demands: TimeDemand[];
  onAddDemand: (origin: string, destination: string, count: number, arrivalInterval: number) => void;
  onRemoveDemand: (index: number) => void;
}

/** FR-001, FR-002: node/road(+signal) drawing controls and the time-based demand form for Signals mode. */
export function SignalToolbar({
  mode,
  onModeChange,
  nodeTypeToPlace,
  onNodeTypeChange,
  travelTime,
  onTravelTimeChange,
  hasSignal,
  onHasSignalChange,
  greenDuration,
  redDuration,
  onGreenDurationChange,
  onRedDurationChange,
  connectingFrom,
  hasSelection,
  onRemoveSelected,
  nodes,
  demands,
  onAddDemand,
  onRemoveDemand,
}: SignalToolbarProps) {
  const [origin, setOrigin] = useState("");
  const [destination, setDestination] = useState("");
  const [count, setCount] = useState(50);
  const [arrivalInterval, setArrivalInterval] = useState(1);

  return (
    <div style={{ display: "flex", flexDirection: "column", gap: "var(--space-4)" }}>
      <div className="card">
        <h3>Mode</h3>
        <Segmented options={MODES} value={mode} onChange={onModeChange} />
        {mode === "connect" && (
          <p className="text-muted" style={{ marginTop: "var(--space-2)" }}>
            {connectingFrom
              ? `Click a second node to connect it to "${connectingFrom}" — the road runs one way, from the first node to the second.`
              : "Click a node to start drawing a one-way road from it."}
          </p>
        )}
      </div>

      {mode === "place" && (
        <div className="card">
          <h3>Node type</h3>
          <Segmented options={NODE_TYPES.map((t) => ({ value: t, label: t }))} value={nodeTypeToPlace} onChange={onNodeTypeChange} />
        </div>
      )}

      {mode === "connect" && (
        <div className="card">
          <h3>Road behavior</h3>
          <Segmented
            options={[
              { value: "constant", label: "Constant" },
              { value: "linear", label: "Gets slower with traffic" },
            ]}
            value={travelTime.type}
            onChange={(t) =>
              // freeFlow: 2, not 0 — a signal-controlled edge with a
              // zero free-flow time is degenerate (queuesim/run.go
              // treats it as instantaneous rather than dividing by it,
              // but a road that takes no time at all isn't a meaningful
              // scenario to hand a person by default).
              onTravelTimeChange(t === "constant" ? { type: "constant", value: 5 } : { type: "linear", freeFlow: 2, slope: 0.1 })
            }
          />

          <h3 style={{ marginTop: "var(--space-4)" }}>Traffic signal</h3>
          <label style={{ display: "flex", alignItems: "center", gap: "var(--space-2)" }}>
            <input type="checkbox" checked={hasSignal} onChange={(e) => onHasSignalChange(e.target.checked)} />
            Attach a signal to this road
          </label>
          {hasSignal && (
            <div style={{ display: "flex", gap: "var(--space-3)", marginTop: "var(--space-2)" }}>
              <label style={{ display: "flex", flexDirection: "column", gap: "var(--space-1)" }}>
                <span className="text-muted" style={{ fontSize: "var(--text-xs)" }}>
                  Green (s)
                </span>
                <input type="number" min={0} value={greenDuration} onChange={(e) => onGreenDurationChange(Number(e.target.value))} style={{ width: "4.5rem" }} />
              </label>
              <label style={{ display: "flex", flexDirection: "column", gap: "var(--space-1)" }}>
                <span className="text-muted" style={{ fontSize: "var(--text-xs)" }}>
                  Red (s)
                </span>
                <input type="number" min={0} value={redDuration} onChange={(e) => onRedDurationChange(Number(e.target.value))} style={{ width: "4.5rem" }} />
              </label>
            </div>
          )}
        </div>
      )}

      {mode === "select" && (
        <div className="card">
          <h3>Selection</h3>
          <button type="button" disabled={!hasSelection} onClick={onRemoveSelected}>
            Remove selected
          </button>
        </div>
      )}

      <div className="card">
        <h3>Time-based demands</h3>
        <p className="text-muted">Agents arrive at a steady rate over the run, not all at once.</p>
        <div style={{ display: "flex", gap: "var(--space-2)", flexWrap: "wrap", margin: "var(--space-2) 0 var(--space-3)" }}>
          <select aria-label="Origin" value={origin} onChange={(e) => setOrigin(e.target.value)}>
            <option value="">origin…</option>
            {nodes.map((n) => (
              <option key={n.id} value={n.id}>
                {nodeDisplayLabel(nodes, n.id)}
              </option>
            ))}
          </select>
          <select aria-label="Destination" value={destination} onChange={(e) => setDestination(e.target.value)}>
            <option value="">destination…</option>
            {nodes.map((n) => (
              <option key={n.id} value={n.id}>
                {nodeDisplayLabel(nodes, n.id)}
              </option>
            ))}
          </select>
          <label style={{ display: "flex", flexDirection: "column", gap: "var(--space-1)" }}>
            <span className="text-muted" style={{ fontSize: "var(--text-xs)" }}>
              Total agents
            </span>
            <input aria-label="Total agents" type="number" min={0} value={count} onChange={(e) => setCount(Number(e.target.value))} style={{ width: "4.5rem" }} />
          </label>
          <label style={{ display: "flex", flexDirection: "column", gap: "var(--space-1)" }}>
            <span className="text-muted" style={{ fontSize: "var(--text-xs)" }}>
              Arrival every (s)
            </span>
            <input
              aria-label="Seconds between arrivals"
              type="number"
              min={0.1}
              step={0.1}
              value={arrivalInterval}
              onChange={(e) => setArrivalInterval(Number(e.target.value))}
              style={{ width: "4.5rem" }}
            />
          </label>
          <button
            type="button"
            className="primary"
            disabled={!origin || !destination}
            onClick={() => {
              onAddDemand(origin, destination, count, arrivalInterval);
              setOrigin("");
              setDestination("");
            }}
          >
            Add
          </button>
        </div>
        <ul style={{ display: "flex", flexDirection: "column", gap: "var(--space-1)" }}>
          {demands.map((d, i) => (
            <li key={`${d.origin}-${d.destination}-${i}`} style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
              <span>
                {nodeDisplayLabel(nodes, d.origin)} → {nodeDisplayLabel(nodes, d.destination)} ({d.count} agents, one every {d.arrivalInterval}s)
              </span>
              <button type="button" className="icon" aria-label={`Remove demand ${d.origin} to ${d.destination}`} onClick={() => onRemoveDemand(i)}>
                ✕
              </button>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}
