import { describe, expect, it } from "vitest";
import {
  addDirectedRoad,
  addNode,
  addTimeDemand,
  emptyQueueNetwork,
  isQueueNetworkError,
  removeDirectedRoad,
  removeNodeFromQueueNetwork,
  removeTimeDemand,
} from "./queueNetwork";

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
  it("cascades to roads and demands referencing the removed node", () => {
    let state = addNode(emptyQueueNetwork, "house", 0, 0);
    state = addNode(state, "company", 100, 0);
    const [a, b] = state.nodes;

    const withRoad = addDirectedRoad(state, a.id, b.id, { type: "constant", value: 5 });
    if (isQueueNetworkError(withRoad)) throw new Error("unreachable");
    const withDemand = addTimeDemand(withRoad, a.id, b.id, 10, 1);

    const result = removeNodeFromQueueNetwork(withDemand, a.id);

    expect(result.nodes.map((n) => n.id)).toEqual([b.id]);
    expect(result.roads).toHaveLength(0);
    expect(result.demands).toHaveLength(0);
  });
});

describe("addTimeDemand / removeTimeDemand", () => {
  it("adds a time-based demand with an arrival interval", () => {
    const state = addTimeDemand(emptyQueueNetwork, "n1", "n2", 20, 2.5);
    expect(state.demands).toHaveLength(1);
    expect(state.demands[0]).toMatchObject({ origin: "n1", destination: "n2", count: 20, arrivalInterval: 2.5 });
  });

  it("removes the demand at the given index", () => {
    let state = addTimeDemand(emptyQueueNetwork, "n1", "n2", 20, 2.5);
    state = addTimeDemand(state, "n2", "n3", 5, 1);

    const result = removeTimeDemand(state, 0);

    expect(result.demands).toHaveLength(1);
    expect(result.demands[0]).toMatchObject({ origin: "n2", destination: "n3" });
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
