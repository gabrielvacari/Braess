import { describe, expect, it } from "vitest";
import { positionAtTime } from "./useAgentAnimation";

describe("positionAtTime", () => {
  const path = [
    { x: 0, y: 0 },
    { x: 10, y: 0 },
    { x: 10, y: 10 },
  ];

  it("is the route's start at elapsed <= 0", () => {
    expect(positionAtTime(path, 100, 0)).toEqual({ x: 0, y: 0 });
    expect(positionAtTime(path, 100, -50)).toEqual({ x: 0, y: 0 });
  });

  it("is the route's end at elapsed >= totalTime", () => {
    expect(positionAtTime(path, 100, 100)).toEqual({ x: 10, y: 10 });
    expect(positionAtTime(path, 100, 500)).toEqual({ x: 10, y: 10 });
  });

  it("is partway along the first segment at the halfway point of that segment's share", () => {
    // Two segments, equal share each (50 time units per segment):
    // at t=25 (halfway through segment 1) we should be halfway from (0,0) to (10,0).
    expect(positionAtTime(path, 100, 25)).toEqual({ x: 5, y: 0 });
  });

  it("is at the shared node exactly at the segment boundary", () => {
    expect(positionAtTime(path, 100, 50)).toEqual({ x: 10, y: 0 });
  });

  it("progresses monotonically in x-then-y across both segments", () => {
    const samples = [0, 10, 25, 40, 50, 60, 75, 90, 100].map((t) => positionAtTime(path, 100, t));
    for (let i = 1; i < samples.length; i++) {
      const prev = samples[i - 1];
      const cur = samples[i];
      // Distance from the start must never decrease.
      const prevDist = Math.abs(prev.x) + Math.abs(prev.y);
      const curDist = Math.abs(cur.x) + Math.abs(cur.y);
      expect(curDist).toBeGreaterThanOrEqual(prevDist);
    }
  });

  it("handles a single-point path by staying put", () => {
    expect(positionAtTime([{ x: 3, y: 4 }], 100, 50)).toEqual({ x: 3, y: 4 });
  });

  it("handles an empty path without throwing", () => {
    expect(positionAtTime([], 100, 50)).toEqual({ x: 0, y: 0 });
  });

  it("treats totalTime <= 0 as already arrived", () => {
    expect(positionAtTime(path, 0, 0)).toEqual({ x: 10, y: 10 });
  });
});
