import { useMemo, useState } from "react";
import { usePlaybackClock } from "./animation/usePlaybackClock";
import { NetworkCanvas, type AgentMarker } from "./components/NetworkCanvas";
import { PlaybackControls } from "./components/PlaybackControls";
import { QueueResultsPanel } from "./components/QueueResultsPanel";
import { SignalToolbar, type SignalMode } from "./components/SignalToolbar";
import { humanizeErrorMessage, type ClientRoad } from "./model/network";
import { runQueueSimulation, toQueueRunRequest, type QueueRunResponse } from "./model/queueApi";
import { groupPositionsByAgent, positionAtSimTime, signalGreenAtSimTime } from "./model/queuePlayback";
import {
  addDirectedRoad,
  addTimeDemand,
  emptyQueueNetwork,
  isQueueNetworkError,
  removeDirectedRoad,
  removeNodeFromQueueNetwork,
  removeTimeDemand,
  addNode as addQueueNode,
  type NodeType,
  type QueueNetworkState,
  type TravelTimeSpec,
} from "./model/queueNetwork";

const DEFAULT_DURATION = 60;
const DEFAULT_TICK = 0.25;

// Simulated-seconds-per-real-second for playback — a visual pacing
// choice (research.md decision #4), not derived from anything the
// engine reports. 1 means roughly real-time (a 60-70s run takes about a
// minute to watch) — slow enough to actually follow an agent queuing
// and resuming, which a faster pace (this was 4x) made hard to see.
const PLAYBACK_SPEED = 1;

// Defensive cap on rendered agent icons, matching the equilibrium mode's
// own MAX_ICONS_PER_ROUTE precedent (App.tsx) — this tool targets tens
// to a few hundred hand-drawn agents (spec.md Assumptions), not
// unbounded populations.
const MAX_ANIMATED_AGENTS = 300;

const SIGNAL_GREEN_COLOR = "#2f9e44";
const SIGNAL_RED_COLOR = "#e03131";

type Selection = { kind: "node" | "road"; id: string } | null;

/**
 * The Signals-mode app: a fully separate slice from App.tsx's equilibrium
 * mode (FR-006, FR-007, research.md decision #2) — its own
 * QueueNetworkState, its own run flow against POST /api/queue-run, and
 * its own result view. Alongside the existing QueueResultsPanel
 * (chart + arrival/wait summary, feature 009, kept per FR-008), a
 * playback now shows agents actually moving and each signal-controlled
 * road's live green/red state (feature 010) — reusing NetworkCanvas for
 * rendering only, mapping DirectedRoad.from/.to to the {a, b} shape it
 * expects (research.md decision #3).
 */
export function SignalApp() {
  const [network, setNetwork] = useState<QueueNetworkState>(emptyQueueNetwork);
  const [result, setResult] = useState<QueueRunResponse | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [running, setRunning] = useState(false);
  const [playing, setPlaying] = useState(false);

  const [mode, setMode] = useState<SignalMode>("select");
  const [nodeTypeToPlace, setNodeTypeToPlace] = useState<NodeType>("house");
  const [travelTime, setTravelTime] = useState<TravelTimeSpec>({ type: "constant", value: 5 });
  const [hasSignal, setHasSignal] = useState(false);
  const [greenDuration, setGreenDuration] = useState(5);
  const [redDuration, setRedDuration] = useState(5);
  const [connectingFrom, setConnectingFrom] = useState<string | null>(null);
  const [selection, setSelection] = useState<Selection>(null);
  const [editError, setEditError] = useState<string | null>(null);

  const clock = usePlaybackClock(DEFAULT_DURATION, PLAYBACK_SPEED, playing, result);

  function applyEdit(update: (net: QueueNetworkState) => QueueNetworkState) {
    setResult(null);
    setPlaying(false);
    setError(null);
    setEditError(null);
    setNetwork(update);
  }

  function handleCanvasClick(x: number, y: number) {
    if (mode === "place") {
      applyEdit((net) => addQueueNode(net, nodeTypeToPlace, x, y));
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
      const signal = hasSignal ? { greenDuration, redDuration } : null;
      const outcome = addDirectedRoad(network, from, id, travelTime, signal);
      if (isQueueNetworkError(outcome)) {
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
      applyEdit((net) => removeNodeFromQueueNetwork(net, selection.id));
    } else {
      applyEdit((net) => removeDirectedRoad(net, selection.id));
    }
    setSelection(null);
  }

  async function handleRun() {
    setError(null);
    setRunning(true);
    setPlaying(false);
    try {
      const request = toQueueRunRequest(network, DEFAULT_DURATION, DEFAULT_TICK);
      const response = await runQueueSimulation(request);
      setResult(response);
      setPlaying(true);
    } catch (e) {
      const raw = e instanceof Error ? e.message : String(e);
      setError(humanizeErrorMessage(raw, network.nodes));
    } finally {
      setRunning(false);
    }
  }

  // NetworkCanvas only knows the equilibrium mode's bidirectional
  // ClientRoad shape (a/b) — map each DirectedRoad to it purely for
  // rendering (research.md decision #3); the direction and signal live
  // only in `network.roads`, never derived back out of this.
  const canvasRoads: ClientRoad[] = network.roads.map((r) => ({ id: r.id, a: r.from, b: r.to, travelTime: r.travelTime }));

  const nodesById = useMemo(() => new Map(network.nodes.map((n) => [n.id, { x: n.x, y: n.y }])), [network.nodes]);
  const edgesById = useMemo(() => new Map(network.roads.map((r) => [r.id, { from: r.from, to: r.to }])), [network.roads]);
  const positionsByAgent = useMemo(() => (result ? groupPositionsByAgent(result.positions) : new Map()), [result]);

  // One AgentMarker per agent, placed at the current playback time
  // (feature 010, User Story 1) — recomputed every frame as `clock.time`
  // advances, purely from data the run already returned.
  const agents: AgentMarker[] = useMemo(() => {
    const markers: AgentMarker[] = [];
    for (const [agentId, samples] of positionsByAgent) {
      if (markers.length >= MAX_ANIMATED_AGENTS) break;
      const point = positionAtSimTime(samples, edgesById, nodesById, clock.time);
      if (point) markers.push({ id: `agent-${agentId}`, x: point.x, y: point.y });
    }
    return markers;
  }, [positionsByAgent, edgesById, nodesById, clock.time]);

  // Every signal-controlled road's live color at the current playback
  // time (feature 010, User Story 2) — a road with no signal never
  // appears in this map, so NetworkCanvas renders it unchanged.
  const roadColors: Record<string, string> = useMemo(() => {
    if (!result) return {};
    const colors: Record<string, string> = {};
    for (const road of network.roads) {
      if (!road.signal) continue;
      const signalId = `${road.id}-signal`;
      const green = signalGreenAtSimTime(result.queueSamples, signalId, clock.time);
      colors[road.id] = green ? SIGNAL_GREEN_COLOR : SIGNAL_RED_COLOR;
    }
    return colors;
  }, [result, network.roads, clock.time]);

  const runDisabled = network.demands.length === 0 || running;

  return (
    <div style={{ maxWidth: "1200px", margin: "0 auto", padding: "var(--space-5)" }}>
      <header style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-start", gap: "var(--space-4)", marginBottom: "var(--space-4)" }}>
        <div>
          <h1>Braess — Signals</h1>
          <p className="text-muted">
            A separate, discrete-time mode: roads are one-way, can carry a traffic signal, and agents arrive over
            time instead of all at once — showing real, sustained queuing rather than a converged equilibrium.
          </p>
        </div>
        <div style={{ display: "flex", flexDirection: "column", alignItems: "flex-end", gap: "var(--space-1)" }}>
          <button type="button" className="primary" onClick={() => void handleRun()} disabled={runDisabled}>
            {running ? "Running…" : "Run"}
          </button>
          {network.demands.length === 0 && <small style={{ maxWidth: "260px", textAlign: "right" }}>Add at least one time-based demand before running.</small>}
          {editError && <small style={{ color: "var(--color-danger)" }}>⚠ {editError}</small>}
          {error && <small style={{ color: "var(--color-danger)" }}>⚠ {error}</small>}
        </div>
      </header>

      {result && (
        <div style={{ marginBottom: "var(--space-4)" }}>
          <QueueResultsPanel samples={result.queueSamples} agents={result.agents} roads={network.roads} nodes={network.nodes} />
        </div>
      )}

      <div style={{ display: "flex", gap: "var(--space-4)", alignItems: "flex-start" }}>
        <div style={{ width: "300px", flexShrink: 0 }}>
          <SignalToolbar
            mode={mode}
            onModeChange={(m) => {
              setMode(m);
              setConnectingFrom(null);
              setSelection(null);
              setEditError(null);
            }}
            nodeTypeToPlace={nodeTypeToPlace}
            onNodeTypeChange={setNodeTypeToPlace}
            travelTime={travelTime}
            onTravelTimeChange={setTravelTime}
            hasSignal={hasSignal}
            onHasSignalChange={setHasSignal}
            greenDuration={greenDuration}
            redDuration={redDuration}
            onGreenDurationChange={setGreenDuration}
            onRedDurationChange={setRedDuration}
            connectingFrom={connectingFrom}
            hasSelection={selection !== null}
            onRemoveSelected={handleRemoveSelected}
            nodes={network.nodes}
            demands={network.demands}
            onAddDemand={(origin, destination, count, arrivalInterval) =>
              applyEdit((net) => addTimeDemand(net, origin, destination, count, arrivalInterval))
            }
            onRemoveDemand={(index) => applyEdit((net) => removeTimeDemand(net, index))}
          />
        </div>

        <div style={{ display: "flex", flexDirection: "column", gap: "var(--space-3)" }}>
          {result && (
            <PlaybackControls
              playing={playing}
              onTogglePlaying={() => setPlaying((p) => !p)}
              onRestart={() => {
                clock.restart();
                setPlaying(true);
              }}
              time={clock.time}
              duration={DEFAULT_DURATION}
            />
          )}
          <NetworkCanvas
            nodes={network.nodes}
            roads={canvasRoads}
            agents={agents}
            roadColors={roadColors}
            directed
            onCanvasClick={handleCanvasClick}
            onNodeClick={handleNodeClick}
            onRoadClick={handleRoadClick}
            selectedId={selection?.id ?? connectingFrom ?? null}
          />
        </div>
      </div>
    </div>
  );
}
