import { NODE_TYPE_LEGEND } from "../content/copy";

const SWATCH_COLOR: Record<string, string> = {
  house: "#2f9e44",
  company: "#f08c00",
  intersection: "#868e96",
};

/** FR-002: maps every node marker to its name and meaning. */
export function Legend() {
  return (
    <div className="card">
      <h3>Legend</h3>
      <ul style={{ display: "flex", flexDirection: "column", gap: "var(--space-2)" }}>
        {NODE_TYPE_LEGEND.map((entry) => (
          <li key={entry.type} style={{ display: "flex", alignItems: "flex-start", gap: "var(--space-2)" }}>
            <span
              aria-hidden="true"
              style={{
                display: "inline-flex",
                alignItems: "center",
                justifyContent: "center",
                width: "1.5rem",
                height: "1.5rem",
                borderRadius: entry.type === "company" ? "50%" : "var(--radius-sm)",
                background: SWATCH_COLOR[entry.type],
                color: "#fff",
                fontSize: "var(--text-xs)",
                fontWeight: 700,
                flexShrink: 0,
              }}
            >
              {entry.label}
            </span>
            <span>
              <strong>{entry.name}</strong>
              <br />
              <span className="text-muted">{entry.description}</span>
            </span>
          </li>
        ))}
      </ul>
    </div>
  );
}
