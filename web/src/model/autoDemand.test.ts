import { describe, expect, it } from "vitest";
import type { ClientNode } from "./network";
import { assignHousesToCompanies, AUTO_DEMAND_ARRIVAL_INTERVAL, AUTO_DEMAND_COUNT, autoTimeDemands } from "./autoDemand";

function node(id: string, type: ClientNode["type"]): ClientNode {
  return { id, type, x: 0, y: 0 };
}

describe("assignHousesToCompanies", () => {
  it("returns no assignments when there are no companies", () => {
    const nodes = [node("h1", "house"), node("h2", "house")];
    expect(assignHousesToCompanies(nodes)).toEqual([]);
  });

  it("assigns every house to the single company when there is only one", () => {
    const nodes = [node("h1", "house"), node("h2", "house"), node("c1", "company")];
    const assignments = assignHousesToCompanies(nodes);
    expect(assignments).toEqual([
      { houseId: "h1", companyId: "c1" },
      { houseId: "h2", companyId: "c1" },
    ]);
  });

  it("uses at least 2 distinct companies when there are 2+ houses and 2+ companies", () => {
    const nodes = [node("h1", "house"), node("h2", "house"), node("h3", "house"), node("c1", "company"), node("c2", "company")];
    const assignments = assignHousesToCompanies(nodes);
    const distinctCompanies = new Set(assignments.map((a) => a.companyId));
    expect(distinctCompanies.size).toBeGreaterThanOrEqual(2);
  });

  it("round-robins deterministically over companies in node order", () => {
    const nodes = [
      node("h1", "house"),
      node("h2", "house"),
      node("h3", "house"),
      node("c1", "company"),
      node("c2", "company"),
    ];
    expect(assignHousesToCompanies(nodes)).toEqual([
      { houseId: "h1", companyId: "c1" },
      { houseId: "h2", companyId: "c2" },
      { houseId: "h3", companyId: "c1" },
    ]);
  });

  it("is deterministic across repeated calls with the same nodes (FR-006)", () => {
    const nodes = [node("h1", "house"), node("h2", "house"), node("c1", "company"), node("c2", "company")];
    expect(assignHousesToCompanies(nodes)).toEqual(assignHousesToCompanies(nodes));
  });

  it("ignores intersection nodes entirely", () => {
    const nodes = [node("i1", "intersection"), node("c1", "company")];
    expect(assignHousesToCompanies(nodes)).toEqual([]);
  });
});

describe("autoTimeDemands", () => {
  it("maps each assignment to the fixed-rate wire shape", () => {
    const nodes = [node("h1", "house"), node("c1", "company")];
    expect(autoTimeDemands(nodes)).toEqual([
      { origin: "h1", destination: "c1", count: AUTO_DEMAND_COUNT, arrivalInterval: AUTO_DEMAND_ARRIVAL_INTERVAL },
    ]);
  });

  it("returns an empty array when no house can be assigned a company", () => {
    expect(autoTimeDemands([node("h1", "house")])).toEqual([]);
  });
});
