# Go Code Review Plugin Specification

## MODIFIED Requirements

### Requirement: Plugin Architecture
The Go Code Review plugin SHALL use a unified architecture with automatic agent selection and parallel execution (v3.0), replacing the previous multi-skill architecture (v2.x).

#### Scenario: Single command triggers comprehensive review
- **WHEN** user invokes "Review my Go code"
- **THEN** the system automatically selects applicable agents and runs them in parallel

#### Scenario: Automatic agent selection based on code patterns
- **WHEN** code contains GORM patterns (gorm tags, db.Model, db.Where)
- **THEN** the system activates gorm-review agent

#### Scenario: Parallel agent execution for performance
- **WHEN** multiple agents are applicable
- **THEN** all agents execute simultaneously and results are aggregated

### Requirement: Rule Coverage
The plugin SHALL enforce 137+ FUTU Go coding standards across 8 categories (expanded from 97 rules in v2.x).

#### Scenario: P0 critical rules enforcement
- **WHEN** reviewing Go code
- **THEN** the system checks all P0 rules (1.1-1.5: error handling, nil safety, GORM, concurrency, JSON)

#### Scenario: P1 code quality recommendations
- **WHEN** reviewing Go code
- **THEN** the system checks all P1 rules (2.1-2.5, 4.1-4.8: naming, logging, organization, design philosophies)

#### Scenario: P2 optimization suggestions
- **WHEN** reviewing Go code
- **THEN** the system checks all P2 rules (3.1-3.3: project structure, testing, configuration)

#### Scenario: Design philosophy validation
- **WHEN** reviewing Go code
- **THEN** the system validates against 8 design philosophies (KISS, DRY, YAGNI, SOLID, LoD, Composition, Less is More, Explicit over Implicit)

## REMOVED Requirements

### Requirement: Manual skill selection (v2.x)
**Reason**: User experience was fragmented requiring manual invocation of 5 separate skills
**Migration**: Users now use single command "Review my Go code" instead of invoking go-code-review-gorm, go-code-review-error-safety, etc.

### Requirement: Sequential skill execution (v2.x)
**Reason**: Sequential execution was slow (~4x slower than parallel)
**Migration**: Automatic parallel execution in v3.0 provides same comprehensive coverage with better performance

## ADDED Requirements

### Requirement: Specialist Agent System
The plugin SHALL provide 4 specialist agents that can be invoked independently or automatically through the orchestrator.

#### Scenario: Explicit agent invocation
- **WHEN** user invokes "Use gorm-review agent"
- **THEN** only the GORM review agent executes, checking rules 1.3.*

#### Scenario: Agent independence
- **WHEN** multiple agents execute in parallel
- **THEN** each agent operates independently with no shared state

#### Scenario: Result aggregation by priority
- **WHEN** multiple agents complete
- **THEN** results are merged and sorted by priority (P0, P1, P2)

### Requirement: Design Philosophy Rules (4.*)
The plugin SHALL validate Go code against 8 fundamental design philosophies with 40 specific rules.

#### Scenario: KISS principle validation (4.1.*)
- **WHEN** reviewing Go code
- **THEN** the system checks for unnecessary abstractions, complex functions, and over-engineering

#### Scenario: DRY principle validation (4.2.*)
- **WHEN** reviewing Go code
- **THEN** the system identifies code duplication, magic numbers, and repeated logic

#### Scenario: YAGNI principle validation (4.3.*)
- **WHEN** reviewing Go code
- **THEN** the system flags unused code, premature optimization, and unnecessary features

#### Scenario: SOLID principles validation (4.4.*)
- **WHEN** reviewing Go code
- **THEN** the system checks single responsibility, open/closed, Liskov substitution, interface segregation, and dependency inversion

#### Scenario: Law of Demeter validation (4.5.*)
- **WHEN** reviewing Go code
- **THEN** the system identifies excessive method chaining and tight coupling

#### Scenario: Composition validation (4.6.*)
- **WHEN** reviewing Go code
- **THEN** the system validates use of embedding and interface composition over inheritance

#### Scenario: Minimalism validation (4.7.*)
- **WHEN** reviewing Go code
- **THEN** the system flags reflection overuse, generic abuse, and over-abstraction

#### Scenario: Explicitness validation (4.8.*)
- **WHEN** reviewing Go code
- **THEN** the system ensures explicit error handling, type conversion, and control flow
