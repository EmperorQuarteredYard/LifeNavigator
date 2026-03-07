package controller

import (
	"LifeNavigator/internal/models"
	"LifeNavigator/internal/service"
	"LifeNavigator/pkg/dto"
	"LifeNavigator/pkg/errcode"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	taskServ service.TaskService
	*BaseController
}

func NewTaskController(taskServ service.TaskService) *TaskController {
	return &TaskController{
		taskServ:       taskServ,
		BaseController: &BaseController{},
	}
}

// CreateTask 创建任务
func (ctl *TaskController) CreateTask(c *gin.Context) {
	var req dto.CreateTaskRequest
	if !ctl.BindJSON(c, &req) {
		return
	}
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	deadline, err := time.Parse(time.RFC3339, *req.Deadline)
	if err != nil {
		ctl.ServerError(c)
		return
	}
	task := &models.Task{
		Name:           req.Name,
		ProjectID:      req.ProjectID,
		Description:    req.Description,
		AutoCalculated: req.AutoCalculated,
		Type:           req.Type,
		Status:         req.Status,
		Category:       req.Category,
		Deadline:       &deadline,
	}
	err = ctl.taskServ.Create(task, authUser.UserID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, task)
}

// DeleteTask 删除任务
func (ctl *TaskController) DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctl.Code(c, errcode.StatusInvalidParams)
		return
	}
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	err = ctl.taskServ.Delete(id, authUser.UserID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, gin.H{"message": "task deleted successfully"})
}

// GetTask 获取单个任务
func (ctl *TaskController) GetTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctl.Code(c, errcode.StatusInvalidParams)
		return
	}
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	task, err := ctl.taskServ.GetByID(id, authUser.UserID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, task)
}

// UpdateTask 更新任务
func (ctl *TaskController) UpdateTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctl.Code(c, errcode.StatusInvalidParams)
		return
	}
	var req dto.UpdateTaskRequest
	if !ctl.BindJSON(c, &req) {
		return
	}
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	task := &models.Task{
		ID:             id,
		Name:           req.Name,
		ProjectID:      req.ProjectID,
		Description:    req.Description,
		AutoCalculated: req.AutoCalculated,
		Type:           req.Type,
		Status:         req.Status,
		Category:       req.Category,
	}
	if req.Deadline != "" {
		deadline, err := time.Parse(time.RFC3339, req.Deadline)
		if err != nil {
			ctl.ServerError(c)
			return
		}
		task.Deadline = &deadline
	}

	err = ctl.taskServ.Update(task, authUser.UserID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, task)
}

// ListTasks 列出任务（支持按项目过滤和分页）
func (ctl *TaskController) ListTasks(c *gin.Context) {
	page, pageSize := ctl.parsePagination(c)
	projectIDStr := c.Query("project_id")
	var projectID uint64
	var err error
	if projectIDStr != "" {
		projectID, err = strconv.ParseUint(projectIDStr, 10, 64)
		if err != nil {
			ctl.Code(c, errcode.StatusInvalidParams)
			return
		}
	}

	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	var tasks []models.Task
	var total int64

	if projectID > 0 {
		tasks, total, err = ctl.taskServ.ListByProjectID(projectID, page, pageSize, authUser.UserID)
		if err != nil {
			ctl.Error(c, err)
			return
		}
	} else {
		offset := (page - 1) * pageSize
		tasks, total, err = ctl.taskServ.ListByUserID(authUser.UserID, offset, pageSize)
		if err != nil {
			ctl.Error(c, err)
			return
		}
	}

	ctl.Success(c, gin.H{
		"list":  tasks,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// FinishTask 完成任务
func (ctl *TaskController) FinishTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctl.Code(c, errcode.StatusInvalidParams)
		return
	}
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}
	// 获取现有任务
	task, err := ctl.taskServ.GetByID(id, authUser.UserID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	var req dto.FinishTaskRequest
	if !ctl.BindJSON(c, &req) {
		return
	}
	now, err := time.Parse(time.RFC3339, req.Time)
	if err != nil {
		log.Printf("fail to parse time %s", req.Time)
		ctl.ServerError(c)
		return
	}
	// 假设状态 2 表示已完成，可提取为常量
	task.Status = 2
	task.CompletedAt = &now

	err = ctl.taskServ.Update(task, authUser.UserID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, task)
}

// GetPostrequisites 获取后置任务
func (ctl *TaskController) GetPostrequisites(c *gin.Context) {
	idStr := c.Param("id")
	prerequisiteID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctl.Code(c, errcode.StatusInvalidParams)
		return
	}
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}
	deps, err := ctl.taskServ.GetPostrequisite(prerequisiteID, authUser.UserID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, deps)
}

// GetPrerequisites 获取前置任务
func (ctl *TaskController) GetPrerequisites(c *gin.Context) {
	idStr := c.Param("id")
	taskID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctl.Code(c, errcode.StatusInvalidParams)
		return
	}
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}
	deps, err := ctl.taskServ.GetPrerequisites(taskID, authUser.UserID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, deps)
}

// SetPrerequisites 设置前置任务
func (ctl *TaskController) SetPrerequisites(c *gin.Context) {
	idStr := c.Param("id")
	taskID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctl.Code(c, errcode.StatusInvalidParams)
		return
	}
	var req dto.PrerequisitesRequest
	if !ctl.BindJSON(c, &req) {
		return
	}
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}
	// 先检查任务所有权
	if _, err = ctl.taskServ.GetByID(taskID, authUser.UserID); err != nil {
		ctl.Error(c, err)
		return
	}
	_, err = ctl.taskServ.SetPrerequisiteTask(req.PrerequisiteID, taskID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, gin.H{"message": "prerequisite set successfully"})
}

// UnsetPrerequisites 取消前置任务
func (ctl *TaskController) UnsetPrerequisites(c *gin.Context) {
	idStr := c.Param("id")
	taskID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctl.Code(c, errcode.StatusInvalidParams)
		return
	}
	var req dto.PrerequisitesRequest
	if !ctl.BindJSON(c, &req) {
		return
	}
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}
	err = ctl.taskServ.UnsetPrerequisiteTask(req.PrerequisiteID, taskID, authUser.UserID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, gin.H{"message": "prerequisite unset successfully"})
}

// SetPayment 添加任务付款
func (ctl *TaskController) SetPayment(c *gin.Context) {
	idStr := c.Param("id")
	taskID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctl.Code(c, errcode.StatusInvalidParams)
		return
	}
	var req dto.CreateTaskPaymentRequest
	if !ctl.BindJSON(c, &req) {
		return
	}
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}
	payment := &models.TaskPayment{
		TaskID:   taskID,
		BudgetID: req.BudgetID,
		Amount:   req.Amount,
	}
	err = ctl.taskServ.AddPayment(taskID, payment, authUser.UserID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, payment)
}

// UpdatePayment 更新任务付款
func (ctl *TaskController) UpdatePayment(c *gin.Context) {
	idStr := c.Param("id")
	budgetID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctl.Code(c, errcode.StatusInvalidParams)
		return
	}
	var req dto.UpdateTaskPaymentRequest
	if !ctl.BindJSON(c, &req) {
		return
	}
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}
	payment := &models.TaskPayment{
		ID:     budgetID,
		Amount: req.Amount,
	}
	err = ctl.taskServ.UpdatePayment(payment, authUser.UserID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, payment)
}

// DeletePayment 删除任务付款
func (ctl *TaskController) DeletePayment(c *gin.Context) {
	idStr := c.Param("id")
	budgetID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctl.Code(c, errcode.StatusInvalidParams)
		return
	}
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}
	err = ctl.taskServ.DeletePayment(budgetID, authUser.UserID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, gin.H{"message": "payment deleted successfully"})
}

// GetPayments 获取任务的所有付款
func (ctl *TaskController) GetPayments(c *gin.Context) {
	idStr := c.Param("id")
	taskID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctl.Code(c, errcode.StatusInvalidParams)
		return
	}
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}
	payments, err := ctl.taskServ.GetPaymentByTaskID(taskID, authUser.UserID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, payments)
}
