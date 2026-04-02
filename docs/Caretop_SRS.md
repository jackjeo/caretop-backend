# 凯普软件官方网站 · 产品需求规格书

| 版本 | 日期 | 说明 |
|------|------|------|
| V1.0 | 2026-04-01 | 初稿 |

---

## 1. 项目概述

### 1.1 项目背景

凯普软件（Caretop.com）旗下拥有三大产品：MindLink、HenryIway、Remote Desktop WPF。其中 MindLink 和 HenryIway 为核心产品，Remote Desktop WPF 为辅助小工具。为统一对外展示品牌形象，建立集产品展示、博客、论坛于一体的品牌官方网站。

### 1.2 项目目标

- 建设凯普软件统一品牌门户（Caretop.com）
- 展示 MindLink、HenryIway、Remote Desktop WPF 三大产品
- 提供博客系统，发布产品动态与技术文章
- 提供论坛系统，支持用户交流与客服支持
- 实现部分功能需注册后使用（发帖/评论/下载/工单）

### 1.3 技术选型

| 层级 | 技术 | 说明 |
|------|------|------|
| 后端 | Go + Gin | 轻量、内存占用极低、编译为单一二进制 |
| 前端 | Vue 3 + Nuxt 3 | SSR 模式，天然支持 SEO 与响应式 |
| 数据库 | PostgreSQL 15+ | 关系型数据库，稳定可靠 |
| 缓存 | Redis | Session 缓存 + 数据缓存 |
| 存储 | 本地文件系统 / OSS | 产品资料/头像/附件存储 |
| 部署 | Docker + Docker Compose | 一键部署，全容器化 |

### 1.4 部署架构

```
                    ┌──────────────┐
                    │   Nginx      │
                    │  (反向代理+SSL)│
                    └──────┬───────┘
                           │
         ┌─────────────────┼─────────────────┐
         │                 │                 │
    ┌────▼────┐      ┌────▼────┐       ┌────▼────┐
    │ Frontend│      │ Backend │       │   DB    │
    │(Nuxt 3) │◄────►│ (Go API)│◄────►│PostgreSQL│
    │ SSR:3000│      │ :8080   │       │  :5432  │
    └─────────┘      └────┬────┘       └─────────┘
                          │
                    ┌─────▼─────┐
                    │   Redis   │
                    │  :6379    │
                    └───────────┘
```

---

## 2. 信息架构

### 2.1 网站结构

```
Caretop.com
├── 首页 /
│   ├── 导航栏（顶部固定）
│   ├── Hero 区域（品牌标语 + 行动按钮）
│   ├── 产品展示区（MindLink / HenryIway / Remote Desktop）
│   ├── 最新博客文章（3篇）
│   ├── 论坛精选帖（3帖）
│   └── 页脚（版权 / 链接 / 社交媒体）
│
├── 产品 /products
│   ├── /products/mindlink     — MindLink 产品页
│   ├── /products/henryiway     — HenryIway 产品页
│   └── /products/remote-desktop — Remote Desktop 产品页
│
├── 博客 /blog
│   ├── 博客列表页 /blog
│   ├── 文章详情页 /blog/[slug]
│   └── 博客分类 /blog/category/[category]
│
├── 论坛 /forum
│   ├── 论坛首页 /forum
│   ├── 板块列表 /forum/[board]
│   ├── 帖子详情 /forum/[board]/[thread]
│   └── 发帖/回复（需登录）
│
├── 工单 /tickets（需登录）
│   ├── 我的工单 /tickets
│   └── 提交工单 /tickets/new
│
├── 用户
│   ├── 注册 /auth/register
│   ├── 登录 /auth/login
│   └── 个人中心 /user/profile（需登录）
│
└── 管理后台 /admin（管理员专用）
    ├── /admin/dashboard       — 数据概览
    ├── /admin/products         — 产品管理
    ├── /admin/blog             — 博客管理（写文章/编辑）
    ├── /admin/forum             — 论坛管理（板块/精华/删帖）
    ├── /admin/tickets           — 工单管理
    └── /admin/users             — 用户管理
```

---

## 3. 功能需求

### 3.1 产品展示

**3.1.1 首页产品区**
- 三个产品卡片展示：图标 + 名称 + 一句话简介
- 点击跳转至对应产品详情页

**3.1.2 产品详情页（每个产品独立页面）**
- 产品名称、标语、详细描述
- 产品截图/图片轮播
- 主要功能特性列表（图文并茂）
- 适用场景
- 下载/试用按钮（点击触发下载或跳转）
- 相关文档链接（需登录）

**3.1.3 Remote Desktop 辅助工具区**
- 提供小工具快速下载入口
- 不作为主推产品，但保留完整产品介绍

---

### 3.2 博客系统

**3.2.1 博客列表页**
- 分页展示文章列表（每页10篇）
- 每篇文章卡片：封面图 + 标题 + 摘要 + 作者 + 发布时间 + 阅读量
- 支持按分类筛选（技术文章 / 产品动态 / 行业洞察）

**3.2.2 文章详情页**
- 完整文章内容（支持 Markdown 渲染）
- 作者头像 + 名称 + 发布时间
- 阅读量统计
- 点赞功能（需登录）
- 评论区（需登录）
- 相关文章推荐（同一分类下的其他文章）

**3.2.3 文章管理（管理员）**
- 在线富文本编辑器（写/编辑/发布/下架文章）
- 分类管理（增删分类）
- 设置封面图

---

### 3.3 论坛系统

**3.3.1 论坛板块**
预设板块：
| 板块 | 说明 |
|------|------|
| MindLink 产品交流 | MindLink 使用讨论、技巧分享 |
| HenryIway 产品交流 | HenryIway 使用讨论 |
| Remote Desktop 交流 | 远程控制工具使用讨论 |
| 经验分享 | 用户原创教程、案例 |
| 官方公告 | 管理员发布产品更新、公告 |

**3.3.2 帖子列表页**
- 分页展示（每页20帖）
- 每帖：标题 + 作者 + 发布时间 + 回复数 + 浏览数
- 支持精华帖置顶
- 支持按最新/最热/精华筛选
- 搜索功能（按标题关键词）

**3.3.3 帖子详情页**
- 帖子正文（Markdown 渲染）
- 作者信息 + 发布时间
- 浏览量
- 楼层式回复（嵌套结构，最多3层）
- 点赞帖子功能（需登录）
- 收藏帖子功能（需登录）

**3.3.4 发帖/回复（需登录）**
- 富文本编辑器（支持 Markdown + 图片上传）
- 发帖时选择板块
- @提及功能（回复中 @用户名）
- 编辑/删除自己的帖子

**3.3.5 论坛管理（管理员）**
- 创建/编辑/删除板块
- 精华/取消精华帖
- 置顶/取消置顶
- 删除/编辑任意帖子和回复
- 用户禁言（临时/永久）

---

### 3.4 用户体系

**3.4.1 注册**
- 邮箱 + 用户名 + 密码
- 邮箱激活验证（发送激活链接）
- 密码强度校验（8位以上，含大小写+数字）

**3.4.2 登录**
- 用户名/邮箱 + 密码
- 记住登录状态（7天有效）
- 登录失败超过5次，账号锁定15分钟

**3.4.3 个人中心**
- 基本信息：头像、用户名、邮箱、个人简介
- 修改密码
- 我发布的帖子
- 我收藏的帖子
- 我的工单记录

**3.4.4 用户角色**
| 角色 | 权限 |
|------|------|
| 游客 | 浏览产品、博客、论坛（部分功能受限） |
| 注册用户 | 发帖/回复/评论/收藏/提交工单/下载 |
| 版主 | 精华/置顶/删帖/禁言本板块用户 |
| 管理员 | 全站管理权限 |

---

### 3.5 工单系统（需登录）

**3.5.1 提交工单**
- 工单类型：Bug报告 / 功能建议 / 使用咨询 / 商务合作
- 标题 + 详细描述（支持附件上传）
- 提交后生成工单编号

**3.5.2 我的工单列表**
- 显示所有我提交的工单
- 状态筛选：待处理 / 处理中 / 已解决 / 已关闭

**3.5.3 工单详情**
- 工单信息 + 附件
- 对话式回复记录
- 管理员/客服处理状态变更

**3.5.4 工单管理（管理员）**
- 查看所有工单
- 分配工单状态
- 回复工单

---

### 3.6 管理后台

**3.6.1 数据概览**
- 注册用户数（总/日/周）
- 帖子数 / 评论数 / 工单数
- 博客文章数 / 阅读量
- 热门产品 / 热门文章

**3.6.2 产品管理**
- 编辑三个产品页面的所有内容
- 上传产品图片/截图
- 设置下载链接

**3.6.3 博客管理**
- 发布/编辑/删除文章
- 管理文章分类
- 查看文章阅读/点赞数据

**3.6.4 论坛管理**
- 板块 CRUD
- 帖子管理（精华/置顶/删除）
- 用户管理（角色调整/禁言）

**3.6.5 工单管理**
- 处理工单、回复用户

---

## 4. 数据库设计（ER 概要）

```
users
├── id (PK, UUID)
├── username
├── email
├── password_hash
├── role (enum: user/moderator/admin)
├── avatar_url
├── bio
├── is_active
├── created_at / updated_at

products
├── id (PK)
├── slug (唯一)
├── name
├── tagline
├── description (markdown)
├── features (jsonb)
├── screenshots (jsonb, 数组)
├── downloads (jsonb, [{name, url}])
├── is_published
├── created_at / updated_at

blog_posts
├── id (PK, UUID)
├── title
├── slug (唯一)
├── summary
├── content (markdown)
├── cover_image_url
├── category (enum: tech/product/industry)
├── author_id (FK → users)
├── view_count
├── like_count
├── is_published
├── published_at
├── created_at / updated_at

blog_comments
├── id (PK)
├── post_id (FK → blog_posts)
├── user_id (FK → users)
├── content
├── created_at

forum_boards
├── id (PK)
├── name
├── slug (唯一)
├── description
├── sort_order

forum_threads
├── id (PK, UUID)
├── board_id (FK → forum_boards)
├── user_id (FK → users)
├── title
├── content (markdown)
├── view_count
├── is_pinned
├── is_essential (精华)
├── is_locked
├── created_at / updated_at

forum_posts (楼层回复)
├── id (PK)
├── thread_id (FK → forum_threads)
├── user_id (FK → users)
├── parent_id (FK → forum_posts, nullable, 嵌套回复)
├── content
├── like_count
├── created_at / updated_at

forum_likes
├── user_id (PK, FK)
├── thread_id (PK, FK)
├── created_at

forum_collections
├── user_id (PK, FK)
├── thread_id (PK, FK)
├── created_at

tickets
├── id (PK, UUID)
├── user_id (FK → users)
├── type (enum: bug/feature/consult/business)
├── title
├── content
├── status (enum: pending/processing/resolved/closed)
├── created_at / updated_at

ticket_replies
├── id (PK)
├── ticket_id (FK → tickets)
├── user_id (FK → users, 回复人)
├── content
├── created_at

attachments
├── id (PK, UUID)
├── user_id (FK → users)
├── ref_type (enum: ticket/post)
├── ref_id (FK)
├── filename
├── file_url
├── file_size
├── created_at
```

---

## 5. API 设计（RESTful）

### 5.1 认证
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/auth/register | 注册 |
| POST | /api/v1/auth/login | 登录 |
| POST | /api/v1/auth/logout | 登出 |
| GET | /api/v1/auth/me | 当前用户信息 |

### 5.2 产品
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/products | 列表 |
| GET | /api/v1/products/:slug | 详情 |

### 5.3 博客
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/blog | 列表（分页） |
| GET | /api/v1/blog/:slug | 详情 |
| POST | /api/v1/blog | 创建文章（管理员） |
| PUT | /api/v1/blog/:id | 编辑文章（管理员） |
| DELETE | /api/v1/blog/:id | 删除文章（管理员） |
| POST | /api/v1/blog/:slug/like | 点赞（需登录） |

### 5.4 论坛
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/forum/boards | 板块列表 |
| GET | /api/v1/forum/boards/:slug/threads | 帖子列表 |
| GET | /api/v1/forum/threads/:id | 帖子详情 |
| POST | /api/v1/forum/threads | 发帖（需登录） |
| PUT | /api/v1/forum/threads/:id | 编辑（需登录，作者或管理员） |
| DELETE | /api/v1/forum/threads/:id | 删除（需登录，作者或管理员） |
| POST | /api/v1/forum/threads/:id/reply | 回复（需登录） |
| POST | /api/v1/forum/threads/:id/like | 点赞（需登录） |
| POST | /api/v1/forum/threads/:id/collect | 收藏（需登录） |

### 5.5 工单
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/tickets | 我的工单列表（需登录） |
| POST | /api/v1/tickets | 提交工单（需登录） |
| GET | /api/v1/tickets/:id | 工单详情 |
| POST | /api/v1/tickets/:id/reply | 回复工单（需登录） |

### 5.6 管理后台
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/admin/stats | 数据统计 |
| CRUD | /api/v1/admin/products | 产品管理 |
| CRUD | /api/v1/admin/blog | 博客管理 |
| CRUD | /api/v1/admin/forum/boards | 板块管理 |
| CRUD | /api/v1/admin/forum/threads | 帖子管理 |
| CRUD | /api/v1/admin/users | 用户管理 |
| CRUD | /api/v1/admin/tickets | 工单管理 |

---

## 6. 前端页面清单

| 页面 | 路径 | 模板 |
|------|------|------|
| 首页 | / | 首页 |
| MindLink 产品页 | /products/mindlink | 产品详情 |
| HenryIway 产品页 | /products/henryiway | 产品详情 |
| Remote Desktop 页 | /products/remote-desktop | 产品详情 |
| 博客列表 | /blog | 博客列表 |
| 博客文章 | /blog/[slug] | 文章详情 |
| 论坛首页 | /forum | 论坛首页 |
| 板块帖子列表 | /forum/[board] | 帖子列表 |
| 帖子详情 | /forum/[board]/[thread] | 帖子详情 |
| 发帖页 | /forum/new | 发帖（需登录） |
| 工单列表 | /tickets | 工单列表（需登录） |
| 提交工单 | /tickets/new | 提交工单（需登录） |
| 注册 | /auth/register | 注册 |
| 登录 | /auth/login | 登录 |
| 个人中心 | /user/profile | 个人中心（需登录） |
| 管理后台 | /admin/* | 管理后台（管理员） |

---

## 7. 非功能需求

### 7.1 性能
- 页面首屏加载 < 2秒（国内网络）
- API 响应时间 < 200ms（P95）
- 支持 1000 并发用户

### 7.2 安全
- 所有用户密码 bcrypt/SHA256 哈希存储
- JWT Token 认证，7天有效期
- SQL 注入/XSS/CSRF 防护
- 文件上传类型/大小校验
- 敏感操作（日志记录）

### 7.3 SEO
- 全站 SSR 渲染，支持搜索引擎收录
- Open Graph / Twitter Card 元数据
- Sitemap 自动生成
- 文章/帖子支持自定义 SEO 描述

### 7.4 国际化（预留）
- 当前版本仅支持中文
- 设计上预留 i18n 扩展能力

---

## 8. 项目交付物

1. **产品规格书**（本文档）
2. **技术设计文档**（数据库详计 + API 文档）
3. **前端源码**（Vue 3 + Nuxt 3）
4. **后端源码**（Go + Gin）
5. **Docker 部署文件**（Dockerfile + docker-compose.yml）
6. **数据库迁移脚本**
7. **测试报告**
8. **部署运维手册**

---

*版权所有 凯普软件 | https://www.katop.com*
