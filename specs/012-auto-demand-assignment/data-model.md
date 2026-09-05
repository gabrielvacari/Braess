# Phase 1 Data Model: Automatic House Demand

## Removed

- `TimeDemand` interface (`queueNetwork.ts`)
- `QueueNetworkState.demands` field
- `addTimeDemand`, `removeTimeDemand` functions

`QueueNetworkState` becomes `{ nodes: ClientNode[]; roads: DirectedRoad[] }`
— no `demands` field at all, since nothing is ever stored there anymore.

## New: `web/src/model/autoDemand.ts`

### Constants

| Name | Value | Notes |
|---|---|---|
| `AUTO_DEMAND_COUNT` | `12` | FR-002, fixed. |
| `AUTO_DEMAND_ARRIVAL_INTERVAL` | `0.8` | FR-002, fixed. |

### `HouseAssignment`

| Field | Type | Notes |
|---|---|---|
| `houseId` | `string` | A house node's id. |
| `companyId` | `string` | The company node it's assigned to (FR-003). |

### `assignHousesToCompanies`

```ts
function assignHousesToCompanies(nodes: ClientNode[]): HouseAssignment[];
```

One entry per house node, in the order houses appear in `nodes`, each
assigned `companies[i % companies.length]` from the company nodes in
the order they appear in `nodes` (research.md decision #1). Returns an
empty array when there are no company nodes (FR-007) or no house nodes.

### `autoTimeDemands`

```ts
function autoTimeDemands(nodes: ClientNode[]): TimeDemandDTO[];
```

Maps `assignHousesToCompanies(nodes)` to the existing wire shape
(`TimeDemandDTO` from `queueApi.ts`, unchanged): `{ origin: houseId,
destination: companyId, count: AUTO_DEMAND_COUNT, arrivalInterval:
AUTO_DEMAND_ARRIVAL_INTERVAL }`.

## Changed call sites

- `queueApi.ts`'s `toQueueRunRequest(state, duration, tick)`: `demands`
  field now built via `autoTimeDemands(state.nodes)` instead of mapping
  `state.demands` (research.md decision #4). Signature unchanged.
- `SignalToolbar.tsx`: no longer takes `demands`/`onAddDemand`/
  `onRemoveDemand` props; renders `assignHousesToCompanies(nodes)` as a
  read-only list instead (FR-010, research.md decision #3).
- `SignalApp.tsx`: drops all demand-related state and handlers; "Run"
  is disabled based on `autoTimeDemands(network.nodes).length === 0`
  instead of `network.demands.length === 0`.
