package main

import (
	"fmt"
	"net/http"
	// 违反 Rule 2.3.2 - import 分组不规范
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

// 违反 Rule 2.3.1 - 嵌入字段未放在开头
type Server struct {
	port int
	name string
	http.Server  // 嵌入字段应该在最前面
	startTime time.Time
}

// Config 配置 - 违反 Rule 3.3.1
type Config struct {
	Host string
	Port int
	Timeout int  // 应该使用 time.Duration
}

// APIResponse API 响应 - 违反 Rule 2.5.1
type APIResponse struct {
	Code int
	Msg string  // 拼写错误,应该是 Message
	Data interface{}
}

// HandleRequest 处理请求 - 违反 Rule 2.5.4
func HandleRequest(c *gin.Context) {  // 公共函数缺少注释
	// 违反 Rule 2.3.3 - 未使用 make 初始化
	var users []string
	users = append(users, "user1")

	// 违反 Rule 2.3.5 - switch 缺少 default
	action := c.Query("action")
	switch action {
	case "list":
		fmt.Println("list users")
	case "get":
		fmt.Println("get user")
	// 缺少 default case
	}

	c.JSON(http.StatusOK, users)
}

// ProcessData 处理数据 - 违反 Rule 2.5.3
func ProcessData(data []byte) error {
	// 函数过长,超过 50 行
	// ... 假设这里有很多代码 ...
	step1 := string(data)
	step2 := strings.ToLower(step1)
	step3 := strings.TrimSpace(step2)
	// ... 更多处理步骤 ...
	_ = step3

	// 违反 Rule 2.5.7 - 复杂条件未提取
	if (len(data) > 0 && step2 != "" && step3 != "invalid") || (len(step1) < 100 && time.Now().Hour() > 9) {
		fmt.Println("valid data")
	}

	return nil
}

// UserStatus 用户状态 - 违反 Rule 2.5.5
const (
	StatusActive = 1  // Active user
	StatusInactive = 0  // Inactive user
	// 枚举注释应该用中文
)

// 违反 Rule 3.2.2, 3.2.3 - 测试文件命名不规范
// 这个应该在 _test.go 文件中

// testUserCreation 测试用户创建
func testUserCreation() {
	// 测试代码...
}

// 违反 Rule 2.3.8 - 函数顺序混乱
// 私有函数应该在公共函数之后

func internalHelper() string {
	return "helper"
}

// PublicAPI 公共 API
func PublicAPI() {
	_ = internalHelper()
}

// init 函数 - 违反 Rule 2.3.6
func init() {
	// 复杂的初始化逻辑应该移到明确的初始化函数
	globalConfig := &Config{
		Host: "localhost",
		Port: 8080,
		Timeout: 30,
	}
	_ = globalConfig
}
