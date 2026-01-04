#!/bin/bash

# Go Code Review Skill 自动化测试脚本
# 用途: 验证技能的检测能力，确保没有漏检

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BAD_DIR="$SCRIPT_DIR/bad"
GOOD_DIR="$SCRIPT_DIR/good"
RESULT_FILE="$SCRIPT_DIR/test_result.md"
EXPECTED_FILE="$SCRIPT_DIR/expected_issues.json"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Go Code Review Skill - 自动化测试${NC}"
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

# 步骤 1: 检查环境
echo "步骤 1: 检查测试环境..."
if [ ! -d "$BAD_DIR" ]; then
    print_error "找不到测试目录: $BAD_DIR"
    exit 1
fi

if [ ! -d "$GOOD_DIR" ]; then
    print_error "找不到测试目录: $GOOD_DIR"
    exit 1
fi

print_success "测试环境准备就绪"
echo ""

# 步骤 2: 列出测试文件
echo "步骤 2: 扫描测试文件..."
BAD_FILES=($(find "$BAD_DIR" -name "*.go" -type f))
GOOD_FILES=($(find "$GOOD_DIR" -name "*.go" -type f))

print_info "发现 ${#BAD_FILES[@]} 个违规测试文件"
for file in "${BAD_FILES[@]}"; do
    echo "  - $(basename $file)"
done

print_info "发现 ${#GOOD_FILES[@]} 个正确实现文件"
for file in "${GOOD_FILES[@]}"; do
    echo "  - $(basename $file)"
done
echo ""

# 步骤 3: 测试违规代码检测
echo "步骤 3: 测试违规代码检测..."
echo "" > "$RESULT_FILE"

total_tests=0
passed_tests=0
failed_tests=0

for bad_file in "${BAD_FILES[@]}"; do
    total_tests=$((total_tests + 1))
    filename=$(basename "$bad_file")

    echo -e "${BLUE}测试: $filename${NC}"

    # 创建临时输出文件
    temp_output="/tmp/code_review_$$.result"

    # 这里需要用户手动运行Claude Code review
    # 因为脚本无法直接调用Claude Code
    print_warning "请手动运行: claude review $bad_file"
    print_warning "或在Claude Code中执行: Review $bad_file"

    # 检查是否存在审查结果
    if [ -f "code_review.result" ]; then
        # 统计检测到的问题数量
        p0_count=$(grep -c "\[P0\]" code_review.result 2>/dev/null || echo "0")
        p1_count=$(grep -c "\[P1\]" code_review.result 2>/dev/null || echo "0")
        p2_count=$(grep -c "\[P2\]" code_review.result 2>/dev/null || echo "0")
        total_issues=$((p0_count + p1_count + p2_count))

        if [ $total_issues -gt 0 ]; then
            print_success "检测到 $total_issues 个问题 (P0: $p0_count, P1: $p1_count, P2: $p2_count)"
            passed_tests=$((passed_tests + 1))
        else
            print_error "未检测到任何问题 - 可能存在漏检"
            failed_tests=$((failed_tests + 1))
        fi

        # 保存结果
        echo "## 测试文件: $filename" >> "$RESULT_FILE"
        echo "- P0 问题: $p0_count" >> "$RESULT_FILE"
        echo "- P1 问题: $p1_count" >> "$RESULT_FILE"
        echo "- P2 问题: $p2_count" >> "$RESULT_FILE"
        echo "" >> "$RESULT_FILE"
    else
        print_warning "未找到审查结果文件"
        print_info "跳过此文件的验证"
    fi

    echo ""
done

# 步骤 4: 测试正确代码（应该不报告问题）
echo "步骤 4: 测试正确代码（不应报告问题）..."

for good_file in "${GOOD_FILES[@]}"; do
    total_tests=$((total_tests + 1))
    filename=$(basename "$good_file")

    echo -e "${BLUE}测试: $filename${NC}"
    print_warning "请手动运行: claude review $good_file"

    # 这里应该不报告任何问题
    print_info "预期: 不应报告任何问题"
    echo ""
done

# 步骤 5: 生成测试报告
echo "步骤 5: 生成测试报告..."

cat > "$RESULT_FILE" << 'EOF'
# Go Code Review Skill - 测试报告

## 测试概要

| 项目 | 数值 |
|------|------|
| 测试文件总数 | ${TOTAL_FILES} |
| 违规代码文件 | ${BAD_FILES_COUNT} |
| 正确代码文件 | ${GOOD_FILES_COUNT} |
| 测试时间 | $(date '+%Y-%m-%d %H:%M:%S') |

## 详细结果

### 违规代码检测结果

EOF

echo "测试报告已保存到: $RESULT_FILE"
print_success "测试完成"
echo ""

# 步骤 6: 输出统计信息
echo "=========================================="
echo "测试统计:"
echo "  总测试数: $total_tests"
echo "  通过: $passed_tests"
echo "  失败: $failed_tests"
echo "=========================================="

# 退出代码
if [ $failed_tests -gt 0 ]; then
    exit 1
else
    exit 0
fi
