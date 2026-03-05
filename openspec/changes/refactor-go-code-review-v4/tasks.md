## 1. 工具层（Tier 1 + Tier 2）

- [x] 1.1 创建 `skills/go-code-review/tools/` 目录
- [x] 1.2 实现 `tools/analyze-go.sh` — 统计文件行数、函数行数、嵌套深度，输出 `metrics.json`
- [x] 1.3 创建 `skills/go-code-review/rules/` 目录
- [x] 1.4 编写 `rules/safety.yaml` — SAFE-001～SAFE-010（nil/error/concurrency 模式规则）
- [x] 1.5 编写 `rules/data.yaml` — DATA-001～DATA-010（GORM/JSON 模式规则）
- [x] 1.6 编写 `rules/quality.yaml` — QUAL-001～QUAL-010（命名/行数/嵌套规则）
- [x] 1.7 编写 `rules/observability.yaml` — OBS-001～OBS-008（日志格式规则）
- [x] 1.8 实现 `tools/scan-rules.sh` — 读取 `rules/*.yaml`，grep 扫描 Go 文件，输出 `rule-hits.json`

## 2. Agent 层（Tier 3）

- [x] 2.1 删除旧 agent 文件（gorm-review.md、error-safety.md、naming-logging.md、organization-agent.md）
- [x] 2.2 创建 `agents/safety.md` — 安全性专家（并发安全上下文、防御性编程完整性）
- [x] 2.3 创建 `agents/data.md` — 数据层专家（N+1、序列化策略、类型语义）
- [x] 2.4 创建 `agents/design.md` — 设计哲学专家（UNIX 7 原则 + 代码变坏根源）
- [x] 2.5 创建 `agents/quality.md` — 代码质量专家（综合 metrics.json + rule-hits.json）
- [x] 2.6 创建 `agents/observability.md` — 可观测性专家（日志分层策略、错误消息质量）

## 3. Orchestrator

- [x] 3.1 重写 `skills/go-code-review/SKILL.md`（v4.0.0）
  - 移除内联规则，只保留调度逻辑
  - 新流程：git diff → analyze-go.sh → scan-rules.sh → 5 个并行 agent → 聚合输出
  - 更新架构说明（三层图示）

## 4. 测试用例

- [x] 4.1 补充 `test-cases/go-code-review/bad/design_patterns_bad.go` — 覆盖 UNIX 哲学违规场景
- [x] 4.2 更新 `test-cases/go-code-review/expected_issues.json` — 加入新 agent 的预期输出
- [x] 4.3 更新 `test-cases/go-code-review/run_test.sh` — 加入 analyze-go.sh 和 scan-rules.sh 的执行步骤
- [x] 4.4 验证 scan-rules.sh 对现有 bad/ 测试用例的命中率

## 5. 文档

- [x] 5.1 更新 `README.md` — v4.0.0 架构说明、Breaking Changes
- [x] 5.2 更新 `CLAUDE.md` — 新目录结构
