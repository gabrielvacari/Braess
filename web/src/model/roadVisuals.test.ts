import { describe, expect, it } from "vitest";
import { roadChevronPositions } from "./roadVisuals";

describe("roadChevronPositions", () => {
  it("places marks evenly between the two nodes' edges, pointing a -> b", () => {
    const a = { x: 0, y: 0 };
    const b = { x: 200, y: 0 };
    const marks = roadChevronPositions(a, b, 18, 40);

    // start = 18 + 20 = 38, end = 200 - 18 = 182 -> 38, 78, 118, 158 (198 > 182, stops)
    expect(marks.map((m) => m.x)).toEqual([38, 78, 118, 158]);
    for (const m of marks) {
      expect(m.y).toBe(0);
      expect(m.dx).toBeCloseTo(1);
      expect(m.dy).toBeCloseTo(0);
    }
  });

  it("orients marks along a diagonal road", () => {
    const a = { x: 0, y: 0 };
    const b = { x: 0, y: 100 };
    const marks = roadChevronPositions(a, b, 10, 30);

    expect(marks.length).toBeGreaterThan(0);
    for (const m of marks) {
      expect(m.dx).toBeCloseTo(0);
      expect(m.dy).toBeCloseTo(1);
      expect(m.x).toBeCloseTo(0);
    }
  });

  it("returns no marks when the road is too short to fit one between the nodes' edges", () => {
    const a = { x: 0, y: 0 };
    const b = { x: 30, y: 0 }; // 30px apart, node radius 18 each -> no room at all
    expect(roadChevronPositions(a, b, 18, 40)).toEqual([]);
  });
});
