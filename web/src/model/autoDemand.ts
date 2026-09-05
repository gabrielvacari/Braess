// Automatic house-to-company demand assignment (feature 012) — Signals
// mode no longer lets a person declare who goes where; every house
// always generates traffic, and this module decides which company it
// targets. This is a derived computation over the network's nodes, not
// an edit to state, which is why it lives apart from queueNetwork.ts's
// state CRUD (research.md decision #2).
import type { ClientNode } from "./network";
import type { TimeDemandDTO } from "./queueApi";

/** FR-002: fixed for every house, not user-configurable. */
export const AUTO_DEMAND_COUNT = 12;
/** FR-002: fixed for every house, not user-configurable. */
export const AUTO_DEMAND_ARRIVAL_INTERVAL = 0.8;

export interface HouseAssignment {
  houseId: string;
  companyId: string;
}

/**
 * Assigns each house node to a company node, round-robin over the
 * companies in the order they appear in `nodes` (research.md decision
 * #1): the i-th house (in `nodes` order) goes to `companies[i %
 * companies.length]`. Deterministic given the same `nodes` (FR-006),
 * and guarantees at least 2 distinct companies are used whenever there
 * are 2+ houses and 2+ companies (FR-004) — the first two houses always
 * land on two different companies.
 *
 * Returns an empty array when there are no company nodes at all
 * (FR-007) — there is nowhere for a house's traffic to be assigned to.
 */
export function assignHousesToCompanies(nodes: ClientNode[]): HouseAssignment[] {
  const houses = nodes.filter((n) => n.type === "house");
  const companies = nodes.filter((n) => n.type === "company");
  if (companies.length === 0) return [];

  return houses.map((house, i) => ({
    houseId: house.id,
    companyId: companies[i % companies.length].id,
  }));
}

/** Maps assignHousesToCompanies's result to the wire shape (FR-002, FR-003). */
export function autoTimeDemands(nodes: ClientNode[]): TimeDemandDTO[] {
  return assignHousesToCompanies(nodes).map((a) => ({
    origin: a.houseId,
    destination: a.companyId,
    count: AUTO_DEMAND_COUNT,
    arrivalInterval: AUTO_DEMAND_ARRIVAL_INTERVAL,
  }));
}
