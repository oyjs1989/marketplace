# Golang工程规范

## 错误处理
category: "错误处理"
description: "业务函数禁止直接panic，必须返回error"

# 证据链（必须可追溯）
evidence:
  internal/app/biz/param_template.go:224 评论锚点位置

# 教学示例
good_example: |
  // 正确：显式返回error
  func Process(data []byte) error {
      if len(data) == 0 {
          return fmt.Errorf("invalid input: data is empty")
      }
      return nil
  }
  
bad_example: |
  // 错误：不可预测的panic
  func Process(data []byte) {
      if len(data) == 0 {
          panic("empty data")  // 调用方无法优雅处理
      }
  }
  
rationale: "panic会导致整个服务崩溃，而error可被调用链捕获和处理"

---

## 日志记录
category: "日志记录"
description: "日志字段使用明确的类型，field name使用snake_case"

# 证据链（必须可追溯）
evidence:
  internal/app/biz/approval_mgr.go:25 评论锚点位置

# 教学示例
good_example: |
  // 正确：使用明确的日志字段类型和snake_case命名
  log.Info(ctx, "get uin list success", log.Int("size", len(uinList)))
  
bad_example: |
  // 错误：类型不明确或命名不规范
  log.Info(ctx, "get uin list success", log.Any("size", len(uinList)))
  log.Info(ctx, "get uin list success", log.Int("Size", len(uinList)))

rationale: "明确的类型和snake_case命名有助于日志的一致性和可读性"

---

## 代码结构
category: "代码结构"
description: "所有switch/case都必须携带default分支"

# 证据链（必须可追溯）
evidence:
  internal/app/biz/approval_mgr.go:36 评论锚点位置

# 教学示例
good_example: |
  // 正确：包含default分支
  switch status {
  case StatusSuccess:
      return "success"
  case StatusFailed:
      return "failed"
  default:
      return "unknown"
  }
  
bad_example: |
  // 错误：缺少default分支
  switch status {
  case StatusSuccess:
      return "success"
  case StatusFailed:
      return "failed"
  }

rationale: "default分支确保所有可能的情况都被处理，避免未定义行为"

---

## 数据结构
category: "数据结构"
description: "结构体字段名称尽量跟tag命名保持一致"

# 证据链（必须可追溯）
evidence:
  internal/app/biz/bo/param_template.go:16 评论锚点位置

# 教学示例
good_example: |
  // 正确：字段名与tag命名一致
  type ParamTemplate struct {
      ID            int64                         `gorm:"column:id"`
      TemplateDesc  *string                       `gorm:"column:template_desc"`
      ProviderCode  string                        `gorm:"column:provider_code"`
  }
  
bad_example: |
  // 错误：字段名与tag不一致
  type ParamTemplate struct {
      ID            int64                         `gorm:"column:id"`
      Desc          *string                       `gorm:"column:template_desc"`
      ProviderCode  string                        `gorm:"column:provider_code"`
  }

rationale: "保持字段名与tag命名一致有助于代码可读性和维护性"

---

## 常量命名
category: "常量命名"
description: "常量使用大写字母和下划线组合，命名需有意义"

# 证据链（必须可追溯）
evidence:
  internal/app/constants/param_template.go:35 评论锚点位置

# 教学示例
good_example: |
  // 正确：使用大写字母和下划线，命名有意义
  const (
      TEMP_TEMPLATE_NAME_PREFIX = "fcmp_temp_template_%s"
      TEMP_TEMPLATE_DESC        = "获取系统参数而创建的临时模版"
  )
  
bad_example: |
  // 错误：命名不规范或无意义
  const (
      TencentCloudType = "tencent"
      AWS_CLOUD_TYPE   = "aws"
  )

rationale: "统一的常量命名规范有助于代码的一致性和可读性"

---

## 接口设计
category: "接口设计"
description: "获取列表的方法名应使用GetList或Search而非Get"

# 证据链（必须可追溯）
evidence:
  internal/app/biz/repo/param_template.go:21 评论锚点位置

# 教学示例
good_example: |
  // 正确：使用GetList或Search表示获取多个元素
  type ParamTemplateRepo interface {
      GetList(ctx context.Context, productCode *string) ([]*bo.ParamTemplate, int64, error)
      Search(ctx context.Context, productCode *string) ([]*bo.ParamTemplate, int64, error)
  }
  
bad_example: |
  // 错误：使用Get表示获取列表
  type ParamTemplateRepo interface {
      Get(ctx context.Context, productCode *string) ([]*bo.ParamTemplate, int64, error)
  }

rationale: "Get通常表示获取单个元素，GetList或Search更明确表示获取多个元素"

---

## 参数处理
category: "参数处理"
description: "创建和更新操作中，参数命名应与具体操作匹配"

# 证据链（必须可追溯）
evidence:
  internal/app/biz/repo/param_template.go:42 评论锚点位置

# 教学示例
good_example: |
  // 正确：参数命名与操作匹配
  type ParamTemplateRepo interface {
      // 创建时不需要templateID
      Create(ctx context.Context, template *bo.CreateParamTemplateParam, status constants.ParamTemplateStatus) (int64, error)
      // 更新时需要templateID
      Update(ctx context.Context, id int64, templateID *string, name *string, desc *string) error
  }
  
bad_example: |
  // 错误：参数命名不明确或冗余
  type ParamTemplateRepo interface {
      Create(ctx context.Context, template *bo.CreateParamTemplateParam, templateID *string, status constants.ParamTemplateStatus) (int64, error)
      Update(ctx context.Context, id int64, name *string, desc *string) error
  }

rationale: "明确的参数命名有助于理解操作意图，避免混淆"

---

## 代码组织
category: "代码组织"
description: "业务字段类型的命名应避免歧义，使用明确的后缀"

# 证据链（必须可追溯）
evidence:
  internal/app/biz/bo/param_template.go:49 评论锚点位置

# 教学示例
good_example: |
  // 正确：使用明确后缀避免歧义
  type ParamTemplateSyncInfo struct {
      ID               int64                `gorm:"column:id"`
      TemplateID       string               `gorm:"column:template_id"`
      Uin              string               `gorm:"column:uin"`
      RegionCode       string               `gorm:"column:region_code"`
      ParamTemplateTid int64                `gorm:"column:param_template_tid"`
  }
  
bad_example: |
  // 错误：命名可能导致歧义
  type ParamTemplateSync struct {
      ID               int64                `gorm:"column:id"`
      TemplateID       string               `gorm:"column:template_id"`
      Uin              string               `gorm:"column:uin"`
      RegionCode       string               `gorm:"column:region_code"`
      ParamTemplateTid int64                `gorm:"column:param_template_tid"`
  }

rationale: "明确的命名避免了'ParamTemplateSync'看起来像动作的歧义"

---

## 注释规范
category: "注释规范"
description: "注释应有每种类型的中文描述，对齐且位置规范"

# 证据链（必须可追溯）
evidence:
  internal/app/biz/bo/approval_mgr.go:16 评论锚点位置

# 教学示例
good_example: |
  // 正确：注释有中文描述，对齐且位置规范
  type ApprovalStatus int32
  
  const (
      ApprovalStatusInit     ApprovalStatus = 0 // 待审批
      ApprovalStatusApproved ApprovalStatus = 1 // 已同意
      ApprovalStatusRejected ApprovalStatus = 2 // 已驳回
  )
  
bad_example: |
  // 错误：注释缺失或位置不规范
  type ApprovalStatus int32
  
  const (
      ApprovalStatusInit     ApprovalStatus = 0
      ApprovalStatusApproved ApprovalStatus = 1 // 已同意
      ApprovalStatusRejected ApprovalStatus = 2
  )

rationale: "规范的注释提高代码可读性和可维护性"

---

## 命名规范
category: "命名规范"
description: "避免一个功能使用多种命名方式，保持命名一致性"

# 证据链（必须可追溯）
evidence:
  internal/app/constants/param_template.go:27 评论锚点位置

# 教学示例
good_example: |
  // 正确：保持命名一致性
  const (
      CloudTencent    = "tencent"
      CloudAliyun     = "aliyun"
      CloudAws        = "aws"
  )
  
bad_example: |
  // 错误：命名方式不一致
  const (
      TencentCloud    = "tencent"
      CloudAliyun     = "aliyun"
      AWS_CLOUD       = "aws"
  )

rationale: "一致的命名方式减少混淆，提高代码可读性"