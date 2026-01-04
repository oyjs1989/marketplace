# CI/CD 集成指南

## 概述

本文档说明如何将 Go Code Review Skill 测试集成到 CI/CD 流水线中。

## GitLab CI 集成

### .gitlab-ci.yml 示例

```yaml
# .gitlab-ci.yml
stages:
  - test
  - validate

# 测试 Go Code Review Skill
code_review_test:
  stage: test
  image: alpine:latest
  before_script:
    - apk add --no-cache bash grep
  script:
    # 注意: 此处需要Claude Code环境，实际CI中需要特殊配置
    - echo "运行 Code Review 测试..."
    - cd test-cases/go-code-review

    # 方法1: 如果可以在CI中使用Claude Code
    # - claude review bad/*.go

    # 方法2: 使用预先生成的结果验证
    - |
      if [ -f "code_review.result" ]; then
        bash validate_results.sh code_review.result
      else
        echo "警告: 未找到code_review.result，跳过验证"
        exit 0
      fi
  artifacts:
    paths:
      - test-cases/go-code-review/validation_report.md
      - test-cases/go-code-review/code_review.result
    expire_in: 30 days
    when: always
  allow_failure: true  # 允许测试失败但不阻塞流水线

# 技能文档验证
skill_docs_check:
  stage: test
  image: alpine:latest
  before_script:
    - apk add --no-cache bash grep
  script:
    - echo "验证技能文档完整性..."
    # 检查所有SKILL.md是否包含中文输出要求
    - |
      cd skills/go-code-review
      for skill in orchestrator gorm-review error-safety naming-logging organization; do
        if ! grep -q "必须使用中文" $skill/SKILL.md; then
          echo "错误: $skill/SKILL.md 缺少中文输出要求"
          exit 1
        fi
      done
    - echo "✓ 所有技能文档格式正确"
  rules:
    - changes:
        - skills/**/*.md
```

## GitHub Actions 集成

### .github/workflows/test-skill.yml

```yaml
name: Test Go Code Review Skill

on:
  push:
    branches: [ main, develop ]
    paths:
      - 'skills/go-code-review/**'
      - 'test-cases/go-code-review/**'
  pull_request:
    branches: [ main ]
    paths:
      - 'skills/go-code-review/**'
      - 'test-cases/go-code-review/**'

jobs:
  test:
    runs-on: ubuntu-latest

    steps:
    - name: Checkout代码
      uses: actions/checkout@v3

    - name: 设置环境
      run: |
        chmod +x test-cases/go-code-review/*.sh

    - name: 验证技能文档
      run: |
        echo "检查技能文档格式..."
        cd skills/go-code-review
        for skill in orchestrator gorm-review error-safety naming-logging organization; do
          echo "检查 $skill/SKILL.md"
          if ! grep -q "必须使用中文" $skill/SKILL.md; then
            echo "::error::$skill/SKILL.md 缺少中文输出要求"
            exit 1
          fi
        done
        echo "✓ 技能文档验证通过"

    - name: 检查测试用例存在性
      run: |
        echo "验证测试文件..."
        test -f test-cases/go-code-review/bad/user_service_bad.go || exit 1
        test -f test-cases/go-code-review/good/user_service_good.go || exit 1
        test -f test-cases/go-code-review/expected_issues.json || exit 1
        echo "✓ 测试文件完整"

    - name: 验证预期结果配置
      run: |
        echo "检查 expected_issues.json 格式..."
        # 需要 jq 工具
        sudo apt-get update && sudo apt-get install -y jq
        jq empty test-cases/go-code-review/expected_issues.json
        echo "✓ JSON 格式正确"

    - name: 运行测试脚本检查
      run: |
        echo "验证测试脚本语法..."
        bash -n test-cases/go-code-review/run_test.sh
        bash -n test-cases/go-code-review/validate_results.sh
        echo "✓ 脚本语法正确"

    - name: 上传测试报告
      if: always()
      uses: actions/upload-artifact@v3
      with:
        name: test-reports
        path: |
          test-cases/go-code-review/validation_report.md
          test-cases/go-code-review/test_result.md
        retention-days: 30

  # 可选: 如果能够在CI中运行Claude Code
  # integration-test:
  #   runs-on: ubuntu-latest
  #   needs: test
  #   steps:
  #   - name: 运行集成测试
  #     run: |
  #       # 这需要特殊的Claude Code CI环境
  #       claude review test-cases/go-code-review/bad/*.go
  #       bash test-cases/go-code-review/validate_results.sh
```

## Jenkins Pipeline

### Jenkinsfile

```groovy
pipeline {
    agent any

    stages {
        stage('准备') {
            steps {
                checkout scm
                sh 'chmod +x test-cases/go-code-review/*.sh'
            }
        }

        stage('验证文档') {
            steps {
                script {
                    sh '''
                        cd skills/go-code-review
                        for skill in orchestrator gorm-review error-safety naming-logging organization; do
                            echo "检查 $skill/SKILL.md"
                            grep -q "必须使用中文" $skill/SKILL.md || exit 1
                        done
                    '''
                }
            }
        }

        stage('检查测试文件') {
            steps {
                sh '''
                    test -f test-cases/go-code-review/bad/user_service_bad.go
                    test -f test-cases/go-code-review/expected_issues.json
                '''
            }
        }

        stage('验证测试结果') {
            when {
                expression {
                    fileExists('test-cases/go-code-review/code_review.result')
                }
            }
            steps {
                sh '''
                    cd test-cases/go-code-review
                    bash validate_results.sh code_review.result
                '''
            }
        }
    }

    post {
        always {
            archiveArtifacts artifacts: 'test-cases/go-code-review/*.md', allowEmptyArchive: true
        }
        failure {
            emailext (
                subject: "测试失败: ${env.JOB_NAME} - ${env.BUILD_NUMBER}",
                body: "Go Code Review Skill 测试失败，请检查构建日志",
                to: "team@example.com"
            )
        }
    }
}
```

## 本地 Pre-commit Hook

### .git/hooks/pre-commit

```bash
#!/bin/bash

# Pre-commit hook: 在提交前验证技能文档

echo "运行 pre-commit 检查..."

# 检查是否修改了技能文件
if git diff --cached --name-only | grep -q "skills/go-code-review"; then
    echo "检测到技能文件修改，验证文档格式..."

    # 验证中文输出要求
    for skill in orchestrator gorm-review error-safety naming-logging organization; do
        file="skills/go-code-review/$skill/SKILL.md"
        if [ -f "$file" ]; then
            if ! grep -q "必须使用中文" "$file"; then
                echo "错误: $file 缺少中文输出要求"
                echo "请在文档中添加输出格式说明"
                exit 1
            fi
        fi
    done

    echo "✓ 技能文档验证通过"
fi

# 提醒运行测试
if git diff --cached --name-only | grep -q "skills/go-code-review/.*SKILL.md"; then
    echo ""
    echo "⚠️  提醒: 你修改了技能文档"
    echo "建议运行测试: Review test-cases/go-code-review/bad/*.go"
    echo ""
fi

exit 0
```

### 安装 Pre-commit Hook

```bash
# 在仓库根目录执行
cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash
# (上面的pre-commit脚本内容)
EOF

chmod +x .git/hooks/pre-commit
```

## 定期回归测试

### Cron 任务 (每周测试)

```bash
# 添加到 crontab
# 每周一上午 9:00 运行回归测试

0 9 * * 1 cd /path/to/marketplace && bash test-cases/go-code-review/run_test.sh && mail -s "Go Code Review 周报" team@example.com < test-cases/go-code-review/validation_report.md
```

### GitLab Scheduled Pipeline

```yaml
# .gitlab-ci.yml
# 添加定期任务

weekly_regression_test:
  stage: test
  script:
    - bash test-cases/go-code-review/run_test.sh
  only:
    - schedules
  variables:
    SCHEDULE_TYPE: "weekly_regression"
```

在 GitLab UI 中设置:
- CI/CD → Schedules → New schedule
- 间隔: 每周一次
- 目标分支: main
- 变量: SCHEDULE_TYPE=weekly_regression

## 测试结果通知

### Slack 通知

```bash
#!/bin/bash
# notify_slack.sh

WEBHOOK_URL="https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
RESULT_FILE="validation_report.md"

if [ -f "$RESULT_FILE" ]; then
    # 提取测试结果
    if grep -q "验证通过" "$RESULT_FILE"; then
        STATUS="success"
        MESSAGE="✅ Go Code Review 测试通过"
        COLOR="good"
    else
        STATUS="failure"
        MESSAGE="❌ Go Code Review 测试失败"
        COLOR="danger"
    fi

    # 发送到Slack
    curl -X POST -H 'Content-type: application/json' \
        --data "{
            \"attachments\": [{
                \"color\": \"$COLOR\",
                \"title\": \"$MESSAGE\",
                \"text\": \"查看完整报告: $CI_JOB_URL\"
            }]
        }" \
        $WEBHOOK_URL
fi
```

### 企业微信通知

```bash
#!/bin/bash
# notify_wechat.sh

WEBHOOK_URL="https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=YOUR-KEY"
RESULT_FILE="validation_report.md"

if grep -q "验证通过" "$RESULT_FILE"; then
    CONTENT="✅ Go Code Review Skill 测试通过"
else
    CONTENT="❌ Go Code Review Skill 测试失败，请检查"
fi

curl -X POST -H 'Content-Type: application/json' \
    --data "{
        \"msgtype\": \"text\",
        \"text\": {
            \"content\": \"$CONTENT\"
        }
    }" \
    $WEBHOOK_URL
```

## Docker 集成

### Dockerfile.test

```dockerfile
FROM alpine:latest

# 安装依赖
RUN apk add --no-cache bash grep coreutils

# 复制测试文件
COPY test-cases/go-code-review /workspace/test-cases/go-code-review
WORKDIR /workspace

# 运行测试
CMD ["bash", "test-cases/go-code-review/run_test.sh"]
```

### 使用Docker运行测试

```bash
# 构建测试镜像
docker build -f Dockerfile.test -t go-review-test .

# 运行测试
docker run --rm -v $(pwd):/workspace go-review-test bash test-cases/go-code-review/validate_results.sh
```

## 测试覆盖率追踪

### 创建覆盖率追踪脚本

```bash
#!/bin/bash
# track_coverage.sh

EXPECTED_FILE="test-cases/go-code-review/expected_issues.json"
RESULT_FILE="code_review.result"

# 提取预期规则总数
TOTAL_RULES=$(jq '.rules_coverage.total_rules' $EXPECTED_FILE)
COVERED_RULES=$(jq '.rules_coverage.covered_rules' $EXPECTED_FILE)

# 提取实际检测到的规则
DETECTED_RULES=$(grep -oP "规则 \K[0-9]+\.[0-9]+\.[0-9]+" $RESULT_FILE | sort -u | wc -l)

# 计算覆盖率
COVERAGE_PERCENT=$((DETECTED_RULES * 100 / TOTAL_RULES))

echo "规则覆盖率追踪"
echo "==============="
echo "总规则数: $TOTAL_RULES"
echo "已覆盖: $COVERED_RULES (配置)"
echo "本次检测: $DETECTED_RULES"
echo "覆盖率: $COVERAGE_PERCENT%"

# 保存到历史文件
echo "$(date +%Y-%m-%d),$DETECTED_RULES,$COVERAGE_PERCENT%" >> test-cases/go-code-review/coverage_history.csv
```

## 最佳实践

### 1. 分阶段集成

```
阶段1: 本地开发
  └─ Pre-commit hook (快速检查)

阶段2: Pull Request
  └─ CI 文档验证 (必须通过)

阶段3: Merge to Main
  └─ 完整测试 (允许失败，但记录)

阶段4: 定期回归
  └─ 周报告 (监控趋势)
```

### 2. 渐进式要求

- **初期**: 文档格式检查（必须通过）
- **中期**: 测试用例存在性（必须通过）
- **后期**: 完整测试验证（强制通过）

### 3. 监控指标

追踪这些指标：
- 规则覆盖率 (目标: >80%)
- P0 检测率 (目标: 100%)
- 误报率 (目标: <5%)
- 测试执行时间 (目标: <60秒)

### 4. 失败处理

```yaml
# CI配置中的失败策略
allow_failure: true  # 初期允许失败
manual_verification: true  # 需要人工确认

# 逐步过渡到
allow_failure: false  # 测试失败阻塞流水线
```

## 故障排查

### CI中Claude Code不可用

**解决方案**:
- 使用本地生成的 code_review.result 文件
- 在开发环境运行测试，将结果提交到仓库
- CI只验证已有结果的格式和完整性

### 测试超时

**解决方案**:
```yaml
# 增加超时时间
timeout: 5m

# 或分批测试
script:
  - timeout 2m claude review bad/user_service_bad.go
  - timeout 2m claude review bad/project_structure_bad.go
```

### 环境差异导致结果不一致

**解决方案**:
- 使用Docker容器统一环境
- 固定Claude Code版本
- 记录环境信息到测试报告

---

**维护者**: Jason Ouyang
**更新日期**: 2026-01-04
**版本**: 1.0.0
