import { useCallback, useMemo, useState } from "react";
import { AgentTicker } from "./animation/AgentTicker";
import type { Point } from "./animation/useAgentAnimation";
import { IntroBanner } from "./components/IntroBanner";
import { Legend } from "./components/Legend";
import { NetworkCanvas, type AgentMarker } from "./components/NetworkCanvas";
import { ResultsPanel } from "./components/ResultsPanel";
import { Toolbar, type EditMode } from "./components/Toolbar";
import { missingDemandHint } from "./content/copy";
import { runSimulation, toRunRequest, type RunResponse } from "./model/api";
import {
  addDemand,
  addNode,
  addRoad,
  emptyNetwork,
  expandRoadsToDirectedEdges,
  humanizeErrorMessage,
  isNetworkError,
  removeDemand,
  removeNode,
  removeRoad,
  type ClientNode,
  type ClientRoad,
  type NetworkState,
  type NodeType,
  type TravelTimeSpec,
} from "./model/network";
import { emptyRunHistory, recordRun, resetRunHistory } from "./model/runHistory";

// A small preset network — the same shape the Phase 1-4 CLI (cmd/graphcli)
// already validated in text: house-1 -> intersection-1 -> {company-1,
// company-2}, plus a slower direct house-1 -> company-2 road. Reusing it
// here is the whole point of this feature: nothing new is being decided,
// only made visible. Every road is bidirectional (feature 007), so this
// still finds the same shortest route the CLI's original demo did.
function exampleNetwork(): NetworkState {
  const nodes: ClientNode[] = [
    { id: "house-1", type: "house", x: 100, y: 300 },
    { id: "intersection-1", type: "intersection", x: 400, y: 150 },
    { id: "company-1", type: "company", x: 700, y: 80 },
    { id: "company-2", type: "company", x: 700, y: 320 },
  ];
  const roads: ClientRoad[] = [
    { id: "road-1", a: "house-1", b: "intersection-1", travelTime: { type: "linear", freeFlow: 2, slope: 0.1 } },
    { id: "road-2", a: "intersection-1", b: "company-1", travelTime: { type: "constant", value: 6 } },
    { id: "road-3", a: "intersection-1", b: "company-2", travelTime: { type: "linear", freeFlow: 3, slope: 0.2 } },
    { id: "road-4", a: "house-1", b: "company-2", travelTime: { type: "constant", value: 9 } },
  ];
  return { nodes, roads, demands: [{ origin: "house-1", destination: "company-2", size: 3 }] };
}

// A route group's `count` can be in the thousands (e.g. the classic
// Braess demo's 4000 agents) — rendering one icon per agent would be
// both unreadable and wasteful, so the response is already grouped by
// distinct route (feature 005, research.md decision #4) and the browser
// only animates a bounded, representative sample per group. 200 is high
// enough to show every agent individually for any realistically-sized
// hand-drawn demand (tens to a couple hundred), while still capping the
// pathological large-population case.
const MAX_ICONS_PER_ROUTE = 200;
const MS_PER_EDGE = 1200;

interface AgentSpec {
  id: string;
  path: Point[];
  durationMs: number;
  startDelayMs: number;
}

// Resolves the server's directed-edge ids (from expandRoadsToDirectedEdges,
// echoed back via RouteGroupDTO.edgeIds) to node positions for animation.
function edgeIdsToPath(edgeIds: string[], network: NetworkState): Point[] {
  const nodeById = new Map(network.nodes.map((n) => [n.id, n]));
  const edgeById = new Map(expandRoadsToDirectedEdges(network.roads).map((e) => [e.id, e]));

  const path: Point[] = [];
  edgeIds.forEach((edgeId, i) => {
    const edge = edgeById.get(edgeId);
    if (!edge) return;
    if (i === 0) {
      const from = nodeById.get(edge.from);
      if (from) path.push({ x: from.x, y: from.y });
    }
    const to = nodeById.get(edge.to);
    if (to) path.push({ x: to.x, y: to.y });
  });
  return path;
}

function buildAgentSpecs(result: RunResponse, network: NetworkState): AgentSpec[] {
  const specs: AgentSpec[] = [];
  result.perDemand.forEach((demand, demandIndex) => {
    demand.routeGroups.forEach((group, groupIndex) => {
      const path = edgeIdsToPath(group.edgeIds, network);
      const iconCount = Math.min(group.count, MAX_ICONS_PER_ROUTE);
      for (let i = 0; i < iconCount; i++) {
        specs.push({
          id: `d${demandIndex}-g${groupIndex}-a${i}`,
          path,
          durationMs: Math.max(1500, group.edgeIds.length * MS_PER_EDGE),
          startDelayMs: i * 150,
        });
      }
    });
  });
  return specs;
}

type Selection = { kind: "node" | "road"; id: string } | null;

function App() {
  const [network, setNetwork] = useState<NetworkState>(emptyNetwork);
  const [history, setHistory] = useState(emptyRunHistory);
  const result = history.current;
  const [error, setError] = useState<string | null>(null);
  const [playing, setPlaying] = useState(false);
  const [positions, setPositions] = useState<Record<string, Point>>({});

  const [mode, setMode] = useState<EditMode>("select");
  const [nodeTypeToPlace, setNodeTypeToPlace] = useState<NodeType>("house");
  const [travelTimeToUse, setTravelTimeToUse] = useState<TravelTimeSpec>({ type: "constant", value: 5 });
  const [connectingFrom, setConnectingFrom] = useState<string | null>(null);
  const [selection, setSelection] = useState<Selection>(null);
  const [editError, setEditError] = useState<string | null>(null);

  const agentSpecs = useMemo(() => (result ? buildAgentSpecs(result, network) : []), [result, network]);

  const handleUpdate = useCallback((id: string, position: Point) => {
    setPositions((prev) => ({ ...prev, [id]: position }));
  }, []);

  const agents: AgentMarker[] = agentSpecs.map((spec) => ({
    id: spec.id,
    x: positions[spec.id]?.x ?? spec.path[0]?.x ?? 0,
    y: positions[spec.id]?.y ?? spec.path[0]?.y ?? 0,
  }));

  // Any edit stops a currently-playing animation (research.md decision
  // #6, feature 005) — the prior run's result is now stale until "Run"
  // is pressed again.
  function applyEdit(update: (net: NetworkState) => NetworkState) {
    setPlaying(false);
    setEditError(null);
    setNetwork(update);
  }

  function handleCanvasClick(x: number, y: number) {
    if (mode === "place") {
      applyEdit((net) => addNode(net, nodeTypeToPlace, x, y));
    }
  }

  function handleNodeClick(id: string) {
    if (mode === "select") {
      setSelection({ kind: "node", id });
      return;
    }
    if (mode === "connect") {
      if (!connectingFrom) {
        setConnectingFrom(id);
        return;
      }
      const from = connectingFrom;
      setConnectingFrom(null);
      // One click-pair creates one bidirectional road (feature 007,
      // FR-001) — the user never has to draw a second connection to get
      // the reverse direction.
      const outcome = addRoad(network, from, id, travelTimeToUse);
      if (isNetworkError(outcome)) {
        setEditError("A road cannot connect a node to itself.");
        return;
      }
      applyEdit(() => outcome);
    }
  }

  function handleRoadClick(id: string) {
    if (mode === "select") {
      setSelection({ kind: "road", id });
    }
  }

  function handleRemoveSelected() {
    if (!selection) return;
    if (selection.kind === "node") {
      applyEdit((net) => removeNode(net, selection.id));
    } else {
      // Removes both directions at once — a road is a single record
      // (feature 007, FR-005).
      applyEdit((net) => removeRoad(net, selection.id));
    }
    setSelection(null);
  }

  async function handleRun(net: NetworkState) {
    setError(null);
    setPlaying(false);
    setPositions({});
    try {
      const response = await runSimulation(toRunRequest(net));
      // FR-007/SC-004 (feature 005): keep the prior result around so
      // before/after stays visible together once a second run happens.
      setHistory((h) => recordRun(h, response));
      setPlaying(true);
    } catch (e) {
      // A 422 (e.g. an unroutable demand, FR-008) or a network failure —
      // surfaced clearly, near the action that triggered it. The server
      // only knows raw node ids, never this client's friendly names, so
      // translate them back before showing the message.
      const raw = e instanceof Error ? e.message : String(e);
      setError(humanizeErrorMessage(raw, net.nodes));
    }
  }

  function loadExample() {
    const example = exampleNetwork();
    setNetwork(example);
    setHistory(resetRunHistory());
    void handleRun(example);
  }

  const runHint = missingDemandHint(network.demands.length); // FR-009

  return (
    <div style={{ maxWidth: "1200px", margin: "0 auto", padding: "var(--space-5)" }}>
      <header style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-start", gap: "var(--space-4)", marginBottom: "var(--space-4)" }}>
        <div>
          <h1>Braess</h1>
          <p className="text-muted">An interactive demonstration of Braess's Paradox</p>
        </div>
        <div style={{ display: "flex", flexDirection: "column", alignItems: "flex-end", gap: "var(--space-1)" }}>
          <div style={{ display: "flex", gap: "var(--space-2)" }}>
            <button type="button" onClick={loadExample}>
              Load example
            </button>
            <button type="button" className="primary" onClick={() => void handleRun(network)} disabled={network.demands.length === 0}>
              Run
            </button>
          </div>
          {runHint && <small style={{ maxWidth: "260px", textAlign: "right" }}>{runHint}</small>}
          {editError && <small style={{ color: "var(--color-danger)" }}>⚠ {editError}</small>}
          {error && <small style={{ color: "var(--color-danger)" }}>⚠ {error}</small>}
        </div>
      </header>

      <div style={{ marginBottom: "var(--space-4)" }}>
        <IntroBanner />
      </div>

      <ResultsPanel current={history.current} previous={history.previous} />

      <div style={{ display: "flex", gap: "var(--space-4)", marginTop: "var(--space-4)", alignItems: "flex-start" }}>
        <div style={{ width: "280px", flexShrink: 0 }}>
          <Toolbar
            mode={mode}
            onModeChange={(m) => {
              setMode(m);
              setConnectingFrom(null);
              setSelection(null);
              setEditError(null);
            }}
            nodeTypeToPlace={nodeTypeToPlace}
            onNodeTypeChange={setNodeTypeToPlace}
            travelTimeToUse={travelTimeToUse}
            onTravelTimeChange={setTravelTimeToUse}
            hasSelection={selection !== null}
            onRemoveSelected={handleRemoveSelected}
            connectingFrom={connectingFrom}
            nodes={network.nodes}
            demands={network.demands}
            onAddDemand={(origin, destination, size) => applyEdit((net) => addDemand(net, origin, destination, size))}
            onRemoveDemand={(index) => applyEdit((net) => removeDemand(net, index))}
          />
        </div>

        <div style={{ display: "flex", flexDirection: "column", gap: "var(--space-4)" }}>
          <NetworkCanvas
            nodes={network.nodes}
            roads={network.roads}
            agents={agents}
            onCanvasClick={handleCanvasClick}
            onNodeClick={handleNodeClick}
            onRoadClick={handleRoadClick}
            selectedId={selection?.id ?? connectingFrom ?? null}
          />
          <Legend />
        </div>
      </div>

      {agentSpecs.map((spec) => (
        <AgentTicker
          key={spec.id}
          id={spec.id}
          path={spec.path}
          durationMs={spec.durationMs}
          startDelayMs={spec.startDelayMs}
          playing={playing}
          onUpdate={handleUpdate}
        />
      ))}
    </div>
  );
}

export default App;
