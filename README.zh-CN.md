# Enterprise Admin — Go + React 企业级后台脚手架

> **语言：** 中文 · [English](./README.md)

参考 `go-admin` / `go-admin-ui` 设计的全栈后台脚手架，核心特性是 **Actions 模式（零代码 CRUD）** 与 **常规模式（手写逻辑）** 双轨制：

- 简单 CRUD → 声明一个资源即可，**零代码**自动获得列表/详情/创建/编辑/删除 + 类型安全过滤。
- 复杂业务逻辑 → 手写 handler / controller，配合 **类型安全查询构造器**。

## 技术栈

| 层 | 技术 |
|----|------|
| 后端 | Go 1.23 · Gin · GORM（SQLite 默认，支持 MySQL/Postgres）· Swagger(swag) · JWT |
| 前端 | React 18 · refine(@refinedev/core + antd) · Ant Design 5 · Vite 8 · TypeScript |
| 特性 | 统一响应体 `{code,data,msg}` · 类型安全查询 · 多标签页(keep-alive) · 轻量 RBAC · 操作日志 |

## 目录结构

```
enterprise-admin/
├── Dockerfile                 # 极简单镜像部署（sqlite，Go 直 serve 前端）
├── docker-compose.yml         # 一键部署：单端口 :8080，sqlite 落 volume
├── Makefile                   # 常用命令：make gen / build / run / tidy / docker-*
├── backend/                 # Go 服务
│   ├── bin/                 # make build 产物（git 忽略）
│   ├── main.go              # 启动入口 + swagger 注解
│   ├── config/              # yaml 配置（可用 CONFIG_PATH / PORT / SQLITE_PATH 覆盖）
│   ├── db/                  # gorm 连接 + 自动迁移 + 种子数据
│   ├── models/              # 实体(BaseModel / 用户·角色·菜单·部门·字典·日志·上传·商品·订单)
│   ├── core/
│   │   ├── response/        # 统一响应体
│   │   ├── query/           # Actions 模式的 HTTP 过滤器（字段名/类型校验）
│   │   ├── actions/         # ★ Actions 模式：零代码 CRUD 资源
│   │   ├── auth/            # JWT + 轻量 RBAC
│   │   └── middleware/      # CORS / 操作日志
│   ├── cmd/gen/             # gorm/gen 代码生成器（生成 genquery 包）
│   ├── genquery/            # ★ gorm/gen 生成的类型安全查询包（gorm 自带）
│   │   ├── actions/         # ★ Actions 模式：零代码 CRUD 资源
│   │   ├── auth/            # JWT + 轻量 RBAC
│   │   └── middleware/      # CORS / 操作日志
│   ├── controllers/         # 常规模式：登录 / 订单(复杂逻辑) / 仪表盘 / 上传 / 菜单树
│   ├── router/              # 路由聚合（双轨并存）+ SPA 静态服务（WEB_DIR）
│   └── docs/                # swag 生成的 swagger 文档
└── frontend/                # React 前端
    └── src/
        ├── providers/       # 自定义 dataProvider(解析 envelope) + authProvider(JWT)
        ├── context/         # 多标签页状态 TabsContext
        ├── layout/          # AdminLayout + TabsView(keep-alive 多标签)
        ├── resources/meta.ts# 资源元数据（驱动自动 CRUD 页面）
        ├── registry.tsx     # 菜单 + 页面注册表（单一数据源）
        └── pages/           # Login / Dashboard / ResourcePage(自动) / OrderPage(手写)
```

## Docker 极简部署（sqlite，无需 nginx）

Go 二进制单端口同时提供 API 与前端生产构建，无需 nginx、无需独立静态服务、无需外部数据库：

```bash
docker compose up -d --build
# 打开 http://localhost:8080（admin / admin123）
```

或直接用 docker：

```bash
docker build -t enterprise-admin .
docker run -p 8080:8080 -v ea-data:/app/data enterprise-admin
```

sqlite 文件持久化在 `ea-data` volume（`/app/data`），上传文件在 `./uploads`。环境变量覆盖：

| 变量            | 默认值                        | 说明          |
|-----------------|-------------------------------|---------------|
| `PORT`          | `8080`（配置 `server.port`）  | 监听端口      |
| `WEB_DIR`       | `./web`                       | 前端 `dist` 目录 |
| `SQLITE_PATH`   | 配置 `sqlite.path`            | sqlite 文件位置 |
| `CONFIG_PATH`   | `config/settings.yml`         | 配置文件位置  |

## 常用命令（根目录 Makefile）

```bash
make help          # 查看全部命令
make gen           # 生成 gorm 类型安全查询包（改 models/ 后重跑）
make build         # 编译后端二进制到 backend/bin/enterprise-admin
make run           # 本地运行后端（默认 sqlite :8080）
make tidy          # 整理 Go 依赖
make dev-backend   # 启动后端开发服务（:8080）
make dev-frontend  # 启动前端开发服务（:5173，代理 /api 到 :8080）
make dev           # 同时启动前后端开发服务
make docker-build  # 构建极简单镜像（sqlite，Go 直 serve 前端）
make docker-up     # 一键极简部署（sqlite 持久化到 volume，:8080）
```

构建产物运行（需在 `backend/` 下执行，配置文件/数据库均走相对路径）：

```bash
cd backend && ./bin/enterprise-admin   # 监听 :8080，自动建表 + 种子数据
```

## 运行

### 后端
```bash
cd backend
go mod tidy
go run ./cmd/gen         # 生成 gorm/gen 类型安全查询包（genquery/），改模型后重跑（等价 make gen）
go run .                 # 默认 sqlite，监听 :8080，自动建表 + 种子数据
# Swagger: http://localhost:8080/swagger/index.html
```
默认账号：`admin / admin123`

数据库切换（`config/settings.yml`）：`driver: mysql|postgres`，并填写对应 DSN。
容器里建议用 `SQLITE_PATH` 指定 sqlite 文件位置。

重新生成 swagger 文档：
```bash
swag init -g main.go -o docs
```

### 前端
```bash
cd frontend
pnpm install
pnpm dev                 # http://localhost:5173 (代理 /api -> :8080)
```
代理已配置：`/api`、`/login`、`/uploads`、`/swagger`。

前端通过 swagger 生成的 typescript-axios 客户端访问后端（`src/api/`，不要手改），
统一封装在 `src/api/client.ts`（axios 实例 + JWT + `{code,data,msg}` 解包）：
```bash
pnpm gen:api             # 后端接口变更后：swag init → regen → typecheck
pnpm typecheck
```
前端 URL 均为相对路径，同一份构建既可走 Vite 开发代理，也可由 docker 镜像里的 Go 静态服务直接 serve。

## 双轨制实战

### 1. Actions 模式（零代码）
后端只需声明一个资源：
```go
actions.NewResource[models.Product](db).
    WithName("demo/product").
    WithHooks(actions.Hooks[models.Product]{
        BeforeCreate: func(c *gin.Context, m *models.Product) error {
            m.SetCreateBy(auth.GetUserID(c)); return nil
        },
    }).
    Register(v1)   // 自动得到 GET/POST /demo/product, GET/PUT/DELETE /demo/product/:id
```
前端只需在 `resources/meta.ts` 加一条 `ResourceMeta`，并在 `registry.tsx` 注册，
即可获得列表/搜索/创建/编辑表单，**无需写任何页面组件**。

类型安全的 HTTP 过滤（字段名/类型在 schema 层校验，非法字段直接 400）：
```
GET /api/v1/demo/product?name__like=键盘&price__gte=1000&stock__lte=30&sort=-price
```

### 2. 常规模式（手写复杂逻辑）
`controllers/order.go` 使用 **gorm 自带的类型安全查询**（`gorm/gen` 生成的 `genquery` 包），
字段访问器在编译期类型安全，拼错字段名或传错值类型都会编译失败：
```go
var conds []gen.Condition
if minAmount > 0 {
    conds = append(conds, genquery.Order.Amount.Gte(minAmount))   // float64 字段
}
if len(statuses) > 0 {
    conds = append(conds, genquery.Order.Status.In(statuses...))  // string 字段
}
list, _ := genquery.Order.Where(conds...).
    Order(genquery.Order.ID.Desc()).
    Offset((page-1)*size).Limit(size).Find()
total, _ := genquery.Order.Where(conds...).Count()
```
新模型加入类型安全查询只需重新生成：`go run ./cmd/gen`（基于 `models` 包生成 `genquery`）。
手写路由（可携带任意聚合/校验，Actions 模式覆盖不到的逻辑）：
```go
v1.GET("/order", order.List)            // 列表 + 营收汇总（gen 类型安全过滤）
v1.GET("/order/stats", order.Stats)     // 按状态分组
v1.PUT("/order/:id/paid", order.MarkPaid) // 业务动作
```

### 类型安全查询要点
- **手写逻辑用 `gorm/gen`**（`genquery.Order.Amount.Gte(...)` 等）：gorm 自带代码生成器，
  生成带字段类型信息的查询方法，条件在编译期校验。
- Actions 模式的 `core/query.ParseFilters[T]` 在 API 边界校验 HTTP 过滤参数的字段名与值类型。

## API 契约
统一响应：`{ "code": 200, "data": ..., "msg": "ok" }`。`code != 200` 即错误。
- 分页列表：`data = { list, total, page, pageSize }`
- 详情/创建/编辑：`data = 记录对象`
- 删除：`data = null`

## 扩展
- **新增一个零代码资源**：定义 `models.Xxx`，在 `router.go` 调 `actions.NewResource[models.Xxx](db).WithName("demo/xxx").Register(v1)`，前端加 `ResourceMeta` + 注册。
- **新增一个手写资源**：在 `controllers/` 写 handler，在 `router.go` 注册，前端加页面组件。
- **权限**：`core/auth/rbac.go` 的 `RequireRoles(...)` 中间件可挂载到任意路由（admin 角色默认放行）。
