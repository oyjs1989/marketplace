# Go Code Review Plugin

Comprehensive Go code review with 4 specialist agents and 137+ FUTU coding standards.

## Quick Start

Just trigger the main skill:
```
Review my Go code
```

The plugin automatically:
- Analyzes changed files
- Selects applicable agents
- Runs parallel reviews
- Delivers prioritized findings (in Chinese)

## Specialist Agents

This plugin includes 4 independent review agents:

- **gorm-review** (blue) - Database operations and GORM best practices (rules 1.3.*)
- **error-safety** (red) - Error handling and concurrency safety (rules 1.1.*, 1.2.*, 1.4.*, 1.5.*)
- **naming-logging** (green) - Naming conventions and logging standards (rules 2.1.*, 2.2.*)
- **organization** (purple) - Code organization and quality (rules 2.3.*, 2.4.*, 2.5.*, 3.*)

## Features

- **Smart Scope Detection**: Automatically determines which agents apply
- **Parallel Agent Execution**: All agents run simultaneously for speed
- **97+ Rules**: Comprehensive FUTU Go coding standards
- **Chinese Output**: All findings reported in Chinese (per FUTU requirement)
- **Priority Classification**: P0 (must fix), P1 (recommended), P2 (suggested)
- **Actionable Recommendations**: Specific fix suggestions with code examples

## Architecture

```
go-code-review/
├── .claude-plugin/
│   └── plugin.json              # Plugin metadata
├── SKILL.md                     # Main orchestrator
├── agents/                      # 4 specialist agents
│   ├── gorm-review.md
│   ├── error-safety.md
│   ├── naming-logging.md
│   └── organization.md
└── references/
    └── FUTU_GO_STANDARDS.md    # 97+ coding standards
```

## Usage Examples

### Automatic Review (Recommended)
```
Review my Go code
```
Automatically selects and runs applicable agents.

### Manual Agent Invocation
```
Use gorm-review agent
Use error-safety agent
Use naming-logging agent
Use organization agent
```
Explicitly call specific agents when needed.

### PR Review
```
Review this PR
```
Reviews all changes in current pull request.

## Rule Categories

### P0 - Must Fix (Critical Issues)
- Missing or improper error handling
- Unchecked nil pointers
- Improper database operations
- Concurrency safety issues
- Business logic errors
- Unused struct fields

### P1 - Strongly Recommended (Code Quality)
- Non-standard naming conventions
- Improper logging
- Missing necessary comments
- Poor function encapsulation
- Code duplication

### P2 - Suggested Optimization (Style & Best Practices)
- Inconsistent formatting
- Readability improvements
- Performance optimization opportunities

## Output Format

All findings are reported in Chinese (per FUTU standards):

```markdown
## 文件: path/to/file.go

### 问题 1 - [P0] 规则 1.3.1
**位置**: 第 42 行
**类别**: GORM/数据库
**原始代码**:
```go
db.Find(&users)
```
**问题描述**: 缺少WHERE条件和显式列选择
**修改建议**: 添加WHERE条件和Select指定列...
```

## Testing

Test with bad code examples:
```
Review test-cases/go-code-review/bad/user_service_bad.go
```

Run automated validation:
```bash
cd test-cases/go-code-review
./run_test.sh
./validate_results.sh
```

## Version History

- **v3.0.0**: Unified plugin with automatic agent orchestration
- **v2.0.0**: Multi-skill architecture with 5 independent skills
- **v1.0.0**: Initial monolithic skill

## Standards Reference

Complete coding standards documentation: `references/FUTU_GO_STANDARDS.md`

## Support

For issues or questions, refer to the marketplace documentation or FUTU development team.
