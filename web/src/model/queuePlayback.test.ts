import { describe, expect, it } from "vitest";
import type { PositionSampleDTO, QueueSampleDTO } from "./queueApi";
import { groupPositionsByAgent, positionAtSimTime, signalGreenAtSimTime } from "./queuePlayback";

function sample(overrides: Partial<PositionSampleDTO>): PositionSampleDTO {
  return { time: 0, agentId: 0, edgeId: "road", progress: 0, queued: false, ...overrides };
}

describe("groupPositionsByAgent", () => {
  it("groups samples by agentId, preserving recorded order", () => {
    const positions = [
      sample({ agentId: 0, time: 0 }),
      sample({ agentId: 1, time: 0 }),
      sample({ agentId: 0, time: 1 }),
      sample({ agentId: 1, time: 1 }),
    ];

    const grouped = groupPositionsByAgent(positions);

    expect(grouped.size).toBe(2);
    expect(grouped.get(0)?.map((p) => p.time)).toEqual([0, 1]);
    expect(grouped.get(1)?.map((p) => p.time)).toEqual([0, 1]);
  });
});

describe("positionAtSimTime", () => {
  const edgesById = new Map([["road", { from: "a", to: "b" }]]);
  const nodesById = new Map([
    ["a", { x: 0, y: 0 }],
    ["b", { x: 100, y: 0 }],
  ]);

  it("returns null before the agent's first sample", () => {
    const samples = [sample({ time: 5, progress: 0.5 })];
    expect(positionAtSimTime(samples, edgesById, nodesById, 2)).toBeNull();
  });

  it("interpolates along the edge from -> to by progress while moving", () => {
    const samples = [sample({ time: 0, progress: 0.25, queued: false })];
    const point = positionAtSimTime(samples, edgesById, nodesById, 0);
    expect(point).toEqual({ x: 25, y: 0 });
  });

  it("holds at the edge's start while queued, ignoring progress", () => {
    const samples = [sample({ time: 0, progress: 0, queued: true })];
    const point = positionAtSimTime(samples, edgesById, nodesById, 0);
    expect(point).toEqual({ x: 0, y: 0 });
  });

  it("uses the last sample at or before the queried time", () => {
    const samples = [
      sample({ time: 0, progress: 0 }),
      sample({ time: 1, progress: 0.5 }),
      sample({ time: 2, progress: 1 }),
    ];
    expect(positionAtSimTime(samples, edgesById, nodesById, 1.9)).toEqual({ x: 50, y: 0 });
    expect(positionAtSimTime(samples, edgesById, nodesById, 2)).toEqual({ x: 100, y: 0 });
    expect(positionAtSimTime(samples, edgesById, nodesById, 100)).toEqual({ x: 100, y: 0 });
  });

  it("returns null if the edge or its endpoints are unknown", () => {
    const samples = [sample({ time: 0, edgeId: "missing" })];
    expect(positionAtSimTime(samples, edgesById, nodesById, 0)).toBeNull();
  });
});

describe("signalGreenAtSimTime", () => {
  const samples: QueueSampleDTO[] = [
    { time: 0, signalId: "s1", length: 0, green: false },
    { time: 1, signalId: "s1", length: 0, green: true },
    { time: 2, signalId: "s1", length: 0, green: false },
  ];

  it("returns the nearest recorded sample at or before the queried time", () => {
    expect(signalGreenAtSimTime(samples, "s1", 0.5)).toBe(false);
    expect(signalGreenAtSimTime(samples, "s1", 1.5)).toBe(true);
    expect(signalGreenAtSimTime(samples, "s1", 2)).toBe(false);
  });

  it("falls back to the earliest sample when queried before it", () => {
    expect(signalGreenAtSimTime(samples, "s1", -1)).toBe(false);
  });

  it("returns false for a signal with no recorded samples", () => {
    expect(signalGreenAtSimTime(samples, "unknown", 1)).toBe(false);
  });
});
