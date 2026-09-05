// The client-side network model: what the user is drawing/editing in the
// browser. This file contains no pathfinding, no cost functions, nothing
// that mirrors the Go engine's graph/agent/simulation packages — it only
// holds and edits plain node/road/demand data. Every route and every
// travel time the app ever shows comes from the server's POST /api/run
// (see api.ts); the client never decides a route itself (constitution
// Principle III boundary — see plan.md research.md decision #3).
//
// The editable primitive is the ROAD (bidirectional — matches how a user
// thinks about "a road"), not the directed edge the engine models
// (feature 001). A road expands into two directed edges only at the API
// boundary — see expandRoadsToDirectedEdges — so drawing, displaying, and
// removing a road are each exactly one operation on exactly one record
// (feature 007, FR-004/FR-005).

export type NodeType = "house" | "company" | "intersection";

export type TravelTimeSpec =
  | { type: "linear"; freeFlow: number; slope: number }
  | { type: "constant"; value: number };

export interface ClientNode {
  id: string;
  type: NodeType;
  x: number;
  y: number;
}

/** A bidirectional road between two nodes — no "from"/"to" direction. */
export interface ClientRoad {
  id: string;
  a: string;
  b: string;
  /** Shared by both directions (feature 007 Assumptions). */
  travelTime: TravelTimeSpec;
}

export interface ClientDemand {
  origin: string;
  destination: string;
  size: number;
}

export interface NetworkState {
  nodes: ClientNode[];
  roads: ClientRoad[];
  demands: ClientDemand[];
}

export const emptyNetwork: NetworkState = { nodes: [], roads: [], demands: [] };

let nextId = 1;

/** Generates a short, unique-within-this-page-load id. Exported only for tests. */
export function generateId(prefix: string): string {
  const id = `${prefix}-${nextId}`;
  nextId += 1;
  return id;
}

/** Adds a node of the given type at (x, y). */
export function addNode(state: NetworkState, type: NodeType, x: number, y: number): NetworkState {
  const node: ClientNode = { id: generateId("node"), type, x, y };
  return { ...state, nodes: [...state.nodes, node] };
}

export interface NetworkError {
  error: "self-loop";
}

function isNetworkError(result: NetworkState | NetworkError): result is NetworkError {
  return "error" in result;
}

/**
 * Adds a bidirectional road between two existing nodes. Rejects a
 * self-loop (a === b) rather than ever sending one to the server
 * (FR-008, feature 006) — a road cannot connect a node to itself.
 */
export function addRoad(
  state: NetworkState,
  a: string,
  b: string,
  travelTime: TravelTimeSpec,
): NetworkState | NetworkError {
  if (a === b) {
    return { error: "self-loop" };
  }
  const road: ClientRoad = { id: generateId("road"), a, b, travelTime };
  return { ...state, roads: [...state.roads, road] };
}

/**
 * Removes a node, cascading to every road that references it and every
 * demand whose origin or destination is it (FR-005) — never leaves a
 * road or demand referencing a node that no longer exists.
 */
export function removeNode(state: NetworkState, id: string): NetworkState {
  return {
    nodes: state.nodes.filter((n) => n.id !== id),
    roads: state.roads.filter((r) => r.a !== id && r.b !== id),
    demands: state.demands.filter((d) => d.origin !== id && d.destination !== id),
  };
}

/** Removes a single road — both its directions disappear with it (FR-005). */
export function removeRoad(state: NetworkState, id: string): NetworkState {
  return { ...state, roads: state.roads.filter((r) => r.id !== id) };
}

/** Adds a demand (origin, destination, size). */
export function addDemand(state: NetworkState, origin: string, destination: string, size: number): NetworkState {
  return { ...state, demands: [...state.demands, { origin, destination, size }] };
}

/** Removes the demand at the given index. */
export function removeDemand(state: NetworkState, index: number): NetworkState {
  return { ...state, demands: state.demands.filter((_, i) => i !== index) };
}

/**
 * A human-readable label for a node, e.g. "House 1", "Company 2" — used
 * anywhere a node is picked or displayed (the demand form and list), so
 * a person never has to make sense of a raw generated id like "node-4".
 * Numbering is per type, in the order those nodes were added; falls back
 * to the raw id if it doesn't match any node (e.g. a stale reference).
 */
export function nodeDisplayLabel(nodes: ClientNode[], id: string): string {
  const node = nodes.find((n) => n.id === id);
  if (!node) return id;
  const ordinal = nodes.filter((n) => n.type === node.type).findIndex((n) => n.id === id) + 1;
  return `${node.type.charAt(0).toUpperCase()}${node.type.slice(1)} ${ordinal}`;
}

/**
 * Replaces every raw node id appearing in a message with its display
 * label — for errors that come from the server (cmd/server, feature
 * 005), which only ever knows raw ids, never this client's friendly
 * names. Longer ids are replaced first so one id being a substring of
 * another (e.g. "node-1" inside "node-10") can't corrupt the longer one.
 */
export function humanizeErrorMessage(message: string, nodes: ClientNode[]): string {
  const byLongestId = [...nodes].sort((a, b) => b.id.length - a.id.length);
  return byLongestId.reduce((text, node) => text.split(node.id).join(nodeDisplayLabel(nodes, node.id)), message);
}

/** One directed connection derived from a ClientRoad — API-boundary only, never stored. */
export interface DirectedEdge {
  id: string;
  from: string;
  to: string;
  travelTime: TravelTimeSpec;
}

/**
 * Expands every road into its two directed edges (a->b and b->a), with
 * deterministic ids so the same `roads` value always expands the same
 * way — required so a request built from a NetworkState and a later
 * response interpreted against that same NetworkState agree on what each
 * edge id means (feature 007 research.md decisions #2, #3). This is the
 * only place a "road" becomes a directed "edge" — everywhere else in the
 * client, roads are the only thing that exists.
 */
export function expandRoadsToDirectedEdges(roads: ClientRoad[]): DirectedEdge[] {
  return roads.flatMap((road) => [
    { id: `${road.id}-ab`, from: road.a, to: road.b, travelTime: road.travelTime },
    { id: `${road.id}-ba`, from: road.b, to: road.a, travelTime: road.travelTime },
  ]);
}

export { isNetworkError };
