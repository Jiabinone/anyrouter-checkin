# AnyRouter 签到系统

AnyRouter 签到中台管理系统，采用前后端分离架构：
- 后端：Go + Gin + GORM + SQLite
- 前端：Vue 3 + TypeScript + shadcn-vue + TailwindCSS

提供账号管理、定时签到、推送通知与日志统计能力，支持单镜像一键部署。

## 目录结构

```
project-root/
├── backend/              # Go 后端服务
│   ├── cmd/server/       # 程序入口
│   ├── internal/         # 内部业务逻辑
│   │   ├── handler/      # HTTP 处理器
│   │   ├── service/      # 业务逻辑层
│   │   ├── repository/   # 数据访问层
│   │   └── model/        # 数据模型
│   └── pkg/              # 可复用包
├── frontend/             # Vue 前端
│   └── src/
│       ├── api/          # API 请求
│       ├── components/   # 组件
│       ├── pages/        # 页面
│       └── stores/       # Pinia 状态
├── docker/               # Nginx / Supervisor 配置
├── scripts/              # 构建/发布脚本
└── docs/                 # 技术文档
```

## 开发环境

- Go 1.23+
- Node.js 20+
- SQLite 3

## 本地开发

后端：

```bash
cd backend
make dev
```

前端：

```bash
cd frontend
npm ci
npm run dev
```

访问地址：
- 前端：`http://localhost:5173`
- 后端：`http://localhost:8080`

## 配置说明

后端配置文件：`backend/config.yaml`（需从示例复制）

```bash
cp backend/config.example.yaml backend/config.yaml
```

示例结构：

```yaml
server:
  port: 8080
  mode: debug

database:
  path: ./data/app.db

jwt:
  secret: <32字节Base64>
  expire: 24h

aes:
  key: <32字符十六进制>
```

## Docker 单镜像运行

```bash
docker run -d \
  --name anyrouter-checkin \
  -p 5173:80 \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  fjiabinc/anyrouter-checkin:latest
```

使用 Compose：

```bash
cp backend/config.example.yaml backend/config.yaml
docker compose up -d
```

## 镜像构建与推送

使用 buildx 多架构构建并推送（需要 Docker Buildx）：

```bash
docker buildx build --platform "linux/amd64,linux/arm64" \
  -t fjiabinc/anyrouter-checkin:latest \
  --push .
```

## CI（GitHub Actions）

- 使用 GitHub 官方 Runner
- 需要配置 Secrets：`DOCKERHUB_USERNAME`、`DOCKERHUB_TOKEN`

## 常用命令

后端：

```bash
cd backend
make lint
make build
make gen-docs
```

前端：

```bash
cd frontend
npm run lint
npm run type-check
npm run build
```

## License

如需开源许可，请补充 LICENSE 文件。
