import { describe, expect, it } from "vitest";
import { emptyCanvasHint, missingDemandHint, modeDescription, noNodesForDemandHint, NODE_TYPE_LEGEND } from "./copy";

describe("NODE_TYPE_LEGEND", () => {
  it("has exactly one entry per node type", () => {
    const types = NODE_TYPE_LEGEND.map((e) => e.type);
    expect(types).toEqual(["house", "company", "intersection"]);
  });
});

describe("modeDescription", () => {
  // FR-003: every mode has a non-empty, distinct description.
  it("returns a distinct, non-empty description for each mode", () => {
    const place = modeDescription("place", null);
    const connect = modeDescription("connect", null);
    const select = modeDescription("select", null);

    for (const text of [place, connect, select]) {
      expect(text.length).toBeGreaterThan(0);
    }
    expect(new Set([place, connect, select]).size).toBe(3);
  });

  it("describes 'connect' differently before and after the first click", () => {
    const before = modeDescription("connect", null);
    const after = modeDescription("connect", "house-1");

    expect(before).not.toBe(after);
    expect(after).toContain("house-1");
  });
});

describe("emptyCanvasHint", () => {
  it("suggests a first action when there are no nodes", () => {
    expect(emptyCanvasHint(0)).not.toBeNull();
  });

  it("is null once at least one node exists", () => {
    expect(emptyCanvasHint(1)).toBeNull();
    expect(emptyCanvasHint(5)).toBeNull();
  });
});

describe("missingDemandHint", () => {
  it("explains what's missing when there are no demands", () => {
    expect(missingDemandHint(0)).not.toBeNull();
  });

  it("is null once at least one demand exists", () => {
    expect(missingDemandHint(1)).toBeNull();
  });
});

describe("noNodesForDemandHint", () => {
  it("explains why the Add button is disabled with fewer than 2 nodes", () => {
    expect(noNodesForDemandHint(0)).not.toBeNull();
    expect(noNodesForDemandHint(1)).not.toBeNull();
  });

  it("is null once there are at least 2 nodes to choose from", () => {
    expect(noNodesForDemandHint(2)).toBeNull();
    expect(noNodesForDemandHint(5)).toBeNull();
  });
});
