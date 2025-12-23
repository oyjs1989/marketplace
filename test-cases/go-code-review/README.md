# Go Code Review Test Cases

这个目录包含用于测试 Go Code Review Skill 的测试用例。

## 目录结构

```
test-cases/go-code-review/
├── bad/                          # 包含违规代码的测试文件
│   ├── user_service_bad.go      # 用户服务 - P0 违规示例
│   └── project_structure_bad.go # 项目结构 - P1/P2 违规示例
├── good/                         # 符合规范的代码示例
│   └── user_service_good.go     # 用户服务 - 正确实现
└── README.md                     # 本文件
```

## 测试文件说明

### 1. user_service_bad.go - P0 关键问题

包含**必须修复**的严重问题,涵盖以下规则:

#### 错误处理 (1.1.*)
- ❌ **Rule 1.1.1** - 使用 `fmt.Errorf` 而不是 `errors.Wrap` (无错误栈)
- ❌ **Rule 1.1.3** - 使用 `panic` 而不是返回错误
- ❌ **Rule 1.1.4** - 未检查错误 (`err` 被忽略)

#### Nil 安全 (1.2.*)
- ❌ **Rule 1.2.1** - 未检查 nil 指针就直接使用

#### GORM 数据库操作 (1.3.*)
- ❌ **Rule 1.3.1** - 缺少显式 WHERE 条件
- ❌ **Rule 1.3.2** - 未显式指定 SELECT 列
- ❌ **Rule 1.3.4** - 使用 `Save()` 方法 (应使用 `Updates`)
- ❌ **Rule 1.3.5** - 使用 `First()` 而不是 `Take()`
- ❌ **Rule 1.3.6** - 缺少 `gorm:"column:xxx"` 标签
- ❌ **Rule 1.3.8** - 字段顺序不规范

#### 并发安全 (1.4.*)
- ❌ **Rule 1.4.1** - 未使用 `errgroup.Group`
- ❌ **Rule 1.4.2** - 未使用 limiter 控制并发数
- ❌ 并发写入 slice 未加锁

#### JSON 处理 (1.5.*)
- ❌ **Rule 1.5.1** - 手动拼接 JSON 字符串

### 2. project_structure_bad.go - P1/P2 质量问题

包含**强烈建议修复**和**建议优化**的问题:

#### 代码组织 (2.3.*)
- ❌ **Rule 2.3.1** - 嵌入字段未放在结构体开头
- ❌ **Rule 2.3.2** - import 分组不规范
- ❌ **Rule 2.3.3** - 未使用 `make()` 初始化 slice/map
- ❌ **Rule 2.3.5** - switch 语句缺少 default
- ❌ **Rule 2.3.6** - init() 函数包含复杂逻辑
- ❌ **Rule 2.3.7** - 使用全局变量
- ❌ **Rule 2.3.8** - 函数顺序混乱
- ❌ **Rule 2.3.10** - TODO 未标注负责人

#### 命名规范 (2.1.*)
- ❌ **Rule 2.1.6** - 魔法数字 (100, 1 等)
- ❌ **Rule 2.1.7** - 命名不一致 (`userId` vs `userID`)

#### 日志规范 (2.2.*)
- ❌ **Rule 2.2.1** - 日志字段未使用 snake_case
- ❌ **Rule 2.2.2** - 未显式指定字段类型
- ❌ **Rule 2.2.7** - 缺少入口/出口日志

#### 代码质量 (2.5.*)
- ❌ **Rule 2.5.1** - 拼写错误 (`Msg` 应为 `Message`)
- ❌ **Rule 2.5.3** - 函数过长
- ❌ **Rule 2.5.4** - 公共函数缺少注释
- ❌ **Rule 2.5.5** - 枚举注释未使用中文
- ❌ **Rule 2.5.7** - 复杂条件未提取为变量/函数

#### 配置 (3.3.*)
- ❌ **Rule 3.3.1** - 超时配置应使用 `time.Duration`

### 3. user_service_good.go - 正确实现

展示符合所有规范的正确实现:

#### ✅ 错误处理
- 使用 `errors.Wrap/Wrapf` 包装错误并保留栈
- 返回错误而不是 panic
- 立即检查所有错误

#### ✅ GORM 操作
- 显式 WHERE 条件
- 显式 SELECT 列
- 使用 `Updates()` 而不是 `Save()`
- 使用 `Take()` 而不是 `First()`
- 完整的 `gorm` 标签
- 正确的字段顺序

#### ✅ 并发安全
- 使用 `errgroup.WithContext`
- 使用 `SetLimit` 控制并发
- 正确的错误传播

#### ✅ 命名和日志
- CamelCase 命名
- snake_case 日志字段
- 显式类型的日志字段
- 入口/出口日志

#### ✅ 代码组织
- 清晰的接口定义
- 常量替代魔法数字
- 请求结构体 (4+ 参数)
- TODO 标注负责人

## 如何使用

### 运行 Code Review

#### 方法 1: 使用 Claude Code 直接审查

```bash
# 在 Claude Code 中执行
Review the bad code files in test-cases/go-code-review/bad/
```

#### 方法 2: 创建 Git Diff

```bash
# 模拟 PR 场景
cd test-cases/go-code-review
git add bad/*.go
git diff --cached bad/*.go > test.diff

# 让 Claude Code 审查 diff
Review this diff file: test.diff
```

#### 方法 3: 对比好坏代码

```bash
# 对比两个版本
diff -u bad/user_service_bad.go good/user_service_good.go
```

### 预期结果

运行 Code Review Skill 后,应该生成 `code_review.result` 文件,包含:

1. **文件路径和行号**
2. **违反的规则编号** (如 Rule 1.1.1)
3. **优先级** (P0/P1/P2)
4. **问题描述**
5. **修改建议**
6. **正确的代码示例**

### 示例输出格式 (中文输出)

```markdown
## 文件: test-cases/go-code-review/bad/user_service_bad.go

### 问题 1 - [P0] 规则 1.1.1 错误包装
**位置**: 第 35 行
**原始代码**:
\`\`\`go
return nil, fmt.Errorf("get user failed: %v", err)
\`\`\`
**问题描述**: 使用 fmt.Errorf 而不是 errors.Wrap,丢失了错误堆栈信息
**修改建议**:
\`\`\`go
return nil, errors.Wrapf(err, "failed to get user, user_id=%d", id)
\`\`\`

### 问题 2 - [P0] 规则 1.3.1 缺少 WHERE 条件
**位置**: 第 32 行
**原始代码**:
\`\`\`go
err := s.db.Find(&user).Error
\`\`\`
**问题描述**: 缺少显式 WHERE 条件,可能导致查询所有记录
**修改建议**:
\`\`\`go
err := s.db.Where("id = ?", id).Take(&user).Error
\`\`\`
```

## 测试检查清单

使用以下清单验证 Code Review Skill 的检测能力:

### P0 - 必须修复 (Critical)

- [ ] 检测到错误未正确包装 (Rule 1.1.1)
- [ ] 检测到 panic 使用 (Rule 1.1.3)
- [ ] 检测到错误未检查 (Rule 1.1.4)
- [ ] 检测到 nil 未检查 (Rule 1.2.1)
- [ ] 检测到缺少 WHERE 条件 (Rule 1.3.1)
- [ ] 检测到缺少 SELECT 列 (Rule 1.3.2)
- [ ] 检测到使用 Save() (Rule 1.3.4)
- [ ] 检测到使用 First() (Rule 1.3.5)
- [ ] 检测到缺少 gorm tag (Rule 1.3.6)
- [ ] 检测到未使用 ErrorGroup (Rule 1.4.1)
- [ ] 检测到未使用 limiter (Rule 1.4.2)
- [ ] 检测到手动构造 JSON (Rule 1.5.1)

### P1 - 强烈推荐 (Quality)

- [ ] 检测到命名不一致 (Rule 2.1.7)
- [ ] 检测到魔法数字 (Rule 2.1.6)
- [ ] 检测到日志格式不规范 (Rule 2.2.1, 2.2.2)
- [ ] 检测到缺少日志 (Rule 2.2.7)
- [ ] 检测到嵌入字段位置错误 (Rule 2.3.1)
- [ ] 检测到 import 分组不规范 (Rule 2.3.2)
- [ ] 检测到未使用 make() (Rule 2.3.3)
- [ ] 检测到缺少 default case (Rule 2.3.5)
- [ ] 检测到全局变量 (Rule 2.3.7)
- [ ] 检测到 TODO 未标注负责人 (Rule 2.3.10)

### P2 - 建议优化 (Style)

- [ ] 检测到拼写错误 (Rule 2.5.1)
- [ ] 检测到函数过长 (Rule 2.5.3)
- [ ] 检测到缺少注释 (Rule 2.5.4)
- [ ] 检测到枚举注释语言错误 (Rule 2.5.5)
- [ ] 检测到复杂条件未提取 (Rule 2.5.7)
- [ ] 检测到配置类型不当 (Rule 3.3.1)

## 统计信息

### user_service_bad.go
- **P0 问题**: 约 20+ 个
- **P1 问题**: 约 10+ 个
- **P2 问题**: 约 5+ 个
- **涉及规则**: 25+ 条

### project_structure_bad.go
- **P1 问题**: 约 8+ 个
- **P2 问题**: 约 6+ 个
- **涉及规则**: 15+ 条

### user_service_good.go
- **违规数**: 0
- **展示规则**: 30+ 条最佳实践

## 注意事项

1. **实际使用**: 这些测试文件**不能**直接编译运行,因为缺少依赖包
2. **目的**: 仅用于测试 Code Review Skill 的检测能力
3. **覆盖率**: 涵盖 FUTU Go Coding Standards 的大部分核心规则
4. **维护**: 随着规则更新,测试用例也需要相应更新

## 扩展测试

### 创建自定义测试

1. 复制 `bad/user_service_bad.go` 作为模板
2. 添加要测试的特定违规
3. 运行 Code Review
4. 验证是否正确检测

### 性能测试

```bash
# 测试大文件审查
time claude-code review bad/user_service_bad.go

# 测试批量文件审查
time claude-code review bad/*.go
```

## 反馈

如果发现 Code Review Skill:
- ✅ 正确检测到的问题
- ❌ 漏检的问题
- ⚠️ 误报的问题

请记录下来并反馈给开发团队,帮助改进检测准确性。

---

**创建日期**: 2025-12-19
**版本**: 1.0.0
**维护者**: Jason Ouyang
