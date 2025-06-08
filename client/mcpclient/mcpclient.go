package mcpclient

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

// 聊天消息结构
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// 创建会话请求结构
type CreateSessionRequest struct {
	SessionID string        `json:"session_id"`
	Context   []ChatMessage `json:"context,omitempty"`
}

// 创建会话响应结构
type CreateSessionResponse struct {
	Success bool `json:"success"`
}

// 聊天请求结构
type ChatRequest struct {
	SessionID string `json:"session_id"`
	Content   string `json:"content"`
}

// 聊天响应结构
type ChatResponse struct {
	Response string `json:"response"`
}

// MCPClient 是客户端结构体
type MCPClient struct {
	baseURL string
	client  *resty.Client
}

// NewMCPClient 构造函数
func NewMCPClient(baseURL string) *MCPClient {
	return &MCPClient{
		baseURL: baseURL,
		client:  resty.New(),
	}
}

// CreateSession 创建会话
func (c *MCPClient) CreateSession(sessionID string, context []ChatMessage) error {
	reqBody := CreateSessionRequest{
		SessionID: sessionID,
		Context:   context,
	}

	var resp CreateSessionResponse
	r, err := c.client.R().
		SetBody(reqBody).
		SetResult(&resp).
		Post(c.baseURL + "/sessions")

	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}

	if r.IsError() {
		return fmt.Errorf("服务端错误: %s", r.Status())
	}

	if !resp.Success {
		return fmt.Errorf("创建会话失败")
	}

	return nil
}

// Chat 向某个会话发送消息并获取回复
func (c *MCPClient) Chat(sessionID string, content string) (string, error) {
	reqBody := ChatRequest{
		SessionID: sessionID,
		Content:   content,
	}

	var resp ChatResponse
	r, err := c.client.R().
		SetBody(reqBody).
		SetResult(&resp).
		Post(c.baseURL + "/chat")

	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}

	if r.IsError() {
		return "", fmt.Errorf("服务端错误: %s", r.Status())
	}

	return resp.Response, nil
}
