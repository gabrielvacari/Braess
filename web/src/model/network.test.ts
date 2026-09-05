import { describe, expect, it } from "vitest";
import {
  addDemand,
  addNode,
  addRoad,
  emptyNetwork,
  expandRoadsToDirectedEdges,
  humanizeErrorMessage,
  isNetworkError,
  nodeDisplayLabel,
  removeNode,
  removeRoad,
} from "./network";

describe("addNode", () => {
  it("adds a node of the given type at the given position", () => {
    const state = addNode(emptyNetwork, "house", 10, 20);
    expect(state.nodes).toHaveLength(1);
    expect(state.nodes[0]).toMatchObject({ type: "house", x: 10, y: 20 });
  });
});

describe("addRoad", () => {
  it("adds a bidirectional road between two existing nodes", () => {
    let state = addNode(emptyNetwork, "house", 0, 0);
    state = addNode(state, "company", 100, 0);
    const [a, b] = state.nodes;

    const outcome = addRoad(state, a.id, b.id, { type: "constant", value: 5 });
    expect(isNetworkError(outcome)).toBe(false);
    if (isNetworkError(outcome)) throw new Error("unreachable");

    expect(outcome.roads).toHaveLength(1);
    expect(outcome.roads[0]).toMatchObject({ a: a.id, b: b.id });
  });

  // FR-008 (feature 006)
  it("rejects a road from a node to itself", () => {
    const state = addNode(emptyNetwork, "house", 0, 0);
    const [a] = state.nodes;

    const outcome = addRoad(state, a.id, a.id, { type: "constant", value: 5 });
    expect(isNetworkError(outcome)).toBe(true);
  });

  it("allows a second, parallel road between the same two nodes", () => {
    let state = addNode(emptyNetwork, "house", 0, 0);
    state = addNode(state, "company", 100, 0);
    const [a, b] = state.nodes;

    let outcome = addRoad(state, a.id, b.id, { type: "constant", value: 5 });
    if (isNetworkError(outcome)) throw new Error("unreachable");
    outcome = addRoad(outcome, a.id, b.id, { type: "constant", value: 9 });
    if (isNetworkError(outcome)) throw new Error("unreachable");

    expect(outcome.roads).toHaveLength(2);
  });
});

describe("removeNode", () => {
  // FR-005 (feature 005)
  it("cascades to every road and demand referencing the removed node", () => {
    let state = addNode(emptyNetwork, "house", 0, 0);
    state = addNode(state, "company", 100, 0);
    const [a, b] = state.nodes;

    let withRoad = addRoad(state, a.id, b.id, { type: "constant", value: 5 });
    if (isNetworkError(withRoad)) throw new Error("unreachable");
    withRoad = addDemand(withRoad, a.id, b.id, 10);

    const after = removeNode(withRoad, a.id);

    expect(after.nodes.map((n) => n.id)).toEqual([b.id]);
    expect(after.roads).toHaveLength(0);
    expect(after.demands).toHaveLength(0);
  });

  it("leaves unrelated roads and demands untouched", () => {
    let state = addNode(emptyNetwork, "house", 0, 0);
    state = addNode(state, "company", 100, 0);
    state = addNode(state, "intersection", 50, 50);
    const [a, b, c] = state.nodes;

    let withRoads = addRoad(state, a.id, b.id, { type: "constant", value: 5 });
    if (isNetworkError(withRoads)) throw new Error("unreachable");
    const unrelated = addRoad(withRoads, b.id, c.id, { type: "constant", value: 3 });
    if (isNetworkError(unrelated)) throw new Error("unreachable");

    const after = removeNode(unrelated, a.id);

    expect(after.roads).toHaveLength(1);
    expect(after.roads[0]).toMatchObject({ a: b.id, b: c.id });
  });
});

describe("removeRoad", () => {
  it("removes only the targeted road", () => {
    let state = addNode(emptyNetwork, "house", 0, 0);
    state = addNode(state, "company", 100, 0);
    const [a, b] = state.nodes;

    let withRoads = addRoad(state, a.id, b.id, { type: "constant", value: 5 });
    if (isNetworkError(withRoads)) throw new Error("unreachable");
    withRoads = addRoad(withRoads, a.id, b.id, { type: "constant", value: 9 });
    if (isNetworkError(withRoads)) throw new Error("unreachable");

    const toRemove = withRoads.roads[0].id;
    const after = removeRoad(withRoads, toRemove);

    expect(after.roads).toHaveLength(1);
    expect(after.roads[0].id).not.toBe(toRemove);
  });
});

describe("expandRoadsToDirectedEdges", () => {
  // FR-001, FR-002, FR-003 (feature 007)
  it("expands one road into exactly two directed edges, reversed, sharing travelTime", () => {
    let state = addNode(emptyNetwork, "house", 0, 0);
    state = addNode(state, "company", 100, 0);
    const [a, b] = state.nodes;
    const travelTime = { type: "constant" as const, value: 7 };

    const withRoad = addRoad(state, a.id, b.id, travelTime);
    if (isNetworkError(withRoad)) throw new Error("unreachable");

    const edges = expandRoadsToDirectedEdges(withRoad.roads);

    expect(edges).toHaveLength(2);
    const forward = edges.find((e) => e.from === a.id && e.to === b.id);
    const backward = edges.find((e) => e.from === b.id && e.to === a.id);
    expect(forward).toBeDefined();
    expect(backward).toBeDefined();
    expect(forward!.travelTime).toEqual(travelTime);
    expect(backward!.travelTime).toEqual(travelTime);
    expect(forward!.id).not.toBe(backward!.id);
  });

  it("produces exactly 2 * roads.length edges for several roads", () => {
    let state = addNode(emptyNetwork, "house", 0, 0);
    state = addNode(state, "company", 100, 0);
    state = addNode(state, "intersection", 50, 50);
    const [a, b, c] = state.nodes;

    let withRoads = addRoad(state, a.id, b.id, { type: "constant", value: 1 });
    if (isNetworkError(withRoads)) throw new Error("unreachable");
    withRoads = addRoad(withRoads, b.id, c.id, { type: "constant", value: 2 });
    if (isNetworkError(withRoads)) throw new Error("unreachable");

    expect(expandRoadsToDirectedEdges(withRoads.roads)).toHaveLength(4);
  });

  it("is deterministic: the same roads value always expands to the same edge ids", () => {
    let state = addNode(emptyNetwork, "house", 0, 0);
    state = addNode(state, "company", 100, 0);
    const [a, b] = state.nodes;
    const withRoad = addRoad(state, a.id, b.id, { type: "constant", value: 5 });
    if (isNetworkError(withRoad)) throw new Error("unreachable");

    const first = expandRoadsToDirectedEdges(withRoad.roads).map((e) => e.id);
    const second = expandRoadsToDirectedEdges(withRoad.roads).map((e) => e.id);

    expect(first).toEqual(second);
  });
});

describe("nodeDisplayLabel", () => {
  it("numbers nodes per type, in the order they were added", () => {
    let state = addNode(emptyNetwork, "house", 0, 0);
    state = addNode(state, "company", 0, 0);
    state = addNode(state, "house", 0, 0);
    const [h1, c1, h2] = state.nodes;

    expect(nodeDisplayLabel(state.nodes, h1.id)).toBe("House 1");
    expect(nodeDisplayLabel(state.nodes, c1.id)).toBe("Company 1");
    expect(nodeDisplayLabel(state.nodes, h2.id)).toBe("House 2");
  });

  it("falls back to the raw id for an unknown node", () => {
    expect(nodeDisplayLabel([], "node-999")).toBe("node-999");
  });
});

describe("humanizeErrorMessage", () => {
  it("replaces every raw node id in a server error message with its display label", () => {
    let state = addNode(emptyNetwork, "house", 0, 0);
    state = addNode(state, "company", 0, 0);
    const [h1, c1] = state.nodes;

    const raw = `demand 0 (origin="${h1.id}" destination="${c1.id}"): agent: no route exists: from "${h1.id}" to "${c1.id}"`;
    const humanized = humanizeErrorMessage(raw, state.nodes);

    expect(humanized).toBe('demand 0 (origin="House 1" destination="Company 1"): agent: no route exists: from "House 1" to "Company 1"');
  });

  it("doesn't let a shorter id corrupt a longer id that contains it", () => {
    let state = addNode(emptyNetwork, "house", 0, 0);
    for (let i = 0; i < 9; i++) state = addNode(state, "house", 0, 0);
    const tenth = state.nodes[9]; // shares the "node-1" prefix with the first node's id

    const raw = `no route exists: from "${tenth.id}" to "x"`;
    expect(humanizeErrorMessage(raw, state.nodes)).toBe('no route exists: from "House 10" to "x"');
  });
});
