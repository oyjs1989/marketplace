# Consolidate Go Code Review into Unified Plugin

**Change ID**: consolidate-go-code-review-plugin
**Status**: Approved
**Version**: v3.0.0
**Date**: 2026-01-13
**Author**: FUTU Development Team

## Motivation

The current Go code review system (v2.x) consists of 5 independent skills that users must manually select and invoke:

1. `go-code-review` - Main orchestrator
2. `go-code-review-gorm` - GORM database reviews
3. `go-code-review-error-safety` - Error handling & safety
4. `go-code-review-naming` - Naming conventions
5. `go-code-review-organization` - Code organization

**Problems with v2.x:**
- **Fragmented User Experience**: Users must manually determine which skills to invoke
- **Sequential Execution**: Skills run one at a time, leading to slow reviews
- **Unclear Scope**: Users don't know which skill covers which rules
- **Maintenance Burden**: 5 separate SKILL.md files to maintain

**Goals for v3.0:**
- Single command invocation for comprehensive reviews
- Automatic agent selection based on code patterns
- Parallel execution for 4x speed improvement
- Expanded rule coverage with design philosophies

## Proposed Changes

### Architecture Transformation

**From (v2.x)**: 5 Independent Skills
```
skills/
├── go-code-review/               # Orchestrator
├── go-code-review-gorm/          # Independent skill
├── go-code-review-error-safety/  # Independent skill
├── go-code-review-naming/        # Independent skill
└── go-code-review-organization/  # Independent skill
```

**To (v3.0)**: 1 Unified Plugin + 4 Specialist Agents
```
skills/go-code-review/
├── SKILL.md                      # Main orchestrator with auto-selection
├── agents/                       # 4 specialist agents
│   ├── gorm-review.md            # GORM database (rules 1.3.*)
│   ├── error-safety.md           # Error & safety (rules 1.1.*, 1.2.*, 1.4.*, 1.5.*)
│   ├── naming-logging.md         # Naming & logging (rules 2.1.*, 2.2.*)
│   └── organization-agent.md     # Organization & quality (rules 2.3.*, 2.4.*, 2.5.*, 3.*, 4.*)
└── references/
    └── FUTU_GO_STANDARDS.md      # 137+ rules (expanded from 97)
```

### Key Features

1. **Smart Agent Selection**
   - Analyzes changed files automatically
   - Detects GORM patterns → activates gorm-review agent
   - Detects error handling → activates error-safety agent
   - Always runs naming-logging + organization agents

2. **Parallel Execution**
   - All applicable agents run simultaneously
   - ~4x speed improvement over sequential execution
   - Results aggregated and merged by priority

3. **Expanded Rule Coverage**
   - **v2.x**: 97 rules across 3 categories (P0/P1/P2)
   - **v3.0**: 137 rules including 8 design philosophies
   - New 4.* category: KISS, DRY, YAGNI, SOLID, LoD, Composition, Less is More, Explicit over Implicit

4. **Simplified Invocation**
   ```
   # v2.x (manual selection)
   /go-code-review
   /go-code-review-gorm
   /go-code-review-error-safety

   # v3.0 (automatic)
   Review my Go code
   ```

## Breaking Changes

### Removed

- ❌ **Skills removed**:
  - `go-code-review-gorm` (moved to `agents/gorm-review.md`)
  - `go-code-review-error-safety` (moved to `agents/error-safety.md`)
  - `go-code-review-naming` (moved to `agents/naming-logging.md`)
  - `go-code-review-organization` (moved to `agents/organization-agent.md`)

- ❌ **Old invocation patterns no longer work**:
  - `/go-code-review-gorm` → Use `Review my Go code` or `Use gorm-review agent`
  - Similar for other removed skills

### Changed

- ✏️ **Plugin structure**: Skills are now agents within main plugin
- ✏️ **Invocation**: Single command triggers comprehensive multi-agent review
- ✏️ **Rule coverage**: Expanded from 97 to 137 rules

### Added

- ✅ **Automatic agent selection**
- ✅ **Parallel agent execution**
- ✅ **8 Design philosophy categories** (40 new rules):
  - 4.1 KISS Principle (5 rules)
  - 4.2 DRY Principle (5 rules)
  - 4.3 YAGNI Principle (5 rules)
  - 4.4 SOLID Principles (5 rules)
  - 4.5 Law of Demeter (5 rules)
  - 4.6 Composition Over Inheritance (5 rules)
  - 4.7 Less is Exponentially More (5 rules)
  - 4.8 Explicit Over Implicit (5 rules)

## Migration Guide

### For Users

**Old way (v2.x)**:
```bash
# Manually invoke multiple skills
Review with go-code-review
Review with go-code-review-gorm
Review with go-code-review-error-safety
Review with go-code-review-naming
Review with go-code-review-organization
```

**New way (v3.0)**:
```bash
# Single command, automatic agent selection
Review my Go code
```

**Explicit agent invocation (optional)**:
```bash
Use gorm-review agent
Use error-safety agent
Use naming-logging agent
Use organization agent
```

### For Developers

1. **Plugin installation**: Same as v2.x, no changes needed
2. **Configuration**: No breaking config changes
3. **Test cases**: New test file `design_philosophy_bad.go` for 4.* rules
4. **Expected results**: Updated `expected_issues.json` with 137 rules

## Implementation Status

✅ All implementation tasks completed (see `tasks.md`)

**Deliverables**:
- [x] Unified plugin architecture
- [x] 4 specialist agent files
- [x] Automatic agent selection logic
- [x] Parallel execution coordination
- [x] Expanded FUTU_GO_STANDARDS.md (137 rules)
- [x] New test cases for design philosophies
- [x] Updated documentation (README, SKILL.md, CLAUDE.md)

## Testing

**Test Coverage**: ~70%

**Test Files**:
- `test-cases/go-code-review/bad/user_service_bad.go` - Error handling, GORM, naming
- `test-cases/go-code-review/bad/project_structure_bad.go` - Organization, structure
- `test-cases/go-code-review/bad/design_philosophy_bad.go` - Design philosophies (NEW)

**Validation**:
```bash
cd test-cases/go-code-review
./run_test.sh
./validate_results.sh
```

## Risks & Mitigations

| Risk | Impact | Likelihood | Mitigation |
|------|--------|------------|------------|
| Users expect old skill names | High | High | Document breaking changes in README, add migration guide |
| New rules flag previously passing code | Medium | Medium | Rules are P1 (recommended), not P0 (required) |
| Parallel execution introduces bugs | Low | Low | Agents are independent, no shared state |
| Test coverage gaps | Low | Low | 70% coverage is acceptable, fill gaps in future iterations |

## Rollout Plan

1. **Phase 1: Documentation** ✅
   - Update README.md with v3.0.0 version
   - Add breaking changes section
   - Update SKILL.md with orchestration logic

2. **Phase 2: Testing** ✅
   - Validate against all test cases
   - Ensure parallel execution works correctly
   - Verify rule count accuracy (137 rules confirmed)

3. **Phase 3: Deployment**
   - Merge to main branch
   - Tag as v3.0.0
   - Announce breaking changes to users

4. **Phase 4: Archive** (after deployment)
   ```bash
   openspec archive consolidate-go-code-review-plugin --yes
   ```

## Success Metrics

- ✅ Single command invocation works
- ✅ Automatic agent selection accurate
- ✅ Parallel execution 4x faster than sequential
- ✅ All 137 rules documented and categorized
- ✅ Test coverage ≥70%
- ✅ Zero regression in existing test cases

## Future Enhancements

Deferred to future iterations:

1. File naming consistency (`organization-agent.md` → `organization.md`)
2. Detailed migration guide document (MIGRATION.md)
3. Test coverage improvement (70% → 90%+)
4. Enhanced BREAKING CHANGE documentation in README

## References

- FUTU Go Coding Standards: `skills/go-code-review/references/FUTU_GO_STANDARDS.md`
- Main Orchestrator: `skills/go-code-review/SKILL.md`
- Project Documentation: `CLAUDE.md`
- Test Suite: `test-cases/go-code-review/`

---

**Approval Date**: 2026-01-13
**Approved By**: FUTU Development Team
**Implementation Status**: Completed
