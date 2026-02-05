# 后端开发规范

开发 Go 后端代码时，**必须**遵循以下规范。完整文档：`docs/backend/技术规范.md`

## 技术栈

| 组件 | 技术 |
|------|------|
| 语言 | Go 1.23+ |
| HTTP 框架 | Gin |
| ORM | GORM |
| 数据库 | SQLite |
| 认证 | JWT |
| 时间处理 | dromara/carbon v2 |
| API 文档 | Swagger |

## 分层架构

```
Router → Handler → Service → Repository → Model
```

| 层级 | 职责 | 禁止 |
|------|------|------|
| Handler | 参数绑定、校验、返回响应 | 直接操作数据库 |
| Service | 业务逻辑编排 | 使用 `*gin.Context` |
| Repository | 数据库 CRUD | 业务逻辑 |
| Model | 数据结构定义 | 方法逻辑 |

## 时间处理（carbon）

```go
import "github.com/dromara/carbon/v2"

// ✅ 正确
now := carbon.Now()
datetime := carbon.DateTime{Carbon: carbon.Now()}
stdTime := carbon.Now().StdTime()  // 与标准库交互

// ❌ 禁止
time.Now()
```

**例外**：`time.Duration` 及其常量可用于标准库接口（如 `http.Client.Timeout`）

## 统一响应

```go
response.Success(c, data)          // 成功
response.Error(c, 400, "参数错误")  // 错误
response.Unauthorized(c)           // 未授权
```

## Swagger 注解

所有 Handler 必须添加完整注解：

```go
// CreateAccount 添加账号
// @Summary 添加账号
// @Tags 账号管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateAccountRequest true "参数"
// @Success 200 {object} response.Response{data=model.Account}
// @Router /accounts [post]
func CreateAccount(c *gin.Context) {}
```

## 数据库规范

```go
// 模型定义
type User struct {
    ID        uint            `gorm:"primarykey" json:"id"`
    Username  string          `gorm:"uniqueIndex;size:50" json:"username"`
    Password  string          `gorm:"size:255" json:"-"`  // 不返回敏感字段
    CreatedAt carbon.DateTime `json:"created_at"`
}

// 查询 - 使用参数化，禁止字符串拼接
repository.DB.Where("status = ?", 1).Find(&accounts)
```

## 禁止事项

1. Handler 层直接操作 `repository.DB`
2. 硬编码敏感信息
3. 使用 `fmt.Println` 调试
4. 忽略错误返回值 `_ = someFunc()`
5. Service 层使用 `*gin.Context`
6. 字符串拼接 SQL（防注入）
7. 返回未脱敏敏感数据
8. 使用 `time.Now()`（用 `carbon.Now()`）

## 开发命令

```bash
make dev          # 开发模式
make build        # 编译
make gen-docs     # 生成 Swagger
make fmt          # 格式化
make lint         # 代码检查
```
