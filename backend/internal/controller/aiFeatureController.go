package controller

import (
	"LifeNavigator/internal/service"
	"LifeNavigator/pkg/dto"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type AIFeatureController struct {
	aiFeatureServ service.AIFeatureService
	accountServ   service.AccountService
	*BaseController
}

func NewAIFeatureController(aiFeatureServ service.AIFeatureService, accountServ service.AccountService) *AIFeatureController {
	return &AIFeatureController{
		aiFeatureServ:  aiFeatureServ,
		accountServ:    accountServ,
		BaseController: &BaseController{},
	}
}

func (ctl *AIFeatureController) ReduceProject(c *gin.Context) {
	var req dto.AIReduceProjectRequest
	if !ctl.BindJSON(c, &req) {
		return
	}

	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	accountSummary, err := ctl.buildAccountSummary(authUser.UserID)
	if err != nil {
		ctl.Error(c, err)
		return
	}

	c.Header("Content-Unit", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	eventChan := make(chan service.StreamEvent, 100)

	go func() {
		err := ctl.aiFeatureServ.ReduceProject(
			c.Request.Context(),
			authUser.UserID,
			req.ProjectDescription,
			accountSummary,
			eventChan,
		)
		if err != nil {
			c.Error(err)
		}
	}()

	for event := range eventChan {
		eventData, err := json.Marshal(event)
		if err != nil {
			continue
		}

		if _, err := fmt.Fprintf(c.Writer, "data: %s\n\n", string(eventData)); err != nil {
			log.Printf("Client disconnected during ReduceProject: %v", err)
			return
		}
		c.Writer.Flush()
	}

	if _, err := fmt.Fprintf(c.Writer, "data: %s\n\n", service.StreamEndMarker); err != nil {
		log.Printf("Failed to write stream end marker: %v", err)
		return
	}
	c.Writer.Flush()
}

func (ctl *AIFeatureController) buildAccountSummary(userID uint64) (string, error) {
	accounts, err := ctl.accountServ.ListByUserID(userID)
	if err != nil {
		return "", err
	}

	if len(accounts) == 0 {
		return "", nil
	}

	summary := "用户账户信息：\n"
	for _, acc := range accounts {
		summary += fmt.Sprintf("- %s: 余额 %.2f\n", acc.Type, acc.Balance)
	}

	return summary, nil
}

func (ctl *AIFeatureController) Summary(c *gin.Context) {
	var req dto.AISummaryRequest
	if !ctl.BindJSON(c, &req) {
		return
	}

	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		ctl.Code(c, 400)
		return
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		ctl.Code(c, 400)
		return
	}

	c.Header("Content-Unit", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	eventChan := make(chan service.StreamEvent, 100)
	contentChan := make(chan string, 1)

	go func() {
		content, err := ctl.aiFeatureServ.Summary(
			c.Request.Context(),
			authUser.UserID,
			startTime,
			endTime,
			eventChan,
		)
		if err != nil {
			c.Error(err)
			contentChan <- ""
			return
		}
		contentChan <- content
	}()

	for event := range eventChan {
		eventData, err := json.Marshal(event)
		if err != nil {
			continue
		}

		if _, err := fmt.Fprintf(c.Writer, "data: %s\n\n", string(eventData)); err != nil {
			log.Printf("Client disconnected during Summary: %v", err)
			return
		}
		c.Writer.Flush()
	}

	fullContent := <-contentChan

	go func() {
		profileUpdate := ctl.extractProfileUpdate(fullContent)
		if profileUpdate != "" {
			if err := ctl.aiFeatureServ.UpdateUserProfile(authUser.UserID, profileUpdate); err != nil {
				log.Printf("Failed to update user profile asynchronously: %v", err)
			} else {
				log.Printf("User profile updated successfully for user %d", authUser.UserID)
			}
		}
	}()

	if _, err := fmt.Fprintf(c.Writer, "data: %s\n\n", service.StreamEndMarker); err != nil {
		log.Printf("Failed to write stream end marker: %v", err)
		return
	}
	c.Writer.Flush()
}

func (ctl *AIFeatureController) extractProfileUpdate(content string) string {
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

	jsonContent := content[jsonStart:jsonEnd]

	var result struct {
		Summary       string `json:"summary"`
		ProfileUpdate string `json:"profile_update"`
	}

	if err := json.Unmarshal([]byte(jsonContent), &result); err != nil {
		return ""
	}

	return result.ProfileUpdate
}
