// A hand-drawn line chart for one signal's queue length over time — no
// charting library (research.md decision #4). The geometry itself lives
// in model/queueChart.ts (pure, tested independently of rendering); this
// file is only the SVG wrapper around it.
import { queueChartPoints, type ChartSample } from "../model/queueChart";

interface QueueChartProps {
  title: string;
  samples: ChartSample[];
  width?: number;
  height?: number;
}

/** One signal's queue-length-over-time chart, labeled with the road/signal it belongs to (FR-003). */
export function QueueChart({ title, samples, width = 360, height = 140 }: QueueChartProps) {
  const points = queueChartPoints(samples, width, height);
  const maxLength = samples.length > 0 ? Math.max(...samples.map((s) => s.length)) : 0;

  return (
    <div className="card">
      <h3>{title}</h3>
      {points.length === 0 ? (
        <p className="text-muted">No queue samples recorded for this signal.</p>
      ) : (
        <>
          <svg width={width} height={height} role="img" aria-label={`${title}: queue length over time, peaking at ${maxLength} agents`}>
            <polyline
              points={points.map((p) => `${p.x},${p.y}`).join(" ")}
              fill="none"
              stroke="#2f6fed"
              strokeWidth={2}
            />
          </svg>
          <p className="text-muted">Peak queue: {maxLength} agents</p>
        </>
      )}
    </div>
  );
}
