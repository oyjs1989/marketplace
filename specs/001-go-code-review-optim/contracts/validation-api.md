# API Contract: Validation Framework

**Feature**: Go Code Review Optimization
**Date**: 2026-01-12
**Purpose**: Define interfaces and contracts for the validation infrastructure

## Overview

This document specifies the contracts for:
1. `validate_results_v2.sh` - Main validation script
2. `expected_issues_v2.json` - Test expectations configuration
3. Code review output format (input to validation)

---

## Contract 1: validate_results_v2.sh

### Purpose
Parse code review output, compare against expected issues, calculate detection rates, and generate comprehensive JSON report.

### Command Line Interface

```bash
validate_results_v2.sh [OPTIONS] <code_review_result_file>
```

### Arguments

| Argument | Type | Required | Description |
|----------|------|----------|-------------|
| `code_review_result_file` | path | Yes | Path to code review output file (Chinese text format) |

### Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `--config` | path | `./expected_issues_v2.json` | Path to expectations config file |
| `--output` | path | stdout | Path to write JSON report (default: stdout) |
| `--verbose` | flag | false | Print detailed progress messages to stderr |
| `--strict` | flag | false | Exit with error if any target not met |

### Input File Format

**File**: `code_review_result_file` (Chinese text output from Go code review skill)

**Expected Patterns**:
```markdown
## 文件: path/to/file.go

### 问题 1 - [P0] 规则 1.3.7 未使用的字段
**位置**: 第 42 行
**原始代码**:
```go
type User struct {
    ID int
    UnusedField string  // This field
}
```
**问题描述**: 字段 `UnusedField` 在Create/Update/Select操作中未使用
**修改建议**: 如确实未使用，请移除该字段

### 问题 2 - [P1] 规则 2.2.7 分层日志
...
```

**Required Pattern Elements**:
- Priority marker: `[P0]`, `[P1]`, or `[P2]`
- Rule ID: `规则 X.X.X` where X are digits
- Problem number: `问题 N` where N is integer

### Output Format

**JSON Structure**:
```json
{
  "validation_metadata": {
    "timestamp": "2026-01-12T10:30:00Z",
    "validator_version": "3.0.0",
    "input_file": "/path/to/code_review.result",
    "config_file": "/path/to/expected_issues_v2.json",
    "validation_duration_seconds": 3.2
  },
  "summary": {
    "total_rules_expected": 97,
    "total_rules_detected": 91,
    "overall_detection_rate": 93.8,
    "meets_all_targets": true,
    "p0_detection": {
      "expected": 23,
      "detected": 22,
      "rate": 95.7,
      "meets_target": true,
      "target": 95.0
    },
    "p1_detection": {
      "expected": 43,
      "detected": 40,
      "rate": 93.0,
      "meets_target": true,
      "target": 92.0
    },
    "p2_detection": {
      "expected": 31,
      "detected": 29,
      "rate": 93.5,
      "meets_target": true,
      "target": 88.0
    },
    "composite_score": 94.2,
    "composite_target": 93.0,
    "false_positive_rate": 1.5,
    "false_positive_target": 2.0
  },
  "per_file_results": [
    {
      "file": "gorm_p0_issues.go",
      "priority": "P0",
      "expected_rules": ["1.3.1", "1.3.2", "1.3.3", "1.3.4", "1.3.5", "1.3.6", "1.3.7", "1.3.8", "1.3.9", "1.3.10", "1.3.11", "1.3.12", "1.3.13"],
      "detected_rules": ["1.3.1", "1.3.2", "1.3.3", "1.3.4", "1.3.5", "1.3.6", "1.3.7", "1.3.8", "1.3.9", "1.3.10", "1.3.11", "1.3.12"],
      "undetected_rules": ["1.3.13"],
      "expected_count": 13,
      "detected_count": 12,
      "detection_rate": 92.3,
      "meets_target": false
    }
    // ... more files
  ],
  "per_rule_results": [
    {
      "rule_id": "1.3.7",
      "rule_name": "未使用的struct字段",
      "priority": "P0",
      "test_file": "gorm_p0_issues.go",
      "detected": true,
      "confidence_claimed": 65,
      "detection_method": "field_usage_tracking",
      "notes": null
    }
    // ... all 97 rules
  ],
  "undetected_rules_detail": [
    {
      "rule_id": "1.3.13",
      "rule_name": "Rule name from expected_issues_v2.json",
      "priority": "P0",
      "test_file": "gorm_p0_issues.go",
      "expected_but_not_found": true,
      "possible_reasons": ["Detection algorithm needs refinement", "Test case issue"]
    }
  ],
  "false_positives": [
    {
      "issue_description": "问题 1 - [P1] 规则 2.1.5 命名规范",
      "file": "good/all_best_practices.go",
      "line": 42,
      "reason": "Flagged correct code as violation"
    }
  ]
}
```

### Exit Codes

| Code | Meaning | Description |
|------|---------|-------------|
| 0 | Success | All detection rate targets met |
| 1 | Target Not Met | One or more detection rate targets not met |
| 2 | Invalid Input | Input file not found or invalid format |
| 3 | Invalid Config | expected_issues_v2.json not found or malformed |
| 4 | Parsing Error | Could not extract rule IDs from input file |

### Error Handling

**Missing Input File**:
```bash
$ validate_results_v2.sh missing.file
Error: Input file 'missing.file' not found
Exit code: 2
```

**Invalid JSON Config**:
```bash
$ validate_results_v2.sh code_review.result --config invalid.json
Error: Config file 'invalid.json' is not valid JSON
Exit code: 3
```

**No Rules Detected**:
```bash
$ validate_results_v2.sh empty.file
Error: No rule IDs detected in input file. Expected pattern: '规则 X.X.X'
Exit code: 4
```

### Performance Requirements

- **Execution Time**: ≤5 seconds for full 97-rule validation
- **Memory Usage**: ≤50MB peak
- **File Size Limit**: Input file up to 10MB

### Dependencies

- **Required**: bash 4.0+, grep, awk, jq 1.6+
- **Optional**: None

### Example Usage

```bash
# Basic usage (output to stdout)
./validate_results_v2.sh code_review.result

# Save report to file
./validate_results_v2.sh code_review.result --output report.json

# Strict mode (exit 1 if targets not met)
./validate_results_v2.sh code_review.result --strict

# Custom config location
./validate_results_v2.sh code_review.result --config ../custom_expected.json

# Verbose mode with progress
./validate_results_v2.sh code_review.result --verbose 2> validation.log
```

---

## Contract 2: expected_issues_v2.json

### Purpose
Configuration file defining expected rule violations per test file and success criteria targets.

### File Format

**JSON Schema**:
```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["version", "test_files", "rules_coverage", "success_criteria"],
  "properties": {
    "version": {
      "type": "string",
      "pattern": "^[0-9]+\\.[0-9]+\\.[0-9]+$",
      "description": "Semantic version matching main SKILL.md version"
    },
    "timestamp": {
      "type": "string",
      "format": "date",
      "description": "Last update date (YYYY-MM-DD)"
    },
    "test_files": {
      "type": "object",
      "patternProperties": {
        "^[a-z_]+\\.go$": {
          "type": "object",
          "required": ["priority", "expected_count", "rules"],
          "properties": {
            "priority": {
              "type": "string",
              "enum": ["P0", "P1", "P2"]
            },
            "expected_count": {
              "type": "integer",
              "minimum": 1
            },
            "rules": {
              "type": "array",
              "items": {
                "type": "string",
                "pattern": "^[0-9]+\\.[0-9]+\\.[0-9]+$"
              },
              "minItems": 1
            },
            "description": {
              "type": "string"
            }
          }
        }
      }
    },
    "rules_coverage": {
      "type": "object",
      "required": ["total_rules", "covered_rules", "coverage_percentage"],
      "properties": {
        "total_rules": {
          "type": "integer",
          "const": 97
        },
        "covered_rules": {
          "type": "integer",
          "const": 97
        },
        "coverage_percentage": {
          "type": "number",
          "const": 100.0
        },
        "p0_coverage": {
          "type": "integer",
          "const": 23
        },
        "p1_coverage": {
          "type": "integer",
          "const": 43
        },
        "p2_coverage": {
          "type": "integer",
          "const": 31
        }
      }
    },
    "success_criteria": {
      "type": "object",
      "required": ["p0_target", "p1_target", "p2_target", "composite_target", "false_positive_max"],
      "properties": {
        "p0_target": {
          "type": "number",
          "const": 95.0
        },
        "p1_target": {
          "type": "number",
          "const": 92.0
        },
        "p2_target": {
          "type": "number",
          "const": 88.0
        },
        "composite_target": {
          "type": "number",
          "const": 93.0
        },
        "false_positive_max": {
          "type": "number",
          "const": 2.0
        }
      }
    }
  }
}
```

### Example Content

```json
{
  "version": "3.0.0",
  "timestamp": "2026-01-12",
  "test_files": {
    "gorm_p0_issues.go": {
      "priority": "P0",
      "expected_count": 13,
      "rules": [
        "1.3.1", "1.3.2", "1.3.3", "1.3.4", "1.3.5", "1.3.6",
        "1.3.7", "1.3.8", "1.3.9", "1.3.10", "1.3.11", "1.3.12", "1.3.13"
      ],
      "description": "GORM database operations violations including missing WHERE/SELECT, incorrect field usage, and batch operation issues"
    },
    "error_safety_p0_issues.go": {
      "priority": "P0",
      "expected_count": 10,
      "rules": [
        "1.1.1", "1.1.2", "1.1.3", "1.1.4", "1.2.1",
        "1.4.1", "1.4.2", "1.5.1", "1.5.2"
      ],
      "description": "Error handling and safety violations including error wrapping, nil checks, concurrency patterns, and JSON processing"
    },
    "naming_logging_p1_issues.go": {
      "priority": "P1",
      "expected_count": 22,
      "rules": [
        "2.1.1", "2.1.2", "2.1.3", "2.1.4", "2.1.5", "2.1.6", "2.1.7",
        "2.1.8", "2.1.9", "2.1.10", "2.1.11", "2.1.12", "2.1.13", "2.1.14",
        "2.2.1", "2.2.2", "2.2.3", "2.2.4", "2.2.5", "2.2.6", "2.2.7", "2.2.8"
      ],
      "description": "Naming conventions and logging standards violations"
    },
    "organization_p1_issues.go": {
      "priority": "P1",
      "expected_count": 16,
      "rules": [
        "2.3.1", "2.3.2", "2.3.3", "2.3.4", "2.3.5", "2.3.6", "2.3.7",
        "2.3.8", "2.3.9", "2.3.10", "2.3.11", "2.3.12", "2.3.13", "2.3.14",
        "2.3.15", "2.3.16"
      ],
      "description": "Code organization violations including global variables, initialization, struct design"
    },
    "quality_p1_issues.go": {
      "priority": "P1",
      "expected_count": 11,
      "rules": [
        "2.4.1", "2.4.2", "2.4.3", "2.4.4", "2.4.5", "2.4.6",
        "2.4.7", "2.4.8", "2.4.9", "2.4.10",
        "2.5.1"
      ],
      "description": "Interface design and code quality violations"
    },
    "structure_p2_issues.go": {
      "priority": "P2",
      "expected_count": 18,
      "rules": [
        "2.5.2", "2.5.3", "2.5.4", "2.5.5", "2.5.6", "2.5.7", "2.5.8",
        "2.5.9", "2.5.10", "2.5.11",
        "3.1.1", "3.1.2", "3.1.3", "3.1.4", "3.1.5", "3.1.6",
        "3.2.1", "3.2.2"
      ],
      "description": "Project structure and testing standards violations"
    },
    "edge_cases_combined.go": {
      "priority": "Mixed",
      "expected_count": 7,
      "rules": [
        "1.3.7", "2.2.7", "2.5.3", "1.2.1", "3.1.2", "2.4.4", "1.3.10"
      ],
      "description": "Complex multi-rule scenarios testing edge cases and rule interactions"
    }
  },
  "rules_coverage": {
    "total_rules": 97,
    "covered_rules": 97,
    "coverage_percentage": 100.0,
    "p0_coverage": 23,
    "p1_coverage": 43,
    "p2_coverage": 31
  },
  "success_criteria": {
    "p0_target": 95.0,
    "p1_target": 92.0,
    "p2_target": 88.0,
    "composite_target": 93.0,
    "false_positive_max": 2.0
  }
}
```

### Validation Rules

1. **Completeness**: All 97 rules must appear in at least one test file
2. **No Duplicates**: Each rule should appear only once across all test files (except edge_cases)
3. **Priority Consistency**: Rules in a P0 file must all be P0 rules (per FUTU_GO_STANDARDS.md)
4. **Count Accuracy**: `expected_count` must equal length of `rules` array for each file
5. **Version Match**: `version` should match main SKILL.md version

---

## Contract 3: Code Review Output Format

### Purpose
Define the expected format of code review output generated by the Go code review skill (input to validation).

### File Format

**Language**: Chinese (中文) per FUTU standards
**Structure**: Markdown with specific heading patterns

### Required Elements

**File Header**:
```markdown
## 文件: path/to/file.go
```

**Issue Block**:
```markdown
### 问题 N - [PRIORITY] 规则 X.X.X RULE_NAME
**位置**: 第 LINE 行
**原始代码**:
```go
// Code snippet
```
**问题描述**: [Chinese description of the issue]
**修改建议**: [Chinese suggestion for fix]
```

### Pattern Specifications

| Element | Pattern | Example |
|---------|---------|---------|
| File header | `## 文件: <filepath>` | `## 文件: models/user.go` |
| Problem header | `### 问题 <N> - [<PX>] 规则 <X.X.X> <name>` | `### 问题 1 - [P0] 规则 1.3.7 未使用的字段` |
| Location | `**位置**: 第 <N> 行` | `**位置**: 第 42 行` |
| Code block | ` ```go\n<code>\n``` ` | See example above |
| Description | `**问题描述**: <text>` | `**问题描述**: 字段未使用` |
| Suggestion | `**修改建议**: <text>` | `**修改建议**: 移除该字段` |

### Mandatory Fields

- Priority marker: `[P0]`, `[P1]`, or `[P2]` (required for categorization)
- Rule ID: `规则 X.X.X` (required for detection tracking)
- Problem number: `问题 N` (required for counting)

### Optional Fields

- Confidence level: `(置信度: 65%)` (recommended for semantic rules)
- Known limitations: `**已知限制**: <text>` (recommended for low-confidence rules)

### Example Complete Output

```markdown
# Go 代码审查结果

## 文件: models/user.go

### 问题 1 - [P0] 规则 1.3.7 未使用的struct字段 (置信度: 65%)
**位置**: 第 15 行
**原始代码**:
```go
type User struct {
    ID          int
    Name        string
    UnusedField string  // This field is never used
    CreateTime  time.Time
}
```
**问题描述**: 字段 `UnusedField` 在Create/Update/Select操作中未使用
**已知限制**: ⚠️ 可能漏检：反射访问、动态map查询、第三方序列化
**修改建议**: 如确实未使用，请移除该字段；如动态使用，请添加注释说明

### 问题 2 - [P1] 规则 2.2.7 分层日志记录 (置信度: 60%)
**位置**: 第 45 行和第 52 行
**原始代码**:
```go
// Inner function
func processUser(id int) error {
    if err := validate(id); err != nil {
        log.Error("validation failed", err)
        return errors.Wrap(err, "process failed")
    }
}

// Outer function
func handleRequest(id int) error {
    if err := processUser(id); err != nil {
        log.Error("request failed", err)  // Redundant
        return err
    }
}
```
**问题描述**: 内层函数 `processUser` 和外层函数 `handleRequest` 都记录了同一错误，造成日志冗余
**已知限制**: ⚠️ 可能漏检：闭包、匿名函数、复杂控制流
**修改建议**: 移除外层函数的日志记录，或者内层函数只返回错误不记录日志

## 文件: services/order.go

### 问题 3 - [P0] 规则 1.3.10 循环中的批量操作
...
```

### Validation by validate_results_v2.sh

The validation script will:
1. Extract all `规则 X.X.X` patterns → detected_rules list
2. Count `[P0]`, `[P1]`, `[P2]` occurrences → priority counts
3. Compare detected_rules against expected_issues_v2.json → detection rates
4. Generate JSON report with per-rule results

---

## Integration Flow

```
┌──────────────────────┐
│  Go Code Review Skill│
│  (generates output)  │
└──────────┬───────────┘
           │ Produces
           ▼
┌──────────────────────┐
│ code_review.result   │◀─────┐
│ (Chinese text)       │      │
└──────────┬───────────┘      │
           │ Input to         │
           ▼                  │ References
┌──────────────────────┐      │
│validate_results_v2.sh│      │
│  (bash script)       │      │
└──────────┬───────────┘      │
           │ Reads            │
           ▼                  │
┌──────────────────────┐      │
│expected_issues_v2.json│──────┘
│  (configuration)     │
└──────────┬───────────┘
           │ Produces
           ▼
┌──────────────────────┐
│  validation_report   │
│  (JSON output)       │
└──────────────────────┘
```

## Summary

This API contract provides:
- **Clear interfaces** for validation script invocation and output
- **Structured configuration** for test expectations and success criteria
- **Defined output format** for code review results (enables reliable parsing)
- **Comprehensive error handling** with specific exit codes
- **Performance guarantees** (≤5 seconds, ≤50MB memory)

All contracts support FR-006 (validation config) and FR-007 (validation script) from the feature specification.
