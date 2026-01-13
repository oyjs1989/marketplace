package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// ============================================================
// 4.1 KISS 原则违规示例
// ============================================================

// 4.1.1 违规：不必要的接口抽象（只有一个实现）
type ConfigLoader interface {
	Load(path string) (*Config, error)
}

type FileConfigLoader struct{}

func (f *FileConfigLoader) Load(path string) (*Config, error) {
	return &Config{}, nil // 唯一实现，接口没有价值
}

// 4.1.2 违规：函数过长（超过50行）
func CreateUserWithManyResponsibilities(ctx context.Context, user *User) error {
	// 验证逻辑 (10行)
	if user == nil {
		return errors.New("user is nil")
	}
	if user.Name == "" {
		return errors.New("name empty")
	}
	if user.Email == "" {
		return errors.New("email empty")
	}

	// 权限检查 (10行)
	if user.Role == "admin" {
		// check admin permission
	}

	// 数据转换 (10行)
	data := convertUserData(user)

	// 数据库操作 (10行)
	db.Create(&data)

	// 发送邮件 (10行)
	sendEmail(user.Email, "Welcome")

	// 记录日志 (10行)
	log.Info("user created")

	// 缓存更新 (10行)
	cache.Set("user", user)

	return nil
}

// 4.1.3 违规：简单转换使用不必要的 goroutine
func ProcessUsersWithUnnecessaryConcurrency(users []*User) []*UserDTO {
	ch := make(chan *UserDTO, len(users))

	for _, user := range users {
		go func(u *User) {
			ch <- &UserDTO{ID: u.ID, Name: u.Name}
		}(user)
	}

	results := make([]*UserDTO, 0, len(users))
	for i := 0; i < len(users); i++ {
		results = append(results, <-ch)
	}
	return results
}

// 4.1.4 违规：复杂的函数式链
func ProcessOrderWithComplexChain(order *Order) error {
	return Optional(order).
		Filter(func(o *Order) bool { return o.Status == "pending" }).
		Map(func(o *Order) (*Order, error) { return validate(o) }).
		FlatMap(func(o *Order) error { return save(o) }).
		OrElse(errors.New("processing failed"))
}

// 4.1.5 违规：使用 map[string]interface{} 而非 struct
func LoadConfigWithMap(data []byte) (map[string]interface{}, error) {
	var config map[string]interface{}
	json.Unmarshal(data, &config)
	name := config["name"].(string)        // 类型断言，容易出错
	age := int(config["age"].(float64))    // JSON数字是float64
	return config, nil
}

// ============================================================
// 4.2 DRY 原则违规示例
// ============================================================

// 4.2.1 违规：重复的验证逻辑
func CreateUser(ctx context.Context, user *User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	if user.Name == "" {
		return errors.New("user name is empty")
	}
	if user.Email == "" {
		return errors.New("user email is empty")
	}
	return repo.Create(ctx, user)
}

func UpdateUser(ctx context.Context, user *User) error {
	// 重复的验证逻辑
	if user == nil {
		return errors.New("user is nil")
	}
	if user.Name == "" {
		return errors.New("user name is empty")
	}
	if user.Email == "" {
		return errors.New("user email is empty")
	}
	return repo.Update(ctx, user)
}

// 4.2.2 违规：魔法数字和字符串
func ListUsers(page, pageSize int) ([]*User, error) {
	if pageSize > 100 { // 魔法数字
		pageSize = 100
	}
	return users, nil
}

func CheckUserRole(role string) bool {
	return role == "admin" || role == "user" // 硬编码字符串
}

// 4.2.3 违规：重复的字段定义
type User struct {
	ID         int64     `gorm:"column:id"`
	CreateTime time.Time `gorm:"column:create_time"`
	UpdateTime time.Time `gorm:"column:update_time"`
	Creator    string    `gorm:"column:creator"`
	Modifier   string    `gorm:"column:modifier"`
	Name       string    `gorm:"column:name"`
	Email      string    `gorm:"column:email"`
}

type Product struct {
	ID         int64     `gorm:"column:id"`        // 重复
	CreateTime time.Time `gorm:"column:create_time"` // 重复
	UpdateTime time.Time `gorm:"column:update_time"` // 重复
	Creator    string    `gorm:"column:creator"`     // 重复
	Modifier   string    `gorm:"column:modifier"`    // 重复
	Name       string    `gorm:"column:name"`
	Price      float64   `gorm:"column:price"`
}

// 4.2.4 违规：重复的错误处理模式
func GetUser(ctx context.Context, id int64) (*User, error) {
	var user User
	err := db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrap(err, "get user: record not found")
		}
		return nil, errors.Wrap(err, "get user failed")
	}
	return &user, nil
}

func GetProduct(ctx context.Context, id int64) (*Product, error) {
	var product Product
	err := db.Where("id = ?", id).First(&product).Error
	// 相同的错误处理模式
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrap(err, "get product: record not found")
		}
		return nil, errors.Wrap(err, "get product failed")
	}
	return &product, nil
}

// 4.2.5 违规：重新实现已有功能
func CreateUserDTO(user *User) *UserDTO {
	// 重新实现时间转换，但项目中已有 convert.TimeToString
	timeStr := user.CreateTime.Format("2006-01-02 15:04:05")
	return &UserDTO{
		ID:   user.ID,
		Name: user.Name,
		Time: timeStr,
	}
}

// ============================================================
// 4.3 YAGNI 原则违规示例
// ============================================================

// 4.3.1 违规：实现未来可能需要的功能
type UserRepository struct {
	db *gorm.DB
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*User, error) {
	return &User{}, nil // 当前需求
}

// 以下方法没有调用者 - "可能以后需要"
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	return &User{}, nil // 没有当前需求
}

func (r *UserRepository) GetByPhone(ctx context.Context, phone string) (*User, error) {
	return &User{}, nil // 没有当前需求
}

func (r *UserRepository) SearchByKeyword(ctx context.Context, keyword string) ([]*User, error) {
	return nil, nil // 没有当前需求
}

// 4.3.2 违规：过度参数化
func CreateUserWithTooManyParams(
	ctx context.Context,
	name, email string,
	enableCache bool,                       // "可能需要"
	auditLog bool,                          // "可能需要"
	notifyAdmin bool,                       // "可能需要"
	customMetadata map[string]interface{}, // "可能需要"
) error {
	// 实际上这些参数都没有被使用
	user := &User{Name: name, Email: email}
	return saveUser(ctx, user)
}

// 4.3.3 违规：预先设计扩展点
type NotificationChannel interface {
	Send(ctx context.Context, recipient string, message string) error
}

type EmailChannel struct{}

func (e *EmailChannel) Send(ctx context.Context, recipient string, message string) error {
	return nil
}

type NotificationStrategy struct {
	channels map[string]NotificationChannel
}

func (s *NotificationStrategy) RegisterChannel(name string, channel NotificationChannel) {
	// 实际上永远只用 email
}

// 4.3.4 违规：未使用的代码
type UserService struct {
	repo UserRepository
}

func (s *UserService) CreateUser(ctx context.Context, user *User) error {
	return s.repo.Create(ctx, user)
}

func (s *UserService) GetUser(ctx context.Context, id int64) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

// 未使用的方法 - 应该删除
func (s *UserService) UpdateUser(ctx context.Context, user *User) error {
	return s.repo.Update(ctx, user) // 没有调用者
}

func (s *UserService) DeleteUser(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id) // 没有调用者
}

// 注释掉的旧实现 - 应该删除
// func (s *UserService) OldCreateUser(ctx context.Context, user *User) error {
//     // 旧逻辑...
//     return nil
// }

// 4.3.5 违规：不必要的配置选项
type ServerConfig struct {
	Port                int
	Host                string
	Timeout             int
	EnableDebugMode     bool              // "可能需要"
	MaxConnectionsPerIP int               // "可能需要"
	EnableMetrics       bool              // "可能需要"
	MetricsPort         int               // "如果启用监控"
	EnableRateLimit     bool              // "可能需要"
	RateLimitPerSecond  int               // "如果启用限流"
	CustomHeader        map[string]string // "可能需要"
	EnableCompression   bool              // "可能提高性能"
}

// ============================================================
// 4.4 SOLID 原则违规示例
// ============================================================

// 4.4.1 违规：单一职责原则 - 多重职责
type UserServiceWithMultipleResponsibilities struct {
	db *gorm.DB
}

func (s *UserServiceWithMultipleResponsibilities) CreateUser(ctx context.Context, user *User) error {
	// 职责 1: 验证
	if user.Name == "" {
		return errors.New("name is required")
	}

	// 职责 2: 发送邮件
	emailContent := fmt.Sprintf("Welcome %s!", user.Name)
	sendEmail(user.Email, emailContent)

	// 职责 3: 数据库操作
	s.db.Create(user)

	// 职责 4: 日志记录
	log.Info(ctx, "user created", log.Int64("user_id", user.ID))

	// 职责 5: 缓存更新
	cache.Set(fmt.Sprintf("user:%d", user.ID), user)

	return nil
}

// 4.4.2 违规：开放封闭原则 - 通过修改扩展
type PaymentService struct{}

func (s *PaymentService) ProcessPayment(ctx context.Context, paymentType string, amount float64) error {
	// 每增加一种支付方式都要修改这个函数
	switch paymentType {
	case "credit_card":
		return s.processCreditCard(amount)
	case "alipay": // 新增需要修改
		return s.processAlipay(amount)
	case "wechat": // 新增需要修改
		return s.processWechat(amount)
	default:
		return errors.New("unsupported payment type")
	}
}

// 4.4.3 违规：接口隔离原则 - 大接口
type UserRepositoryLarge interface {
	GetUser(ctx context.Context, id int64) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	ListUsers(ctx context.Context, offset, limit int) ([]*User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id int64) error
	BatchCreate(ctx context.Context, users []*User) error
	BatchDelete(ctx context.Context, ids []int64) error
	CountUsers(ctx context.Context) (int64, error)
}

// 被迫依赖整个接口，但只使用一个方法
type UserQueryService struct {
	repo UserRepositoryLarge // 只使用 GetUser
}

// 4.4.4 违规：依赖反转原则 - 依赖具体实现
type MySQLUserRepository struct {
	db *gorm.DB
}

type UserServiceDependsOnConcrete struct {
	repo *MySQLUserRepository // 直接依赖具体实现
}

func NewUserServiceWithConcreteDepend(db *gorm.DB) *UserServiceDependsOnConcrete {
	return &UserServiceDependsOnConcrete{
		repo: &MySQLUserRepository{db: db}, // 直接创建具体实现
	}
}

// 4.4.5 违规：里氏替换原则 - 违反接口约定
type Storage interface {
	Save(ctx context.Context, key string, data []byte) error
}

type ReadOnlyStorage struct{}

func (r *ReadOnlyStorage) Save(ctx context.Context, key string, data []byte) error {
	// 违反约定：声称实现 Save 但实际不支持
	return errors.New("not supported")
}

// ============================================================
// 4.5 LoD (迪米特法则) 违规示例
// ============================================================

// 4.5.1 违规：链式调用
type User struct {
	Profile *Profile // 导出
}

type Profile struct {
	Address *Address // 导出
}

type Address struct {
	City string // 导出
}

func PrintUserCity(user *User) {
	// 链式调用暴露内部实现
	city := user.Profile.Address.City // 任何 nil 都会 panic
	fmt.Println(city)
}

// 4.5.2 违规：访问返回对象的内部
type Order struct {
	Items []*OrderItem // 导出
}

type OrderItem struct {
	Product *Product // 导出
	Price   float64  // 导出
}

func ProcessOrder(order *Order) {
	var total float64
	for _, item := range order.Items {
		total += item.Price // 访问返回对象的内部
	}

	for _, item := range order.Items {
		if item.Product != nil {
			fmt.Println(item.Product.Name) // 深层访问
		}
	}
}

// 4.5.3 违规：过度依赖
type OrderService struct {
	userRepo *UserRepository
}

func (s *OrderService) CreateOrder(ctx context.Context, order *Order) error {
	// 错误：通过多层调用访问远程对象
	user, _ := s.userRepo.GetByID(ctx, order.UserID)
	address := user.GetProfile().GetAddress() // 多层链式调用
	city := address.GetCity()                 // 依赖太深

	// 错误：访问全局对象
	emailService := GlobalEmailService
	emailService.Send(user.Email, "Order created")

	return nil
}

// 4.5.4 违规：内部创建依赖
type UserServiceCreatesInternally struct {
	db *gorm.DB
}

func (s *UserServiceCreatesInternally) GetUser(ctx context.Context, id int64) (*User, error) {
	// 在内部创建依赖 - 难以测试
	cache := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	logger := zap.NewProduction()

	repo := &UserRepository{db: s.db}

	user, err := repo.GetByID(ctx, id)
	if err != nil {
		logger.Error("failed to get user", zap.Error(err))
		return nil, err
	}

	cache.Set(ctx, fmt.Sprintf("user:%d", id), user, 0)
	return user, nil
}

// 4.5.5 违规：直接耦合
type OrderServiceTightlyCoupled struct {
	// 直接依赖具体实现
	smtpClient  *SMTPClient
	pushService *PushNotificationService
	smsGateway  *TwilioSMSGateway
}

func (s *OrderServiceTightlyCoupled) ProcessOrder(order *Order) error {
	// 紧密耦合到具体实现
	s.smtpClient.ConnectAndSend(order.Email, "subject", "body")
	s.pushService.SendPushNotification(order.UserID, "message")
	s.smsGateway.SendSMS(order.Phone, "sms content")
	return nil
}

// ============================================================
// 4.6 Composition Over Inheritance 违规示例
// ============================================================

// 4.6.1 违规：重复字段而非嵌入
// User 和 Product 重复定义基础字段（已在 4.2.3 展示）

// 4.6.2 违规：大接口而非组合
type StorageLarge interface {
	Read(ctx context.Context, id int64) ([]byte, error)
	Write(ctx context.Context, id int64, data []byte) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]int64, error)
	Close() error
	Flush() error
	Reset() error
}

type ReadOnlyStorageForced struct{}

func (r *ReadOnlyStorageForced) Read(ctx context.Context, id int64) ([]byte, error) {
	return nil, nil
}

// 不应该存在的方法必须返回错误
func (r *ReadOnlyStorageForced) Write(ctx context.Context, id int64, data []byte) error {
	return errors.New("not supported")
}

func (r *ReadOnlyStorageForced) Delete(ctx context.Context, id int64) error {
	return errors.New("not supported")
}

// 4.6.3 违规：过度使用嵌入
type OrderServiceOverEmbedding struct {
	OrderValidator      // 嵌入
	PriceCalculator     // 嵌入
	NotificationService // 嵌入
}

func (s *OrderServiceOverEmbedding) CreateOrder(ctx context.Context, order *Order) error {
	// 不清楚方法来自哪里
	if err := s.Validate(order); err != nil { // 来自 OrderValidator?
		return err
	}

	total := s.Calculate(order.Items) // 来自 PriceCalculator?
	order.Total = total

	s.Notify(order) // 来自 NotificationService?

	return nil
}

// 4.6.4 违规：深层嵌入层次
type Entity struct {
	ID int64
}

type TimestampedEntity struct {
	Entity
	CreateTime time.Time
}

type AuditedEntity struct {
	TimestampedEntity
	Creator string
}

type UserDeepEmbedding struct {
	AuditedEntity // 层次 4 - 太深了！
	Name          string
}

// 4.6.5 违规：大接口而非多个小接口
type EntityLarge interface {
	Validate() error
	Save(ctx context.Context) error
	Delete(ctx context.Context) error
	Notify(message string) error
	Serialize() ([]byte, error)
	Deserialize(data []byte) error
	GetID() int64
	SetID(id int64)
}

type OrderWithLargeInterface struct {
	ID int64
}

// 必须实现所有方法，即使不需要
func (o *OrderWithLargeInterface) Validate() error                    { return nil }
func (o *OrderWithLargeInterface) Save(ctx context.Context) error     { return nil }
func (o *OrderWithLargeInterface) Delete(ctx context.Context) error   { return nil }
func (o *OrderWithLargeInterface) Notify(message string) error        { return nil }
func (o *OrderWithLargeInterface) Serialize() ([]byte, error)         { return nil, nil }
func (o *OrderWithLargeInterface) Deserialize(data []byte) error      { return nil }
func (o *OrderWithLargeInterface) GetID() int64                        { return o.ID }
func (o *OrderWithLargeInterface) SetID(id int64)                      { o.ID = id }

// ============================================================
// 4.7 Less is Exponentially More 违规示例
// ============================================================

// 4.7.1 违规：使用反射
func ProcessValueWithReflection(v interface{}) error {
	val := reflect.ValueOf(v)

	switch val.Kind() {
	case reflect.String:
		str := val.String()
		return processString(str)
	case reflect.Int:
		i := val.Int()
		return processInt(int(i))
	case reflect.Ptr:
		if val.Elem().Type().Name() == "User" {
			return processUser(val.Interface().(*User))
		}
	}
	return nil
}

// 4.7.2 违规：过度使用泛型
type UserServiceOverGeneric[T User | *User, R UserRepository | *UserRepository] struct {
	repo R
}

func (s *UserServiceOverGeneric[T, R]) Create[C context.Context](ctx C, user T) error {
	// 泛型没有带来任何好处，只是增加复杂性
	return s.repo.Create(ctx, user)
}

// 4.7.3 违规：不必要的第三方依赖
import (
	"github.com/some/json-lib"      // 标准库可以做
	"github.com/another/time-lib"   // 标准库可以做
	"github.com/fancy/http-wrapper" // 标准库可以做
)

// 4.7.4 违规：过度抽象
type MessageChannel interface {
	Send(recipient Recipient, content Content) error
}

type Recipient interface {
	GetDestination() string
	GetType() RecipientType
}

type Content interface {
	GetBody() string
	GetMetadata() map[string]interface{}
}

type MessageFactory interface {
	CreateMessage(channelType ChannelType, recipient Recipient, content Content) Message
}

type MessageDispatcher struct {
	factory  MessageFactory
	channels map[ChannelType]MessageChannel
}

// 4.7.5 违规：大而全的包
// package notification 包含邮件、短信、推送、历史、分析等所有功能
// 数十个类型和函数都在一个包里

// ============================================================
// 4.8 Explicit Over Implicit 违规示例
// ============================================================

// 4.8.1 违规：使用 panic/recover 隐藏错误
func ProcessUserWithPanic(ctx context.Context, userID int64) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("panic: %v", r)
		}
	}()

	// 错误被 panic 隐藏
	user := mustGetUser(ctx, userID)  // 内部 panic
	mustValidate(user)                // 内部 panic
	mustProcess(user)                 // 内部 panic

	return nil
}

func mustGetUser(ctx context.Context, userID int64) *User {
	user, err := repo.GetByID(ctx, userID)
	if err != nil {
		panic(err) // 隐藏错误
	}
	return user
}

// 4.8.2 违规：依赖隐式类型转换
func CalculateTotalWithImplicit(items []interface{}) interface{} {
	var total float64
	for _, item := range items {
		// 隐式依赖类型断言
		itemMap := item.(map[string]interface{})
		price := itemMap["price"].(float64)
		quantity := itemMap["quantity"].(float64)
		total += price * quantity
	}
	return total // 返回类型不明确
}

// 4.8.3 违规：依赖零值
func NewUserServiceRelyOnZeroValue(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
		// cache: nil - 依赖零值
		// mu: sync.RWMutex{} - 依赖零值
		// timeout: 0 - 依赖零值
		// enabled: false - 依赖零值，语义不明确
	}
}

// 4.8.4 违规：全局变量和 init()
var (
	globalUserRepo   UserRepository
	globalProductRepo ProductRepository
	globalLogger     Logger
)

func init() {
	db := connectDB()
	globalUserRepo = NewUserRepository(db)
	globalProductRepo = NewProductRepository(db)
	globalLogger = NewLogger()
}

type OrderServiceWithGlobals struct{}

func (s *OrderServiceWithGlobals) CreateOrder(ctx context.Context, order *Order) error {
	// 使用全局变量 - 依赖不清晰
	user, err := globalUserRepo.GetByID(ctx, order.UserID)
	if err != nil {
		globalLogger.Error(ctx, "failed", log.ErrorField(err))
		return err
	}
	return nil
}

// 4.8.5 违规：使用 goto
func ProcessTransactionWithGoto(tx *Transaction) error {
	if err := validateTransaction(tx); err != nil {
		goto handleError
	}

	if err := checkBalance(tx); err != nil {
		goto handleError
	}

	if err := deductBalance(tx); err != nil {
		goto rollback
	}

	if err := recordTransaction(tx); err != nil {
		goto rollback
	}

	return nil

rollback:
	refundBalance(tx)
	// 错误信息丢失

handleError:
	return errors.New("transaction failed") // 原始错误丢失
}
