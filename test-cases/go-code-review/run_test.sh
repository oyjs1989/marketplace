#!/bin/bash
# Go Code Review v4.0.0 测试脚本
# 验证三层架构：Tier 1（量化）+ Tier 2（规则扫描）+ Tier 3（AI 审查指南）

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MARKETPLACE_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
SKILL_DIR="$MARKETPLACE_DIR/skills/go-code-review"
BAD_DIR="$SCRIPT_DIR/bad"
GOOD_DIR="$SCRIPT_DIR/good"
ANALYZE_SH="$SKILL_DIR/tools/analyze-go.sh"
SCAN_SH="$SKILL_DIR/tools/scan-rules.sh"
RULES_DIR="$SKILL_DIR/rules"
METRICS_OUT="/tmp/go_review_metrics_$$.json"
HITS_OUT="/tmp/go_review_hits_$$.json"
RESULT_FILE="$SCRIPT_DIR/test_result.md"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Go Code Review v4.0.0 - 自动化测试${NC}"
echo -e "${BLUE}  三层架构验证${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 函数: 打印彩色消息
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

# 计数器
tier1_passed=0
tier1_failed=0
tier2_passed=0
tier2_failed=0

# ============================================================
# 步骤 1: 验证工具文件和测试目录存在
# ============================================================
echo "步骤 1: 检查测试环境和工具文件..."

if [ ! -d "$BAD_DIR" ]; then
    print_error "找不到测试目录: $BAD_DIR"
    exit 1
fi

if [ ! -d "$GOOD_DIR" ]; then
    print_warning "找不到正确实现目录: $GOOD_DIR"
fi

print_success "测试目录准备就绪"

# 检查工具文件
if [ ! -f "$ANALYZE_SH" ]; then
    print_warning "未找到 Tier 1 工具: $ANALYZE_SH"
    print_info "Tier 1 测试将被跳过"
    SKIP_TIER1=1
else
    print_success "Tier 1 工具: $ANALYZE_SH"
    SKIP_TIER1=0
fi

if [ ! -f "$SCAN_SH" ]; then
    print_warning "未找到 Tier 2 工具: $SCAN_SH"
    print_info "Tier 2 测试将被跳过"
    SKIP_TIER2=1
else
    print_success "Tier 2 工具: $SCAN_SH"
    SKIP_TIER2=0
fi

if [ ! -d "$RULES_DIR" ]; then
    print_warning "未找到规则目录: $RULES_DIR"
    SKIP_TIER2=1
else
    RULE_COUNT=$(find "$RULES_DIR" -name "*.yaml" 2>/dev/null | wc -l | tr -d ' ')
    print_success "规则目录: $RULES_DIR ($RULE_COUNT 个 YAML 文件)"
fi

echo ""

# 扫描测试文件
BAD_FILES=($(find "$BAD_DIR" -name "*.go" -type f 2>/dev/null))
print_info "发现 ${#BAD_FILES[@]} 个违规测试文件:"
for f in "${BAD_FILES[@]}"; do
    echo "  - $(basename $f)"
done
echo ""

# ============================================================
# 步骤 2: 运行 Tier 1 分析 (analyze-go.sh)
# ============================================================
echo "步骤 2: 运行 Tier 1 分析 (analyze-go.sh)..."

if [ "${SKIP_TIER1:-0}" = "1" ]; then
    print_warning "跳过 Tier 1（工具不存在）"
else
    # 对 bad/*.go 运行 analyze-go.sh
    BAD_GO_FILES=()
    for f in "${BAD_FILES[@]}"; do
        BAD_GO_FILES+=("$f")
    done

    if [ ${#BAD_GO_FILES[@]} -eq 0 ]; then
        print_warning "bad/ 目录中没有 .go 文件"
    else
        print_info "分析文件: ${BAD_GO_FILES[*]}"
        if bash "$ANALYZE_SH" "${BAD_GO_FILES[@]}" > "$METRICS_OUT" 2>&1; then
            # 验证输出 JSON 格式
            if python3 -c "import json,sys; d=json.load(open('$METRICS_OUT')); print('summary:', d.get('summary', {}))" 2>/dev/null; then
                print_success "Tier 1 分析完成，metrics.json 格式有效"
                # 报告关键指标
                FILES_OVER=$(python3 -c "import json; d=json.load(open('$METRICS_OUT')); print(d.get('summary',{}).get('files_over_800',0))" 2>/dev/null || echo "N/A")
                FUNCS_OVER=$(python3 -c "import json; d=json.load(open('$METRICS_OUT')); print(d.get('summary',{}).get('functions_over_80',0))" 2>/dev/null || echo "N/A")
                NESTING=$(python3 -c "import json; d=json.load(open('$METRICS_OUT')); print(d.get('summary',{}).get('nesting_violations',0))" 2>/dev/null || echo "N/A")
                print_info "  files_over_800: $FILES_OVER"
                print_info "  functions_over_80: $FUNCS_OVER"
                print_info "  nesting_violations: $NESTING"
                tier1_passed=$((tier1_passed + 1))
            else
                # 备用：用 python3 -m json.tool 验证
                if cat "$METRICS_OUT" | python3 -m json.tool > /dev/null 2>&1; then
                    print_success "Tier 1 分析完成，JSON 格式有效"
                    tier1_passed=$((tier1_passed + 1))
                else
                    print_error "Tier 1 输出 JSON 格式无效"
                    print_info "输出内容: $(head -5 $METRICS_OUT 2>/dev/null)"
                    tier1_failed=$((tier1_failed + 1))
                fi
            fi
        else
            print_error "analyze-go.sh 执行失败"
            print_info "错误输出: $(cat $METRICS_OUT 2>/dev/null | head -10)"
            tier1_failed=$((tier1_failed + 1))
        fi
    fi
fi
echo ""

# ============================================================
# 步骤 3: 运行 Tier 2 规则扫描 (scan-rules.sh)
# ============================================================
echo "步骤 3: 运行 Tier 2 规则扫描 (scan-rules.sh)..."

if [ "${SKIP_TIER2:-0}" = "1" ]; then
    print_warning "跳过 Tier 2（工具或规则目录不存在）"
else
    BAD_GO_FILES=()
    for f in "${BAD_FILES[@]}"; do
        BAD_GO_FILES+=("$f")
    done

    if [ ${#BAD_GO_FILES[@]} -eq 0 ]; then
        print_warning "bad/ 目录中没有 .go 文件"
    else
        print_info "扫描文件: ${BAD_GO_FILES[*]}"
        if bash "$SCAN_SH" "$RULES_DIR" "${BAD_GO_FILES[@]}" > "$HITS_OUT" 2>&1; then
            # 验证输出 JSON 格式
            if cat "$HITS_OUT" | python3 -m json.tool > /dev/null 2>&1; then
                print_success "Tier 2 扫描完成，rule-hits.json 格式有效"
                # 报告命中统计
                TOTAL=$(python3 -c "import json; d=json.load(open('$HITS_OUT')); print(d.get('summary',{}).get('total',0))" 2>/dev/null || echo "N/A")
                P0_COUNT=$(python3 -c "import json; d=json.load(open('$HITS_OUT')); print(d.get('summary',{}).get('p0_count',0))" 2>/dev/null || echo "N/A")
                P1_COUNT=$(python3 -c "import json; d=json.load(open('$HITS_OUT')); print(d.get('summary',{}).get('p1_count',0))" 2>/dev/null || echo "N/A")
                print_info "  总命中数: $TOTAL"
                print_info "  P0 命中: $P0_COUNT"
                print_info "  P1 命中: $P1_COUNT"
                tier2_passed=$((tier2_passed + 1))
            else
                print_error "Tier 2 输出 JSON 格式无效"
                print_info "输出内容: $(head -5 $HITS_OUT 2>/dev/null)"
                tier2_failed=$((tier2_failed + 1))
            fi
        else
            print_error "scan-rules.sh 执行失败"
            print_info "错误输出: $(cat $HITS_OUT 2>/dev/null | head -10)"
            tier2_failed=$((tier2_failed + 1))
        fi
    fi
fi
echo ""

# ============================================================
# 步骤 4: 指导 AI 审查（Tier 3）
# ============================================================
echo "步骤 4: Tier 3 AI 审查指南..."
echo ""
print_info "Tier 3 需要在 Claude Code 中手动运行以下命令："
echo ""
echo "  # 完整审查所有违规测试文件"
echo "  Review test-cases/go-code-review/bad/*.go"
echo ""
echo "  # 审查单个文件"
echo "  Review test-cases/go-code-review/bad/user_service_bad.go"
echo "  Review test-cases/go-code-review/bad/project_structure_bad.go"
echo "  Review test-cases/go-code-review/bad/design_philosophy_bad.go"
echo "  Review test-cases/go-code-review/bad/early_return_bad.go"
echo ""
print_info "预期行为："
echo "  - 5 个领域专家 agent 并行运行（safety/data/design/quality/observability）"
echo "  - 发现问题按 P0/P1/P2 分级，输出中文"
echo "  - 引用 SAFE-/DATA-/QUAL-/OBS- 规则 ID"
echo ""

# ============================================================
# 步骤 5: 汇总报告
# ============================================================
echo "步骤 5: 生成测试报告..."

cat > "$RESULT_FILE" << EOF
# Go Code Review v4.0.0 - 测试报告

生成时间: $(date '+%Y-%m-%d %H:%M:%S')

## 三层架构测试结果

| 层级 | 状态 | 说明 |
|------|------|------|
| Tier 1 (量化分析) | $([ "${SKIP_TIER1:-0}" = "1" ] && echo "跳过" || ([ $tier1_failed -eq 0 ] && echo "通过" || echo "失败")) | analyze-go.sh → metrics.json |
| Tier 2 (规则扫描) | $([ "${SKIP_TIER2:-0}" = "1" ] && echo "跳过" || ([ $tier2_failed -eq 0 ] && echo "通过" || echo "失败")) | scan-rules.sh → rule-hits.json |
| Tier 3 (AI 审查) | 待手动执行 | 5 个领域专家 agent |

## 测试文件

$(for f in "${BAD_FILES[@]}"; do echo "- $(basename $f)"; done)

## Tier 1 指标

$(if [ -f "$METRICS_OUT" ] && cat "$METRICS_OUT" | python3 -m json.tool > /dev/null 2>&1; then
    python3 -c "
import json
d = json.load(open('$METRICS_OUT'))
s = d.get('summary', {})
print('- files_over_800:', s.get('files_over_800', 'N/A'))
print('- functions_over_80:', s.get('functions_over_80', 'N/A'))
print('- nesting_violations:', s.get('nesting_violations', 'N/A'))
" 2>/dev/null || echo "- 详细指标见 $METRICS_OUT"
else
    echo "- Tier 1 工具未运行或输出无效"
fi)

## Tier 2 规则命中

$(if [ -f "$HITS_OUT" ] && cat "$HITS_OUT" | python3 -m json.tool > /dev/null 2>&1; then
    python3 -c "
import json
d = json.load(open('$HITS_OUT'))
s = d.get('summary', {})
print('- 总命中:', s.get('total', 'N/A'))
print('- P0:', s.get('p0_count', 'N/A'))
print('- P1:', s.get('p1_count', 'N/A'))
" 2>/dev/null || echo "- 详细命中见 $HITS_OUT"
else
    echo "- Tier 2 工具未运行或输出无效"
fi)

## Tier 3 手动审查命令

\`\`\`
Review test-cases/go-code-review/bad/*.go
\`\`\`
EOF

print_success "测试报告已保存: $RESULT_FILE"
echo ""

# ============================================================
# 最终统计
# ============================================================
echo -e "${BLUE}==========================================${NC}"
echo -e "${BLUE}  测试统计汇总${NC}"
echo -e "${BLUE}==========================================${NC}"
echo ""

if [ "${SKIP_TIER1:-0}" = "1" ]; then
    print_warning "Tier 1 (analyze-go.sh): 已跳过（工具不存在）"
elif [ $tier1_failed -eq 0 ]; then
    print_success "Tier 1 (analyze-go.sh): 通过"
else
    print_error "Tier 1 (analyze-go.sh): 失败"
fi

if [ "${SKIP_TIER2:-0}" = "1" ]; then
    print_warning "Tier 2 (scan-rules.sh): 已跳过（工具不存在）"
elif [ $tier2_failed -eq 0 ]; then
    print_success "Tier 2 (scan-rules.sh): 通过"
else
    print_error "Tier 2 (scan-rules.sh): 失败"
fi

print_info "Tier 3 (AI 审查): 需手动在 Claude Code 中执行"
echo ""

# 清理临时文件
rm -f "$METRICS_OUT" "$HITS_OUT"

# 退出代码
TOTAL_FAILED=$((tier1_failed + tier2_failed))
if [ $TOTAL_FAILED -gt 0 ]; then
    print_error "测试完成，存在 $TOTAL_FAILED 个失败项"
    exit 1
else
    print_success "测试完成"
    exit 0
fi
