# Galaxy Empire Manager

游戏管理后台系统，提供用户管理、玩家信息查询、道具发放等功能。

## 功能特性

- 用户管理（创建、禁用、重置密码）
- 玩家信息查询
- 道具发放
- 权限控制
- 操作日志记录

## 技术栈

- Go 1.20+
- MongoDB 7.0+
- Gin Web Framework
- JWT 身份验证
- Logrus 日志系统

## 快速开始

### 前置要求
- Go 1.20 或更新版本
- MongoDB 7.0+

### 安装步骤
1. 克隆项目
```bash
git clone https://github.com/your-username/galaxy-empire-manager.git
cd galaxy-empire-manager
```

2. 安装依赖
```bash
go mod tidy
```

3. 初始化数据库
```bash
cd script
init_data.bat
```

4. 启动服务
```bash
start.bat
```

## 默认管理员账号
- 用户名：admin
- 密码：admin888

## 项目结构
```
├── config/         # 配置管理
├── handlers/      # 请求处理器
├── middleware/    # 中间件
├── models/        # 数据模型
├── script/        # 数据库脚本
├── templates/     # 前端模板
├── utils/         # 工具函数
│
├── main.go        # 入口文件
├── go.mod         # Go 模块定义
├── README.md      # 项目文档
└── .gitignore     # Git 忽略规则
```

## 日志配置
- 日志文件按天分割
- 信息日志：logs/info.log
- 错误日志：logs/error.log
- 自动保留最近7天日志

## 许可证
[MIT License](LICENSE)
