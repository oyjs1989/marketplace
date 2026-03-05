# Change: Refactor Go Code Review to v4.0 — Three-Tier Expert Architecture

**Change ID**: refactor-go-code-review-v4
**Status**: Draft
**Version**: v4.0.0
**Date**: 2026-03-04

## Why

v3.0 存在三个根本性问题：

1. **Agent 职责不清** — 按规则编号（1.3.\*、2.1.\*）划分 agent，不是按专业视角，导致每个 agent 知识域太宽、判断准确率低
2. **AI 在做机械活** — 大量可以用正则匹配的规则（`fmt.Errorf`、`db.Save`、`UserId`…）都交给 AI 判断，浪费 token 且精度不稳定
3. **规则是自由文本** — 规则埋在 prompt 里，无法按需开关、无法统计漏报率、不可维护

核心问题与 spec.md 识别的一致：**规则太自由文本，缺乏结构 → 不被遵循**。

## What Changes

### 架构：从两层变三层

**旧架构（v3.0）**：
```
工具层（analyze-go.sh）
    ↓
4 个 AI agents（按规则编号划分）
```

**新架构（v4.0）**：
```
Tier 1: analyze-go.sh    → metrics.json    （量化：文件/函数行数、嵌套深度）
Tier 2: scan-rules.sh    → rule-hits.json  （模式匹配：YAML 结构化规则扫描）
Tier 3: 5 个 AI agents   → 判断报告         （只处理需要理解上下文的判断）
```

**职责分离原则（来自 spec.md）**：
- 脚本做机械（确定性检查）→ Tier 1 + Tier 2
- AI 做智能（需要理解上下文的判断）→ Tier 3

### 新目录结构

```
skills/go-code-review/
├── SKILL.md              # Orchestrator（只负责调度，不内联规则）
├── agents/               # 5 个领域专家 agent
│   ├── safety.md         # 安全性专家
│   ├── data.md           # 数据层专家
│   ├── design.md         # 架构设计哲学专家
│   ├── quality.md        # 代码质量专家
│   └── observability.md  # 可观测性专家
├── rules/                # 结构化 YAML 规则（单一数据源）
│   ├── safety.yaml       # nil / error / concurrency 模式规则
│   ├── data.yaml         # GORM / JSON 模式规则
│   ├── quality.yaml      # 命名 / 行数 / 嵌套规则
│   └── observability.yaml # 日志格式规则
└── tools/
    ├── analyze-go.sh     # 量化分析（行数、嵌套）→ metrics.json
    └── scan-rules.sh     # YAML 规则扫描 → rule-hits.json
```

### 新 Agent 划分（5 个，按专业视角）

| Agent | 专业视角 | 只处理…（Tier 2 覆盖不了的） |
|-------|---------|---------------------------|
| 🔴 `safety` | 安全与正确性 | 并发安全上下文判断、防御性编程完整性 |
| 🗄️ `data` | 数据层 | N+1 查询、序列化策略、类型语义合理性 |
| 🏗️ `design` | 架构设计哲学 | UNIX 7 原则、领域模型质量、代码变坏根源 |
| 📐 `quality` | 代码质量 | 综合 metrics.json + rule-hits.json，补充命名语义判断 |
| 👁️ `observability` | 可观测性 | 日志分层策略、错误消息描述质量 |

### YAML 规则结构（来自 spec.md 设计）

```yaml
# rules/safety.yaml
rules:
  - id: SAFE-001
    name: 错误必须用 errors.Wrap 包装
    severity: P0
    pattern:
      type: regex
      match: 'fmt\.Errorf\('
    message: "使用 fmt.Errorf 会丢失错误堆栈，请改用 errors.Wrapf"
    examples:
      bad: 'return fmt.Errorf("get user failed: %v", err)'
      good: 'return errors.Wrapf(err, "get user failed")'

  - id: SAFE-002
    name: 禁止在业务函数中使用 panic
    severity: P0
    pattern:
      type: regex
      match: '\bpanic\('
    message: "业务函数不应使用 panic，请改用 error 返回"
```

### 主流程（SKILL.md）

```
Step 1: git diff --name-only → 获取变更文件列表
Step 2: tools/analyze-go.sh  → metrics.json
Step 3: tools/scan-rules.sh  → rule-hits.json
Step 4: 并行启动 5 个 agent，传入 metrics.json + rule-hits.json + 代码内容
Step 5: 聚合所有 agent 输出，按 P0→P1→P2 排序，去重，输出中文报告
```

## Impact

- **BREAKING**: 4 个旧 agent 文件（gorm-review、error-safety、naming-logging、organization）替换为 5 个新 agent
- **NEW**: `rules/` 目录（4 个 YAML 文件）
- **NEW**: `tools/scan-rules.sh`
- Affected files:
  - `skills/go-code-review/SKILL.md` — 重写
  - `skills/go-code-review/agents/` — 全部替换
  - `skills/go-code-review/tools/` — 新增 scan-rules.sh
  - `skills/go-code-review/rules/` — 全新目录
  - `test-cases/go-code-review/` — 补充设计哲学测试用例
