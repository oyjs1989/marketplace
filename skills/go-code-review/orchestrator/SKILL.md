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

### Step 4: Parallel Agent Execution

**重要**: 为了提高审查效率，对每个变更文件**并行启动**4个独立的审查 agent，每个 agent 只负责检查自己的规则类别。

#### 4.1 并行启动审查 Agents

对于每个变更的文件，**同时启动**以下 4 个独立的 agent：

1. **GORM Database Review Agent** (`agents/go-code-review/gorm-agent.md`)
   - **职责范围**: 仅检查 GORM 数据库相关规则 (规则 1.3.*)
   - **检查内容**: `gorm:"column:"`, `db.Model()`, `db.Where()`, `db.Find()`, `db.Create()`, struct definitions with GORM tags
   - **输出**: 该文件中的 GORM/数据库相关问题列表

2. **Error & Safety Review Agent** (`agents/go-code-review/error-safety-agent.md`)
   - **职责范围**: 仅检查错误处理和并发安全规则 (规则 1.1.*, 1.2.*, 1.4.*, 1.5.*)
   - **检查内容**: `error` returns, `fmt.Errorf()`, `errors.Wrap()`, `nil` checks, `go func()`, channels, `json.Marshal()`
   - **输出**: 该文件中的错误处理/并发安全相关问题列表

3. **Naming & Logging Review Agent** (`agents/go-code-review/naming-logging-agent.md`)
   - **职责范围**: 仅检查命名规范和日志标准 (规则 2.1.*, 2.2.*)
   - **检查内容**: struct/function declarations, `log.Info()`, `log.Error()`, constants, variable names
   - **输出**: 该文件中的命名/日志相关问题列表

4. **Organization & Quality Review Agent** (`agents/go-code-review/organization-agent.md`)
   - **职责范围**: 仅检查代码组织和质量规则 (规则 2.3.*, 2.4.*, 2.5.*, 3.*)
   - **检查内容**: interfaces, function signatures, project structure, code organization patterns
   - **输出**: 该文件中的代码组织/质量问题列表

#### 4.2 并行执行指令

**执行方式**:
```
对于文件 <file_path>:
  1. 同时调用 4 个 agent，每个 agent 接收相同的文件内容作为输入
  2. 每个 agent 独立执行，只检查自己的规则范围
  3. 等待所有 agent 完成审查
  4. 收集每个 agent 的输出结果
```

**关键原则**:
- **并行执行**: 4 个 agent 同时工作，不等待其他 agent
- **职责分离**: 每个 agent 只检查自己的规则，不检查其他类别
- **独立输出**: 每个 agent 返回自己类别的问题列表

#### 4.3 结果合并

等待所有 agent 完成后，合并结果：

1. **收集所有 agent 的输出**
   - 从 GORM Agent 获取数据库相关问题
   - 从 Error & Safety Agent 获取错误处理/并发安全问题
   - 从 Naming & Logging Agent 获取命名/日志问题
   - 从 Organization Agent 获取代码组织问题

2. **按优先级排序**
   - 先按优先级排序 (P0 → P1 → P2)
   - 同优先级内按文件路径排序
   - 同文件内按行号排序

3. **去重和汇总**
   - 检查是否有重复问题（理论上不应该有，因为职责分离）
   - 统计各类别问题数量
   - 生成汇总报告

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

## Specialized Review Agents

This orchestrator **parallelly delegates** to 4 independent review agents:

1. **GORM Database Review Agent** (`agents/go-code-review/gorm-agent.md`)
   - **Rules**: 1.3.* (Database operations, GORM patterns)
   - **Executes in parallel** with other agents

2. **Error & Concurrency Safety Agent** (`agents/go-code-review/error-safety-agent.md`)
   - **Rules**: 1.1.*, 1.2.*, 1.4.*, 1.5.* (Error handling, nil safety, concurrency, JSON)
   - **Executes in parallel** with other agents

3. **Naming & Logging Standards Agent** (`agents/go-code-review/naming-logging-agent.md`)
   - **Rules**: 2.1.*, 2.2.* (Naming conventions, logging standards)
   - **Executes in parallel** with other agents

4. **Code Organization & Quality Agent** (`agents/go-code-review/organization-agent.md`)
   - **Rules**: 2.3.*, 2.4.*, 2.5.*, 3.* (Code organization, interfaces, quality, structure)
   - **Executes in parallel** with other agents

**Execution Model**: All 4 agents run **simultaneously** for each file, each focusing only on their assigned rule categories. Results are merged after all agents complete.

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
