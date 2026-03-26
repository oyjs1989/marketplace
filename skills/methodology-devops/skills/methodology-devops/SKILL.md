---
name: methodology-devops
description: 'This skill should be used when the user asks to "design a CI/CD pipeline", "set up deployment", "infrastructure as code", "monitoring strategy", "incident response plan", or mentions DevOps, continuous integration, continuous delivery, blue-green deployment, canary release, or observability. Provides DevOps methodology guidance for CI/CD, IaC, and operational excellence.'
version: 1.0.0
disable-model-invocation: true
---

# DevOps方法论技能

提供DevOps软件开发和运维的系统化实践指导,帮助团队实现持续集成、持续交付和快速反馈。

## 何时使用此技能

当用户请求以下任何内容时使用此技能:
- "DevOps"
- "CI/CD"
- "持续集成"
- "持续交付"
- "自动化部署"
- "监控和告警"
- "事件响应"
- "基础设施即代码"
- "部署管道"

## 核心原则

基于**DevOps文化**的3个核心支柱:

1. **文化 (Culture)**: 协作、信任、共同责任
2. **自动化 (Automation)**: 减少手动操作、提高可靠性
3. **度量 (Measurement)**: 数据驱动决策、持续改进
4. **共享 (Sharing)**: 知识共享、透明沟通

### CALMS框架

```markdown
C - Culture (文化): 打破壁垒,Dev和Ops协作
A - Automation (自动化): 自动化一切可自动化的
L - Lean (精益): 减少浪费,优化流程
M - Measurement (度量): 度量一切,数据驱动
S - Sharing (共享): 知识共享,开放透明
```

## DevOps实践

### 1. 持续集成 (CI)

**定义**: 频繁将代码集成到主干,自动构建和测试

**CI管道阶段**:
```markdown
代码提交 → 代码检查 → 单元测试 → 构建 → 集成测试 → 打包 → 制品存储
```

**实践清单**:
```markdown
- [ ] 使用版本控制 (Git)
- [ ] 主干开发或短期分支
- [ ] 自动化构建 (每次提交触发)
- [ ] 自动化测试 (单元、集成)
- [ ] 快速反馈 (<10分钟)
- [ ] 保持构建通过
- [ ] 每日至少集成一次
- [ ] 测试覆盖率 >80%
```

**CI工具链**:
- **版本控制**: Git, GitLab, GitHub
- **CI服务器**: Jenkins, GitLab CI, GitHub Actions, CircleCI
- **构建工具**: Maven, Gradle, npm, Make
- **代码检查**: SonarQube, ESLint, golangci-lint
- **制品仓库**: Nexus, Artifactory, Docker Registry

### 2. 持续交付/部署 (CD)

**持续交付 (Continuous Delivery)**:
- 代码随时可部署
- 需要人工批准部署

**持续部署 (Continuous Deployment)**:
- 全自动部署
- 测试通过自动上线

**CD管道阶段**:
```markdown
制品 → 部署到测试环境 → 自动化测试 → 部署到预发布 → 冒烟测试 → 部署到生产
```

**部署策略**:

#### 蓝绿部署 (Blue-Green)
```markdown
优点:
- 零停机时间
- 快速回滚

实施:
1. 准备新版本 (绿)
2. 新旧版本并行运行
3. 流量切换到绿
4. 蓝环境待命

适用: 有状态服务、关键业务
```

#### 金丝雀发布 (Canary)
```markdown
优点:
- 渐进式验证
- 风险可控

实施:
1. 部署新版本到少量节点 (5-10%)
2. 监控关键指标
3. 逐步扩大流量 (25% → 50% → 100%)
4. 异常则回滚

适用: 高风险变更、用户体验敏感
```

#### 滚动发布 (Rolling)
```markdown
优点:
- 资源利用高
- 逐步替换

实施:
1. 逐个节点更新
2. 更新后健康检查
3. 继续下一个节点

适用: 无状态服务、标准更新
```

**部署清单**:
```markdown
部署前:
- [ ] 代码审查通过
- [ ] 所有测试通过
- [ ] 变更文档准备
- [ ] 回滚计划明确
- [ ] 数据库迁移脚本
- [ ] 配置项准备

部署中:
- [ ] 健康检查通过
- [ ] 监控指标正常
- [ ] 日志无异常
- [ ] 关键功能验证

部署后:
- [ ] 冒烟测试通过
- [ ] 性能指标符合预期
- [ ] 告警无异常
- [ ] 用户反馈正常
```

### 3. 基础设施即代码 (IaC)

**定义**: 用代码管理基础设施,实现版本化、可重复、自动化

**IaC工具**:
```markdown
- Terraform: 跨云平台,声明式
- CloudFormation: AWS原生
- Ansible: 配置管理,过程式
- Puppet/Chef: 传统配置管理
- Pulumi: 使用编程语言 (TypeScript, Python)
```

**IaC实践**:
```markdown
- [ ] 所有基础设施代码化
- [ ] 使用版本控制
- [ ] 声明式配置优于过程式
- [ ] 模块化和复用
- [ ] 环境隔离 (dev/test/prod)
- [ ] 密钥管理 (Vault, KMS)
- [ ] 代码审查基础设施变更
```

**IaC模板示例**:
```hcl
# Terraform示例
resource "aws_instance" "web" {
  ami           = var.ami_id
  instance_type = "t3.medium"

  tags = {
    Name = "WebServer"
    Environment = var.environment
  }
}

resource "aws_autoscaling_group" "web" {
  desired_capacity = var.desired_capacity
  max_size         = var.max_size
  min_size         = var.min_size

  launch_template {
    id      = aws_launch_template.web.id
    version = "$Latest"
  }
}
```

### 4. 监控和可观测性

**可观测性三大支柱**:

#### Metrics (指标)
```markdown
基础设施指标:
- CPU使用率、内存使用率
- 磁盘I/O、网络带宽
- 容器/Pod健康状态

应用指标:
- QPS (每秒请求数)
- 响应时间 (P50, P95, P99)
- 错误率
- 并发连接数

业务指标:
- 注册用户数
- 订单量
- GMV
- 转化率
```

#### Logs (日志)
```markdown
日志类型:
- 应用日志 (INFO, WARN, ERROR)
- 访问日志 (Nginx, Apache)
- 系统日志 (syslog)
- 审计日志

日志实践:
- [ ] 结构化日志 (JSON)
- [ ] 统一日志格式
- [ ] 日志聚合 (ELK, Loki)
- [ ] 日志保留策略
- [ ] 敏感信息脱敏
```

#### Traces (追踪)
```markdown
分布式追踪:
- 请求链路追踪
- 调用关系可视化
- 性能瓶颈定位

工具:
- Jaeger
- Zipkin
- OpenTelemetry
- SkyWalking
```

**监控工具栈**:
```markdown
指标监控:
- Prometheus + Grafana (开源)
- Datadog (商业)
- New Relic (商业)
- CloudWatch (AWS)

日志聚合:
- ELK Stack (Elasticsearch + Logstash + Kibana)
- Loki + Grafana
- Splunk (商业)

追踪:
- Jaeger
- Zipkin
- APM工具 (New Relic, Datadog)
```

**告警策略**:
```markdown
告警原则:
- [ ] 可操作性 (收到告警能立即行动)
- [ ] 上下文丰富 (包含足够诊断信息)
- [ ] 避免告警疲劳 (减少误报)
- [ ] 分级响应 (P0/P1/P2/P3)
- [ ] 升级机制

告警类型:
- 阈值告警 (CPU > 80%)
- 趋势告警 (错误率持续上升)
- 异常检测 (ML based)
- 服务不可用

告警渠道:
- P0: 电话 + SMS + 即时消息
- P1: SMS + 即时消息
- P2: 即时消息
- P3: 邮件
```

### 5. 事件响应和SRE

**事件级别**:

| 级别 | 描述 | 响应时间 | 示例 |
|------|------|---------|------|
| P0 | 严重故障,核心业务不可用 | 15分钟 | 支付系统宕机 |
| P1 | 重大故障,主要功能受影响 | 1小时 | 用户无法登录 |
| P2 | 中等故障,部分功能受影响 | 4小时 | 辅助功能异常 |
| P3 | 轻微问题,用户体验下降 | 1天 | 界面显示错误 |

**事件响应流程**:
```markdown
1. 识别 (Detection)
   - [ ] 监控告警触发
   - [ ] 用户报告
   - [ ] 自动健康检查

2. 响应 (Response)
   - [ ] 确认事件级别
   - [ ] 启动响应团队
   - [ ] 建立沟通渠道 (War Room)
   - [ ] 记录时间线

3. 诊断 (Diagnosis)
   - [ ] 收集日志和指标
   - [ ] 分析错误模式
   - [ ] 识别根因
   - [ ] 评估影响范围

4. 缓解 (Mitigation)
   - [ ] 执行紧急修复
   - [ ] 或回滚到上个版本
   - [ ] 或启用降级措施
   - [ ] 验证恢复

5. 恢复 (Recovery)
   - [ ] 全面功能验证
   - [ ] 监控稳定性
   - [ ] 通知利益相关者

6. 事后分析 (Post-Mortem)
   - [ ] 编写事件报告
   - [ ] 根因分析 (5 Whys)
   - [ ] 改进行动计划
   - [ ] 知识库更新
```

**SRE实践**:
```markdown
- Error Budget (错误预算)
  - SLO: Service Level Objective (99.9%)
  - Error Budget = 100% - SLO = 0.1%
  - 用完则减缓发布,聚焦稳定性

- Toil管理 (减少重复手工劳动)
  - 目标: Toil < 50%工作时间
  - 自动化重复任务
  - 度量和减少Toil

- On-call轮值
  - 轮换机制 (1-2周)
  - 补偿时间
  - Runbook准备
```

### 6. 配置管理

**配置分类**:
```markdown
应用配置:
- 数据库连接串
- API端点
- 功能开关
- 业务参数

环境配置:
- dev / test / staging / prod
- 区域配置 (多区域部署)
- 租户配置 (SaaS)

密钥配置:
- API密钥
- 数据库密码
- 证书
```

**配置管理工具**:
```markdown
- Consul: 服务发现 + 配置中心
- etcd: 分布式键值存储
- Spring Cloud Config: Java生态
- Apollo: 携程开源
- Vault: 密钥管理

配置实践:
- [ ] 外部化配置 (不打包到代码)
- [ ] 环境变量优先
- [ ] 配置版本化
- [ ] 配置变更审计
- [ ] 热加载支持
- [ ] 密钥加密存储
```

## DevOps流程

### 完整DevOps流水线

```
计划 → 开发 → 构建 → 测试 → 发布 → 部署 → 运营 → 监控
 ↑                                                    ↓
 └────────────────── 反馈和改进 ←─────────────────────┘
```

**各阶段活动**:

#### 1. 计划 (Plan)
```markdown
- [ ] Sprint规划
- [ ] 用户故事编写
- [ ] 技术设计评审
- [ ] 风险评估
```

#### 2. 开发 (Code)
```markdown
- [ ] 使用Git分支策略 (GitFlow, Trunk-based)
- [ ] 代码审查 (Pull Request)
- [ ] 单元测试编写
- [ ] 本地开发环境
```

#### 3. 构建 (Build)
```markdown
- [ ] 自动化构建
- [ ] 依赖管理
- [ ] 编译优化
- [ ] 制品生成
```

#### 4. 测试 (Test)
```markdown
- [ ] 单元测试
- [ ] 集成测试
- [ ] API测试
- [ ] UI测试 (Selenium, Cypress)
- [ ] 性能测试 (JMeter, k6)
- [ ] 安全测试 (SAST, DAST)
```

#### 5. 发布 (Release)
```markdown
- [ ] 版本标签
- [ ] Release Notes
- [ ] 制品签名
- [ ] 变更审批
```

#### 6. 部署 (Deploy)
```markdown
- [ ] 部署到测试环境
- [ ] 自动化部署脚本
- [ ] 蓝绿/金丝雀部署
- [ ] 健康检查
```

#### 7. 运营 (Operate)
```markdown
- [ ] 服务监控
- [ ] 日志聚合
- [ ] 容量规划
- [ ] 成本优化
```

#### 8. 监控 (Monitor)
```markdown
- [ ] 应用性能监控 (APM)
- [ ] 业务指标监控
- [ ] 用户体验监控
- [ ] 告警配置
```

## 中文输出格式

**所有输出必须使用中文**,遵循以下结构:

```markdown
# DevOps实施指南: [项目名称]

## CI/CD管道设计

### 管道阶段
```mermaid
graph LR
    A[代码提交] --> B[代码检查]
    B --> C[单元测试]
    C --> D[构建]
    D --> E[集成测试]
    E --> F[部署到测试环境]
    F --> G[自动化测试]
    G --> H[部署到生产]
```

### 工具链
| 阶段 | 工具 | 配置 |
|------|------|------|
| 源代码管理 | GitLab | [配置详情] |
| CI服务器 | Jenkins | [配置详情] |
| 制品仓库 | Nexus | [配置详情] |
| 部署工具 | Ansible | [配置详情] |

## 部署策略

### 选择的策略: 蓝绿部署
**理由**: [为什么选择这个策略]

**实施步骤**:
1. [ ] 准备绿环境
2. [ ] 部署新版本
3. [ ] 健康检查
4. [ ] 流量切换
5. [ ] 监控观察
6. [ ] 蓝环境保留24小时

**回滚计划**:
- 触发条件: [错误率 > 1%, 响应时间 > 500ms]
- 回滚步骤: [流量切回蓝环境]
- 预期时间: 5分钟

## 监控方案

### 关键指标
**基础设施指标**:
- CPU使用率: 目标 < 70%, 告警 > 80%
- 内存使用率: 目标 < 80%, 告警 > 90%
- 磁盘使用率: 目标 < 70%, 告警 > 85%

**应用指标**:
- QPS: 正常范围 [1000-5000]
- P99响应时间: < 500ms
- 错误率: < 0.1%
- 可用性: > 99.9%

**业务指标**:
- [自定义业务指标]

### 告警配置
| 告警名称 | 级别 | 条件 | 响应时间 | 通知渠道 |
|---------|------|------|---------|---------|
| 服务不可用 | P0 | 连续3次健康检查失败 | 15分钟 | 电话+SMS |
| CPU过高 | P1 | CPU > 90% 持续5分钟 | 1小时 | SMS |
| 错误率上升 | P1 | 错误率 > 1% | 30分钟 | 即时消息 |

## 事件响应计划

### 响应团队
- **Incident Commander**: [姓名]
- **技术负责人**: [姓名]
- **沟通负责人**: [姓名]

### 联系方式
- 紧急电话: [电话]
- Slack频道: #incident-response
- War Room: [Zoom链接]

### Runbook
#### P0事件: 服务完全不可用
1. **确认**:
   - [ ] 检查监控仪表盘
   - [ ] 确认影响范围
   - [ ] 建立War Room

2. **诊断**:
   - [ ] 查看日志聚合平台
   - [ ] 检查最近部署
   - [ ] 分析错误模式

3. **缓解**:
   - [ ] 选项A: 回滚到上个版本
   - [ ] 选项B: 执行紧急修复
   - [ ] 选项C: 启用降级模式

4. **验证**:
   - [ ] 健康检查通过
   - [ ] 错误率恢复正常
   - [ ] 用户反馈正常

## 基础设施即代码

### IaC工具: Terraform
### 模块结构
```
terraform/
├── modules/
│   ├── vpc/
│   ├── compute/
│   ├── database/
│   └── monitoring/
├── environments/
│   ├── dev/
│   ├── test/
│   └── prod/
└── shared/
```

### 部署流程
```bash
# 1. 规划变更
terraform plan -var-file=environments/prod/terraform.tfvars

# 2. 应用变更
terraform apply -var-file=environments/prod/terraform.tfvars

# 3. 验证资源
terraform show
```

## 改进建议

### 短期 (1个月内)
- [ ] [改进项1]
- [ ] [改进项2]

### 中期 (3个月内)
- [ ] [改进项3]
- [ ] [改进项4]

### 长期 (6个月内)
- [ ] [改进项5]
- [ ] [改进项6]

## 度量指标

### DORA四大指标
| 指标 | 当前值 | 目标值 | 行业基准 |
|------|--------|--------|---------|
| 部署频率 | [X次/天] | [Y次/天] | Elite: 多次/天 |
| 变更前置时间 | [X小时] | [Y小时] | Elite: < 1小时 |
| 平均恢复时间 (MTTR) | [X分钟] | [Y分钟] | Elite: < 1小时 |
| 变更失败率 | [X%] | [Y%] | Elite: 0-15% |
```

## 语言要求

- ✅ **所有文档和配置**: 必须使用中文
- ✅ **技术术语**: 首次使用提供英文,后续使用中文
- ✅ **工具名称**: 保留英文原名 + 中文注释
- ✅ **代码和脚本**: 保持原语言,注释使用中文

## 质量标准

- [ ] CI管道每次提交自动触发
- [ ] 构建时间 < 10分钟
- [ ] 测试覆盖率 > 80%
- [ ] 自动化测试通过率 > 95%
- [ ] 部署全自动化 (测试环境)
- [ ] 生产部署有审批流程
- [ ] 回滚时间 < 5分钟
- [ ] 关键服务可用性 > 99.9%
- [ ] P0事件平均恢复时间 < 1小时
- [ ] 所有基础设施代码化

## DevOps成熟度评估

### Level 1: 初始级
```markdown
特征:
- 手动部署
- 无自动化测试
- 无持续集成
- 反应式运维

改进建议:
- 引入版本控制
- 自动化构建
- 建立基础监控
```

### Level 2: 管理级
```markdown
特征:
- 基本CI/CD管道
- 部分自动化测试
- 代码审查流程
- 基础监控

改进建议:
- 提高测试覆盖率
- 自动化部署
- 完善监控告警
```

### Level 3: 定义级
```markdown
特征:
- 完整CI/CD管道
- 高测试覆盖率
- 自动化部署
- 全面监控

改进建议:
- 引入IaC
- 蓝绿/金丝雀部署
- 混沌工程
```

### Level 4: 量化级
```markdown
特征:
- 基础设施即代码
- 多种部署策略
- 可观测性完善
- 数据驱动决策

改进建议:
- 机器学习运维 (AIOps)
- 自愈系统
- 预测性维护
```

### Level 5: 优化级
```markdown
特征:
- 全面自动化
- 自服务平台
- AIOps实践
- 持续优化文化

保持并发展:
- 创新实验
- 行业领先实践
```

## 常见挑战

### 挑战1: 文化阻力
```markdown
症状: Dev和Ops缺乏协作
解决:
- 建立联合团队
- 共享On-call职责
- 统一目标和激励
- 定期沟通会议
```

### 挑战2: 工具链复杂
```markdown
症状: 工具过多,维护困难
解决:
- 整合工具链
- 标准化流程
- 集中式平台
- 培训和文档
```

### 挑战3: 遗留系统
```markdown
症状: 老系统难以自动化
解决:
- 渐进式改造
- Strangler Pattern
- API抽象层
- 容器化封装
```

## 参考文档

详细实践指南请参考:
- `references/DEVOPS_FRAMEWORK.md` - DevOps框架详解

## 版本历史

- v1.0.0 (2026-03-03): 初始版本
  - CI/CD管道设计
  - 部署策略 (蓝绿、金丝雀、滚动)
  - 基础设施即代码
  - 监控和可观测性
  - 事件响应流程
  - 中文模板
