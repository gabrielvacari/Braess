import { assignHousesToCompanies, AUTO_DEMAND_ARRIVAL_INTERVAL, AUTO_DEMAND_COUNT } from "../model/autoDemand";
import { nodeDisplayLabel } from "../model/network";
import type { ClientNode, NodeType, TravelTimeSpec } from "../model/queueNetwork";

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
}

/**
 * FR-001: node/road(+signal) drawing controls for Signals mode. There is
 * no demand form (feature 012) — every house always sends traffic
 * automatically; this only shows, read-only, which company each house
 * has been assigned to (FR-010), computed straight from `nodes`.
 */
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
}: SignalToolbarProps) {
  const assignments = assignHousesToCompanies(nodes);

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
        <h3>Automatic demand</h3>
        <p className="text-muted">
          Every house sends {AUTO_DEMAND_COUNT} agents every {AUTO_DEMAND_ARRIVAL_INTERVAL}s — there's nothing to
          configure. The code picks which company each house's agents go to, spreading across every company you've
          drawn rather than piling onto just one.
        </p>
        {assignments.length === 0 ? (
          <p className="text-muted">
            {nodes.some((n) => n.type === "house")
              ? "Add a company for houses to send traffic to."
              : "Add a house and a company to see traffic here."}
          </p>
        ) : (
          <ul style={{ display: "flex", flexDirection: "column", gap: "var(--space-1)" }}>
            {assignments.map((a) => (
              <li key={a.houseId}>
                {nodeDisplayLabel(nodes, a.houseId)} → {nodeDisplayLabel(nodes, a.companyId)}
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}
