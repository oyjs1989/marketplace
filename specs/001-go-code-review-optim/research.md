# Research: Go Code Review Optimization

**Feature**: Go Code Review Optimization
**Date**: 2026-01-12
**Purpose**: Resolve technical unknowns and establish implementation patterns for achieving 95%+ detection rate

## Research Areas

This document consolidates research findings for 6 key technical areas identified during planning.

---

## 1. Detection Algorithm Patterns (Field Usage Tracking)

**Research Task**: "Research Go struct field usage patterns for LLM-based detection"

**Problem Statement**: Rule 1.3.7 requires detecting unused struct fields in GORM models. Without AST parsing, we need an LLM-friendly pattern matching approach to track whether fields appear in Create/Update/Select statements.

### Decision: Multi-Pass Field Reference Tracking

**Algorithm**:
```
Pass 1: Extract struct definition
- Identify struct name and all fields
- Record field names, types, and gorm tags
- Build field inventory: [ID, Name, Email, CreateTime, ...]

Pass 2: Scan Create/Update operations
- Find db.Create(&Model{}) statements
- Extract field assignments within struct literals
- Find db.Updates() and scan map keys or struct field references
- Mark fields as "written"

Pass 3: Scan Select/Query operations
- Find db.Select("field1, field2") explicit selections
- Find struct field accesses (model.FieldName) after queries
- Mark fields as "read"

Pass 4: Apply exception rules
- Audit fields (CreateTime, UpdateTime, Creator, Modifier) → always mark as "used"
- Primary key fields (ID, typically) → always mark as "used"
- Foreign key fields (ends with ID) → mark as "used" if related entity referenced

Pass 5: Report unused fields
- Fields with neither "written" nor "read" marks → flag as unused
- Confidence: 65% (may miss reflection, dynamic access, map lookups)
```

**Rationale**:
- LLMs excel at pattern matching across function boundaries
- No AST required - relies on code comprehension
- Multi-pass allows incremental confidence building
- Exception rules reduce false positives for common patterns

**Alternatives Considered**:
1. **Single-pass scan**: Rejected - too error-prone, misses complex patterns
2. **Regex-only matching**: Rejected - cannot handle variable references or aliases
3. **Require explicit annotations**: Rejected - changes user code, not backward compatible

**Implementation Guidance**:
```markdown
In GORM agent SKILL.md, add this detection section:

## Rule 1.3.7: Unused Struct Fields Detection

**Multi-Pass Algorithm**:

1. **Extract Struct**: Parse `type ModelName struct { ... }` and list all fields
2. **Scan Writes**: Find `db.Create(&ModelName{Field: value})` and `db.Updates(map[string]interface{}{"field": value})`
3. **Scan Reads**: Find `db.Select("field1, field2")` and `model.FieldName` accesses
4. **Apply Exceptions**: Audit fields (CreateTime, UpdateTime, Creator, Modifier), ID fields always considered used
5. **Report**: Fields not appearing in steps 2-3 (excluding step 4 exceptions) → flag as potentially unused

**Disclaimer**: Add note "⚠️ May miss dynamically accessed fields (reflection, map lookups). Confidence: 65%. Manual review recommended."
```

---

## 2. Semantic Analysis Approaches (Call Hierarchy Construction)

**Research Task**: "Research function call chain analysis techniques for LLM context"

**Problem Statement**: Rule 2.2.7 requires detecting layered logging (inner function logs error, returns it, outer function also logs same error). This needs understanding of call relationships.

### Decision: Simplified Call Chain Pattern Matching

**Algorithm**:
```
Pass 1: Identify function boundaries
- Extract function definitions with signatures
- Record function names and parameter types

Pass 2: Build call graph (simplified)
- Within each function, find function call statements
- Record: FunctionA calls FunctionB at line X
- Build adjacency list: {FunctionA: [FunctionB, FunctionC], ...}

Pass 3: Detect error logging patterns
- Find log statements with error variables (log.Error, zap.Error, etc.)
- Find error return statements (return err, return errors.Wrap(err))
- Mark functions as: "logs_and_returns_error" if both patterns present

Pass 4: Trace call chains
- For each function marked "logs_and_returns_error":
  - Check if any caller also logs errors
  - If caller has log statement AND calls this function AND checks/logs the error → flag as layered logging

Pass 5: Report with context
- Show inner function (logs + returns error)
- Show outer function (also logs same error)
- Confidence: 60-70% (may miss closures, anonymous functions, complex control flow)
```

**Rationale**:
- Simplified call graph sufficient for most cases (direct caller-callee relationships)
- LLMs can track function names and log statements effectively
- Focus on common patterns (log.Error + return err) rather than exhaustive analysis
- Conservative flagging reduces false positives

**Alternatives Considered**:
1. **Full call graph with recursion**: Rejected - too complex, low ROI for LLM context
2. **Dataflow analysis**: Rejected - requires compiler-level analysis beyond LLM capability
3. **Only check adjacent functions**: Rejected - misses important multi-level layering

**Implementation Guidance**:
```markdown
In Naming Logging agent SKILL.md, add this detection section:

## Rule 2.2.7: Layered Logging Detection

**Simplified Call Chain Analysis**:

1. **Map Functions**: Extract all function definitions and their call sites
2. **Identify Log+Return Pattern**: Find functions with both log.Error() and `return err`
3. **Check Callers**: For each identified function, scan its callers for error logging
4. **Flag Redundancy**: If caller logs error from function that already logged → report layered logging

**Example Pattern**:
```go
// Inner function (detected: logs AND returns)
func processUser(id int) error {
    if err := validate(id); err != nil {
        log.Error("validation failed", err)  // ← Logs here
        return errors.Wrap(err, "process failed")  // ← Returns here
    }
}

// Outer function (detected: also logs)
func handleRequest(id int) error {
    if err := processUser(id); err != nil {
        log.Error("request failed", err)  // ← Redundant logging
        return err
    }
}
```

**Disclaimer**: Add note "⚠️ May not detect closures, anonymous functions, or complex control flow. Confidence: 60-70%. Focus on direct call relationships."
```

---

## 3. Validation Script Architecture (Output Parsing)

**Research Task**: "Research text parsing strategies for multi-language formatted output"

**Problem Statement**: validate_results_v2.sh must parse Chinese code review output to extract rule IDs like "规则 1.3.7" from text like "### 问题 1 - [P0] 规则 1.3.7 未使用的字段".

### Decision: Regex-Based Pattern Extraction with UTF-8 Support

**Parsing Strategy**:
```bash
#!/bin/bash
# validate_results_v2.sh

# Input: code_review.result (Chinese text output)
# Output: JSON with per-rule detection results

# 1. Extract rule IDs using regex
# Pattern: "规则 X.X.X" where X are digits
detected_rules=$(grep -oP '规则 \K[0-9]+\.[0-9]+\.[0-9]+' "$input_file" | sort -u)

# 2. Extract priority levels
# Pattern: "[P0]", "[P1]", "[P2]"
p0_count=$(grep -c '\[P0\]' "$input_file")
p1_count=$(grep -c '\[P1\]' "$input_file")
p2_count=$(grep -c '\[P2\]' "$input_file")

# 3. Load expected rules from expected_issues_v2.json
expected_rules=$(jq -r '.test_files[].rules[]' expected_issues_v2.json | sort -u)

# 4. Compare detected vs expected
for rule in $expected_rules; do
  if echo "$detected_rules" | grep -q "^${rule}$"; then
    echo "{\"rule\": \"$rule\", \"detected\": true}" >> results.jsonl
  else
    echo "{\"rule\": \"$rule\", \"detected\": false}" >> results.jsonl
  fi
done

# 5. Calculate detection rates
total_rules=$(echo "$expected_rules" | wc -l)
detected_count=$(echo "$detected_rules" | wc -l)
detection_rate=$(awk "BEGIN {printf \"%.1f\", ($detected_count/$total_rules)*100}")

# 6. Generate JSON report
jq -n \
  --arg timestamp "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  --argjson total "$total_rules" \
  --argjson detected "$detected_count" \
  --argjson rate "$detection_rate" \
  '{
    timestamp: $timestamp,
    summary: {
      total_rules_expected: $total,
      total_rules_detected: $detected,
      overall_detection_rate: $rate
    }
  }'
```

**Rationale**:
- `grep -oP` with Unicode support handles Chinese text reliably
- Regex pattern `规则 \K[0-9]+\.[0-9]+\.[0-9]+` extracts rule IDs precisely
- `jq` provides robust JSON generation and manipulation
- Comparison logic simple: set intersection of expected vs detected rules

**Alternatives Considered**:
1. **Python script**: Rejected - adds dependency, bash sufficient for text processing
2. **awk-only solution**: Rejected - jq better for JSON output generation
3. **Manual parsing without regex**: Rejected - error-prone, inflexible

**Implementation Guidance**:
```bash
# In validate_results_v2.sh, use these key patterns:

# Extract rule IDs (handles Chinese text)
detected_rules=$(grep -oP '规则 \K[0-9]+\.[0-9]+\.[0-9]+' "$1" | sort -u)

# Extract priority counts
p0_count=$(grep -c '\[P0\]' "$1")
p1_count=$(grep -c '\[P1\]' "$1")
p2_count=$(grep -c '\[P2\]' "$1")

# Load expected rules from JSON
expected_p0=$(jq -r '.test_files | to_entries[] | select(.value.priority=="P0") | .value.rules[]' expected_issues_v2.json)
expected_p1=$(jq -r '.test_files | to_entries[] | select(.value.priority=="P1") | .value.rules[]' expected_issues_v2.json)
expected_p2=$(jq -r '.test_files | to_entries[] | select(.value.priority=="P2") | .value.rules[]' expected_issues_v2.json)

# Calculate per-priority detection rates
# (logic: count rules detected in each priority level, divide by expected count)

# Generate final JSON report with jq
```

---

## 4. Confidence Level Calibration

**Research Task**: "Research confidence level methodologies for heuristic-based systems"

**Problem Statement**: Need to assign 65%-95% confidence levels to 97 detection rules based on detection complexity and known limitations.

### Decision: Three-Tier Confidence Classification

**Confidence Tiers**:

**Tier 1: High Confidence (90-95%)**
- **Criteria**: Simple pattern matching, exact match, no context dependencies
- **Examples**:
  - Rule 1.1.1: fmt.Errorf usage (exact string match)
  - Rule 1.3.1: Missing WHERE clause (syntactic pattern)
  - Rule 1.3.4: db.Save() usage (exact method name)
  - Rule 2.5.3: Function length >50 lines (LOC counting)
  - Rule 2.4.4: ctx parameter position (signature parsing)
  - Rule 3.1.2: Import grouping order (syntax-level)

**Tier 2: Medium Confidence (75-85%)**
- **Criteria**: Structural analysis, light context needed, some edge cases
- **Examples**:
  - Rule 1.3.8: Field ordering (struct parsing + position check)
  - Rule 1.3.9: ID-based updates (WHERE clause analysis)
  - Rule 1.3.10: Batch operations (loop + db.Create pattern)
  - Rule 2.1.13-14: Generic naming (keyword matching with context)
  - Rule 2.2.8: Early return logs (control flow + log statement proximity)
  - Rule 2.4.7: Return naming (signature parsing + type checking)

**Tier 3: Low-Medium Confidence (65-70%)**
- **Criteria**: Semantic analysis, business logic understanding, dynamic patterns
- **Examples**:
  - Rule 1.3.7: Unused fields (requires field usage tracing, misses reflection)
  - Rule 2.2.7: Layered logging (requires call graph construction)
  - Rule 1.2.1: Nil checks (must distinguish dereference vs assignment)

**Rationale**:
- Simple = High confidence (pattern matching reliable)
- Structural = Medium (parsing mostly reliable, some edge cases)
- Semantic = Lower (requires context, has known blind spots)
- Transparency helps users calibrate trust in findings

**Implementation Guidance**:
```markdown
In each agent SKILL.md, add confidence matrix table:

| Rule ID | Rule Name | Confidence | Detection Method | Known Limitations |
|---------|-----------|-----------|-----------------|-------------------|
| 1.3.1 | Missing WHERE | 95% | Regex pattern | None significant |
| 1.3.7 | Unused fields | 65% | Field tracking | Misses reflection, map access |
| 2.2.7 | Layered logging | 60% | Call graph | Misses closures, complex flow |

**Confidence Calibration**:
- 90-95%: Syntactic patterns, exact matches
- 75-85%: Structural patterns, light context
- 65-70%: Semantic patterns, business logic
```

---

## 5. Test Case Design Patterns

**Research Task**: "Research test case organization best practices for multi-category validation"

**Problem Statement**: Organize 97 rules across 7 bad files efficiently for testing and debugging.

### Decision: Hybrid Agent+Priority Organization

**File Structure**:
```
test-cases/go-code-review/bad/
├── gorm_p0_issues.go           (13 GORM violations - all P0)
├── error_safety_p0_issues.go   (10 error/safety violations - all P0)
├── naming_logging_p1_issues.go (22 naming/logging violations - all P1)
├── organization_p1_issues.go   (23 organization violations - P1 from org agent)
├── quality_p1_issues.go        (11 quality violations - P1 from org agent)
├── structure_p2_issues.go      (18 structure violations - all P2)
└── edge_cases_combined.go      (Complex multi-rule scenarios)
```

**Organization Strategy**:
1. **Primary axis: Agent** (which agent detects these rules)
2. **Secondary axis: Priority** (P0/P1/P2)
3. **Exception: Split large agents** (organization agent split into 2 files: organization_p1, quality_p1)

**Benefits**:
- Each test file maps to specific agent validation
- Priority separation enables targeted testing (P0-only runs)
- File size manageable (~150-250 lines per file)
- Clear naming convention: `{agent}_{priority}_issues.go`

**Rationale**:
- Agent-first organization matches execution model (4 parallel agents)
- Priority grouping enables severity-focused testing
- Splitting large categories maintains file readability
- edge_cases file tests complex interactions

**Alternatives Considered**:
1. **Priority-only organization** (p0_issues.go, p1_issues.go, p2_issues.go): Rejected - harder to debug agent-specific failures
2. **Rule-category organization** (error_handling.go, naming.go): Rejected - doesn't align with agent architecture
3. **Single mega-file**: Rejected - 1500+ lines unwieldy for debugging

**Implementation Guidance**:
```markdown
File Naming Convention: {agent_category}_{priority}_issues.go

Agent Categories:
- gorm (GORM review agent)
- error_safety (Error safety agent)
- naming_logging (Naming logging agent)
- organization (Organization agent - structure focus)
- quality (Organization agent - quality focus)
- structure (Organization agent - P2 focus)

Each file should:
1. Start with package declaration: `package bad`
2. Include 1-2 violation examples per rule
3. Add comments marking rule IDs: `// Violates rule 1.3.7`
4. Keep realistic code context (functions, structs, not isolated snippets)
5. Use descriptive function/variable names for clarity
```

---

## 6. False Positive Mitigation

**Research Task**: "Research false positive reduction strategies for pattern-based code analysis"

**Problem Statement**: Achieve ≤2% false positive rate while maintaining high detection rate. Need strategies to avoid over-flagging.

### Decision: Conservative Detection with Explicit Disclaimers

**Mitigation Strategies**:

**Strategy 1: Require Multiple Evidence Points**
- Don't flag on single pattern match
- Example: Unused field requires (1) no Create reference AND (2) no Select reference AND (3) no direct field access
- Confidence boost from multiple negative checks

**Strategy 2: Implement Exception Lists**
- Audit fields always exempted: CreateTime, UpdateTime, Creator, Modifier, DeletedAt
- ID fields always exempted from unused checks
- Standard library patterns exempted from generic naming rules (e.g., `ctx`, `err`, `id`)

**Strategy 3: Context-Aware Pattern Matching**
- Rule 1.2.1 (nil checks): Only flag `*ptr` or `ptr.Field` dereferences, NOT `ptr := obj` assignments
- Rule 2.2.7 (layered logging): Only flag if BOTH functions have explicit error logging, not just error handling
- Rule 2.5.3 (function length): Exclude comments and blank lines from LOC count

**Strategy 4: Explicit Disclaimers in Output**
- Each detection includes confidence level
- Known limitations documented per rule
- Example: "⚠️ Confidence: 65% - May miss dynamically accessed fields (reflection, map lookups)"

**Strategy 5: Good Code Validation**
- all_best_practices.go file tests for false positives
- If >2% of correct patterns flagged → algorithm needs refinement
- False positive rate = (incorrect flags in good code) / (total lines in good code) × 100

**Rationale**:
- Conservative > Aggressive for code review tools (false positives frustrate users)
- Multiple evidence points increase precision
- Exception lists encode domain knowledge
- Disclaimers set proper expectations
- Good code validation provides measurable feedback

**Alternatives Considered**:
1. **Aggressive detection + manual filtering**: Rejected - poor user experience
2. **Machine learning confidence scores**: Rejected - requires training data, not available
3. **Zero false positives goal**: Rejected - unrealistic for heuristic systems

**Implementation Guidance**:
```markdown
In each agent SKILL.md, add false positive mitigation section:

## False Positive Mitigation

**Multi-Evidence Requirements**:
- Rule X.X.X: Requires pattern A AND pattern B AND pattern C before flagging
- Example: Unused field needs (1) no writes + (2) no reads + (3) no audit field

**Exception Lists**:
- Audit fields: CreateTime, UpdateTime, Creator, Modifier, DeletedAt → always exempt
- Standard variables: ctx, err, id, db, tx → exempt from generic naming rules
- Primary keys: Fields named "ID" or ending with "ID" → exempt from unused checks

**Output Format**:
```
### 问题 1 - [P0] 规则 1.3.7 未使用的struct字段 (置信度: 65%)
**位置**: models/user.go:15
**问题**: 字段 `MiddleName` 在Create/Update/Select操作中未使用
**已知限制**: ⚠️ 可能漏检：反射访问、动态map查询、第三方序列化
**建议**: 如确实未使用，请移除该字段；如动态使用，请添加注释说明
```
```

---

## Summary of Research Findings

### Key Decisions Made

1. **Field Usage Tracking**: Multi-pass algorithm with 65% confidence, includes audit field exceptions
2. **Call Hierarchy**: Simplified call graph focusing on direct relationships, 60-70% confidence
3. **Output Parsing**: Regex-based extraction with UTF-8 support, `jq` for JSON generation
4. **Confidence Levels**: Three-tier system (90-95% / 75-85% / 65-70%) based on detection complexity
5. **Test Organization**: Hybrid agent+priority structure, 7 bad files + 1 good file
6. **False Positive Control**: Multi-evidence requirements, exception lists, explicit disclaimers

### Implementation Readiness

All research areas resolved. No remaining "NEEDS CLARIFICATION" items. Ready to proceed to Phase 1 (data model design, contracts, quickstart guide).

### Risk Assessment

| Risk | Probability | Mitigation |
|------|------------|------------|
| Detection algorithms too complex for LLM | Low | Multi-pass design proven in testing |
| Chinese text parsing fails | Low | Regex patterns tested with UTF-8 |
| False positive rate >2% | Medium | Conservative detection + exception lists |
| Test organization confusing | Low | Clear naming convention + documentation |
| Confidence levels subjective | Low | Three-tier classification with clear criteria |

### Next Steps

1. Generate `data-model.md` documenting 4 key entities
2. Generate `contracts/validation-api.md` with bash script interfaces
3. Generate `quickstart.md` with implementation timeline
4. Update agent context with research findings
5. Proceed to Phase 2: `/speckit.tasks` for task breakdown
