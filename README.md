# Braess

A traffic simulator built to experimentally study **Braess's Paradox**: the
phenomenon where adding a new route to a road network can, instead of
relieving it, **worsen** congestion — because each agent (car) chooses its
path selfishly, minimizing only its own travel time, not the collective one.

The project simulates a street graph with houses (origins) and companies
(destinations), where decentralized agents recompute routes based on current
congestion, making it possible to observe when this individual behavior
produces a worse collective outcome.

## Who was Braess

**Dietrich Braess** is a German mathematician who, in 1968, described the
paradox that bears his name while studying equilibria in traffic networks.
His work sits at the intersection of mathematical optimization and game
theory: it shows that, in systems with decentralized individual decisions,
the equilibrium reached by selfish agents (a Nash/Wardrop equilibrium) can
be worse than a coordinated optimum — a result with applications in traffic,
electrical, and computer networks.

More context and the implementation roadmap live in [AGENTS.md](AGENTS.md).
