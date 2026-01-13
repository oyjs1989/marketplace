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
│   ├── go-code-review/            # Orchestrator skill (v2.0.0)
│   │   ├── SKILL.md               # Main orchestrator
│   │   └── references/            # FUTU_GO_STANDARDS.md (97+ rules)
│   ├── go-code-review-gorm/       # GORM database review skill
│   │   ├── SKILL.md
│   │   └── references/ (symlink)
│   ├── go-code-review-error-safety/# Error & safety review skill
│   │   ├── SKILL.md
│   │   └── references/ (symlink)
│   ├── go-code-review-naming/     # Naming & logging review skill
│   │   ├── SKILL.md
│   │   └── references/ (symlink)
│   ├── go-code-review-organization/# Organization review skill
│   │   ├── SKILL.md
│   │   └── references/ (symlink)
│   └── gitlab-ai-summary/         # GitLab MR summary skill
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

Go Code Review plugin uses a **unified orchestration model** with automatic agent selection:
- **Main Orchestrator**: Single skill coordinates comprehensive review
- **4 Specialist Agents**: Independent agents in `agents/` directory
  - Each agent has its own agent.md file with specific focus
  - Agents auto-triggered based on code patterns
  - Can be invoked independently or through orchestrator
- **Shared References**: FUTU_GO_STANDARDS.md with 97+ rules
- **Smart Selection**: Automatically determines which agents apply
- **Parallel Execution**: All agents run simultaneously for speed

Example: Go Code Review Plugin (v3.0.0):
1. **go-code-review** - Main orchestrator with automatic agent selection
2. **gorm-review** agent (blue) - GORM database review (rules 1.3.*)
3. **error-safety** agent (red) - Error & safety review (rules 1.1.*, 1.2.*, 1.4.*, 1.5.*)
4. **naming-logging** agent (green) - Naming & logging review (rules 2.1.*, 2.2.*)
5. **organization** agent (purple) - Organization & quality review (rules 2.3.*, 2.4.*, 2.5.*, 3.*)

**Key Benefit**: Single invocation (`Review my Go code`) automatically triggers all applicable agents

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

The Go code review skill enforces 73+ rules across 5 categories:
- P0 (Critical): Error handling, nil checks, GORM operations, concurrency, JSON processing
- P1 (Recommended): Naming conventions, logging standards, code organization, interface design, code quality
- P2 (Optimization): Project structure, testing standards, configuration management

Reference: `skills/go-code-review/shared/FUTU_GO_STANDARDS.md`

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
