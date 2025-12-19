---
name: Go Code Organization & Quality
description: Reviews code structure, interface design, function quality, project organization, and testing standards. Ensures maintainable, well-organized, and testable Go code. Automatically used during Go code reviews.
---

# Go Code Organization & Quality Review

## Purpose

Specialized skill for reviewing code organization, interface design, function structure, project layout, and testing practices.

## When to Use

Automatically invoked when reviewing:
- Code structure and organization
- Interface definitions
- Function signatures and parameters
- Project directory structure
- Test files and test quality
- Import statements

## P1 Rules (Strongly Recommended)

### 2.3 Code Organization

#### 2.3.1 Embedded Structs at Top

```go
// ✅ Correct - embedded struct first
type Model struct {
    orm.BaseModel // Embedded at top
    Name   string
    Status int
}

// ❌ Incorrect - embedded in middle
type Model struct {
    Name   string
    orm.BaseModel // Should be at top
    Status int
}
```

#### 2.3.2 JSON Tags Should Relate to Field Names

```go
// ✅ Correct - clear relationship
type User struct {
    UserID   int64  `json:"user_id"`
    UserName string `json:"user_name"`
    IsActive bool   `json:"is_active"`
}

// ❌ Confusing - mismatched naming
type User struct {
    UserID   int64  `json:"id"` // Lost context
    UserName string `json:"name"` // Ambiguous
}
```

#### 2.3.3 Use make() to Initialize map/slice

```go
// ✅ Correct
slice := make([]string, 0, capacity)
m := make(map[string]int, capacity)

// ❌ Incorrect (except for unmarshal/scan)
var slice []string // Don't rely on nil slice
var m map[string]int // Dangerous for writes
```

**Exception**: It's OK to use `var` for:
- Unmarshal targets: `var result Response; json.Unmarshal(data, &result)`
- Database scan targets: `var user User; db.Find(&user)`

#### 2.3.4 Specify Capacity When Known

```go
// ✅ Correct - pre-allocate capacity
users := make([]*User, 0, len(userIDs))
for _, id := range userIDs {
    users = append(users, &User{ID: id})
}

// ❌ Suboptimal - will reallocate
users := make([]*User, 0) // Capacity unknown, may grow
```

**Why**: Prevents repeated memory reallocations.

#### 2.3.5 Switch Must Have default Branch

```go
// ✅ Correct
switch status {
case StatusActive:
    return "active"
case StatusInactive:
    return "inactive"
default:
    return "unknown" // Handle unexpected values
}

// ❌ Missing default
switch status {
case StatusActive:
    return "active"
case StatusInactive:
    return "inactive"
// What if status is neither?
}
```

#### 2.3.6 Don't Use iota - Specify Values Explicitly

```go
// ✅ Correct - explicit values
const (
    StatusInit    = 0 // 初始化
    StatusPending = 1 // 待处理
    StatusActive  = 2 // 已激活
)

// ❌ Avoid iota - fragile when reordering
const (
    StatusInit = iota // Value changes if order changes
    StatusPending
    StatusActive
)
```

**Why**: Explicit values are stable across refactorings and database storage.

#### 2.3.7 No Global Variables

```go
// ❌ Forbidden - global mutable state
var userCache = make(map[int64]*User)
var connectionPool *sql.DB

// ✅ Correct - use dependency injection
type Service struct {
    cache *Cache
    db    *sql.DB
}

func NewService(cache *Cache, db *sql.DB) *Service {
    return &Service{cache: cache, db: db}
}
```

**Exceptions**:
- Unit test setup: `var testDB *sql.DB` in test files
- Pre-compiled regex: `var emailRegex = regexp.MustCompile(...)`

#### 2.3.8 No init() Functions

```go
// ❌ Avoid init() functions
func init() {
    setupDatabase()
    loadConfig()
}

// ✅ Correct - explicit initialization
func main() {
    db := setupDatabase()
    config := loadConfig()
    // ...
}
```

**Exception**: `init()` allowed in `server` package for service registration.

#### 2.3.9 Function Encapsulation Must Be Meaningful

```go
// ❌ Meaningless wrapper
func GetUserName(user *User) string {
    return user.Name // Just accessing field
}

// ✅ Meaningful encapsulation - adds logic
func GetDisplayName(user *User) string {
    if user.Nickname != "" {
        return user.Nickname
    }
    return user.Name
}
```

#### 2.3.10 NOTE/TODO Must Include Owner

```go
// ✅ Correct - has owner
// TODO(zhangsan): optimize database query performance
// NOTE(lisi): temporary solution until API v2 is ready

// ❌ Missing owner
// TODO: optimize this
// NOTE: temporary solution
```

#### 2.3.11 Chain Calls Should Break Lines

```go
// ✅ Readable - one method per line
result := db.Select("id, name, status").
    Where("status = ?", activeStatus).
    Order("created_at DESC").
    Limit(100).
    Find(&users)

// ❌ Hard to read - all on one line
result := db.Select("id, name, status").Where("status = ?", activeStatus).Order("created_at DESC").Limit(100).Find(&users)
```

### 2.4 Interface Design

#### 2.4.1 Interfaces Follow Minimal Principle

```go
// ✅ Minimal - single responsibility
type UserReader interface {
    GetUser(ctx context.Context, id int64) (*User, error)
}

type UserWriter interface {
    CreateUser(ctx context.Context, user *User) error
    UpdateUser(ctx context.Context, user *User) error
}

// ❌ Too broad - violates ISP
type UserRepository interface {
    GetUser(ctx context.Context, id int64) (*User, error)
    CreateUser(ctx context.Context, user *User) error
    UpdateUser(ctx context.Context, user *User) error
    DeleteUser(ctx context.Context, id int64) error
    ListUsers(ctx context.Context) ([]*User, error)
    // Many methods make interface hard to mock/implement
}
```

#### 2.4.2 Repeated Parameters in Constructor

```go
// ✅ Correct - repeated params in constructor
type Service struct {
    apiKey string
    region string
}

func NewService(apiKey, region string) *Service {
    return &Service{apiKey: apiKey, region: region}
}

func (s *Service) GetUser(ctx context.Context, id int64) (*User, error) {
    // Uses s.apiKey, s.region - no need to pass every time
}

// ❌ Avoid - passing same params repeatedly
func GetUser(ctx context.Context, apiKey, region string, id int64) (*User, error)
func UpdateUser(ctx context.Context, apiKey, region string, user *User) error
```

#### 2.4.3 ctx Parameter Required If Using ctx

```go
// ✅ Correct - function uses ctx, so it's a parameter
func ProcessOrder(ctx context.Context, order *Order) error {
    log.Info(ctx, "processing order") // Uses ctx
    // ...
}

// ❌ Inconsistent - uses context without parameter
func ProcessOrder(order *Order) error {
    ctx := context.Background() // Creating ctx inside
    log.Info(ctx, "processing order")
}
```

#### 2.4.4 ctx as First Parameter

```go
// ✅ Correct
func Process(ctx context.Context, data []byte) error
func GetUser(ctx context.Context, id int64) (*User, error)

// ❌ Wrong position
func Process(data []byte, ctx context.Context) error
func GetUser(id int64, ctx context.Context) (*User, error)
```

#### 2.4.5 Parameter Naming Matches Operation

```go
// ✅ Correct - create doesn't need ID
func CreateUser(ctx context.Context, data *CreateUserData) (int64, error)

// ✅ Correct - update needs ID
func UpdateUser(ctx context.Context, id int64, data *UpdateUserData) error

// ❌ Confusing - create with ID parameter
func CreateUser(ctx context.Context, id int64, data *CreateUserData) error
```

#### 2.4.6 More Than 4 Parameters - Use Struct

```go
// ❌ Too many parameters
func CreateUser(ctx context.Context, name, email, phone, address, role string) error

// ✅ Use struct for clarity
type CreateUserParams struct {
    Name    string
    Email   string
    Phone   string
    Address string
    Role    string
}

func CreateUser(ctx context.Context, params *CreateUserParams) error
```

#### 2.4.7 Required Parameters Consider Zero Value

Use pointer for parameters where zero value is valid business input.

```go
// ✅ Correct - Status 0 is valid, use pointer
type UpdateUserRequest struct {
    ID     int64  `json:"id" validate:"required"`
    Status *int   `json:"status" validate:"required"` // 0 = disabled is valid
    Age    *int   `json:"age"`                        // 0 years may be valid
}

// ✅ Correct - Empty string invalid, use value type
type CreateUserRequest struct {
    Name   string `json:"name" validate:"required"`  // Empty name invalid
    Email  string `json:"email" validate:"required"` // Empty email invalid
    RoleID int64  `json:"role_id" validate:"required"` // 0 invalid by rule
}
```

### 2.5 Code Quality

#### 2.5.1 No Spelling Errors

```go
// ❌ Spelling errors
func ProcessPaymnet(payment *Payment) error // "Paymnet"
type UserAccont struct{} // "Accont"

// ✅ Correct spelling
func ProcessPayment(payment *Payment) error
type UserAccount struct{}
```

#### 2.5.2 Avoid Magic Numbers/Strings

Already covered in 2.1.6 (Naming section).

#### 2.5.3 Moderate Function Length

```go
// ❌ Too long - should split
func ProcessOrder(order *Order) error {
    // 200+ lines of code
    // Validate order
    // Check inventory
    // Calculate price
    // Process payment
    // Update database
    // Send notifications
    // Generate invoice
    // ...
}

// ✅ Split into focused functions
func ProcessOrder(order *Order) error {
    if err := validateOrder(order); err != nil {
        return err
    }
    if err := checkInventory(order); err != nil {
        return err
    }
    if err := processPayment(order); err != nil {
        return err
    }
    return finalizeOrder(order)
}
```

**Guideline**: Functions > 100 lines should be considered for splitting.

#### 2.5.4 Public Functions Must Have Comments

```go
// ✅ Correct - documented public function
// GetUserByID retrieves a user by their unique identifier.
// Parameters:
//   ctx: request context for cancellation and deadlines
//   id: unique user identifier
// Returns:
//   *User: the user object if found
//   error: nil on success, error otherwise
func GetUserByID(ctx context.Context, id int64) (*User, error) {
    // ...
}

// ❌ Missing documentation
func GetUserByID(ctx context.Context, id int64) (*User, error) {
    // ...
}
```

#### 2.5.5 Enum Constants Have Chinese Comments

```go
// ✅ Correct - aligned comments
const (
    StatusInit     = 0 // 初始化
    StatusPending  = 1 // 待处理
    StatusApproved = 2 // 已通过
    StatusRejected = 3 // 已拒绝
)

// ❌ Missing or misaligned comments
const (
    StatusInit = 0
    StatusPending = 1 // 待处理
    StatusApproved = 2 //已通过 (no space)
)
```

#### 2.5.6 Inject Dependencies via Constructor

Already covered in 2.3.7 (avoid globals).

#### 2.5.7 Use Struct Embedding for Code Reuse

```go
// ✅ Good use of embedding
type BaseModel struct {
    ID         int64
    CreateTime time.Time
    UpdateTime time.Time
}

type User struct {
    BaseModel // Reuses common fields
    Name  string
    Email string
}

type Product struct {
    BaseModel // Reuses common fields
    Name  string
    Price float64
}
```

#### 2.5.8 Allow Pointer Passing When Already Pointer

```go
// ✅ Correct - data already pointer, pass directly
type User struct {
    Profile *UserProfile
}

func ProcessUser(user *User) {
    if user.Profile != nil {
        SaveProfile(user.Profile) // OK to pass nil-checked pointer
    }
}

// ❌ Unnecessary complexity
func ProcessUser(user *User) {
    if user.Profile != nil {
        profile := *user.Profile // Unnecessary dereference
        SaveProfile(&profile)    // Then re-reference
    }
}
```

### 3.1 Project Structure (P2)

#### 3.1.1 Package Name Matches Directory

```
✅ Correct:
internal/service/user.go → package service
internal/model/user.go → package model

❌ Incorrect:
internal/service/user.go → package svc
internal/model/user.go → package models
```

#### 3.1.2 Import Grouping and Sorting

```go
// ✅ Correct grouping
import (
    // Standard library
    "context"
    "fmt"
    "time"

    // Third-party
    "github.com/pkg/errors"
    "gorm.io/gorm"

    // Internal
    "gitlab.futunn.com/project/internal/model"
    "gitlab.futunn.com/project/pkg/log"
)

// ❌ No grouping or mixed order
import (
    "github.com/pkg/errors"
    "context"
    "gitlab.futunn.com/project/internal/model"
    "fmt"
)
```

#### Other Project Structure Rules

- 3.1.3: Business logic in `biz/`, data access in `data/`
- 3.1.4: Cloud provider code organized by provider
- 3.1.5: Configuration in `config/`
- 3.1.6: Utilities in `pkg/`
- 3.1.7: Initialization in `server/server.go`
- 3.1.8: Service implementation in `service/`

### 3.2 Testing Standards (P2)

#### 3.2.1 Tests Verify Structure, Not Just err != nil

```go
// ❌ Weak test
func TestGetUser(t *testing.T) {
    user, err := GetUser(123)
    assert.NoError(t, err) // Only checks no error
}

// ✅ Strong test - verifies actual data
func TestGetUser(t *testing.T) {
    user, err := GetUser(123)
    assert.NoError(t, err)
    assert.Equal(t, int64(123), user.ID)
    assert.Equal(t, "John Doe", user.Name)
    assert.NotEmpty(t, user.Email)
}
```

#### 3.2.2 Test File Naming

```
✅ Correct: user_service_test.go
❌ Incorrect: user_service_tests.go, test_user_service.go
```

#### 3.2.3 Test Function Naming

```go
// ✅ Correct
func TestGetUser(t *testing.T) {}
func TestCreateUser_ValidInput(t *testing.T) {}
func TestUpdateUser_NotFound(t *testing.T) {}

// ❌ Incorrect
func Test_GetUser(t *testing.T) {} // Extra underscore
func testGetUser(t *testing.T) {}  // Lowercase
func GetUserTest(t *testing.T) {}  // Wrong order
```

## Review Checklist

### Code Organization
- [ ] Embedded structs at top
- [ ] JSON tags match field names
- [ ] Using make() for slices/maps
- [ ] Capacity specified when known
- [ ] Switch has default branch
- [ ] No iota, values explicit
- [ ] No global variables (except allowed)
- [ ] No init() functions (except allowed)
- [ ] Meaningful function encapsulation
- [ ] TODO/NOTE has owner
- [ ] Chain calls broken into readable lines

### Interface Design
- [ ] Interfaces are minimal
- [ ] Repeated params in constructor
- [ ] ctx parameter present if used
- [ ] ctx as first parameter
- [ ] Parameter naming matches operation
- [ ] Struct used for 4+ parameters
- [ ] Pointers used for valid zero values

### Code Quality
- [ ] No spelling errors
- [ ] No magic numbers/strings
- [ ] Functions reasonable length
- [ ] Public functions commented
- [ ] Enum constants have Chinese comments
- [ ] Dependencies injected
- [ ] Struct embedding used appropriately

### Project Structure
- [ ] Package names match directories
- [ ] Imports grouped and sorted correctly
- [ ] Follows project layout conventions

### Testing
- [ ] Tests verify actual values
- [ ] Test files named correctly
- [ ] Test functions named correctly

## Output Format

```markdown
### Issue - [P1/P2] Rule X.Y.Z
**Location**: path/to/file.go:67
**Category**: Organization / Interface / Quality / Testing
**Original Code**:
```go
func CreateUser(ctx context.Context, name, email, phone, address, role string) error
```
**Problem**: Too many parameters (>4), reduces readability
**Suggestion**:
```go
type CreateUserParams struct {
    Name, Email, Phone, Address, Role string
}
func CreateUser(ctx context.Context, params *CreateUserParams) error
```
```

## Reference

For complete standards, see: `../shared/FUTU_GO_STANDARDS.md` (Sections 2.3, 2.4, 2.5, 3.*)
