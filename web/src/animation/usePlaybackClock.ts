import { useEffect, useRef, useState } from "react";

export interface PlaybackClock {
  /** Current simulated time, in [0, duration]. */
  time: number;
  /** True once `time` has reached `duration` and playback stopped itself. */
  finished: boolean;
  /** Jumps back to t=0 and re-arms playback (feature 010, FR-006: replay without re-running the simulation). */
  restart: () => void;
}

/**
 * A single shared simulated-time clock for Signals-mode playback
 * (research.md decision #4) — unlike the equilibrium mode's
 * useAgentAnimation (one independent real-time tween per agent), every
 * agent's position and every road's signal color are derived from this
 * one `time` value, because they were all recorded against the same
 * simulated timeline by queuesim.Run.
 *
 * `speed` is simulated-seconds-per-real-second: a visual pacing choice
 * (how pleasant this is to watch), not a claim about wall-clock accuracy
 * — mirrors how the equilibrium mode's per-edge animation duration is
 * also a pacing choice, not derived from the engine's own time units.
 *
 * `resetKey` identifies the current run result: when it changes, this
 * hook resets to t=0 rather than continuing wherever the previous run's
 * clock happened to stop.
 */
export function usePlaybackClock(duration: number, speed: number, playing: boolean, resetKey: unknown): PlaybackClock {
  const [time, setTime] = useState(0);

  // A ref mirror of `time`, kept in sync after every render (never read
  // or written during render itself) so the animation-frame loop below
  // can read the latest value without re-arming every time it changes.
  const timeRef = useRef(time);
  useEffect(() => {
    timeRef.current = time;
  });

  const lastFrame = useRef<number | null>(null);
  const [restartTick, setRestartTick] = useState(0);

  // Reset during render when resetKey changes, rather than in an effect
  // — this is React's own recommended pattern for adjusting state in
  // response to a prop change: it applies before the browser paints,
  // with no extra render pass and no synchronous setState-in-effect.
  const [prevResetKey, setPrevResetKey] = useState(resetKey);
  if (resetKey !== prevResetKey) {
    setPrevResetKey(resetKey);
    setTime(0);
  }

  useEffect(() => {
    if (!playing || duration <= 0) {
      lastFrame.current = null;
      return;
    }

    let raf: number;
    const tick = (now: number) => {
      if (lastFrame.current === null) {
        lastFrame.current = now;
      }
      const deltaSeconds = (now - lastFrame.current) / 1000;
      lastFrame.current = now;

      const next = Math.min(duration, timeRef.current + deltaSeconds * speed);
      timeRef.current = next;
      setTime(next);

      if (next < duration) {
        raf = requestAnimationFrame(tick);
      }
    };

    raf = requestAnimationFrame(tick);
    return () => {
      cancelAnimationFrame(raf);
      lastFrame.current = null;
    };
    // restartTick has no meaning of its own — bumping it is only ever
    // used to force this effect to re-arm the animation frame loop even
    // when playing/duration/speed didn't change (e.g. restarting after
    // it already ran to completion).
  }, [playing, duration, speed, restartTick]);

  return {
    time,
    finished: time >= duration,
    restart: () => {
      setTime(0);
      setRestartTick((n) => n + 1);
    },
  };
}
