#!/bin/bash

# 验证 Code Review 结果的脚本
# 用途: 对比实际检测结果和预期结果，识别漏检和误报

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
RESULT_FILE="${1:-code_review.result}"
EXPECTED_FILE="$SCRIPT_DIR/expected_issues.json"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_success() { echo -e "${GREEN}✓ $1${NC}"; }
print_error() { echo -e "${RED}✗ $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠ $1${NC}"; }
print_info() { echo -e "${BLUE}ℹ $1${NC}"; }

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Code Review 结果验证工具${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 检查结果文件是否存在
if [ ! -f "$RESULT_FILE" ]; then
    print_error "找不到审查结果文件: $RESULT_FILE"
    print_info "请先运行 Code Review 生成结果文件"
    exit 1
fi

if [ ! -f "$EXPECTED_FILE" ]; then
    print_warning "找不到预期结果配置文件: $EXPECTED_FILE"
    print_info "将使用基本验证模式"
fi

print_success "找到审查结果文件: $RESULT_FILE"
echo ""

# 统计实际检测到的问题
echo "步骤 1: 统计检测结果..."

P0_COUNT=$(grep -c "\[P0\]" "$RESULT_FILE" 2>/dev/null || echo "0")
P1_COUNT=$(grep -c "\[P1\]" "$RESULT_FILE" 2>/dev/null || echo "0")
P2_COUNT=$(grep -c "\[P2\]" "$RESULT_FILE" 2>/dev/null || echo "0")
TOTAL_ISSUES=$((P0_COUNT + P1_COUNT + P2_COUNT))

echo "实际检测结果:"
echo "  P0 (必须修复): $P0_COUNT"
echo "  P1 (强烈建议): $P1_COUNT"
echo "  P2 (建议优化): $P2_COUNT"
echo "  总计: $TOTAL_ISSUES"
echo ""

# 提取所有检测到的规则编号
echo "步骤 2: 提取规则编号..."
DETECTED_RULES=$(grep -oP "规则 \K[0-9]+\.[0-9]+\.[0-9]+" "$RESULT_FILE" 2>/dev/null | sort -u || echo "")
DETECTED_COUNT=$(echo "$DETECTED_RULES" | grep -c "." || echo "0")

if [ $DETECTED_COUNT -gt 0 ]; then
    print_success "检测到 $DETECTED_COUNT 个不同的规则违规"
    echo "检测到的规则:"
    echo "$DETECTED_RULES" | while read rule; do
        echo "  - Rule $rule"
    done
else
    print_warning "未能提取规则编号"
fi
echo ""

# 分析文件覆盖
echo "步骤 3: 分析文件覆盖..."
FILES=$(grep -oP "## 文件: \K.*" "$RESULT_FILE" 2>/dev/null || echo "")
FILE_COUNT=$(echo "$FILES" | grep -c "." || echo "0")

if [ $FILE_COUNT -gt 0 ]; then
    print_success "审查了 $FILE_COUNT 个文件"
    echo "审查的文件:"
    echo "$FILES" | while read file; do
        echo "  - $file"
    done
else
    print_warning "未能识别审查的文件"
fi
echo ""

# 基本合理性检查
echo "步骤 4: 合理性检查..."

validation_passed=true

# 检查 1: user_service_bad.go 应该有 P0 问题
if echo "$FILES" | grep -q "user_service_bad.go"; then
    if [ $P0_COUNT -ge 10 ]; then
        print_success "user_service_bad.go: 检测到足够的 P0 问题 ($P0_COUNT >= 10)"
    else
        print_error "user_service_bad.go: P0 问题数量不足 ($P0_COUNT < 10)"
        print_warning "可能存在漏检"
        validation_passed=false
    fi
else
    print_warning "未审查 user_service_bad.go"
fi

# 检查 2: 关键规则是否被检测
echo ""
echo "关键规则检测验证:"

critical_rules=("1.1.1" "1.3.1" "1.3.4" "1.2.1" "1.4.1")
for rule in "${critical_rules[@]}"; do
    if echo "$DETECTED_RULES" | grep -q "^$rule$"; then
        print_success "Rule $rule: 已检测"
    else
        print_error "Rule $rule: 未检测 (可能漏检)"
        validation_passed=false
    fi
done

echo ""

# 检查 3: user_service_good.go 不应该有问题
if echo "$FILES" | grep -q "user_service_good.go"; then
    # 计算 good 文件的问题数
    # (这需要解析每个文件的问题数，简化处理)
    print_info "user_service_good.go: 已审查 (手动验证是否误报)"
fi

# 生成详细报告
echo ""
echo "步骤 5: 生成验证报告..."

REPORT_FILE="$SCRIPT_DIR/validation_report.md"

cat > "$REPORT_FILE" << EOF
# Code Review 验证报告

**生成时间**: $(date '+%Y-%m-%d %H:%M:%S')
**结果文件**: $RESULT_FILE

## 检测统计

| 优先级 | 数量 |
|--------|------|
| P0 (必须修复) | $P0_COUNT |
| P1 (强烈建议) | $P1_COUNT |
| P2 (建议优化) | $P2_COUNT |
| **总计** | **$TOTAL_ISSUES** |

## 规则覆盖

检测到 $DETECTED_COUNT 个不同的规则违规:

EOF

echo "$DETECTED_RULES" | while read rule; do
    echo "- Rule $rule" >> "$REPORT_FILE"
done

cat >> "$REPORT_FILE" << EOF

## 文件覆盖

审查了 $FILE_COUNT 个文件:

EOF

echo "$FILES" | while read file; do
    echo "- $file" >> "$REPORT_FILE"
done

cat >> "$REPORT_FILE" << EOF

## 验证结果

EOF

if [ "$validation_passed" = true ]; then
    echo "**状态**: ✅ 通过" >> "$REPORT_FILE"
    echo "" >> "$REPORT_FILE"
    echo "所有关键规则都被正确检测。" >> "$REPORT_FILE"
else
    echo "**状态**: ❌ 失败" >> "$REPORT_FILE"
    echo "" >> "$REPORT_FILE"
    echo "发现潜在的漏检问题，请检查技能实现。" >> "$REPORT_FILE"
fi

cat >> "$REPORT_FILE" << EOF

## 建议

1. 检查是否所有预期问题都被检测到
2. 验证good文件是否有误报
3. 对比 expected_issues.json 中的规则列表
4. 更新技能文档以覆盖缺失的规则

---

*此报告由 validate_results.sh 自动生成*
EOF

print_success "验证报告已保存: $REPORT_FILE"
echo ""

# 最终结果
echo "=========================================="
if [ "$validation_passed" = true ]; then
    print_success "验证通过！"
    echo ""
    echo "检测能力: 正常"
    echo "漏检风险: 低"
    exit 0
else
    print_error "验证失败！"
    echo ""
    echo "检测能力: 需要改进"
    echo "漏检风险: 中-高"
    echo ""
    print_info "请查看验证报告了解详情: $REPORT_FILE"
    exit 1
fi
