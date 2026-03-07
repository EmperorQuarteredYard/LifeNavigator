package controller

import (
	"LifeNavigator/internal/models"
	"LifeNavigator/internal/service"
	"LifeNavigator/pkg/dto"
	"LifeNavigator/pkg/errcode"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProjectController struct {
	projectServ service.ProjectService
	*BaseController
}

func NewProjectController(projectServ service.ProjectService) *ProjectController {
	return &ProjectController{
		projectServ:    projectServ,
		BaseController: &BaseController{},
	}
}

// CreateProject 创建项目
func (ctl *ProjectController) CreateProject(c *gin.Context) {
	var req dto.CreateProjectRequest
	if !ctl.BindJSON(c, &req) {
		return
	}
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	project := &models.Project{
		Name:            req.Name,
		Description:     req.Description,
		RefreshInterval: req.RefreshInterval,
	}
	err := ctl.projectServ.Create(project, authUser.UserID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, project)
}

// GetProject 获取单个项目详情
func (ctl *ProjectController) GetProject(c *gin.Context) {
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

	project, err := ctl.projectServ.GetByID(id, authUser.UserID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, project)
}

// GetProjectsByUser 获取当前用户的项目列表（分页）
func (ctl *ProjectController) GetProjectsByUser(c *gin.Context) {
	page, pageSize := ctl.parsePagination(c)
	offset := (page - 1) * pageSize
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	projects, err := ctl.projectServ.ListByUserID(authUser.UserID, offset, pageSize)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	// 注意：ListByUserID 只返回列表，没有总数，可能需要扩展服务层返回总数
	// 这里简化处理，只返回列表，前端可自行决定是否显示总数
	ctl.Success(c, gin.H{
		"list": projects,
		"page": page,
		"size": pageSize,
	})
}

// UpdateProject 更新项目
func (ctl *ProjectController) UpdateProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctl.Code(c, errcode.StatusInvalidParams)
		return
	}
	var req dto.UpdateProjectRequest
	if !ctl.BindJSON(c, &req) {
		return
	}
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	project := &models.Project{
		ID:              id,
		Name:            req.Name,
		Description:     req.Description,
		RefreshInterval: req.RefreshInterval,
	}
	err = ctl.projectServ.Update(project, authUser.UserID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, project)
}

// DeleteProject 删除项目（级联删除相关预算和任务）
func (ctl *ProjectController) DeleteProject(c *gin.Context) {
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

	err = ctl.projectServ.Delete(id, authUser.UserID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, gin.H{"message": "project deleted successfully"})
}

// AddBudget 添加项目预算
func (ctl *ProjectController) AddBudget(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		ctl.Code(c, errcode.StatusInvalidParams)
		return
	}
	var req dto.ProjectBudgetRequest
	if !ctl.BindJSON(c, &req) {
		return
	}
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	budget := &models.ProjectBudget{
		AccountID: req.AccountID,
		Budget:    req.Budget,
		Used:      req.Used,
	}
	err = ctl.projectServ.AddBudget(projectID, budget, authUser.UserID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, budget)
}

// UpdateBudget 更新项目预算
func (ctl *ProjectController) UpdateBudget(c *gin.Context) {
	budgetIDStr := c.Param("budgetId")
	budgetID, err := strconv.ParseUint(budgetIDStr, 10, 64)
	if err != nil {
		ctl.Code(c, errcode.StatusInvalidParams)
		return
	}
	var req dto.ProjectBudgetRequest
	if !ctl.BindJSON(c, &req) {
		return
	}
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	budget := &models.ProjectBudget{
		ID:        budgetID,
		AccountID: req.AccountID,
		Budget:    req.Budget,
		Used:      req.Used,
	}
	err = ctl.projectServ.UpdateBudget(budget, authUser.UserID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, budget)
}

// DeleteBudget 删除项目预算
func (ctl *ProjectController) DeleteBudget(c *gin.Context) {
	budgetIDStr := c.Param("budgetId")
	budgetID, err := strconv.ParseUint(budgetIDStr, 10, 64)
	if err != nil {
		ctl.Code(c, errcode.StatusInvalidParams)
		return
	}
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	err = ctl.projectServ.DeleteBudget(budgetID, authUser.UserID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, gin.H{"message": "budget deleted successfully"})
}

// GetBudgetSummary 获取项目预算汇总
func (ctl *ProjectController) GetBudgetSummary(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		ctl.Code(c, errcode.StatusInvalidParams)
		return
	}
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	budgets, totalBudget, totalUsed, err := ctl.projectServ.GetBudgetSummary(projectID, authUser.UserID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, gin.H{
		"budgets":      budgets,
		"total_budget": totalBudget,
		"total_used":   totalUsed,
	})
}
