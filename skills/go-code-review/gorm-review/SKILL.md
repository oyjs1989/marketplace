---
name: GORM Database Review
description: Specialized review for GORM database operations and struct definitions. Checks query patterns, field usage, struct ordering, and database best practices. Automatically used during Go code reviews when GORM code is detected.
---

# GORM Database Review

## Purpose

Specialized skill for reviewing GORM database layer code, ensuring query safety, performance, and proper struct design.

## When to Use

This skill is automatically invoked when reviewing code containing:
- GORM struct tags: `gorm:"column:..."`
- Database operations: `db.Model()`, `db.Where()`, `db.Find()`, `db.Create()`, etc.
- ORM model definitions
- Database query builders

## P0 Rules (Must Fix)

### 1.3.1 Explicit Where Conditions

Queries must explicitly specify Where conditions to prevent accidental full table scans.

```go
// ✅ Correct
db.Where("id = ?", id).Find(&result)
db.Where("status = ? AND user_id = ?", status, userID).Find(&results)

// ❌ Forbidden
db.Find(&result) // Missing WHERE clause
```

### 1.3.2 Explicit Column Selection

Queries must explicitly specify columns (except for constant definitions).

```go
// ✅ Correct
db.Select("id, name, status").Find(&result)
db.Select("id, key, value, create_time").Where(query).Find(&list)

// ❌ Forbidden
db.Find(&result) // SELECT * is implicit
```

**Why**: Reduces data transfer, prevents breaking changes when table schema evolves.

### 1.3.3 Fluent Method Chaining

Chain GORM methods together for readability and assign error at the end.

```go
// ✅ Recommended
err = db.Model(&model.TagValue{}).
    Select("id, `key`, value, status, sync_status, create_time, update_time, creator, modifier").
    Where(query).
    Offset(int(offset)).
    Limit(int(limit)).
    Find(&tagValueList).Error

// ❌ Avoid: Breaking chain unnecessarily
query := db.Model(&model.TagValue{})
query = query.Select("id, name")
if err := query.Where(condition).Error; err != nil {
    return err
}
```

### 1.3.4 Avoid Save() Method

Do not use `Save()` as it has ambiguous behavior.

```go
// ❌ Forbidden
db.Save(&model) // Unclear if creating or updating

// ✅ Use specific methods
db.Create(&model)  // For new records
db.Updates(&model) // For updates
```

### 1.3.5 Prefer Take() Over First()

Use `Take()` instead of `First()` for better performance.

```go
// ✅ Recommended
db.Where("id = ?", id).Take(&user).Error

// ❌ Avoid
db.Where("id = ?", id).First(&user).Error // Adds ORDER BY primary_key
```

**Why**: `First()` automatically adds `ORDER BY primary_key`, causing unnecessary sorting overhead.

### 1.3.6 Explicit GORM Column Tags

All struct fields must explicitly specify column names.

```go
// ✅ Correct
type TagKey struct {
    ID     int64  `gorm:"column:id"`
    Name   string `gorm:"column:name"`
    Status int    `gorm:"column:status"`
}

// ❌ Missing explicit column tags
type TagKey struct {
    ID     int64  `gorm:"column:id"`
    Name   string // ⚠️ No column tag
    Status int
}
```

### 1.3.7 All Struct Fields Must Be Used

**Critical**: Every field defined in a struct MUST be used in business logic.

**Review Process**:
1. Identify all fields in struct definition
2. Trace business flow from creation to completion
3. Verify each field is read/written at least once
4. Flag unused fields for removal

```go
// ❌ Problem: UnusedField and Description never used
type TagKey struct {
    ID          int64  `gorm:"column:id"`
    Key         string `gorm:"column:key"`
    Name        string `gorm:"column:name"`
    UnusedField string `gorm:"column:unused_field"` // ⚠️ Never used
    Description string `gorm:"column:description"`   // ⚠️ Never used
}

func CreateTagKey(key, name string) error {
    tagKey := &TagKey{
        Key:  key,
        Name: name,
        // UnusedField and Description never set
    }
    return db.Create(tagKey).Error
}

// ✅ Solution: Remove unused fields
type TagKey struct {
    ID   int64  `gorm:"column:id"`
    Key  string `gorm:"column:key"`
    Name string `gorm:"column:name"`
}
```

**Flag as issue when**:
- Field defined but never set in Create/Update
- Field queried but never used in business logic
- Field exists only for "future use" (YAGNI violation)

**Exceptions** (don't require usage tracking):
- Audit fields: `CreateTime`, `UpdateTime`, `Creator`, `Modifier`
- Soft delete: `DeleteTime`, `IsDeleted`

### 1.3.8 Struct Field Ordering

Fields must follow priority order for consistency and readability.

**Order**:
1. **ID field** - Always first
2. **Business fields** - By usage frequency (most used first)
3. **Creator/Modifier** - User tracking fields
4. **CreateTime/UpdateTime** - Timestamp fields last

```go
// ✅ Correct ordering
type TagKeyForQuery struct {
    // 1. ID first
    ID         int64                      `gorm:"column:id"`

    // 2. Business fields (by priority)
    Key        string                     `gorm:"column:key"`
    Name       string                     `gorm:"column:name"`
    Status     constants.TagStatus        `gorm:"column:status"`
    SyncStatus constants.TagKeySyncStatus `gorm:"column:sync_status"`
    Editable   bool                       `gorm:"column:editable"`
    Remark     *string                    `gorm:"column:remark"`

    // 3. User tracking
    Creator    string                     `gorm:"column:creator"`
    Modifier   string                     `gorm:"column:modifier"`

    // 4. Timestamps last
    CreateTime time.Time                  `gorm:"column:create_time"`
    UpdateTime time.Time                  `gorm:"column:update_time"`
}

// ❌ Random ordering
type TagKey struct {
    Name       string    `gorm:"column:name"`
    UpdateTime time.Time `gorm:"column:update_time"`
    ID         int64     `gorm:"column:id"` // Should be first
    Creator    string    `gorm:"column:creator"`
}
```

## Review Checklist

When reviewing GORM code, check:

### Query Safety
- [ ] All queries have explicit WHERE conditions
- [ ] All queries explicitly select columns
- [ ] No use of `Save()` method
- [ ] Using `Take()` instead of `First()` where appropriate

### Struct Design
- [ ] All fields have explicit `gorm:"column:"` tags
- [ ] Field ordering follows ID → Business → Tracking → Timestamps
- [ ] All defined fields are actually used in business logic
- [ ] No unused or "future use" fields

### Performance
- [ ] No SELECT * queries
- [ ] No N+1 query patterns
- [ ] Proper use of indexes (check if WHERE columns are indexed)
- [ ] Fluent chaining for readability

### Security
- [ ] No SQL injection risks (always use parameterized queries)
- [ ] Proper input validation before database operations

## Common Issues

### Issue: Unused Fields

```go
// ❌ Problem
type User struct {
    ID          int64  `gorm:"column:id"`
    Name        string `gorm:"column:name"`
    Email       string `gorm:"column:email"`
    PhoneNumber string `gorm:"column:phone_number"` // Never used
    Address     string `gorm:"column:address"`      // Never used
}
```

**Fix**: Remove unused fields or implement the business logic that uses them.

### Issue: Missing Explicit Select

```go
// ❌ Problem
db.Where("user_id = ?", userID).Find(&orders)

// ✅ Fix
db.Select("id, user_id, amount, status, create_time").
    Where("user_id = ?", userID).
    Find(&orders)
```

### Issue: Using Save() Ambiguously

```go
// ❌ Problem
db.Save(&user) // Is this create or update?

// ✅ Fix - Be explicit
if user.ID == 0 {
    db.Create(&user)
} else {
    db.Updates(&user)
}
```

## Output Format

**重要**: 所有问题描述和建议必须使用中文输出。

Report issues in this format (用中文):

```markdown
### 问题 - [P0] 规则 1.3.X
**位置**: path/to/file.go:123
**类别**: GORM/数据库
**原始代码**:
```go
db.Find(&users)
```
**问题描述**: 缺少 WHERE 条件和显式列选择
**修改建议**:
```go
db.Select("id, name, email, status").
    Where("status = ?", activeStatus).
    Find(&users)
```
```

## Reference

For complete standards, see: `../shared/FUTU_GO_STANDARDS.md` (Section 1.3)
