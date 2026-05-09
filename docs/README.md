# Mou1ght

Mou1ght 是一个轻量级博客系统，采用 Go 后端 + Vue 管理端 + Next.js 前端的架构。

## 项目结构

```
Mou1ght/
├── cmd/                  # Go 后端入口
├── internal/             # 后端核心代码
│   ├── api/              # API 控制器与中间件
│   ├── config/           # 配置解析
│   ├── domain/           # 数据模型
│   ├── pkg/              # 工具包（数据库、Token 等）
│   ├── repository/       # 数据库操作层
│   └── service/
│       └── router/
│           ├── adminui/  # Admin 前端编译产物（Go embed 嵌入）
│           └── router.go # 路由定义
├── frontend/
│   ├── admin/            # Vue 3 + Naive UI 管理后台
│   └── client/           # Next.js 前端客户端
├── scripts/
│   ├── build.sh          # Linux/macOS 构建脚本
│   └── build.ps1         # Windows 构建脚本
├── config.yaml           # 应用配置文件
├── Dockerfile            # 后端镜像（含 Admin 前端）
├── docker-compose.yaml   # 开发环境编排
├── docker-compose.prod.yaml # 生产环境编排（含 Nginx + SSL）
└── nginx/
    └── conf.d/
        └── default.conf  # Nginx 反向代理配置模板
```

## 前置依赖

| 工具 | 版本要求 | 用途 |
|------|---------|------|
| Go | >= 1.25 | 后端编译 |
| Node.js | >= 20 | 前端编译 |
| pnpm 或 npm | 任意 | 前端包管理 |
| Docker + Docker Compose | Docker >= 20.10 | 容器化部署（可选） |

## 方式一：Docker 部署（推荐）

### 1. 克隆仓库

```bash
git clone --recursive https://github.com/Yoak3n/Mou1ght.git
cd Mou1ght
```

> `--recursive` 会自动拉取 `frontend` 子模块，其中包含 admin 和 client 前端源码。

### 2. 启动服务

```bash
docker compose up -d
```

首次启动会自动完成以下构建：

| 阶段 | 镜像 | 说明 |
|------|------|------|
| admin-builder | `node:20-alpine` | 编译 Vue Admin 前端 |
| backend-builder | `golang:alpine` | 将 Admin 产物嵌入 Go 二进制 |
| runtime | `alpine` | 仅包含最终可执行文件，体积极小 |
| client | `node:20-alpine` | 编译 Next.js standalone |

### 3. 访问

| 服务 | 地址 |
|------|------|
| 前端客户端 | http://localhost:3000 |
| 管理后台 | http://localhost:10420/admin |
| 后端 API | http://localhost:10420/api/v1 |

### 4. 常用命令

```bash
# 查看运行状态
docker compose ps

# 查看日志
docker compose logs -f backend
docker compose logs -f client

# 重建镜像（代码修改后）
docker compose up -d --build

# 仅重建后端
docker compose up -d --build backend

# 停止服务
docker compose down

# 停止并删除数据（SQLite 数据库和上传文件会丢失）
docker compose down -v
```

### 5. 自定义配置

编辑项目根目录的 `config.yaml`，然后重启后端容器：

```bash
docker compose restart backend
```

支持的配置项：

```yaml
blog:
    nav_bar:
        links: []                    # 导航栏链接
        website_information:
            title: "Mou1ght"         # 站点标题
            icon: ""
            logo: ""
            keywords: []             # SEO 关键词
    bottom_extra:
        html: ""                     # 页脚自定义 HTML
        css: ""
    board:
        question: ""                 # 留言板验证问题
        answer: ""                   # 验证答案
        need_reviewed: false         # 留言是否需要审核
database:
    dsn: "Mou1ght"                   # SQLite 数据库文件名
    type: sqlite                     # 数据库类型：sqlite / postgres
security:
    jwt_key: ""                      # JWT 密钥（留空使用环境变量或默认值）
    visitor_jwt_key: ""              # 访客 JWT 密钥
```

JWT 密钥优先级：**配置文件 → 环境变量（`JWT_KEY` / `VISITOR_JWT_KEY`）→ 内置默认值**。

### 6. 使用外部 PostgreSQL（可选）

修改 `config.yaml` 中的数据库配置：

```yaml
database:
    dsn: "host=db-host user=postgres password=your_password dbname=Mou1ght port=5432 sslmode=disable TimeZone=Asia/Shanghai"
    type: postgres
```

同时在 `docker-compose.yaml` 中添加 PostgreSQL 服务并调整 `backend` 的 `depends_on`。

## 方式二：本地编译部署

### 1. 克隆仓库

```bash
git clone --recursive https://github.com/Yoak3n/Mou1ght.git
cd Mou1ght
```

### 2. 使用构建脚本

**Linux / macOS：**

```bash
chmod +x scripts/build.sh
./scripts/build.sh
```

**Windows (PowerShell)：**

```powershell
.\scripts\build.ps1
```

构建脚本会依次执行：

1. 安装前端依赖并编译 `frontend/admin`（产物自动输出到 `internal/service/router/adminui/dist/`）
2. 安装前端依赖并编译 `frontend/client`
3. 运行 `go test`（可通过 `--skip-tests` 跳过）
4. 编译 Go 后端，输出到 `bin/mou1ght`

构建脚本支持以下选项：

```bash
# 跳过测试
./scripts/build.sh --skip-tests

# 仅编译后端（前端已编译好）
./scripts/build.sh --skip-frontend

# 仅编译 Admin 前端 + 后端
./scripts/build.sh --skip-client

# 跳过所有前端编译
./scripts/build.sh --skip-frontend
```

### 3. 运行

```bash
# 将编译产物和配置放到一起
cp config.yaml bin/
cp -r frontend/client bin/   # 如果需要本地运行 Next.js 前端

# 启动后端（内含 Admin 管理界面）
cd bin
./mou1ght
```

### 4. 本地运行 Next.js 前端

```bash
cd frontend/client

# 安装依赖
pnpm install

# 开发模式
pnpm dev

# 或生产模式
pnpm build
pnpm start
```

前端默认通过 `http://localhost:10420` 连接后端 API，可通过环境变量覆盖：

```bash
BASE_URL=http://your-server:10420 pnpm start
```

## 生产环境部署

如需部署到带有域名和 SSL 证书的服务器上，请参阅 [生产环境部署指南](deployment.md)。
