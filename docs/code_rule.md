# Go代码规范规则

## 1. 错误处理规范

### 规则1: 错误包装必须使用errors.Wrap或errors.WithMessage
- **描述**: 任何错误返回必须使用errors.Wrap或errors.WithMessage进行错误包装，不允许直接返回fmt.Errorf
- **证据**: 文件位置 internal/app/acl/aliyun/role_manager.go:21, internal/app/acl/aliyun/role_manager.go:40
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
- **证据**: 文件位置 internal/app/biz/role.go:60, internal/app/biz/role.go:130
- **反面示例**:
```go
log.Error(ctx, "failed to get account config", log.String("uin", account.Uin), log.ErrorField(err))
```
- **正面示例**:
```go
log.Error(ctx, "failed to get account config", log.String("uin", account.Uin), log.ErrorField(err))
```
- **理由**: 保持错误日志的一致性，便于错误追踪和分析

### 规则3: 错误必须放在返回值的最后
- **描述**: 函数返回值中错误必须在最后位置
- **证据**: 文件位置 internal/app/biz/depart_permission.go:63, internal/app/biz/product_permission.go:82
- **反面示例**:
```go
func ChangeToDepartPermissionType(val pb.PermissionType) error {...}
```
- **正面示例**:
```go
func ChangeToDepartPermissionType(val pb.PermissionType) (result type, err error) {...}
```
- **理由**: 符合Go的错误处理标准，便于统一处理错误

## 2. 日志规范

### 规则4: 日志字段使用snake_case命名
- **描述**: 日志字段的key必须使用snake_case命名风格
- **证据**: 文件位置 internal/app/biz/depart_permission.go:124, internal/app/data/db/account.go:36
- **反面示例**:
```go
log.String("UIN", account.Uin)
```
- **正面示例**:
```go
log.String("uin", account.Uin)
```
- **理由**: 保持日志字段命名的一致性，便于日志查询和分析

### 规则5: 避免在data层打日志
- **描述**: data层（数据库操作层）不应该打印日志
- **证据**: 文件位置 internal/app/data/db/role_depart.go:59
- **反面示例**:
```go
log.Info(ctx, "some log in data layer")
```
- **正面示例**:
```go
// data层不打印日志，由上层业务逻辑处理日志
```
- **理由**: data层应该专注于数据操作，日志属于业务逻辑层面，应由上层处理

### 规则6: 日志输出必要的上下文信息
- **描述**: 日志必须包含足够的上下文信息，便于问题定位
- **证据**: 文件位置 internal/app/biz/role_depart.go:54
- **反面示例**:
```go
log.Info(ctx, "some operation")
```
- **正面示例**:
```go
log.Info(ctx, "some operation", log.String("alias", alias), log.Int64("id", id))
```
- **理由**: 充足的上下文信息有助于快速定位问题

## 3. 数据库操作规范

### 规则7: GORM查询必须明确指定查询条件
- **描述**: 使用GORM查询时必须明确使用Where条件
- **证据**: 文件位置 internal/app/data/db/apply.go:31, internal/app/data/db/apply.go:32
- **反面示例**:
```go
db.Model(&model.Role{}).Take(&role).Error
```
- **正面示例**:
```go
db.Where("id = ?", id).Model(&model.Role{}).Take(&role).Error
```
- **理由**: 避免查询条件不明确导致的数据错误

### 规则8: GORM查询必须明确指定查询列
- **描述**: 不要使用SELECT *，必须明确指定查询的列
- **证据**: 文件位置 internal/app/data/db/permission_group.go:76, internal/app/data/db/apply.go:57
- **反面示例**:
```go
db.Find(&pg).Error
```
- **正面示例**:
```go
db.Select("group_code, permission_list").Where(...).Find(&pg).Error
```
- **理由**: 提高查询效率，减少不必要的数据传输

### 规则9: 避免使用GORM的save方法
- **描述**: 项目内不允许使用GORM的save方法，应使用update等其他方法
- **证据**: 文件位置 internal/app/data/db/apply.go:71
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

### 规则10: 指针变量必须进行空指针校验
- **描述**: 使用指针变量前必须进行空指针校验
- **证据**: 文件位置 internal/app/biz/role_depart.go:52, internal/app/biz/role_product.go:36
- **反面示例**:
```go
result := *somePtr  // 没有校验是否为空
```
- **正面示例**:
```go
if somePtr != nil {
    result := *somePtr
}
```
- **理由**: 防止空指针引用导致程序panic

### 规则11: ID统一写成ID而不是Id
- **描述**: 除了生成代码外，ID统一写成ID，而不是Id
- **证据**: 文件位置 internal/app/biz/repo/role_depart.go:10
- **反面示例**:
```go
Id string
```
- **正面示例**:
```go
ID string
```
- **理由**: 保持ID命名的一致性

### 规则12: 接口命名规范
- **描述**: 接口命名应保持简单准确，以-er后缀命名
- **证据**: 文件位置 internal/app/biz/repo/account.go:14
- **反面示例**:
```go
type MainAccountRepo interface {
    AddFavourite(ctx context.Context, alias, uin string) error  // 注释与方法名不匹配
}
```
- **正面示例**:
```go
type MainAccountRepo interface {
    AddFavourite(ctx context.Context, alias, uin string) error  // 名称和注释应匹配
}
```
- **理由**: 接口命名应清晰准确，便于理解

## 5. 代码结构和组织规范

### 规则13: 常量定义规范
- **描述**: 字符串、数字等如果有统一管理、复用等需求，统一定义在常量包内
- **证据**: 文件位置 internal/app/acl/aliyun/role_manager.go:55, internal/app/acl/aws/role_manager.go:67
- **反面示例**:
```go
if provider == "aws" {  // 使用硬编码字符串
```
- **正面示例**:
```go
const (
    ProviderAWS = "aws"
)
if provider == ProviderAWS {
```
- **理由**: 避免魔法字符串，提高代码可维护性

### 规则14: 结构体嵌入规范
- **描述**: 嵌入的结构体永远在最顶层
- **证据**: 文件位置 internal/app/model/db/account.go:11
- **反面示例**:
```go
type AccountFavor struct {
    Alias string `gorm:"column:alias;type:varchar(128);not null;comment:别名"`
    orm.BaseWithCTime
}
```
- **正面示例**:
```go
type AccountFavor struct {
    orm.BaseWithCTime
    Alias string `gorm:"column:alias;type:varchar(128);not null;comment:别名"`
}
```
- **理由**: 保持结构体字段组织的一致性

### 规则15: switch/case必须带default分支
- **描述**: 原则上所有的switch/case都必须带上default分支
- **证据**: 文件位置 internal/app/constants/config.go:43
- **反面示例**:
```go
switch val {
case 1:
    // do something
case 2:
    // do something
}
```
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

## 6. 并发和错误处理规范

### 规则16: 使用common包的recovered.ErrorGroup
- **描述**: 并发操作应该使用common包中的recovered.ErrorGroup
- **证据**: 文件位置 internal/app/biz/depart_permission.go:157, internal/app/biz/policy.go:291
- **反面示例**:
```go
errGroup := &sync.WaitGroup{}
```
- **正面示例**:
```go
errGroup := recovered.NewErrorGroup()
```
- **理由**: recovered.ErrorGroup提供更好的错误处理和取消机制

### 规则17: 并发控制使用limiter
- **描述**: 控制并发数量应使用limiter.NewConcurrentLimiter
- **证据**: 文件位置 internal/app/biz/role.go:77
- **反面示例**:
```go
// 没有限制并发数
```
- **正面示例**:
```go
l := limiter.NewConcurrentLimiter(limiterCount)
```
- **理由**: 防止并发过多导致资源耗尽

## 7. 配置管理规范

### 规则18: 配置更新处理
- **描述**: 配置中心更新后，相关配置应能够正确更新
- **证据**: 文件位置 internal/app/server/server.go:55
- **反面示例**:
```go
// 配置初始化后不再更新
```
- **正面示例**:
```go
// 通过回调函数处理配置更新
func (p *ProxyConfig) OnConfigChange() (err error) {
    // 处理配置更新
}
```
- **理由**: 保证配置能够动态更新，无需重启服务

## 8. 测试规范

### 规则19: 单元测试必须验证具体结构
- **描述**: 单元测试必须验证具体的结构，不能只验证err是否为nil
- **证据**: 文件位置 internal/app/service/cloud_mgr_server_test.go:191
- **反面示例**:
```go
assert.NoError(t, err)  // 只验证错误
```
- **正面示例**:
```go
assert.NoError(t, err)
assert.Equal(t, expectedValue, actualValue)  // 验证具体值
```
- **理由**: 验证具体结构确保功能正确性

## 9. 其他规范

### 规则20: 禁止使用全局变量
- **描述**: 无特殊场景（例如单测、正则预编译），禁止使用全局变量
- **证据**: 文件位置 internal/app/config/workflow.go:19, internal/pkg/workflow/client.go:25
- **反面示例**:
```go
var globalVar = "something"
```
- **正面示例**:
```go
// 通过依赖注入传递
```
- **理由**: 避免全局状态导致的并发问题和测试困难

### 规则21: 枚举类型定义
- **描述**: 一般情况下不使用iota作为自增常量，应明确指定值
- **证据**: 文件位置 internal/app/acl/tencent/polic_manager.go:25
- **反面示例**:
```go
type Status int
const (
    StatusEnable Status = iota  // 使用iota自增
    StatusDisable
)
```
- **正面示例**:
```go
type Status int
const (
    StatusEnable Status = 0
    StatusDisable Status = 1
)
```
- **理由**: 明确指定值可以防止枚举顺序变化导致的错误

### 规则22: JSON标签规范
- **描述**: 如果结构体不涉及json序列化和反序列化，则不添加json标签
- **证据**: 文件位置 internal/app/config/aws_actions.go:12
- **反面示例**:
```go
type Config struct {
    data *Proxy `json:"data"`  // json标签但字段不可导出
}
```
- **正面示例**:
```go
type Config struct {
    data *Proxy  // 不需要json标签
}
```
- **理由**: 避免不必要的标签定义

### 规则23: 切片初始化规范
- **描述**: 初始化map和切片时，除了unmarshal、scan等场景，必须通过make进行初始化
- **证据**: 文件位置 internal/app/constants/config.go:24
- **反面示例**:
```go
var slice []string
```
- **正面示例**:
```go
slice := make([]string, 0)  // 或 make([]string, 0, capacity)
```
- **理由**: 使用make初始化更明确，且可以预设容量优化性能

### 规则24: 链式调用格式
- **描述**: 链式调用可以按类型换行，方便阅读
- **证据**: 文件位置 internal/app/data/db/permission_group.go:45
- **反面示例**:
```go
return db.Model(&model.PermissionGroup{}).Where("id = ? ", id).Updates(map[string]interface{}{"permission_list": PermissionList}).Error
```
- **正面示例**:
```go
return db.Model(&model.PermissionGroup{}).
    Where("id = ? ", id).
    Updates(map[string]interface{}{"permission_list": PermissionList}).
    Error
```
- **理由**: 提高代码可读性