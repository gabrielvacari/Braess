import type { RunResponse } from "./api";

export interface RunHistory {
  current: RunResponse | null;
  previous: RunResponse | null;
}

export const emptyRunHistory: RunHistory = { current: null, previous: null };

/**
 * Records a new run's result: the prior "current" becomes "previous", so
 * a before/after comparison stays visible together across a network edit
 * and re-run (FR-007, SC-004).
 */
export function recordRun(history: RunHistory, next: RunResponse): RunHistory {
  return { current: next, previous: history.current };
}

/** Clears history — used when loading a fresh example, not just re-running the same network. */
export function resetRunHistory(): RunHistory {
  return emptyRunHistory;
}
