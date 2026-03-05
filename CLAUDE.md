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
- Skills: Domain expertise (e.g., Go code review with FUTU standards)
- Commands: Custom slash commands via `.claude/commands/`
- Agents: Specialized configurations in `agents/`
- Test Cases: Validation tests in `test-cases/`

## Development Commands

### Testing Skills

**Go Code Review Skill:**
```bash
# Test with bad code examples
Review test-cases/go-code-review/bad/user_service_bad.go

# Run automated tests (from test-cases/go-code-review/)
./run_test.sh

# Validate test results
./validate_results.sh
```

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
│   ├── go-code-review/            # Go code review plugin (v4.0.0)
│   │   ├── SKILL.md               # Main orchestrator (三层架构)
│   │   ├── agents/                # 5 领域专家 agents
│   │   │   ├── safety.md          # 安全性专家
│   │   │   ├── data.md            # 数据层专家
│   │   │   ├── design.md          # 架构设计哲学专家
│   │   │   ├── quality.md         # 代码质量专家
│   │   │   └── observability.md   # 可观测性专家
│   │   ├── rules/                 # YAML 结构化规则（单一数据源）
│   │   │   ├── safety.yaml        # SAFE-001~010
│   │   │   ├── data.yaml          # DATA-001~010
│   │   │   ├── quality.yaml       # QUAL-001~010
│   │   │   └── observability.yaml # OBS-001~008
│   │   └── tools/
│   │       ├── analyze-go.sh      # 量化分析 → metrics.json
│   │       └── scan-rules.sh      # YAML 规则扫描 → rule-hits.json
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
├── test-cases/          # Test files for validation
│   └── go-code-review/  # Comprehensive test suite with validation scripts
├── openspec/            # Spec-driven development artifacts
│   ├── AGENTS.md        # OpenSpec workflow instructions
│   ├── project.md       # Project conventions and standards
│   ├── specs/           # Current specifications (what IS built)
│   └── changes/         # Active and archived change proposals
└── .claude/             # Claude Code configuration
    └── commands/        # SpecKit and OpenSpec commands
```

### Skill Architecture Pattern

Go Code Review plugin (v4.0.0) uses a **Three-Tier Expert Architecture**:
- **Tier 1 - Quantitative Tools**: `analyze-go.sh` measures file size, function length, and nesting depth, producing `metrics.json` before AI agents run
- **Tier 2 - YAML Rule Scanning**: `scan-rules.sh` reads 38 structured YAML rules (SAFE/DATA/QUAL/OBS) and produces `rule-hits.json` with pattern-matched violations
- **Tier 3 - 5 Domain-Expert Agents**: Independent agents each consuming Tier 1 + Tier 2 output for context-aware review
  - Each agent has its own agent.md with a specific domain focus
  - Agents run in parallel after Tier 1 and Tier 2 complete
  - Can be invoked independently or through the orchestrator
- **YAML as Single Source of Truth**: Rules defined in `rules/*.yaml`, not inline in agent prompts
- **Parallel Execution**: All 5 agents run simultaneously for speed

Example: Go Code Review Plugin (v4.0.0):
1. **go-code-review** (SKILL.md) - Main orchestrator: runs tools, then dispatches agents
2. **safety** agent - Concurrency, nil safety, context propagation (SAFE-001~010)
3. **data** agent - N+1 queries, GORM operations, type semantics (DATA-001~010)
4. **design** agent - UNIX 7 principles, code rot causes, architecture philosophy
5. **quality** agent - Naming, metrics violations, code organization (QUAL-001~010)
6. **observability** agent - Logging strategy, error message quality (OBS-001~008)

**Key Benefit**: Single invocation (`Review my Go code`) runs quantitative tooling then triggers all 5 domain-expert agents in parallel

## File Naming Conventions

- Root documentation: `UPPERCASE.md` (README.md, CLAUDE.md, LICENSE)
- Directories: `kebab-case` (go-code-review, test-cases)
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

## Special Requirements

### Go Code Review Output

All Go code review output **MUST be in Chinese (中文)** per FUTU standards:
- Problem descriptions: 中文
- Suggestions: 中文
- File paths and code: Keep original
- Severity levels: P0 (必须修复), P1 (强烈建议), P2 (建议优化)
- Rule IDs: Use structured IDs — SAFE-001~010, DATA-001~010, QUAL-001~010, OBS-001~008
  - Example: `[P0][SAFE-001] 使用 fmt.Errorf 会丢失错误堆栈，应使用 errors.Wrap`

### Test Case Structure

For skills with validation:
- Create `test-cases/[skill-name]/` directory
- Include `good/` and `bad/` example subdirectories
- Provide `README.md`, `TESTING_GUIDE.md`, and validation scripts
- Document expected issues in `expected_issues.json`

## Version Management

- Use semantic versioning (e.g., v2.0.0)
- Document version in skill frontmatter
- Track changes in git commit messages
- Major version for breaking changes

## Important Context

### FUTU Go Standards

The Go code review skill (v4.0.0) enforces 142+ rules across 4 YAML rule domains plus agent-driven design principles:
- **SAFE-001~010**: Error handling, nil checks, concurrency, context propagation
- **DATA-001~010**: GORM operations, N+1 queries, JSON processing, type semantics
- **QUAL-001~010**: Naming conventions, code organization, magic numbers, metrics thresholds
- **OBS-001~008**: Logging standards, log field naming, error field requirements, data layer logging
- **Design principles**: UNIX 7 principles + code rot causes (enforced by design agent)

Rules are defined in `skills/go-code-review/rules/*.yaml` (single source of truth).

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
