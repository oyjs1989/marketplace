---
name: error-safety
description: This agent should be used when reviewing Go code for error handling, nil pointer safety, concurrency control, and JSON processing. Focuses on error handling patterns, safety checks, and concurrent operations (rules 1.1.*, 1.2.*, 1.4.*, 1.5.*).
model: inherit
color: red
---

# Go Error & Concurrency Safety Review

## Purpose

Specialized skill for reviewing error handling, nil safety, concurrency patterns, and JSON processing to ensure robust and safe Go code.

## When to Use

Automatically invoked when reviewing code containing:
- Error returns and handling
- `fmt.Errorf()`, `errors.New()`, `errors.Wrap()`
- Nil pointer checks
- Goroutines, channels, mutexes
- `json.Marshal()`, JSON string construction

## P0 Rules (Must Fix)

### 1.1 Error Handling

#### 1.1.1 Wrap Errors with Stack Trace

Errors must preserve stack trace information.

```go
// ❌ Forbidden - loses stack trace
return fmt.Errorf("operation failed: %v", err)
return fmt.Errorf("get user failed")

// ✅ Correct - preserves stack trace
return errors.Wrapf(err, "operation failed")
return errors.WithMessage(err, "operation failed")
return errors.New("get user failed") // For new errors
```

**Why**: Stack traces are critical for debugging production issues.

#### 1.1.2 Error in Last Return Position

Errors must always be the last return value.

```go
// ✅ Correct
func GetUser(id int64) (*User, error)
func ProcessData(ctx context.Context, data []byte) (result *Result, err error)

// ❌ Wrong
func GetUser(id int64) (error, *User) // Error not last
```

#### 1.1.3 Return Errors, Don't Panic

Business functions must return errors, not panic.

```go
// ❌ Forbidden in business logic
if user == nil {
    panic("user not found")
}

// ✅ Correct
if user == nil {
    return errors.New("user not found")
}
```

**Exception**: Panic is acceptable in:
- `init()` functions
- Unrecoverable programming errors (e.g., nil interface assertions)

#### 1.1.4 Check Errors Immediately

Errors must be checked and handled right away.

```go
// ✅ Correct
user, err := GetUser(id)
if err != nil {
    return errors.Wrapf(err, "failed to get user %d", id)
}
// Use user here...

// ❌ Forbidden - delayed error checking
user, err := GetUser(id)
result, err2 := ProcessUser(user) // Ignoring err!
if err != nil {
    return err
}
```

### 1.2 Nil Pointer Safety

#### 1.2.1 Check Nil Before Dereferencing

Must check nil before using pointers.

```go
// ✅ Correct
if ptr != nil {
    value := *ptr
    // Use value...
}

if user != nil && user.Profile != nil {
    name := user.Profile.Name
}

// ❌ Dangerous - no nil check
value := *ptr // May panic
name := user.Profile.Name // May panic if user or Profile is nil
```

**Pattern for optional fields**:
```go
func ProcessUser(user *User) error {
    // Check required pointer fields
    if user == nil {
        return errors.New("user cannot be nil")
    }

    // Safe to use user.ID, user.Name...

    // Check optional nested pointers
    if user.Profile != nil {
        // Safe to use user.Profile fields
    }

    return nil
}
```

### 1.4 Concurrency Control

#### 1.4.1 Use recovered.ErrorGroup

Use `recovered.ErrorGroup` for concurrent operations with error handling.

```go
// ✅ Correct
import "your-project/pkg/recovered"

errGroup := recovered.NewErrorGroup()

for _, item := range items {
    item := item // Capture loop variable
    errGroup.Go(func() error {
        return ProcessItem(item)
    })
}

if err := errGroup.Wait(); err != nil {
    return errors.Wrapf(err, "failed to process items")
}

// ❌ Avoid - no error handling or panic recovery
var wg sync.WaitGroup
for _, item := range items {
    wg.Add(1)
    go func(i Item) {
        defer wg.Done()
        ProcessItem(i) // Errors ignored, panics unhandled
    }(item)
}
wg.Wait()
```

**Why**: `recovered.ErrorGroup` provides:
- Automatic panic recovery
- Error collection from all goroutines
- Context cancellation on first error

#### 1.4.2 Use Limiter for Concurrency Control

Use limiter to control concurrent goroutine count.

```go
// ✅ Correct
import "your-project/pkg/limiter"

l := limiter.NewConcurrentLimiter(10) // Max 10 concurrent
errGroup := recovered.NewErrorGroup()

for _, item := range items {
    item := item
    l.Acquire()
    errGroup.Go(func() error {
        defer l.Release()
        return ProcessItem(item)
    })
}

// ❌ Dangerous - unlimited goroutines
for _, item := range items {
    go ProcessItem(item) // May spawn thousands of goroutines
}
```

**Why**: Prevents resource exhaustion from too many concurrent operations.

### 1.5 JSON Processing

#### 1.5.1 Never Construct JSON Manually

Do not construct JSON via string concatenation or formatting.

```go
// ❌ Forbidden - prone to syntax errors and injection
jsonStr := fmt.Sprintf(`{"name": "%s", "age": %d}`, name, age)
jsonStr := `{"key": "` + value + `"}`

// ✅ Correct - use json.Marshal
data := map[string]interface{}{
    "name": name,
    "age":  age,
}
jsonBytes, err := json.Marshal(data)
if err != nil {
    return errors.Wrapf(err, "failed to marshal data")
}

// Or use structs
type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}
user := User{Name: name, Age: age}
jsonBytes, err := json.Marshal(user)
```

**Why**: Manual JSON construction is error-prone and vulnerable to injection attacks.

## Review Checklist

### Error Handling
- [ ] All errors wrapped with `errors.Wrapf()` or `errors.WithMessage()`
- [ ] No `fmt.Errorf()` usage (except for new errors without wrapping)
- [ ] Errors in last return position
- [ ] No panic in business logic
- [ ] All errors checked immediately after function calls

### Nil Safety
- [ ] All pointer dereferences have nil checks
- [ ] Nested pointer access protected with nil checks
- [ ] Function parameters validated for nil

### Concurrency
- [ ] Using `recovered.ErrorGroup` for concurrent operations
- [ ] Using `limiter` to control goroutine count
- [ ] No race conditions (check with `go run -race`)
- [ ] Proper use of mutexes for shared state
- [ ] Loop variables captured correctly in goroutines

### JSON Processing
- [ ] No manual JSON string construction
- [ ] Using `json.Marshal()` for all JSON generation
- [ ] Proper error handling for Marshal/Unmarshal operations

## Common Issues

### Issue 1: Lost Stack Trace

```go
// ❌ Problem
func GetUser(id int64) (*User, error) {
    user, err := db.FindUser(id)
    if err != nil {
        return nil, fmt.Errorf("user not found: %v", err) // Stack lost!
    }
    return user, nil
}

// ✅ Fix
func GetUser(id int64) (*User, error) {
    user, err := db.FindUser(id)
    if err != nil {
        return nil, errors.Wrapf(err, "failed to find user %d", id)
    }
    return user, nil
}
```

### Issue 2: Missing Nil Check

```go
// ❌ Problem
func ProcessOrder(order *Order) error {
    amount := order.Payment.Amount // May panic!
    return ValidateAmount(amount)
}

// ✅ Fix
func ProcessOrder(order *Order) error {
    if order == nil {
        return errors.New("order cannot be nil")
    }
    if order.Payment == nil {
        return errors.New("payment information required")
    }
    amount := order.Payment.Amount
    return ValidateAmount(amount)
}
```

### Issue 3: Uncontrolled Concurrency

```go
// ❌ Problem
func ProcessItems(items []Item) error {
    for _, item := range items {
        go func(i Item) {
            ProcessItem(i) // Errors ignored, no limit on goroutines
        }(item)
    }
    return nil
}

// ✅ Fix
func ProcessItems(items []Item) error {
    limiter := limiter.NewConcurrentLimiter(20)
    errGroup := recovered.NewErrorGroup()

    for _, item := range items {
        item := item
        limiter.Acquire()
        errGroup.Go(func() error {
            defer limiter.Release()
            return ProcessItem(item)
        })
    }

    return errGroup.Wait()
}
```

### Issue 4: Manual JSON Construction

```go
// ❌ Problem
func BuildResponse(user *User) string {
    return fmt.Sprintf(`{"id": %d, "name": "%s"}`, user.ID, user.Name)
}

// ✅ Fix
func BuildResponse(user *User) ([]byte, error) {
    response := map[string]interface{}{
        "id":   user.ID,
        "name": user.Name,
    }
    return json.Marshal(response)
}
```

## Output Format

**重要**: 所有问题描述和建议必须使用中文输出。

Report issues in this format (用中文):

```markdown
### 问题 - [P0] 规则 1.X.Y
**位置**: path/to/file.go:45
**类别**: 错误处理 / 空指针安全 / 并发控制 / JSON处理
**原始代码**:
```go
return fmt.Errorf("operation failed: %v", err)
```
**问题描述**: 使用 fmt.Errorf 会丢失错误堆栈信息
**修改建议**:
```go
return errors.Wrapf(err, "operation failed")
```
```

## Reference

For complete standards, see: `references/FUTU_GO_STANDARDS.md` (Sections 1.1, 1.2, 1.4, 1.5)
