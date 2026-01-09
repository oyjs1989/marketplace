# Cursor Skills 使用指南

## GitLab评论提炼规则 Skill

### 简介
这个skill可以从GitLab Merge Request的评论中自动提取代码审查问题，分析问题背后的逻辑，提炼成通用的代码设计原则和规则。

### 快速开始

#### 基本用法
```
从GitLab评论提炼规则，URL: https://gitlab.futunn.com/fcmp/resource_center/-/merge_requests/39
```

#### 指定输出文件
```
从GitLab评论提炼规则，URL: https://gitlab.futunn.com/fcmp/resource_center/-/merge_requests/39，输出到: my_rules.md
```

#### 只生成通用原则
```
从GitLab评论提炼规则，URL: https://gitlab.futunn.com/fcmp/resource_center/-/merge_requests/39，规则层次: general
```

### 参数说明

- **gitlab_url** (必需): GitLab MR的完整URL
- **output_file** (可选): 输出文件路径，默认为 `code_rules_from_mr.md`
- **rule_level** (可选): 规则提炼层次
  - `specific`: 只生成具体规则（保留原始问题和解决方案）
  - `general`: 只生成通用原则（提炼背后的逻辑）
  - `both`: 两者都包含（默认）

### 工作流程

1. **访问GitLab页面** - 自动打开并登录（如需要）
2. **提取评论** - 从页面中提取所有讨论和评论
3. **分析模式** - 识别问题和解决方案
4. **分类整理** - 按主题分类评论
5. **提炼规则** - 从具体问题提炼通用原则
6. **生成文档** - 输出结构化的规则文档

### 输出示例

生成的文档包含：
- 通用代码设计原则（按分类组织）
- 每个原则包含：
  - 核心思想
  - 具体规则
  - 代码示例（正确✅和错误❌）
  - 原始问题和解决方案（如果选择specific或both）
- 实践检查清单

### 注意事项

1. **权限要求**: 需要GitLab访问权限，可能需要登录
2. **处理时间**: 大量评论可能需要较长时间
3. **人工审核**: 某些规则可能需要人工审核和补充
4. **评论质量**: 规则提炼的质量取决于评论的详细程度

### 技术实现

- 使用浏览器MCP工具访问GitLab
- JavaScript提取评论数据
- Python分析器提炼规则
- Markdown格式化输出

### 扩展功能

未来可能支持：
- 多个MR合并提炼
- 规则去重和合并
- 规则库持续更新
- 不同格式输出（JSON、YAML等）
- 规则验证和测试用例生成

