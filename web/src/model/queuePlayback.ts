// Pure playback math for Signals mode (feature 010): every function
// here reads back data queuesim.Run already computed — none of it
// decides a route, a queue order, or a signal phase itself (constitution
// Principle III boundary). Kept separate from SignalApp.tsx's rendering
// so it can be tested without mounting anything, the same split every
// other model/*.ts file in this project already established.
import type { PositionSampleDTO, QueueSampleDTO } from "./queueApi";

export interface CanvasPoint {
  x: number;
  y: number;
}

/**
 * Groups a run's flat position samples by which agent they belong to,
 * preserving each group's recorded (time) order — the server does not
 * group or sort them further (contracts/playback-contract.md).
 */
export function groupPositionsByAgent(positions: PositionSampleDTO[]): Map<number, PositionSampleDTO[]> {
  const byAgent = new Map<number, PositionSampleDTO[]>();
  for (const p of positions) {
    const group = byAgent.get(p.agentId);
    if (group) {
      group.push(p);
    } else {
      byAgent.set(p.agentId, [p]);
    }
  }
  return byAgent;
}

/**
 * Places one agent on the canvas at a given simulated time, from its own
 * time-ordered position samples: finds the last sample at or before
 * `time` and either holds it at its edge's start (queued) or
 * interpolates from -> to by `progress` (moving). Returns null before
 * the agent's first sample (it hasn't spawned yet at this `time`) — a
 * caller should render nothing for that agent rather than guessing a
 * position (research.md decision #3).
 */
export function positionAtSimTime(
  samples: PositionSampleDTO[],
  edgesById: Map<string, { from: string; to: string }>,
  nodesById: Map<string, CanvasPoint>,
  time: number,
): CanvasPoint | null {
  if (samples.length === 0 || time < samples[0].time) {
    return null;
  }

  // samples are already time-ordered (spawn to arrival/end of run) —
  // find the last one at or before `time`.
  let current = samples[0];
  for (const s of samples) {
    if (s.time > time) break;
    current = s;
  }

  const edge = edgesById.get(current.edgeId);
  if (!edge) return null;
  const from = nodesById.get(edge.from);
  const to = nodesById.get(edge.to);
  if (!from || !to) return null;

  if (current.queued) {
    return { x: from.x, y: from.y };
  }
  return {
    x: from.x + (to.x - from.x) * current.progress,
    y: from.y + (to.y - from.y) * current.progress,
  };
}

/**
 * Whether a signal-controlled road's signal is green at a given
 * simulated time, from its own time-ordered queue samples: the nearest
 * recorded sample at or before `time`, falling back to the earliest
 * sample when `time` precedes all of them (a run always samples from
 * t=0, so this only matters for a `time` argument slightly before the
 * first tick).
 */
export function signalGreenAtSimTime(samples: QueueSampleDTO[], signalId: string, time: number): boolean {
  const forSignal = samples.filter((s) => s.signalId === signalId);
  if (forSignal.length === 0) return false;

  let current = forSignal[0];
  for (const s of forSignal) {
    if (s.time > time) break;
    current = s;
  }
  return current.green;
}
