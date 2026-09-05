import type { AgentReportDTO, QueueSampleDTO } from "../model/queueApi";
import type { DirectedRoad } from "../model/queueNetwork";
import { nodeDisplayLabel, type ClientNode } from "../model/network";
import { QueueChart } from "./QueueChart";

interface QueueResultsPanelProps {
  samples: QueueSampleDTO[];
  agents: AgentReportDTO[];
  roads: DirectedRoad[];
  nodes: ClientNode[];
}

/**
 * The signal-controlled run's result: one queue-length-over-time chart
 * per signal-controlled road, plus an arrived/still-waiting summary and
 * average wait time (FR-005, FR-007) — clearly labeled as a
 * signal-controlled result so it is never confused with the equilibrium
 * mode's ResultsPanel (FR-006).
 */
export function QueueResultsPanel({ samples, agents, roads, nodes }: QueueResultsPanelProps) {
  if (samples.length === 0 && agents.length === 0) return null;

  const samplesBySignal = new Map<string, QueueSampleDTO[]>();
  for (const s of samples) {
    const list = samplesBySignal.get(s.signalId) ?? [];
    list.push(s);
    samplesBySignal.set(s.signalId, list);
  }

  const arrived = agents.filter((a) => a.arrived);
  const stillWaiting = agents.filter((a) => !a.arrived);
  const averageWait = arrived.length > 0 ? arrived.reduce((sum, a) => sum + a.waitTime, 0) / arrived.length : 0;

  return (
    <div>
      <h2>Signal-controlled result</h2>
      <p className="text-muted">
        {arrived.length} of {agents.length} agents arrived before the run ended
        {stillWaiting.length > 0 && ` (${stillWaiting.length} still waiting)`}. Average wait time for agents who
        arrived: {averageWait.toFixed(1)}s.
      </p>
      <div style={{ display: "flex", gap: "var(--space-4)", flexWrap: "wrap" }}>
        {roads
          .filter((r) => r.signal !== null)
          .map((road) => {
            const signalId = `${road.id}-signal`;
            const roadSamples = samplesBySignal.get(signalId) ?? [];
            const title = `${nodeDisplayLabel(nodes, road.from)} → ${nodeDisplayLabel(nodes, road.to)}`;
            return <QueueChart key={road.id} title={title} samples={roadSamples} />;
          })}
      </div>
    </div>
  );
}
