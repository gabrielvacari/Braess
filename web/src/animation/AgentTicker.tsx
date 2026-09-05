import { useEffect } from "react";
import { useAgentAnimation, type Point } from "./useAgentAnimation";

interface AgentTickerProps {
  id: string;
  path: Point[];
  durationMs: number;
  playing: boolean;
  startDelayMs?: number;
  onUpdate: (id: string, position: Point) => void;
}

/**
 * Renders nothing itself — just runs useAgentAnimation for one agent and
 * reports its position up to a parent-held map, so the parent can pass a
 * flat AgentMarker[] into NetworkCanvas regardless of how many agents are
 * animating.
 */
export function AgentTicker({ id, path, durationMs, playing, startDelayMs, onUpdate }: AgentTickerProps) {
  const position = useAgentAnimation(path, durationMs, playing, startDelayMs);

  useEffect(() => {
    onUpdate(id, position);
  }, [id, position, onUpdate]);

  return null;
}
