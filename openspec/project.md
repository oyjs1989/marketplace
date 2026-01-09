# Project Context

## Purpose

This repository is the **Claude Code Marketplace** - a curated collection of reusable skills, commands, agents, and workflows for Claude Code. The project provides pre-built components that extend Claude Code's capabilities and help developers automate common tasks, implement best practices, and accelerate development workflows.

### Goals
- Share reusable Claude Code components across teams and projects
- Provide production-ready skills with comprehensive documentation
- Enable rapid adoption of coding standards and best practices
- Build a community-driven marketplace of AI-powered development tools

## Tech Stack

### Core Technologies
- **Language**: Markdown (documentation-centric project)
- **Platform**: Claude Code (Anthropic's CLI tool)
- **Version Control**: Git (hosted on GitLab at gitlab.futunn.com)
- **Package Management**: npm (project metadata only)
- **Spec System**: OpenSpec (for change management and requirements)

### Development Tools
- Claude Code skills framework
- OpenSpec CLI for spec-driven development
- Markdown for all documentation

## Project Conventions

### Code Style

**File Naming**
- Use UPPERCASE for root-level documentation: `README.md`, `LICENSE`, `CLAUDE.md`
- Use `kebab-case` for directories: `go-code-review`, `test-cases`
- Use `SKILL.md` for Claude Code skill definitions (uppercase)
- Use `spec.md` and `design.md` for OpenSpec files (lowercase)

**Markdown Conventions**
- Use GitHub Flavored Markdown
- Include frontmatter metadata in SKILL.md files using YAML format
- Use code fences with language identifiers (```go, ```bash, etc.)
- Include emoji sparingly and only when enhancing clarity
- Maintain consistent heading hierarchy (no skipping levels)

**Documentation Standards**
- All skills must include clear "When to Use This Skill" sections
- Provide concrete examples for all features
- Include test cases in `test-cases/` directory
- Document all configuration options and parameters
- Use Chinese (中文) for Go code review output as per FUTU standards

### Architecture Patterns

**Component Organization**
```
marketplace/
├── skills/          # Reusable domain expertise and knowledge
├── commands/        # Quick shortcuts with /command syntax
├── agents/          # Autonomous task executors with tool access
├── workflows/       # Multi-step processes orchestrating the above
├── test-cases/      # Validation tests for skills
└── openspec/        # Spec-driven development artifacts
```

**Skill Architecture Pattern**
- **Orchestrator**: Main SKILL.md that coordinates execution
- **Parallel Agents**: Multiple specialized agents for concurrent execution
- **Shared Standards**: Common reference documents (e.g., FUTU_GO_STANDARDS.md)
- **Modular Design**: Each agent focuses on specific rule categories

**Version Management**
- Use semantic versioning (e.g., v2.0.0)
- Document version in skill metadata
- Maintain backwards compatibility when possible
- Clearly mark breaking changes

### Testing Strategy

**Test Organization**
- Create dedicated test directory for each skill: `test-cases/[skill-name]/`
- Include both "good" and "bad" example files
- Provide comprehensive test documentation in README.md
- Include TESTING_GUIDE.md for detailed validation procedures

**Test Case Requirements**
- Cover all P0/P1/P2 severity levels for code review skills
- Include real-world scenarios and edge cases
- Document expected outcomes clearly
- Provide quick start instructions

**Validation Approach**
- Manual testing through Claude Code sessions
- Document test commands explicitly (e.g., "Review test-cases/go-code-review/bad/user_service_bad.go")
- Track issues found and resolved during testing

### Git Workflow

**Branch Strategy**
- `main` branch: stable, production-ready code
- Feature branches: short-lived, named descriptively
- Create separate PRs for proposals vs implementations

**Commit Conventions**
- Write clear, descriptive commit messages
- Use present tense ("Add feature" not "Added feature")
- Keep commits focused on single concerns
- Include context in commit body when needed

**OpenSpec Integration**
- Create proposal before implementing new features
- Validate proposals with `openspec validate --strict`
- Implement changes following tasks.md checklist
- Archive completed changes to `changes/archive/YYYY-MM-DD-[name]/`

**Pull Request Process**
1. Create proposal in `openspec/changes/[change-id]/`
2. Get proposal approval
3. Implement according to tasks.md
4. Test thoroughly
5. Submit PR for review
6. After merge, create separate PR to archive the change

## Domain Context

### Claude Code Skills
Claude Code skills are structured markdown files (SKILL.md) with:
- YAML frontmatter defining name, description, and trigger conditions
- Clear documentation of when and how to use the skill
- Embedded domain expertise and coding standards
- Execution instructions for the AI assistant

### Go Code Review Domain
- **FUTU Standards**: Internal Go coding standards with 73+ rules
- **Priority Levels**: P0 (critical), P1 (recommended), P2 (optimization)
- **Rule Categories**: GORM operations, error handling, naming, logging, code organization
- **Language**: All review output must be in Chinese (中文)
- **Multi-Agent Model**: Parallel execution of 4 specialized review agents for efficiency

### OpenSpec Methodology
- Spec-driven development approach
- Proposals precede implementation for significant changes
- Delta-based change management (ADDED/MODIFIED/REMOVED/RENAMED)
- Strict validation ensures consistency
- Scenarios must use `#### Scenario:` format with WHEN/THEN structure

## Important Constraints

### Technical Constraints
- Skills must be compatible with Claude Code's execution environment
- All documentation must be in Markdown format
- No runtime dependencies beyond Claude Code itself
- Skills should be self-contained and portable

### Language Constraints
- Go code review output: Chinese (中文) mandatory
- Other documentation: English preferred for broader accessibility
- Code examples: Follow language-specific conventions

### Performance Constraints
- Skills should complete reviews within reasonable time
- Parallel agent execution for improved performance
- Avoid unnecessary file operations

### Business Constraints
- MIT License for all components
- Internal hosting on GitLab (gitlab.futunn.com)
- Must align with FUTU's Go coding standards

## External Dependencies

### Claude Code Platform
- **Provider**: Anthropic
- **Purpose**: AI-powered CLI development assistant
- **Integration**: All skills execute within Claude Code environment
- **Documentation**: Refer to Claude Code official docs for skill framework

### OpenSpec CLI
- **Purpose**: Spec-driven development and change management
- **Usage**: Proposal creation, validation, and archiving
- **Commands**: `openspec list`, `openspec validate`, `openspec archive`

### GitLab
- **Host**: gitlab.futunn.com
- **Repository**: jasonouyang/marketplace
- **Purpose**: Version control and collaboration

### FUTU Standards
- **Source**: Internal FUTU Go coding standards
- **Location**: `skills/go-code-review/shared/FUTU_GO_STANDARDS.md`
- **Rules**: 73+ coding rules across 5 major categories
- **Updates**: Standards evolve; skills must stay synchronized
