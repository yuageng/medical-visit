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

## 运行方式一：Docker 一键启动（推荐）

该方式会同时运行 PostgreSQL 数据库、Go 后端 API 和 React 前端，不需要分别启动后端和前端服务。

环境要求：已安装并启动 Docker Desktop，以及支持 Compose v2 的 Docker Compose。

在项目根目录执行：

```bash
docker compose up -d --build
```

确认三个服务均已运行：

```bash
docker compose ps
```

启动完成后访问：

- 前端：http://localhost
- 后端健康检查：http://localhost/healthz
- API Ping：http://localhost/api/v1/ping

查看全部服务日志，或只查看后端/前端日志：

```bash
docker compose logs -f
docker compose logs -f api
docker compose logs -f web
```

说明：前端容器通过 Nginx 对外提供页面，并将 `/api` 请求代理给后端容器；API 会等待数据库健康后启动，并自动执行数据库表结构迁移。

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

## 运行方式二：本地运行后端和前端

该方式适合开发和完整链路调试。需要一个数据库容器，并保持后端、前端两个终端同时运行。

### 0. 前置检查

在项目根目录确认开发工具可用：

```bash
docker --version
docker compose version
go version
node --version
npm --version
```

如果 Homebrew 安装的 Docker CLI 与 Docker Desktop Engine 出现 API 版本不匹配，可直接使用 Docker Desktop 自带命令：

```bash
/Applications/Docker.app/Contents/Resources/bin/docker compose version
```

### 1. 启动 PostgreSQL

在项目根目录只启动数据库容器：

```bash
docker compose up -d db
docker compose ps db
```

等待状态显示为 `healthy`。数据库通过宿主机 `localhost:5432` 暴露给本地后端。

如果本机 `5432` 已被其他 PostgreSQL 占用，可改用：

```bash
DB_EXPOSE_PORT=15432 docker compose up -d db
```

后续启动后端时相应设置 `DB_PORT=15432`。

### 2. 终端一：启动后端 API

在项目根目录执行：

```bash
DB_HOST=localhost \
DB_PORT=5432 \
DB_USER=medical_visit \
DB_PASSWORD=medical_visit_password \
DB_NAME=medical_visit \
DB_SSLMODE=disable \
AUTO_MIGRATE=true \
go run ./cmd/api
```

后端默认监听 `http://localhost:8080`。看到 `Server is listening on 0.0.0.0:8080` 后保持该终端运行。

若上一步使用了 `15432`，这里将 `DB_PORT` 改为 `15432`。

### 3. 导入前端演示主数据

首次启动或重建数据库后，在项目根目录新开一个终端执行：

```bash
docker compose exec -T db psql \
  -U medical_visit \
  -d medical_visit \
  < scripts/seed_demo.sql
```

该脚本可重复执行。它会写入与前端选择项 UUID 一致的 MR、医院、科室、医生、产品和学术资料；医院坐标为 `31.2304, 121.4737`。

验证主数据：

```bash
docker compose exec db psql \
  -U medical_visit \
  -d medical_visit \
  -c "SELECT id, name FROM products ORDER BY name;"
```

### 4. 终端二：启动前端开发服务

在第二个长期运行的终端执行：

```bash
cd web
npm ci
npm run dev
```

前端默认运行在 `http://localhost:5173`，Vite 会将 `/api` 请求代理到本地后端 `http://localhost:8080`。保持该终端运行。

### 5. 确认三个组件均正常

```bash
# 数据库应显示 healthy
docker compose ps db

# 后端健康检查
curl -i http://localhost:8080/healthz

# 前端页面响应
curl -I http://localhost:5173
```

预期后端返回 HTTP `200` 和 `status: healthy`，前端返回 HTTP `200`。

### 6. 停止本地服务

前端和后端终端分别按 `Ctrl+C`，然后停止数据库：

```bash
docker compose stop db
```

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

### 后端测试

在项目根目录执行：

```bash
# 运行全部 Go 单元测试
go test ./...

# 运行竞态检测并生成覆盖率文件
go test -race -coverprofile=coverage.out ./...

# 查看覆盖率报告
go tool cover -html=coverage.out

# 静态检查
go vet ./...
```

### 前端构建测试

在 [`web`](web) 目录执行：

```bash
npm ci
npm run build
```

### 后端完整链路验证

先确认 [`scripts/seed_demo.sql`](scripts/seed_demo.sql) 已导入，然后在项目根目录执行以下步骤。所有操作必须使用创建接口返回的真实 UUID，不要使用 `demo-visit-001`。

#### 1. 健康检查

```bash
curl -s http://localhost:8080/healthz | python3 -m json.tool
```

#### 2. 创建拜访并保存真实 ID

```bash
PLAN_TIME=$(date -u -v+10M '+%Y-%m-%dT%H:%M:%SZ')

curl -sS -o /tmp/medical-visit-create.json \
  -w 'HTTP %{http_code}\n' \
  -X POST http://localhost:8080/api/v1/visits \
  -H 'Content-Type: application/json' \
  -d "{
    \"mr_id\": \"00000000-0000-0000-0000-000000000001\",
    \"hcp_id\": \"00000000-0000-0000-0000-000000000004\",
    \"hospital_id\": \"00000000-0000-0000-0000-000000000002\",
    \"department_id\": \"00000000-0000-0000-0000-000000000003\",
    \"product_id\": \"00000000-0000-0000-0000-000000000005\",
    \"planned_start_at\": \"${PLAN_TIME}\",
    \"plan_note\": \"完整链路验证：沟通最新临床证据\"
  }"

cat /tmp/medical-visit-create.json | python3 -m json.tool
VISIT_ID=$(python3 -c 'import json; print(json.load(open("/tmp/medical-visit-create.json"))["data"]["id"])')
echo "VISIT_ID=${VISIT_ID}"
```

预期 HTTP `201`，状态为 `PLANNED`。

#### 3. 签到

```bash
CHECK_IN_TIME=$(date -u '+%Y-%m-%dT%H:%M:%SZ')

curl -sS -X POST "http://localhost:8080/api/v1/visits/${VISIT_ID}/check-in" \
  -H 'Content-Type: application/json' \
  -d "{
    \"latitude\": 31.2304,
    \"longitude\": 121.4737,
    \"check_in_time\": \"${CHECK_IN_TIME}\"
  }" | python3 -m json.tool
```

预期状态为 `CHECKED_IN`，签到距离接近 0 米。

#### 4. 签退并验证合规结果

立即签退会因默认最短时长为 300 秒而得到 `NON_COMPLIANT` 和 `DURATION_TOO_SHORT`，这正好可以验证异常规则：

```bash
CHECK_OUT_TIME=$(date -u '+%Y-%m-%dT%H:%M:%SZ')

curl -sS -X POST "http://localhost:8080/api/v1/visits/${VISIT_ID}/check-out" \
  -H 'Content-Type: application/json' \
  -d "{
    \"latitude\": 31.2304,
    \"longitude\": 121.4737,
    \"check_out_time\": \"${CHECK_OUT_TIME}\"
  }" | python3 -m json.tool
```

预期状态为 `CHECKED_OUT`，合规状态为 `NON_COMPLIANT`。若要快速验证正常结果，可使用 `COMPLIANCE_MIN_DURATION_SECONDS=0` 重启本地后端后重新创建一条拜访。

#### 5. 提交 Call Report

```bash
curl -sS -X POST "http://localhost:8080/api/v1/visits/${VISIT_ID}/report" \
  -H 'Content-Type: application/json' \
  -d '{
    "conversation_summary": "介绍最新临床证据和标准用药路径",
    "doctor_feedback": "医生关注长期用药的安全性数据",
    "materials_distributed": true,
    "material_ids": ["00000000-0000-0000-0000-000000000007"],
    "additional_notes": "下次跟进真实世界研究数据"
  }' | python3 -m json.tool
```

预期返回 `SUCCESS`。报告保存后拜访状态会变为 `COMPLETED`。

#### 6. 查询详情与月度看板

```bash
curl -sS "http://localhost:8080/api/v1/visits/${VISIT_ID}" \
  | python3 -m json.tool

MONTH=$(date '+%Y-%m')
curl -sS "http://localhost:8080/api/v1/dashboard/visits/monthly?month=${MONTH}" \
  | python3 -m json.tool
```

预期详情状态为 `COMPLETED`，看板中“心宁平”的拜访次数增加。

### 前端完整链路验证

确保数据库、后端和前端均已运行，并已导入 [`scripts/seed_demo.sql`](scripts/seed_demo.sql)，然后执行：

1. 浏览器打开 `http://localhost:5173`。
2. 打开开发者工具的 Network 面板，筛选 `Fetch/XHR`，便于确认所有请求均返回成功。
3. 点击左侧“创建拜访”。不要从初始“拜访列表”或“异常样例”开始，因为其中的 `demo-visit-001`、`demo-visit-002` 只是前端展示数据，不存在于数据库。
4. 选择默认 MR、医生、医院、科室和“心宁平”，将计划时间设置为当前时间之后，填写拜访目标，点击“创建拜访”。
5. 确认页面提示“拜访计划已创建”，详情状态为“已计划”；Network 中 `POST /api/v1/visits` 应返回 HTTP `201`。
6. 在详情页保持默认签到坐标 `31.2304, 121.4737`，点击“现场签到”；状态应变为“已签到”，`POST .../check-in` 返回 HTTP `200`。
7. 点击“去签退”，保持默认坐标并确认签退；`POST .../check-out` 返回 HTTP `200`。立即签退时出现“异常/停留时长不足 5 分钟”是默认合规规则的预期结果，不代表链路失败。
8. 返回详情后点击“填写 Call Report”，填写“谈话要点”和可选反馈。若勾选派发资料，当前前端不会填写资料 UUID，因此建议浏览器链路验证时暂不勾选；然后提交报告。
9. 确认出现“Call Report 已保存”，Network 中 `POST .../report` 返回 HTTP `200`。
10. 点击“月度 Dashboard”，月份选择当前月份，确认产品拜访统计增加且页面没有“后端 Dashboard 暂未连接”提示。
11. 如需确认最终数据库状态，可执行：

```bash
docker compose exec db psql \
  -U medical_visit \
  -d medical_visit \
  -c "SELECT id, status, compliance_status, duration_seconds FROM visits ORDER BY created_at DESC LIMIT 5;"
```

### Docker 部署验证

如果使用 Docker 一键启动，执行：

```bash
docker compose ps
docker compose logs --tail=100 api
docker compose logs --tail=100 web
curl http://localhost/healthz
curl http://localhost/api/v1/ping
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
