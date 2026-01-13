# Data Model: Go Code Review Optimization

**Feature**: Go Code Review Optimization
**Date**: 2026-01-12
**Purpose**: Define key entities and their relationships for the optimization feature

## Overview

This data model describes the core entities involved in the Go code review optimization system. These entities represent conceptual abstractions, not database schemas or code classes.

---

## Entity 1: Detection Rule

**Description**: Represents a single code review rule from FUTU_GO_STANDARDS.md with its detection metadata.

### Fields

| Field Name | Type | Description | Constraints |
|------------|------|-------------|-------------|
| `rule_id` | string | Unique rule identifier | Format: X.X.X (e.g., "1.3.7"), Required |
| `rule_name` | string | Human-readable rule name | Chinese description, Required |
| `priority` | enum | Severity level | Values: P0, P1, P2, Required |
| `detection_algorithm_type` | enum | Classification of detection approach | Values: pattern, structural, semantic, Required |
| `confidence_level` | integer | Detection confidence percentage | Range: 65-95%, Required |
| `known_limitations` | text | Description of detection blind spots | Optional, markdown format |
| `agent_assignment` | string | Which agent handles this rule | Values: gorm, error-safety, naming-logging, organization, Required |
| `example_violation` | text | Code snippet showing violation | Go code block, Optional |
| `example_fix` | text | Code snippet showing correct usage | Go code block, Optional |

### Relationships

- **Belongs to Agent** (many-to-one): Each rule is assigned to exactly one specialized agent
- **Appears in Test Files** (many-to-many): A rule can appear in multiple test files, a test file contains multiple rules

### Validation Rules

1. `rule_id` must match regex pattern: `^\d+\.\d+\.\d+$`
2. `confidence_level` must be in range [65, 95] and match tier guidelines:
   - 90-95%: Pattern-based detection
   - 75-85%: Structural analysis
   - 65-70%: Semantic analysis
3. `priority` distribution matches FUTU standards: 23 P0, 43 P1, 31 P2
4. `agent_assignment` must be one of: gorm, error-safety, naming-logging, organization

### State Transitions

**State**: Immutable (rules are defined in FUTU_GO_STANDARDS.md and don't change state)

Rules are static definitions that don't transition between states. Updates to rules require:
1. Version bump in FUTU_GO_STANDARDS.md
2. Test case updates
3. Agent SKILL.md algorithm updates

### Example

```json
{
  "rule_id": "1.3.7",
  "rule_name": "未使用的struct字段",
  "priority": "P0",
  "detection_algorithm_type": "semantic",
  "confidence_level": 65,
  "known_limitations": "May miss dynamically accessed fields (reflection, map lookups, third-party serialization)",
  "agent_assignment": "gorm",
  "example_violation": "type User struct {\n  ID int\n  Name string\n  UnusedField string  // Never referenced\n}",
  "example_fix": "type User struct {\n  ID int\n  Name string\n  // UnusedField removed\n}"
}
```

---

## Entity 2: Test Case File

**Description**: Represents a Go source file containing deliberate rule violations for testing detection capabilities.

### Fields

| Field Name | Type | Description | Constraints |
|------------|------|-------------|-------------|
| `filename` | string | Test file name | Pattern: {agent}_{priority}_issues.go, Required |
| `file_path` | string | Absolute path to test file | Full path, Required |
| `rule_category` | enum | Primary agent category | Values: gorm, error-safety, naming-logging, organization, quality, structure, Required |
| `priority` | enum | Primary priority level | Values: P0, P1, P2, Required |
| `expected_issue_count` | integer | Total violations expected | Range: 1-30, Required |
| `rule_ids` | array<string> | List of rule IDs violated | Each element matches rule_id format, Required |
| `description` | text | File purpose summary | Brief explanation, Required |
| `lines_of_code` | integer | Approximate file size | Range: 50-300, Optional |
| `last_validated` | datetime | Last successful validation run | ISO 8601 format, Optional |

### Relationships

- **Contains Detection Rules** (many-to-many): A test file deliberately violates multiple rules
- **Produces Validation Reports** (one-to-many): Each test run generates a report referencing this file

### Validation Rules

1. `filename` must match pattern: `^(gorm|error_safety|naming_logging|organization|quality|structure)_(p0|p1|p2)_issues\.go$`
2. `expected_issue_count` must equal length of `rule_ids` array
3. `rule_ids` must all belong to the specified `priority` level
4. `rule_category` must align with `rule_ids` agent assignments (majority)
5. File must exist at `file_path` location

### State Transitions

**States**:
- **Draft**: File created but not yet tested
- **Validated**: File tested and expected issues confirmed
- **Outdated**: Rule changes require test updates
- **Maintained**: Actively used in CI/CD validation

**Transitions**:
1. Draft → Validated: First successful test run with expected_issue_count match
2. Validated → Outdated: Rule addition/removal in FUTU_GO_STANDARDS.md
3. Outdated → Validated: Test file updated and re-validated
4. Validated → Maintained: Integrated into automated testing

### Example

```json
{
  "filename": "gorm_p0_issues.go",
  "file_path": "/data/jason/workspace/marketplace/test-cases/go-code-review/bad/gorm_p0_issues.go",
  "rule_category": "gorm",
  "priority": "P0",
  "expected_issue_count": 13,
  "rule_ids": ["1.3.1", "1.3.2", "1.3.3", "1.3.4", "1.3.5", "1.3.6", "1.3.7", "1.3.8", "1.3.9", "1.3.10", "1.3.11", "1.3.12", "1.3.13"],
  "description": "GORM database operations violations including missing WHERE/SELECT, incorrect field usage, and batch operation issues",
  "lines_of_code": 180,
  "last_validated": "2026-01-12T10:30:00Z"
}
```

---

## Entity 3: Validation Report

**Description**: Represents the output of running validate_results_v2.sh on code review results, containing per-rule detection metrics.

### Fields

| Field Name | Type | Description | Constraints |
|------------|------|-------------|-------------|
| `timestamp` | datetime | When validation ran | ISO 8601 format, Required |
| `detected_rules` | array<string> | Rule IDs successfully detected | Array of rule_id strings, Required |
| `undetected_rules` | array<string> | Rule IDs missed by detection | Array of rule_id strings, Required |
| `p0_detection_rate` | float | P0 detection percentage | Range: 0-100%, Required |
| `p1_detection_rate` | float | P1 detection percentage | Range: 0-100%, Required |
| `p2_detection_rate` | float | P2 detection percentage | Range: 0-100%, Required |
| `composite_score` | float | Weighted average detection rate | Formula: P0×0.5 + P1×0.3 + P2×0.2, Required |
| `false_positive_rate` | float | Incorrect flag percentage | Range: 0-100%, Required |
| `detailed_results` | array<object> | Per-rule detection details | Each: {rule_id, detected, confidence}, Required |
| `test_files_processed` | array<string> | List of test files validated | Array of filenames, Required |
| `validation_duration` | integer | Time taken in seconds | Range: 1-600, Optional |
| `meets_success_criteria` | boolean | Whether SC-001-010 targets met | true/false, Required |

### Relationships

- **References Test Case Files** (many-to-one): Report aggregates results from multiple test files
- **Validates Detection Rules** (many-to-many): Each report checks detection of many rules

### Validation Rules

1. `detected_rules` ∪ `undetected_rules` = complete rule set (97 rules)
2. `detected_rules` ∩ `undetected_rules` = ∅ (no overlaps)
3. `p0_detection_rate` calculated as: (P0 rules detected / 23) × 100
4. `p1_detection_rate` calculated as: (P1 rules detected / 43) × 100
5. `p2_detection_rate` calculated as: (P2 rules detected / 31) × 100
6. `composite_score` = (p0_detection_rate × 0.5) + (p1_detection_rate × 0.3) + (p2_detection_rate × 0.2)
7. `meets_success_criteria` = true if:
   - p0_detection_rate ≥ 95%
   - p1_detection_rate ≥ 92%
   - p2_detection_rate ≥ 88%
   - composite_score ≥ 93%
   - false_positive_rate ≤ 2%

### State Transitions

**States**:
- **Generated**: Report created from test run
- **Passed**: Meets all success criteria
- **Failed**: Below one or more targets
- **Archived**: Historical record

**Transitions**:
1. Generated → Passed: All success criteria met
2. Generated → Failed: One or more criteria not met
3. Passed/Failed → Archived: New test run creates new report

### Example

```json
{
  "timestamp": "2026-01-12T10:30:00Z",
  "detected_rules": ["1.3.1", "1.3.2", ..., "3.2.7"],
  "undetected_rules": ["1.3.13", "2.3.8", "3.1.5"],
  "p0_detection_rate": 95.7,
  "p1_detection_rate": 93.0,
  "p2_detection_rate": 90.3,
  "composite_score": 93.8,
  "false_positive_rate": 1.5,
  "detailed_results": [
    {
      "rule_id": "1.3.7",
      "rule_name": "未使用的struct字段",
      "priority": "P0",
      "test_file": "gorm_p0_issues.go",
      "detected": true,
      "confidence_claimed": 65,
      "detection_method": "field_usage_tracking"
    }
  ],
  "test_files_processed": ["gorm_p0_issues.go", "error_safety_p0_issues.go", ...],
  "validation_duration": 45,
  "meets_success_criteria": true
}
```

---

## Entity 4: Agent Enhancement

**Description**: Represents the configuration and metadata for detection algorithm improvements to a specific agent.

### Fields

| Field Name | Type | Description | Constraints |
|------------|------|-------------|-------------|
| `agent_name` | string | Name of specialized agent | Values: gorm, error-safety, naming-logging, organization, Required |
| `agent_skill_path` | string | Path to agent SKILL.md file | Absolute path, Required |
| `enhanced_rules` | array<string> | Rule IDs receiving algorithm improvements | Array of rule_id strings, Required |
| `detection_algorithms` | text | Markdown documentation of new algorithms | Detailed algorithm descriptions, Required |
| `multi_pass_strategy` | text | Description of multi-pass detection approach | Pass 1, Pass 2, etc. descriptions, Required |
| `confidence_matrix` | array<object> | Per-rule confidence mappings | Each: {rule_id, confidence, method, limitations}, Required |
| `known_limitations` | text | Overall agent detection limitations | Markdown format, Required |
| `pre_enhancement_coverage` | integer | Rules handled before enhancement | Count, Optional |
| `post_enhancement_coverage` | integer | Rules handled after enhancement | Count, Optional |
| `version` | string | Agent version identifier | Semantic version, Optional |

### Relationships

- **Enhances Detection Rules** (one-to-many): Each agent enhancement improves multiple rules
- **Belongs to Agent** (one-to-one): Each enhancement configuration corresponds to one agent

### Validation Rules

1. `agent_name` must match directory structure: `skills/go-code-review/{agent_name}/SKILL.md`
2. `enhanced_rules` must all be assigned to this `agent_name` per Detection Rule entity
3. `confidence_matrix` must have entry for each rule in `enhanced_rules`
4. `post_enhancement_coverage` ≥ `pre_enhancement_coverage`
5. `version` follows semantic versioning format

### State Transitions

**States**:
- **Planned**: Enhancement designed but not implemented
- **In Progress**: SKILL.md being updated
- **Implemented**: Code written, not yet tested
- **Validated**: Tested and meeting targets
- **Deployed**: Merged to main branch

**Transitions**:
1. Planned → In Progress: Implementation begins
2. In Progress → Implemented: SKILL.md updates complete
3. Implemented → Validated: Test suite passes with expected detection rates
4. Validated → Deployed: Merged and version tagged

### Example

```json
{
  "agent_name": "gorm",
  "agent_skill_path": "/data/jason/workspace/marketplace/skills/go-code-review/gorm-review/SKILL.md",
  "enhanced_rules": ["1.3.7", "1.3.8", "1.3.9", "1.3.10"],
  "detection_algorithms": "## Enhanced Detection Algorithms\n\n### Rule 1.3.7: Field Usage Tracking\n[Multi-pass algorithm...]\n\n### Rule 1.3.8: Field Ordering\n[Struct parsing algorithm...]\n\n...",
  "multi_pass_strategy": "Pass 1: Extract struct definitions\nPass 2: Scan Create/Update operations\nPass 3: Scan Select/Query operations\nPass 4: Apply exception rules\nPass 5: Report unused fields",
  "confidence_matrix": [
    {
      "rule_id": "1.3.7",
      "confidence": 65,
      "method": "field_usage_tracking",
      "limitations": "May miss reflection, dynamic map access"
    },
    {
      "rule_id": "1.3.8",
      "confidence": 90,
      "method": "struct_parsing",
      "limitations": "None significant"
    }
  ],
  "known_limitations": "- Cannot detect dynamically accessed fields via reflection\n- May not track fields used in third-party serialization\n- Complex embedded struct patterns may be missed",
  "pre_enhancement_coverage": 9,
  "post_enhancement_coverage": 13,
  "version": "3.0.0"
}
```

---

## Entity Relationships Diagram

```
┌─────────────────┐      assigned to     ┌──────────────────┐
│ Detection Rule  │────────────────────▶│ Agent Enhancement│
│  (97 instances) │                      │  (4 instances)   │
└────────┬────────┘                      └──────────────────┘
         │                                         │
         │ appears in                    enhances │
         │ (many-to-many)                         │
         │                                         │
         ▼                                         │
┌─────────────────┐      produces      ┌──────────▼──────────┐
│ Test Case File  │─────────────────▶│ Validation Report  │
│  (8 instances)  │                    │ (N instances)      │
└─────────────────┘                    └────────────────────┘
```

## Data Flow

1. **Design Phase**: Detection Rules defined with enhanced_rules list
2. **Enhancement Phase**: Agent Enhancements implemented with detection algorithms
3. **Test Creation Phase**: Test Case Files created with deliberate violations
4. **Validation Phase**: Code review executed, output saved
5. **Reporting Phase**: Validation Report generated by parsing output against Test Case File expectations
6. **Iteration Phase**: Validation Report metrics inform Agent Enhancement refinements

## Persistence Strategy

- **Detection Rules**: Documented in `skills/go-code-review/shared/FUTU_GO_STANDARDS.md` (Markdown)
- **Test Case Files**: Stored as `.go` files in `test-cases/go-code-review/bad/` and `good/`
- **Agent Enhancements**: Documented in `skills/go-code-review/{agent}/SKILL.md` (Markdown)
- **Validation Reports**: Generated as JSON by `validate_results_v2.sh`, can be stored in `test-cases/go-code-review/reports/`

No database required - all entities are file-based or in-memory during validation.

## Summary

This data model provides the conceptual framework for:
- Organizing 97 detection rules across 4 agents
- Structuring 8 test files with clear validation expectations
- Generating comprehensive validation reports with per-rule metrics
- Tracking agent enhancements with confidence levels and known limitations

All entities are technology-agnostic (no implementation details) and align with functional requirements FR-001 through FR-012.
