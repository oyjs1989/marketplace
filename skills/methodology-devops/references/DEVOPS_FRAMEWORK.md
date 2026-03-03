# DevOps框架详解

本文档提供DevOps方法论的详细实践指南、模板和检查清单。

> **来源**: 基于 doc.md 第3.3节 - DevOps的反馈循环机制

---

## DevOps文化基础

### 打破壁垒

**传统问题**:
```markdown
开发团队 (Dev):
- 目标: 快速交付新功能
- KPI: 功能上线数量

运维团队 (Ops):
- 目标: 系统稳定运行
- KPI: 可用性和故障率

结果: 利益冲突,互相指责
```

**DevOps解决方案**:
```markdown
融合团队:
- 共同目标: 快速且稳定地交付价值
- 共享责任: 全生命周期负责
- 共同On-call: 谁开发谁运维

文化转变:
- "不是我的问题" → "我们的问题"
- "抛过墙" → "端到端负责"
- "责备文化" → "无责归咎文化"
```

### 无责归咎文化 (Blameless Culture)

**原则**:
```markdown
事后分析 (Post-Mortem):
- ✅ 聚焦系统问题,非个人问题
- ✅ 学习为目的,非惩罚
- ✅ 透明分享,全员学习
- ❌ 不追究个人责任
- ❌ 不隐藏错误
```

**Post-Mortem模板**:
```markdown
# 事件Post-Mortem: [事件标题]

## 事件摘要
- **发生时间**: [开始时间] - [结束时间]
- **持续时间**: [X小时Y分钟]
- **影响范围**: [受影响用户数/功能]
- **事件级别**: P0/P1/P2

## 时间线
| 时间 | 事件 | 操作人 |
|------|------|--------|
| 14:32 | 监控告警触发 | 自动 |
| 14:35 | On-call工程师响应 | 张三 |
| 14:40 | 确认服务不可用 | 张三 |
| 14:45 | 启动War Room | 李四 |
| 15:10 | 识别根因: 数据库连接池耗尽 | 团队 |
| 15:20 | 执行修复: 增加连接池大小 | 王五 |
| 15:30 | 服务恢复正常 | - |
| 16:00 | 验证完成,告警解除 | - |

## 根因分析 (5 Whys)
1. **为什么服务不可用?**
   - 数据库连接池耗尽

2. **为什么连接池耗尽?**
   - 新功能导致数据库查询量突增

3. **为什么查询量突增?**
   - 新功能未进行充分性能测试

4. **为什么未进行性能测试?**
   - 测试流程未包含性能测试

5. **为什么流程未包含?**
   - 缺乏性能测试规范和自动化

## 改进行动
| 行动 | 负责人 | 期限 | 状态 |
|------|--------|------|------|
| 制定性能测试规范 | 张三 | 2周 | 进行中 |
| 引入性能测试工具 (k6) | 李四 | 3周 | 待开始 |
| 添加数据库连接池监控 | 王五 | 1周 | 已完成 |
| 设置连接池使用率告警 | 赵六 | 1周 | 已完成 |

## 经验教训
### 做得好的
- ✅ 监控及时发现问题
- ✅ 响应迅速 (3分钟内响应)
- ✅ 团队协作高效

### 需要改进的
- ❌ 性能测试流程缺失
- ❌ 数据库监控不全面
- ❌ 新功能上线前未进行容量评估

## 附录
- 相关监控截图
- 日志片段
- 配置变更记录
```

---

## CI/CD管道详解

### 阶段1: 源代码管理

**Git分支策略**:

#### GitFlow
```markdown
分支结构:
- main: 生产环境代码
- develop: 开发主干
- feature/*: 功能分支
- release/*: 发布分支
- hotfix/*: 紧急修复

工作流:
feature → develop → release → main
                        ↓
                    hotfix → main + develop

优点: 结构清晰,适合计划发布
缺点: 分支多,合并复杂
```

#### Trunk-Based Development
```markdown
分支结构:
- main: 单一主干
- short-lived feature branches (< 2天)

工作流:
feature (1-2天) → main

优点: 简单,减少合并冲突,支持持续交付
缺点: 需要Feature Toggle

最佳实践:
- 每天至少合并一次
- 使用Feature Flag控制功能可见性
- 强依赖自动化测试
```

**Pull Request流程**:
```markdown
1. 创建PR
   - [ ] 描述变更内容
   - [ ] 关联Issue/Ticket
   - [ ] 自测通过

2. 自动检查
   - [ ] CI构建通过
   - [ ] 代码检查通过 (Lint)
   - [ ] 测试覆盖率达标
   - [ ] 安全扫描通过

3. 人工审查
   - [ ] 至少1人审查 (小改动)
   - [ ] 至少2人审查 (大改动)
   - [ ] 架构师审查 (架构变更)

4. 合并
   - [ ] Squash Merge (保持历史清晰)
   - [ ] 删除feature分支
   - [ ] 自动部署到测试环境
```

### 阶段2: 持续集成

**Jenkins Pipeline示例**:
```groovy
pipeline {
    agent any

    environment {
        DOCKER_REGISTRY = 'registry.example.com'
        APP_NAME = 'myapp'
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Build') {
            steps {
                sh './gradlew build'
            }
        }

        stage('Unit Test') {
            steps {
                sh './gradlew test'
                junit 'build/test-results/**/*.xml'
            }
        }

        stage('Code Quality') {
            steps {
                sh './gradlew sonarqube'
            }
        }

        stage('Build Docker Image') {
            steps {
                script {
                    docker.build("${DOCKER_REGISTRY}/${APP_NAME}:${BUILD_NUMBER}")
                }
            }
        }

        stage('Push to Registry') {
            steps {
                script {
                    docker.withRegistry("https://${DOCKER_REGISTRY}", 'docker-credentials') {
                        docker.image("${DOCKER_REGISTRY}/${APP_NAME}:${BUILD_NUMBER}").push()
                        docker.image("${DOCKER_REGISTRY}/${APP_NAME}:${BUILD_NUMBER}").push('latest')
                    }
                }
            }
        }

        stage('Deploy to Test') {
            steps {
                sh "kubectl set image deployment/${APP_NAME} ${APP_NAME}=${DOCKER_REGISTRY}/${APP_NAME}:${BUILD_NUMBER} -n test"
            }
        }

        stage('Integration Test') {
            steps {
                sh './scripts/run-integration-tests.sh test'
            }
        }
    }

    post {
        always {
            cleanWs()
        }
        success {
            slackSend color: 'good', message: "Build #${BUILD_NUMBER} succeeded"
        }
        failure {
            slackSend color: 'danger', message: "Build #${BUILD_NUMBER} failed"
        }
    }
}
```

**GitLab CI配置示例**:
```yaml
# .gitlab-ci.yml
stages:
  - build
  - test
  - package
  - deploy

variables:
  DOCKER_DRIVER: overlay2
  DOCKER_TLS_CERTDIR: ""

build:
  stage: build
  image: golang:1.21
  script:
    - go build -o myapp
  artifacts:
    paths:
      - myapp
    expire_in: 1 day

unit-test:
  stage: test
  image: golang:1.21
  script:
    - go test -v -coverprofile=coverage.out ./...
    - go tool cover -html=coverage.out -o coverage.html
  coverage: '/coverage: \d+.\d+% of statements/'
  artifacts:
    paths:
      - coverage.html
    reports:
      coverage_report:
        coverage_format: cobertura
        path: coverage.xml

lint:
  stage: test
  image: golangci/golangci-lint:latest
  script:
    - golangci-lint run

docker-build:
  stage: package
  image: docker:latest
  services:
    - docker:dind
  before_script:
    - docker login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY
  script:
    - docker build -t $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA .
    - docker push $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA
    - docker tag $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA $CI_REGISTRY_IMAGE:latest
    - docker push $CI_REGISTRY_IMAGE:latest

deploy-test:
  stage: deploy
  image: bitnami/kubectl:latest
  script:
    - kubectl config use-context test-cluster
    - kubectl set image deployment/myapp myapp=$CI_REGISTRY_IMAGE:$CI_COMMIT_SHA -n test
    - kubectl rollout status deployment/myapp -n test
  environment:
    name: test
    url: https://test.example.com
  only:
    - develop

deploy-prod:
  stage: deploy
  image: bitnami/kubectl:latest
  script:
    - kubectl config use-context prod-cluster
    - kubectl set image deployment/myapp myapp=$CI_REGISTRY_IMAGE:$CI_COMMIT_SHA -n production
    - kubectl rollout status deployment/myapp -n production
  environment:
    name: production
    url: https://example.com
  when: manual
  only:
    - main
```

### 阶段3: 测试自动化

**测试金字塔**:
```
          /\
         /  \  E2E测试 (10%)
        /____\
       /      \
      / 集成测试 \ (20%)
     /___________\
    /             \
   /  单元测试 (70%) \
  /__________________\
```

**单元测试 (Unit Test)**:
```go
// user_service_test.go
package service

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
    // Arrange
    service := NewUserService()
    user := &User{
        Name:  "Test User",
        Email: "test@example.com",
    }

    // Act
    result, err := service.CreateUser(user)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, user.Name, result.Name)
    assert.NotEmpty(t, result.ID)
}
```

**集成测试 (Integration Test)**:
```go
// user_repository_test.go
package repository

import (
    "testing"
    "github.com/stretchr/testify/suite"
)

type UserRepositoryTestSuite struct {
    suite.Suite
    db   *sql.DB
    repo *UserRepository
}

func (suite *UserRepositoryTestSuite) SetupTest() {
    // 初始化测试数据库
    suite.db = setupTestDB()
    suite.repo = NewUserRepository(suite.db)
}

func (suite *UserRepositoryTestSuite) TearDownTest() {
    // 清理测试数据
    suite.db.Close()
}

func (suite *UserRepositoryTestSuite) TestFindByID() {
    // 插入测试数据
    user := &User{Name: "Test", Email: "test@example.com"}
    id, _ := suite.repo.Create(user)

    // 查询
    found, err := suite.repo.FindByID(id)

    // 验证
    suite.NoError(err)
    suite.Equal(user.Name, found.Name)
}

func TestUserRepositoryTestSuite(t *testing.T) {
    suite.Run(t, new(UserRepositoryTestSuite))
}
```

**E2E测试 (Cypress示例)**:
```javascript
// cypress/e2e/user_registration.cy.js
describe('User Registration', () => {
  beforeEach(() => {
    cy.visit('/register')
  })

  it('should register a new user successfully', () => {
    // 填写表单
    cy.get('input[name="name"]').type('Test User')
    cy.get('input[name="email"]').type('test@example.com')
    cy.get('input[name="password"]').type('password123')
    cy.get('input[name="confirmPassword"]').type('password123')

    // 提交
    cy.get('button[type="submit"]').click()

    // 验证
    cy.url().should('include', '/dashboard')
    cy.contains('Welcome, Test User')
  })

  it('should show error for duplicate email', () => {
    cy.get('input[name="email"]').type('existing@example.com')
    cy.get('button[type="submit"]').click()

    cy.contains('Email already exists')
  })
})
```

**性能测试 (k6示例)**:
```javascript
// load_test.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  stages: [
    { duration: '2m', target: 100 }, // 爬坡到100用户
    { duration: '5m', target: 100 }, // 保持100用户
    { duration: '2m', target: 200 }, // 爬坡到200用户
    { duration: '5m', target: 200 }, // 保持200用户
    { duration: '2m', target: 0 },   // 降至0
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95%请求<500ms
    http_req_failed: ['rate<0.01'],   // 错误率<1%
  },
};

export default function () {
  let res = http.get('https://api.example.com/users');

  check(res, {
    'status is 200': (r) => r.status === 200,
    'response time < 500ms': (r) => r.timings.duration < 500,
  });

  sleep(1);
}
```

---

## 部署策略详解

### 蓝绿部署实施

**基础设施要求**:
```markdown
- 2套完全相同的环境 (蓝/绿)
- 负载均衡器 (可切换流量)
- 共享数据库 (或同步机制)
- 足够的资源容量 (2倍)
```

**Kubernetes实现**:
```yaml
# deployment-blue.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp-blue
  labels:
    app: myapp
    version: blue
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
      version: blue
  template:
    metadata:
      labels:
        app: myapp
        version: blue
    spec:
      containers:
      - name: myapp
        image: myapp:v1.0.0
        ports:
        - containerPort: 8080

---
# service.yaml
apiVersion: v1
kind: Service
metadata:
  name: myapp-service
spec:
  selector:
    app: myapp
    version: blue  # 切换到green即可切换流量
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: LoadBalancer
```

**切换脚本**:
```bash
#!/bin/bash
# blue-green-switch.sh

CURRENT=$(kubectl get service myapp-service -o jsonpath='{.spec.selector.version}')
NEW_VERSION=$1

if [ "$NEW_VERSION" != "blue" ] && [ "$NEW_VERSION" != "green" ]; then
    echo "Error: Version must be 'blue' or 'green'"
    exit 1
fi

echo "Current version: $CURRENT"
echo "Switching to: $NEW_VERSION"

# 1. 部署新版本
kubectl apply -f deployment-${NEW_VERSION}.yaml

# 2. 等待新版本就绪
kubectl rollout status deployment/myapp-${NEW_VERSION}

# 3. 健康检查
HEALTH_URL="http://myapp-${NEW_VERSION}.default.svc.cluster.local:8080/health"
for i in {1..10}; do
    if curl -f $HEALTH_URL; then
        echo "Health check passed"
        break
    fi
    if [ $i -eq 10 ]; then
        echo "Health check failed"
        exit 1
    fi
    sleep 5
done

# 4. 切换流量
kubectl patch service myapp-service -p "{\"spec\":{\"selector\":{\"version\":\"${NEW_VERSION}\"}}}"

echo "Traffic switched to $NEW_VERSION"

# 5. 保留旧版本24小时便于回滚
echo "Old version ($CURRENT) kept running for rollback"
```

### 金丝雀发布实施

**Istio配置示例**:
```yaml
# virtualservice.yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: myapp
spec:
  hosts:
  - myapp
  http:
  - match:
    - headers:
        canary:
          exact: "true"
    route:
    - destination:
        host: myapp
        subset: canary
  - route:
    - destination:
        host: myapp
        subset: stable
      weight: 90
    - destination:
        host: myapp
        subset: canary
      weight: 10  # 10%流量到金丝雀

---
# destinationrule.yaml
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: myapp
spec:
  host: myapp
  subsets:
  - name: stable
    labels:
      version: stable
  - name: canary
    labels:
      version: canary
```

**渐进式推进脚本**:
```bash
#!/bin/bash
# canary-rollout.sh

CANARY_WEIGHTS=(5 10 25 50 100)
MONITOR_DURATION=300  # 5分钟监控期

for WEIGHT in "${CANARY_WEIGHTS[@]}"; do
    echo "Setting canary weight to ${WEIGHT}%"

    # 更新流量权重
    kubectl patch virtualservice myapp --type merge -p "
    spec:
      http:
      - route:
        - destination:
            host: myapp
            subset: stable
          weight: $((100-WEIGHT))
        - destination:
            host: myapp
            subset: canary
          weight: ${WEIGHT}
    "

    # 监控关键指标
    echo "Monitoring for ${MONITOR_DURATION} seconds..."

    # 检查错误率
    ERROR_RATE=$(kubectl exec -it prometheus-0 -- promtool query instant \
        'rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m])' \
        | grep -oP '\d+\.\d+')

    if (( $(echo "$ERROR_RATE > 0.01" | bc -l) )); then
        echo "Error rate too high: $ERROR_RATE, rolling back"
        kubectl patch virtualservice myapp --type merge -p "
        spec:
          http:
          - route:
            - destination:
                host: myapp
                subset: stable
              weight: 100
        "
        exit 1
    fi

    # 检查响应时间
    P95_LATENCY=$(kubectl exec -it prometheus-0 -- promtool query instant \
        'histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))' \
        | grep -oP '\d+\.\d+')

    if (( $(echo "$P95_LATENCY > 0.5" | bc -l) )); then
        echo "P95 latency too high: $P95_LATENCY, rolling back"
        # 回滚...
        exit 1
    fi

    echo "Metrics healthy at ${WEIGHT}% canary"
    sleep $MONITOR_DURATION
done

echo "Canary rollout completed successfully"
```

---

## 监控和可观测性实施

### Prometheus配置

**prometheus.yml**:
```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

alerting:
  alertmanagers:
  - static_configs:
    - targets:
      - alertmanager:9093

rule_files:
  - "alerts/*.yml"

scrape_configs:
  # Kubernetes节点监控
  - job_name: 'kubernetes-nodes'
    kubernetes_sd_configs:
    - role: node
    relabel_configs:
    - action: labelmap
      regex: __meta_kubernetes_node_label_(.+)

  # Kubernetes Pod监控
  - job_name: 'kubernetes-pods'
    kubernetes_sd_configs:
    - role: pod
    relabel_configs:
    - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
      action: keep
      regex: true
    - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_path]
      action: replace
      target_label: __metrics_path__
      regex: (.+)
    - source_labels: [__address__, __meta_kubernetes_pod_annotation_prometheus_io_port]
      action: replace
      regex: ([^:]+)(?::\d+)?;(\d+)
      replacement: $1:$2
      target_label: __address__

  # 应用指标
  - job_name: 'myapp'
    static_configs:
    - targets: ['myapp:8080']
```

**告警规则**:
```yaml
# alerts/app_alerts.yml
groups:
- name: app_alerts
  interval: 30s
  rules:
  # 高错误率告警
  - alert: HighErrorRate
    expr: |
      rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m]) > 0.01
    for: 5m
    labels:
      severity: P1
    annotations:
      summary: "High error rate on {{ $labels.instance }}"
      description: "Error rate is {{ $value | humanizePercentage }}"

  # 高响应时间告警
  - alert: HighLatency
    expr: |
      histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 0.5
    for: 5m
    labels:
      severity: P1
    annotations:
      summary: "High latency on {{ $labels.instance }}"
      description: "P95 latency is {{ $value }}s"

  # 服务不可用告警
  - alert: ServiceDown
    expr: up{job="myapp"} == 0
    for: 1m
    labels:
      severity: P0
    annotations:
      summary: "Service {{ $labels.instance }} is down"
      description: "Service has been down for more than 1 minute"

  # CPU使用率高
  - alert: HighCPU
    expr: |
      100 - (avg by (instance) (irate(node_cpu_seconds_total{mode="idle"}[5m])) * 100) > 80
    for: 10m
    labels:
      severity: P2
    annotations:
      summary: "High CPU usage on {{ $labels.instance }}"
      description: "CPU usage is {{ $value }}%"

  # 内存使用率高
  - alert: HighMemory
    expr: |
      (node_memory_MemTotal_bytes - node_memory_MemAvailable_bytes) / node_memory_MemTotal_bytes * 100 > 90
    for: 5m
    labels:
      severity: P1
    annotations:
      summary: "High memory usage on {{ $labels.instance }}"
      description: "Memory usage is {{ $value }}%"
```

### Grafana仪表盘

**应用性能仪表盘 (JSON配置)**:
```json
{
  "dashboard": {
    "title": "Application Performance",
    "panels": [
      {
        "title": "QPS",
        "targets": [
          {
            "expr": "sum(rate(http_requests_total[5m]))"
          }
        ],
        "type": "graph"
      },
      {
        "title": "Response Time (P50, P95, P99)",
        "targets": [
          {
            "expr": "histogram_quantile(0.50, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "P50"
          },
          {
            "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "P95"
          },
          {
            "expr": "histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "P99"
          }
        ],
        "type": "graph"
      },
      {
        "title": "Error Rate",
        "targets": [
          {
            "expr": "sum(rate(http_requests_total{status=~\"5..\"}[5m])) / sum(rate(http_requests_total[5m]))"
          }
        ],
        "type": "singlestat",
        "format": "percentunit"
      },
      {
        "title": "Request by Status Code",
        "targets": [
          {
            "expr": "sum by (status) (rate(http_requests_total[5m]))"
          }
        ],
        "type": "piechart"
      }
    ]
  }
}
```

### 应用埋点

**Go应用Prometheus客户端**:
```go
package main

import (
    "net/http"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    httpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )

    httpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
        },
        []string{"method", "endpoint"},
    )
)

func instrumentHandler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // 包装ResponseWriter以捕获状态码
        wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

        next.ServeHTTP(wrapped, r)

        duration := time.Since(start).Seconds()
        status := wrapped.statusCode

        httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(status)).Inc()
        httpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

func main() {
    http.Handle("/metrics", promhttp.Handler())
    http.Handle("/api/", instrumentHandler(http.HandlerFunc(apiHandler)))

    http.ListenAndServe(":8080", nil)
}
```

---

## 基础设施即代码最佳实践

### Terraform模块化结构

```
terraform/
├── modules/                    # 可复用模块
│   ├── vpc/
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   └── outputs.tf
│   ├── compute/
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   └── outputs.tf
│   └── database/
│       ├── main.tf
│       ├── variables.tf
│       └── outputs.tf
├── environments/               # 环境配置
│   ├── dev/
│   │   ├── main.tf
│   │   ├── terraform.tfvars
│   │   └── backend.tf
│   ├── test/
│   │   ├── main.tf
│   │   ├── terraform.tfvars
│   │   └── backend.tf
│   └── prod/
│       ├── main.tf
│       ├── terraform.tfvars
│       └── backend.tf
└── shared/                     # 共享资源
    ├── main.tf
    └── terraform.tfvars
```

**VPC模块示例**:
```hcl
# modules/vpc/main.tf
resource "aws_vpc" "main" {
  cidr_block           = var.vpc_cidr
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = merge(
    var.tags,
    {
      Name = "${var.environment}-vpc"
    }
  )
}

resource "aws_subnet" "public" {
  count                   = length(var.public_subnet_cidrs)
  vpc_id                  = aws_vpc.main.id
  cidr_block              = var.public_subnet_cidrs[count.index]
  availability_zone       = var.availability_zones[count.index]
  map_public_ip_on_launch = true

  tags = merge(
    var.tags,
    {
      Name = "${var.environment}-public-${count.index + 1}"
      Type = "Public"
    }
  )
}

# modules/vpc/variables.tf
variable "environment" {
  description = "Environment name"
  type        = string
}

variable "vpc_cidr" {
  description = "VPC CIDR block"
  type        = string
}

variable "public_subnet_cidrs" {
  description = "Public subnet CIDR blocks"
  type        = list(string)
}

variable "availability_zones" {
  description = "Availability zones"
  type        = list(string)
}

variable "tags" {
  description = "Common tags"
  type        = map(string)
  default     = {}
}

# modules/vpc/outputs.tf
output "vpc_id" {
  description = "VPC ID"
  value       = aws_vpc.main.id
}

output "public_subnet_ids" {
  description = "Public subnet IDs"
  value       = aws_subnet.public[*].id
}
```

**生产环境配置**:
```hcl
# environments/prod/main.tf
terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "s3" {
    bucket         = "mycompany-terraform-state"
    key            = "prod/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "terraform-locks"
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Environment = "production"
      ManagedBy   = "Terraform"
      Project     = "MyApp"
    }
  }
}

module "vpc" {
  source = "../../modules/vpc"

  environment          = "prod"
  vpc_cidr            = "10.0.0.0/16"
  public_subnet_cidrs = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  availability_zones  = ["us-east-1a", "us-east-1b", "us-east-1c"]

  tags = {
    CostCenter = "Engineering"
  }
}

module "compute" {
  source = "../../modules/compute"

  environment    = "prod"
  vpc_id         = module.vpc.vpc_id
  subnet_ids     = module.vpc.public_subnet_ids
  instance_type  = "t3.large"
  min_size       = 3
  max_size       = 10
  desired_size   = 5
}

# environments/prod/terraform.tfvars
aws_region = "us-east-1"
```

**最佳实践**:
```markdown
- [ ] 使用远程状态存储 (S3 + DynamoDB锁)
- [ ] 状态文件加密
- [ ] 模块化和可复用
- [ ] 使用变量和输出
- [ ] 环境隔离
- [ ] 使用workspaces或独立目录
- [ ] 敏感信息使用secrets管理 (Vault, AWS Secrets Manager)
- [ ] 代码审查IaC变更
- [ ] 使用terraform plan预览变更
- [ ] 使用terraform fmt格式化代码
- [ ] 使用terraform validate验证配置
```

---

## DORA指标追踪

### 1. 部署频率 (Deployment Frequency)

**度量方法**:
```sql
-- 统计每天部署次数
SELECT
    DATE(deployment_time) as date,
    COUNT(*) as deployment_count
FROM deployments
WHERE environment = 'production'
  AND deployment_time >= NOW() - INTERVAL '30 days'
GROUP BY DATE(deployment_time)
ORDER BY date DESC;
```

**分级标准**:
```markdown
Elite: 多次/天
High: 1次/天 - 1次/周
Medium: 1次/周 - 1次/月
Low: < 1次/月
```

### 2. 变更前置时间 (Lead Time for Changes)

**度量方法**:
```sql
-- 从代码提交到生产部署的时间
SELECT
    AVG(EXTRACT(EPOCH FROM (d.deployment_time - c.commit_time)) / 3600) as avg_lead_time_hours
FROM deployments d
JOIN commits c ON d.commit_sha = c.sha
WHERE d.environment = 'production'
  AND d.deployment_time >= NOW() - INTERVAL '30 days';
```

**分级标准**:
```markdown
Elite: < 1小时
High: 1天 - 1周
Medium: 1周 - 1月
Low: 1月 - 6月
```

### 3. 平均恢复时间 (Mean Time to Recovery)

**度量方法**:
```sql
-- 从事件发生到恢复的平均时间
SELECT
    AVG(EXTRACT(EPOCH FROM (resolved_time - detected_time)) / 60) as avg_mttr_minutes
FROM incidents
WHERE severity IN ('P0', 'P1')
  AND resolved_time IS NOT NULL
  AND detected_time >= NOW() - INTERVAL '30 days';
```

**分级标准**:
```markdown
Elite: < 1小时
High: < 1天
Medium: 1天 - 1周
Low: > 1周
```

### 4. 变更失败率 (Change Failure Rate)

**度量方法**:
```sql
-- 导致事件或回滚的部署比例
SELECT
    (COUNT(CASE WHEN d.failed = true OR d.rollback = true THEN 1 END) * 100.0 / COUNT(*)) as change_failure_rate
FROM deployments d
WHERE d.environment = 'production'
  AND d.deployment_time >= NOW() - INTERVAL '30 days';
```

**分级标准**:
```markdown
Elite: 0-15%
High: 16-30%
Medium: 31-45%
Low: > 45%
```

---

**版本**: 1.0.0
**最后更新**: 2026-03-03
