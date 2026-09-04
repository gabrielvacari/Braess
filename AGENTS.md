# Project: Traffic Simulator / Braess's Paradox

## Objective

Build a web-based simulation where you can draw roads (a road graph) and
observe agents (cars) moving between origin points (houses) and destination
points (companies), with the goal of experimentally investigating whether
adding more roads/routes reduces or worsens congestion — Braess's Paradox.

This is also an AI Engineering portfolio project, serving as hands-on
practice with agents, decentralized decision-making, and emergent behavior
in multi-agent systems.

## Core concept

Each car is an agent that:
- Spawns at a house (origin point) and has a fixed destination (a company).
- Chooses its route trying to minimize its own travel time (selfish/
  individual behavior, not collective).
- Recalculates or reacts to current road congestion.

Congestion emerges because roads have limited capacity: the travel time of
an edge increases as more agents use it simultaneously.

## Architecture in blocks

### 1. Graph (simulation engine)
- Nodes = intersections, houses, companies.
- Edges = roads, each with a length/capacity and a travel-time function that
  grows with traffic volume.
- This block should not depend on the UI — it should run standalone (e.g. a
  terminal script) so the logic can be validated before any visuals.

### 2. Agents (cars)
- Each agent has a defined origin (house) and destination (company).
- Route choice/recalculation logic based on current traffic.
- This is where individual selfish behavior generates (or doesn't generate)
  Braess's Paradox when many agents make the same decision.

### 3. Frontend (browser)
- Interactively draw the map/graph.
- Mark nodes as houses or companies.
- Dynamically add/remove roads to test the effect on traffic formation.
- Visualize agents moving in real time.

## Implementation order (roadmap)

1. **Pure graph in code**, no UI — basic nodes and edges.
2. **Simple agents** going from A to B with shortest-path routing (no
   congestion yet).
3. **Congestion**: travel time increasing with volume — validate via
   console/logs that Braess's Paradox appears when a road is added.
4. **Multiple origins/destinations**: distinct houses and companies, not
   just a single fixed A→B pair.
5. **Web frontend**: visualize what's already validated underneath, allow
   interactively drawing roads and marking houses/companies.

> Golden rule: don't jump to the frontend before the graph + agents +
> congestion are working and demonstrating the paradox in text/logs.

## Stack (to be defined/adjusted as development progresses)

- Simulation logic: language Go.
- Frontend: React and Konva.js
- Convetional Commits & Gitflow

## Context notes

- Author is transitioning into AI Engineering, focused on agent fundamentals
  (harness engineering, computational vs. inferential controls, feedforward
  vs. feedback) and graph engineering (nodes, edges, shared state).
- Prefers learning through dialogue and conceptual synthesis — when
  proposing code, explain the reasoning behind design decisions, not just
  hand over a finished solution.