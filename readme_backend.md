# 蓝铃铛演唱会售票系统后端文档

## 项目概述

蓝铃铛(Bluebell)是一个高并发演唱会售票系统，专注于解决大规模用户同时抢购演唱会门票的场景。本系统采用Go语言开发，基于Gin框架构建Web服务，使用MySQL作为持久化存储，Redis作为缓存和高并发处理的支持。

## 技术栈

- **语言**: Go 
- **Web框架**: Gin
- **数据库**: MySQL
- **缓存**: Redis
- **日志**: Zap
- **认证**: JWT
- **ID生成**: Snowflake算法

## 系统架构

系统采用经典的分层架构设计:

```
┌─────────────┐
│    Routes   │ ← 路由层，定义API路径和HTTP方法
└─────┬───────┘
      │
┌─────▼───────┐
│ Controllers │ ← 控制器层，处理HTTP请求和响应
└─────┬───────┘
      │
┌─────▼───────┐
│    Logic    │ ← 业务逻辑层，实现核心业务逻辑
└─────┬───────┘
      │
┌─────▼───────┐
│     DAO     │ ← 数据访问层，处理数据库操作
└─────┬───────┘
      │
┌─────▼───────┐
│   Models    │ ← 数据模型层，定义数据结构
└─────────────┘
```

## 目录结构

```
bluebell/
├── conf/                  # 配置文件目录
│   └── config.yaml        # 主配置文件
├── controller/            # 控制器层，处理HTTP请求
│   ├── code.go            # 状态码定义
│   ├── concert.go         # 演唱会相关控制器
│   ├── response.go        # 响应处理
│   ├── ticket.go          # 票务相关控制器
│   ├── user.go            # 用户相关控制器
│   └── validator.go       # 请求验证器
├── dao/                   # 数据访问层
│   ├── mysql/             # MySQL数据访问
│   │   ├── concert.go     # 演唱会数据操作
│   │   ├── error_code.go  # 错误码定义
│   │   ├── mysql.go       # MySQL连接管理
│   │   ├── order.go       # 订单数据操作
│   │   ├── seat.go        # 座位数据操作
│   │   ├── ticket.go      # 票务数据操作
│   │   └── user.go        # 用户数据操作
│   └── redis/             # Redis数据访问
├── logger/                # 日志处理
├── logic/                 # 业务逻辑层
│   ├── buy.go             # 购票逻辑
│   ├── concert.go         # 演唱会逻辑
│   ├── order.go           # 订单逻辑
│   ├── ticker.go          # 定时任务逻辑
│   └── user.go            # 用户逻辑
├── middlewares/           # 中间件
│   └── auth.go            # JWT认证中间件
├── models/                # 数据模型
│   ├── concert.go         # 演唱会模型
│   ├── create_table.sql   # 数据库表创建SQL
│   ├── order.go           # 订单模型
│   ├── params.go          # 请求参数模型
│   ├── seat.go            # 座位模型
│   ├── ticket.go          # 票务模型
│   └── user.go            # 用户模型
├── pkg/                   # 工具包
│   └── snowflake/         # ID生成器
├── routes/                # 路由定义
│   └── routes.go          # 路由配置
├── settings/              # 配置加载
├── go.mod                 # Go模块定义
├── go.sum                 # 依赖版本锁定
└── main.go                # 应用入口
```

## 核心功能

### 1. 用户管理

- **注册**: 用户通过提供用户名、密码等信息注册账号
- **登录**: 用户登录后获取JWT令牌，用于后续请求认证

### 2. 演唱会管理

- **演唱会列表**: 获取所有可用演唱会信息
- **演唱会详情**: 获取特定演唱会的详细信息
- **创建演唱会**: 管理员创建新的演唱会信息

### 3. 票务管理

- **座位查询**: 查询演唱会可用座位信息
- **购票**: 用户选择座位并购买票
- **订单管理**: 创建、查询和取消订单
- **支付处理**: 处理订单支付流程

## API接口文档

### 用户相关

#### 注册

- **URL**: `/api/v1/signup`
- **方法**: POST
- **请求体**:
  ```json
  {
    "username": "string",
    "password": "string",
    "re_password": "string",
    "email": "string"
  }
  ```
- **响应**:
  ```json
  {
    "code": 1000,
    "message": "success",
    "data": null
  }
  ```

#### 登录

- **URL**: `/api/v1/login`
- **方法**: POST
- **请求体**:
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```
- **响应**:
  ```json
  {
    "code": 1000,
    "message": "success",
    "data": {
      "token": "string",
      "user_id": "int64"
    }
  }
  ```

### 演唱会相关

#### 获取演唱会列表

- **URL**: `/api/v1/concert-list`
- **方法**: GET
- **响应**:
  ```json
  {
    "code": 1000,
    "message": "success",
    "data": [
      {
        "id": "int64",
        "name": "string",
        "artist": "string",
        "venue": "string",
        "time": "string",
        "price": "float",
        "status": "int"
      }
    ]
  }
  ```

#### 获取演唱会详情

- **URL**: `/api/v1/concert/:id`
- **方法**: GET
- **响应**:
  ```json
  {
    "code": 1000,
    "message": "success",
    "data": {
      "id": "int64",
      "name": "string",
      "artist": "string",
      "venue": "string",
      "time": "string",
      "price": "float",
      "description": "string",
      "image": "string",
      "status": "int",
      "total_seats": "int",
      "available_seats": "int"
    }
  }
  ```

#### 创建演唱会

- **URL**: `/api/v1/create_concert`
- **方法**: POST
- **请求头**: Authorization: Bearer {token}
- **请求体**:
  ```json
  {
    "name": "string",
    "artist": "string",
    "venue": "string",
    "time": "string",
    "price": "float",
    "description": "string",
    "image": "string",
    "total_seats": "int"
  }
  ```
- **响应**:
  ```json
  {
    "code": 1000,
    "message": "success",
    "data": {
      "concert_id": "int64"
    }
  }
  ```

#### 获取演唱会座位信息

- **URL**: `/api/v1/concert/:id/seats`
- **方法**: GET
- **响应**:
  ```json
  {
    "code": 1000,
    "message": "success",
    "data": [
      {
        "id": "int64",
        "row": "int",
        "column": "int",
        "status": "int"
      }
    ]
  }
  ```

### 票务相关

#### 购买票

- **URL**: `/api/v1/buy`
- **方法**: POST
- **请求头**: Authorization: Bearer {token}
- **请求体**:
  ```json
  {
    "concert_id": "int64",
    "seat_ids": ["int64"]
  }
  ```
- **响应**:
  ```json
  {
    "code": 1000,
    "message": "success",
    "data": {
      "order_id": "int64"
    }
  }
  ```

#### 支付订单

- **URL**: `/api/v1/pay`
- **方法**: POST
- **请求头**: Authorization: Bearer {token}
- **请求体**:
  ```json
  {
    "order_id": "int64",
    "payment_method": "string"
  }
  ```
- **响应**:
  ```json
  {
    "code": 1000,
    "message": "success",
    "data": {
      "payment_id": "string"
    }
  }
  ```

#### 获取订单列表

- **URL**: `/api/v1/order-list`
- **方法**: POST
- **请求头**: Authorization: Bearer {token}
- **请求体**:
  ```json
  {
    "page": "int",
    "size": "int"
  }
  ```
- **响应**:
  ```json
  {
    "code": 1000,
    "message": "success",
    "data": {
      "orders": [
        {
          "id": "int64",
          "concert_id": "int64",
          "concert_name": "string",
          "user_id": "int64",
          "status": "int",
          "total_price": "float",
          "created_at": "string"
        }
      ],
      "total": "int"
    }
  }
  ```

#### 获取订单详情

- **URL**: `/api/v1/order/:id`
- **方法**: POST
- **请求头**: Authorization: Bearer {token}
- **响应**:
  ```json
  {
    "code": 1000,
    "message": "success",
    "data": {
      "id": "int64",
      "concert_id": "int64",
      "concert_name": "string",
      "user_id": "int64",
      "status": "int",
      "total_price": "float",
      "created_at": "string",
      "tickets": [
        {
          "id": "int64",
          "seat_id": "int64",
          "row": "int",
          "column": "int"
        }
      ]
    }
  }
  ```

#### 取消订单

- **URL**: `/api/v1/cancel-order`
- **方法**: POST
- **请求头**: Authorization: Bearer {token}
- **请求体**:
  ```json
  {
    "order_id": "int64"
  }
  ```
- **响应**:
  ```json
  {
    "code": 1000,
    "message": "success",
    "data": null
  }
  ```

## 高并发处理策略

### 1. Redis缓存

系统使用Redis缓存演唱会座位信息，减轻数据库压力，提高响应速度。

### 2. 分布式锁

在购票过程中使用Redis实现分布式锁，防止超卖和并发冲突。

### 3. 异步处理

订单创建后，使用异步方式处理支付和票务生成，提高系统吞吐量。

### 4. 定时任务

系统包含定时任务处理未支付订单的自动取消，释放座位资源。

## 数据库设计

### 用户表(users)

```sql
CREATE TABLE `users` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) NOT NULL,
  `username` varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
  `password` varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
  `email` varchar(64) COLLATE utf8mb4_general_ci,
  `gender` tinyint(4) NOT NULL DEFAULT '0',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`),
  UNIQUE KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
```

### 演唱会表(concerts)

```sql
CREATE TABLE `concerts` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `concert_id` bigint(20) NOT NULL,
  `name` varchar(128) NOT NULL,
  `artist` varchar(64) NOT NULL,
  `venue` varchar(128) NOT NULL,
  `time` datetime NOT NULL,
  `price` decimal(10,2) NOT NULL,
  `description` text,
  `image` varchar(255),
  `status` tinyint(4) NOT NULL DEFAULT '0',
  `total_seats` int(11) NOT NULL,
  `available_seats` int(11) NOT NULL,
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_concert_id` (`concert_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
```

### 座位表(seats)

```sql
CREATE TABLE `seats` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `seat_id` bigint(20) NOT NULL,
  `concert_id` bigint(20) NOT NULL,
  `row` int(11) NOT NULL,
  `column` int(11) NOT NULL,
  `status` tinyint(4) NOT NULL DEFAULT '0',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_seat_id` (`seat_id`),
  KEY `idx_concert_id` (`concert_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
```

### 订单表(orders)

```sql
CREATE TABLE `orders` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `order_id` bigint(20) NOT NULL,
  `user_id` bigint(20) NOT NULL,
  `concert_id` bigint(20) NOT NULL,
  `status` tinyint(4) NOT NULL DEFAULT '0',
  `total_price` decimal(10,2) NOT NULL,
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `expire_time` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_order_id` (`order_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_concert_id` (`concert_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
```

### 票表(tickets)

```sql
CREATE TABLE `tickets` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `ticket_id` bigint(20) NOT NULL,
  `order_id` bigint(20) NOT NULL,
  `concert_id` bigint(20) NOT NULL,
  `seat_id` bigint(20) NOT NULL,
  `status` tinyint(4) NOT NULL DEFAULT '0',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_ticket_id` (`ticket_id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_seat_id` (`seat_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
```

## 配置说明

系统配置文件位于 `conf/config.yaml`，主要配置项包括：

```yaml
name: "bluebell"
mode: "dev"  # 运行模式，可选值：dev, release
port: 8081   # 服务端口
version: "v0.0.1"
start_time: "2020-07-01"  # 用于Snowflake ID生成
machine_id: 1  # 机器ID，用于Snowflake ID生成

auth:
  jwt_expire: 8760  # JWT令牌过期时间（小时）

log:
  level: "debug"  # 日志级别
  filename: "web_app.log"  # 日志文件名
  max_size: 200  # 单个日志文件最大大小（MB）
  max_age: 30  # 日志保留天数
  max_backups: 7  # 保留的旧日志文件数量

mysql:
  host: "127.0.0.1"  # MySQL主机地址
  port: 8086  # MySQL端口
  user: "root"  # MySQL用户名
  password: "123456"  # MySQL密码
  dbname: "concert-ticket-sales"  # 数据库名
  max_open_conns: 200  # 最大打开连接数
  max_idle_conns: 50  # 最大空闲连接数

redis:
  host: "127.0.0.1"  # Redis主机地址
  port: 6380  # Redis端口
  password: ""  # Redis密码
  db: 0  # Redis数据库索引
  pool_size: 100  # 连接池大小
```

## 启动流程

系统启动流程如下：

1. 加载配置文件
2. 初始化日志系统
3. 连接MySQL数据库
4. 连接Redis
5. 初始化Snowflake ID生成器
6. 初始化验证器翻译器
7. 注册路由
8. 启动HTTP服务
9. 启动定时任务

## 部署指南

### 环境要求

- Go 1.16+
- MySQL 5.7+
- Redis 5.0+

### 编译

```bash
go build -o bluebell main.go
```

### 运行

```bash
./bluebell
```

### 使用Docker部署

1. 构建Docker镜像

```bash
docker build -t bluebell:latest .
```

2. 运行容器

```bash
docker run -d -p 8081:8081 --name bluebell bluebell:latest
```

## 性能优化

1. **连接池管理**：合理配置MySQL和Redis连接池大小，避免频繁建立连接
2. **缓存策略**：对热点数据如演唱会信息、座位状态等进行缓存
3. **异步处理**：将非关键流程如日志记录、邮件通知等异步处理
4. **分布式锁**：使用Redis实现分布式锁，保证数据一致性
5. **批量操作**：批量处理数据库操作，减少数据库交互次数

## 错误处理

系统定义了统一的错误码和响应格式，便于前端处理和展示错误信息。主要错误类型包括：

- 参数错误
- 用户认证错误
- 资源不存在
- 权限错误
- 业务逻辑错误
- 系统内部错误

## 安全措施

1. **JWT认证**：使用JWT进行用户身份验证
2. **密码加密**：用户密码使用bcrypt算法加密存储
3. **输入验证**：对所有用户输入进行严格验证，防止注入攻击
4. **CORS配置**：合理配置跨域资源共享策略
5. **限流措施**：对API接口实施限流，防止恶意请求

## 日志系统

系统使用Zap日志库，支持多级别日志记录和日志轮转。日志配置可在config.yaml中调整。

## 测试

### 单元测试

```bash
go test -v ./...
```

### 性能测试

使用wrk工具进行性能测试：

```bash
wrk -t12 -c400 -d30s http://localhost:8081/api/v1/concert-list
```

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开Pull Request

## 许可证

MIT 