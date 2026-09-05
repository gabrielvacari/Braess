// The POST /api/run request/response shapes, mirroring cmd/server/dto.go
// exactly (see specs/005-web-frontend/contracts/api-contract.md). Keeping
// these types hand-written (not generated) is why the project uses
// TypeScript here — a shape mismatch against the Go side becomes a
// compile error instead of a runtime surprise (research.md decision #5).
import { expandRoadsToDirectedEdges, type NetworkState, type TravelTimeSpec } from "./network";

export interface NodeDTO {
  id: string;
  type: string;
}

export interface EdgeDTO {
  id: string;
  from: string;
  to: string;
  length: number;
  capacity: number;
  travelTime: TravelTimeSpec;
}

export interface DemandDTO {
  origin: string;
  destination: string;
  size: number;
}

export interface RunRequest {
  nodes: NodeDTO[];
  edges: EdgeDTO[];
  demands: DemandDTO[];
  maxRounds?: number;
}

export interface RouteGroupDTO {
  edgeIds: string[];
  count: number;
  travelTime: number;
}

export interface DemandResultDTO {
  origin: string;
  destination: string;
  size: number;
  routeGroups: RouteGroupDTO[];
  totalTravelTime: number;
  averageTravelTime: number;
}

export interface RunResponse {
  converged: boolean;
  rounds: number;
  totalTravelTime: number;
  averageTravelTime: number;
  perDemand: DemandResultDTO[];
}

export interface ApiError {
  error: string;
}

/**
 * Maps the client's edit model to the wire request shape. Every road
 * expands into two directed edges (research.md decision #1/#2) — the
 * server itself has no notion of "road," only the arbitrary edge list it
 * already accepted before this feature (FR-007: no server change).
 */
export function toRunRequest(state: NetworkState, maxRounds?: number): RunRequest {
  return {
    nodes: state.nodes.map((n) => ({ id: n.id, type: n.type })),
    edges: expandRoadsToDirectedEdges(state.roads).map((e) => ({
      id: e.id,
      from: e.from,
      to: e.to,
      length: 0,
      capacity: 0,
      travelTime: e.travelTime,
    })),
    demands: state.demands.map((d) => ({ origin: d.origin, destination: d.destination, size: d.size })),
    ...(maxRounds ? { maxRounds } : {}),
  };
}

const API_BASE = import.meta.env.VITE_API_BASE ?? "http://localhost:8090";

/**
 * Calls POST /api/run. Throws with the server's own error message on a
 * non-2xx response (a 400 validation failure or a 422 RunDemands
 * failure, e.g. naming an unroutable demand — FR-009).
 */
export async function runSimulation(request: RunRequest): Promise<RunResponse> {
  const res = await fetch(`${API_BASE}/api/run`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(request),
  });

  const body = (await res.json()) as RunResponse | ApiError;

  if (!res.ok) {
    throw new Error("error" in body ? body.error : `request failed with status ${res.status}`);
  }

  return body as RunResponse;
}
