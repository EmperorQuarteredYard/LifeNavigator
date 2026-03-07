# LifeNavigator Client

Flutter 客户端项目 - 咫尺生涯

## 项目结构

```
lib/
├── main.dart              # 应用入口
├── app/
│   └── app.dart           # 应用配置、路由
├── models/                # 数据模型（与后端 DTO 对应）
│   ├── common.dart
│   ├── user.dart
│   ├── project.dart
│   ├── task.dart
│   └── account.dart
├── services/              # API 调用服务
│   ├── api_client.dart
│   ├── auth_service.dart
│   ├── project_service.dart
│   ├── task_service.dart
│   └── account_service.dart
├── providers/             # Riverpod 状态管理
│   ├── auth_provider.dart
│   ├── project_provider.dart
│   ├── task_provider.dart
│   └── account_provider.dart
├── views/                 # 页面视图
│   ├── auth/              # 登录/注册
│   ├── dashboard/         # 仪表盘
│   ├── projects/          # 项目管理
│   └── tasks/             # 任务管理
├── widgets/               # 可复用组件
│   ├── battery_gauge.dart
│   ├── task_card.dart
│   ├── project_card.dart
│   └── ...
└── utils/                 # 工具函数
    ├── date_utils.dart
    ├── currency_utils.dart
    ├── notification_utils.dart
    └── theme.dart
```

## 开发指南

### 环境要求
- Flutter SDK >= 3.0.0
- Dart SDK >= 3.0.0

### 安装依赖
```bash
flutter pub get
```

### 生成代码（freezed/json_serializable）
```bash
flutter pub run build_runner build --delete-conflicting-outputs
```

### 运行应用
```bash
flutter run
```

## 架构说明

### 分层架构
1. **Models**: 数据模型定义，使用 freezed 生成不可变类
2. **Services**: API 调用封装，使用 Dio 进行网络请求
3. **Providers**: Riverpod 状态管理，管理应用状态
4. **Views**: 页面视图，负责 UI 展示和用户交互
5. **Widgets**: 可复用组件，提供通用 UI 组件

### 状态管理
使用 Riverpod 进行状态管理，主要 Provider：
- `authStateProvider`: 用户认证状态
- `projectStateProvider`: 项目数据和状态
- `taskStateProvider`: 任务数据和状态
- `accountStateProvider`: 账户数据和状态

### 路由
使用 go_router 进行路由管理，支持：
- 声明式路由配置
- 路由守卫（认证检查）
- 深度链接

## 扩展指南

### 添加新功能模块
1. 在 `models/` 创建数据模型
2. 在 `services/` 创建 API 服务
3. 在 `providers/` 创建状态管理
4. 在 `views/` 创建视图页面
5. 在 `app/app.dart` 添加路由配置

### 添加新组件
在 `widgets/` 目录下创建新组件，遵循以下原则：
- 组件应该是可复用的
- 通过参数配置组件行为
- 使用 `const` 构造函数优化性能
