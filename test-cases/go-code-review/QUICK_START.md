# 快速开始 - 5分钟掌握测试系统

## 🚀 一分钟快速测试

```bash
# 在 Claude Code 中运行这个命令:
Review test-cases/go-code-review/bad/user_service_bad.go
```

就这么简单！检查生成的 `code_review.result` 文件，应该看到 20+ 个问题被检测到。

## 📋 完整测试流程（5分钟）

### 步骤 1: 测试违规代码检测（2分钟）

```bash
# 在 Claude Code 中
Review test-cases/go-code-review/bad/user_service_bad.go
```

**预期结果**: 检测到 20+ 个问题
- P0 问题: 15+ 个
- P1 问题: 5+ 个
- P2 问题: 3+ 个

### 步骤 2: 验证检测结果（1分钟）

```bash
# 在终端中
cd test-cases/go-code-review
bash validate_results.sh code_review.result
```

**预期输出**:
```
✓ 找到审查结果文件: code_review.result
✓ 检测到 23 个不同的规则违规
✓ user_service_bad.go: 检测到足够的 P0 问题 (20 >= 10)
✓ Rule 1.1.1: 已检测
✓ Rule 1.3.1: 已检测
...
✓ 验证通过！
```

### 步骤 3: 查看详细报告（1分钟）

```bash
cat validation_report.md
```

### 步骤 4: 测试误报检测（1分钟）

```bash
# 在 Claude Code 中
Review test-cases/go-code-review/good/user_service_good.go
```

**预期结果**: 不应报告任何问题（如果有，那是误报）

## 🎯 核心文件说明

| 文件 | 用途 | 优先级 |
|------|------|--------|
| `bad/user_service_bad.go` | 包含20+个违规的测试代码 | ⭐⭐⭐ |
| `good/user_service_good.go` | 正确实现，用于测试误报 | ⭐⭐⭐ |
| `validate_results.sh` | 自动验证检测结果 | ⭐⭐⭐ |
| `TESTING_GUIDE.md` | 完整测试指南 | ⭐⭐ |
| `expected_issues.json` | 预期结果配置 | ⭐⭐ |
| `CI_INTEGRATION.md` | CI/CD集成指南 | ⭐ |

## 🔍 关键检查点

测试后，确认这些关键规则被检测到：

```bash
# 快速检查
grep "规则 1.1.1" code_review.result  # 错误包装
grep "规则 1.3.1" code_review.result  # WHERE条件
grep "规则 1.3.4" code_review.result  # Save方法
grep "规则 1.2.1" code_review.result  # nil检查
grep "规则 1.4.1" code_review.result  # ErrorGroup
```

如果这5条都检测到了，基本功能就正常了 ✅

## 📊 测试决策树

```
开始测试
    ↓
运行 Review bad/user_service_bad.go
    ↓
检查 code_review.result
    ↓
    ├─ 文件存在?
    │   ├─ 否 → 检查技能是否加载
    │   └─ 是 → 继续
    ↓
    ├─ P0问题 >= 10?
    │   ├─ 否 → ⚠️ 可能漏检
    │   └─ 是 → ✅ 基本正常
    ↓
运行 validate_results.sh
    ↓
    ├─ 验证通过?
    │   ├─ 是 → ✅ 测试通过
    │   └─ 否 → ⚠️ 需要改进
    ↓
测试 good/user_service_good.go
    ↓
    ├─ 报告问题?
    │   ├─ 是 → ❌ 误报
    │   └─ 否 → ✅ 完美
    ↓
测试完成！
```

## 🛠️ 常用命令速查

### 测试命令

```bash
# 测试单个文件
Review test-cases/go-code-review/bad/user_service_bad.go

# 测试所有违规代码
Review test-cases/go-code-review/bad/*.go

# 验证结果
bash test-cases/go-code-review/validate_results.sh

# 查看报告
cat test-cases/go-code-review/validation_report.md
```

### 统计命令

```bash
# 统计P0问题数
grep -c "\[P0\]" code_review.result

# 统计所有问题
grep -c "### 问题" code_review.result

# 列出检测到的规则
grep -oP "规则 \K[0-9]+\.[0-9]+\.[0-9]+" code_review.result | sort -u

# 查看某个规则的所有实例
grep "规则 1.3.1" code_review.result
```

## 🎓 学习路径

### 新手（第1天）
1. ✅ 运行一次快速测试
2. ✅ 查看 code_review.result
3. ✅ 理解输出格式

### 进阶（第1周）
1. ✅ 阅读 TESTING_GUIDE.md
2. ✅ 使用 validate_results.sh
3. ✅ 理解 expected_issues.json

### 高级（第1月）
1. ✅ 集成到 Pre-commit hook
2. ✅ 添加自定义测试用例
3. ✅ 配置 CI/CD

## 💡 提示和技巧

### 提示 1: 快速验证修改

```bash
# 修改技能后，快速测试
Review test-cases/go-code-review/bad/user_service_bad.go && \
bash test-cases/go-code-review/validate_results.sh
```

### 提示 2: 对比两次测试结果

```bash
# 保存第一次结果
cp code_review.result result_v1.md

# 修改技能后再测试
Review test-cases/go-code-review/bad/user_service_bad.go

# 对比差异
diff -u result_v1.md code_review.result
```

### 提示 3: 批量测试性能

```bash
# 测试执行时间
time Review test-cases/go-code-review/bad/*.go

# 预期: < 60秒
```

### 提示 4: 调试特定规则

如果某个规则一直漏检：

```bash
# 1. 检查规则在文档中的定义
grep -A 20 "Rule 1.3.1" skills/go-code-review/gorm-review/SKILL.md

# 2. 确认测试代码包含该违规
grep -n "Rule 1.3.1" test-cases/go-code-review/bad/user_service_bad.go

# 3. 查看是否被检测到
grep "规则 1.3.1" code_review.result
```

## ❓ 常见问题 FAQ

### Q1: 测试结果为空怎么办？

**A**: 检查：
1. 技能是否正确加载？`cat .claude/settings.local.json`
2. 文件路径是否正确？
3. 是否在正确的目录运行？

### Q2: P0问题数量少于预期

**A**:
1. 运行 `bash validate_results.sh` 查看具体缺失的规则
2. 检查对应 SKILL.md 的检测逻辑
3. 确认测试文件包含该违规代码

### Q3: good文件报告了问题

**A**: 这是误报，需要：
1. 查看具体是哪条规则
2. 检查规则定义是否过于严格
3. 或者 good 文件本身有问题

### Q4: 如何添加新的测试用例？

**A**:
```bash
# 1. 创建新文件
cp bad/user_service_bad.go bad/new_feature_bad.go

# 2. 添加特定违规代码
vim bad/new_feature_bad.go

# 3. 测试
Review bad/new_feature_bad.go

# 4. 更新 expected_issues.json
```

## 🎉 成功标准

测试成功的标志：

- ✅ bad 文件检测到 20+ 个问题
- ✅ 关键规则（1.1.1, 1.3.1, 1.3.4, 1.2.1, 1.4.1）都被检测
- ✅ good 文件不报告任何问题
- ✅ validate_results.sh 返回"验证通过"
- ✅ 执行时间 < 60秒

## 📞 获取帮助

如果遇到问题：

1. 📖 查看 TESTING_GUIDE.md 的详细说明
2. 🔍 检查 validation_report.md 的错误信息
3. 🐛 在 GitHub Issues 报告问题
4. 💬 联系维护者: Jason Ouyang

## 🔗 相关文档

- [完整测试指南](TESTING_GUIDE.md) - 详细的测试说明
- [CI集成指南](CI_INTEGRATION.md) - 集成到CI/CD流水线
- [测试用例说明](README.md) - 测试文件的详细信息
- [预期结果配置](expected_issues.json) - 每个测试的预期结果

---

**就这么简单！开始测试吧 🚀**

*创建日期: 2026-01-04*
*维护者: Jason Ouyang*
