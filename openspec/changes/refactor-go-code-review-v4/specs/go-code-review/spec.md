## MODIFIED Requirements

### Requirement: Go Code Review Plugin Architecture
The plugin SHALL use a three-tier architecture that separates deterministic checks from AI judgment.

**Tier 1 — Quantitative Tools**: Scripts that produce structured JSON with zero AI involvement.
**Tier 2 — Rule Pattern Scanning**: YAML-defined rules scanned via regex, producing structured hits.
**Tier 3 — AI Expert Agents**: Five domain-expert agents that handle only judgment-based review.

#### Scenario: Full review execution
- **WHEN** the user invokes `go-code-review`
- **THEN** the orchestrator runs `analyze-go.sh` → `scan-rules.sh` → 5 parallel agents → aggregated Chinese report

#### Scenario: Tier 1 quantitative violation detected
- **WHEN** a Go file exceeds 800 lines, a function exceeds 80 lines, or nesting exceeds 4 levels
- **THEN** `analyze-go.sh` outputs the violation in `metrics.json` and the quality-agent reports it as P1

#### Scenario: Tier 2 pattern rule matched
- **WHEN** `scan-rules.sh` finds a match (e.g., `fmt.Errorf(` in source)
- **THEN** the match is recorded in `rule-hits.json` with rule-id, severity, file, line, and message
- **AND** the responsible agent confirms the hit and provides fix suggestion

#### Scenario: Agent handles judgment-only issue
- **WHEN** an issue requires semantic understanding (e.g., N+1 query pattern, UNIX philosophy violation)
- **THEN** the relevant expert agent reports it with context not capturable by regex

## ADDED Requirements

### Requirement: Structured YAML Rules
The plugin SHALL store all deterministic rules in YAML files under `rules/` with a standard schema.

Each rule MUST include: `id`, `name`, `severity` (P0/P1/P2), `pattern.match` (regex), `message`, `examples.bad`, `examples.good`.

Rule ID prefixes: `SAFE-` (safety), `DATA-` (data), `QUAL-` (quality), `OBS-` (observability).

#### Scenario: Rule file loaded
- **WHEN** `scan-rules.sh` executes
- **THEN** it reads all `rules/*.yaml` files and applies each rule's `pattern.match` to the changed Go files

#### Scenario: Rule produces structured output
- **WHEN** a rule matches
- **THEN** `rule-hits.json` contains the rule-id, file path, line number, matched text, and Chinese message

### Requirement: Design Philosophy Review Dimension
The plugin SHALL include a `design-agent` that reviews code against UNIX design philosophy and code rot root causes.

The agent SHALL evaluate all seven UNIX principles: KISS, Composition, Parsimony, Transparency, Least Surprise, Silence, and Repair.

The agent SHALL diagnose the six code rot root causes: duplicate code, no domain model, OOP abuse, lack of rigor, no upfront design, premature optimization.

#### Scenario: UNIX philosophy violation detected
- **WHEN** a function has deeply nested responsibilities (transparency violation) or inheritance where composition fits
- **THEN** `design-agent` reports the violation with the specific principle name and a Go code fix example

#### Scenario: Code rot root cause identified
- **WHEN** the same business logic appears in multiple locations
- **THEN** `design-agent` reports it as a DRY/duplicate-code violation with suggested abstraction

### Requirement: Five Domain-Expert Agents
The plugin SHALL use five agents, each covering one orthogonal expert perspective.

| Agent | Perspective | Handles |
|-------|------------|---------|
| safety | Safety & correctness | Concurrency context, defensive programming completeness |
| data | Data layer | N+1 queries, serialization strategy, type semantics |
| design | Architecture & design philosophy | UNIX principles, domain model quality, code rot causes |
| quality | Code quality | Synthesizes metrics.json + rule-hits.json, naming semantics |
| observability | Observability | Logging layer strategy, error message descriptive quality |

#### Scenario: All five agents run in parallel
- **WHEN** orchestrator launches agents
- **THEN** all five start simultaneously, each receiving metrics.json, their relevant rule-hits subset, and the code

#### Scenario: Agents do not duplicate Tier 2 findings
- **WHEN** a finding was already captured in rule-hits.json
- **THEN** the agent confirms and adds context but does NOT re-report the same raw violation

## REMOVED Requirements

### Requirement: Rule-Number-Based Agent Division
**Reason**: Dividing agents by rule numbers (1.3.\*, 2.1.\*) creates unfocused agents that span unrelated concerns. Replaced by domain-expert perspective division.
**Migration**: Old agents (gorm-review, error-safety, naming-logging, organization) are replaced by the five new domain-expert agents.
