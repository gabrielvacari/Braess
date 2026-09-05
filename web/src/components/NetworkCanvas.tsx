import { Fragment } from "react";
import { Arrow, Circle, Layer, Line, RegularPolygon, Rect, Stage, Text } from "react-konva";
import { emptyCanvasHint } from "../content/copy";
import type { ClientNode, ClientRoad } from "../model/network";
import { roadChevronPositions } from "../model/roadVisuals";

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
   * Overrides a road's surface color by id (feature 010, research.md
   * decision #5) — used by Signals mode to show a signal-controlled
   * road's live green/red state. Omitted entirely, rendering is
   * unchanged from before this prop existed; selection highlighting
   * still takes precedence over it.
   */
  roadColors?: Record<string, string>;
  /**
   * Adds repeated directional marks along every road (feature 011,
   * FR-004) so a one-way road's direction is actually visible —
   * omitted (the default) for the equilibrium mode, whose roads are
   * genuinely bidirectional and would be misrepresented by a direction.
   * Signals mode passes true, since every DirectedRoad there really
   * does run one way.
   */
  directed?: boolean;
}

const NODE_RADIUS = 18;
const SELECTION_COLOR = "#2f6fed";

// Feature 011: roads read as roads, not lines — a thick strip (rounded
// at the ends, via lineCap) with a dashed center-line marking. The
// marking colors stay fixed regardless of selection/signal state
// (research.md decision #2/#3) — only the strip itself changes color,
// the same way real lane paint doesn't change color with traffic.
const ROAD_WIDTH = 10;
const ROAD_WIDTH_SELECTED = 14;
const ROAD_BASE_COLOR = "#495057";
const CENTERLINE_COLOR = "#f1f3f5";
const CENTERLINE_WIDTH = 2;
const CENTERLINE_DASH = [10, 8];
const DIRECTION_MARK_SPACING = 40;
const DIRECTION_MARK_HALF_LENGTH = 6;

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
            // The surface color carries selection/signal state
            // (FR-006/FR-007); the center line and direction marks
            // below always stay CENTERLINE_COLOR, like real lane paint
            // that doesn't change with traffic (research.md decision #2).
            const surfaceColor = isSelected ? SELECTION_COLOR : (overrideColor ?? ROAD_BASE_COLOR);
            const strokeWidth = isSelected ? ROAD_WIDTH_SELECTED : ROAD_WIDTH;
            const handleClick = (e: { cancelBubble: boolean }) => {
              e.cancelBubble = true;
              onRoadClick?.(road.id);
            };
            // FR-005: a bidirectional (non-directed) road gets no marks
            // at all — only Signals mode's one-way roads do.
            const marks = directed ? roadChevronPositions(a, b, NODE_RADIUS, DIRECTION_MARK_SPACING) : [];

            return (
              <Fragment key={road.id}>
                {/* The road surface: a thick, rounded-cap strip (FR-001,
                    FR-003) — the only piece that owns clicks/hit area,
                    so selecting/removing a road is unchanged (FR-008). */}
                <Line
                  points={[a.x, a.y, b.x, b.y]}
                  stroke={surfaceColor}
                  strokeWidth={strokeWidth}
                  lineCap="round"
                  hitStrokeWidth={20}
                  onClick={handleClick}
                />
                {/* The center-line lane marking (FR-002) — purely
                    decorative, never intercepts clicks. */}
                <Line
                  points={[a.x, a.y, b.x, b.y]}
                  stroke={CENTERLINE_COLOR}
                  strokeWidth={CENTERLINE_WIDTH}
                  dash={CENTERLINE_DASH}
                  lineCap="round"
                  listening={false}
                />
                {/* Pavement-style direction marks (FR-004) — small
                    Arrows reuse the same primitive already verified to
                    point correctly, rather than a hand-rotated polygon
                    (research.md decision #3). */}
                {marks.map((m, i) => (
                  <Arrow
                    key={`${road.id}-mark-${i}`}
                    points={[
                      m.x - m.dx * DIRECTION_MARK_HALF_LENGTH,
                      m.y - m.dy * DIRECTION_MARK_HALF_LENGTH,
                      m.x + m.dx * DIRECTION_MARK_HALF_LENGTH,
                      m.y + m.dy * DIRECTION_MARK_HALF_LENGTH,
                    ]}
                    stroke={CENTERLINE_COLOR}
                    fill={CENTERLINE_COLOR}
                    strokeWidth={2}
                    pointerLength={6}
                    pointerWidth={6}
                    listening={false}
                  />
                ))}
              </Fragment>
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
