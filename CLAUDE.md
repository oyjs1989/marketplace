# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

<!-- OPENSPEC:START -->
# OpenSpec Instructions

These instructions are for AI assistants working in this project.

Always open `@/openspec/AGENTS.md` when the request:
- Mentions planning or proposals (words like proposal, spec, change, plan)
- Introduces new capabilities, breaking changes, architecture shifts, or big performance/security work
- Sounds ambiguous and you need the authoritative spec before coding

Use `@/openspec/AGENTS.md` to learn:
- How to create and apply change proposals
- Spec format and conventions
- Project structure and guidelines

Keep this managed block so 'openspec update' can refresh the instructions.

<!-- OPENSPEC:END -->

## Project Overview

This is the **Claude Code Marketplace** - a curated collection of reusable skills, commands, agents, and workflows for Claude Code. The repository provides production-ready components with comprehensive documentation and testing.

**Key Components:**
- Skills: Domain expertise (e.g., problem-solving, decision-support)
- Commands: Custom slash commands via `.claude/commands/`
- Agents: Specialized configurations in `agents/`

## Development Commands

### OpenSpec Workflow

OpenSpec is used for spec-driven development. Always create proposals before implementing significant changes.

```bash
# List active changes
openspec list

# List all specifications
openspec list --specs

# Show details of a change or spec
openspec show [change-id or spec-id]

# Validate a change proposal
openspec validate [change-id] --strict

# Archive completed change (after deployment)
openspec archive <change-id> --yes
```

**OpenSpec workflow:**
1. Create proposal in `openspec/changes/[change-id]/` (use kebab-case, verb-led naming)
2. Validate with `openspec validate [change-id] --strict`
3. Get approval before implementation
4. Implement following `tasks.md` checklist
5. After merge and deploy, archive to `changes/archive/YYYY-MM-DD-[name]/`

### Version Control

```bash
# View recent changes
git log --oneline -10

# Check diff against main
git diff main

# View current status
git status
```

## Architecture

### Directory Structure

```
marketplace/
├── skills/                        # Reusable domain expertise
│   ├── gitlab-ai-summary/         # GitLab MR summary skill
│   ├── problem-solving/           # Problem-solving orchestrator (v1.0.0)
│   │   ├── SKILL.md               # Main orchestrator with 5 agents
│   │   ├── agents/                # 5 cognitive agents
│   │   │   ├── systems-thinking.md
│   │   │   ├── modeling-abstraction.md
│   │   │   ├── decomposition.md
│   │   │   ├── iteration.md
│   │   │   └── pattern-recognition.md
│   │   └── references/            # Thinking methods & problem patterns
│   ├── decision-support/          # Multi-criteria decision analysis
│   │   ├── SKILL.md
│   │   └── references/            # Decision frameworks (AHP, TOPSIS)
│   ├── risk-assessment/           # Risk evaluation and management
│   │   ├── SKILL.md
│   │   └── references/            # Risk frameworks (COSO ERM, ISO 31000)
│   ├── cost-benefit-analysis/     # Economic evaluation
│   │   ├── SKILL.md
│   │   └── references/            # CBA methods (ROI, NPV, TCO)
│   ├── methodology-agile/         # Agile/Scrum methodology
│   │   ├── SKILL.md
│   │   └── references/            # Agile framework (Sprint, user stories)
│   ├── methodology-devops/        # DevOps methodology
│   │   ├── SKILL.md
│   │   └── references/            # DevOps framework (CI/CD, IaC, monitoring)
│   └── methodology-waterfall/     # Waterfall methodology
│       ├── SKILL.md
│       └── references/            # Waterfall framework (6 stages, SRS, SDD)
├── agents/              # Agent configurations (agent.md files)
├── commands/            # Custom slash commands (empty, coming soon)
├── workflows/           # Multi-step workflows (coming soon)
├── openspec/            # Spec-driven development artifacts
│   ├── AGENTS.md        # OpenSpec workflow instructions
│   ├── project.md       # Project conventions and standards
│   ├── specs/           # Current specifications (what IS built)
│   └── changes/         # Active and archived change proposals
└── .claude/             # Claude Code configuration
    └── commands/        # SpecKit and OpenSpec commands
```

## File Naming Conventions

- Root documentation: `UPPERCASE.md` (README.md, CLAUDE.md, LICENSE)
- Directories: `kebab-case` (problem-solving, decision-support)
- Skill definitions: `SKILL.md` (uppercase, required name)
- OpenSpec files: `spec.md`, `design.md`, `proposal.md`, `tasks.md` (lowercase)
- Change IDs: `kebab-case` with verb prefix (add-feature, update-logic, remove-legacy)

## Documentation Standards

### SKILL.md Files

Every skill must include:
- YAML frontmatter with `name` and `description`
- "When to Use This Skill" section with trigger keywords
- Clear architecture description if using agents
- Workflow with step-by-step instructions
- Output format specifications
- Examples and test commands

### OpenSpec Proposals

When creating change proposals:
- Must have: `proposal.md`, `tasks.md`
- Optional: `design.md` (only if cross-cutting, new dependencies, or security/performance concerns)
- Required for: new features, breaking changes, architecture changes
- Not required for: bug fixes, typos, dependency updates, config changes

**Spec Delta Format:**
```markdown
## ADDED Requirements
### Requirement: Feature Name
Description with SHALL/MUST

#### Scenario: Success case
- **WHEN** action
- **THEN** result

## MODIFIED Requirements
### Requirement: Existing Feature
(Complete updated requirement with all scenarios)

## REMOVED Requirements
### Requirement: Old Feature
**Reason**: Why removing
**Migration**: How to handle
```

## Version Management

- Use semantic versioning (e.g., v2.0.0)
- Document version in skill frontmatter
- Track changes in git commit messages
- Major version for breaking changes

## Important Context

### Problem-Solving Framework

The marketplace includes a comprehensive problem-solving framework with 7 skills:

**Core Orchestrator:**
- `problem-solving` - Main orchestrator with 5 cognitive agents that run in parallel:
  - systems-thinking (🔵): Global perspective, emergent properties, leverage points
  - modeling-abstraction (🟢): Conceptual models, abstraction layers, DDD patterns
  - decomposition (🟣): Hierarchical breakdown, WBS, composition strategies
  - iteration (🟠): Incremental improvement, feedback loops, MVP approach
  - pattern-recognition (🔷): Pattern discovery/refactoring, GoF patterns, anti-patterns

**Analysis Modules:**
- `decision-support` - Multi-criteria decision analysis (MCDA, AHP, TOPSIS)
- `risk-assessment` - Risk evaluation using COSO ERM and ISO 31000 frameworks
- `cost-benefit-analysis` - Economic evaluation (ROI, NPV, IRR, TCO)

**Methodology Skills:**
- `methodology-agile` - Scrum framework with Sprint planning and iterative development
- `methodology-devops` - CI/CD pipelines, IaC, monitoring, and incident response
- `methodology-waterfall` - Traditional SDLC with 6 stages and document templates

**Usage Examples:**
```
# Comprehensive problem analysis
"Help me analyze this architecture design problem"
→ Triggers problem-solving orchestrator with all 5 agents

# Decision support
"Help me choose between MongoDB and PostgreSQL"
→ Uses decision-support with weighted criteria

# Risk assessment
"Evaluate the risks of migrating to microservices"
→ Uses risk-assessment with risk matrix

# Economic evaluation
"Is it worth building this feature in-house vs buying?"
→ Uses cost-benefit-analysis with ROI calculation

# Agile planning
"Help me plan the next Sprint"
→ Uses methodology-agile with Scrum guidance

# DevOps design
"Design a CI/CD pipeline for my Node.js app"
→ Uses methodology-devops with pipeline stages

# Waterfall project
"Create a requirements specification for the banking system"
→ Uses methodology-waterfall with SRS template
```

**Integration Patterns:**
- Problem-solving → Risk assessment → Decision support (comprehensive analysis)
- Agile methodology + DevOps practices (modern development)
- Waterfall methodology + Cost-benefit analysis (traditional projects)

**Output Language:** All problem-solving skills output in Chinese (中文) following FUTU standards.

### OpenSpec Integration

This project uses OpenSpec for spec-driven development:
- Read `openspec/AGENTS.md` for complete workflow
- Always validate with `--strict` flag before submitting
- Scenarios must use `#### Scenario:` format (4 hashtags)
- Every requirement needs at least one scenario
- Use ADDED/MODIFIED/REMOVED/RENAMED sections for deltas

### Available Commands

SpecKit commands (in `.claude/commands/`):
- `/speckit.specify` - Create feature specification
- `/speckit.plan` - Execute implementation planning
- `/speckit.tasks` - Generate dependency-ordered tasks
- `/speckit.implement` - Execute implementation plan
- `/speckit.analyze` - Cross-artifact consistency analysis
- `/speckit.clarify` - Identify underspecified areas
- `/speckit.checklist` - Generate custom checklist
- `/speckit.constitution` - Create/update project constitution
- `/speckit.taskstoissues` - Convert tasks to GitHub issues

OpenSpec commands:
- `/openspec:proposal` - Scaffold new change proposal
- `/openspec:apply` - Implement approved change
- `/openspec:archive` - Archive deployed change

## Testing and Validation

Before submitting changes:
1. Test skills with provided test cases
2. Run validation scripts if available
3. Validate OpenSpec proposals with `openspec validate --strict`
4. Ensure all documentation is updated
5. Verify examples work as documented

## License

MIT License - Free to use, modify, and distribute.
