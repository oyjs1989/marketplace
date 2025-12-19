package service

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

// UserService 用户服务接口 - 遵循 Rule 2.1.10, 2.4.1
type UserService interface {
	GetUser(ctx context.Context, id uint) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
}

// userService 用户服务实现
type userService struct {
	db *gorm.DB
}

// User 用户模型 - 遵循 Rule 1.3.6, 1.3.8
type User struct {
	ID        uint   `gorm:"column:id;primaryKey"`
	Name      string `gorm:"column:name"`
	Email     string `gorm:"column:email"`
	Age       int    `gorm:"column:age"`
	Status    int    `gorm:"column:status"`
	CreatedAt int64  `gorm:"column:created_at"`
	UpdatedAt int64  `gorm:"column:updated_at"`
}

// TableName 表名
func (User) TableName() string {
	return "users"
}

// 常量定义 - 遵循 Rule 2.1.4, 2.1.6
const (
	USER_STATUS_ACTIVE   = 1
	USER_STATUS_INACTIVE = 0
	MAX_PAGE_SIZE        = 100
	DEFAULT_PAGE         = 1
)

// NewUserService 创建用户服务 - 遵循命名规范
func NewUserService(db *gorm.DB) UserService {
	return &userService{
		db: db,
	}
}

// GetUser 获取用户 - 遵循 Rule 1.3.1, 1.3.2, 1.1.1, 2.4.4
func (s *userService) GetUser(ctx context.Context, id uint) (*User, error) {
	var user User

	// 遵循 Rule 1.3.1, 1.3.2 - 显式 WHERE 和 SELECT
	err := s.db.WithContext(ctx).
		Select("id", "name", "email", "age", "status", "created_at", "updated_at").
		Where("id = ?", id).
		Take(&user).Error // 遵循 Rule 1.3.5 - 使用 Take()

	if err != nil {
		// 遵循 Rule 1.1.1 - 使用 errors.Wrap 包装错误
		return nil, errors.Wrapf(err, "failed to get user, user_id=%d", id)
	}

	return &user, nil
}

// UpdateUser 更新用户 - 遵循 Rule 1.3.4, 1.1.1
func (s *userService) UpdateUser(ctx context.Context, user *User) error {
	// 遵循 Rule 1.3.4 - 使用 Updates() 而不是 Save()
	err := s.db.WithContext(ctx).
		Model(&User{}).
		Where("id = ?", user.ID).
		Updates(map[string]interface{}{
			"name":   user.Name,
			"email":  user.Email,
			"age":    user.Age,
			"status": user.Status,
		}).Error

	if err != nil {
		// 遵循 Rule 1.1.1 - 正确包装错误
		return errors.Wrapf(err, "failed to update user, user_id=%d", user.ID)
	}

	return nil
}

// BatchGetUsersRequest 批量获取用户请求 - 遵循 Rule 2.4.6
type BatchGetUsersRequest struct {
	UserIDs []uint
	Limit   int
}

// BatchGetUsers 批量获取用户 - 遵循 Rule 1.4.1, 1.4.2
func (s *userService) BatchGetUsers(ctx context.Context, req *BatchGetUsersRequest) ([]*User, error) {
	users := make([]*User, 0, len(req.UserIDs))

	// 遵循 Rule 1.4.1 - 使用 ErrorGroup
	eg, egCtx := errgroup.WithContext(ctx)

	// 遵循 Rule 1.4.2 - 使用 limiter 控制并发
	eg.SetLimit(10)

	userChan := make(chan *User, len(req.UserIDs))

	for _, id := range req.UserIDs {
		userID := id
		eg.Go(func() error {
			user, err := s.GetUser(egCtx, userID)
			if err != nil {
				// 遵循 Rule 1.1.4 - 正确检查错误
				return errors.Wrapf(err, "failed to get user in batch, user_id=%d", userID)
			}

			// 遵循 Rule 1.2.1 - nil 检查
			if user != nil {
				userChan <- user
			}

			return nil
		})
	}

	// 等待所有 goroutine 完成
	if err := eg.Wait(); err != nil {
		close(userChan)
		return nil, err
	}
	close(userChan)

	// 收集结果
	for user := range userChan {
		users = append(users, user)
	}

	return users, nil
}

// BuildUserJSON 构建用户 JSON - 遵循 Rule 1.5.1
func (s *userService) BuildUserJSON(user *User) (string, error) {
	// 遵循 Rule 1.5.1 - 使用 json.Marshal 而不是手动构造
	data, err := json.Marshal(user)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal user")
	}

	return string(data), nil
}

// GetUserByEmail 通过邮箱获取用户 - 遵循 Rule 1.1.3
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User

	err := s.db.WithContext(ctx).
		Select("id", "name", "email", "age", "status", "created_at", "updated_at").
		Where("email = ?", email).
		Take(&user).Error

	if err != nil {
		// 遵循 Rule 1.1.3 - 返回错误而不是 panic
		return nil, errors.Wrapf(err, "failed to get user by email, email=%s", email)
	}

	return &user, nil
}

// ListUsersRequest 列出用户请求 - 遵循 Rule 2.4.6
type ListUsersRequest struct {
	Page     int
	PageSize int
}

// ListUsers 列出用户 - 遵循命名和日志规则
func (s *userService) ListUsers(ctx context.Context, req *ListUsersRequest) ([]*User, error) {
	// 遵循 Rule 2.2.7 - 入口日志
	// log.Info("list users start", log.Int("page", req.Page), log.Int("page_size", req.PageSize))

	var users []*User

	// 遵循 Rule 2.1.6 - 使用常量而不是魔法数字
	page := req.Page
	if page < DEFAULT_PAGE {
		page = DEFAULT_PAGE
	}

	pageSize := req.PageSize
	if pageSize > MAX_PAGE_SIZE {
		pageSize = MAX_PAGE_SIZE
	}

	offset := (page - 1) * pageSize

	// 遵循 Rule 1.3.2 - 显式指定 SELECT
	err := s.db.WithContext(ctx).
		Select("id", "name", "email", "age", "status", "created_at", "updated_at").
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error

	if err != nil {
		// 遵循 Rule 2.2.1, 2.2.2, 2.2.6 - 正确的日志格式
		// log.Error("list users failed",
		//     log.Int("page", page),
		//     log.Int("page_size", pageSize),
		//     log.ErrorField(err))
		return nil, errors.Wrapf(err, "failed to list users, page=%d, page_size=%d", page, pageSize)
	}

	// 遵循 Rule 2.2.7 - 出口日志
	// log.Info("list users success", log.Int("count", len(users)))

	return users, nil
}

// DeleteUser 删除用户 - 遵循 Rule 2.1.7, 1.3.1
func (s *userService) DeleteUser(ctx context.Context, userID uint) error {
	// 遵循 Rule 2.2.7 - 入口日志
	// log.Info("delete user start", log.Uint("user_id", userID))

	// 遵循 Rule 1.3.1 - 显式 WHERE 条件
	result := s.db.WithContext(ctx).
		Where("id = ?", userID).
		Delete(&User{})

	if result.Error != nil {
		// 遵循 Rule 1.1.1 - 包装错误
		return errors.Wrapf(result.Error, "failed to delete user, user_id=%d", userID)
	}

	// 遵循 Rule 2.2.7 - 出口日志
	// log.Info("delete user success", log.Uint("user_id", userID))

	return nil
}

// CreateUserRequest 创建用户请求 - 遵循 Rule 2.4.6
type CreateUserRequest struct {
	Name  string
	Email string
	Age   int
}

// CreateUser 创建用户 - 遵循多个规则
func (s *userService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
	user := &User{
		Name:   req.Name,
		Email:  req.Email,
		Age:    req.Age,
		Status: USER_STATUS_ACTIVE,
	}

	err := s.db.WithContext(ctx).Create(user).Error
	if err != nil {
		// 遵循 Rule 1.1.1 - 包装错误
		return nil, errors.Wrapf(err, "failed to create user, email=%s", req.Email)
	}

	return user, nil
}

// TODO(jasonouyang): 优化批量查询性能 - 遵循 Rule 2.3.10
func (s *userService) optimizeBatchQuery() {
	// Implementation pending
}
