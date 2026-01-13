# FUTU Go Coding Standards

**Version**: 1.2.0
**Last Updated**: 2026-01-13
**Owner**: FUTU Development Team

This document contains the complete Go coding standards for FUTU projects. All specialized review skills reference these standards.

## Table of Contents

- [P0 Rules - Must Follow](#p0-rules---must-follow)
  - [1.1 Error Handling](#11-error-handling)
  - [1.2 Nil Pointer Safety](#12-nil-pointer-safety)
  - [1.3 Database Operations (GORM)](#13-database-operations-gorm)
  - [1.4 Concurrency Control](#14-concurrency-control)
  - [1.5 JSON Processing](#15-json-processing)
- [P1 Rules - Strongly Recommended](#p1-rules---strongly-recommended)
  - [2.1 Naming Conventions](#21-naming-conventions)
  - [2.2 Logging Standards](#22-logging-standards)
  - [2.3 Code Organization](#23-code-organization)
  - [2.4 Interface Design](#24-interface-design)
  - [2.5 Code Quality](#25-code-quality)
- [P2 Rules - Suggested](#p2-rules---suggested)
  - [3.1 Project Structure](#31-project-structure)
  - [3.2 Testing Standards](#32-testing-standards)
  - [3.3 Configuration Management](#33-configuration-management)
- [P1 Rules - Design Philosophies](#p1-rules---design-philosophies)
  - [4.1 KISS Principle](#41-kiss-principle)
  - [4.2 DRY Principle](#42-dry-principle)
  - [4.3 YAGNI Principle](#43-yagni-principle)
  - [4.4 SOLID Principles](#44-solid-principles)
  - [4.5 Law of Demeter](#45-law-of-demeter)
  - [4.6 Composition Over Inheritance](#46-composition-over-inheritance)

---

## P0 Rules - Must Follow

### 1.1 Error Handling

**1.1.1** Errors must be wrapped and preserve stack trace
- ❌ Forbidden: `fmt.Errorf("operation failed")`
- ✅ Use: `errors.Errorf()` for creating new errors, `errors.Wrapf(err, "operation failed")` or `errors.WithMessage(err, "operation failed")` for wrapping errors
- ✅ Use `errors.Errorf()` instead of `fmt.Errorf()` for creating new errors
- ❌ Forbidden: `fmt.Sprintf()` or `fmt.Errorf()` for simple error formatting
- If outer layer already uses `errors.WithMessage`, inner layer doesn't need to wrap again

**1.1.2** Error must be in the last return position
- ✅ `func Foo() (result Type, err error)`

**1.1.3** Business functions must return error, not panic
- ❌ Forbidden: `panic("error occurred")`
- ✅ Use: `return errors.New("error occurred")`

**1.1.4** Errors must be checked and handled immediately

**1.1.5** All error types from external interfaces must be handled
When calling external interfaces or APIs, all error types and error codes returned by the interface must be handled. Cannot ignore any error cases.
- Parse all error types and error codes provided by the interface
- Each error type must have corresponding handling logic
- Cannot handle only part of errors, ignoring other error cases
- For unknown errors, must have default handling logic

**Applicable scenarios:**
- Calling third-party APIs or service interfaces
- Processing error responses from external systems
- Scenarios requiring distinction between different error types

---

### 1.2 Nil Pointer Safety

**1.2.1** Must check nil before using pointers
```go
if ptr != nil {
    value := *ptr
}
```

---

### 1.3 Database Operations (GORM)

**1.3.1** Queries must explicitly specify Where conditions
- ✅ `db.Where("id = ?", id).Find(&result)`
- ✅ If there are no query conditions, Where clause can be omitted
- Query conditions should be explicit, cannot rely on implicit dependencies

**1.3.2** Queries must explicitly specify Select columns (except constant definitions)
- ❌ Forbidden: `SELECT *`
- ✅ Use: `db.Select("id, name, status").Find(&result)`
- ✅ Only query fields that are actually used (e.g., if only `Detail` is used, only query `detail`)
- ✅ Use different structs to read different select contents
- ✅ If only `role_id` is used, only query `role_id`, not the entire row

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

// ✅ Solution: Remove unused fields or use them properly
type TagKey struct {
    ID   int64  `gorm:"column:id"`
    Key  string `gorm:"column:key"`
    Name string `gorm:"column:name"`
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

**1.3.9** Use ID to locate records for updates and deletes
When updating or deleting records, use primary key ID to locate records, not other fields.
- ❌ Should not use other fields to locate records, should use `id`
- ✅ Using primary key ensures uniqueness and performance
- ❌ Avoid using non-primary key fields for update or delete operations

```go
// ✅ Correct: Use ID to locate record
err := db.Where("id = ?", id).Updates(&user).Error

// ❌ Incorrect: Use other field to locate record
err := db.Where("name = ?", name).Updates(&user).Error
```

**1.3.10** Use batch operations with lists
When creating or deleting records in batch, use list operations at once, not individual operations in loops.
- ✅ Use a list to batch create, no need for `create` in loop
- ✅ Batch delete: `db.Where("id in ?", idList).Delete(&xxx)`
- ⚠️ CRUD operations in transactions can also fail, all need error handling

```go
// ✅ Correct: Batch create
users := []User{user1, user2, user3}
err := db.Create(&users).Error

// ❌ Incorrect: Create in loop
for _, user := range users {
    err := db.Create(&user).Error  // Poor performance
    if err != nil {
        return err
    }
}
```

**1.3.11** Index naming conventions
Database index names should follow unified conventions, regular indexes and unique indexes use different prefixes.
- ✅ Regular index: `idx_xxx`
- ✅ Unique index: `uk_xxx`
- Index names should clearly express the purpose of the index

**1.3.12** Use orm.Slice for list operations
When processing list data, use `orm.Slice` provided methods instead of manual processing.
- ✅ List already provides `orm.Slice`, should use it
- ❌ Avoid manually implementing list operations
- ✅ Use framework-provided utility methods

**1.3.13** Use `tx` naming for transactions
When opening a transaction, the transaction object should use `tx` naming, not `db`.
- ✅ Opening transaction uses `tx` naming
- ✅ Operations within transaction should use `tx` not `db`
- ✅ Clearly distinguish transaction operations from normal operations

```go
// ✅ Correct: Use tx naming
err := db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&user).Error; err != nil {
        return err
    }
    return tx.Create(&profile).Error
})

// ❌ Incorrect: Use db in transaction
err := db.Transaction(func(tx *gorm.DB) error {
    if err := db.Create(&user).Error; err != nil {  // Should use tx
        return err
    }
    return nil
})
```

---

### 1.4 Concurrency Control

**1.4.1** Use `recovered.ErrorGroup` for concurrent operations
```go
errGroup := recovered.NewErrorGroup()
```

**1.4.2** Use limiter for concurrency control
```go
l := limiter.NewConcurrentLimiter(count)
```

---

### 1.5 JSON Processing

**1.5.1** Do not construct JSON via string concatenation or formatting
- ❌ Forbidden: `fmt.Sprintf(\`{"name": "%s"}\`, name)`
- ✅ Use: `json.Marshal(data)`

**1.5.2** Do not directly fmt JSON structure for structs
Struct JSON serialization should use standard `json` package, cannot directly use `fmt` to print JSON structure.
- ✅ Struct serialization must use `json.Marshal` or `json.Encoder`
- ❌ Forbidden: `fmt.Printf("%+v", struct)` to output struct
- ✅ JSON serialization should be done through standard library

```go
// ✅ Correct: Use json.Marshal
data, err := json.Marshal(user)
if err != nil {
    return err
}
log.Info(string(data))

// ❌ Incorrect: Use fmt to print struct
log.Info(fmt.Sprintf("%+v", user))
```

---

## P1 Rules - Strongly Recommended

### 2.1 Naming Conventions

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

**2.1.13** Struct receiver names should maintain struct characteristics
Struct method receiver names should reflect the struct's characteristics, use meaningful abbreviations.
- ✅ Receiver names maintain struct characteristics, e.g., `mgr` or `m` (for Manager)
- ✅ Receiver names should be concise but meaningful
- ❌ Avoid overly short names (like single letters `a`, `b`)

```go
// ✅ Correct: Receiver name reflects struct characteristics
type UserManager struct {}

func (mgr *UserManager) CreateUser() error {
    // ...
}

// ❌ Incorrect: Receiver name has no meaning
type UserManager struct {}

func (u *UserManager) CreateUser() error {
    // ...
}
```

**2.1.14** Names must match comments
Code naming must be consistent with comments, comments should clearly describe the code's purpose.
- ❌ When name doesn't match comment, need to fix name or comment
- ❌ Comments cannot be too simple, should explain clearly
- ✅ Names and comments should complement each other and remain consistent

---

### 2.2 Logging Standards

**2.2.1** Log fields use snake_case
- ❌ `log.String("UserName", name)`
- ✅ `log.String("user_name", name)`

**2.2.2** Log fields use explicit types (except arrays, structs, maps use log.Any)
- ❌ `log.Any("count", len(list))`
- ✅ `log.Int("count", len(list))`

**2.2.3** Logs must include key context
Logs must include necessary context fields for problem tracking and debugging. Logs are for developers to read, should clearly show what operation is being performed and what the context is.
- ✅ Logs missing necessary fields must be supplemented
- ✅ Use typed log methods like `log.Int32("role_id", role.ID)`
- ✅ Logs should include key information about the operation

```go
// ✅ Correct: Include necessary fields
log.Info(ctx, "operation success",
    log.String("user_id", userID),
    log.Int64("resource_id", resourceID),
    log.Int32("role_id", role.ID))

// ❌ Incorrect: Missing necessary fields
log.Info(ctx, "operation success")
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

**2.2.8** Early returns that skip logic must log (only when NOT returning error)
When a method skips subsequent logic due to conditions and returns successfully (nil error), log the reason before returning. **Exception**: If returning an error, no additional log needed as outer layer will log the error.
```go
// ❌ Problem - silent successful early return
func ProcessOrder(ctx context.Context, order *Order) error {
    if order.Status == StatusCompleted {
        return nil // Why did we skip processing? No visibility!
    }

    // main processing logic...
}

// ✅ Fix - log before successful early return
func ProcessOrder(ctx context.Context, order *Order) error {
    if order.Status == StatusCompleted {
        log.Info(ctx, "order already completed, skipping processing",
            log.Int64("order_id", order.ID),
            log.String("status", order.Status))
        return nil // Returning nil = success, need log
    }

    // main processing logic...
}

// ✅ No log needed - returning error (outer layer will log)
func ProcessOrder(ctx context.Context, order *Order) error {
    if order == nil {
        return errors.New("order is nil") // Error return, no log needed here
    }

    if order.Status == StatusInvalid {
        return errors.Errorf("invalid order status: %s", order.Status) // Error logged by caller
    }

    // main processing logic...
}

// ✅ Another example - validation skip
func ValidateAndProcess(ctx context.Context, data *Data) error {
    if data.SkipValidation {
        log.Info(ctx, "validation skipped by flag",
            log.String("data_id", data.ID),
            log.Bool("skip_validation", true))
        return processDirectly(ctx, data) // Skipping normal flow, need log
    }

    // normal validation flow...
}
```

---

### 2.3 Code Organization

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

**2.3.15** Clear business layering with clear responsibilities
Code must be organized according to business layers, each layer has clear responsibilities, cannot arbitrarily add new layers.
- ✅ Business layering does not add new layers (e.g., ACL layer), should be placed in `pkg` or `biz`
- ❌ `pkg` package cannot call `repo` layer
- ❌ `pb` cannot directly read/write database, must go through data layer conversion
- ✅ Each layer only handles things within its own responsibility scope

```go
// ✅ Correct: Business logic in biz layer
// biz/user.go
type UserService struct {
    repo UserRepo
}

func (s *UserService) CreateUser(ctx context.Context, req *CreateUserReq) error {
    // business logic
    return s.repo.Create(ctx, user)
}

// ❌ Incorrect: Call repo in pkg
// pkg/user.go
func CreateUser(repo UserRepo) error {
    return repo.Create(ctx, user) // pkg should not depend on repo
}
```

**2.3.16** Interface definitions don't need separate package
Interface definitions should be placed together with implementations, or where the interface is used, don't need to create a separate package for interfaces.
- ❌ `interface` doesn't need a separate `package`
- ✅ Interface definitions should be placed where they are used, or together with implementations
- ❌ Avoid creating additional package structure just for interfaces

---

### 2.4 Interface Design

**2.4.1** Interfaces follow minimal principle

**2.4.2** Repeated parameters can be defined in constructor

**2.4.3** ctx parameter required if using ctx
- ✅ If function uses `ctx`, upper layer must pass it in, cannot ignore
- ✅ Context should be passed throughout the entire call chain
- ❌ Cannot ignore Context parameter

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

**2.4.7** Multiple return values of same type must be named
When a function returns multiple values of the same type (except error), use named return values to improve clarity and prevent confusion.
```go
// ❌ Problem - unclear which string is which
func ParseUser(data string) (string, string, error) {
    // returns name and email, but which is which?
}

// ✅ Fix - named return values make it clear
func ParseUser(data string) (name string, email string, err error) {
    // clear what each return value represents
    return "John", "john@example.com", nil
}

// ✅ Also acceptable - different types don't need naming
func GetUserCount() (int, error) {
    // only one non-error return, naming not required
}

// ✅ Single type return doesn't need naming
func GetUserName() (string, error) {
    // only one non-error return, clear purpose
}
```

**2.4.8** Required parameters must consider zero value validity
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

**2.4.9** Avoid ineffective encapsulation
Don't create ineffective encapsulation, directly calling underlying methods may be more convenient.
- ❌ Don't create ineffective encapsulation, directly calling `db` query is more convenient than calling this method
- ✅ Encapsulation should provide real value
- ❌ Avoid encapsulation just for the sake of encapsulation

**2.4.10** Parameter design should be reasonable
Function parameter design should be reasonable, cannot have logical contradictions.
- ❌ Why does getting account need to pass account list? Parameter design is unreasonable
- ✅ Function parameters should conform to business logic
- ✅ Parameter names and purposes should be consistent

---

### 2.5 Code Quality

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

**2.5.9** Use proto utility functions
When constructing protobuf messages, should use `proto.XXX()` utility functions.
- ✅ Construct pb messages using `proto.XXX()`, e.g., `proto.Int64(roleDepart.ID)`
- ✅ Use `proto.String()`, `proto.Int64()` and other methods
- ❌ Avoid direct assignment, using utility functions is safer

```go
// ✅ Correct: Use proto utility functions
pb := &UserPB{
    ID:   proto.Int64(user.ID),
    Name: proto.String(user.Name),
}

// ❌ Incorrect: Direct assignment
pb := &UserPB{
    ID:   &user.ID,  // May have issues
    Name: &user.Name,
}
```

**2.5.10** Avoid duplicate implementations
If there are already other conversion functions or utility functions, should reuse them instead of re-implementing.
- ✅ If there are already other conversion functions, should use them
- ❌ Avoid re-implementing the same functionality
- ✅ Maintain code consistency

**2.5.11** Use cache reasonably
Not all data needs caching, should decide whether to use cache based on business requirements.
- ✅ Cache should be used for scenarios that truly need performance improvement
- ❌ Avoid overusing cache
- ✅ For example, favorites don't need redis cache, can directly read/write db

---

## P2 Rules - Suggested

### 3.1 Project Structure

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

**3.1.9** Use framework-provided utilities
Should prioritize using utility methods provided by frameworks or libraries, instead of implementing yourself.
- ✅ List already provides `orm.Slice`, should use it
- ✅ Prioritize using standard library and framework-provided functionality
- ❌ Avoid reinventing the wheel

**3.1.10** Don't commit debug code
Debug code and initialization content should not be committed to code repository.
- ❌ Debug code should not be committed
- ❌ Initialization content doesn't need to be committed
- ✅ Keep code repository clean

---

### 3.2 Testing Standards

**3.2.1** Unit tests must verify specific structure, not just err != nil
Unit tests must verify specific business logic, cannot just verify that interface doesn't error.
- ✅ Tests need to verify specific logic, not just that interface doesn't error
- ✅ Tests should verify correctness of business logic
- ✅ Need to verify return values and state changes

```go
// ✅ Correct: Verify specific structure
assert.NoError(t, err)
assert.Equal(t, expectedValue, actualValue)
assert.Equal(t, expectedStatus, result.Status)

// ❌ Incorrect: Only check err != nil
if err != nil {
    t.Fatal(err)
}
// Missing verification of actual business logic
```

**3.2.2** Test file naming: `xxx_test.go`

**3.2.3** Test function naming: `TestXxx`

**3.2.4** Unit tests should be lightweight
Unit tests should be lightweight, directly associated with specific functionality, don't need to start complete application.
- ✅ No `biz_test`, directly associate with specific functionality
- ❌ Non-interface unit tests should not use `application.Run`
- ✅ Unit tests should execute quickly, not depend on external services

**3.2.5** Mock implementation standards
Mock implementations should be placed in independent mock files, completed through interface implementation, not hardcode.
- ✅ `mock` should be placed in independent `mock` files
- ✅ `mock` should be completed through interface implementation, not hardcode
- ✅ Mock implementations should be clear and maintainable

**3.2.6** Interface test location standards
Interface tests should be placed in `internal/test` directory, separated from unit tests.
- ✅ Interface tests should be placed in `internal/test` directory
- ✅ Interface tests and unit tests should be separated
- ✅ Keep test code organization clear

**3.2.7** Unit tests must not start services
Unit tests must not start complete services, should use mock database.
- ❌ Unit tests must not start services, should use mock database
- ✅ Tests should run independently, not depend on external services
- ✅ Use mocks to simulate external dependencies

---

### 3.3 Configuration Management

**3.3.1** Configuration updates should take effect dynamically
```go
func (c *Config) OnConfigChange() error {
    // handle config update
}
```

---

## Rule Reference by Category

### By Priority
- **P0 (Must Fix)**: Rules 1.1.* through 1.5.*
- **P1 (Strongly Recommended)**: Rules 2.1.* through 2.5.*
- **P2 (Suggested)**: Rules 3.1.* through 3.3.*

### By Domain
- **Error Handling**: 1.1.*
- **Safety**: 1.2.*
- **Database**: 1.3.*
- **Concurrency**: 1.4.*
- **JSON**: 1.5.*
- **Naming**: 2.1.*
- **Logging**: 2.2.*
- **Organization**: 2.3.*
- **Interfaces**: 2.4.*
- **Quality**: 2.5.*
- **Structure**: 3.1.*
- **Testing**: 3.2.*
- **Config**: 3.3.*

---

## Document History

- **v1.1.0** (2025-01-09): Added rules from GitLab MR reviews
  - Added 7 new P0 rules (1.1.5, 1.3.9-1.3.13, 1.5.2)
  - Added 8 new P1 rules (2.1.13-2.1.14, 2.3.15-2.3.16, 2.4.9-2.4.10, 2.5.9-2.5.11)
  - Added 5 new P2 rules (3.1.9-3.1.10, 3.2.4-3.2.7)
  - Enhanced existing rules with additional clarifications
  - Rules extracted from GitLab MR !7, !17, !27 code reviews

- **v1.0.0** (2025-01-16): Initial comprehensive standards document
  - Consolidated from original Go Code Review skill
  - Added rule numbering system
  - Organized by priority (P0/P1/P2) and category

---

**Maintained by**: FUTU Development Team
**Questions**: Contact team leads or post in #go-standards channel
