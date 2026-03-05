package controller

import (
	"LifeNavigator/internal/models"
	"LifeNavigator/internal/service"
	"LifeNavigator/pkg/dto"
	"LifeNavigator/pkg/errcode"
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	BaseController
	taskService service.TaskService
}

func NewTaskController(taskService service.TaskService) *TaskController {
	return &TaskController{
		taskService: taskService,
	}
}

func (ctl *TaskController) GetPostrequisite(c *gin.Context) {
	var req dto.GetPostrequisiteRequest
	_, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}
	if !ctl.BindJSON(c, &req) {
		return
	}
	dependencies, err := ctl.taskService.GetPostrequisite(req.PrerequisiteID)
	if err != nil {
		if errors.Is(err, service.ErrTaskDependencyNotFound) {
			ctl.HandleCode(c, errcode.StatusPrerequisiteNotFound)
			return
		}
		ctl.HandleCode(c, errcode.StatusServerError)
		return
	}
	var responses []dto.DependencyResponse
	for _, item := range dependencies {
		responses = append(responses, dto.DependencyResponse{
			PrerequisiteID: item.PrerequisiteID,
			TaskID:         item.TaskID,
		})
	}
	ctl.Success(c, responses)
}

func (ctl *TaskController) GetPrerequisite(c *gin.Context) {
	var req dto.GetPrerequisiteRequest
	_, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}
	if !ctl.BindJSON(c, &req) {
		return
	}
	dependencies, err := ctl.taskService.GetPrerequisites(req.TaskID)
	if err != nil {
		if errors.Is(err, service.ErrTaskDependencyNotFound) {
			ctl.HandleCode(c, errcode.StatusPrerequisiteNotFound)
			return
		}
		ctl.HandleCode(c, errcode.StatusServerError)
		return
	}
	var responses []dto.DependencyResponse
	for _, item := range dependencies {
		responses = append(responses, dto.DependencyResponse{
			PrerequisiteID: item.PrerequisiteID,
			TaskID:         item.TaskID,
		})
	}
	ctl.Success(c, responses)
}

// UnsetPrerequisiteTask 取消设置前置任务
func (ctl *TaskController) UnsetPrerequisiteTask(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}
	var req dto.SetPrerequisiteRequest
	if !ctl.BindJSON(c, &req) {
		return
	}
	err := ctl.taskService.UnsetPrerequisiteTask(req.PrerequisiteID, req.TaskID, authUser.UserID)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			ctl.HandleCode(c, errcode.StatusInsufficientPerm)
		} else if errors.Is(err, service.ErrTaskNotFound) {
			ctl.HandleCode(c, errcode.StatusTaskNotFound)
		} else if errors.Is(err, service.ErrTaskDependencyNotFound) {
			ctl.HandleCode(c, errcode.StatusPrerequisiteNotFound)
		} else {
			ctl.HandleCode(c, errcode.StatusServerError)
		}
		return
	}
	ctl.Success(c, "success")
}

// SetPrerequisiteTask 设置前置任务
func (ctl *TaskController) SetPrerequisiteTask(c *gin.Context) {
	var req dto.SetPrerequisiteRequest
	if !ctl.BindJSON(c, &req) {
		return
	}
	_, err := ctl.taskService.SetPrerequisiteTask(req.PrerequisiteID, req.TaskID)
	if err != nil {
		ctl.HandleCode(c, errcode.StatusServerError)
		return
	}
	ctl.Success(c, "success")
}

// CreateTask 创建任务
func (ctl *TaskController) CreateTask(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	var req dto.CreateTaskRequest
	if !ctl.BindJSON(c, &req) {
		return
	}

	var deadline *time.Time
	if req.Deadline != nil {
		t, err := time.Parse(time.RFC3339, *req.Deadline)
		if err != nil {
			ctl.HandleCode(c, errcode.StatusInvalidParams)
			return
		}
		deadline = &t
	}

	task := &models.Task{
		ProjectID:      req.ProjectID,
		Name:           req.Name,
		Description:    req.Description,
		AutoCalculated: req.AutoCalculated,
		Type:           req.Type,
		Status:         req.Status,
		Category:       req.Category,
		ForWhom:        req.ForWhom,
		Deadline:       deadline,
	}

	err := ctl.taskService.Create(task, authUser.UserID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrForbidden):
			ctl.HandleCode(c, errcode.StatusInsufficientPermissions)
		default:
			ctl.HandleCode(c, errcode.StatusServerError)
		}
		return
	}
	ctl.Success(c, task)
}

// GetTask 获取任务详情
func (ctl *TaskController) GetTask(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		ctl.HandleCode(c, errcode.StatusInvalidParams)
		return
	}

	task, err := ctl.taskService.GetByID(id, authUser.UserID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrForbidden):
			ctl.HandleCode(c, errcode.StatusInsufficientPermissions)
		case errors.Is(err, service.ErrTaskNotFound):
			ctl.HandleCode(c, errcode.StatusTaskNotFound)
		default:
			ctl.HandleCode(c, errcode.StatusServerError)
		}
		return
	}
	ctl.Success(c, task)
}

// ListProjectTasks 列出指定项目的任务
func (ctl *TaskController) ListProjectTasks(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		ctl.HandleCode(c, errcode.StatusInvalidParams)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	tasks, total, err := ctl.taskService.ListByProjectID(projectID, page, pageSize, authUser.UserID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrForbidden):
			ctl.HandleCode(c, errcode.StatusInsufficientPermissions)
		case errors.Is(err, service.ErrInvalidInput):
			ctl.HandleCode(c, errcode.StatusInvalidParams)
		default:
			ctl.HandleCode(c, errcode.StatusServerError)
		}
		return
	}
	ctl.Success(c, gin.H{
		"tasks": tasks,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// ListUserTasks 列出当前用户的所有任务
func (ctl *TaskController) ListUserTasks(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	tasks, err := ctl.taskService.ListByUserID(authUser.UserID, offset, limit)
	if err != nil {
		ctl.HandleCode(c, errcode.StatusServerError)
		return
	}
	ctl.Success(c, tasks)
}

// UpdateTask 更新任务
func (ctl *TaskController) UpdateTask(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		ctl.HandleCode(c, errcode.StatusInvalidParams)
		return
	}

	var req dto.UpdateTaskRequest
	if !ctl.BindJSON(c, &req) {
		return
	}

	var deadline *time.Time
	if req.Deadline != nil {
		t, err := time.Parse(time.RFC3339, *req.Deadline)
		if err != nil {
			ctl.HandleCode(c, errcode.StatusInvalidParams)
			return
		}
		deadline = &t
	}

	task := &models.Task{
		ID:             id,
		ProjectID:      req.ProjectID,
		Name:           req.Name,
		Description:    req.Description,
		AutoCalculated: req.AutoCalculated,
		Type:           req.Type,
		Status:         req.Status,
		Category:       req.Category,
		ForWhom:        req.ForWhom,
		Deadline:       deadline,
	}

	err = ctl.taskService.Update(task, authUser.UserID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrForbidden):
			ctl.HandleCode(c, errcode.StatusInsufficientPermissions)
		case errors.Is(err, service.ErrTaskNotFound):
			ctl.HandleCode(c, errcode.StatusTaskNotFound)
		default:
			ctl.HandleCode(c, errcode.StatusServerError)
		}
		return
	}
	ctl.Success(c, gin.H{"message": "更新成功"})
}

// DeleteTask 删除任务（级联删除任务预算）
func (ctl *TaskController) DeleteTask(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		ctl.HandleCode(c, errcode.StatusInvalidParams)
		return
	}

	err = ctl.taskService.Delete(id, authUser.UserID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrForbidden):
			ctl.HandleCode(c, errcode.StatusInsufficientPermissions)
		case errors.Is(err, service.ErrTaskNotFound):
			ctl.HandleCode(c, errcode.StatusTaskNotFound)
		default:
			ctl.HandleCode(c, errcode.StatusServerError)
		}
		return
	}
	ctl.Success(c, gin.H{"message": "删除成功"})
}

// GetTasksByStatus 按状态查询项目任务
func (ctl *TaskController) GetTasksByStatus(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		ctl.HandleCode(c, errcode.StatusInvalidParams)
		return
	}

	status, err := strconv.ParseUint(c.Param("status"), 10, 8)
	if err != nil {
		ctl.HandleCode(c, errcode.StatusInvalidParams)
		return
	}

	tasks, err := ctl.taskService.GetByStatus(projectID, uint8(status), authUser.UserID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrForbidden):
			ctl.HandleCode(c, errcode.StatusInsufficientPermissions)
		default:
			ctl.HandleCode(c, errcode.StatusServerError)
		}
		return
	}
	ctl.Success(c, tasks)
}

// GetTasksByTimePeriod 查询指定时间区间内的任务
func (ctl *TaskController) GetTasksByTimePeriod(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		ctl.HandleCode(c, errcode.StatusInvalidParams)
		return
	}

	startStr := c.Query("start")
	endStr := c.Query("end")
	if startStr == "" || endStr == "" {
		ctl.HandleCode(c, errcode.StatusInvalidParams)
		return
	}
	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		ctl.HandleCode(c, errcode.StatusInvalidParams)
		return
	}
	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		ctl.HandleCode(c, errcode.StatusInvalidParams)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	tasks, total, err := ctl.taskService.GetByTimePeriod(projectID, start, end, page, pageSize, authUser.UserID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrForbidden):
			ctl.HandleCode(c, errcode.StatusInsufficientPermissions)
		case errors.Is(err, service.ErrInvalidInput):
			ctl.HandleCode(c, errcode.StatusInvalidParams)
		default:
			ctl.HandleCode(c, errcode.StatusServerError)
		}
		return
	}
	ctl.Success(c, gin.H{
		"tasks": tasks,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// AddTaskBudget 添加任务预算
func (ctl *TaskController) AddTaskBudget(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		ctl.HandleCode(c, errcode.StatusInvalidParams)
		return
	}

	var req dto.TaskBudgetRequest
	if !ctl.BindJSON(c, &req) {
		return
	}

	budget := &models.TaskBudget{
		Type:   req.Type,
		Budget: req.Budget,
		Used:   req.Used,
	}

	err = ctl.taskService.AddBudget(taskID, budget, authUser.UserID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrForbidden):
			ctl.HandleCode(c, errcode.StatusInsufficientPermissions)
		default:
			ctl.HandleCode(c, errcode.StatusServerError)
		}
		return
	}
	ctl.Success(c, budget)
}

// UpdateTaskBudget 更新任务预算
func (ctl *TaskController) UpdateTaskBudget(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	budgetID, err := strconv.ParseUint(c.Param("budgetId"), 10, 64)
	if err != nil {
		ctl.HandleCode(c, errcode.StatusInvalidParams)
		return
	}

	var req dto.TaskBudgetRequest
	if !ctl.BindJSON(c, &req) {
		return
	}

	budget := &models.TaskBudget{
		ID:     budgetID,
		Type:   req.Type,
		Budget: req.Budget,
		Used:   req.Used,
	}

	err = ctl.taskService.UpdateBudget(budget, authUser.UserID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrForbidden):
			ctl.HandleCode(c, errcode.StatusInsufficientPermissions)
		case errors.Is(err, service.ErrBudgetNotFound):
			ctl.HandleCode(c, errcode.StatusBudgetNotFound)
		default:
			ctl.HandleCode(c, errcode.StatusServerError)
		}
		return
	}
	ctl.Success(c, gin.H{"message": "更新成功"})
}

// DeleteTaskBudget 删除任务预算
func (ctl *TaskController) DeleteTaskBudget(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	budgetID, err := strconv.ParseUint(c.Param("budgetId"), 10, 64)
	if err != nil {
		ctl.HandleCode(c, errcode.StatusInvalidParams)
		return
	}

	err = ctl.taskService.DeleteBudget(budgetID, authUser.UserID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrForbidden):
			ctl.HandleCode(c, errcode.StatusInsufficientPermissions)
		case errors.Is(err, service.ErrBudgetNotFound):
			ctl.HandleCode(c, errcode.StatusBudgetNotFound)
		default:
			ctl.HandleCode(c, errcode.StatusServerError)
		}
		return
	}
	ctl.Success(c, gin.H{"message": "删除成功"})
}

// GetTaskBudgets 获取任务预算列表
func (ctl *TaskController) GetTaskBudgets(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		ctl.HandleCode(c, errcode.StatusInvalidParams)
		return
	}

	budgets, err := ctl.taskService.GetBudgetByTaskID(taskID, authUser.UserID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrForbidden):
			ctl.HandleCode(c, errcode.StatusInsufficientPermissions)
		case errors.Is(err, service.ErrTaskNotFound):
			ctl.HandleCode(c, errcode.StatusTaskNotFound)
		default:
			ctl.HandleCode(c, errcode.StatusServerError)
		}
		return
	}
	ctl.Success(c, budgets)
}
