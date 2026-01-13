---
name: go-code-review
description: This skill should be used when the user asks to "review Go code", "check Go code quality", "review this PR", "code review", or mentions Go code standards, GORM best practices, error handling patterns, or design philosophies. Orchestrates comprehensive Go code reviews based on FUTU's coding standards (137+ rules) using parallel specialized agents.
version: 3.0.0
---

# Go Code Review Skill

## When to Use This Skill

This skill activates when users need help with:
- Reviewing Go code changes against coding standards
- Checking code quality and identifying potential issues
- Performing PR reviews for Go projects
- Validating GORM database operations
- Checking error handling patterns
- Verifying naming conventions and code structure
- Analyzing concurrency safety

## Architecture

This plugin uses a **unified orchestration model** with automatic agent selection:

- **Main Orchestrator**: This skill (go-code-review) coordinates the review process
- **4 Specialist Agents** in `agents/` directory:
  - **gorm-review** (blue) - Database operations and GORM best practices (rules 1.3.*)
  - **error-safety** (red) - Error handling and concurrency safety (rules 1.1.*, 1.2.*, 1.4.*, 1.5.*)
  - **naming-logging** (green) - Naming conventions and logging standards (rules 2.1.*, 2.2.*)
  - **organization** (purple) - Code organization, quality, and design philosophies (rules 2.3.*, 2.4.*, 2.5.*, 3.*, 4.*)

**Key Features**:
- **Smart Agent Selection**: Automatically determines which agents apply based on code changes
- **Parallel Execution**: All applicable agents run simultaneously for maximum speed
- **Single Invocation**: One command triggers comprehensive multi-agent review
- **137+ Rules**: Complete FUTU Go coding standards coverage (includes 8 design philosophies)

## Automatic Agent Selection

When invoked, this skill automatically:

1. **Analyzes Changed Files**: Uses `git diff --name-only` to identify modified Go files
2. **Scans Code Content**: Detects patterns that indicate which agents should run:
   - GORM patterns (gorm tags, `db.Model`, `db.Where`) → **gorm-review agent**
   - Error handling (`error` returns, nil checks, `errors.Wrap`) → **error-safety agent**
   - Naming/logging (`func`, `type`, `log.Info`) → **naming-logging agent**
   - Structure (interfaces, function organization) → **organization agent**
3. **Launches Agents in Parallel**: Selected agents run simultaneously
4. **Aggregates Results**: Collects and merges findings by priority (P0/P1/P2)
5. **Delivers Report**: Unified review with Chinese output (per FUTU standards)

**Manual Override**: You can explicitly invoke individual agents:
```
Use gorm-review agent
Use error-safety agent
Use naming-logging agent
Use organization agent
```

## Agent Coordination

**Default Behavior**: Parallel execution for speed
- All 4 agents launch simultaneously
- Each agent checks only its assigned rule categories
- Results merge after all agents complete
- ~4x faster than sequential execution

**Agent Independence**:
- Agents share no state
- Each focuses on specific rule categories
- No conflicts due to clear boundaries
- Can be invoked separately or through orchestrator

## Review Workflow

### Step 1: Get Code Changes
```bash
git diff master
```
Use this command to get all changes between the current branch and master branch.

### Step 2: Systematic Analysis
Analyze all modified files systematically to identify issues based on the priority levels:

**P0 - Must Fix (Critical Issues)**
- Missing or improper error handling
- Unchecked nil pointers
- Improper database operations
- Concurrency safety issues
- Business logic errors

**P1 - Strongly Recommended (Code Quality)**
- Non-standard naming conventions
- Improper logging
- Missing necessary comments
- Poor function encapsulation
- Code duplication

**P2 - Suggested Optimization (Style & Best Practices)**
- Inconsistent formatting
- Readability improvements
- Performance optimization opportunities

### Step 3: Output Format

**IMPORTANT - 重要输出要求**:
- **必须使用中文输出所有审查结果** (All review results MUST be in Chinese)
- Only report code that has problems. Do NOT mention code that follows standards.
- 所有的问题描述、建议都必须用中文表达
- 文件路径、代码片段保持原样

For each issue found, output (用中文):

```markdown
## 文件: <file_path>

### 问题 1 - [严重程度] 规则编号
**位置**: 第 X 行
**原始代码**:
```go
// original code
```
**问题描述**: 用中文详细描述问题
**修改建议**: 用中文说明如何修复
```

### Step 4: Save Results
Output all findings to `code_review.result` file.

## Coding Standards Reference

### P0 Rules - Must Follow

#### 1.1 Error Handling

**1.1.1** Errors must be wrapped and preserve stack trace
- ❌ Forbidden: `fmt.Errorf("operation failed")`
- ✅ Use: `errors.Wrapf(err, "operation failed")` or `errors.WithMessage(err, "operation failed")`

**1.1.2** Error must be in the last return position
- ✅ `func Foo() (result Type, err error)`

**1.1.3** Business functions must return error, not panic
- ❌ Forbidden: `panic("error occurred")`
- ✅ Use: `return errors.New("error occurred")`

**1.1.4** Errors must be checked and handled immediately

**1.1.5** Error message accuracy - Error messages should describe observed facts, not absolute assertions
- ✅ Correct: "instance has no subnet, may be terminating"
- ❌ Wrong: "instance is terminating" (too absolute)
```go
// ✅ Correct: Describe facts, hint possible causes
if instance.SubnetId == nil {
    log.Warn(ctx, "instance has no subnet, may be terminating",
        log.String("instance_id", *instance.InstanceId))
    continue
}
```

**1.1.6** Error message format consistency - Use `failed to <action>` format, include sufficient context
- ✅ Correct: `errors.Wrap(err, "failed to describe aws ec2 instance types")`
- ❌ Wrong: `errors.Wrap(err, "get aws ec2 instance type list failed")`
```go
// ✅ Correct: Unified format
return nil, errors.Wrap(err, "failed to describe aws ec2 instance types")

// ❌ Wrong: Inconsistent format
return nil, errors.Wrap(err, "get aws ec2 instance type list failed")
```

#### 1.2 Nil Pointer Checks

**1.2.1** Must check nil before dereferencing pointers
```go
// ❌ Problem: Dereferencing without nil check
value := *ptr

// ✅ Correct: Check before dereferencing
if ptr != nil {
    value := *ptr
}
```

**重要判断规则**:
- **需要检查**: 当代码对指针进行解引用操作 (使用 `*ptr` 获取值) 时
- **不需要检查**: 当只是将指针值赋给另一个指针类型时 (如 `dest.Field = src.Field`,两者都是指针)
- **不需要检查**: 当传递指针给接受指针类型参数的函数时 (由被调用函数负责处理 nil)

```go
// ✅ 这些情况不需要 nil 检查
type Source struct {
    Data *string
}
type Dest struct {
    Field *string  // 指针类型
}

// Case 1: 指针到指针的赋值
dest.Field = src.Data  // ✅ 不需要检查,两者都是指针类型

// Case 2: 传递给接受指针的函数
func ProcessData(data *string) {
    // 函数内部负责处理 nil
    if data != nil {
        // 使用 data
    }
}
ProcessData(src.Data)  // ✅ 不需要检查,由 ProcessData 处理

// ❌ 这些情况需要 nil 检查
if src.Data != nil {
    value := *src.Data  // 解引用操作
}
```

#### 1.3 Database Operations (GORM)

**1.3.1** Queries must explicitly specify Where conditions
- ✅ `db.Where("id = ?", id).Find(&result)`

**1.3.2** Queries must explicitly specify Select columns (except constant definitions)
- ❌ Forbidden: `SELECT *`
- ✅ Use: `db.Select("id, name, status").Find(&result)`

**1.3.3** Use fluent method chaining for readability and error handling
Chain GORM methods together and assign error at the end for better readability:
```go
// ✅ Recommended: Fluent chaining with error at end
err = db.Model(&model.TagValue{}).
    Select("id, `key`, value, status, sync_status, create_time, update_time, creator, modifier").
    Where(query).
    Offset(int(offset)).
    Limit(int(limit)).
    Find(&tagValueList).Error

// ❌ Avoid: Breaking chain to check intermediate errors (unless necessary)
query := db.Model(&model.TagValue{})
query = query.Select("id, name")
if err := query.Where(condition).Error; err != nil {
    return err
}
```

**1.3.4** Do not use `Save()` method
- ❌ Forbidden: `db.Save(&model)`
- ✅ Use: `db.Updates(&model)` or `db.Create(&model)`

**1.3.5** Prefer `Take()` over `First()`
- `First()` sorts by primary key (worse performance)
- `Take()` doesn't sort (better performance)

**1.3.6** GORM tags must explicitly specify column
```go
type Model struct {
    ID   int64  `gorm:"column:id"`
    Name string `gorm:"column:name"`
}
```

**1.3.7** All struct fields must be used in business logic
Every field defined in a struct MUST be used somewhere in the business flow from creation to completion. During code review, trace the entire business flow to verify all fields are utilized.

**Critical Review Process:**
1. Identify all fields in the struct definition
2. Trace business flow from creation to end
3. Verify each field is read/written at least once
4. Flag any unused fields for removal

```go
// ❌ Problem: UnusedField and Description are never used
type TagKey struct {
    ID          int64  `gorm:"column:id"`
    Key         string `gorm:"column:key"`
    Name        string `gorm:"column:name"`
    UnusedField string `gorm:"column:unused_field"` // ⚠️ Never used in business logic
    Description string `gorm:"column:description"`   // ⚠️ Never read or written
}

func CreateTagKey(key, name string) error {
    tagKey := &TagKey{
        Key:  key,
        Name: name,
        // UnusedField and Description are never set or used
    }
    return db.Create(tagKey).Error
}

func GetTagKey(id int64) (*TagKey, error) {
    var tagKey TagKey
    err := db.Where("id = ?", id).Find(&tagKey).Error
    // Only Key and Name are used later, UnusedField and Description are wasted
    return &tagKey, err
}

// ✅ Solution: Remove unused fields or use them properly
type TagKey struct {
    ID   int64  `gorm:"column:id"`
    Key  string `gorm:"column:key"`
    Name string `gorm:"column:name"`
    // Removed: UnusedField, Description (not needed)
}
```

**When to flag as issue:**
- Field is defined but never set in Create/Update operations
- Field is queried but never used in business logic
- Field exists only for "future use" (violates YAGNI principle)

**Exceptions (fields that don't require usage tracking):**
- Audit fields: `CreateTime`, `UpdateTime`, `Creator`, `Modifier` (system-managed)
- Soft delete: `DeleteTime`, `IsDeleted` (framework-managed)

**1.3.8** GORM struct field ordering for database models
Fields should follow this priority order for readability and consistency:
1. **ID field** - Always first
2. **Business fields** - Ordered by usage frequency (most used first)
3. **Creator/Modifier** - User tracking fields
4. **CreateTime/UpdateTime** - Timestamp fields last

```go
// ✅ Recommended: Proper field ordering
type TagKeyForQuery struct {
    // 1. ID field first
    ID         int64                      `gorm:"column:id"`

    // 2. Business fields by priority (most commonly used first)
    Key        string                     `gorm:"column:key"`
    Name       string                     `gorm:"column:name"`
    Status     constants.TagStatus        `gorm:"column:status"`
    SyncStatus constants.TagKeySyncStatus `gorm:"column:sync_status"`
    Editable   bool                       `gorm:"column:editable"`
    Remark     *string                    `gorm:"column:remark"`

    // 3. User tracking fields
    Creator    string                     `gorm:"column:creator"`
    Modifier   string                     `gorm:"column:modifier"`

    // 4. Timestamp fields last
    CreateTime time.Time                  `gorm:"column:create_time"`
    UpdateTime time.Time                  `gorm:"column:update_time"`
}

// ❌ Avoid: Random field ordering
type TagKeyForQuery struct {
    Name       string    `gorm:"column:name"`
    UpdateTime time.Time `gorm:"column:update_time"`
    ID         int64     `gorm:"column:id"`
    Creator    string    `gorm:"column:creator"`
    // ... fields in random order
}
```

**1.3.9** Serialization strategy - Different layers should use different serialization strategies
- **Storage layer (Mapping)**: Do not use `omitempty`, preserve all field information
- **Business layer (Detail)**: Use `omitempty` for optional fields to reduce serialization size
- **API layer**: Decide based on API specification
```go
// ✅ Correct: Different strategies for different layers
// Storage layer: Preserve all information
type AliyunECSMapping struct {
    SystemDiskSize *int64 `json:"system_disk_size"`  // No omitempty
}

// Business layer: Use omitempty for optional fields
type AliyunECSDetail struct {
    AdditionalFields *AdditionalFields `json:"additional_fields,omitempty"`
}
```

#### 1.4 Concurrency Control

**1.4.1** Use `recovered.ErrorGroup` for concurrent operations
```go
errGroup := recovered.NewErrorGroup()
```

**1.4.2** Use limiter for concurrency control
```go
l := limiter.NewConcurrentLimiter(count)
```

#### 1.5 JSON Processing

**1.5.1** Do not construct JSON via string concatenation or formatting
- ❌ Forbidden: `fmt.Sprintf(`{"name": "%s"}`, name)`
- ✅ Use: `json.Marshal(data)`

### P1 Rules - Strongly Recommended

#### 2.1 Naming Conventions

**2.1.1** Structs: CamelCase, first letter uppercase
- ✅ `type UserAccount struct {}`

**2.1.2** Struct fields: CamelCase matching tags
```go
type Model struct {
    UserName  string `gorm:"column:user_name"`
    ProductID int64  `gorm:"column:product_id"`
}
```

**2.1.3** Functions/Methods: CamelCase
- Public: First letter uppercase `GetUserByID()`
- Private: First letter lowercase `getUserByID()`

**2.1.4** Constants: ALL_CAPS with underscores
```go
const (
    MAX_RETRY_COUNT = 3
    DEFAULT_TIMEOUT = 30
)
```

**2.1.5** ID abbreviation always uppercase
- ✅ `UserID`, `ProductID`
- ❌ `UserId`, `ProductId`

**2.1.6** Avoid magic numbers and strings - define constants

**2.1.7** Avoid multiple naming styles for same concept
- ✅ Consistent: `CloudTencent`, `CloudAliyun`, `CloudAws`
- ❌ Mixed: `TencentCloud`, `AliyunCloud`, `CloudAws`

**2.1.8** Business field naming should be unambiguous
```go
type Model struct {
    ID         int64  `gorm:"column:id"`          // table primary key
    TemplateID string `gorm:"column:template_id"` // cloud provider ID
    TemplateTid int64 `gorm:"column:template_tid"` // related table ID
}
```

**2.1.9** Variable names should be descriptive
- ❌ Avoid: `a`, `tmp`, `x`
- ✅ Use: `userCount`, `templateList`, `configMap`

**2.1.10** Interface names usually end with -er
- ✅ `Reader`, `Writer`, `Processor`

**2.1.11** Opposite operations must correspond strictly
- ✅ `Enable` ↔ `Disable`
- ✅ `Start` ↔ `Stop`
- ✅ `Open` ↔ `Close`

**2.1.12** Method names for getting lists use GetList or Search
- ❌ `Get()` - implies single item
- ✅ `GetList()`, `Search()` - implies multiple items

**2.1.13** Self-documenting names - Function names should clearly describe functionality and data source
- ✅ `ExtractResourceNameFromTags()` - clearly indicates extracting from tags
- ❌ `embedAwsEC2Name2ResourceName()` - unclear what "embed" and "2" mean
- Variable names should express type and purpose (e.g., `privateIPList` for list types)
- Similar field names should be clearly distinguished (e.g., `PublicIPAddress` vs `PublicIPList`)
```go
// ✅ Correct: Self-documenting
func ExtractResourceNameFromTags(tags []types.Tag) string {
    // Function name already explains: extract resource name from tags
}

// ❌ Wrong: Need to read implementation to understand
func embedAwsEC2Name2ResourceName(tags []types.Tag) string {
    // Name unclear: what does "embed" mean? what does "2" mean?
}
```

**2.1.14** Naming consistency - Use standard terminology abbreviations, maintain consistent naming style
- Use standard abbreviations: `ARN` not `Arn`
- Same type entities should follow same naming style
- Go field names and JSON tags should follow consistent conversion rules
```go
// ✅ Correct: Standard terminology
RoleARN string `json:"role_arn"`  // ARN is standard abbreviation

// ❌ Wrong: Inconsistent terminology
RoleArn string `json:"role_arn"`  // Should use ARN
```

#### 2.2 Logging Standards

**2.2.1** Log fields use snake_case
- ❌ `log.String("UserName", name)`
- ✅ `log.String("user_name", name)`

**2.2.2** Log fields use explicit types (except arrays, structs, maps use log.Any)
- ❌ `log.Any("count", len(list))`
- ✅ `log.Int("count", len(list))`

**2.2.3** Logs must include key context
```go
log.Info(ctx, "operation success",
    log.String("user_id", userID),
    log.Int64("resource_id", resourceID))
```

**2.2.4** Functions should have necessary info/debug logs (non-trivial functions)

**2.2.5** Data layer must not log - logging should be in business logic layer

**2.2.6** Error logs must include ErrorField at the end
```go
log.Error(ctx, "operation failed",
    log.String("user_id", userID),
    log.ErrorField(err))  // error field at the end
```

**2.2.7** Layered logging strategy - avoid redundant logs
遵循分层日志记录原则,避免重复日志:

**外层函数(Service/Handler)**: 记录入参和最终结果
- 函数入口: 记录关键输入参数 (非敏感)
- 函数出口: 记录执行结果和关键输出
- 错误处理: 记录错误日志

**内层函数(Helper/Utility)**: 只记录外部调用结果
- ✅ 记录: 调用云 API、第三方服务、数据库的返回值
- ❌ 不记录: 内部逻辑处理过程
- ❌ 不记录: 错误日志 (使用 error wrap 传递上下文)

**错误传递**: 使用 error wrap,不使用日志
- ✅ 使用: `errors.Wrapf(err, "context info")`
- ❌ 避免: `log.Error() + return err` (会导致重复日志)

```go
// ✅ 正确示例: 分层日志
// 外层函数 - Service 层
func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    // 入口日志: 记录关键入参
    log.Info(ctx, "creating user",
        log.String("user_name", req.Name),
        log.String("email", req.Email))

    // 调用内层函数
    user, err := s.createUserInternal(ctx, req)
    if err != nil {
        // 外层记录错误日志
        log.Error(ctx, "create user failed",
            log.String("user_name", req.Name),
            log.ErrorField(err))
        return nil, errors.Wrap(err, "create user failed")
    }

    // 出口日志: 记录结果
    log.Info(ctx, "user created successfully",
        log.Int64("user_id", user.ID))

    return user, nil
}

// 内层函数 - 不记录日志,只 wrap 错误
func (s *UserService) createUserInternal(ctx context.Context, req *CreateUserRequest) (*User, error) {
    // ❌ 不记录入口日志 (外层已记录)

    // 验证逻辑
    if err := s.validateUser(req); err != nil {
        // ❌ 不记录错误日志,只 wrap
        return nil, errors.Wrap(err, "validate user failed")
    }

    // 调用数据库
    user, err := s.userRepo.Create(ctx, req)
    if err != nil {
        // ❌ 不记录错误日志,只 wrap
        return nil, errors.Wrap(err, "insert user to db failed")
    }

    // ✅ 记录外部调用结果 (如果有)
    if req.SendWelcomeEmail {
        err := s.emailClient.SendWelcome(ctx, user.Email)
        if err != nil {
            // 记录外部调用失败 (非致命)
            log.Warn(ctx, "send welcome email failed",
                log.String("email", user.Email),
                log.ErrorField(err))
        } else {
            log.Info(ctx, "welcome email sent",
                log.String("email", user.Email))
        }
    }

    // ❌ 不记录出口日志 (外层会记录)
    return user, nil
}

// ❌ 错误示例: 多层重复日志
func (s *UserService) CreateUserBad(ctx context.Context, req *CreateUserRequest) (*User, error) {
    log.Info(ctx, "creating user") // 外层记录

    user, err := s.createUserInternalBad(ctx, req)
    if err != nil {
        log.Error(ctx, "create user failed", log.ErrorField(err)) // 外层记录错误
        return nil, err
    }

    log.Info(ctx, "user created") // 外层记录
    return user, nil
}

func (s *UserService) createUserInternalBad(ctx context.Context, req *CreateUserRequest) (*User, error) {
    log.Info(ctx, "internal create user") // ❌ 重复: 内层也记录入口

    if err := s.validateUser(req); err != nil {
        log.Error(ctx, "validate failed", log.ErrorField(err)) // ❌ 重复: 内层记录错误
        return nil, errors.Wrap(err, "validate failed")
    }

    user, err := s.userRepo.Create(ctx, req)
    if err != nil {
        log.Error(ctx, "db insert failed", log.ErrorField(err)) // ❌ 重复: 内层记录错误
        return nil, errors.Wrap(err, "db insert failed")
    }

    log.Info(ctx, "internal user created") // ❌ 重复: 内层也记录出口
    return user, nil
}
// 结果: 一次操作产生 6 条日志,信息高度重复
```

**关键原则**:
1. **单一职责**: 每条信息只在一个层次记录
2. **外层负责**: Service/Handler 层记录完整的请求-响应链路
3. **内层简化**: Helper/Utility 层只 wrap 错误,不记录日志
4. **外部调用例外**: 调用云 API/第三方服务必须记录结果 (便于排查外部问题)
5. **错误链传递**: 使用 `errors.Wrap` 构建错误上下文,最终在外层统一记录

#### 2.3 Code Organization

**2.3.1** Embedded structs must be at the top
```go
type Model struct {
    orm.BaseModel  // embedded at top
    Name   string
    Status int
}
```

**2.3.2** JSON tags should relate to field names

**2.3.3** Use make to initialize map/slice (except unmarshal/scan)
```go
slice := make([]string, 0, capacity)
m := make(map[string]int, capacity)
```

**2.3.4** Specify capacity when size is known
```go
result := make([]Item, 0, len(sourceList))
```

**2.3.5** Switch must have default branch
```go
switch status {
case StatusActive:
    return "active"
default:
    return "unknown"
}
```

**2.3.6** Don't use iota - specify values explicitly
```go
const (
    StatusEnable  = 0
    StatusDisable = 1
)
```

**2.3.7** No global variables (except unit tests, pre-compiled regex)

**2.3.8** No init functions (except server package service registration)

**2.3.9** Function encapsulation must be meaningful

**2.3.10** NOTE/TODO must include owner
```go
// TODO(zhangsan): optimize performance
// NOTE(lisi): temporary solution, refactor later
```

**2.3.11** Chain calls should break lines for readability
```go
db.Select("id, name").
    Where("status = ?", 1).
    Order("created_at DESC").
    Find(&result)
```

**2.3.12** Separation of concerns - Different business concepts should use independent structs
- Use naming conventions (e.g., `Parent_Child`) to organize related types
- Embedded structs should be at top level, clearly expressing composition
```go
// ✅ Correct: Separation of concerns
type AliyunECSDetail struct {
    *BaseMapping  // At top, clearly embedded
    SystemDisk *SystemDisk
    DataDisks  []*DataDisk
}

// ❌ Wrong: Mixed concerns
type Disk struct {
    IsSystem bool  // Distinguish by flag, error-prone
}
```

**2.3.13** Reduce cognitive load - Related fields should be grouped together
- **Field grouping**: Related fields together (compute resources: CPU, memory; storage: system disk, data disks; network: VPC, subnet, IP)
- **Logical grouping**: Use comments for grouping, but don't over-comment
- **Order consistency**: Same type structs should maintain consistent field ordering
```go
// ✅ Correct: Clear logical grouping
type ResourceDetail struct {
    // Compute resources
    CPU    int32 `json:"cpu"`
    Memory int32 `json:"memory"`
    
    // Storage resources
    SystemDisk *SystemDisk `json:"system_disk"`
    DataDisks  []*DataDisk `json:"data_disks"`
}

// ❌ Wrong: Chaotic field ordering
type ResourceDetail struct {
    CPU        int32
    SystemDisk *SystemDisk
    Memory     int32  // Memory separated
    DataDisks  []*DataDisk
}
```

**2.3.14** Avoid duplication (DRY) - Common information should not be repeated
- **Base fields**: Common fields should be in base structs, subclasses should not repeat
- **Constant extraction**: Magic values should be extracted as constants in common packages
- **Utility functions**: Repeated conversion logic should be extracted as utility functions
```go
// ✅ Correct: Avoid duplication
type BaseMapping struct {
    ID     string `json:"id"`
    Name   string `json:"name"`
    Region string `json:"region"`
}

type AliyunECSMapping struct {
    *BaseMapping  // Reuse base fields
    // Specific fields
}

// ❌ Wrong: Repeated definition
type AliyunECSMapping struct {
    ID     string `json:"id"`     // Duplicate
    Name   string `json:"name"`   // Duplicate
    Region string `json:"region"` // Duplicate
}
```

#### 2.4 Interface Design

**2.4.1** Interfaces follow minimal principle

**2.4.2** Repeated parameters can be defined in constructor

**2.4.3** ctx parameter required if using ctx

**2.4.4** ctx as first parameter
```go
func Process(ctx context.Context, data []byte) error
```

**2.4.5** Parameter naming matches operation
```go
Create(ctx context.Context, data *CreateParam) (int64, error)  // create doesn't need ID
Update(ctx context.Context, id int64, data *UpdateParam) error // update needs ID
```

**2.4.6** More than 4 parameters - use struct
```go
type CreateUserParams struct {
    Name  string
    Email string
    Phone string
    Role  string
}
func CreateUser(ctx context.Context, params *CreateUserParams) error
```

**2.4.7** Required parameters must consider zero value validity
- If zero value is valid, use pointer type to distinguish "not provided" vs "zero value provided"
- If zero value is invalid, use value type
```go
// ✅ Zero value valid - use pointer
type UpdateUserRequest struct {
    ID     int64   `json:"id" validate:"required"`
    Status *int    `json:"status" validate:"required"` // 0 is valid
    Age    *int    `json:"age" validate:"required"`    // 0 years may be valid
}

// ✅ Zero value invalid - use value type
type CreateUserRequest struct {
    Name   string `json:"name" validate:"required"`    // empty string invalid
    Email  string `json:"email" validate:"required"`   // empty string invalid
    RoleID int64  `json:"role_id" validate:"required"` // 0 invalid (by business rule)
}
```

#### 2.5 Code Quality

**2.5.1** No spelling errors (except whitelist)

**2.5.2** Avoid magic numbers and strings

**2.5.3** Function length should be moderate - split if too long

**2.5.4** Public functions must have comments
```go
// GetUserByID retrieves user by ID
// Parameters:
//   ctx: context
//   id: user ID
// Returns:
//   *User: user info
//   error: error if any
func GetUserByID(ctx context.Context, id int64) (*User, error)
```

**2.5.5** Enum constants must have Chinese comments and be aligned
```go
const (
    StatusInit     = 0 // 初始化
    StatusPending  = 1 // 待处理
    StatusApproved = 2 // 已通过
)
```

**2.5.6** Inject dependencies via constructor

**2.5.7** Use struct embedding for code reuse

**2.5.8** Allow pointer passing when data is already pointer type
```go
// ✅ Data is already pointer, pass directly
type User struct {
    Profile *UserProfile
}

func ProcessUser(user *User) {
    SaveProfile(user.Profile) // ✅ OK to pass nil if function accepts it
}
```

**2.5.9** Types should accurately express business semantics
- **Nullable modeling**: Use pointer types `*T` for nullable fields, value types `T` for required fields
- **Type distinction**: Semantically different data should use different types
- **Judgment basis**: Based on business logic and API documentation, not technical convenience
```go
// ✅ Correct: Types clearly express semantics
type SystemDisk struct {
    ID   *string `json:"id"`   // May be null
    Size int     `json:"size"` // Must have value
}

type DataDisk struct {
    ID   *string `json:"id"`
    Size int     `json:"size"`
}

// ❌ Wrong: Types cannot distinguish semantics
type Disk struct {
    Type string  // Distinguish by field, error-prone
    ID   *string
}
```

**2.5.10** Avoid type abuse - Don't sacrifice type safety for technical convenience
- Different types of data should not be in the same container (e.g., array)
- Use structs instead of arrays to express composite data
- Avoid using `interface{}` or `map[string]any` to escape type checking
```go
// ✅ Correct: Use struct
type InstanceSpec struct {
    Memory int64 `json:"memory"`
    Cores  int64 `json:"cores"`
}

// ❌ Wrong: Use array, type unsafe
spec := []*int64{memory, &cores}  // Cannot distinguish which is memory, which is cores
```

**2.5.11** Unit consistency - Same type of measurements should use unified units
- **Unit standardization**: Same type measurements use same unit (memory: GB, time: seconds or milliseconds)
- **Transparent conversion**: Unit conversion should be done at data entry point, internally use standard units
- **Document clearly**: Field comments must clearly state units
```go
// ✅ Correct: Unified units
type InstanceSpec struct {
    Memory float32 `json:"memory"` // Unit: GB
}

// Conversion done at entry point
memory := float32(apiResponse.MemoryInMB) / 1024.0  // MB -> GB

// ❌ Wrong: Unit confusion
type InstanceSpec struct {
    Memory int32 `json:"memory"`  // Unit unclear, may be MB or GB
}
```

**2.5.12** Data source abstraction - Differences between data sources should be handled in adapter layer
- **Unified interface**: API differences between cloud providers should be handled in adapter layer
- **Boundary handling**: Boundary conditions (empty string vs nil) should be unified in adapter layer
- **Error handling**: Different provider error formats should be unified in adapter layer
```go
// ✅ Correct: Unified in adapter layer
func (a *AliyunAdapter) GetNextToken(resp *APIResponse) *string {
    if resp.NextToken == nil || *resp.NextToken == "" {
        return nil  // Unified return nil
    }
    return resp.NextToken
}

// Business layer uses unified interface
token := adapter.GetNextToken(resp)
if token == nil {
    break
}
```

**2.5.13** Minimalism principle - Only keep necessary code, remove redundant and unused parts
- **Remove unused code**: Regularly clean unused structs, functions, constants, imports
- **Remove redundant fields**: Only keep fields truly needed by business
- **Simplify structure**: Nested structures without clear purpose should be simplified

**Checklist**:
- [ ] Unused struct definitions
- [ ] Unused functions
- [ ] Unused constants
- [ ] Unused import packages
- [ ] Duplicate field definitions
- [ ] Unnecessary nested structures

**2.5.14** Defensive programming - Code should gracefully handle boundary conditions and exceptions
- **Boundary checks**: Loops and conditionals must consider boundary cases (initial values, nil, empty strings)
- **Error descriptions**: Error messages should describe actual situation, not absolute assertions
  - ✅ "instance has no subnet, may be terminating"
  - ❌ "instance is terminating" (too absolute)
- **Testability**: New logic must be verified through actual execution
```go
// ✅ Correct: Defensive checks
var nextToken *string = nil
for {
    // Processing logic
    if nextToken == nil || *nextToken == "" {
        break  // Handle nil and empty string
    }
}

// ❌ Wrong: Missing boundary check
for nextToken != nil {  // If initial value is nil, never executes
    // ...
}
```

**2.5.15** Verifiability - Code correctness should be verifiable through testing
- **New logic must be tested**: New loops, conditionals must be verified through actual execution
- **Boundary testing**: Must test boundary conditions (nil, empty strings, empty arrays)
- **Integration testing**: Logic involving external APIs should have integration tests

**2.5.16** Appropriate abstraction level - Code should operate at appropriate abstraction level
- **Utility function encapsulation**: Repeated conversion logic should be encapsulated as utility functions
- **Context awareness**: Tool function usage should consider context (e.g., conversions for different data sources)
- **Abstraction consistency**: Operations at same abstraction level should use same patterns
```go
// ✅ Correct: Use utility functions
privateIPList := tea.StringSliceValue(instance.VpcAttributes.PrivateIpAddress.IpAddress)

// Or custom conversion function
func ConvertIPAddresses(ips []*string) []string {
    result := make([]string, 0, len(ips))
    for _, ip := range ips {
        if ip != nil {
            result = append(result, *ip)
        }
    }
    return result
}
```

**2.5.17** Performance optimization principle - Reasonable performance optimization while ensuring readability
- **Pre-allocate capacity**: If slice capacity is known, should pre-allocate
- **Avoid unnecessary allocation**: Avoid repeated allocation in loops
- **Reuse utility functions**: Use optimized utility functions
```go
// ✅ Correct: Pre-allocate capacity
list := make([]string, 0, len(source))

// ❌ Wrong: No pre-allocation
list := make([]string, 0)  // May cause multiple reallocations
```

**2.5.18** Field semantics clarity - Field meanings and sources should be clear and explicit
- **Field comments**: Similar field names should be clearly distinguished through comments
- **Naming distinction**: Express different meanings through naming
- **Data source**: Field data sources should be reflected in comments or naming
```go
// ✅ Correct: Clear semantics
// PublicIPAddress is the original public IP field returned by AWS API
PublicIPAddress string `json:"public_ip_address"`

// PublicIPList is the public IP list collected from all network interfaces
PublicIPList []string `json:"public_ip_list"`
```

### P2 Rules - Suggested

#### 3.1 Project Structure

**3.1.1** Package name matches directory name

**3.1.2** Import grouping and sorting
```go
import (
    // Standard library
    "context"
    "fmt"

    // Third-party
    "github.com/pkg/errors"

    // Internal
    "project/internal/model"
)
```

**3.1.3** Business logic in `biz/`, data access in `data/`

**3.1.4** Cloud provider code organized by provider

**3.1.5** Configuration code in `config/`

**3.1.6** Utility code in `pkg/`

**3.1.7** Initialization in `server/server.go`

**3.1.8** Service implementation in `service/`

#### 3.2 Testing Standards

**3.2.1** Unit tests must verify specific structure, not just err != nil
```go
assert.NoError(t, err)
assert.Equal(t, expectedValue, actualValue)
```

**3.2.2** Test file naming: `xxx_test.go`

**3.2.3** Test function naming: `TestXxx`

#### 3.3 Configuration Management

**3.3.1** Configuration updates should take effect dynamically
```go
func (c *Config) OnConfigChange() error {
    // handle config update
}
```

## Review Checklist

When performing code review, check:

### Security
- [ ] SQL injection risks
- [ ] Nil pointer risks
- [ ] Concurrency safety issues
- [ ] Proper handling of sensitive information

### Correctness
- [ ] Business logic is correct
- [ ] Error handling is complete
- [ ] Edge cases are considered
- [ ] Data consistency is guaranteed

### Performance
- [ ] N+1 query issues
- [ ] Repeated calculations in loops
- [ ] Capacity pre-allocation
- [ ] Unnecessary memory allocations

### Maintainability
- [ ] Clear naming
- [ ] Adequate comments
- [ ] No code duplication
- [ ] Reasonable function length

### Standards Compliance
- [ ] Follows project coding standards
- [ ] Correct log format
- [ ] Adequate test coverage

## Example Output (示例输出 - 使用中文)

```markdown
# 代码审查结果

## 文件: internal/service/user_service.go

### 问题 1 - [P0] 规则 1.1.1 错误包装
**位置**: 第 45 行
**原始代码**:
```go
return fmt.Errorf("get user failed: %v", err)
```
**问题描述**: 使用 fmt.Errorf 而不是 errors.Wrap,会丢失错误堆栈信息
**修改建议**:
```go
return errors.Wrapf(err, "get user failed")
```

### 问题 2 - [P1] 规则 2.2.1 日志字段命名
**位置**: 第 67 行
**原始代码**:
```go
log.Info(ctx, "user created", log.String("UserID", userID))
```
**问题描述**: 日志字段使用了 PascalCase,应该使用 snake_case
**修改建议**:
```go
log.Info(ctx, "user created", log.String("user_id", userID))
```
```

## Important Notes

- **所有输出必须使用中文** (All output MUST be in Chinese)
- **Only report issues** - Don't mention code that already follows standards
- **Be specific** - Include exact line numbers from source files
- **Prioritize by severity** - P0 issues first, then P1, then P2
- **Provide clear suggestions** - Show exactly how to fix each issue (用中文说明)
- **Focus on patterns** - Look for systematic issues across files
- **Consider context** - Business requirements may justify exceptions

## Best Practices

1. **Systematic Review**: Check all modified files thoroughly
2. **Priority Focus**: Address P0 issues before P1/P2
3. **Clear Communication**: Use exact file paths and line numbers
4. **Actionable Feedback**: Provide specific code suggestions
5. **Standards Adherence**: Reference specific rule numbers
6. **Save Results**: Always output to `code_review.result` file
