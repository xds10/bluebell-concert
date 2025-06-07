# 演唱会票务系统 API 文档

## 基础信息
- 基础路径: `/api/v1`
- 认证方式: JWT Token (Bearer Token)

## 1. 用户认证接口

### 1.1 用户注册
- **接口**: `POST /signup`
- **功能**: 新用户注册
- **请求体**: 
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```
- **响应**: 成功/失败信息

### 1.2 用户登录
- **接口**: `POST /login`
- **功能**: 用户登录获取 token
- **请求体**: 
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```
- **响应**: JWT Token

## 2. 演唱会管理接口

### 2.1 获取演唱会列表
- **接口**: `GET /concert-list`
- **功能**: 获取所有演唱会信息
- **权限**: 无需认证
- **响应**: 
  ```json
  {
    "data": [
      {
        "id": "int64",
        "title": "string",
        "venue_id": "int",
        "date": "string",
        "enddate": "string",
        "saletime": "string",
        "status": "int"
      }
    ]
  }
  ```

### 2.2 创建演唱会
- **接口**: `POST /create_concert`
- **功能**: 创建新的演唱会
- **权限**: 需要认证
- **请求体**: 
  ```json
  {
    "title": "string",
    "venue_id": "int",
    "date": "string",
    "enddate": "string",
    "saletime": "string",
    "status": "int"
  }
  ```
- **响应**: 成功/失败信息

### 2.3 获取演唱会详情
- **接口**: `GET /concert/:id`
- **功能**: 获取指定演唱会的详细信息
- **权限**: 需要认证
- **参数**: 
  - id: 演唱会ID (路径参数)
- **响应**: 演唱会详细信息

## 3. 票务接口

### 3.1 购买票
- **接口**: `POST /buy`
- **功能**: 购买演唱会票
- **权限**: 需要认证
- **请求体**: 
  ```json
  {
    "concert_id": "int64",
    "user_id": "int64",
    "seat_idx": {
      "section": "string"
    }
  }
  ```
- **响应**: 座位信息

### 3.2 支付订单
- **接口**: `POST /pay`
- **功能**: 支付已创建的订单
- **权限**: 需要认证
- **请求体**: 
  ```json
  {
    "id": "int64",
    "user_id": "int64",
    "concert_id": "int64",
    "seat_id": "int64",
    "price": "float64",
    "status": "int",
    "create_time": "string"
  }
  ```
- **响应**: 成功/失败信息

## 数据库表结构

### 1. 演唱会表 (concert)
```sql
CREATE TABLE `concert` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '演唱会ID',
  `title` varchar(100) DEFAULT NULL COMMENT '演唱会名称',
  `venue_id` int DEFAULT NULL COMMENT '场馆ID',
  `start_time` datetime DEFAULT NULL COMMENT '演出开始时间',
  `end_time` datetime DEFAULT NULL COMMENT '演出结束时间',
  `release_time` datetime DEFAULT NULL COMMENT '放票时间',
  `status` tinyint(1) DEFAULT '0' COMMENT '状态(0:未开始,1:进行中,2:已结束)',
  PRIMARY KEY (`id`),
  KEY `idx_start_time` (`start_time`),
  KEY `idx_status` (`status`)
)
```

### 2. 订单表 (orders)
```sql
CREATE TABLE `orders` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '订单ID',
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `concert_id` bigint unsigned NOT NULL COMMENT '演唱会ID',
  `seat_id` bigint NOT NULL COMMENT '座位ID',
  `price` decimal(10,2) NOT NULL COMMENT '订单金额',
  `status` int NOT NULL COMMENT '订单状态',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_concert_id` (`concert_id`),
  KEY `idx_status` (`status`)
)
```

## 状态码说明
- 0: 未开始售票
- 1: 售票中
- 2: 已结束

## 注意事项
1. 所有需要认证的接口都需要在请求头中携带 `Authorization: Bearer <token>`
2. 订单支付有5分钟的超时限制
3. 系统使用 Redis 进行座位锁定和并发控制
4. 系统每2小时会自动检查并更新未来演唱会的座位信息 