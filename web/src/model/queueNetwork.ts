// The Signals-mode client model — fully separate from model/network.ts's
// equilibrium-mode NetworkState (FR-006, FR-007; research.md decision
// #2). Nodes mean the same thing in both modes, so ClientNode is reused;
// roads and demands do not, because a traffic signal inherently governs
// one direction of travel (research.md decision #1) and a time-based
// demand has an arrival rate the equilibrium mode's ClientDemand doesn't.
// This file contains no pathfinding or queuing logic — every number this
// mode ever shows comes from the server's POST /api/queue-run (see
// queueApi.ts), never computed here.
import { addNode, generateId, type ClientNode, type NodeType, type TravelTimeSpec } from "./network";

export interface SignalSpec {
  greenDuration: number;
  redDuration: number;
}

/** One direction of travel, optionally governed by a signal (FR-001). */
export interface DirectedRoad {
  id: string;
  from: string;
  to: string;
  travelTime: TravelTimeSpec;
  signal: SignalSpec | null;
}

/** Agents arriving at a steady rate over time (FR-002), unlike the equilibrium mode's fixed-count ClientDemand. */
export interface TimeDemand {
  origin: string;
  destination: string;
  count: number;
  arrivalInterval: number;
}

export interface QueueNetworkState {
  nodes: ClientNode[];
  roads: DirectedRoad[];
  demands: TimeDemand[];
}

export const emptyQueueNetwork: QueueNetworkState = { nodes: [], roads: [], demands: [] };

// Re-exported so SignalToolbar/SignalApp have one place to add nodes,
// consistent with the fact that a node means the same thing in both modes.
export { addNode };
export type { ClientNode, NodeType, TravelTimeSpec };

export interface QueueNetworkError {
  error: "self-loop";
}

export function isQueueNetworkError(result: QueueNetworkState | QueueNetworkError): result is QueueNetworkError {
  return "error" in result;
}

/** Adds a directed road, optionally with a signal. Rejects a self-loop (from === to), mirroring addRoad's rule. */
export function addDirectedRoad(
  state: QueueNetworkState,
  from: string,
  to: string,
  travelTime: TravelTimeSpec,
  signal: SignalSpec | null = null,
): QueueNetworkState | QueueNetworkError {
  if (from === to) {
    return { error: "self-loop" };
  }
  const road: DirectedRoad = { id: generateId("qroad"), from, to, travelTime, signal };
  return { ...state, roads: [...state.roads, road] };
}

/** Removes a single directed road. */
export function removeDirectedRoad(state: QueueNetworkState, id: string): QueueNetworkState {
  return { ...state, roads: state.roads.filter((r) => r.id !== id) };
}

/** Removes a node, cascading to every road and demand referencing it (mirrors removeNode). */
export function removeNodeFromQueueNetwork(state: QueueNetworkState, id: string): QueueNetworkState {
  return {
    nodes: state.nodes.filter((n) => n.id !== id),
    roads: state.roads.filter((r) => r.from !== id && r.to !== id),
    demands: state.demands.filter((d) => d.origin !== id && d.destination !== id),
  };
}

/** Adds a time-based demand. */
export function addTimeDemand(
  state: QueueNetworkState,
  origin: string,
  destination: string,
  count: number,
  arrivalInterval: number,
): QueueNetworkState {
  return { ...state, demands: [...state.demands, { origin, destination, count, arrivalInterval }] };
}

/** Removes the time-based demand at the given index. */
export function removeTimeDemand(state: QueueNetworkState, index: number): QueueNetworkState {
  return { ...state, demands: state.demands.filter((_, i) => i !== index) };
}
