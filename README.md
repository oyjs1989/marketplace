# Claude Code Marketplace

A curated collection of reusable skills, commands, agents, and workflows for Claude Code.

## Overview

This repository provides a marketplace of pre-built components that extend Claude Code's capabilities. Whether you're looking to automate common tasks, implement best practices, or accelerate your development workflow, you'll find ready-to-use solutions here.

## Project Structure

```
marketplace/
├── skills/          # Individual reusable skills
│   ├── go-code-review/        # Go code review skill (v4.0.0)
│   ├── problem-solving/       # Problem-solving orchestrator (v1.0.0)
│   ├── decision-support/      # Multi-criteria decision analysis (v1.0.0)
│   ├── risk-assessment/       # Risk evaluation (v1.0.0)
│   ├── cost-benefit-analysis/ # Economic evaluation (v1.0.0)
│   ├── methodology-agile/     # Agile/Scrum methodology (v1.0.0)
│   ├── methodology-devops/    # DevOps methodology (v1.0.0)
│   └── methodology-waterfall/ # Waterfall methodology (v1.0.0)
├── commands/        # Custom slash commands (coming soon)
├── agents/          # Specialized agent configurations (coming soon)
├── workflows/       # Complex multi-step workflows (coming soon)
├── test-cases/      # Test files for validating skills
├── package.json     # Project metadata
└── README.md        # This file
```

## Core Concepts

**Skills**: Domain knowledge and expertise (e.g., coding standards, design patterns)
**Commands**: Quick shortcuts with `/command` syntax
**Agents**: Autonomous task executors with tool access
**Workflows**: Multi-step processes orchestrating the above

Note: Claude Code has built-in support for Skills, Commands, and Agents. Workflows can be implemented through these components.

## Getting Started

### Installation

1. Clone this repository:
```bash
git clone https://gitlab.futunn.com/jasonouyang/marketplace.git
cd marketplace
```

2. Browse the available components in each directory

3. Copy the desired components to your Claude Code project

### Usage

Each component type has its own usage pattern:

- **Skills**: Invoke using skill notation in your Claude Code sessions
- **Commands**: Use with slash notation (e.g., `/command-name`)
- **Agents**: Configure as specialized assistants for specific tasks
- **Workflows**: Follow the documented process for complex operations

## Available Components

### Skills

#### Go Code Review Plugin (v4.0.0)

Three-Tier Expert Architecture for Go code review, enforcing FUTU's 142+ coding standards with quantitative metrics, YAML rule scanning, and domain-expert AI agents.

> **Breaking Change from v3.0.0**: The old 4 specialist agents (gorm-review, error-safety, naming-logging, organization) are replaced by a new three-tier architecture with 5 domain-expert agents and automated tooling.

**Quick Start**:
```
Review my Go code
```

**One command triggers everything** - the plugin automatically:
- Runs quantitative analysis (Tier 1: metrics.json)
- Scans YAML rules for pattern violations (Tier 2: rule-hits.json)
- Dispatches 5 domain-expert agents in parallel (Tier 3)
- Delivers comprehensive findings in Chinese

**Three-Tier Expert Architecture**:
- **Tier 1 - Quantitative Analysis**: `analyze-go.sh` produces `metrics.json` (file size, function length, nesting depth)
- **Tier 2 - Rule Scanning**: `scan-rules.sh` reads 38 YAML rules (SAFE/DATA/QUAL/OBS) and produces `rule-hits.json`
- **Tier 3 - 5 Domain-Expert Agents**:
  - **safety** - Concurrency, nil safety, context propagation (SAFE-001~010)
  - **data** - N+1 queries, GORM operations, type semantics (DATA-001~010)
  - **design** - UNIX 7 principles, code rot causes, architecture
  - **quality** - Naming, metrics violations, code organization (QUAL-001~010)
  - **observability** - Logging strategy, error message quality (OBS-001~008)

**Features**:
- **Automated Tooling**: Tier 1 + Tier 2 run before AI agents to provide quantitative context
- **YAML Rule Source of Truth**: 38 structured rules replace inline rule documentation
- **Single Invocation**: One call reviews with all 142+ rules
- **Parallel Execution**: All 5 agents run simultaneously for maximum speed
- **Chinese Output**: All findings reported in Chinese (FUTU requirement)
- **Priority Classification**: P0 (must fix), P1 (recommended), P2 (suggested)
- **Structured Rule IDs**: SAFE-/DATA-/QUAL-/OBS- prefix naming for traceability

**Documentation**:
- Main: [skills/go-code-review/README.md](skills/go-code-review/README.md)
- Details: [skills/go-code-review/SKILL.md](skills/go-code-review/SKILL.md)

**Test**: Run `Review test-cases/go-code-review/bad/user_service_bad.go` to see it in action.

---

#### Problem-Solving Framework (v1.0.0)

Comprehensive problem-solving capabilities with cognitive agents and specialized analysis modules.

**Core Orchestrator**:
- **problem-solving** - Main orchestrator with 5 cognitive agents:
  - 🔵 **systems-thinking** - Global perspective, emergent properties, leverage points
  - 🟢 **modeling-abstraction** - Conceptual models, DDD patterns, abstraction layers
  - 🟣 **decomposition** - Hierarchical breakdown, WBS, composition strategies
  - 🟠 **iteration** - Incremental improvement, feedback loops, MVP approach
  - 🔷 **pattern-recognition** - Pattern discovery/refactoring, GoF patterns

**Analysis Modules**:
- **decision-support** - Multi-criteria decision analysis (MCDA, AHP, TOPSIS, weighted scoring)
- **risk-assessment** - Risk evaluation using COSO ERM and ISO 31000 frameworks
- **cost-benefit-analysis** - Economic evaluation (ROI, NPV, IRR, TCO, payback period)

**Methodology Skills**:
- **methodology-agile** - Scrum framework (Sprint planning, user stories, velocity, burndown charts)
- **methodology-devops** - CI/CD pipelines, IaC, monitoring, incident response, DORA metrics
- **methodology-waterfall** - Traditional SDLC (6 stages, SRS, SDD, gate reviews, change management)

**Quick Start Examples**:
```
# Comprehensive problem analysis
"Help me analyze this microservices architecture design"
→ Triggers all 5 cognitive agents in parallel

# Decision support
"Compare MongoDB vs PostgreSQL for my use case"
→ Uses MCDA with weighted criteria matrix

# Risk assessment
"Evaluate risks of cloud migration"
→ Generates risk matrix with mitigation strategies

# Cost-benefit analysis
"Should we build or buy a CRM system?"
→ Calculates ROI, payback period, NPV

# Agile planning
"Plan the next Sprint for user authentication"
→ Provides Scrum guidance and templates

# DevOps design
"Design a CI/CD pipeline for Node.js microservices"
→ Provides stage-by-stage pipeline configuration

# Waterfall project
"Create requirements specification for banking system"
→ Generates IEEE 830 SRS template
```

**Integration Patterns**:
- Problem-solving → Risk assessment → Decision support (comprehensive analysis)
- Agile + DevOps (modern software development)
- Waterfall + Cost-benefit analysis (traditional projects)

**Output Language**: All skills output in Chinese (中文) following FUTU standards.

**Documentation**:
- [problem-solving README](skills/problem-solving/README.md)
- [decision-support README](skills/decision-support/README.md)
- [risk-assessment README](skills/risk-assessment/README.md)
- [cost-benefit-analysis README](skills/cost-benefit-analysis/README.md)
- [methodology-agile README](skills/methodology-agile/README.md)
- [methodology-devops README](skills/methodology-devops/README.md)
- [methodology-waterfall README](skills/methodology-waterfall/README.md)

---

### Commands
Coming soon. See [commands/README.md](commands/README.md) for planned commands.

### Agents
Coming soon. See [agents/README.md](agents/README.md) for planned agents.

### Workflows
Coming soon. See [workflows/README.md](workflows/README.md) for planned workflows.

## Contributing

Contributions are welcome! To contribute:

1. Fork this repository
2. Create a new branch for your component
3. Add your skill/command/agent/workflow with proper documentation
4. Test thoroughly
5. Submit a merge request

### Contribution Guidelines

- Follow the existing structure and documentation patterns
- Include clear examples and usage instructions
- Test your components before submitting
- Document any dependencies
- Follow best practices for Claude Code development

## Best Practices

When creating or using components from this marketplace:

1. **Documentation**: Always include clear, comprehensive documentation
2. **Testing**: Thoroughly test components before use in production
3. **Dependencies**: Clearly document any external dependencies
4. **Versioning**: Use semantic versioning for tracking changes
5. **Security**: Review code for security implications
6. **Maintenance**: Keep components up-to-date with Claude Code updates

## Support

For questions, issues, or suggestions:

- Open an issue in the repository
- Review existing documentation in each directory
- Check the examples provided with each component

## License

MIT License - feel free to use, modify, and distribute these components.

## Project Status

Active development. New skills, commands, agents, and workflows are being added regularly.

## Testing

Test cases are available in the `test-cases/` directory to validate skill functionality:

- **Go Code Review Tests**: `test-cases/go-code-review/`
  - Comprehensive test files covering P0/P1/P2 violations
  - Good vs bad code examples
  - Detailed test documentation and quick start guide
  - See [test-cases/go-code-review/README.md](test-cases/go-code-review/README.md)

## Roadmap

- [x] Add initial set of common skills (Go Code Review v4.0.0)
- [x] Create testing framework (Go Code Review test cases)
- [x] Problem-solving framework (v1.0.0) - 7 skills with cognitive agents
- [ ] Integration tests for problem-solving framework
- [ ] Create useful slash commands
- [ ] Develop specialized agents
- [ ] Document complex workflows
- [ ] Build example integrations
- [ ] Add CI/CD pipeline

## Acknowledgments

Built for the Claude Code community to share and reuse valuable development tools.
