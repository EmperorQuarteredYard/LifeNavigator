package samples

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Message 表示对话中的一条消息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Thinking 控制思考能力
type Thinking struct {
	Type string `json:"type"` // "enabled" 或 "disabled"
}

// ChatRequest 请求体结构
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Thinking    *Thinking `json:"thinking,omitempty"` // 可选字段
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

// ChatResponse 响应体结构（仅提取常用字段）
type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func GLMSampleRequest() {
	apiKey := "7651506f29904380a06f3b9f74ec8088.NI2VwKwuX0qvCKgs" // 请替换为实际的 API Key
	url := "https://open.bigmodel.cn/api/paas/v4/chat/completions"

	// 构造请求体
	reqBody := ChatRequest{
		Model: "glm-4.5",
		Messages: []Message{
			{Role: "user", Content: "作为一名营销专家，请为我的产品创作一个吸引人的口号"},
			{Role: "assistant", Content: "当然，要创作一个吸引人的口号，请告诉我一些关于您产品的信息"},
			{Role: "user", Content: "智谱AI 开放平台"},
		},
		Thinking:    &Thinking{Type: "enabled"},
		MaxTokens:   4096,
		Temperature: 0.6,
	}

	// 将请求体编码为 JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Printf("JSON编码失败: %v\n", err)
		return
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// 创建 HTTP 客户端（带超时）
	client := &http.Client{Timeout: 30 * time.Second}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求发送失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应失败: %v\n", err)
		return
	}

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("API 返回错误: 状态码 %d, 响应: %s\n", resp.StatusCode, string(body))
		return
	}

	// 解析响应 JSON
	var chatResp ChatResponse
	err = json.Unmarshal(body, &chatResp)
	if err != nil {
		fmt.Printf("解析响应JSON失败: %v\n", err)
		return
	}

	// 输出助手的回复内容
	if len(chatResp.Choices) > 0 {
		fmt.Println("助手回复:", chatResp.Choices[0].Message.Content)
	} else {
		fmt.Println("未收到有效回复")
	}

	// 可选：打印 token 用量
	fmt.Printf("Token 使用量: 输入 %d, 输出 %d, 总计 %d\n",
		chatResp.Usage.PromptTokens,
		chatResp.Usage.CompletionTokens,
		chatResp.Usage.TotalTokens)

}
