import { describe, expect, it } from "vitest";
import type { RunResponse } from "./api";
import { emptyRunHistory, recordRun } from "./runHistory";

function fakeResult(averageTravelTime: number): RunResponse {
  return { converged: true, rounds: 1, totalTravelTime: averageTravelTime, averageTravelTime, perDemand: [] };
}

describe("recordRun", () => {
  it("has no previous result after the first run", () => {
    const history = recordRun(emptyRunHistory, fakeResult(10));
    expect(history.current).toEqual(fakeResult(10));
    expect(history.previous).toBeNull();
  });

  // FR-007/SC-004: both results stay visible together after a second run.
  it("moves the prior current result into previous on a second run", () => {
    let history = recordRun(emptyRunHistory, fakeResult(10));
    history = recordRun(history, fakeResult(15));

    expect(history.previous).toEqual(fakeResult(10));
    expect(history.current).toEqual(fakeResult(15));
  });

  it("keeps sliding the window on a third run", () => {
    let history = recordRun(emptyRunHistory, fakeResult(10));
    history = recordRun(history, fakeResult(15));
    history = recordRun(history, fakeResult(8));

    expect(history.previous).toEqual(fakeResult(15));
    expect(history.current).toEqual(fakeResult(8));
  });
});
