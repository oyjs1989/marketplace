一、我们讨论了什么

   1. 分析了 spec-kit 项目
      - 规格驱动开发 (SDD) 方法论
      - 多 Agent 支持 (20+)
      - 命令模板 + 脚本的协作模式
      - 分阶段流水线：specify → clarify → plan → tasks → implement

   2. 提取了核心设计经验
      - 配置驱动（单一数据源）
      - 模板约束 LLM 行为
      - 职责分离（脚本做机械，命令做智能）
      - Gate + Checklist 质量机制
      - 工具中立，可扩展

   3. 确定 Code Review 工具定位
      - ✅ 先做单独 Agent 工具
      - ❌ 不急于做方法论框架
      - 原因：Code Review 是检查点，不是端到端流程

   4. 识别了当前问题的本质
      - 规则不被遵循 → 规则太自由文本，缺乏结构
      - 业务逻辑判断难 → 业务知识没有显式化

  ---

  二、方案：结构化 Code Review Agent

  目录结构

   .review/
   ├── commands/
   │   ├── review.md          # 主命令模板
   │   ├── rules.md           # 规则管理命令
   │   └── learn.md           # 学习命令
   ├── rules/
   │   ├── security.yaml      # 安全规则
   │   ├── performance.yaml   # 性能规则
   │   ├── architecture.yaml  # 架构规则
   │   └── business.yaml      # 业务规则
   ├── contracts/
   │   ├── workflows.yaml     # 业务流程契约
   │   └── dependencies.yaml  # 依赖关系契约
   └── templates/
       └── report.md          # 审查报告模板

  规则结构设计

   # rules/security.yaml
   rules:
     - id: SEC-001
       name: SQL注入防护
       severity: error
       description: 所有数据库查询必须使用参数化
       pattern:
         type: regex
         match: "(execute|query)\\s*\\(['\"]"
       check:
         - type: pattern_absent
           pattern: "execute\\s*\\(['\"].*\\+"
       message: "检测到可能的SQL注入风险，请使用参数化查询"
       examples:
         bad: 'db.execute("SELECT * FROM users WHERE id = " + userId)'
         good: 'db.execute("SELECT * FROM users WHERE id = ?", [userId])'

     - id: SEC-002
       name: 敏感数据暴露
       severity: warning
       pattern:
         type: regex
         match: "(password|token|secret|key)\\s*[:=]\\s*['\"]"
       message: "检测到硬编码的敏感数据"

  业务契约设计

   # contracts/workflows.yaml
   workflows:
     - name: 用户支付流程
       steps:
         - id: validate_cart
           must_before: [check_inventory]
         - id: check_inventory
           must_before: [create_order]
         - id: create_order
           must_before: [process_payment]
         - id: process_payment
           must_before: [send_confirmation]

       constraints:
         - "process_payment 必须在 create_order 成功后"
         - "check_inventory 失败时不能执行 create_order"

     - name: 数据导出流程
       steps:
         - id: check_permission
           must_before: [validate_query]
         - id: validate_query
           must_before: [execute_query]
         - id: execute_query
           must_before: [format_output]

  命令设计

   # 主命令
   /review                    # 审查当前文件/目录
   /review --diff HEAD~1      # 审查最近的变更
   /review --pr 123           # 审查指定 PR

   # 规则管理
   /review.rules list         # 列出所有规则
   /review.rules add          # 添加新规则
   /review.rules enable SEC-001
   /review.rules disable SEC-002

   # 学习命令
   /review.learn --from-pr 123,124,125  # 从历史 PR 学习
   /review.learn --from-commit abc123   # 从指定提交学习

  命令模板核心逻辑

   ---
   description: 执行代码审查，检查规则和契约
   scripts:
     sh: scripts/review.sh --json
   ---

   ## Outline

   1. 加载规则文件
      - 读取 rules/*.yaml
      - 解析 pattern 和 check 条件
      - 加载 contracts/*.yaml

   2. 分析目标代码
      - 解析代码结构（AST/正则）
      - 提取关键模式
      - 识别代码块和调用关系

   3. 执行检查
      - 遍历每条规则
      - 执行 pattern 匹配
      - 执行 check 条件
      - 检查契约约束

   4. 生成报告
      - 汇总违规项
      - 按严重程度排序
      - 提供修复建议

   5. 输出格式
      - 控制台：简洁摘要
      - 文件：详细报告

  ---

  三、你需要审视的内容

  回去看你的项目时，请关注：

   1. 现有规则是如何定义的？
      □ 自由文本？
      □ 有一定结构？
      □ 可以转换为 YAML 格式吗？

   2. 业务逻辑是否有明确的流程？
      □ 能否列出关键业务流程？
      □ 每个流程的步骤顺序是什么？
      □ 步骤之间有什么约束？

   3. 现有工具的输出是什么？
      □ 输出格式是什么样的？
      □ 漏报/误报的模式是什么？
      □ 哪些规则经常不被遵循？

   4. 技术栈和代码模式
      □ 主要用什么语言？
      □ 代码结构有什么特点？
      □ 常见的代码模式有哪些？

  ---

  四、下一步行动

   1. 检查现有项目的规则定义方式
      → 评估能否结构化

   2. 梳理核心业务流程
      → 写成 workflows.yaml

   3. 分析典型代码模式
      → 设计 pattern 规则

   4. 选择一个最小范围测试
      → 验证结构化规则是否更有效

   5. 根据结果迭代
      → 逐步完善规则库