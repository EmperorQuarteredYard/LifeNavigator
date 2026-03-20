package controller

import (
	"LifeNavigator/internal/models"
	"LifeNavigator/internal/service"
	"LifeNavigator/pkg/dto"
	"LifeNavigator/pkg/errcode"
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

func (ctl *TaskController) CreateTask(c *gin.Context) {
	var req dto.CreateTaskRequest
	if !ctl.BindJSON(c, &req) {
		return
	}
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	var deadline *time.Time
	if req.Deadline != nil && *req.Deadline != "" {
		t, err := time.Parse(time.RFC3339, *req.Deadline)
		if err != nil {
			ctl.Code(c, errcode.StatusInvalidParams)
			return
		}
		deadline = &t
	}

	task := &models.Task{
		Name:        req.Name,
		ProjectID:   req.ProjectID,
		Description: req.Description,
		Type:        req.Type,
		Status:      req.Status,
		Category:    req.Category,
		Deadline:    deadline,
	}
	created, err := ctl.taskServ.Create(authUser.UserID, task)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, created)
}

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

	err = ctl.taskServ.Delete(authUser.UserID, id)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, gin.H{"message": "task deleted successfully"})
}

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

	task, err := ctl.taskServ.GetByID(authUser.UserID, id)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, task)
}

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

	var deadline *time.Time
	if req.Deadline != "" {
		t, err := time.Parse(time.RFC3339, req.Deadline)
		if err != nil {
			ctl.Code(c, errcode.StatusInvalidParams)
			return
		}
		deadline = &t
	}

	task := &models.Task{
		ID:          id,
		Name:        req.Name,
		ProjectID:   req.ProjectID,
		Description: req.Description,
		Type:        req.Type,
		Status:      req.Status,
		Category:    req.Category,
		Deadline:    deadline,
	}

	err = ctl.taskServ.Update(authUser.UserID, task)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, gin.H{"message": "task updated successfully"})
}

func (ctl *TaskController) ListTasks(c *gin.Context) { // TODO 这里应当增加一个开始，结束时间的筛选，如果没有指定则为在此方向上无边界
	page, pageSize := ctl.parsePagination(c)
	projectIDStr := c.Query("project_id")
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	if projectIDStr != "" {
		projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
		if err != nil {
			ctl.Code(c, errcode.StatusInvalidParams)
			return
		}
		tasks, err := ctl.taskServ.ListByProjectID(authUser.UserID, projectID, page, pageSize)
		if err != nil {
			ctl.Error(c, err)
			return
		}
		ctl.Success(c, tasks)
		return
	}

	offset := (page - 1) * pageSize
	tasks, err := ctl.taskServ.ListByUserID(authUser.UserID, offset, pageSize)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, tasks)
}

func (ctl *TaskController) UpdateTaskStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctl.Code(c, errcode.StatusInvalidParams)
		return
	}
	var req dto.UpdateTaskStatusRequest
	if !ctl.BindJSON(c, &req) {
		return
	}
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	err = ctl.taskServ.UpdateStatus(authUser.UserID, id, req.Status)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, gin.H{"message": "task status updated successfully"})
}

func (ctl *TaskController) GetPostrequisites(c *gin.Context) {
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
	deps, err := ctl.taskServ.GetPostrequisite(authUser.UserID, taskID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, deps)
}

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
	deps, err := ctl.taskServ.GetPrerequisites(authUser.UserID, taskID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, deps)
}

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
	_, err = ctl.taskServ.SetPrerequisiteTask(authUser.UserID, req.PrerequisiteID, taskID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, gin.H{"message": "prerequisite set successfully"})
}

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
	err = ctl.taskServ.UnsetPrerequisiteTask(authUser.UserID, req.PrerequisiteID, taskID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, gin.H{"message": "prerequisite unset successfully"})
}
