# Go Code Review Skill - 测试指南

## 快速开始

### 方法 1: 快速测试（推荐）

在 Claude Code 中直接测试：

```bash
# 测试违规代码检测
Review test-cases/go-code-review/bad/user_service_bad.go

# 测试正确代码（不应报告问题）
Review test-cases/go-code-review/good/user_service_good.go
```

### 方法 2: 批量测试

```bash
# 测试所有违规代码
Review test-cases/go-code-review/bad/*.go

# 查看结果
cat code_review.result
```

### 方法 3: 自动化验证

```bash
# 1. 运行审查（在Claude Code中）
Review test-cases/go-code-review/bad/user_service_bad.go

# 2. 验证结果
bash test-cases/go-code-review/validate_results.sh code_review.result

# 3. 查看详细报告
cat test-cases/go-code-review/validation_report.md
```

## 测试工作流

### 完整测试流程

```
┌─────────────────────────────────────────────────┐
│ 1. 修改技能代码                                  │
│    (更新 SKILL.md 或规则)                        │
└────────────────┬────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────┐
│ 2. 运行 Code Review                             │
│    Review test-cases/go-code-review/bad/*.go   │
└────────────────┬────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────┐
│ 3. 验证检测结果                                  │
│    bash validate_results.sh code_review.result │
└────────────────┬────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────┐
│ 4. 检查报告                                      │
│    - 是否所有预期问题都被检测？                   │
│    - 是否有漏检？                                 │
│    - good 文件是否误报？                          │
└────────────────┬────────────────────────────────┘
                 │
        ┌────────┴────────┐
        │                 │
        ▼                 ▼
    ✅ 通过            ❌ 失败
    提交代码          修复问题
```

## 预期结果参考

### user_service_bad.go

**预期检测**:
- **P0 问题**: 15+ 个
  - Rule 1.1.1 (fmt.Errorf)
  - Rule 1.1.3 (panic)
  - Rule 1.1.4 (未检查错误)
  - Rule 1.2.1 (nil检查)
  - Rule 1.3.1 (缺少WHERE)
  - Rule 1.3.2 (缺少SELECT)
  - Rule 1.3.4 (Save方法)
  - Rule 1.3.5 (First vs Take)
  - Rule 1.3.6 (gorm tag)
  - Rule 1.3.8 (字段顺序)
  - Rule 1.4.1 (ErrorGroup)
  - Rule 1.4.2 (limiter)
  - Rule 1.5.1 (JSON拼接)

- **P1 问题**: 5+ 个
  - Rule 2.1.6 (魔法数字)
  - Rule 2.1.7 (命名不一致)
  - Rule 2.2.1 (日志字段命名)
  - Rule 2.2.2 (日志字段类型)
  - Rule 2.2.7 (缺少日志)

- **P2 问题**: 3+ 个
  - Rule 2.5.3 (函数过长)
  - Rule 2.5.4 (缺少注释)
  - Rule 2.5.7 (复杂条件)

**总计**: 23+ 个问题

### project_structure_bad.go

**预期检测**:
- **P1 问题**: 8+ 个
- **P2 问题**: 5+ 个
- **总计**: 13+ 个问题

### user_service_good.go

**预期检测**: **0 个问题** （任何检测都是误报）

## 常见检查点

### ✅ 必检项

运行测试后，确认以下关键规则被检测：

```bash
# 使用 validate_results.sh 会自动检查这些
- [ ] Rule 1.1.1 - 错误包装 (fmt.Errorf)
- [ ] Rule 1.3.1 - 缺少 WHERE 条件
- [ ] Rule 1.3.4 - Save() 方法
- [ ] Rule 1.2.1 - nil 指针检查
- [ ] Rule 1.4.1 - ErrorGroup 使用
```

### ⚠️ 漏检检查

如果以下情况出现，说明存在漏检：

1. **P0 问题数 < 10** (user_service_bad.go)
2. **关键规则未检测** (1.1.1, 1.3.1, 1.3.4)
3. **good 文件报告了问题** (误报)
4. **某个规则类别完全缺失** (如所有 GORM 规则)

## 测试文件说明

### 测试文件结构

```
test-cases/go-code-review/
├── README.md                    # 测试用例总览
├── TESTING_GUIDE.md            # 本文件 - 测试指南
├── expected_issues.json        # 预期结果配置
├── run_test.sh                 # 自动化测试脚本
├── validate_results.sh         # 结果验证脚本
├── bad/                        # 违规代码
│   ├── user_service_bad.go     # P0 问题集合
│   └── project_structure_bad.go # P1/P2 问题集合
└── good/                       # 正确实现
    └── user_service_good.go    # 零问题参考
```

### 每个测试文件的用途

| 文件 | 目的 | 预期问题数 |
|------|------|-----------|
| user_service_bad.go | 测试 P0 关键问题检测 | 20+ |
| project_structure_bad.go | 测试 P1/P2 质量问题检测 | 13+ |
| user_service_good.go | 测试误报率（应为0） | 0 |

## 脚本使用说明

### validate_results.sh

**用途**: 自动验证审查结果是否符合预期

**使用方法**:
```bash
# 默认验证 code_review.result
bash validate_results.sh

# 验证指定文件
bash validate_results.sh my_custom_result.md

# 查看验证报告
cat validation_report.md
```

**输出**:
- ✅ 通过 - 所有关键规则都被检测
- ❌ 失败 - 存在漏检，需要改进

### run_test.sh

**用途**: 完整的测试流程指导（需要手动运行审查）

**使用方法**:
```bash
bash run_test.sh
```

**输出**: 测试报告和统计信息

### expected_issues.json

**用途**: 定义每个测试文件的预期检测结果

**结构**:
```json
{
  "test_files": {
    "user_service_bad.go": {
      "expected_issues": {
        "P0": {
          "min": 15,
          "rules": ["1.1.1", "1.3.1", ...]
        }
      }
    }
  }
}
```

## 回归测试

### 什么时候运行测试？

1. **修改技能代码后** - 确保没有引入漏检
2. **添加新规则后** - 验证新规则能被检测
3. **更新 SKILL.md 后** - 确保文档和检测一致
4. **发布新版本前** - 质量保证
5. **每周定期** - 持续监控

### 回归测试清单

```bash
# 1. 修改技能代码
vim skills/go-code-review/xxx/SKILL.md

# 2. 运行完整测试
Review test-cases/go-code-review/bad/*.go

# 3. 验证结果
bash test-cases/go-code-review/validate_results.sh

# 4. 检查报告
cat test-cases/go-code-review/validation_report.md

# 5. 如果通过，提交代码
git add .
git commit -m "feat: 更新规则检测逻辑"

# 6. 如果失败，修复问题并重新测试
```

## 扩展测试用例

### 添加新的测试文件

如果发现某个规则没有测试覆盖：

1. **创建新的测试文件**:
```bash
cd test-cases/go-code-review/bad/
cp user_service_bad.go new_feature_bad.go
```

2. **添加特定违规代码**:
```go
// 针对特定规则添加违规代码
// Rule X.Y.Z - 规则描述
func BadExample() {
    // 违规代码
}
```

3. **更新 expected_issues.json**:
```json
{
  "test_files": {
    "new_feature_bad.go": {
      "expected_issues": {
        "P0": {
          "min": 3,
          "rules": ["X.Y.Z"]
        }
      }
    }
  }
}
```

4. **运行测试验证**

### 覆盖缺失的规则

当前测试覆盖率: **70%** (35/50 规则)

**缺失的规则**:
- 1.4.3, 1.4.4 (并发控制)
- 1.5.2 (JSON处理)
- 2.4.* (接口设计)
- 3.1.* (项目结构)
- 3.2.* (测试规范)
- 3.3.2 (配置管理)

**TODO**: 创建专门的测试文件覆盖这些规则

## 故障排查

### 问题 1: 检测结果为空

**症状**: code_review.result 文件不存在或为空

**原因**:
- 技能未正确加载
- Git diff 为空（没有修改）
- 技能执行出错

**解决**:
```bash
# 检查技能是否加载
cat .claude/settings.local.json

# 确保文件已修改（或强制审查）
git add test-cases/go-code-review/bad/*.go

# 直接审查文件而不是 diff
Review test-cases/go-code-review/bad/user_service_bad.go
```

### 问题 2: P0 问题数量不足

**症状**: 检测到的 P0 问题少于预期

**原因**:
- 某些规则未被触发
- 规则检测逻辑有bug
- 测试文件不包含某些违规

**解决**:
```bash
# 1. 查看检测到的规则
grep "规则" code_review.result | sort -u

# 2. 对比预期规则列表
cat expected_issues.json | grep rules

# 3. 找出缺失的规则
# 4. 检查对应的 SKILL.md 是否正确定义了检测逻辑
```

### 问题 3: good 文件报告了问题（误报）

**症状**: user_service_good.go 被标记了问题

**原因**:
- 规则检测过于严格
- good 文件本身有问题
- 规则定义不清晰

**解决**:
1. 检查报告的具体问题
2. 验证 good 文件是否真的正确
3. 调整规则的检测逻辑或异常情况
4. 更新 good 文件或规则文档

## 性能基准

### 预期性能

| 操作 | 预期时间 |
|------|---------|
| 单文件审查 (bad) | < 30秒 |
| 批量审查 (3个文件) | < 60秒 |
| 验证脚本运行 | < 5秒 |

### 性能测试

```bash
# 测试审查速度
time Review test-cases/go-code-review/bad/user_service_bad.go

# 测试批量审查
time Review test-cases/go-code-review/bad/*.go
```

## 持续改进

### 记录测试结果

每次测试后，记录：
- 日期和版本号
- 检测到的问题总数
- P0/P1/P2 分布
- 是否有漏检或误报
- 需要改进的地方

### 提交测试报告

```bash
# 保存测试结果
cp code_review.result test-results/$(date +%Y%m%d).result
cp validation_report.md test-results/$(date +%Y%m%d).report

# 提交到 Git（可选）
git add test-results/
git commit -m "test: 记录 $(date +%Y-%m-%d) 测试结果"
```

## 总结

使用这套测试系统，你可以：

✅ **快速验证** - 每次修改后立即测试
✅ **防止回归** - 确保新改动不破坏现有检测
✅ **提升质量** - 持续改进检测准确性
✅ **自动化** - 减少手动验证工作量
✅ **可追踪** - 记录测试历史和趋势

---

**维护者**: Jason Ouyang
**更新日期**: 2026-01-04
**版本**: 1.0.0
