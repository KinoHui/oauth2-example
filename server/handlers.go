package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

// RegisterClientHandler 客户端注册处理器
func RegisterClientHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ClientRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 验证必填字段
	if req.Name == "" || req.RedirectURL == "" || req.GrantType == "" || req.Scope == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 验证授权模式
	if req.GrantType != "authorization_code" {
		http.Error(w, "Only authorization_code grant type is supported", http.StatusBadRequest)
		return
	}

	// 验证权限范围
	scopes := strings.Split(req.Scope, " ")
	for _, scope := range scopes {
		if scope != "userid" && scope != "profile" {
			http.Error(w, "Invalid scope: "+scope, http.StatusBadRequest)
			return
		}
	}

	// 创建客户端
	client, err := db.CreateClient(&req)
	if err != nil {
		http.Error(w, "Failed to create client", http.StatusInternalServerError)
		return
	}

	response := ClientRegistrationResponse{
		ClientID:     client.ID,
		ClientSecret: client.Secret,
		Message:      "Client registered successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UserInfoHandler 用户信息处理器
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// 从Authorization header获取token
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	// 验证Bearer token格式
	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	// 验证token并获取用户信息
	userID, scope, err := validateToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// 根据scope返回相应的用户信息
	userInfo := UserInfo{
		UserID:   userID,
		Username: "test_user", // 这里应该从数据库获取真实用户信息
		Phone:    "13800138000",
	}

	// 根据请求的scope过滤返回的信息
	response := make(map[string]interface{})
	scopes := strings.Split(scope, " ")

	for _, s := range scopes {
		switch s {
		case "userid":
			response["userid"] = userInfo.UserID
		case "profile":
			response["username"] = userInfo.Username
			response["phone"] = userInfo.Phone
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// validateToken 验证token并返回用户ID和scope
func validateToken(tokenString string) (string, string, error) {
	// 解析JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return "", "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, _ := claims["user_id"].(string)
		scope, _ := claims["scope"].(string)
		return userID, scope, nil
	}

	return "", "", fmt.Errorf("invalid token")
}

// SetAutoApproveHandler 设置客户端自动授权处理器
func SetAutoApproveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ClientID    string `json:"client_id"`
		AutoApprove bool   `json:"auto_approve"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := db.SetAutoApprove(req.ClientID, req.AutoApprove); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]string{
		"message": fmt.Sprintf("Client %s auto approve set to %v", req.ClientID, req.AutoApprove),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListClientsHandler 列出所有客户端处理器
func ListClientsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	clients := db.GetAllClients()

	// 隐藏敏感信息
	var safeClients []map[string]interface{}
	for _, client := range clients {
		safeClient := map[string]interface{}{
			"id":           client.ID,
			"name":         client.Name,
			"redirect_url": client.RedirectURL,
			"grant_type":   client.GrantType,
			"scope":        client.Scope,
			"auto_approve": client.AutoApprove,
			"created_at":   client.CreatedAt,
			"updated_at":   client.UpdatedAt,
		}
		safeClients = append(safeClients, safeClient)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(safeClients)
}
