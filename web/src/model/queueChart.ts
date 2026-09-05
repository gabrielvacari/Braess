// Pure chart geometry, kept separate from QueueChart.tsx's rendering so
// it can be tested without rendering anything — the same split
// model/network.ts already established for edit operations and
// content/copy.ts for explanatory strings.

export interface ChartSample {
  time: number;
  length: number;
}

export interface ChartPoint {
  x: number;
  y: number;
}

/**
 * Maps time-ordered samples to SVG coordinates inside a width x height
 * box: time -> x (left to right, earliest first), queue length -> y
 * (0 at the bottom, taller queues higher up). Returns an empty array for
 * no samples so callers can render an empty-state instead of a chart.
 */
export function queueChartPoints(samples: ChartSample[], width: number, height: number): ChartPoint[] {
  if (samples.length === 0) return [];

  const minTime = samples[0].time;
  const maxTime = samples[samples.length - 1].time;
  const timeSpan = maxTime - minTime || 1;
  const maxLength = Math.max(1, ...samples.map((s) => s.length));

  return samples.map((s) => ({
    x: ((s.time - minTime) / timeSpan) * width,
    y: height - (s.length / maxLength) * height,
  }));
}
