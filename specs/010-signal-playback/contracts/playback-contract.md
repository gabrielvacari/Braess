# Contract: `POST /api/queue-run` response additions

This feature does not add a new endpoint — it extends the existing
`POST /api/queue-run` (feature 009) response with additive fields. The
request shape is unchanged.

## Response body (additions highlighted)

```json
{
  "queueSamples": [
    { "time": 0.25, "signalId": "road-a-signal", "length": 3, "green": false }
  ],
  "agents": [
    { "id": 0, "demandIndex": 0, "spawnTime": 0, "travelTime": 4.5, "waitTime": 1.25, "arrived": true }
  ],
  "positions": [
    { "time": 0.25, "agentId": 0, "edgeId": "road-a", "progress": 0, "queued": true },
    { "time": 0.5,  "agentId": 0, "edgeId": "road-a", "progress": 0, "queued": true },
    { "time": 0.75, "agentId": 0, "edgeId": "road-a", "progress": 0.25, "queued": false },
    { "time": 1.0,  "agentId": 0, "edgeId": "road-a", "progress": 1.0, "queued": false }
  ]
}
```

- `queueSamples[].green` (new): whether that signal was green at that
  tick.
- `agents[].id` (new): stable identifier, matching `positions[].agentId`
  for that same agent.
- `positions` (new): one entry per in-flight agent per tick, in
  recording order — not grouped or sorted further by the server; the
  client groups by `agentId` (see `model/queuePlayback.ts`).

## Client model operations (`web/src/model/queuePlayback.ts`)

```ts
function groupPositionsByAgent(positions: PositionSampleDTO[]): Map<number, PositionSampleDTO[]>;

function positionAtSimTime(
  samples: PositionSampleDTO[],
  edgesById: Map<string, { from: string; to: string }>,
  nodesById: Map<string, { x: number; y: number }>,
  time: number,
): { x: number; y: number } | null;

function signalGreenAtSimTime(
  samples: QueueSampleDTO[],
  signalId: string,
  time: number,
): boolean;
```

## Playback clock (`web/src/animation/usePlaybackClock.ts`)

```ts
function usePlaybackClock(duration: number, speed: number, playing: boolean): {
  time: number;
  finished: boolean;
};
```

`speed` is simulated-seconds-per-real-second (a visual pacing choice,
not a claim about wall-clock accuracy — mirrors how `App.tsx`'s
`durationMs` is a pacing choice, not derived from the engine's abstract
travel-time units).

## `NetworkCanvas` prop addition

```ts
interface NetworkCanvasProps {
  // ...existing fields unchanged...
  roadColors?: Record<string, string>; // roadId -> stroke color override
}
```
