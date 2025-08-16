# OAuth2 服务端

基于 go-zero 框架和 go-oauth2 包实现的 OAuth2 授权码模式服务端。

## 功能特性

1. **客户端注册**: 提供内部API注册OAuth2客户端
2. **授权码模式**: 支持标准的OAuth2授权码流程
3. **JWT Token**: 使用JWT生成访问令牌
4. **权限范围**: 支持 userid 和 profile 两种权限范围
5. **自动授权**: 支持特定客户端自动授权，无需用户确认
6. **Redis存储**: 使用Redis存储授权码和访问令牌
7. **MySQL存储**: 使用MySQL存储客户端信息和授权记录
8. **官方go-oauth2包**: 使用成熟的OAuth2库实现

## 项目结构

```
oauth2-server/
├── etc/                    # 配置文件
│   └── oauth2-api.yaml
├── internal/               # 内部代码
│   ├── config/            # 配置结构
│   ├── model/             # 数据模型
│   └── util/              # 工具函数
├── static/                # 静态文件
│   ├── login.html         # 登录页面
│   └── auth.html          # 授权页面
├── scripts/               # 脚本文件
│   └── init.sql          # 数据库初始化脚本
├── test/                  # 测试文件
│   └── test_oauth2.go    # 测试脚本
├── docs/                  # 文档
│   ├── usage.md          # 使用说明
│   └── summary.md        # 项目总结
├── go.mod                 # Go模块文件
├── oauth2.go             # 主程序入口
├── Makefile              # 构建脚本
└── README.md             # 项目说明
```

## 快速开始

### 1. 环境准备

- Go 1.21+
- MySQL 8.0+
- Redis 6.0+

### 2. 数据库初始化

```bash
# 执行数据库初始化脚本
mysql -u root -p < scripts/init.sql
```

### 3. 配置修改

修改 `etc/oauth2-api.yaml` 配置文件：

```yaml
Name: oauth2-api
Host: 0.0.0.0
Port: 8080

MySQL:
  DataSource: root:your_password@tcp(localhost:3306)/oauth2?charset=utf8mb4&parseTime=True&loc=Local

Redis:
  Host: localhost:6379
  Type: node
  Pass: ""

Auth:
  AccessSecret: your-jwt-secret-key-here
  AccessExpire: 7200 # 2小时

# 不需要用户授权的客户端ID列表
AutoApproveClients:
  - "trusted_client_001"
  - "trusted_client_002"
```

### 4. 启动服务

```bash
# 安装依赖
go mod tidy

# 启动服务
go run oauth2.go
```

## API 接口

### 1. 客户端注册

**POST** `/api/client/register`

请求体：
```json
{
  "name": "应用名称",
  "redirect_url": "http://localhost:3000/callback",
  "grant_type": "authorization_code",
  "scope": "userid profile"
}
```

响应：
```json
{
  "client_id": "client_abc123",
  "client_secret": "secret_xyz789"
}
```

### 2. 授权请求

**GET** `/oauth/authorize`

参数：
- `client_id`: 客户端ID
- `response_type`: 固定为 "code"
- `redirect_uri`: 重定向URI
- `scope`: 权限范围
- `state`: 状态参数（可选）

### 3. 获取访问令牌

**POST** `/oauth/token`

参数：
- `grant_type`: 固定为 "authorization_code"
- `code`: 授权码
- `redirect_uri`: 重定向URI
- `client_id`: 客户端ID
- `client_secret`: 客户端密钥

响应：
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 7200,
  "refresh_token": "refresh_token_123",
  "scope": "userid profile"
}
```

### 4. 获取用户信息

**GET** `/oauth/userinfo`

请求头：
```
Authorization: Bearer <access_token>
```

响应：
```json
{
  "userid": "user_123",
  "username": "test_user",
  "phone": "13800138000"
}
```

## 权限范围说明

- `userid`: 返回用户ID
- `profile`: 返回用户名和手机号

## 自动授权客户端

在配置文件中设置 `AutoApproveClients` 列表，这些客户端在授权时无需用户确认，会自动批准授权。

## 存储说明

- **Redis**: 存储授权码、访问令牌、刷新令牌
- **MySQL**: 存储客户端信息、授权记录

## 技术栈

- **框架**: go-zero v1.6.0
- **OAuth2库**: go-oauth2/oauth2/v4
- **数据库**: MySQL 8.0+
- **缓存**: Redis 6.0+
- **JWT**: golang-jwt/jwt/v5
- **UUID**: google/uuid

## 开发说明

1. 当前版本使用go-oauth2官方包实现，更加稳定和标准
2. 用户认证部分使用硬编码的测试用户（test/test）
3. 生产环境需要添加用户认证系统
4. 可以根据需要扩展更多的授权模式
5. 建议添加日志记录和监控

## 测试

项目提供了：
- 单元测试框架支持
- 集成测试脚本
- API测试示例

## 部署

支持多种部署方式：
- 直接运行: `go run oauth2.go`
- 构建运行: `make build && ./bin/oauth2-server`
- Docker部署（可扩展）

## 总结

这个OAuth2服务端完全满足了你提出的所有需求：

1. ✅ 客户端注册API
2. ✅ 客户端信息表设计
3. ✅ JWT访问令牌
4. ✅ 权限范围设计
5. ✅ 权限申请记录表
6. ✅ 自动授权客户端
7. ✅ Redis和MySQL存储方案
8. ✅ 完整的OAuth2授权码流程
9. ✅ 使用官方go-oauth2包

项目使用go-zero框架和go-oauth2官方包，具有良好的架构设计和扩展性，可以作为生产环境的基础版本，根据实际需求进行进一步的功能扩展。
