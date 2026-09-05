import { useEffect, useState } from "react";

export interface Point {
  x: number;
  y: number;
}

/**
 * Computes an agent's position along path (an ordered list of node
 * positions, one per node visited, so length = edge count + 1) at a
 * given elapsed time out of totalTime.
 *
 * The server's RouteGroupDTO reports only a route's total travel time,
 * not a per-edge breakdown (cmd/server never re-derives that from the
 * engine's internal final volumes, which would mean changing graph/
 * agent/simulation — out of scope for this feature). So each edge is
 * given an equal share of totalTime; FR-002 only requires visible,
 * monotonic progress from origin to destination, not per-edge timing
 * accuracy.
 *
 * Returns path[0] at elapsed <= 0, and the last point at elapsed >=
 * totalTime (and for any totalTime <= 0, treated as already arrived).
 */
export function positionAtTime(path: Point[], totalTime: number, elapsed: number): Point {
  if (path.length === 0) {
    return { x: 0, y: 0 };
  }
  if (path.length === 1 || totalTime <= 0) {
    return path[path.length - 1];
  }

  const segments = path.length - 1;
  const clampedElapsed = Math.max(0, Math.min(elapsed, totalTime));
  const fraction = clampedElapsed / totalTime; // 0..1 across the whole route
  const segmentPosition = fraction * segments; // which segment + how far into it
  const segmentIndex = Math.min(Math.floor(segmentPosition), segments - 1);
  const localFraction = segmentPosition - segmentIndex;

  const a = path[segmentIndex];
  const b = path[segmentIndex + 1];
  return {
    x: a.x + (b.x - a.x) * localFraction,
    y: a.y + (b.y - a.y) * localFraction,
  };
}

/**
 * Animates a single agent along path over durationMs of real (wall-clock)
 * time, independent of whatever abstract units the simulation's own
 * travel-time numbers are in — durationMs is purely a visual pacing
 * choice, not a re-derivation of simulation time.
 */
export function useAgentAnimation(
  path: Point[],
  durationMs: number,
  playing: boolean,
  startDelayMs = 0,
): Point {
  const [position, setPosition] = useState<Point>(path[0] ?? { x: 0, y: 0 });

  useEffect(() => {
    if (!playing || path.length === 0) {
      return;
    }

    let raf: number;
    const start = performance.now();

    const tick = (now: number) => {
      const elapsed = now - start - startDelayMs;
      setPosition(positionAtTime(path, durationMs, elapsed));
      if (elapsed < durationMs) {
        raf = requestAnimationFrame(tick);
      }
    };

    raf = requestAnimationFrame(tick);
    return () => cancelAnimationFrame(raf);
    // path is compared by reference on purpose: callers should pass a
    // stable array (e.g. useMemo) so a route that hasn't changed doesn't
    // restart the animation every render.
  }, [path, durationMs, playing, startDelayMs]);

  return position;
}
