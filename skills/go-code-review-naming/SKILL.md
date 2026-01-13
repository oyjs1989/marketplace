---
name: go-code-review-naming
description: This skill should be used when reviewing Go code for naming conventions and logging standards. Focuses exclusively on struct/function naming, variable naming, constant definitions, and logging format compliance (rules 2.1.*, 2.2.*).
version: 2.0.0
---

# Go Naming & Logging Standards Review

## Purpose

Specialized skill for reviewing naming conventions and logging standards to ensure consistency, readability, and proper observability.

## When to Use

Automatically invoked when reviewing code containing:
- Struct, function, method, constant declarations
- Variable names
- Log statements: `log.Info()`, `log.Error()`, `log.Debug()`, etc.
- Log field definitions

## P1 Rules (Strongly Recommended)

### 2.1 Naming Conventions

#### 2.1.1 Struct Names: CamelCase, Uppercase First Letter

```go
// ✅ Correct
type UserAccount struct{}
type OrderProcessor struct{}
type HTTPClient struct{}

// ❌ Incorrect
type user_account struct{} // snake_case
type orderprocessor struct{} // no camel case
type httpClient struct{} // should be HTTPClient for acronym
```

#### 2.1.2 Struct Fields: CamelCase Matching Tags

```go
// ✅ Correct - field names match tag semantics
type Model struct {
    UserName  string `gorm:"column:user_name"`  // UserName → user_name
    ProductID int64  `gorm:"column:product_id"` // ProductID → product_id
    IsActive  bool   `json:"is_active"`
}

// ❌ Incorrect - mismatched naming
type Model struct {
    Username string `gorm:"column:user_name"` // Should be UserName
    Productid int64  `gorm:"column:product_id"` // Should be ProductID
}
```

#### 2.1.3 Functions/Methods: CamelCase

```go
// ✅ Correct
// Public functions: uppercase first letter
func GetUserByID(id int64) (*User, error)
func ProcessOrder(order *Order) error

// Private functions: lowercase first letter
func validateEmail(email string) bool
func parseTimestamp(ts string) (time.Time, error)

// ❌ Incorrect
func get_user_by_id(id int64) (*User, error) // snake_case
func Processorder(order *Order) error // not camel case
```

#### 2.1.4 Constants: ALL_CAPS with Underscores

```go
// ✅ Correct
const (
    MAX_RETRY_COUNT = 3
    DEFAULT_TIMEOUT = 30
    API_VERSION     = "v1"
)

// ❌ Incorrect
const (
    MaxRetryCount = 3  // Should use ALL_CAPS
    defaultTimeout = 30 // Should use ALL_CAPS
    apiVersion = "v1"
)
```

#### 2.1.5 ID Abbreviation Always Uppercase

```go
// ✅ Correct
UserID, ProductID, OrderID
func GetUserByID(userID int64)

// ❌ Incorrect
UserId, ProductId, OrderId
func GetUserById(userId int64)
```

#### 2.1.6 No Magic Numbers/Strings - Define Constants

```go
// ❌ Problem - magic numbers
func ProcessOrder(order *Order) error {
    if order.Status == 1 { // What is 1?
        time.Sleep(30 * time.Second) // Why 30?
    }
}

// ✅ Fix - use constants
const (
    ORDER_STATUS_PENDING  = 1
    ORDER_STATUS_APPROVED = 2
    RETRY_DELAY_SECONDS   = 30
)

func ProcessOrder(order *Order) error {
    if order.Status == ORDER_STATUS_PENDING {
        time.Sleep(RETRY_DELAY_SECONDS * time.Second)
    }
}
```

#### 2.1.7 Consistent Naming for Same Concept

```go
// ✅ Consistent
CloudTencent, CloudAliyun, CloudAws

// ❌ Inconsistent
TencentCloud, AliyunCloud, CloudAws // Mixed patterns
```

#### 2.1.8 Unambiguous Business Field Names

```go
// ✅ Clear - different purposes clearly named
type Model struct {
    ID         int64  `gorm:"column:id"`          // This table's PK
    TemplateID string `gorm:"column:template_id"` // Cloud provider template ID
    TemplateTid int64 `gorm:"column:template_tid"` // Related template table FK
}

// ❌ Ambiguous
type Model struct {
    ID         int64  `gorm:"column:id"`
    TemplateID string `gorm:"column:template_id"` // Which ID is this?
    Template   int64  `gorm:"column:template"`     // Confusing
}
```

#### 2.1.9 Descriptive Variable Names

```go
// ✅ Clear
userCount := len(users)
templateList := GetTemplates()
configMap := make(map[string]Config)

// ❌ Vague
a := len(users) // What is 'a'?
tmp := GetTemplates() // Avoid 'tmp'
x := make(map[string]Config) // What is 'x'?
```

**Exception**: Short names OK in limited scope:
```go
// Acceptable in small scope
for i, item := range items { // i is clear in loop context
    process(item)
}
```

#### 2.1.10 Interface Names End with -er

```go
// ✅ Correct
type Reader interface { Read(p []byte) (n int, err error) }
type Writer interface { Write(p []byte) (n int, err error) }
type Processor interface { Process(data []byte) error }

// ❌ Incorrect
type ReadInterface interface {} // Avoid "Interface" suffix
type DataProcess interface {} // Should be DataProcessor
```

#### 2.1.11 Opposite Operations Correspond Strictly

```go
// ✅ Correct pairs
Enable() / Disable()
Start() / Stop()
Open() / Close()
Acquire() / Release()
Lock() / Unlock()

// ❌ Mismatched pairs
Enable() / TurnOff() // Should be Enable()/Disable()
Start() / End() // Should be Start()/Stop()
Open() / CloseConnection() // Should be Open()/Close()
```

#### 2.1.12 List Methods Use GetList or Search

```go
// ✅ Correct - clear that multiple items returned
func GetUserList(query Query) ([]*User, error)
func SearchOrders(criteria Criteria) ([]*Order, error)

// ❌ Confusing - sounds like single item
func GetUser(query Query) ([]*User, error) // Misleading
func GetOrders(criteria Criteria) (*Order, error) // Should return []*Order
```

### 2.2 Logging Standards

#### 2.2.1 Log Fields Use snake_case

```go
// ✅ Correct
log.Info(ctx, "user created",
    log.String("user_id", userID),
    log.String("user_name", userName),
    log.Int64("product_id", productID))

// ❌ Incorrect - using PascalCase
log.Info(ctx, "user created",
    log.String("UserID", userID),
    log.String("UserName", userName))
```

**Why**: Consistent with JSON field naming and query filters.

#### 2.2.2 Log Fields Use Explicit Types

```go
// ✅ Correct - explicit types
log.Info(ctx, "processing complete",
    log.Int("count", len(list)),
    log.Int64("user_id", userID),
    log.String("status", status),
    log.Bool("success", true))

// ❌ Incorrect - using log.Any for primitives
log.Info(ctx, "processing complete",
    log.Any("count", len(list)), // Should be log.Int
    log.Any("user_id", userID),  // Should be log.Int64
    log.Any("status", status))   // Should be log.String
```

**Exception**: Use `log.Any()` for:
- Arrays: `log.Any("items", []string{"a", "b"})`
- Structs: `log.Any("config", configStruct)`
- Maps: `log.Any("metadata", map[string]string{})`

#### 2.2.3 Logs Must Include Key Context

```go
// ✅ Correct - includes identifying context
log.Info(ctx, "order processed",
    log.String("user_id", userID),
    log.Int64("order_id", orderID),
    log.String("status", newStatus))

// ❌ Insufficient context
log.Info(ctx, "order processed") // Which order? Which user?
```

#### 2.2.4 Non-Trivial Functions Should Log

```go
// ✅ Good logging practice
func ProcessPayment(ctx context.Context, payment *Payment) error {
    log.Info(ctx, "processing payment",
        log.String("payment_id", payment.ID),
        log.Int64("amount", payment.Amount))

    result, err := gateway.Charge(payment)
    if err != nil {
        log.Error(ctx, "payment failed",
            log.String("payment_id", payment.ID),
            log.ErrorField(err))
        return err
    }

    log.Info(ctx, "payment succeeded",
        log.String("payment_id", payment.ID),
        log.String("transaction_id", result.TransactionID))
    return nil
}
```

#### 2.2.5 Data Layer Must Not Log

Logging should be in **business logic layer**, not data access layer.

```go
// ❌ Incorrect - logging in data layer
func (r *UserRepo) GetUser(ctx context.Context, id int64) (*User, error) {
    log.Info(ctx, "getting user", log.Int64("user_id", id)) // NO!
    return r.db.Where("id = ?", id).Take(&user).Error
}

// ✅ Correct - logging in business layer
func (s *UserService) GetUser(ctx context.Context, id int64) (*User, error) {
    log.Info(ctx, "fetching user", log.Int64("user_id", id)) // YES!
    return s.repo.GetUser(ctx, id)
}
```

#### 2.2.6 Error Logs Must Include ErrorField at End

```go
// ✅ Correct - ErrorField at end
log.Error(ctx, "failed to process order",
    log.String("user_id", userID),
    log.Int64("order_id", orderID),
    log.ErrorField(err)) // Error field last

// ❌ Incorrect - ErrorField not at end
log.Error(ctx, "failed to process order",
    log.ErrorField(err),
    log.String("user_id", userID)) // Error should be last
```

#### 2.2.7 Log Method Entry and Exit

```go
// ✅ Good practice
func ProcessOrder(ctx context.Context, orderID int64) error {
    // Entry log with inputs
    log.Info(ctx, "processing order started",
        log.Int64("order_id", orderID))

    // ... business logic ...

    // Exit log with key results
    log.Info(ctx, "processing order completed",
        log.Int64("order_id", orderID),
        log.String("result_status", status))

    return nil
}
```

#### 2.2.8 Successful Early Returns Must Log (Error Returns Don't Need Logs)

When a method skips subsequent logic and returns successfully (nil error), log the reason. **Exception**: Error returns don't need additional logs - the outer layer will log them.

```go
// ❌ Problem - silent successful early return
func ProcessOrder(ctx context.Context, order *Order) error {
    if order.Status == StatusCompleted {
        return nil // Why did we skip processing? No visibility!
    }

    // main processing logic...
}

// ✅ Correct - log before successful early return
func ProcessOrder(ctx context.Context, order *Order) error {
    if order.Status == StatusCompleted {
        log.Info(ctx, "order already completed, skipping processing",
            log.Int64("order_id", order.ID),
            log.String("status", order.Status))
        return nil // Returning nil = success, need log
    }

    // main processing logic...
}

// ✅ No log needed - returning error (outer layer logs it)
func ProcessOrder(ctx context.Context, order *Order) error {
    if order == nil {
        return errors.New("order is nil") // Error return, no log needed
    }

    if order.Status == StatusInvalid {
        return errors.Errorf("invalid status: %s", order.Status) // Caller logs
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

**Why**: Successful early returns skip business logic silently, making debugging difficult. Error returns are already logged by callers, so duplicate logging is unnecessary.

## Review Checklist

### Naming
- [ ] Structs use CamelCase with uppercase first letter
- [ ] Struct fields match tag naming conventions
- [ ] Functions/methods use CamelCase appropriately
- [ ] Constants use ALL_CAPS with underscores
- [ ] ID always uppercase, not Id
- [ ] No magic numbers or strings
- [ ] Consistent naming for related concepts
- [ ] Descriptive variable names
- [ ] Interface names end with -er
- [ ] Opposite operations match strictly
- [ ] List methods use GetList/Search

### Logging
- [ ] Log fields use snake_case
- [ ] Explicit types for log fields (not log.Any for primitives)
- [ ] Logs include sufficient context
- [ ] Non-trivial functions have entry/exit logs
- [ ] No logging in data layer
- [ ] Error logs have ErrorField at end
- [ ] Successful early returns (nil error) have log statements before returning
- [ ] Error returns do NOT have redundant logs (outer layer logs them)
- [ ] Sensitive data not logged

## Common Issues

### Issue: Inconsistent Naming

```go
// ❌ Problem
func GetUserById(userId int64) (*User, error) {
    // ...
}

// ✅ Fix
func GetUserByID(userID int64) (*User, error) {
    // ...
}
```

### Issue: Wrong Log Field Case

```go
// ❌ Problem
log.Info(ctx, "user created",
    log.String("UserID", userID),
    log.String("UserName", userName))

// ✅ Fix
log.Info(ctx, "user created",
    log.String("user_id", userID),
    log.String("user_name", userName))
```

### Issue: Using log.Any for Primitives

```go
// ❌ Problem
log.Info(ctx, "processing",
    log.Any("count", count),
    log.Any("user_id", userID))

// ✅ Fix
log.Info(ctx, "processing",
    log.Int("count", count),
    log.Int64("user_id", userID))
```

### Issue: Silent Successful Early Return

```go
// ❌ Problem - no log before successful early return
func ProcessData(ctx context.Context, data *Data) error {
    if data.IsProcessed {
        return nil // Silent skip - why was processing skipped?
    }

    // main logic...
}

// ✅ Fix - log before successful early return
func ProcessData(ctx context.Context, data *Data) error {
    if data.IsProcessed {
        log.Info(ctx, "data already processed, skipping",
            log.String("data_id", data.ID),
            log.Bool("is_processed", true))
        return nil // Success return needs log
    }

    // main logic...
}

// ✅ Correct - error returns don't need extra log
func ProcessData(ctx context.Context, data *Data) error {
    if data == nil {
        return errors.New("data is nil") // No log needed, caller logs
    }

    if data.IsProcessed {
        log.Info(ctx, "data already processed, skipping",
            log.String("data_id", data.ID))
        return nil // Success return needs log
    }

    // main logic...
}
```

## Output Format

**重要**: 所有问题描述和建议必须使用中文输出。

Report issues in this format (用中文):

```markdown
### 问题 - [P1] 规则 2.X.Y
**位置**: path/to/file.go:89
**类别**: 命名规范 / 日志规范
**原始代码**:
```go
func GetUserById(userId int64) (*User, error)
```
**问题描述**: ID 应该全部大写,参数应该使用驼峰命名法
**修改建议**:
```go
func GetUserByID(userID int64) (*User, error)
```
```

## Reference

For complete standards, see: `references/FUTU_GO_STANDARDS.md` (Sections 2.1, 2.2)
