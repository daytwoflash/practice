# 博客 API 作业项目

基于 Gin + GORM + SQLite 的 RESTful API，实现用户注册登录、文章 CRUD、评论管理，并包含统一错误处理与结构化日志（`log/slog`）。

## 运行环境

| 项目 | 要求 |
|------|------|
| Go | 1.21+（项目 `go.mod` 为 1.25.0） |
| 操作系统 | Windows / macOS / Linux |
| 数据库 | SQLite（纯 Go 驱动，无需单独安装数据库服务） |
| 测试工具 | Postman、curl 或 Apifox 等 |

## 项目结构

```
project/
├── main.go              # 程序入口
├── config.yaml          # 配置文件
├── go.mod / go.sum      # 依赖管理
├── configs/             # 配置加载（Viper）
├── handlers/            # HTTP 处理器
├── services/            # 业务逻辑
├── models/              # 数据模型与 DTO
├── middlewares/         # 中间件（Auth、Logger、Recovery）
├── logger/              # slog 日志封装
└── utils/               # 响应、JWT、错误处理
```

## 依赖安装

进入项目目录：

```bash
cd 04-homework/project
```

安装主要依赖（若已有 `go.mod` / `go.sum`，可跳过 `go get`，直接 tidy）：

```bash
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get github.com/glebarez/sqlite
go get github.com/golang-jwt/jwt/v5
go get github.com/spf13/viper

go mod tidy
```

## 配置说明

编辑 `config.yaml`：

```yaml
server:
  port: "8080"
  host: "0.0.0.0"
  mode: "debug"   # debug：文本日志；release：JSON 日志

database:
  dbname: "myapp" # SQLite 数据库文件名

jwt:
  secret: "your-secret-key-change-in-production"
  expire: "24h"
```

启动后会在当前目录生成 SQLite 文件（默认 `myapp`）。

## 启动方式

```bash
go run main.go
```

成功启动后，终端可见类似日志：

```
level=INFO msg="config loaded" mode=debug
level=INFO msg="database connected" db=myapp
level=INFO msg="database migrated"
level=INFO msg="server starting" addr=0.0.0.0:8080
```

服务默认地址：`http://localhost:8080`

## 统一响应格式

**成功：**

```json
{
  "code": 200,
  "message": "success",
  "data": { }
}
```

**错误：**

```json
{
  "code": 200,
  "message": "Post not found",
  "error": "Post not found"
}
```

HTTP 状态码与业务含义对应，例如 404、401、403、422、500。

## API 接口

### 健康检查

| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| GET | `/health` | 否 | 服务存活检查 |

### 用户认证

| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| POST | `/api/v1/auth/register` | 否 | 用户注册 |
| POST | `/api/v1/auth/login` | 否 | 用户登录，返回 JWT |

**注册请求体：**

```json
{
  "username": "alice",
  "email": "alice@example.com",
  "password": "123456"
}
```

**登录请求体：**

```json
{
  "username": "alice",
  "password": "123456"
}
```

**登录成功响应示例：**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJ1c2VybmFtZSI6ImFsaWNlIiwiZXhwIjoxNzgxMzQ0Njg5LCJuYmYiOjE3ODEyNTgyODksImlhdCI6MTc4MTI1ODI4OX0.j7rJZhwsw8jEIbIn1cGEVb__uS2bnDB8Ye89HpzXG2A",
    "user": {
      "id": 1,
      "username": "alice",
      "email": "alice@example.com",
      "created_at": "2026-06-12T10:00:00Z"
    }
  }
}
```

### 文章

| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| GET | `/api/v1/posts` | 否 | 文章列表 |
| GET | `/api/v1/posts/:id` | 否 | 文章详情 |
| POST | `/api/v1/posts` | 是 | 创建文章 |
| PUT | `/api/v1/posts/:id` | 是 | 更新文章（仅作者） |
| DELETE | `/api/v1/posts/:id` | 是 | 删除文章（仅作者） |

**鉴权 Header：**

```
Authorization: Bearer <token>
```

**创建文章请求体（作者 id 从 JWT 获取，无需在 body 中传递）：**

```json
{
  "title": "我的第一篇文章",
  "content": "文章内容..."
}
```

**更新文章请求体（id 通过 URL 路径传递）：**

```
PUT /api/v1/posts/1
```

```json
{
  "title": "更新后的标题",
  "content": "更新后的内容"
}
```

### 评论

| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| GET | `/api/v1/posts/:id/comments` | 否 | 获取某文章下的评论列表 |
| POST | `/api/v1/posts/:id/comments` | 是 | 在某文章下创建评论 |

**创建评论请求体（文章 id 通过 URL 路径传递）：**

```
POST /api/v1/posts/1/comments
```

```json
{
  "content": "写得不错！"
}
```

## 错误处理说明

| 场景 | HTTP 状态码 | 示例消息 |
|------|-------------|----------|
| 参数校验失败 | 422 | validation failed |
| 未登录 / Token 无效 | 401 | Invalid token |
| 无权限操作他人文章 | 403 | You can only delete your own posts |
| 用户/文章/评论不存在 | 404 | Post not found |
| 用户名或邮箱重复 | 409 | Username already exists |
| 服务器内部错误 | 500 | Internal server error |

错误由 `utils.HandleError` 统一处理并返回 JSON；同时在日志中记录（4xx 为 Warn，5xx 为 Error）。

## 日志说明

- 使用标准库 `log/slog`，封装于 `logger/` 包
- **运行信息**：启动、数据库连接、每个 HTTP 请求（method、path、status、latency）
- **错误信息**：业务错误、未知错误、panic（Recovery 中间件捕获）
- 日志输出到控制台（`os.Stdout`），不写入数据库

## Postman 测试用例

建议按以下顺序测试，后序步骤依赖前序返回的 `token` 和 `post_id`。

### 用例 1：健康检查

- **请求：** `GET http://localhost:8080/health`
- **预期：** HTTP 200，`data.status` 为 `"ok"`

### 用例 2：用户注册

- **请求：** `POST http://localhost:8080/api/v1/auth/register`
- **Body：**

```json
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "123456"
}
```

- **预期：** HTTP 200，返回用户信息

### 用例 3：重复注册

- **请求：** 同上，再次提交相同用户名
- **预期：** HTTP 409，`Username already exists`

### 用例 4：用户登录

- **请求：** `POST http://localhost:8080/api/v1/auth/login`
- **Body：**

```json
{
  "username": "testuser",
  "password": "123456"
}
```

- **预期：** HTTP 200，返回 `token`，保存供后续使用

### 用例 5：错误密码登录

- **请求：** 同上，`password` 改为错误值
- **预期：** HTTP 401，`Invalid credentials`

### 用例 6：创建文章（需鉴权）

- **请求：** `POST http://localhost:8080/api/v1/posts`
- **Header：** `Authorization: Bearer {{token}}`
- **Body：**

```json
{
  "title": "测试文章",
  "content": "这是测试内容"
}
```

- **预期：** HTTP 200，返回 `post_id`

### 用例 7：获取文章列表

- **请求：** `GET http://localhost:8080/api/v1/posts`
- **预期：** HTTP 200，返回文章数组

### 用例 8：获取文章详情

- **请求：** `GET http://localhost:8080/api/v1/posts/1`
- **预期：** HTTP 200，包含 `title`、`content`

### 用例 9：文章不存在

- **请求：** `GET http://localhost:8080/api/v1/posts/9999`
- **预期：** HTTP 404，`Post not found`

### 用例 10：未带 Token 创建文章

- **请求：** `POST http://localhost:8080/api/v1/posts`（无 Authorization）
- **预期：** HTTP 401

### 用例 11：创建评论

- **请求：** `POST http://localhost:8080/api/v1/posts/1/comments`
- **Header：** `Authorization: Bearer {{token}}`
- **Body：**

```json
{
  "content": "第一条评论"
}
```

- **预期：** HTTP 200，返回评论信息（文章 id 由 URL `/posts/1/comments` 指定）

### 用例 12：获取评论列表

- **请求：** `GET http://localhost:8080/api/v1/posts/1/comments`
- **预期：** HTTP 200，返回该文章下的评论列表

### 用例 13：更新文章

- **请求：** `PUT http://localhost:8080/api/v1/posts/1`
- **Header：** `Authorization: Bearer {{token}}`
- **Body：**

```json
{
  "title": "更新标题",
  "content": "更新内容"
}
```

- **预期：** HTTP 200，`option` 为 `"update"`（文章 id 由 URL `/posts/1` 指定）

### 用例 14：删除文章

- **请求：** `DELETE http://localhost:8080/api/v1/posts/1`
- **Header：** `Authorization: Bearer {{token}}`
- **预期：** HTTP 200，`option` 为 `"delete"`

## 测试结果记录

测试环境：`go run main.go`，`http://localhost:8080`，SQLite 数据库 `myapp`。

| 用例 | 预期 | 实际结果 | 状态 |
|------|------|----------|------|
| 1. 健康检查 | 200 | `{"code":200,"data":{"status":"ok"}}` | ✅ 通过 |
| 2. 用户注册 | 200 | 返回 `id=1, username=testuser` | ✅ 通过 |
| 3. 重复注册 | 409 | `Username already exists` | ✅ 通过 |
| 4. 用户登录 | 200 | 返回 JWT token | ✅ 通过 |
| 5. 错误密码 | 401 | `Invalid credentials` | ✅ 通过 |
| 6. 创建文章 | 200 | 返回 `post_id=1` | ✅ 通过 |
| 7. 无 Token 创建 | 401 | `Authorization header required` | ✅ 通过 |
| 8. 文章列表 | 200 | 返回 1 条记录 | ✅ 通过 |
| 9. 文章详情 | 200 | 包含 `title`、`content` | ✅ 通过 |
| 10. 文章不存在 | 404 | `Post not found` | ✅ 通过 |
| 11. 创建评论 | 200 | 返回评论 `content=nice post` | ✅ 通过 |
| 12. 评论列表 | 200 | 返回 1 条评论 | ✅ 通过 |
| 13. 更新文章 | 200 | `title=updated title` | ✅ 通过 |
| 14. 删除文章 | 200 | `option=delete` | ✅ 通过 |

**合计：14/14 通过**

### 响应示例（节选）

**登录成功：**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": { "id": 1, "username": "testuser", "email": "test@example.com" }
  }
}
```

**文章不存在：**

```json
{
  "code": 200,
  "message": "Post not found",
  "error": "Post not found"
}
```

（HTTP 状态码为 404）

> 可选：附上 Postman 截图或导出的 Collection JSON，作为辅助证明材料。

## 常用命令

```bash
# 编译
go build -o app .

# 运行
go run main.go

# 整理依赖
go mod tidy
```
