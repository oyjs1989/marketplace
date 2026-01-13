<!--
Sync Impact Report:
Version: 0.0.0 → 1.0.0
Change Type: MAJOR (Initial constitution ratification)
Modified Principles: N/A (initial creation)
Added Sections:
  - Core Principles (I-VI covering quality, testing, UX, performance, documentation, versioning)
  - Quality Gates
  - Development Workflow
  - Governance
Templates Status:
  ✅ spec-template.md - Already aligned with constitution requirements (User Stories, Requirements, Success Criteria)
  ✅ plan-template.md - Already includes Constitution Check section
  ✅ tasks-template.md - Already organized by user story with test-first approach
  ⚠ checklist-template.md - May need review to ensure quality gates alignment (deferred)
  ⚠ agent-file-template.md - May need review for consistency requirements (deferred)
Follow-up TODOs:
  - Review checklist-template.md for quality gates alignment
  - Review agent-file-template.md for parallel execution standards
  - Consider adding performance benchmarking guidelines in future amendment
-->

# Claude Code Marketplace Constitution

## Core Principles

### I. Code Quality (NON-NEGOTIABLE)

All components in the marketplace MUST meet production-ready standards:

- **Comprehensive Documentation**: Every skill, command, agent, or workflow MUST include:
  - YAML frontmatter with name and description
  - Clear "When to Use" section with trigger keywords
  - Step-by-step workflow instructions
  - Concrete examples with expected outcomes
  - File path references (e.g., `file.go:42`)

- **Zero Defects Policy**: Components MUST be thoroughly tested before submission. No untested
  code SHALL be merged. Use test-cases directory structure with validation scripts.

- **Consistency Standards**: All components MUST follow established patterns:
  - File naming: UPPERCASE for root docs (README.md, CLAUDE.md), kebab-case for directories,
    SKILL.md for skills (uppercase required)
  - Markdown: GitHub Flavored Markdown with code fences, consistent heading hierarchy
  - Version management: Semantic versioning (MAJOR.MINOR.PATCH)

- **Self-Contained Components**: Skills MUST be portable and executable without external
  runtime dependencies beyond Claude Code itself. Document any platform-specific constraints.

**Rationale**: The marketplace represents the quality standard for Claude Code community.
Every component is a reference implementation that users will copy and adapt. Poor quality
damages trust and adoption.

### II. Testing Standards (NON-NEGOTIABLE)

Testing is mandatory for all skills with executable validation requirements:

- **Test Structure**: Create `test-cases/[skill-name]/` with:
  - `good/` and `bad/` example subdirectories
  - `README.md` documenting test scenarios
  - `TESTING_GUIDE.md` with detailed validation procedures
  - Validation scripts (e.g., `run_test.sh`, `validate_results.sh`)
  - `expected_issues.json` documenting expected findings

- **Coverage Requirements**: Tests MUST cover:
  - All severity levels (P0/P1/P2 for code review skills)
  - Real-world scenarios and edge cases
  - Both positive (good examples) and negative (bad examples) cases
  - Integration with Claude Code execution environment

- **Documentation of Expected Outcomes**: Each test MUST document:
  - What input is provided
  - What output is expected
  - How to verify the test passes
  - Quick start command (e.g., "Review test-cases/go-code-review/bad/file.go")

- **Manual Testing Protocol**: Before submission, execute tests through Claude Code sessions
  and document results. Automated scripts are encouraged but manual validation is required.

**Rationale**: Skills execute in user environments with real code. Untested skills create
support burden, damage reputation, and risk incorrect guidance. Testing ensures reliability.

### III. User Experience Consistency

All marketplace components MUST provide consistent, predictable user experiences:

- **Activation Triggers**: Skills MUST document clear trigger keywords and activation
  conditions in the description field. Users should know exactly when to invoke each skill.

- **Output Format Standardization**: Skills that produce structured output MUST:
  - Document the exact output format
  - Use consistent severity levels (P0/P1/P2 or equivalent)
  - Include file paths and line numbers for actionable feedback
  - Save results to documented output files (e.g., `code_review.result`)

- **Language Consistency**: Follow established language requirements:
  - Go code review output: Chinese (中文) mandatory per FUTU standards
  - General documentation: English for broader accessibility
  - Code examples: Language-specific conventions

- **Error Handling**: Skills MUST handle common error scenarios gracefully:
  - Missing files or directories
  - Invalid input formats
  - Unexpected code structures
  - Provide clear error messages with remediation guidance

- **Architecture Transparency**: Multi-agent skills MUST document their execution model:
  - Orchestrator responsibilities
  - Parallel vs sequential agent execution
  - Rule category assignments per agent
  - Expected completion time ranges

**Rationale**: Inconsistent experiences frustrate users and reduce adoption. Users invest
time learning patterns; breaking those patterns wastes their investment and damages trust.

### IV. Performance Requirements

Skills MUST complete execution within reasonable timeframes:

- **Parallel Execution**: When multiple independent analyses are required, skills MUST use
  parallel agent execution patterns (e.g., 4 concurrent review agents for Go code review).

- **Resource Efficiency**: Skills MUST:
  - Avoid unnecessary file reads or writes
  - Cache shared reference documents (e.g., FUTU_GO_STANDARDS.md)
  - Process files once per execution
  - Minimize redundant analysis

- **Scalability Considerations**: Skills MUST handle:
  - Large files (>1000 lines)
  - Multiple files in single execution
  - Complex codebases with deep nesting
  - Document known limitations explicitly

- **Performance Benchmarks**: Complex skills (>3 agents, >50 rules) SHOULD document:
  - Typical execution time for standard test cases
  - Performance characteristics with large inputs
  - Optimization strategies employed

**Rationale**: Slow skills interrupt developer flow. The marketplace competes with manual
processes; if skills are slower than human review, adoption fails.

### V. Documentation Excellence

Documentation is a first-class deliverable, not an afterthought:

- **Architecture Documentation**: For multi-component skills:
  - Document parallel agent execution model
  - Explain rule category distribution
  - Show coordination patterns
  - Reference agent configuration files

- **Runnable Examples**: Every skill MUST provide:
  - Copy-pasteable test commands
  - Expected output samples
  - Links to test files in repository
  - Quick start instructions (<5 steps to first execution)

- **Maintenance Documentation**: Skills MUST document:
  - How to update rule sets (e.g., adding rules to FUTU_GO_STANDARDS.md)
  - Where standards are stored (e.g., `shared/` directory)
  - How to test changes
  - Version compatibility requirements

- **Context for Reviewers**: README files MUST include:
  - Project overview and goals
  - Directory structure with annotations
  - Development commands (build, test, validate)
  - OpenSpec workflow integration
  - Available slash commands

**Rationale**: Undocumented components are unusable. The marketplace is self-service;
documentation is the only interface between component authors and users.

### VI. Versioning & Breaking Changes

All components MUST follow semantic versioning and change management protocols:

- **Version Format**: MAJOR.MINOR.PATCH (e.g., v2.0.0)
  - MAJOR: Breaking changes (incompatible with previous versions)
  - MINOR: New features, backwards-compatible additions
  - PATCH: Bug fixes, documentation updates, non-functional changes

- **Breaking Change Requirements**: When introducing breaking changes, components MUST:
  - Document the change in version metadata
  - Provide migration guide
  - Mark deprecated features clearly
  - Allow transition period when possible (MINOR version with deprecation warnings)

- **OpenSpec Integration**: Significant changes MUST follow OpenSpec workflow:
  - Create proposal in `openspec/changes/[change-id]/`
  - Validate with `openspec validate --strict`
  - Get approval before implementation
  - Implement following `tasks.md` checklist
  - Archive to `changes/archive/YYYY-MM-DD-[name]/` after deployment

- **Backwards Compatibility**: MINOR and PATCH versions MUST maintain:
  - Existing skill activation triggers
  - Output format structure
  - File naming conventions
  - API contracts (for commands/agents)

**Rationale**: Users build workflows around marketplace components. Breaking changes without
migration paths damage production systems and erode trust.

## Quality Gates

Before merging any component, the following gates MUST pass:

### Pre-Submission Checklist

- [ ] **Documentation Complete**: SKILL.md/README.md with all required sections
- [ ] **Examples Provided**: Runnable commands with expected outputs documented
- [ ] **Tests Created**: Test cases in `test-cases/[component-name]/` with validation
- [ ] **Tests Executed**: Manual validation through Claude Code completed successfully
- [ ] **Validation Scripts**: Automated validation scripts pass (if applicable)
- [ ] **Version Specified**: Semantic version documented in metadata
- [ ] **File Naming**: Follows conventions (UPPERCASE, kebab-case, SKILL.md)
- [ ] **Performance Verified**: Execution completes within reasonable timeframes
- [ ] **Error Handling**: Common error scenarios handled gracefully

### Code Review Requirements

All merge requests MUST be reviewed for:

- [ ] **Standards Compliance**: Adheres to all Core Principles (I-VI)
- [ ] **Documentation Quality**: Clear, comprehensive, and accurate
- [ ] **Test Coverage**: Adequate test cases with expected outcomes
- [ ] **Performance**: No unnecessary blocking operations or redundant processing
- [ ] **Consistency**: Matches existing patterns in the marketplace
- [ ] **OpenSpec Compliance**: Proposals approved, tasks completed (for significant changes)

### Post-Merge Validation

After merging:

- [ ] **Integration Test**: Verify component works in clean Claude Code environment
- [ ] **Documentation Build**: Ensure all links and references resolve correctly
- [ ] **Announcement**: Update README.md with new component (if major addition)
- [ ] **Archive OpenSpec Change**: Move proposal to archive with `openspec archive`

## Development Workflow

### For New Components

1. **Plan**: Review existing components for patterns. Draft SKILL.md structure.
2. **Implement**: Create component following Core Principles I-VI.
3. **Test**: Create comprehensive test cases in `test-cases/`.
4. **Document**: Write complete documentation with examples.
5. **Validate**: Execute tests manually through Claude Code.
6. **Review**: Self-review against Quality Gates checklist.
7. **Submit**: Create merge request with completed checklist.

### For Significant Changes (OpenSpec Required)

Changes requiring OpenSpec proposals include:
- New capabilities or features
- Breaking changes to existing components
- Architecture modifications
- Performance-impacting changes

**Workflow**:
1. Create proposal in `openspec/changes/[change-id]/` (use kebab-case, verb-led naming)
2. Write `proposal.md` (why, what, impact) and `tasks.md` (implementation checklist)
3. Create spec deltas in `specs/[capability]/spec.md` (ADDED/MODIFIED/REMOVED)
4. Validate: `openspec validate [change-id] --strict`
5. Get approval before implementing
6. Implement following `tasks.md`
7. After merge and deployment, archive: `openspec archive <change-id> --yes`

**Not requiring OpenSpec proposals**:
- Bug fixes restoring intended behavior
- Typos, formatting, comments
- Dependency updates (non-breaking)
- Configuration changes
- Tests for existing behavior

### For Minor Updates (Direct Implementation)

For minor changes (bug fixes, documentation improvements):
1. Make changes following existing patterns
2. Test changes
3. Update version (PATCH increment)
4. Submit merge request with clear description
5. Review for Quality Gates compliance

## Governance

### Constitution Authority

This constitution supersedes all other practices and guidelines. In case of conflict between
this constitution and other documentation, the constitution takes precedence.

### Amendment Process

Amendments require:
1. **Proposal**: Document proposed change with rationale
2. **Impact Analysis**: Assess effect on existing components
3. **Community Input**: Allow review period for feedback
4. **Approval**: Maintainer approval required
5. **Migration Plan**: Provide guidance for adapting existing components
6. **Version Update**: Increment constitution version (MAJOR/MINOR/PATCH)

### Compliance Review

All merge requests MUST verify compliance with this constitution. Reviewers SHALL:
- Check Core Principles adherence
- Validate Quality Gates completion
- Ensure Development Workflow followed
- Request changes for non-compliance

### Complexity Justification

Any violation of simplicity principles (e.g., >4 parallel agents, multiple execution models)
MUST be justified in the proposal with:
- Specific problem being solved
- Simpler alternatives considered and rejected
- Measurable benefits of added complexity

### Runtime Guidance

For day-to-day development guidance beyond constitutional requirements, refer to:
- **CLAUDE.md**: Instructions for Claude Code instances working in this repository
- **openspec/AGENTS.md**: OpenSpec workflow and conventions
- **openspec/project.md**: Detailed project conventions and standards

**Version**: 1.0.0 | **Ratified**: 2026-01-12 | **Last Amended**: 2026-01-12
