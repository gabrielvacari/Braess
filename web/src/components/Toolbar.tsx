import { useState } from "react";
import { modeDescription, noNodesForDemandHint } from "../content/copy";
import { nodeDisplayLabel, type ClientDemand, type ClientNode, type NodeType, type TravelTimeSpec } from "../model/network";

export type EditMode = "place" | "connect" | "select";

interface ToolbarProps {
  mode: EditMode;
  onModeChange: (mode: EditMode) => void;
  nodeTypeToPlace: NodeType;
  onNodeTypeChange: (type: NodeType) => void;
  travelTimeToUse: TravelTimeSpec;
  onTravelTimeChange: (spec: TravelTimeSpec) => void;
  hasSelection: boolean;
  onRemoveSelected: () => void;
  connectingFrom: string | null;
  nodes: ClientNode[];
  demands: ClientDemand[];
  onAddDemand: (origin: string, destination: string, size: number) => void;
  onRemoveDemand: (index: number) => void;
}

const MODES: { value: EditMode; label: string }[] = [
  { value: "place", label: "Place node" },
  { value: "connect", label: "Draw road" },
  { value: "select", label: "Select / remove" },
];

const NODE_TYPES: NodeType[] = ["house", "company", "intersection"];

function SegmentedControl<T extends string>({
  options,
  value,
  onChange,
}: {
  options: { value: T; label: string }[];
  value: T;
  onChange: (v: T) => void;
}) {
  return (
    <div
      role="radiogroup"
      style={{ display: "flex", border: "1px solid var(--color-border)", borderRadius: "var(--radius-md)", overflow: "hidden" }}
    >
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

export function Toolbar({
  mode,
  onModeChange,
  nodeTypeToPlace,
  onNodeTypeChange,
  travelTimeToUse,
  onTravelTimeChange,
  hasSelection,
  onRemoveSelected,
  connectingFrom,
  nodes,
  demands,
  onAddDemand,
  onRemoveDemand,
}: ToolbarProps) {
  const [origin, setOrigin] = useState("");
  const [destination, setDestination] = useState("");
  const [size, setSize] = useState(10);

  return (
    <div style={{ display: "flex", flexDirection: "column", gap: "var(--space-4)" }}>
      <div className="card">
        <h3>Mode</h3>
        <SegmentedControl options={MODES} value={mode} onChange={onModeChange} />
        {/* FR-003: what clicking will do in the active mode, right now. */}
        <p className="text-muted" style={{ marginTop: "var(--space-2)" }}>
          {modeDescription(mode, connectingFrom)}
        </p>
      </div>

      {mode === "place" && (
        <div className="card">
          <h3>Node type</h3>
          <SegmentedControl
            options={NODE_TYPES.map((t) => ({ value: t, label: t }))}
            value={nodeTypeToPlace}
            onChange={onNodeTypeChange}
          />
        </div>
      )}

      {mode === "connect" && (
        <div className="card">
          <h3>Road behavior</h3>
          <SegmentedControl
            options={[
              { value: "constant", label: "Constant" },
              { value: "linear", label: "Gets slower with traffic" },
            ]}
            value={travelTimeToUse.type}
            onChange={(t) =>
              onTravelTimeChange(t === "constant" ? { type: "constant", value: 5 } : { type: "linear", freeFlow: 0, slope: 0.1 })
            }
          />
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
        <h3>Demands</h3>
        {noNodesForDemandHint(nodes.length) && (
          <p className="text-muted" style={{ marginBottom: "var(--space-2)" }}>
            {noNodesForDemandHint(nodes.length)}
          </p>
        )}
        <div style={{ display: "flex", gap: "var(--space-3)", flexWrap: "wrap", alignItems: "flex-end", marginBottom: "var(--space-3)" }}>
          <label style={{ display: "flex", flexDirection: "column", gap: "var(--space-1)" }}>
            <span className="text-muted" style={{ fontSize: "var(--text-xs)" }}>
              From
            </span>
            <select aria-label="Origin" value={origin} onChange={(e) => setOrigin(e.target.value)}>
              <option value="">choose a node…</option>
              {nodes.map((n) => (
                <option key={n.id} value={n.id}>
                  {nodeDisplayLabel(nodes, n.id)}
                </option>
              ))}
            </select>
          </label>
          <label style={{ display: "flex", flexDirection: "column", gap: "var(--space-1)" }}>
            <span className="text-muted" style={{ fontSize: "var(--text-xs)" }}>
              To
            </span>
            <select aria-label="Destination" value={destination} onChange={(e) => setDestination(e.target.value)}>
              <option value="">choose a node…</option>
              {nodes.map((n) => (
                <option key={n.id} value={n.id}>
                  {nodeDisplayLabel(nodes, n.id)}
                </option>
              ))}
            </select>
          </label>
          <label style={{ display: "flex", flexDirection: "column", gap: "var(--space-1)" }}>
            <span className="text-muted" style={{ fontSize: "var(--text-xs)" }}>
              Agents
            </span>
            <input
              aria-label="Number of agents making this trip"
              title="How many agents will travel from the origin to the destination"
              type="number"
              min={0}
              value={size}
              onChange={(e) => setSize(Number(e.target.value))}
              style={{ width: "4.5rem" }}
            />
          </label>
          <button
            type="button"
            className="primary"
            disabled={!origin || !destination}
            onClick={() => {
              onAddDemand(origin, destination, size);
              setOrigin("");
              setDestination("");
              setSize(10);
            }}
          >
            Add
          </button>
        </div>
        <ul style={{ display: "flex", flexDirection: "column", gap: "var(--space-1)" }}>
          {demands.map((d, i) => (
            <li
              key={`${d.origin}-${d.destination}-${i}`}
              style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}
            >
              <span>
                {nodeDisplayLabel(nodes, d.origin)} → {nodeDisplayLabel(nodes, d.destination)} ({d.size} agents)
              </span>
              <button
                type="button"
                className="icon"
                aria-label={`Remove demand ${nodeDisplayLabel(nodes, d.origin)} to ${nodeDisplayLabel(nodes, d.destination)}`}
                onClick={() => onRemoveDemand(i)}
              >
                ✕
              </button>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}
