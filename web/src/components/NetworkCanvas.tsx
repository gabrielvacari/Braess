import { Arrow, Circle, Layer, Line, RegularPolygon, Rect, Stage, Text } from "react-konva";
import { emptyCanvasHint } from "../content/copy";
import type { ClientNode, ClientRoad } from "../model/network";

/** One agent icon's current position, for the animation layer (US1/T011). */
export interface AgentMarker {
  id: string;
  x: number;
  y: number;
}

export interface NetworkCanvasProps {
  nodes: ClientNode[];
  /** Each road is bidirectional and rendered as one line (feature 007, FR-004). */
  roads: ClientRoad[];
  agents?: AgentMarker[];
  width?: number;
  height?: number;
  /** Fired with the clicked point's canvas coordinates (US2: placing a node). */
  onCanvasClick?: (x: number, y: number) => void;
  /** Fired with the clicked node's id (US2: selecting nodes to draw a road, or to remove). */
  onNodeClick?: (id: string) => void;
  /** Fired with the clicked road's id (US2: selecting a road to remove). */
  onRoadClick?: (id: string) => void;
  /** A node or road id to visually highlight (US2: the current selection). */
  selectedId?: string | null;
  /**
   * Overrides a road's stroke color by id (feature 010, research.md
   * decision #5) — used by Signals mode to show a signal-controlled
   * road's live green/red state. Omitted entirely, rendering is
   * unchanged from before this prop existed; selection highlighting
   * still takes precedence over it.
   */
  roadColors?: Record<string, string>;
  /**
   * Renders every road as an arrow (a -> b) instead of a plain line, so
   * a one-way road's direction is actually visible — omitted (the
   * default) for the equilibrium mode, whose roads are genuinely
   * bidirectional and would be misrepresented by an arrow. Signals mode
   * passes true, since every DirectedRoad there really does run one way.
   */
  directed?: boolean;
}

const NODE_RADIUS = 18;
const SELECTION_COLOR = "#2f6fed";

// FR-002/SC-001: each node type gets a distinct shape, color, AND text
// label, used for no other purpose — see content/copy.ts's
// NODE_TYPE_LEGEND, which documents these same three labels for the
// on-screen legend (FR-007: type isn't color-only).
const NODE_STYLE: Record<ClientNode["type"], { fill: string; label: string }> = {
  house: { fill: "#2f9e44", label: "H" },
  company: { fill: "#f08c00", label: "C" },
  intersection: { fill: "#868e96", label: "I" },
};

function NodeShape({ node, selected = false }: { node: ClientNode; selected?: boolean }) {
  const style = NODE_STYLE[node.type];
  const stroke = selected ? SELECTION_COLOR : undefined;
  const strokeWidth = selected ? 3 : 0;

  if (node.type === "house") {
    return (
      <Rect
        x={node.x - NODE_RADIUS}
        y={node.y - NODE_RADIUS}
        width={NODE_RADIUS * 2}
        height={NODE_RADIUS * 2}
        fill={style.fill}
        stroke={stroke}
        strokeWidth={strokeWidth}
        cornerRadius={4}
      />
    );
  }

  if (node.type === "company") {
    return (
      <Circle x={node.x} y={node.y} radius={NODE_RADIUS} fill={style.fill} stroke={stroke} strokeWidth={strokeWidth} />
    );
  }

  return (
    <RegularPolygon
      x={node.x}
      y={node.y}
      sides={4}
      radius={NODE_RADIUS}
      rotation={45}
      fill={style.fill}
      stroke={stroke}
      strokeWidth={strokeWidth}
    />
  );
}

export function NetworkCanvas({
  nodes,
  roads,
  agents = [],
  width = 800,
  height = 600,
  onCanvasClick,
  onNodeClick,
  onRoadClick,
  selectedId = null,
  roadColors,
  directed = false,
}: NetworkCanvasProps) {
  const nodeById = new Map(nodes.map((n) => [n.id, n]));
  const hint = emptyCanvasHint(nodes.length); // FR-004

  return (
    <div className="canvas-frame">
      {hint && (
        <div className="canvas-hint" role="status">
          {hint}
        </div>
      )}
      <Stage
        width={width}
        height={height}
        onClick={(e) => {
          if (e.target === e.target.getStage() && onCanvasClick) {
            const pos = e.target.getStage()!.getPointerPosition();
            if (pos) onCanvasClick(pos.x, pos.y);
          }
        }}
      >
        <Layer>
          {roads.map((road) => {
            const a = nodeById.get(road.a);
            const b = nodeById.get(road.b);
            if (!a || !b) return null;
            const isSelected = road.id === selectedId;
            const overrideColor = roadColors?.[road.id];
            const color = isSelected ? SELECTION_COLOR : (overrideColor ?? "#495057");
            const handleClick = (e: { cancelBubble: boolean }) => {
              e.cancelBubble = true;
              onRoadClick?.(road.id);
            };

            if (directed) {
              // Pull the arrow's tip back to b's edge rather than its
              // center — drawn at the exact center, the arrowhead would
              // render underneath the node shape (painted afterward) and
              // never be visible.
              const dx = b.x - a.x;
              const dy = b.y - a.y;
              const dist = Math.hypot(dx, dy) || 1;
              const tipX = b.x - (dx / dist) * NODE_RADIUS;
              const tipY = b.y - (dy / dist) * NODE_RADIUS;
              return (
                <Arrow
                  key={road.id}
                  points={[a.x, a.y, tipX, tipY]}
                  stroke={color}
                  fill={color}
                  strokeWidth={isSelected ? 4 : 2}
                  pointerLength={10}
                  pointerWidth={10}
                  hitStrokeWidth={16}
                  onClick={handleClick}
                />
              );
            }

            return (
              <Line
                key={road.id}
                points={[a.x, a.y, b.x, b.y]}
                stroke={color}
                strokeWidth={isSelected ? 4 : 2}
                hitStrokeWidth={16}
                onClick={handleClick}
              />
            );
          })}

          {nodes.map((node) => (
            <NodeShape key={node.id} node={node} selected={node.id === selectedId} />
          ))}

          {nodes.map((node) => (
            <Text
              key={`${node.id}-label`}
              x={node.x - NODE_RADIUS}
              y={node.y - 6}
              width={NODE_RADIUS * 2}
              align="center"
              text={NODE_STYLE[node.type].label}
              fill="#fff"
              fontStyle="bold"
              listening={false}
            />
          ))}

          {nodes.map((node) => (
            <Circle
              key={`${node.id}-hit`}
              x={node.x}
              y={node.y}
              radius={NODE_RADIUS}
              fill="transparent"
              onClick={(e) => {
                e.cancelBubble = true;
                onNodeClick?.(node.id);
              }}
            />
          ))}

          {agents.map((a) => (
            <Circle key={a.id} x={a.x} y={a.y} radius={6} fill="#e03131" />
          ))}
        </Layer>
      </Stage>
    </div>
  );
}
