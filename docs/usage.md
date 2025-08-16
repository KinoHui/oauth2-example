# OAuth2 服务端使用示例

## 1. 环境准备

确保你已经安装了以下软件：
- Go 1.21+
- MySQL 8.0+
- Redis 6.0+

## 2. 数据库初始化

```bash
# 创建数据库和表
mysql -u root -p < scripts/init.sql
```

## 3. 配置修改

修改 `etc/oauth2-api.yaml` 文件中的数据库和Redis连接信息：

```yaml
MySQL:
  DataSource: root:your_password@tcp(localhost:3306)/oauth2?charset=utf8mb4&parseTime=True&loc=Local

Redis:
  Host: localhost:6379
  Type: node
  Pass: ""
```

## 4. 启动服务

```bash
# 方式1：直接运行
go run oauth2.go

# 方式2：使用Makefile
make run

# 方式3：构建后运行
make build
./bin/oauth2-server
```

## 5. API 使用示例

### 5.1 客户端注册

```bash
curl -X POST http://localhost:8080/api/client/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "我的应用",
    "redirect_url": "http://localhost:3000/callback",
    "grant_type": "authorization_code",
    "scope": "userid profile"
  }'
```

响应示例：
```json
{
  "client_id": "client_abc123",
  "client_secret": "secret_xyz789"
}
```

### 5.2 授权请求

```bash
curl "http://localhost:8080/oauth/authorize?client_id=trusted_client_001&response_type=code&redirect_uri=http://localhost:3000/callback&scope=userid%20profile&state=test_state"
```

响应示例：
```json
{
  "code": "auth_code_123",
  "state": "test_state"
}
```

### 5.3 获取访问令牌

```bash
curl -X POST http://localhost:8080/oauth/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=authorization_code&code=auth_code_123&redirect_uri=http://localhost:3000/callback&client_id=trusted_client_001&client_secret=trusted_secret_001"
```

响应示例：
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 7200,
  "refresh_token": "refresh_token_123",
  "scope": "userid profile"
}
```

### 5.4 获取用户信息

```bash
curl -X GET http://localhost:8080/oauth/userinfo \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

响应示例：
```json
{
  "userid": "test_user",
  "username": "test_user",
  "phone": "13800138000"
}
```

## 6. 完整的OAuth2流程示例

### 6.1 使用预定义的可信客户端

```bash
# 1. 授权请求
AUTH_RESPONSE=$(curl -s "http://localhost:8080/oauth/authorize?client_id=trusted_client_001&response_type=code&redirect_uri=http://localhost:3000/callback&scope=userid%20profile&state=test_state")

# 2. 提取授权码
CODE=$(echo $AUTH_RESPONSE | jq -r '.code')

# 3. 获取访问令牌
TOKEN_RESPONSE=$(curl -s -X POST http://localhost:8080/oauth/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=authorization_code&code=$CODE&redirect_uri=http://localhost:3000/callback&client_id=trusted_client_001&client_secret=trusted_secret_001")

# 4. 提取访问令牌
ACCESS_TOKEN=$(echo $TOKEN_RESPONSE | jq -r '.access_token')

# 5. 获取用户信息
curl -X GET http://localhost:8080/oauth/userinfo \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

### 6.2 使用新注册的客户端

```bash
# 1. 注册新客户端
REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/client/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "测试应用",
    "redirect_url": "http://localhost:3000/callback",
    "grant_type": "authorization_code",
    "scope": "userid profile"
  }')

# 2. 提取客户端信息
CLIENT_ID=$(echo $REGISTER_RESPONSE | jq -r '.client_id')
CLIENT_SECRET=$(echo $REGISTER_RESPONSE | jq -r '.client_secret')

# 3. 授权请求
AUTH_RESPONSE=$(curl -s "http://localhost:8080/oauth/authorize?client_id=$CLIENT_ID&response_type=code&redirect_uri=http://localhost:3000/callback&scope=userid%20profile&state=test_state")

# 4. 提取授权码
CODE=$(echo $AUTH_RESPONSE | jq -r '.code')

# 5. 获取访问令牌
TOKEN_RESPONSE=$(curl -s -X POST http://localhost:8080/oauth/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=authorization_code&code=$CODE&redirect_uri=http://localhost:3000/callback&client_id=$CLIENT_ID&client_secret=$CLIENT_SECRET")

# 6. 提取访问令牌
ACCESS_TOKEN=$(echo $TOKEN_RESPONSE | jq -r '.access_token')

# 7. 获取用户信息
curl -X GET http://localhost:8080/oauth/userinfo \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

## 7. 测试脚本

项目提供了一个简单的测试脚本：

```bash
# 运行测试脚本
make test-client
```

## 8. 注意事项

1. **自动授权客户端**：在配置文件中设置的 `AutoApproveClients` 列表中的客户端会自动授权，无需用户确认。

2. **权限范围**：
   - `userid`：返回用户ID
   - `profile`：返回用户名和手机号

3. **令牌过期**：
   - 访问令牌：2小时（可在配置中修改）
   - 刷新令牌：30天
   - 授权码：10分钟

4. **安全性**：
   - 生产环境请修改JWT密钥
   - 建议使用HTTPS
   - 定期轮换客户端密钥

5. **扩展性**：
   - 当前版本使用硬编码的测试用户（test/test）
   - 生产环境需要集成真实的用户认证系统
   - 可以根据需要添加更多的授权模式

6. **go-oauth2包特性**：
   - 使用官方go-oauth2包，更加稳定和标准
   - 支持完整的OAuth2规范
   - 内置JWT令牌生成
   - 支持多种存储后端 