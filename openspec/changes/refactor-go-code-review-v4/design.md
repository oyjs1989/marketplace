# Design: Go Code Review v4.0 — Three-Tier Expert Architecture

## Context

当前 v3.0 已完成从多技能到单插件的整合（4 个并行 agent），但存在根本性的设计缺陷：AI 在做大量本该由脚本完成的确定性检查，且规则以自由文本形式嵌入 prompt，不可配置、不可维护。

spec.md 提出了关键工程经验：**职责分离（脚本做机械，命令做智能）** + **配置驱动（规则作为数据）**。v4.0 将这两个原则落地到 Go Code Review 插件中。

## Goals / Non-Goals

**Goals:**
- 规则结构化：从自由文本改为 YAML，可独立维护、可按需开关
- AI 专注判断：Tier 1+2 覆盖所有确定性检查，AI 只处理需要上下文理解的部分
- Agent 专业化：按领域专家视角划分，而非按规则编号划分
- 纳入设计哲学：UNIX 7 条原则 + 代码变坏根源成为正式 review 维度

**Non-Goals:**
- 不构建完整的 SDD 方法论框架（Code Review 是检查点，不是端到端流程）
- 不实现 `/review.learn`（学习机制）——留给 v5.0
- 不实现 `/review.rules add/disable`——留给 v5.0

## Architecture Decisions

### Decision 1: 三层架构

**What**: 在工具层（Tier 1）和 AI 层（Tier 3）之间插入规则扫描层（Tier 2）。

**Why**:
- 确定性规则（fmt.Errorf、db.Save、UserId 等）用 regex 检测精度 = 100%，成本 ≈ 0
- AI 对这类检查的精度 < 100%，且浪费 token
- 规则作为 YAML 数据后，未来可以扩展 enable/disable、统计漏报率

**Alternatives considered**:
- 全部交给 AI：精度不稳定，token 浪费
- 全部用脚本：无法处理需要语义理解的判断（设计哲学、N+1 模式等）

### Decision 2: 5 个专家 Agent（按领域视角划分）

**What**: 用"一个资深工程师 review 时的平行思路"替代"按规则编号分组"。

**5 个正交视角**:
- 🔴 **safety** — "这段代码会崩/死锁/数据损坏吗？"
- 🗄️ **data** — "数据存取正确高效吗？"
- 🏗️ **design** — "设计合理吗？能活过 100 万行代码吗？"
- 📐 **quality** — "新人 5 分钟能看懂吗？"
- 👁️ **observability** — "凌晨 3 点出故障，能快速定位吗？"

**Why better than v3.0**:
- 每个 agent 知识域聚焦，判断更准确
- 视角正交，结果不重叠
- 设计哲学有了专属 agent（design），不再分散在各处

### Decision 3: YAML 规则结构

**What**: 参照 spec.md 的规则结构设计。

```yaml
rules:
  - id: <PREFIX>-<NNN>       # 唯一 ID，便于引用和开关
    name: <中文规则名>
    severity: P0|P1|P2
    pattern:
      type: regex
      match: '<正则表达式>'
    message: "<中文问题描述>"
    examples:
      bad: '<违规示例>'
      good: '<修复示例>'
```

**规则 ID 前缀约定**:
- `SAFE-` → safety.yaml (nil/error/concurrency)
- `DATA-` → data.yaml (GORM/JSON)
- `QUAL-` → quality.yaml (naming/size)
- `OBS-` → observability.yaml (logging)

**Why YAML, not embedded in agent prompts**:
- 单一数据源：规则只定义一次
- 可维护：修改规则不用修改 agent prompt
- 可扩展：未来支持 enable/disable

### Decision 4: scan-rules.sh 设计

**What**: 读取 `rules/*.yaml`，对 Go 文件执行 grep 扫描，输出 `rule-hits.json`。

```json
{
  "hits": [
    {
      "rule_id": "SAFE-001",
      "severity": "P0",
      "file": "service/user.go",
      "line": 45,
      "matched": "return fmt.Errorf(\"get user failed: %v\", err)",
      "message": "使用 fmt.Errorf 会丢失错误堆栈，请改用 errors.Wrapf"
    }
  ],
  "summary": {
    "P0": 2,
    "P1": 5,
    "P2": 1
  }
}
```

**Implementation**: bash + grep，解析 YAML 中的 `pattern.match` 字段，逐文件扫描。

### Decision 5: design-agent 的内容来源

**What**: design-agent 覆盖两个来源的内容：

1. **UNIX 设计哲学**（来自 code_review.md）:
   - KISS（大道至简）
   - 组合原则（组合优于继承）
   - 吝啬原则（代码越少越好）
   - 透明性原则（函数行为可在大脑中构建完整过程）
   - 通俗原则（接口避免标新立异）
   - 缄默原则（没有好说的就沉默）
   - 补救原则（异常立即高调退出）

2. **代码变坏根源诊断**（来自 code_review.md）:
   - 重复代码（同一逻辑多处实现）
   - 无领域模型（上手就写 if/else）
   - OOP 滥用（不合理继承树）
   - 对合理性缺乏苛求（该用 defer 不用，该 early return 不 return）
   - 无设计就写代码

## Tier Details

### Tier 1: analyze-go.sh

输入：Go 文件列表（来自 git diff）
输出：`metrics.json`

```json
{
  "files": [
    {
      "path": "service/user.go",
      "lines": 342,
      "violations": {
        "over_800_lines": false
      },
      "functions": [
        {
          "name": "CreateUser",
          "start_line": 45,
          "lines": 111,
          "max_nesting": 5,
          "violations": {
            "over_80_lines": true,
            "over_4_nesting": true
          }
        }
      ]
    }
  ],
  "summary": {
    "files_over_800": 0,
    "functions_over_80": 2,
    "nesting_violations": 1
  }
}
```

### Tier 2: scan-rules.sh

输入：`rules/*.yaml` + Go 文件列表
输出：`rule-hits.json`（见 Decision 4）

### Tier 3: 5 个 AI Agents

每个 agent 接收：
- 变更代码内容
- `metrics.json`（Tier 1 结果）
- `rule-hits.json` 中属于自己负责的命中项
- 自己专属的判断规则（无法被 regex 捕获的部分）

每个 agent 输出：
- 确认 Tier 2 命中项（标注行号、补充上下文）
- 补充 Tier 2 无法检测的判断性问题

## Risks / Trade-offs

| Risk | Mitigation |
|------|-----------|
| scan-rules.sh 正则误报 | 规则设计时加 `examples.bad` 验证，先覆盖 P0 规则 |
| agent 与 YAML 规则重复报告 | agent 明确说明"Tier 2 已报告的不重复输出，只补充上下文" |
| YAML 维护成本 | 初期只覆盖 30-40 条最高频的确定性规则，其余仍由 AI 处理 |
| design-agent 过于主观 | 每条 UNIX 原则配具体 Go 代码判据，减少主观判断 |

## Migration Plan

1. 创建 `rules/` 目录，先迁移 10-15 条最确定的 P0 规则到 YAML
2. 实现 `tools/scan-rules.sh`（bash，解析 YAML 中的 pattern.match）
3. 重写 5 个 agent 文件
4. 重写 `SKILL.md` orchestrator
5. 更新测试用例
6. 归档旧 agent 文件

## Open Questions

- scan-rules.sh 解析 YAML：直接用 grep 还是引入 `yq`？（建议先用 grep 解析，保持零依赖）
- quality-agent 与 design-agent 在"函数过长"上有重叠：quality 报量化违规，design 分析为何过长（设计问题）→ 允许这种分层补充
