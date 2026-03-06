# Go Code Review v5.0 工具链升级实现计划

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 用 `go build`/`go vet`/`staticcheck` 替换 grep 正则扫描，消除 90%+ 假阳性，让 Agent 聚焦真正有价值的判断性分析。

**Architecture:**
- 废弃 Tier 2（scan-rules.sh grep 扫描），合并为单一工具 `run-go-tools.sh`
- 新 Tier 1 = 编译器诊断（`go build` + `go vet`，0 假阳性）
- 新 Tier 2 = SSA 静态分析（`staticcheck`，可选，高置信）
- YAML 规则修复最严重的假阳性作为兜底（保留工具但修正 pattern）
- Agent 输入从"557条命中过滤"变为"10条高置信命中确认"

**Tech Stack:** bash, go vet (内置), staticcheck (可选安装), JSON

---

## Task 1: 创建 `tools/run-go-tools.sh`

**Files:**
- Create: `skills/go-code-review/tools/run-go-tools.sh`

**背景：** 替换现有两个工具脚本。接收 Go 文件路径（stdin），提取 package 路径，运行 `go vet` 和可选的 `staticcheck`，输出统一的 `diagnostics.json`。

**Step 1: 创建脚本文件**

```bash
cat > skills/go-code-review/tools/run-go-tools.sh << 'EOF'
#!/usr/bin/env bash
# run-go-tools.sh - Run Go compiler tools on changed files
# Usage:
#   git diff HEAD~1 --name-only --diff-filter=AM | grep '\.go$' | bash run-go-tools.sh
#   echo "service/user.go" | bash run-go-tools.sh
#
# Output: JSON to stdout
# {
#   "build_errors": [...],
#   "vet_issues": [...],
#   "staticcheck_issues": [...],
#   "large_files": [...],
#   "summary": {"build_errors": 0, "vet_issues": 0, "staticcheck_issues": 0}
# }

set -euo pipefail

# ---------------------------------------------------------------------------
# Collect file paths from stdin
# ---------------------------------------------------------------------------
declare -a GO_FILES=()
while IFS= read -r line; do
    [[ -n "$line" ]] && GO_FILES+=("$line")
done

if [[ ${#GO_FILES[@]} -eq 0 ]]; then
    echo '{"build_errors":[],"vet_issues":[],"staticcheck_issues":[],"large_files":[],"summary":{"build_errors":0,"vet_issues":0,"staticcheck_issues":0}}'
    exit 0
fi

# ---------------------------------------------------------------------------
# Extract unique packages from file paths
# e.g. "service/user.go" -> "./service"
# ---------------------------------------------------------------------------
declare -A SEEN_PKGS
declare -a PACKAGES=()
for f in "${GO_FILES[@]}"; do
    pkg="./$(dirname "$f")"
    if [[ -z "${SEEN_PKGS[$pkg]+x}" ]]; then
        SEEN_PKGS["$pkg"]=1
        PACKAGES+=("$pkg")
    fi
done

# ---------------------------------------------------------------------------
# Helper: json_escape
# ---------------------------------------------------------------------------
json_escape() {
    printf '%s' "$1" | python3 -c "import json,sys; print(json.dumps(sys.stdin.read()))" 2>/dev/null \
    || printf '%s' "$1" | sed 's/\\/\\\\/g; s/"/\\"/g; s/$/\\n/g' | tr -d '\n' | sed 's/\\n$//'
}

# ---------------------------------------------------------------------------
# Tier 0: go build - compilation errors
# ---------------------------------------------------------------------------
BUILD_ERRORS_JSON="[]"
BUILD_COUNT=0

if command -v go &>/dev/null; then
    BUILD_OUT=$(go build "${PACKAGES[@]}" 2>&1 || true)
    if [[ -n "$BUILD_OUT" ]]; then
        ENTRIES=""
        while IFS= read -r line; do
            [[ -z "$line" ]] && continue
            # Format: "file.go:line:col: message"
            if [[ "$line" =~ ^([^:]+):([0-9]+):([0-9]+):[[:space:]](.+)$ ]]; then
                file="${BASH_REMATCH[1]}"
                lineno="${BASH_REMATCH[2]}"
                msg="${BASH_REMATCH[4]}"
                entry=$(printf '{"file":"%s","line":%s,"message":"%s","severity":"P0"}' \
                    "$(json_escape "$file")" "$lineno" "$(json_escape "$msg")")
                ENTRIES="${ENTRIES:+$ENTRIES,}$entry"
                BUILD_COUNT=$((BUILD_COUNT + 1))
            fi
        done <<< "$BUILD_OUT"
        [[ -n "$ENTRIES" ]] && BUILD_ERRORS_JSON="[$ENTRIES]"
    fi
fi

# ---------------------------------------------------------------------------
# Tier 1: go vet - type-level issues (built-in, ~0 false positives)
# ---------------------------------------------------------------------------
VET_ISSUES_JSON="[]"
VET_COUNT=0

if command -v go &>/dev/null; then
    VET_OUT=$(go vet "${PACKAGES[@]}" 2>&1 || true)
    if [[ -n "$VET_OUT" ]]; then
        ENTRIES=""
        while IFS= read -r line; do
            [[ -z "$line" ]] && continue
            [[ "$line" == *"#"* ]] && continue  # skip package header lines
            if [[ "$line" =~ ^([^:]+):([0-9]+):([0-9]+):[[:space:]](.+)$ ]]; then
                file="${BASH_REMATCH[1]}"
                lineno="${BASH_REMATCH[2]}"
                msg="${BASH_REMATCH[4]}"
                entry=$(printf '{"file":"%s","line":%s,"message":"%s","severity":"P1"}' \
                    "$(json_escape "$file")" "$lineno" "$(json_escape "$msg")")
                ENTRIES="${ENTRIES:+$ENTRIES,}$entry"
                VET_COUNT=$((VET_COUNT + 1))
            fi
        done <<< "$VET_OUT"
        [[ -n "$ENTRIES" ]] && VET_ISSUES_JSON="[$ENTRIES]"
    fi
fi

# ---------------------------------------------------------------------------
# Tier 2: staticcheck - SSA analysis (optional, install: go install honnef.co/go/tools/cmd/staticcheck@latest)
# ---------------------------------------------------------------------------
STATICCHECK_JSON="[]"
STATICCHECK_COUNT=0

if command -v staticcheck &>/dev/null; then
    SC_OUT=$(staticcheck "${PACKAGES[@]}" 2>&1 || true)
    if [[ -n "$SC_OUT" ]]; then
        ENTRIES=""
        while IFS= read -r line; do
            [[ -z "$line" ]] && continue
            # Format: "file.go:line:col: CODE message"
            if [[ "$line" =~ ^([^:]+):([0-9]+):([0-9]+):[[:space:]]([A-Z]+[0-9]+):[[:space:]](.+)$ ]]; then
                file="${BASH_REMATCH[1]}"
                lineno="${BASH_REMATCH[2]}"
                code="${BASH_REMATCH[4]}"
                msg="${BASH_REMATCH[5]}"
                # Map staticcheck severity: S1xxx=style(P2), SA=bug(P0/P1)
                sev="P1"
                [[ "$code" == SA* ]] && sev="P0"
                [[ "$code" == S1* ]] && sev="P2"
                entry=$(printf '{"file":"%s","line":%s,"code":"%s","message":"%s","severity":"%s"}' \
                    "$(json_escape "$file")" "$lineno" "$code" "$(json_escape "$msg")" "$sev")
                ENTRIES="${ENTRIES:+$ENTRIES,}$entry"
                STATICCHECK_COUNT=$((STATICCHECK_COUNT + 1))
            fi
        done <<< "$SC_OUT"
        [[ -n "$ENTRIES" ]] && STATICCHECK_JSON="[$ENTRIES]"
    fi
fi

# ---------------------------------------------------------------------------
# Large files (>800 lines) - kept as simple size heuristic
# ---------------------------------------------------------------------------
LARGE_FILES_JSON="[]"
LARGE_ENTRIES=""
for f in "${GO_FILES[@]}"; do
    [[ ! -f "$f" ]] && continue
    lines=$(wc -l < "$f" 2>/dev/null || echo 0)
    if [[ "$lines" -gt 800 ]]; then
        entry=$(printf '{"file":"%s","lines":%s}' "$(json_escape "$f")" "$lines")
        LARGE_ENTRIES="${LARGE_ENTRIES:+$LARGE_ENTRIES,}$entry"
    fi
done
[[ -n "$LARGE_ENTRIES" ]] && LARGE_FILES_JSON="[$LARGE_ENTRIES]"

# ---------------------------------------------------------------------------
# Output
# ---------------------------------------------------------------------------
printf '{"build_errors":%s,"vet_issues":%s,"staticcheck_issues":%s,"large_files":%s,"summary":{"build_errors":%d,"vet_issues":%d,"staticcheck_issues":%d}}\n' \
    "$BUILD_ERRORS_JSON" "$VET_ISSUES_JSON" "$STATICCHECK_JSON" "$LARGE_FILES_JSON" \
    "$BUILD_COUNT" "$VET_COUNT" "$STATICCHECK_COUNT"
EOF
chmod +x skills/go-code-review/tools/run-go-tools.sh
```

**Step 2: 验证脚本语法**

```bash
bash -n skills/go-code-review/tools/run-go-tools.sh && echo "syntax OK"
```
Expected: `syntax OK`

**Step 3: 冒烟测试（在任意 Go 项目目录下）**

```bash
# 在一个有 Go 代码的目录
git diff HEAD~1 --name-only --diff-filter=AM | grep '\.go$' | bash /path/to/run-go-tools.sh | python3 -m json.tool | head -30
```
Expected: 合法 JSON，`summary` 字段显示各工具命中数

**Step 4: Commit**

```bash
git add skills/go-code-review/tools/run-go-tools.sh
git commit -m "feat(go-code-review): add run-go-tools.sh using go vet + staticcheck"
```

---

## Task 2: 修复 YAML 规则最严重的假阳性

**Files:**
- Modify: `skills/go-code-review/rules/safety.yaml`
- Modify: `skills/go-code-review/rules/data.yaml`
- Modify: `skills/go-code-review/rules/observability.yaml`

**背景：** 三个最严重的假阳性来源，通过精化 pattern 解决：
- SAFE-001：`fmt.Errorf("%w", err)` 是合法的 Go 错误包装，不应命中
- SAFE-002：测试文件中的 `panic()` 是正常用法
- DATA-006：for-range 循环体内没有 db 写操作时不应命中

**Step 1: 修复 SAFE-001 - 排除 `%w` 合法包装**

在 `rules/safety.yaml` 中找到 SAFE-001，将 pattern.match 从：
```yaml
match: 'fmt\.Errorf\('
```
改为（排除含 `%w` 的行，只匹配无法包装错误的新建错误）：
```yaml
match: 'fmt\.Errorf\((?![^)]*%w)'
```

**Step 2: 修复 SAFE-002 - 排除 `_test.go` 文件**

SAFE-002 的 pattern 无法在 grep 层面区分测试文件，在 YAML 中添加 exclude 注释供 Agent 参考：
```yaml
  - id: SAFE-002
    name: 业务函数禁用 panic
    severity: P0
    pattern:
      type: regex
      match: '\bpanic\('
    exclude_files: ['_test.go']   # Agent 需跳过测试文件的命中
    message: "禁止在业务逻辑中使用 panic()。注意：测试辅助代码（_test.go）中的 panic 属正常用法，应忽略。"
```

**Step 3: 修复 DATA-006 - 要求循环体内有 DB 操作**

将 DATA-006 的 match 精化，仅匹配紧跟 db/DB/repo 调用的 for-range：
当前：`match: 'for\s+\w+.*:=\s+range'`

添加备注：此规则误报率极高，Agent 需人工确认循环体内是否有实际 DB 写操作。将 severity 从 P1 降为 P2（降低噪音）：
```yaml
    severity: P2   # 改为 P2，Agent 确认循环体有 db 写操作后才升级
    message: "检测到 for-range 循环——需确认循环体内是否有 DB 写操作（N+1 问题）。如无 DB 调用请忽略此命中。"
```

**Step 4: 降低 OBS-006/007 误报（注释代码）**

OBS-006/007 pattern 命中注释行。更新 message 提醒 Agent：
```yaml
    message: "检测到 fmt.Print* 调用。注意：如该行以 // 开头（注释），请忽略此命中。"
```

**Step 5: 验证 YAML 仍可被 scan-rules.sh 解析**

```bash
echo "test/any_file.go" | bash skills/go-code-review/tools/scan-rules.sh skills/go-code-review/rules | python3 -m json.tool | tail -5
```
Expected: 合法 JSON，无 parse 错误

**Step 6: Commit**

```bash
git add skills/go-code-review/rules/
git commit -m "fix(go-code-review): reduce false positives in SAFE-001/002, DATA-006, OBS-006/007"
```

---

## Task 3: 更新 SKILL.md 工作流

**Files:**
- Modify: `skills/go-code-review/SKILL.md`

**背景：** SKILL.md 是编排器，需要反映新的工具链。新的 Step 2 运行 `run-go-tools.sh`，Step 3 保留 `scan-rules.sh` 作为补充（但期望命中数大幅减少）。

**Step 1: 更新架构图**

将架构图中 Tier 1/2 描述改为：

```
┌─────────────────────────────────────────────────────┐
│  Tier 1: tools/run-go-tools.sh                      │  → diagnostics.json
│  go build（编译错误）+ go vet（类型检查）             │    (P0: 编译/类型问题)
│  + staticcheck（SSA分析，可选）                      │    (P1: 高置信逻辑问题)
└─────────────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────┐
│  Tier 2: tools/scan-rules.sh                        │  → rule-hits.json
│  修复后的 YAML 规则（兜底）                           │    (预期 <50 条，假阳性大幅减少)
└─────────────────────────────────────────────────────┘
```

**Step 2: 更新 Step 2（原 Tier 1 量化分析）**

将原来的 analyze-go.sh 调用替换为：

```bash
# Step 2: 运行 Tier 1 工具链分析
git diff master --name-only --diff-filter=AM | grep '\.go$' | bash tools/run-go-tools.sh > /tmp/diagnostics.json
```

输出字段说明：
- `build_errors`：编译错误（P0，必须修复）
- `vet_issues`：go vet 发现的类型/格式问题（P0/P1）
- `staticcheck_issues`：SSA 分析结果（SA* 为 P0，S1* 为 P2）
- `large_files`：行数 > 800 的文件（参考）

**Step 3: 更新 Tier 3 Agent 调用说明**

修改 safety/quality agent 的输入来源：

```
- **safety agent** — 读取 diagnostics.json（build_errors + vet_issues + staticcheck SA*）；
  同时确认 rule-hits.json 中 SAFE-001~010 的命中（过滤假阳性）
- **quality agent** — 读取 diagnostics.json（large_files）；
  确认 rule-hits.json 中 QUAL-001~010 命中 + diagnostics.json 的 S1* 代码风格建议
```

**Step 4: 更新版本号**

```yaml
version: 5.0.0
```

同时更新输出报告的版本：
```markdown
# Go 代码审查报告（v5.0.0）
```

**Step 5: Commit**

```bash
git add skills/go-code-review/SKILL.md
git commit -m "feat(go-code-review): update orchestrator to use run-go-tools.sh (v5.0.0)"
```

---

## Task 4: 更新 safety agent 使用新工具输出

**Files:**
- Modify: `skills/go-code-review/agents/safety.md`

**背景：** safety agent 目前被要求"确认 Tier 2 命中"，现在改为"优先读取 diagnostics.json，再补充判断性分析"。

**Step 1: 更新工具输入部分**

在 agent 中将：
```
读取 /tmp/rule-hits.json，筛选 SAFE-* 命中
```
改为：
```
优先读取 /tmp/diagnostics.json：
  - build_errors → 所有编译错误直接报告为 P0
  - vet_issues 含 "copylocks"/"printf"/"assign" → 报告为 P0/P1
  - staticcheck_issues 含 SA* 代码 → 报告为 P0
补充读取 /tmp/rule-hits.json：
  - 筛选 SAFE-* 命中，人工判断是否假阳性后再报告
```

**Step 2: 在 agent 中明确假阳性过滤规则**

添加段落：
```markdown
## 假阳性过滤规则

- SAFE-001 命中但匹配行含 `%w` → 忽略（fmt.Errorf("%w", err) 是正确包装）
- SAFE-002 命中但文件名以 `_test.go` 结尾 → 忽略（测试辅助 panic 是正常用法）
- vet_issues 中同一问题已由 build_errors 覆盖 → 去重，只报一次
```

**Step 3: Commit**

```bash
git add skills/go-code-review/agents/safety.md
git commit -m "fix(go-code-review): update safety agent to use diagnostics.json, add false-positive filters"
```

---

## Task 5: 更新 quality agent 使用新工具输出

**Files:**
- Modify: `skills/go-code-review/agents/quality.md`

**背景：** quality agent 目前使用 metrics.json（来自 analyze-go.sh），需要改为使用 diagnostics.json。

**Step 1: 更新工具读取说明**

将所有 `metrics.json` 引用改为 `diagnostics.json`：

```markdown
读取 /tmp/diagnostics.json：
- large_files → 超过 800 行的文件（列出文件名和行数）
- staticcheck_issues 含 S1* 代码 → P2 代码风格建议
不再读取 /tmp/metrics.json（已废弃）
```

**Step 2: 移除"函数行数 > 80 行违规"的硬性报告**

当前 quality agent 会报告所有函数 > 80 行的情况。根据分析，这些大多是业务复杂度导致的合理情况，不是问题。改为：

```markdown
函数行数 > 80 行：仅在**同时满足**以下条件时报告为 P2：
  1. 函数行数 > 80 行
  2. 函数内有明显的可拆分逻辑（多个独立的 if-else 块）
  纯业务流程的长函数（如状态机、复杂查询构建）不应报告
```

**Step 3: Commit**

```bash
git add skills/go-code-review/agents/quality.md
git commit -m "fix(go-code-review): update quality agent to use diagnostics.json, soften function-length rule"
```

---

## Task 6: 更新版本文件

**Files:**
- Modify: `skills/go-code-review/.claude-plugin/plugin.json`
- Modify: `.claude-plugin/marketplace.json`
- 同步更新 `~/.claude/plugins/cache/jasonouyang-marketplace/go-code-review/4.0.0/` 缓存（或重装插件）

**Step 1: 更新 plugin.json**

```json
{
  "name": "go-code-review",
  "version": "5.0.0",
  "description": "Go code review with three-tier architecture: go build/vet/staticcheck + rule scanning + 5 domain-expert AI agents",
  "author": {
    "name": "jasonouyang",
    "email": "jasonouyang@futunn.com"
  }
}
```

**Step 2: 更新 marketplace.json**

```json
{
  "name": "jasonouyang-marketplace",
  "plugins": [
    {
      "name": "go-code-review",
      "source": "./skills/go-code-review",
      "version": "5.0.0",
      "description": "Go code review: go build/vet/staticcheck + YAML rule scanning + 5 domain-expert AI agents",
      "author": {
        "name": "jasonouyang",
        "email": "jasonouyang@futunn.com"
      }
    }
  ]
}
```

**Step 3: Final commit + 同步缓存**

```bash
git add skills/go-code-review/.claude-plugin/plugin.json .claude-plugin/marketplace.json
git commit -m "chore(go-code-review): bump version to 5.0.0"

# 同步到插件缓存（重新安装会更干净）
claude plugin uninstall go-code-review@jasonouyang-marketplace
claude plugin install go-code-review@jasonouyang-marketplace
```

---

## 预期效果对比

| 指标 | v4.0 (当前) | v5.0 (目标) |
|------|------------|------------|
| Tier 工具命中数 | 557 条 | ~10-20 条 |
| 假阳性率 | ~80% | <5% |
| go vet 覆盖 | 无 | ✅ printf格式/copylock/assign |
| staticcheck 覆盖 | 无 | ✅ SA4006(值拷贝bug) / SA4023 |
| 编译错误检测 | 无 | ✅ go build 直接报告 |
| Agent token 浪费 | 高（过滤假阳性）| 低（确认高置信命中）|
| SAFE-001 假阳性 | `%w` 全部命中 | ✅ 排除 |
| DATA-006 假阳性 | 7处全假阳性 | ✅ 降为 P2 |

**说明：** `staticcheck SA4006`（赋值给局部变量但未使用）能直接发现上次审查中的 P0-3 值拷贝 bug（`UpdateDepartStatement` 修改 slice 元素但值被丢弃）。

---

## 安装 staticcheck（可选但推荐）

```bash
go install honnef.co/go/tools/cmd/staticcheck@latest
```

如未安装，`run-go-tools.sh` 会自动跳过 staticcheck 步骤，仍使用 `go build` + `go vet`。
