import { APP_INTRO } from "../content/copy";

/** FR-001: what the tool demonstrates, always visible as soon as the app loads. */
export function IntroBanner() {
  return (
    <div className="card" style={{ background: "var(--color-surface-alt)", borderStyle: "dashed" }}>
      <p>{APP_INTRO}</p>
    </div>
  );
}
