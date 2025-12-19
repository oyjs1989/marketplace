---
name: Go Code Review
description: Performs comprehensive Go code reviews based on FUTU's coding standards. ALWAYS use this when user asks to "review Go code", "check Go code quality", "review this PR", "code review", or mentions Go code standards, GORM best practices, or error handling patterns.
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

**IMPORTANT**: Only report code that has problems. Do NOT mention code that follows standards.

For each issue found, output:

```markdown
## File: <file_path>

### Issue 1 - [Severity] Rule Number
**Location**: Line X
**Original Code**:
```go
// original code
```
**Problem**: Detailed problem description
**Suggestion**: How to fix it
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

#### 1.2 Nil Pointer Checks

**1.2.1** Must check nil before using pointers
```go
if ptr != nil {
    value := *ptr
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

**2.2.7** Method entry and exit must log
- Entry: log input parameters (non-sensitive)
- Exit: log key execution information
- Important method calls: log input and return values

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

## Example Output

```markdown
# Code Review Result

## File: internal/service/user_service.go

### Issue 1 - [P0] Rule 1.1.1 Error Wrapping
**Location**: Line 45
**Original Code**:
```go
return fmt.Errorf("get user failed: %v", err)
```
**Problem**: Using fmt.Errorf instead of errors.Wrap, losing error stack information
**Suggestion**:
```go
return errors.Wrapf(err, "get user failed")
```

### Issue 2 - [P1] Rule 2.2.1 Log Field Naming
**Location**: Line 67
**Original Code**:
```go
log.Info(ctx, "user created", log.String("UserID", userID))
```
**Problem**: Log field uses PascalCase, should use snake_case
**Suggestion**:
```go
log.Info(ctx, "user created", log.String("user_id", userID))
```
```

## Important Notes

- **Only report issues** - Don't mention code that already follows standards
- **Be specific** - Include exact line numbers from source files
- **Prioritize by severity** - P0 issues first, then P1, then P2
- **Provide clear suggestions** - Show exactly how to fix each issue
- **Focus on patterns** - Look for systematic issues across files
- **Consider context** - Business requirements may justify exceptions

## Best Practices

1. **Systematic Review**: Check all modified files thoroughly
2. **Priority Focus**: Address P0 issues before P1/P2
3. **Clear Communication**: Use exact file paths and line numbers
4. **Actionable Feedback**: Provide specific code suggestions
5. **Standards Adherence**: Reference specific rule numbers
6. **Save Results**: Always output to `code_review.result` file
