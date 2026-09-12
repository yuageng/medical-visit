# 医药代表学术拜访管理系统

医药代表（MR）学术拜访管理子系统，围绕“创建拜访、现场签到、学术沟通、签退、提交拜访报告”的完整流程，记录拜访过程并进行合规校验，同时提供按产品统计的月度拜访看板。

## 项目简介

系统用于模拟医疗 CRM 中的学术拜访场景：

- MR 预约医院、科室、医生及拜访产品。
- 签到和签退时手动录入 GPS 坐标，模拟移动端定位采集。
- 自动计算拜访停留时长，并校验拜访地点与医院坐标的距离。
- 默认将停留不足 5 分钟或距离医院超过 500 米的拜访标记为异常。
- 签退后补充谈话要点、医生反馈及学术资料派发情况。
- 按产品和月份查看拜访数量统计。

本项目不包含完整的认证、权限和移动端定位能力，适合作为业务流程演示和后端工程实践项目。

## 技术栈

- 后端：Go、Gin、GORM
- 数据库：PostgreSQL 16
- 前端：React、TypeScript、Vite
- 部署：Docker、Docker Compose、Nginx
- 测试：Go 单元测试及竞态检测

## 目录结构

```text
.
├── cmd/api/                 # API 服务入口
├── internal/                # 领域模型、应用服务、持久化和 HTTP 层
├── migrations/              # SQL 迁移脚本
├── web/                     # React 前端
├── Dockerfile               # 后端镜像构建文件
├── docker-compose.yaml      # 一键部署编排文件
└── README.md
```

## 一键部署

环境要求：已安装并启动 Docker Desktop，以及支持 Compose v2 的 Docker Compose。

在项目根目录执行：

```bash
docker compose up -d --build
```

启动完成后访问：

- 前端：http://localhost
- 后端健康检查：http://localhost/healthz
- API Ping：http://localhost/api/v1/ping

查看服务状态和日志：

```bash
docker compose ps
docker compose logs -f api
```

停止服务但保留数据库数据：

```bash
docker compose down
```

停止服务并删除数据库数据（谨慎操作）：

```bash
docker compose down -v
```

### 配置覆盖

Compose 默认使用以下数据库配置：

| 配置项 | 默认值 | 说明 |
| --- | --- | --- |
| `DB_USER` | `medical_visit` | PostgreSQL 用户名 |
| `DB_PASSWORD` | `medical_visit_password` | PostgreSQL 密码 |
| `DB_NAME` | `medical_visit` | 数据库名称 |
| `WEB_PORT` | `80` | 宿主机前端端口 |
| `COMPLIANCE_MIN_DURATION_SECONDS` | `300` | 最短合规拜访时长 |
| `COMPLIANCE_MAX_DISTANCE_METERS` | `500` | 最大合规距离 |
| `LOG_LEVEL` | `info` | 日志级别 |

可以在项目根目录创建 `.env` 覆盖配置，例如：

```dotenv
DB_USER=medical_visit
DB_PASSWORD=change-me
DB_NAME=medical_visit
WEB_PORT=8088
COMPLIANCE_MIN_DURATION_SECONDS=300
COMPLIANCE_MAX_DISTANCE_METERS=500
LOG_LEVEL=info
```

然后重新创建容器：

```bash
docker compose up -d --build
```

首次启动时，API 会在 PostgreSQL 健康后启动，并通过 `AUTO_MIGRATE=true` 自动创建或更新业务表结构。数据存储在 Docker volume `postgres_data` 中。

## 本地开发

### 启动 PostgreSQL

可以只启动数据库：

```bash
docker compose up -d db
```

### 启动后端

确保本地环境变量中的 `DB_HOST` 为 `localhost`，再执行：

```bash
AUTO_MIGRATE=true go run ./cmd/api
```

后端默认监听 `http://localhost:8080`。

### 启动前端

```bash
cd web
npm ci
npm run dev
```

前端开发服务器默认运行在 `http://localhost:5173`，并将 `/api` 请求代理到本地后端 `8080` 端口。

## API 概览

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| `POST` | `/api/v1/visits` | 创建拜访 |
| `GET` | `/api/v1/visits/:id` | 查询拜访详情 |
| `POST` | `/api/v1/visits/:id/check-in` | 拜访签到 |
| `POST` | `/api/v1/visits/:id/check-out` | 拜访签退并执行合规校验 |
| `POST` | `/api/v1/visits/:id/report` | 保存拜访报告 |
| `GET` | `/api/v1/dashboard/visits/monthly?month=2026-09` | 查询月度产品统计 |
| `GET` | `/healthz` | 服务健康检查 |

系统同时保留了 `/api/visits` 和 `/api/dashboard` 兼容路径。

## 测试与构建

在项目根目录执行：

```bash
# 运行后端测试
go test ./...

# 运行竞态检测和覆盖率
go test -race -coverprofile=coverage.out ./...

# 编译后端
 go build ./cmd/api

# 构建前端
cd web
npm ci
npm run build
```

项目也提供了 [`Makefile`](Makefile)，可使用 `make help` 查看常用命令。

## 合规规则

拜访签退时系统根据签到与签退信息计算：

1. 拜访持续时间是否达到 `COMPLIANCE_MIN_DURATION_SECONDS`。
2. 签到位置与医院坐标的球面距离是否超过 `COMPLIANCE_MAX_DISTANCE_METERS`。

任一规则不满足时，拜访会被标记为异常，并保留相应的合规原因，便于后续查看和审计。

## 相关文档

业务需求和场景说明见 [`description.md`](description.md)。环境变量示例见 [`.env.example`](.env.example)。

## 注意事项

- 默认密码仅用于本地演示，生产环境必须通过 `.env` 或部署平台密钥管理进行替换。
- 当前 Compose 配置使用自动迁移便于一键启动；生产环境建议改用版本化迁移工具执行 [`migrations/`](migrations) 中的 SQL 脚本。
- 业务接口当前未实现身份认证和权限隔离，不应直接暴露到公网。
