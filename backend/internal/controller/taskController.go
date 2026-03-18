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

func (ctl *TaskController) ListTasks(c *gin.Context) {
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

	task, err := ctl.taskServ.GetByID(authUser.UserID, id)
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
		ctl.Code(c, errcode.StatusInvalidParams)
		return
	}

	updateTask := &models.Task{
		ID:          task.ID,
		ProjectID:   task.ProjectID,
		Name:        task.Name,
		Description: task.Description,
		Type:        task.Type,
		Status:      2,
		Category:    task.Category,
		Deadline:    task.Deadline,
		CompletedAt: &now,
	}

	err = ctl.taskServ.Update(authUser.UserID, updateTask)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, gin.H{"message": "task finished successfully"})
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
	err = ctl.taskServ.AddPayment(authUser.UserID, taskID, payment)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, gin.H{"message": "payment added successfully", "id": payment.ID})
}

func (ctl *TaskController) UpdatePayment(c *gin.Context) {
	idStr := c.Param("id")
	paymentID, err := strconv.ParseUint(idStr, 10, 64)
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
		ID:     paymentID,
		Amount: req.Amount,
	}
	err = ctl.taskServ.UpdatePayment(authUser.UserID, payment)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, gin.H{"message": "payment updated successfully"})
}

func (ctl *TaskController) DeletePayment(c *gin.Context) {
	idStr := c.Param("id")
	paymentID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctl.Code(c, errcode.StatusInvalidParams)
		return
	}
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}
	err = ctl.taskServ.DeletePayment(authUser.UserID, paymentID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, gin.H{"message": "payment deleted successfully"})
}

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
	payments, err := ctl.taskServ.GetPaymentByTaskID(authUser.UserID, taskID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, payments)
}
