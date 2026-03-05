package bad

import (
	"errors"
	"fmt"
)

// User represents a user entity
type User struct {
	ID         int64
	Name       string
	Email      string
	IsActive   bool
	IsBlocked  bool
	Age        int
	Permission string
}

// ❌ 违规示例 1: 深度嵌套 - 应使用 Early Return
// 规则 2.5.12: 优化代码执行顺序，使用 Early Return
func ProcessUser(user *User) error {
	if user != nil {
		if user.IsActive {
			if user.Permission == "write" {
				if !user.IsBlocked {
					// Main logic deeply nested (4 levels)
					fmt.Printf("Processing user: %s\n", user.Name)
					return nil
				} else {
					return errors.New("user is blocked")
				}
			} else {
				return errors.New("no write permission")
			}
		} else {
			return errors.New("user not active")
		}
	} else {
		return errors.New("user is nil")
	}
}

// ❌ 违规示例 2: 循环中的深度嵌套 - 应使用 Early Continue
// 规则 2.5.12: 循环中应使用 Early Continue 减少嵌套
func ProcessUserList(users []*User) {
	for _, user := range users {
		if user.IsActive {
			if !user.IsBlocked {
				if user.Age >= 18 {
					// Main logic deeply nested (3 levels)
					fmt.Printf("Processing adult user: %s\n", user.Name)
				}
			}
		}
	}
}

// ❌ 违规示例 3: 昂贵操作在前 - 应先做廉价检查
// 规则 2.5.12: 廉价检查应该优先执行
func ValidateAndProcess(id int64, db *MockDB) error {
	// 昂贵的数据库查询在前
	user, err := db.GetUserByID(id)
	if err != nil {
		return err
	}

	// 简单的参数验证应该在数据库查询之前
	if id <= 0 {
		return errors.New("invalid id")
	}

	fmt.Printf("Processing user: %s\n", user.Name)
	return nil
}

// ❌ 违规示例 4: 多重 else - 应使用 Early Return
// 规则 2.5.12: 避免多重 else，使用 Early Return
func CheckUserPermission(user *User, action string) (bool, error) {
	if user == nil {
		return false, errors.New("user is nil")
	} else {
		if !user.IsActive {
			return false, errors.New("user not active")
		} else {
			if user.IsBlocked {
				return false, errors.New("user is blocked")
			} else {
				if user.Permission == action {
					return true, nil
				} else {
					return false, fmt.Errorf("permission denied: need %s", action)
				}
			}
		}
	}
}

// ❌ 违规示例 5: 条件检查顺序不当 - 应按失败概率排序
// 规则 2.5.12: 常见失败条件应该优先检查
func CreateOrder(userID int64, productID int64, quantity int) error {
	// 不常见的失败条件在前
	if quantity > 1000 { // 很少发生
		return errors.New("quantity too large")
	}

	// 更常见的失败条件应该在前面
	if userID <= 0 { // 经常发生
		return errors.New("invalid user id")
	}

	if productID <= 0 { // 经常发生
		return errors.New("invalid product id")
	}

	if quantity <= 0 { // 经常发生
		return errors.New("invalid quantity")
	}

	fmt.Printf("Creating order for user %d\n", userID)
	return nil
}

// ❌ 违规示例 6: 复杂的嵌套 if-else 链 - 应拆分为多个 Early Return
// 规则 2.5.12: 复杂条件应该拆分为多个独立检查
func CalculateDiscount(user *User, orderAmount float64) (float64, error) {
	if user != nil && user.IsActive {
		if !user.IsBlocked && orderAmount > 0 {
			if orderAmount >= 1000 && user.Age >= 18 {
				if user.Permission == "vip" {
					return orderAmount * 0.8, nil // 20% discount
				} else if user.Permission == "premium" {
					return orderAmount * 0.9, nil // 10% discount
				} else {
					return orderAmount * 0.95, nil // 5% discount
				}
			} else {
				return orderAmount, nil // No discount
			}
		} else {
			return 0, errors.New("blocked user or invalid amount")
		}
	} else {
		return 0, errors.New("invalid user")
	}
}

// ❌ 违规示例 7: Switch 中缺少早期返回
// 规则 2.5.12: Switch 语句中也应该使用 Early Return
func HandleUserAction(user *User, action string) error {
	switch action {
	case "read":
		if user != nil {
			if user.IsActive {
				fmt.Println("Reading...")
				return nil
			} else {
				return errors.New("user not active")
			}
		} else {
			return errors.New("user is nil")
		}
	case "write":
		if user != nil {
			if user.IsActive {
				if user.Permission == "write" {
					fmt.Println("Writing...")
					return nil
				} else {
					return errors.New("no write permission")
				}
			} else {
				return errors.New("user not active")
			}
		} else {
			return errors.New("user is nil")
		}
	default:
		return errors.New("unknown action")
	}
}

// MockDB for testing
type MockDB struct{}

func (db *MockDB) GetUserByID(id int64) (*User, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}
	return &User{
		ID:        id,
		Name:      "Test User",
		Email:     "test@example.com",
		IsActive:  true,
		IsBlocked: false,
		Age:       25,
	}, nil
}
