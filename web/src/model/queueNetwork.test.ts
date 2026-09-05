import { describe, expect, it } from "vitest";
import { addDirectedRoad, addNode, emptyQueueNetwork, isQueueNetworkError, removeDirectedRoad, removeNodeFromQueueNetwork } from "./queueNetwork";

describe("addDirectedRoad", () => {
  it("adds a one-way road, optionally with a signal", () => {
    let state = addNode(emptyQueueNetwork, "house", 0, 0);
    state = addNode(state, "company", 100, 0);
    const [a, b] = state.nodes;

    const outcome = addDirectedRoad(state, a.id, b.id, { type: "constant", value: 5 }, { greenDuration: 5, redDuration: 5 });
    expect(isQueueNetworkError(outcome)).toBe(false);
    if (isQueueNetworkError(outcome)) throw new Error("unreachable");

    expect(outcome.roads).toHaveLength(1);
    expect(outcome.roads[0]).toMatchObject({ from: a.id, to: b.id, signal: { greenDuration: 5, redDuration: 5 } });
  });

  it("rejects a road from a node to itself", () => {
    const state = addNode(emptyQueueNetwork, "house", 0, 0);
    const [a] = state.nodes;

    const outcome = addDirectedRoad(state, a.id, a.id, { type: "constant", value: 5 });
    expect(isQueueNetworkError(outcome)).toBe(true);
  });
});

describe("removeNodeFromQueueNetwork", () => {
  it("cascades to roads referencing the removed node", () => {
    let state = addNode(emptyQueueNetwork, "house", 0, 0);
    state = addNode(state, "company", 100, 0);
    const [a, b] = state.nodes;

    const withRoad = addDirectedRoad(state, a.id, b.id, { type: "constant", value: 5 });
    if (isQueueNetworkError(withRoad)) throw new Error("unreachable");

    const result = removeNodeFromQueueNetwork(withRoad, a.id);

    expect(result.nodes.map((n) => n.id)).toEqual([b.id]);
    expect(result.roads).toHaveLength(0);
  });
});

describe("removeDirectedRoad", () => {
  it("removes a single road without affecting others", () => {
    let state = addNode(emptyQueueNetwork, "house", 0, 0);
    state = addNode(state, "company", 100, 0);
    const [a, b] = state.nodes;

    const withRoad = addDirectedRoad(state, a.id, b.id, { type: "constant", value: 5 });
    if (isQueueNetworkError(withRoad)) throw new Error("unreachable");
    const roadId = withRoad.roads[0].id;

    const result = removeDirectedRoad(withRoad, roadId);
    expect(result.roads).toHaveLength(0);
  });
});
