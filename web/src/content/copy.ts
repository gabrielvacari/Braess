// Every explanatory string this app shows, plus the small amount of
// logic that decides which one applies right now. Kept separate from the
// components that render it (App.tsx, Toolbar.tsx, NetworkCanvas.tsx) so
// the *logic* — is every mode covered? does the hint disappear once it
// no longer applies? — can be tested independently of rendering, the
// same split model/network.ts already established for edit operations.
import type { EditMode } from "../components/Toolbar";
import type { NodeType } from "../model/network";

/** What the tool demonstrates — FR-001, shown as soon as the app loads. */
export const APP_INTRO =
  "This is a live demonstration of Braess's Paradox: adding a road to a network can make " +
  "the average trip slower for everyone, because each driver picks their own fastest route " +
  "without coordinating with anyone else. Load the example below, or draw your own network, " +
  "then press Run to see it happen.";

export interface LegendEntry {
  type: NodeType;
  /** Matches NetworkCanvas's NODE_STYLE[type].label exactly — see contracts/copy-api.md. */
  label: string;
  name: string;
  description: string;
}

/** FR-002: one entry per node type, always all three, always in this order. */
export const NODE_TYPE_LEGEND: LegendEntry[] = [
  { type: "house", label: "H", name: "House", description: "Where agents start their trip." },
  { type: "company", label: "C", name: "Company", description: "Where agents are headed." },
  { type: "intersection", label: "I", name: "Intersection", description: "A point where roads meet." },
];

/**
 * What clicking the canvas or a node will do in the given mode right
 * now (FR-003). For "connect" mode, connectingFrom distinguishes the
 * first click (pick a starting node) from the second (pick where it
 * connects to).
 */
export function modeDescription(mode: EditMode, connectingFrom: string | null): string {
  switch (mode) {
    case "place":
      return "Click anywhere on the canvas to place a node of the selected type.";
    case "connect":
      return connectingFrom
        ? `Click a second node to connect it to "${connectingFrom}" with a new road.`
        : "Click a node to start drawing a road from it.";
    case "select":
      return "Click a node or road to select it, then remove it with the button below.";
    default:
      return "";
  }
}

/** FR-004: a first-action suggestion when the canvas is empty, or null once it isn't. */
export function emptyCanvasHint(nodeCount: number): string | null {
  if (nodeCount > 0) return null;
  return 'Nothing here yet — click "Load example" above, or switch to "Place node" and click the canvas to start drawing.';
}

/** FR-009: what's missing when Run is pressed with no demand declared, or null once one exists. */
export function missingDemandHint(demandCount: number): string | null {
  if (demandCount > 0) return null;
  return "Add at least one demand (an origin, a destination, and how many agents) before running — otherwise there's no traffic to simulate.";
}

/**
 * Why the "Add demand" button can't be used yet when there aren't at
 * least two nodes to pick as origin/destination, or null once there are.
 * Without this, the button is disabled with no visible explanation —
 * exactly the kind of unexplained control this feature exists to fix.
 */
export function noNodesForDemandHint(nodeCount: number): string | null {
  if (nodeCount >= 2) return null;
  return "Place at least two nodes on the canvas first — a demand needs an origin and a destination to choose from.";
}
