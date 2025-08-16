package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const baseURL = "http://localhost:8080"

// 测试客户端注册
func testClientRegister() {
	url := baseURL + "/api/client/register"
	data := map[string]string{
		"name":         "测试应用",
		"redirect_url": "http://localhost:3000/callback",
		"grant_type":   "authorization_code",
		"scope":        "userid profile",
	}

	jsonData, _ := json.Marshal(data)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("客户端注册失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("客户端注册响应: %s\n", string(body))
}

// 测试授权请求（使用预定义的可信客户端）
func testAuthorize() {
	params := url.Values{}
	params.Add("client_id", "trusted_client_001")
	params.Add("response_type", "code")
	params.Add("redirect_uri", "http://localhost:3000/callback")
	params.Add("scope", "userid profile")
	params.Add("state", "test_state")

	url := baseURL + "/oauth/authorize?" + params.Encode()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("授权请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("授权响应: %s\n", string(body))
}

// 测试获取访问令牌
func testToken(code string) {
	fullUrl := baseURL + "/oauth/token"
	data := url.Values{}
	data.Add("grant_type", "authorization_code")
	data.Add("code", code)
	data.Add("redirect_uri", "http://localhost:3000/callback")
	data.Add("client_id", "trusted_client_001")
	data.Add("client_secret", "trusted_secret_001")

	resp, err := http.PostForm(fullUrl, data)
	if err != nil {
		fmt.Printf("获取令牌失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("令牌响应: %s\n", string(body))
}

// 测试获取用户信息
func testUserInfo(accessToken string) {
	url := baseURL + "/oauth/userinfo"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("获取用户信息失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("用户信息响应: %s\n", string(body))
}

func main() {
	fmt.Println("开始测试OAuth2服务端...")

	// 1. 测试客户端注册
	fmt.Println("\n1. 测试客户端注册")
	testClientRegister()

	// 2. 测试授权请求
	fmt.Println("\n2. 测试授权请求")
	testAuthorize()

	// 注意：在实际使用中，你需要：
	// 1. 先调用授权接口获取授权码
	// 2. 使用授权码调用令牌接口获取访问令牌
	// 3. 使用访问令牌调用用户信息接口

	fmt.Println("\n测试完成！")
	fmt.Println("注意：这是一个基本测试，实际使用时需要按照OAuth2流程逐步调用各个接口。")
}
