// Pure geometry for the top-down road visual (feature 011) — kept
// separate from NetworkCanvas.tsx's rendering so it can be tested
// without mounting anything, the same split queueChart.ts and
// queuePlayback.ts already established.

export interface RoadMark {
  x: number;
  y: number;
  /** Unit direction vector (a -> b) at this mark's position. */
  dx: number;
  dy: number;
}

/**
 * Places one directional mark every `spacing` pixels along the segment
 * from `a` to `b`, strictly between the two nodes' edges (so no mark
 * overlaps a node) — used to paint pavement-style direction arrows
 * along a one-way road instead of a single end arrowhead (FR-004).
 * Returns an empty array when the road is too short to fit even one
 * mark between the nodes' edges (spec Edge Cases).
 */
export function roadChevronPositions(
  a: { x: number; y: number },
  b: { x: number; y: number },
  nodeRadius: number,
  spacing: number,
): RoadMark[] {
  const dx = b.x - a.x;
  const dy = b.y - a.y;
  const dist = Math.hypot(dx, dy) || 1;
  const ux = dx / dist;
  const uy = dy / dist;

  const start = nodeRadius + spacing / 2;
  const end = dist - nodeRadius;
  if (end <= start) return [];

  const marks: RoadMark[] = [];
  for (let d = start; d <= end; d += spacing) {
    marks.push({ x: a.x + ux * d, y: a.y + uy * d, dx: ux, dy: uy });
  }
  return marks;
}
