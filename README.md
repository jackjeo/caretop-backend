# Caretop Backend

凯普软件官方网站后端服务 (Go + Gin)

## 技术栈

- **语言**: Go 1.22+
- **框架**: Gin + GORM
- **数据库**: PostgreSQL 15+
- **缓存**: Redis
- **认证**: JWT + bcrypt

## 项目结构

```
caretop-backend/
├── config/         # 配置文件
├── database/       # 数据库迁移
├── handlers/       # HTTP 处理器
├── middleware/     # 中间件
├── models/         # 数据模型
├── routes/        # 路由定义
├── utils/         # 工具函数
├── docs/          # 项目文档
├── main.go        # 入口文件
├── go.mod
├── Dockerfile
└── docker-compose.yml
```

## 快速开始

### 环境要求

- Go 1.22+
- PostgreSQL 15+
- Redis

### 本地运行

```bash
# 复制环境配置
cp .env.example .env
# 编辑 .env 填写数据库配置

# 运行服务
go run main.go
```

### Docker 部署

```bash
docker-compose up -d
```

## API 文档

- [产品需求规格书](./docs/Caretop_SRS.md)

## 许可证

MIT License
