# Go代码审查规范

## 1. 错误处理规范

### 规则1: 错误包装必须使用errors.Wrap或errors.WithMessage
- **描述**: 任何错误返回必须使用errors.Wrap或errors.WithMessage进行错误包装，不允许直接返回fmt.Errorf
- **反面示例**:
```go
// 错误示例
return fmt.Errorf("Subscribe %s failed", fileName)
```
- **正面示例**:
```go
// 正确示例
return errors.Wrapf(err, "Subscribe %s failed", fileName)
return errors.WithMessage(err, "get role failed")
```
- **理由**: 通过errors.Wrap和errors.WithMessage可以保留错误栈信息，便于调试和定位问题

### 规则2: 错误信息放在日志字段末尾
- **描述**: log.ErrorField(err)必须放在日志字段的最后
- **正面示例**:
```go
log.Error(ctx, "failed to get account config", log.String("uin", account.Uin), log.ErrorField(err))
```
- **理由**: 保持错误日志的一致性，便于错误追踪和分析

### 规则3: 错误必须放在返回值的最后
- **描述**: 函数返回值中错误必须在最后位置
- **正面示例**:
```go
func ChangeToDepartPermissionType(val pb.PermissionType) (result type, err error) {...}
```
- **理由**: 符合Go的错误处理标准，便于统一处理错误

### 规则4: 业务函数禁止直接panic，必须返回error
- **描述**: 业务函数禁止直接panic，必须返回error
- **正面示例**:
```go
// 正确：显式返回error
func Process(data []byte) error {
    if len(data) == 0 {
        return fmt.Errorf("invalid input: data is empty")
    }
    return nil
}
```
- **反面示例**:
```go
// 错误：不可预测的panic
func Process(data []byte) {
    if len(data) == 0 {
        panic("empty data")  // 调用方无法优雅处理
    }
}
```
- **理由**: panic会导致整个服务崩溃，而error可被调用链捕获和处理

## 2. 日志规范

### 规则5: 日志字段使用snake_case命名
- **描述**: 日志字段的key必须使用snake_case命名风格
- **反面示例**:
```go
log.String("UIN", account.Uin)
```
- **正面示例**:
```go
log.String("uin", account.Uin)
```
- **理由**: 保持日志字段命名的一致性，便于日志查询和分析

### 规则6: 日志字段使用明确的类型
- **描述**: 日志字段使用明确的类型，field name使用snake_case
- **正面示例**:
```go
// 正确：使用明确的日志字段类型和snake_case命名
log.Info(ctx, "get uin list success", log.Int("size", len(uinList)))
```
- **反面示例**:
```go
// 错误：类型不明确或命名不规范
log.Info(ctx, "get uin list success", log.Any("size", len(uinList)))
log.Info(ctx, "get uin list success", log.Int("Size", len(uinList)))
```
- **理由**: 明确的类型和snake_case命名有助于日志的一致性和可读性

### 规则7: 避免在data层打日志
- **描述**: data层（数据库操作层）不应该打印日志
- **理由**: data层应该专注于数据操作，日志属于业务逻辑层面，应由上层处理

### 规则8: 日志输出必要的上下文信息
- **描述**: 日志必须包含足够的上下文信息，便于问题定位
- **正面示例**:
```go
log.Info(ctx, "some operation", log.String("alias", alias), log.Int64("id", id))
```
- **理由**: 充足的上下文信息有助于快速定位问题

## 3. 数据库操作规范

### 规则9: GORM查询必须明确指定查询条件
- **描述**: 使用GORM查询时必须明确使用Where条件
- **正面示例**:
```go
db.Where("id = ?", id).Model(&model.Role{}).Take(&role).Error
```
- **理由**: 避免查询条件不明确导致的数据错误

### 规则10: GORM查询必须明确指定查询列
- **描述**: 不要使用SELECT *，必须明确指定查询的列
- **正面示例**:
```go
db.Select("group_code, permission_list").Where(...).Find(&pg).Error
```
- **理由**: 提高查询效率，减少不必要的数据传输

### 规则11: 避免使用GORM的save方法
- **描述**: 项目内不允许使用GORM的save方法，应使用update等其他方法
- **反面示例**:
```go
db.Save(&model)
```
- **正面示例**:
```go
db.Updates(&model)
```
- **理由**: save方法会更新所有字段，可能导致意外的数据覆盖

## 4. 变量和命名规范

### 规则12: 指针变量必须进行空指针校验
- **描述**: 使用指针变量前必须进行空指针校验
- **正面示例**:
```go
if somePtr != nil {
    result := *somePtr
}
```
- **理由**: 防止空指针引用导致程序panic

### 规则13: ID统一写成ID而不是Id
- **描述**: 除了生成代码外，ID统一写成ID，而不是Id
- **正面示例**:
```go
ID string
```
- **理由**: 保持ID命名的一致性

### 规则14: 接口命名规范
- **描述**: 接口命名应保持简单准确，以-er后缀命名
- **理由**: 接口命名应清晰准确，便于理解

### 规则15: 常量使用大写字母和下划线组合
- **描述**: 常量使用大写字母和下划线组合，命名需有意义
- **正面示例**:
```go
// 正确：使用大写字母和下划线，命名有意义
const (
    TEMP_TEMPLATE_NAME_PREFIX = "fcmp_temp_template_%s"
    TEMP_TEMPLATE_DESC        = "获取系统参数而创建的临时模版"
)
```
- **理由**: 统一的常量命名规范有助于代码的一致性和可读性

### 规则16: 常量统一定义在constants包中
- **描述**: 字符串、数字等如果有统一管理、复用等需求，统一定义在常量包内
- **理由**: 避免魔法字符串，提高代码可维护性

### 规则17: 避免一个功能使用多种命名方式
- **描述**: 避免一个功能使用多种命名方式，保持命名一致性
- **正面示例**:
```go
// 正确：保持命名一致性
const (
    CloudTencent    = "tencent"
    CloudAliyun     = "aliyun"
    CloudAws        = "aws"
)
```
- **理由**: 一致的命名方式减少混淆，提高代码可读性

### 规则18: 结构体字段名称尽量跟tag命名保持一致
- **描述**: 结构体字段名称尽量跟tag命名保持一致
- **正面示例**:
```go
// 正确：字段名与tag命名一致
type ParamTemplate struct {
    ID            int64                         `gorm:"column:id"`
    TemplateDesc  *string                       `gorm:"column:template_desc"`
    ProviderCode  string                        `gorm:"column:provider_code"`
}
```
- **理由**: 保持字段名与tag命名一致有助于代码可读性和维护性

### 规则19: 结构体命名采用驼峰命名法，首字母大写
- **描述**: 所有结构体名称都使用驼峰命名法，首字母大写
- **理由**: 统一结构体命名规范

### 规则20: 结构体字段命名采用驼峰命名法，首字母大写
- **描述**: 结构体字段名使用驼峰命名法，首字母大写
- **理由**: 统一结构体字段命名规范

### 规则21: 函数和方法命名采用驼峰命名法
- **描述**: 私有函数/方法：首字母小写，公共函数/方法：首字母大写
- **理由**: 统一函数和方法命名规范

### 规则22: 业务字段类型的命名应避免歧义
- **描述**: 业务字段类型的命名应避免歧义，使用明确的后缀
- **正面示例**:
```go
// 正确：使用明确后缀避免歧义
type ParamTemplateSyncInfo struct {
    ID               int64                `gorm:"column:id"`
    TemplateID       string               `gorm:"column:template_id"`
    Uin              string               `gorm:"column:uin"`
    RegionCode       string               `gorm:"column:region_code"`
    ParamTemplateTid int64                `gorm:"column:param_template_tid"`
}
```
- **理由**: 明确的命名避免了歧义

## 5. 代码结构和组织规范

### 规则23: switch/case必须带default分支
- **描述**: 原则上所有的switch/case都必须带上default分支
- **正面示例**:
```go
switch val {
case 1:
    // do something
case 2:
    // do something
default:
    // 处理其他情况
}
```
- **理由**: 提高代码健壮性，处理未预期的情况

### 规则24: 所有switch/case都必须携带default分支
- **描述**: 所有switch/case都必须携带default分支
- **正面示例**:
```go
// 正确：包含default分支
switch status {
case StatusSuccess:
    return "success"
case StatusFailed:
    return "failed"
default:
    return "unknown"
}
```
- **理由**: default分支确保所有可能的情况都被处理，避免未定义行为

### 规则25: 嵌入的结构体永远在最顶层
- **描述**: 嵌入的结构体永远在最顶层
- **正面示例**:
```go
type AccountFavor struct {
    orm.BaseWithCTime
    Alias string `gorm:"column:alias;type:varchar(128);not null;comment:别名"`
}
```
- **理由**: 保持结构体字段组织的一致性

### 规则26: 初始化map和切片时必须通过make进行初始化
- **描述**: 初始化map和切片时，除了unmarshal、scan等场景，必须通过make进行初始化
- **正面示例**:
```go
slice := make([]string, 0)  // 或 make([]string, 0, capacity)
```
- **理由**: 使用make初始化更明确，且可以预设容量优化性能

### 规则27: 一般情况下不使用iota作为自增常量
- **描述**: 一般情况下不使用iota作为自增常量，应明确指定值
- **正面示例**:
```go
type Status int
const (
    StatusEnable Status = 0
    StatusDisable Status = 1
)
```
- **理由**: 明确指定值可以防止枚举顺序变化导致的错误

### 规则28: 禁止使用全局变量
- **描述**: 无特殊场景（例如单测、正则预编译），禁止使用全局变量
- **理由**: 避免全局状态导致的并发问题和测试困难

### 规则29: 除了server包需要注册服务以外，其他包禁止使用init函数
- **描述**: 除了server包需要注册服务以外，其他包禁止使用init函数
- **理由**: 避免不可控的初始化行为

## 6. 并发和错误处理规范

### 规则30: 使用common包的recovered.ErrorGroup
- **描述**: 并发操作应该使用common包中的recovered.ErrorGroup
- **正面示例**:
```go
errGroup := recovered.NewErrorGroup()
```
- **理由**: recovered.ErrorGroup提供更好的错误处理和取消机制

### 规则31: 并发控制使用limiter
- **描述**: 控制并发数量应使用limiter.NewConcurrentLimiter
- **正面示例**:
```go
l := limiter.NewConcurrentLimiter(limiterCount)
```
- **理由**: 防止并发过多导致资源耗尽

## 7. 配置管理规范

### 规则32: 配置更新处理
- **描述**: 配置中心更新后，相关配置应能够正确更新
- **正面示例**:
```go
// 通过回调函数处理配置更新
func (p *ProxyConfig) OnConfigChange() (err error) {
    // 处理配置更新
}
```
- **理由**: 保证配置能够动态更新，无需重启服务

## 8. 测试规范

### 规则33: 单元测试必须验证具体结构
- **描述**: 单元测试必须验证具体的结构，不能只验证err是否为nil
- **正面示例**:
```go
assert.NoError(t, err)
assert.Equal(t, expectedValue, actualValue)  // 验证具体值
```
- **理由**: 验证具体结构确保功能正确性

## 9. 接口设计规范

### 规则34: 获取列表的方法名应使用GetList或Search而非Get
- **描述**: 获取列表的方法名应使用GetList或Search而非Get
- **正面示例**:
```go
// 正确：使用GetList或Search表示获取多个元素
type ParamTemplateRepo interface {
    GetList(ctx context.Context, productCode *string) ([]*bo.ParamTemplate, int64, error)
    Search(ctx context.Context, productCode *string) ([]*bo.ParamTemplate, int64, error)
}
```
- **理由**: Get通常表示获取单个元素，GetList或Search更明确表示获取多个元素

### 规则35: 接口设计遵循最小原则
- **描述**: 接口设计遵循最小原则，只包含必要的方法
- **理由**: 减少接口复杂性，提高可维护性

## 10. 其他规范

### 规则36: 创建和更新操作中参数命名应与具体操作匹配
- **描述**: 创建和更新操作中，参数命名应与具体操作匹配
- **正面示例**:
```go
// 正确：参数命名与操作匹配
type ParamTemplateRepo interface {
    // 创建时不需要templateID
    Create(ctx context.Context, template *bo.CreateParamTemplateParam, status constants.ParamTemplateStatus) (int64, error)
    // 更新时需要templateID
    Update(ctx context.Context, id int64, templateID *string, name *string, desc *string) error
}
```
- **理由**: 明确的参数命名有助于理解操作意图，避免混淆

### 规则37: 注释应有每种类型的中文描述
- **描述**: 注释应有每种类型的中文描述，对齐且位置规范
- **正面示例**:
```go
// 正确：注释有中文描述，对齐且位置规范
type ApprovalStatus int32

const (
    ApprovalStatusInit     ApprovalStatus = 0 // 待审批
    ApprovalStatusApproved ApprovalStatus = 1 // 已同意
    ApprovalStatusRejected ApprovalStatus = 2 // 已驳回
)
```
- **理由**: 规范的注释提高代码可读性和可维护性

### 规则38: 函数参数较多时使用结构体封装
- **描述**: 当函数参数较多时，通过结构体封装参数传递，避免长参数列表
- **理由**: 提高代码可读性和可维护性

### 规则39: 上下文（context）作为函数的第一个参数
- **描述**: 所有需要上下文的函数，ctx作为第一个参数传入
- **理由**: 统一上下文传递方式

### 规则40: 返回错误时立即处理，避免延迟处理
- **描述**: 错误检查紧跟在可能产生错误的函数调用之后
- **理由**: 及时处理错误，避免错误传播

### 规则41: 注释格式统一，公共函数必须有注释
- **描述**: 公共函数/方法必须有注释说明功能、参数和返回值
- **理由**: 提高代码可读性和可维护性

### 规则42: 包名与目录名保持一致
- **描述**: Go包名与目录名保持一致
- **理由**: 保持项目结构一致性

### 规则43: 包导入按标准库、第三方库、内部库分组
- **描述**: 包导入按标准库、第三方库、内部库顺序分组，每组之间空一行
- **理由**: 提高代码可读性

### 规则44: 测试文件命名规范
- **描述**: 测试文件以`_test.go`结尾，与被测试文件保持一致的命名
- **理由**: 保持测试文件命名一致性

### 规则45: 测试函数命名规范
- **描述**: 测试函数以`Test`开头，后跟被测试函数名
- **理由**: 保持测试函数命名一致性

### 规则46: 业务逻辑与数据访问分离
- **描述**: 业务逻辑在`internal/app/biz/`目录下实现，数据访问在`internal/app/data/`目录下实现
- **理由**: 分离关注点，提高代码可维护性

### 规则47: 云厂商相关代码按厂商分目录组织
- **描述**: 不同云厂商的实现代码按厂商分目录
- **理由**: 保持项目结构清晰

### 规则48: 配置相关代码统一放在config目录
- **描述**: 配置管理相关代码放在`internal/app/config/`目录下
- **理由**: 统一配置管理

### 规则49: 工具类代码放在pkg目录
- **描述**: 公共工具类代码放在`internal/pkg/`目录下
- **理由**: 统一工具类代码管理

### 规则50: 错误日志必须包含错误字段
- **描述**: 记录错误日志时必须使用`log.ErrorField(err)`包含错误信息
- **理由**: 保持错误日志一致性

### 规则51: 函数文档注释格式统一
- **描述**: 公共函数使用统一的文档注释格式，说明函数功能、参数和返回值
- **理由**: 提高代码可读性

### 规则52: 变量命名具有描述性
- **描述**: 变量命名具有描述性，避免使用单字母或无意义的缩写
- **理由**: 提高代码可读性

### 规则53: 避免魔法数字和字符串
- **描述**: 避免在代码中直接使用魔法数字和字符串，使用常量或配置
- **理由**: 提高代码可维护性

### 规则54: 函数长度控制
- **描述**: 函数长度适中，过长的函数应拆分为多个小函数
- **理由**: 提高代码可读性和可维护性

### 规则55: 依赖注入原则
- **描述**: 通过构造函数注入依赖，避免在函数内部直接创建依赖对象
- **理由**: 提高代码可测试性和可维护性

### 规则56: 错误处理统一包装
- **描述**: 错误处理统一使用errors包进行包装，保留错误堆栈信息
- **理由**: 统一错误处理方式

### 规则57: 日志上下文传递
- **描述**: 通过`log.With()`传递日志上下文信息，确保日志关联性
- **理由**: 保持日志上下文一致性

### 规则58: 结构体嵌入用于扩展
- **描述**: 使用结构体嵌入实现代码复用和扩展
- **理由**: 提高代码复用性

### 规则59: 初始化函数集中管理
- **描述**: 初始化相关代码集中在`internal/app/server/server.go`中管理
- **理由**: 统一初始化管理

### 规则60: 服务实现统一注册
- **描述**: 服务实现在`internal/app/service/`目录下统一管理
- **理由**: 统一服务管理

### 规则61: 链式调用格式
- **描述**: 链式调用可以按类型换行，方便阅读
- **理由**: 提高代码可读性

### 规则62: JSON标签规范
- **描述**: 如果结构体不涉及json序列化和反序列化，则不添加json标签
- **理由**: 避免不必要的标签定义