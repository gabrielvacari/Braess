import { describe, expect, it } from "vitest";
import { queueChartPoints } from "./queueChart";

describe("queueChartPoints", () => {
  it("returns an empty array for no samples", () => {
    expect(queueChartPoints([], 100, 50)).toEqual([]);
  });

  it("maps samples to in-bounds coordinates, preserving time ordering", () => {
    const samples = [
      { time: 0, length: 0 },
      { time: 5, length: 10 },
      { time: 10, length: 4 },
    ];

    const points = queueChartPoints(samples, 100, 50);

    expect(points).toHaveLength(3);
    for (const p of points) {
      expect(p.x).toBeGreaterThanOrEqual(0);
      expect(p.x).toBeLessThanOrEqual(100);
      expect(p.y).toBeGreaterThanOrEqual(0);
      expect(p.y).toBeLessThanOrEqual(50);
    }

    // Time-ordered input stays x-ordered (earliest sample leftmost).
    expect(points[0].x).toBeLessThan(points[1].x);
    expect(points[1].x).toBeLessThan(points[2].x);

    // The highest queue length (index 1) reaches the top of the chart (y = 0).
    expect(points[1].y).toBe(0);
    // A zero-length sample sits at the bottom (y = height).
    expect(points[0].y).toBe(50);
  });

  it("does not divide by zero when every sample shares the same time", () => {
    const samples = [
      { time: 3, length: 1 },
      { time: 3, length: 2 },
    ];
    const points = queueChartPoints(samples, 100, 50);
    expect(points).toHaveLength(2);
    expect(points.every((p) => Number.isFinite(p.x) && Number.isFinite(p.y))).toBe(true);
  });
});
