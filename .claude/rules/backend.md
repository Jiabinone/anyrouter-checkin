---
globs: ["backend/**/*.go", "**/*.go"]
---

# 后端开发规范

开发 Go 后端代码时，**必须**遵循 `docs/backend/技术规范.md` 中定义的所有规范。

## 核心要点速查

### 技术栈
- Go 1.23+ / Gin / GORM / SQLite
- JWT 认证 / Swagger 文档
- **时间处理**：统一使用 `dromara/carbon`，禁止 `time.Now()`

### 分层架构
```
Router → Handler → Service → Repository → Model
```
- **Handler**：参数绑定、校验、调用 Service、返回响应
- **Service**：业务逻辑编排
- **Repository**：数据库 CRUD
- **Model**：数据结构定义

### 时间处理（carbon）
```go
import "github.com/dromara/carbon/v2"

// ✅ 正确
now := carbon.Now()
datetime := carbon.DateTime{Carbon: carbon.Now()}
stdTime := carbon.Now().StdTime()  // 与标准库交互时

// ❌ 禁止
time.Now()
```

### 统一响应
```go
response.Success(c, data)
response.Error(c, 400, "参数错误")
response.Unauthorized(c)
```

### Swagger 注解
所有 Handler 必须添加完整注解：
```go
// @Summary 简要描述
// @Tags 标签
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RequestType true "参数说明"
// @Success 200 {object} response.Response{data=Type}
// @Router /path [method]
```

### 禁止事项
1. Handler 层直接操作 `repository.DB`
2. 硬编码敏感信息
3. 使用 `fmt.Println` 调试
4. 忽略错误返回值
5. Service 层使用 `*gin.Context`
6. 字符串拼接 SQL
7. 返回未脱敏敏感数据
8. 使用 `time.Now()`（用 `carbon.Now()`）
