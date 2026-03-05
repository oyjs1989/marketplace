---
name: go-code-review
description: This skill should be used when the user asks to "review Go code", "check Go code quality", "review this PR", "code review", or mentions Go code standards, GORM best practices, error handling patterns, concurrency safety, design philosophy, or UNIX principles. Orchestrates comprehensive Go code reviews using a three-tier architecture: quantitative tools + YAML pattern scanning + 5 domain-expert AI agents.
version: 4.0.0
---

# Go Code Review Skill (v4.0.0)

## When to Use This Skill

This skill activates when users need help with:
- Reviewing Go code changes against coding standards
- Checking code quality and identifying potential issues
- Performing PR reviews for Go projects
- Validating database operations and data layer correctness
- Checking error handling, concurrency safety, and nil safety
- Analyzing design philosophy and UNIX principles compliance
- Evaluating observability: logging strategy and error message quality
- Reviewing naming conventions, code structure, and readability

## Architecture: Three-Tier Expert Review

```
输入：git diff 变更的 Go 文件
         │
         ▼
┌─────────────────────────────────┐
│  Tier 1: tools/analyze-go.sh   │  → metrics.json
│  量化：文件行数/函数行数/嵌套深度  │    (过 800/80/4 的违规)
└─────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────┐
│  Tier 2: tools/scan-rules.sh   │  → rule-hits.json
│  模式匹配：38 条 YAML 规则 regex │    (SAFE/DATA/QUAL/OBS 命中)
└─────────────────────────────────┘
         │
         ▼
┌──────────────────────────────────────────────────────────┐
│  Tier 3: 5 个领域专家 Agent（并行）                         │
│  🔴 safety      │ 安全与正确性，上下文并发判断              │
│  🗄️  data        │ 数据层，N+1，序列化，类型语义             │
│  🏗️  design      │ UNIX 7 原则，领域模型，代码变坏根源       │
│  📐 quality     │ 综合 metrics.json，命名语义，可读性       │
│  👁️  observability│ 日志分层策略，错误消息质量               │
└──────────────────────────────────────────────────────────┘
         │
         ▼
聚合：P0 → P1 → P2，去重，中文报告输出到 code_review.result
```

### Tier 1 — 量化分析工具

Script: `tools/analyze-go.sh`
Output: `/tmp/metrics.json`

Measures per file and per function:
- File line count (threshold: 800 lines)
- Function line count (threshold: 80 lines)
- Nesting depth (threshold: 4 levels)

### Tier 2 — YAML 规则扫描

Script: `tools/scan-rules.sh`
Output: `/tmp/rule-hits.json`

Scans against 38 deterministic regex rules across four YAML files:
- `rules/safety.yaml` — SAFE-001 to SAFE-010
- `rules/data.yaml` — DATA-001 to DATA-010
- `rules/quality.yaml` — QUAL-001 to QUAL-010
- `rules/observability.yaml` — OBS-001 to OBS-008

### Tier 3 — 5 个领域专家 Agent

| Agent | Expert Perspective |
|-------|--------------------|
| safety (red) | 安全与正确性：会崩/死锁/数据损坏吗？ |
| data (blue) | 数据层：存取正确高效吗？ |
| design (purple) | 架构设计哲学：能活过百万行代码吗？ |
| quality (green) | 代码质量：新人 5 分钟能看懂吗？ |
| observability (yellow) | 可观测性：凌晨 3 点能快速定位吗？ |

Each agent receives the full code diff plus the subset of `rule-hits.json` relevant to its domain. Agents confirm Tier 2 hits with business context and surface additional judgment-based issues that regex cannot detect.

## Review Workflow

### Step 1: 获取变更文件

```bash
git diff master --name-only | grep '\.go$'
# 或针对特定 commit
git diff HEAD~1 --name-only | grep '\.go$'
```

### Step 2: 运行 Tier 1 量化分析

```bash
git diff master --name-only | grep '\.go$' | bash tools/analyze-go.sh > /tmp/metrics.json
```

读取 `/tmp/metrics.json`，记录量化违规（文件 > 800 行、函数 > 80 行、嵌套 > 4 层）。

### Step 3: 运行 Tier 2 规则扫描

```bash
git diff master --name-only | grep '\.go$' | bash tools/scan-rules.sh > /tmp/rule-hits.json
```

读取 `/tmp/rule-hits.json`，记录所有 regex 命中。

### Step 4: 启动 5 个专家 Agent（并行）

所有 agent 同时接收：
- 变更代码内容（通过 `git diff master -- <files>` 读取）
- `/tmp/metrics.json`（Tier 1 结果）
- `/tmp/rule-hits.json` 中属于各自领域的命中项

**Agent 分工与 Tier 2 关系：**

- **safety agent** — 确认 SAFE-001~010 命中；处理并发/context/防御性编程判断
- **data agent** — 确认 DATA-001~010 命中；处理 N+1/序列化/事务边界判断
- **design agent** — 无 Tier 2 规则；专注 UNIX 7 原则 + 5 大代码变坏根源
- **quality agent** — 确认 QUAL-001~010 命中 + 综合 metrics.json；处理命名语义/注释质量
- **observability agent** — 确认 OBS-001~008 命中；处理日志分层策略/错误消息质量

### Step 5: 聚合输出

收集所有 agent 输出后：
1. 合并 Tier 2 命中（已在 rule-hits.json 中）和 agent 补充的判断性问题
2. 去重：同一位置的问题只保留最高严重度
3. 按 P0 → P1 → P2 排序
4. 输出到 `code_review.result`

## Output Format

**重要**：所有审查输出必须使用中文。

```markdown
# Go 代码审查报告（v4.0.0）

## 审查摘要

| 指标 | 数量 |
|------|------|
| P0（必须修复） | X 个 |
| P1（强烈建议） | X 个 |
| P2（建议优化） | X 个 |

## 量化违规（Tier 1）

（来自 metrics.json，由 quality agent 报告）

## P0 问题（必须修复）

### 问题 - [P0] <问题类别>（来自：<agent名称>/<rule-id>）
**位置**: path/to/file.go:行号
**类别**: <具体类别>
**原始代码**:
```go
// 问题代码
```
**问题描述**: <中文说明>
**修改建议**:
```go
// 修复代码
```

## P1 问题（强烈建议）
...

## P2 问题（建议优化）
...
```

## Manual Agent Invocation

Individual agents can be invoked directly without running the full orchestrator:

```
直接调用 safety agent
直接调用 data agent
直接调用 design agent
直接调用 quality agent
直接调用 observability agent
```
