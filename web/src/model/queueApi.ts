// The POST /api/queue-run request/response shapes, mirroring
// cmd/server/queuedto.go exactly (see specs/009-signal-frontend/
// contracts/queue-api-contract.md). Sibling to api.ts's POST /api/run —
// a deliberately separate contract (FR-006, FR-007), not a variant of it.
import type { NodeDTO } from "./api";
import type { QueueNetworkState, TravelTimeSpec } from "./queueNetwork";

export interface EdgeDTO {
  id: string;
  from: string;
  to: string;
  length: number;
  capacity: number;
  travelTime: TravelTimeSpec;
}

export interface SignalDTO {
  id: string;
  edgeId: string;
  greenDuration: number;
  redDuration: number;
  offset?: number;
}

export interface TimeDemandDTO {
  origin: string;
  destination: string;
  count: number;
  arrivalInterval: number;
}

export interface QueueRunRequest {
  nodes: NodeDTO[];
  edges: EdgeDTO[];
  signals: SignalDTO[];
  demands: TimeDemandDTO[];
  duration: number;
  tick: number;
}

export interface QueueSampleDTO {
  time: number;
  signalId: string;
  length: number;
  /** Whether this signal was green at `time` (feature 010, FR-002). */
  green: boolean;
}

export interface AgentReportDTO {
  /** Spawn-order identifier, matching PositionSampleDTO.agentId (feature 010). */
  id: number;
  demandIndex: number;
  spawnTime: number;
  travelTime: number;
  waitTime: number;
  arrived: boolean;
}

/** One agent's location at one simulated tick (feature 010, FR-001). */
export interface PositionSampleDTO {
  time: number;
  agentId: number;
  edgeId: string;
  /** 0..1 along edgeId; always 0 while queued. */
  progress: number;
  queued: boolean;
}

export interface QueueRunResponse {
  queueSamples: QueueSampleDTO[];
  agents: AgentReportDTO[];
  positions: PositionSampleDTO[];
}

export interface QueueApiError {
  error: string;
}

/** Maps a QueueNetworkState (+ run duration/tick) to the wire request shape. */
export function toQueueRunRequest(state: QueueNetworkState, duration: number, tick: number): QueueRunRequest {
  const signals: SignalDTO[] = state.roads
    .filter((r) => r.signal !== null)
    .map((r) => ({
      id: `${r.id}-signal`,
      edgeId: r.id,
      greenDuration: r.signal!.greenDuration,
      redDuration: r.signal!.redDuration,
    }));

  return {
    nodes: state.nodes.map((n) => ({ id: n.id, type: n.type })),
    edges: state.roads.map((r) => ({
      id: r.id,
      from: r.from,
      to: r.to,
      length: 0,
      capacity: 0,
      travelTime: r.travelTime,
    })),
    signals,
    demands: state.demands.map((d) => ({
      origin: d.origin,
      destination: d.destination,
      count: d.count,
      arrivalInterval: d.arrivalInterval,
    })),
    duration,
    tick,
  };
}

const API_BASE = import.meta.env.VITE_API_BASE ?? "http://localhost:8090";

/**
 * Calls POST /api/queue-run. Throws with the server's own error message
 * on a non-2xx response (a 400 validation failure or a 422 queuesim.Run
 * failure — FR-008).
 */
export async function runQueueSimulation(request: QueueRunRequest): Promise<QueueRunResponse> {
  const res = await fetch(`${API_BASE}/api/queue-run`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(request),
  });

  const body = (await res.json()) as QueueRunResponse | QueueApiError;

  if (!res.ok) {
    throw new Error("error" in body ? body.error : `request failed with status ${res.status}`);
  }

  return body as QueueRunResponse;
}
