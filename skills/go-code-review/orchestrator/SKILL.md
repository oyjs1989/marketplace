---
name: Go Code Review
description: Orchestrates comprehensive Go code reviews based on FUTU's coding standards. ALWAYS use this when user asks to "review Go code", "check Go code quality", "review this PR", or "code review". This skill coordinates specialized review skills for different aspects of Go code.
---

# Go Code Review Orchestrator

## Purpose

This is the **main entry point** for Go code reviews. It orchestrates multiple specialized review skills to provide comprehensive code quality feedback.

## When to Use This Skill

Activate when users request:
- "Review Go code"
- "Check Go code quality"
- "Review this PR"
- "Code review"
- "Perform code review"

## Review Workflow

### Step 1: Get Code Changes

```bash
git diff master
```

Get all changes between the current branch and master branch.

### Step 2: Analyze Changes

Identify the scope of changes:
- Database models and GORM operations → Trigger **GORM Database Review**
- Error handling and concurrency → Trigger **Error & Safety Review**
- Variable/function naming and logging → Trigger **Naming & Logging Review**
- Code structure and interfaces → Trigger **Organization & Quality Review**

### Step 3: Priority Levels

**P0 - Must Fix (Critical Issues)**
- Missing or improper error handling
- Unchecked nil pointers
- Improper database operations
- Concurrency safety issues
- Business logic errors
- Unused struct fields

**P1 - Strongly Recommended (Code Quality)**
- Non-standard naming conventions
- Improper logging
- Missing necessary comments
- Poor function encapsulation
- Code duplication

**P2 - Suggested Optimization (Style & Best Practices)**
- Inconsistent formatting
- Readability improvements
- Performance optimization opportunities

### Step 4: Coordinate Reviews

For each changed file:

1. **Check for GORM/Database code**
   - Look for: `gorm:"column:"`, `db.Model()`, `db.Where()`, struct definitions with GORM tags
   - Action: Apply GORM Database Review rules

2. **Check for Error/Concurrency patterns**
   - Look for: `error` returns, `fmt.Errorf()`, `nil` checks, `go func()`, channels
   - Action: Apply Error & Safety Review rules

3. **Check Naming and Logging**
   - Look for: struct/function declarations, `log.Info()`, `log.Error()`, constants
   - Action: Apply Naming & Logging Review rules

4. **Check Code Organization**
   - Look for: interfaces, function signatures, project structure
   - Action: Apply Organization & Quality Review rules

### Step 5: Output Format

**IMPORTANT - 重要输出要求**:
- **必须使用中文输出所有审查结果**
- Only report code that has problems. Do NOT mention code that follows standards.
- 所有的问题描述(Problem)、建议(Suggestion)都必须用中文表达
- 文件路径、代码片段保持原样,不需要翻译

For each issue found:

```markdown
## 文件: <file_path>

### 问题 1 - [P0/P1/P2] 规则 X.Y.Z
**位置**: 第 N 行
**类别**: [GORM/错误处理/命名/日志/代码组织]
**原始代码**:
```go
// problematic code
```
**问题描述**: 用中文清晰描述问题
**修改建议**: 用中文说明如何修复,并提供代码示例
```

### Step 6: Save Results

Output all findings to `code_review.result` file with summary (用中文输出):

```markdown
# 代码审查总结

**问题总数**: X
- P0 (必须修复): X
- P1 (强烈建议): X
- P2 (建议优化): X

**问题分类**:
- GORM/数据库: X 个问题
- 错误处理: X 个问题
- 并发安全: X 个问题
- 命名规范: X 个问题
- 日志规范: X 个问题
- 代码组织: X 个问题

---

[详细问题列表如下...]
```

## Review Checklist

### Security & Correctness
- [ ] SQL injection risks (GORM queries)
- [ ] Nil pointer safety
- [ ] Concurrency safety (race conditions)
- [ ] Error handling completeness
- [ ] Business logic correctness

### Performance
- [ ] N+1 query issues
- [ ] Unused struct fields (memory waste)
- [ ] Missing capacity pre-allocation
- [ ] Inefficient GORM methods (First vs Take)

### Maintainability
- [ ] Clear naming conventions
- [ ] Proper logging
- [ ] Adequate comments
- [ ] No code duplication
- [ ] Reasonable function length

### Standards Compliance
- [ ] FUTU Go coding standards
- [ ] GORM best practices
- [ ] Error handling patterns
- [ ] Logging format

## Specialized Review Skills

This orchestrator delegates to:

1. **GORM Database Review** - Database layer code (rules 1.3.*)
2. **Go Error & Concurrency Safety** - Error handling and concurrency (rules 1.1.*, 1.2.*, 1.4.*)
3. **Go Naming & Logging Standards** - Naming and logging (rules 2.1.*, 2.2.*)
4. **Go Code Organization & Quality** - Structure and quality (rules 2.3.*, 2.4.*, 2.5.*, 3.*)

## Standards Reference

For complete rule details, refer to: `shared/FUTU_GO_STANDARDS.md`

## Best Practices

1. **Systematic**: Review all modified files thoroughly
2. **Prioritized**: Address P0 issues first
3. **Specific**: Include exact file paths and line numbers
4. **Actionable**: Provide concrete code suggestions
5. **Referenced**: Cite specific rule numbers
6. **Comprehensive**: Check all applicable categories

## Important Notes

- **所有输出必须使用中文描述** (All output MUST be in Chinese)
- Only report issues, not compliant code
- Include specific line numbers from actual files
- Prioritize by severity (P0 → P1 → P2)
- Provide clear fix suggestions with examples (用中文说明)
- Consider business context for exceptions
- Always save results to `code_review.result`
