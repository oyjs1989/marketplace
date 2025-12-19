package service

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"
)

// 用户服务 - 包含多个违反编码规范的问题
type userService struct {
	db *gorm.DB
	cache map[string]interface{}
	mu sync.Mutex
}

// User 用户模型 - 违反 Rule 1.3.6, 1.3.8
type User struct {
	ID uint
	name string  // 字段未导出
	Email string
	Age int
	CreatedAt time.Time
	status int  // 未使用 gorm tag
}

// GetUser 获取用户 - 违反多个规则
func (s *userService) GetUser(id uint) (*User, error) {
	var user User

	// 违反 Rule 1.3.1, 1.3.2 - 缺少显式 WHERE 和 SELECT
	err := s.db.Find(&user).Error
	if err != nil {
		// 违反 Rule 1.1.1 - 使用 fmt.Errorf 而不是 errors.Wrap
		return nil, fmt.Errorf("get user failed: %v", err)
	}

	// 违反 Rule 1.2.1 - 未检查 nil
	result := user.Email

	return &user, nil
}

// UpdateUser 更新用户 - 违反多个规则
func (s *userService) UpdateUser(user *User) error {
	// 违反 Rule 1.3.4 - 使用 Save() 方法
	err := s.db.Save(user).Error
	if err != nil {
		// 违反 Rule 1.1.1 - 错误未正确包装
		return fmt.Errorf("update failed")
	}
	return nil
}

// BatchGetUsers 批量获取用户 - 违反并发规则
func (s *userService) BatchGetUsers(ids []uint) ([]*User, error) {
	users := make([]*User, 0)

	// 违反 Rule 1.4.1 - 未使用 ErrorGroup
	var wg sync.WaitGroup
	for _, id := range ids {
		wg.Add(1)
		go func(id uint) {
			defer wg.Done()
			user, _ := s.GetUser(id)  // 违反 Rule 1.1.4 - 未检查错误
			users = append(users, user)  // 并发写入未加锁
		}(id)
	}
	wg.Wait()

	return users, nil
}

// BuildUserJSON 构建用户 JSON - 违反 Rule 1.5.1
func (s *userService) BuildUserJSON(user *User) string {
	// 违反 Rule 1.5.1 - 手动构造 JSON
	return fmt.Sprintf(`{"id":%d,"email":"%s"}`, user.ID, user.Email)
}

// GetUserByEmail 通过邮箱获取用户 - 违反多个规则
func (s *userService) GetUserByEmail(email string) (*User, error) {
	var user User

	// 违反 Rule 1.3.5 - 使用 First() 而不是 Take()
	err := s.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		// 违反 Rule 1.1.3 - 可能会导致 panic
		panic("failed to get user")
	}

	return &user, nil
}

// ListUsers 列出用户 - 违反命名和日志规则
func (s *userService) ListUsers(page, pageSize int) ([]*User, error) {
	var users []*User

	// 违反 Rule 2.1.6 - 魔法数字
	if page < 1 {
		page = 1
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	// 违反 Rule 1.3.2 - 未显式指定 SELECT
	err := s.db.Offset(offset).Limit(pageSize).Find(&users).Error
	if err != nil {
		// 违反 Rule 2.2.1, 2.2.2 - 日志字段不规范
		fmt.Printf("ListUsers failed, Error: %v, Page: %d\n", err, page)
		return nil, err
	}

	return users, nil
}

// deleteUser 删除用户 - 违反多个规则
func (s *userService) deleteUser(userId uint) error {  // 违反 Rule 2.1.7 - 命名不一致
	// 违反 Rule 1.3.1 - 缺少 WHERE 条件
	result := s.db.Delete(&User{})

	if result.Error != nil {
		return result.Error  // 违反 Rule 1.1.1 - 未包装错误
	}

	// 违反 Rule 2.2.7 - 缺少日志

	return nil
}

// CreateUser 创建用户
func (s *userService) CreateUser(name, email string, age int) (*User, error) {
	user := &User{
		name: name,  // 违反: 字段未导出
		Email: email,
		Age: age,
	}

	err := s.db.Create(user).Error
	if err != nil {
		return nil, err  // 违反 Rule 1.1.1 - 未包装错误
	}

	return user, nil
}

// 全局变量 - 违反 Rule 2.3.7
var globalCache = make(map[string]interface{})

// ProcessUsers 处理用户 - 违反多个质量规则
func ProcessUsers(db *gorm.DB, ids []uint, callback func(*User)) error {  // 违反 Rule 2.4.4 - ctx 应为第一参数
	// 违反 Rule 2.5.3 - 函数过长
	// ... 假设这里有很多代码 ...

	for _, id := range ids {
		var user User
		db.Find(&user, id)  // 多个违规
		callback(&user)
	}

	return nil
}

// TODO 修复这个函数  - 违反 Rule 2.3.10 - TODO 未标注负责人
func fixMeLater() {
	// ...
}
