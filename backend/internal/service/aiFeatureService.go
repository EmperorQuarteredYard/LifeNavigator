package service

import (
	"LifeNavigator/internal/models"
	"LifeNavigator/internal/repository"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type AIFeatureService interface {
	ReduceProject(ctx context.Context, userID uint64, projectDescription string, accountSummary string, eventChan chan<- StreamEvent) error
	Summary(ctx context.Context, userID uint64, startTime, endTime time.Time, eventChan chan<- StreamEvent) (string, error)
	UpdateUserProfile(userID uint64, profile string) error
	GetUserCompletedTasks(userID uint64, startTime, endTime time.Time) ([]models.Task, string, error)
}

type StreamEvent struct {
	Type    string
	Content interface{}
}

type aiFeatureService struct {
	transactor repository.Transactor
}

func NewAIFeatureService(transactor repository.Transactor) AIFeatureService {
	return &aiFeatureService{
		transactor: transactor,
	}
}

const (
	EventTypeProjectCreated = "project_created"
	EventTypeTaskCreated    = "task_created"
	EventTypeBudgetCreated  = "budget_created"
	EventTypeComplete       = "stream_complete"
	EventTypeError          = "stream_error"
	EventTypeSummaryContent = "summary_content"
	StreamEndMarker         = "[STREAM_END]"
)

type GLMMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GLMRequest struct {
	Model       string       `json:"model"`
	Messages    []GLMMessage `json:"messages"`
	Stream      bool         `json:"stream"`
	MaxTokens   int          `json:"max_tokens"`
	Temperature float64      `json:"temperature"`
}

type GLMStreamResponse struct {
	ID      string `json:"id"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index int `json:"index"`
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

type ParsedProject struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ParsedTask struct {
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	Category       string  `json:"category"`
	Deadline       *string `json:"deadline"`
	PrerequisiteID *int    `json:"prerequisite_id"`
}

type ParsedBudget struct {
	Type   string  `json:"type"`
	Budget float64 `json:"budget"`
}

func (s *aiFeatureService) ReduceProject(ctx context.Context, userID uint64, projectDescription string, accountSummary string, eventChan chan<- StreamEvent) error {
	defer close(eventChan)

	glmToken := os.Getenv("GLM_TOKEN")
	if glmToken == "" {
		eventChan <- StreamEvent{Type: EventTypeError, Content: "GLM_TOKEN not configured"}
		return fmt.Errorf("GLM_TOKEN not configured")
	}

	prompt := s.buildPrompt(projectDescription, accountSummary)

	reqBody := GLMRequest{
		Model: "glm-4.6",
		Messages: []GLMMessage{
			{Role: "system", Content: s.getSystemPrompt()},
			{Role: "user", Content: prompt},
		},
		Stream:      true,
		MaxTokens:   8192,
		Temperature: 0.7,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		eventChan <- StreamEvent{Type: EventTypeError, Content: "Failed to marshal request"}
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://open.bigmodel.cn/api/paas/v4/chat/completions", strings.NewReader(string(jsonBody)))
	if err != nil {
		eventChan <- StreamEvent{Type: EventTypeError, Content: "Failed to create request"}
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+glmToken)

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		eventChan <- StreamEvent{Type: EventTypeError, Content: "Failed to call GLM API"}
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		eventChan <- StreamEvent{Type: EventTypeError, Content: fmt.Sprintf("GLM API error: %s", string(body))}
		return fmt.Errorf("GLM API returned status %d", resp.StatusCode)
	}

	var fullContent strings.Builder
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var glmResp GLMStreamResponse
		if err := json.Unmarshal([]byte(data), &glmResp); err != nil {
			continue
		}

		if len(glmResp.Choices) > 0 && glmResp.Choices[0].Delta.Content != "" {
			fullContent.WriteString(glmResp.Choices[0].Delta.Content)
		}
	}

	if err := scanner.Err(); err != nil {
		eventChan <- StreamEvent{Type: EventTypeError, Content: "Error reading stream"}
		return err
	}

	return s.parseAndCreateEntities(ctx, userID, fullContent.String(), eventChan)
}

func (s *aiFeatureService) getSystemPrompt() string {
	return `你是一个专业的项目管理助手，擅长将用户的项目描述转化为结构化的项目计划。

你的任务是根据用户提供的项目描述，生成一个完整的项目计划，包括：
1. 项目基本信息（名称、描述）
2. 任务列表（包含任务名称、描述、分类、截止时间、前置任务关系）
3. 项目预算（包含预算类型和金额）

## 预算类型说明
常见的预算类型包括：
- time: 时间预算（单位：小时）
- money: 金钱预算（单位：元）
- token: Token预算（用于AI相关项目）
- energy: 精力预算（抽象单位）

## 任务分类说明
常见的任务分类包括：
- planning: 规划类
- development: 开发类
- design: 设计类
- testing: 测试类
- documentation: 文档类
- deployment: 部署类
- research: 研究类
- communication: 沟通类

## 前置任务关系说明
任务之间可能存在依赖关系，prerequisite_id表示该任务依赖的前置任务的索引（从0开始）。
例如：prerequisite_id为0表示该任务依赖于第一个任务完成后才能开始。

## 输出格式
请严格按照以下JSON格式输出，不要包含任何其他内容：

` + "```json" + `
{
  "project": {
    "name": "项目名称",
    "description": "项目详细描述"
  },
  "tasks": [
    {
      "name": "任务名称",
      "description": "任务描述",
      "category": "任务分类",
      "deadline": "2024-12-31T23:59:59Z",
      "prerequisite_id": null
    }
  ],
  "budgets": [
    {
      "type": "money",
      "budget": 10000.00
    }
  ]
}
` + "```" + `

请确保：
1. 所有任务都有明确的名称和描述
2. 截止时间使用ISO 8601格式
3. 预算金额为正数
4. 前置任务关系合理，不存在循环依赖`
}

func (s *aiFeatureService) buildPrompt(projectDescription string, accountSummary string) string {
	prompt := fmt.Sprintf(`请根据以下信息创建一个完整的项目计划：

## 项目描述
%s

`, projectDescription)

	if accountSummary != "" {
		prompt += fmt.Sprintf(`## 用户账户摘要
%s

`, accountSummary)
	}

	prompt += `请生成项目计划，确保任务安排合理，预算设置适当。`
	return prompt
}

type AIResponse struct {
	Project ParsedProject  `json:"project"`
	Tasks   []ParsedTask   `json:"tasks"`
	Budgets []ParsedBudget `json:"budgets"`
}

func (s *aiFeatureService) parseAndCreateEntities(ctx context.Context, userID uint64, content string, eventChan chan<- StreamEvent) error {
	jsonContent := s.extractJSON(content)
	if jsonContent == "" {
		eventChan <- StreamEvent{Type: EventTypeError, Content: "Failed to extract valid JSON from AI response"}
		return fmt.Errorf("failed to extract valid JSON from AI response")
	}

	var aiResp AIResponse
	if err := json.Unmarshal([]byte(jsonContent), &aiResp); err != nil {
		eventChan <- StreamEvent{Type: EventTypeError, Content: fmt.Sprintf("Failed to parse AI response: %v", err)}
		return err
	}

	var createdProject *models.Project
	taskIDMap := make(map[int]uint64)

	err := s.transactor.WithinTransaction(ctx, func(txRepo repository.TxRepositories) error {
		project := &models.Project{
			UserID:      userID,
			Name:        aiResp.Project.Name,
			Description: aiResp.Project.Description,
		}

		if err := txRepo.Project.Create(project); err != nil {
			return fmt.Errorf("failed to create project: %w", err)
		}
		createdProject = project

		eventChan <- StreamEvent{
			Type: EventTypeProjectCreated,
			Content: map[string]interface{}{
				"id":          project.ID,
				"name":        project.Name,
				"description": project.Description,
			},
		}

		for i, taskData := range aiResp.Tasks {
			var deadline *time.Time
			if taskData.Deadline != nil {
				t, err := time.Parse(time.RFC3339, *taskData.Deadline)
				if err == nil {
					deadline = &t
				}
			}

			task := &models.Task{
				UserID:      userID,
				ProjectID:   project.ID,
				Name:        taskData.Name,
				Description: taskData.Description,
				Category:    taskData.Category,
				Deadline:    deadline,
				Status:      0,
				Type:        0,
			}

			if err := txRepo.Task.Create(task); err != nil {
				return fmt.Errorf("failed to create task: %w", err)
			}

			taskIDMap[i] = task.ID

			var prerequisiteID *uint64
			if taskData.PrerequisiteID != nil {
				if pid, ok := taskIDMap[*taskData.PrerequisiteID]; ok {
					prerequisiteID = &pid
				}
			}

			eventChan <- StreamEvent{
				Type: EventTypeTaskCreated,
				Content: map[string]interface{}{
					"id":              task.ID,
					"project_id":      task.ProjectID,
					"name":            task.Name,
					"description":     task.Description,
					"category":        task.Category,
					"deadline":        taskData.Deadline,
					"prerequisite_id": prerequisiteID,
				},
			}

			if prerequisiteID != nil {
				if _, err := txRepo.Task.SetPrerequisiteTask(*prerequisiteID, task.ID); err != nil {
					log.Printf("Warning: failed to create task dependency: %v", err)
				}
			}
		}

		for _, budgetData := range aiResp.Budgets {
			budget := &models.ProjectBudget{
				ProjectID: project.ID,
				Type:      budgetData.Type,
				Budget:    budgetData.Budget,
				Used:      0,
			}

			if err := txRepo.ProjectBudget.Create(budget); err != nil {
				return fmt.Errorf("failed to create budget: %w", err)
			}

			eventChan <- StreamEvent{
				Type: EventTypeBudgetCreated,
				Content: map[string]interface{}{
					"id":         budget.ID,
					"project_id": budget.ProjectID,
					"type":       budget.Type,
					"budget":     budget.Budget,
				},
			}
		}

		return nil
	})

	if err != nil {
		eventChan <- StreamEvent{Type: EventTypeError, Content: err.Error()}
		return err
	}

	log.Printf("Project created successfully: ID=%d, Name=%s", createdProject.ID, createdProject.Name)

	eventChan <- StreamEvent{
		Type: EventTypeComplete,
		Content: map[string]string{
			"message": "Project created successfully",
		},
	}

	return nil
}

func (s *aiFeatureService) extractJSON(content string) string {
	jsonStart := strings.Index(content, "{")
	if jsonStart == -1 {
		return ""
	}

	braceCount := 0
	jsonEnd := -1

	for i := jsonStart; i < len(content); i++ {
		if content[i] == '{' {
			braceCount++
		} else if content[i] == '}' {
			braceCount--
			if braceCount == 0 {
				jsonEnd = i + 1
				break
			}
		}
	}

	if jsonEnd == -1 {
		return ""
	}

	return content[jsonStart:jsonEnd]
}

func (s *aiFeatureService) Summary(ctx context.Context, userID uint64, startTime, endTime time.Time, eventChan chan<- StreamEvent) (string, error) {
	defer close(eventChan)

	glmToken := os.Getenv("GLM_TOKEN")
	if glmToken == "" {
		eventChan <- StreamEvent{Type: EventTypeError, Content: "GLM_TOKEN not configured"}
		return "", fmt.Errorf("GLM_TOKEN not configured")
	}

	completedTasks, currentProfile, err := s.GetUserCompletedTasks(userID, startTime, endTime)
	if err != nil {
		eventChan <- StreamEvent{Type: EventTypeError, Content: err.Error()}
		return "", err
	}

	if len(completedTasks) == 0 {
		eventChan <- StreamEvent{Type: EventTypeError, Content: "No completed tasks found in the specified time range"}
		return "", fmt.Errorf("no completed tasks found")
	}

	prompt := s.buildSummaryPrompt(completedTasks, currentProfile)

	reqBody := GLMRequest{
		Model: "glm-4.6",
		Messages: []GLMMessage{
			{Role: "system", Content: s.getSummarySystemPrompt()},
			{Role: "user", Content: prompt},
		},
		Stream:      true,
		MaxTokens:   4096,
		Temperature: 0.8,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		eventChan <- StreamEvent{Type: EventTypeError, Content: "Failed to marshal request"}
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://open.bigmodel.cn/api/paas/v4/chat/completions", strings.NewReader(string(jsonBody)))
	if err != nil {
		eventChan <- StreamEvent{Type: EventTypeError, Content: "Failed to create request"}
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+glmToken)

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		eventChan <- StreamEvent{Type: EventTypeError, Content: "Failed to call GLM API"}
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		eventChan <- StreamEvent{Type: EventTypeError, Content: fmt.Sprintf("GLM API error: %s", string(body))}
		return "", fmt.Errorf("GLM API returned status %d", resp.StatusCode)
	}

	var fullContent strings.Builder
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var glmResp GLMStreamResponse
		if err := json.Unmarshal([]byte(data), &glmResp); err != nil {
			continue
		}

		if len(glmResp.Choices) > 0 && glmResp.Choices[0].Delta.Content != "" {
			content := glmResp.Choices[0].Delta.Content
			fullContent.WriteString(content)
			eventChan <- StreamEvent{
				Type:    EventTypeSummaryContent,
				Content: map[string]string{"content": content},
			}
		}
	}

	if err := scanner.Err(); err != nil {
		eventChan <- StreamEvent{Type: EventTypeError, Content: "Error reading stream"}
		return "", err
	}

	eventChan <- StreamEvent{
		Type: EventTypeComplete,
		Content: map[string]string{
			"message": "Summary completed successfully",
		},
	}

	return fullContent.String(), nil
}

func (s *aiFeatureService) GetUserCompletedTasks(userID uint64, startTime, endTime time.Time) ([]models.Task, string, error) {
	var completedTasks []models.Task
	var currentProfile string

	err := s.transactor.WithinTransaction(context.Background(), func(txRepo repository.TxRepositories) error {
		var err error
		completedTasks, err = txRepo.Task.ListCompletedByUserIDAndTimeRange(userID, startTime, endTime)
		if err != nil {
			return fmt.Errorf("failed to get completed tasks: %w", err)
		}

		user, err := txRepo.User.GetByID(userID)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}
		currentProfile = user.Profile

		return nil
	})

	return completedTasks, currentProfile, err
}

func (s *aiFeatureService) UpdateUserProfile(userID uint64, profile string) error {
	return s.transactor.WithinTransaction(context.Background(), func(txRepo repository.TxRepositories) error {
		return txRepo.User.UpdateProfile(userID, profile)
	})
}

func (s *aiFeatureService) getSummarySystemPrompt() string {
	return `你是一个专业的生活助手，擅长总结用户的成就并更新用户画像。

你的任务是：
1. 总结用户在指定时间段内完成的任务和成就
2. 识别用户值得鼓励的行为和进步
3. 基于用户的最新表现，更新用户画像

## 用户画像说明
用户画像是对用户特点、习惯、优势、兴趣等方面的简短描述（不超过2000字符）。
画像应该：
- 突出用户的优点和特长
- 反映用户的工作习惯和偏好
- 记录用户的成长和进步
- 保持积极正面的语调

## 输出格式
请按以下格式输出：

` + "```json" + `
{
  "summary": "对用户成就的详细总结，突出值得鼓励的事情...",
  "profile_update": "更新后的用户画像，简短描述用户的特点和进步..."
}
` + "```" + `

请确保：
1. 总结内容具体、真实，基于提供的任务数据
2. 语气积极鼓励，肯定用户的努力
3. 用户画像更新简洁明了，不超过2000字符
4. 如果当前画像已经很好，可以保持不变或做小幅调整`
}

func (s *aiFeatureService) buildSummaryPrompt(tasks []models.Task, currentProfile string) string {
	prompt := "请总结用户在指定时间段内的成就：\n\n"

	prompt += "## 已完成的任务\n"
	for i, task := range tasks {
		prompt += fmt.Sprintf("%d. %s\n", i+1, task.Name)
		if task.Description != "" {
			prompt += fmt.Sprintf("   描述: %s\n", task.Description)
		}
		if task.Category != "" {
			prompt += fmt.Sprintf("   分类: %s\n", task.Category)
		}
		if task.CompletedAt != nil {
			prompt += fmt.Sprintf("   完成时间: %s\n", task.CompletedAt.Format("2006-01-02 15:04:05"))
		}
		prompt += "\n"
	}

	if currentProfile != "" {
		prompt += fmt.Sprintf("## 当前用户画像\n%s\n\n", currentProfile)
	}

	prompt += "请基于以上信息，总结用户的成就并更新用户画像。"
	return prompt
}

func (s *aiFeatureService) extractProfileUpdate(content string) string {
	jsonContent := s.extractJSON(content)
	if jsonContent == "" {
		return ""
	}

	var result struct {
		Summary       string `json:"summary"`
		ProfileUpdate string `json:"profile_update"`
	}

	if err := json.Unmarshal([]byte(jsonContent), &result); err != nil {
		return ""
	}

	return result.ProfileUpdate
}

func (s *aiFeatureService) updateUserProfile(userID uint64, profile string) error {
	return s.transactor.WithinTransaction(context.Background(), func(txRepo repository.TxRepositories) error {
		return txRepo.User.UpdateProfile(userID, profile)
	})
}
