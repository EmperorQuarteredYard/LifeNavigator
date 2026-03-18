> 一个融合资源管理、任务追踪、游戏化激励与 AI 助手的个人综合效率工具。

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)](https://golang.org)
[![Flutter](https://img.shields.io/badge/Flutter-3.x-02569B?logo=flutter)](https://flutter.dev)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

---

## 📖 项目概述

**LifeNavigator** 旨在帮助个人高效规划生活事项，追踪时间、金钱、精力三大资源的投入与产出，并通过游戏化机制提升执行力。  
它将**待办事项**视为未来的资源承诺，**已办事项**视为实际的资源消耗记录，让用户清晰看到“时间花在哪”、“钱流向哪”、“精力分配给谁”，从而更明智地规划下一步。

项目采用 **Golang 后端 + Flutter 客户端**，已实现核心数据模型与业务层（项目、任务、预算池、账户系统、事务支持），为上层功能奠定坚实基础。

---

## ✨ 愿景与特色

- **资源统一抽象**：时间、金钱、精力被抽象为可量化、可消耗、可恢复的资源，与任务深度关联。
- **Agent 认领任务**：内置多种 Agent 角色（提醒小助手、记账机器人等），可自动认领重复性任务，减轻用户负担。  
  *AI 负责智能分析（任务拆解、总结生成），Agent 负责落地执行，分工明确。*
- **游戏化激励**：经验值、等级、成就系统，随机奖励，让坚持变得有趣。
- **AI 增强**：集成 LLM（DeepSeek），提供智能拆解、每日总结、异常检测等能力。
- **规整的资源管理**：账户、预算、付款联动，自动更新余额，支持事务保证数据一致性。
- **权限设计**：类似 Linux 777 的权限分级，为未来多人协作铺路。
- **陪伴式桌宠**（远期）：LLM 驱动的桌宠，提供情感支持与随机激励。

---

## ✅ 功能状态

### 基础记录（已完成）
- [x] 任务管理（增删改查，状态、截止时间）
- [x] 项目组织（任务按项目分组，支持刷新间隔）
- [x] 项目级预算池（预算总额、已用金额，可关联账户）
- [x] 任务付款（自动更新预算与账户余额）
- [x] 多账户系统（现金、银行卡等，乐观锁并发控制）
- [x] 事务支持（保证多步写入原子性）
- [x] 任务依赖（前置任务约束）

### 任务追踪与 Agent（开发中）
- [ ] 任务预估/实际用时、预算对比报表
- [ ] Agent 框架：定时任务、规则引擎
- [ ] Agent 认领未分配任务（如过期提醒）
- [ ] 可配置的重复性任务自动处理

### 游戏化激励（规划中）
- [ ] 经验值、等级、成就系统
- [ ] 随机奖励（虚拟积分）
- [ ] 积分兑换现实奖品（未来结合 AutoGLM 自动化购买）

### AI 增强（规划中）
- [ ] 智能任务拆解
- [ ] 每日/每周总结生成
- [ ] 支出异常检测
- [ ] 情感支持聊天
- [ ] 预算分配建议、优先级排序、趣味挑战……

### 数据可视化与仪表盘
- [ ] 今日概览（待办、预算剩余、经验进度）
- [ ] 电池图展示资源占用
- [ ] 时间/金钱趋势图、分类饼图

### 通知提醒
- [ ] 任务截止提醒、预算超支预警、每日总结推送

### 扩展自动化（远期）
- [ ] AutoGLM 集成：奖励金累计后自动完成奖品下单（支付环节人工确认）

---

## 🛠️ 技术栈

| 层次     | 技术选型                          | 说明                                   |
|----------|-----------------------------------|----------------------------------------|
| 后端     | Golang + Gin                      | 高性能 RESTful API                     |
| 数据库   | PostgreSQL / SQLite（初期）        | 关系型数据库，支持复杂查询              |
| ORM      | GORM                              | 简化数据库操作，支持事务与软删除        |
| 客户端   | Flutter 3.x + Riverpod             | 跨平台（iOS/Android/桌面），状态管理    |
| 本地通知 | flutter_local_notifications       | 多平台本地通知                          |
| 图表     | fl_chart                          | 电池图、统计图                          |
| LLM 集成 | DeepSeek API / AutoGLM            | 免费、长上下文；自动化操作               |

---

## 🚀 快速开始

### 环境要求
- Go 1.21+
- Flutter 3.x
- PostgreSQL（可选，开发可用 SQLite）

### 后端启动
```bash
git clone https://github.com/EmperorQuarteredYard/LifeNavigator.git
cd lifenavigator/backend
cp config.example.yaml config.yaml   # 修改数据库等配置
go mod download
go run cmd/main.go
```

### 前端启动
```bash
cd ../frontend
flutter pub get
flutter run
```

> 详细部署文档请参阅 [docs/DEPLOY.md](docs/DEPLOY.md)（待完善）

---

## 📁 项目结构

### 后端（Golang）
```
cmd/
  main.go                     # 应用入口
internal/
  controller/                 # 控制器（HTTP 请求处理）
  service/                    # 业务逻辑层
  repository/                 # 数据访问层
  models/                     # 数据库模型
  router/                     # 路由注册
  database/                   # 数据库初始化
pkg/
  config/                     # 配置加载
  dto/                        # 数据传输对象
  errcode/                    # 错误码
  jwt/                        # JWT 工具
  response/                   # 统一响应
```

### 前端（Flutter）
```
lib/
  main.dart
  app/                        # 应用入口、路由
  models/                     # 数据模型（与后端 DTO 对应）
  views/                      # 页面（认证、仪表盘、项目、任务、统计、Agent）
  widgets/                    # 可复用组件（电池图、任务卡片等）
  services/                   # API 调用（dio）
  providers/                  # Riverpod 状态管理
  utils/                      # 工具函数
```

---

---

## 🤝 贡献指南

欢迎任何形式的贡献！
- 提交 [Issue](https://github.com/EmperorQuarteredYard/LifeNavigator/issues) 报告 bug 或提出新功能
- Fork 项目并提交 Pull Request
- 参与讨论，帮助改进设计

---

## 📄 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

---

## 📬 联系我们

- 作者：[EQuarterY](https://github.com/EmperorQuarteredYard)
- 项目主页：[https://github.com/EmperorQuarteredYard/LifeNavigator](https://github.com/EmperorQuarteredYard/LifeNavigator)

---

**LifeNavigator** —— 从个人需求出发，一步一个脚印，像经营 RPG 一样经营自己的人生。
# 如何开始
在/backend下新建一个docker-compose.yml，内容如下所示：
```yaml
version: '3.8'

services:
  web:
    build:
      context: .
    stdin_open: true   # 等价于 docker run -i
    tty: true          # 等价于 docker run -t
    ports:
      - "5083:5083"               # 映射宿主机端口到容器端口（根据你的应用端口调整）
    environment:
      - DB_HOST=mysql              # MySQL 服务名
      - DB_PORT=3306
      - DB_USER=root
      - DB_PASSWORD=CCEl8shKS1Byu0QFbSwy0   # 请修改为强密码
      - DB_NAME=MySQLDatabase               # 数据库名
      - JWT_ACCESS_SECRET=S4f7XnZh5zha9dLFbbjj1TpT77LBft8aDc4DOj34hluLaEZaG8pzERXKuZQQMqCTEZlLiZgtv_RAIUXvyzKFkjjhbA5Nr9yRV_JuibeIl9jChwzC7Of8FGg7HHXdiTTf #这里是随机生成的字符可以不用管
      - JWT_REFRESH_SECRET=yONJ7qBqz9mDGSufSQZs30bT5h2lumRBtaK_M8XD766NkSU73xK8uWeat1MfpUmWPTFo_Xrg7ybKHConVxYcyfOrOgboWdrhBdnYfBe6bUoxLmed9OBCKddymfgVU8t8 #这里是随机生成的字符可以不用管
      - JWT_ISSUER=LifeNavigatorDemo
      - JWT_ACCESS_TTL=7200000000000      #2小时
      - JWT_REFRESH_TTL=259200000000000   #3天
      - GLM_TOKEN=这里硬编码了GLM的Token
    depends_on:
      - mysql
    networks:
      - app-network

  mysql:
    image: mysql:9.6
#    ports:
#      - "3306:3306"                # 可选：暴露到宿主机方便调试
    environment:
      MYSQL_ROOT_PASSWORD: CCEl8shKS1Byu0QFbSwy0   # root 密码
      MYSQL_DATABASE: LifeNavigatorData                 # 数据库
      MYSQL_USER: user                       # 普通用户
      MYSQL_PASSWORD: xPBPH4UzbQ           # 普通用户密码
    volumes:
      - mysql-data:/var/lib/mysql            # 持久化数据
    networks:
      - app-network
    command:
      - --character-set-server=utf8mb4
      - --collation-server=utf8mb4_unicode_ci
    healthcheck:
      test: [ "CMD", "mysqladmin", "ping", "-h", "localhost", "-u", "root", "-p$$MYSQL_ROOT_PASSWORD" ]
      interval: 5s       # 每 5 秒检查一次
      timeout: 3s        # 每次检查超时 3 秒
      retries: 5         # 连续失败 5 次后标记为 unhealthy
      start_period: 10s  # 启动后等待 10 秒再开始检查（给 MySQL 初始化时间）

networks:
  app-network:
    driver: bridge

volumes:
  mysql-data:
```
在/backend下运行 
```bash
docker-compose up --build
```
即可