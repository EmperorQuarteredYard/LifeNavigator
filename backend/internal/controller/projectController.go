package controller

import (
	"LifeNavigator/internal/models"
	"LifeNavigator/internal/service"
	"LifeNavigator/pkg/dto"
	"LifeNavigator/pkg/errcode"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProjectController struct {
	BaseController
	projectService service.ProjectService
}

func NewProjectController(projectService service.ProjectService) *ProjectController {
	return &ProjectController{
		projectService: projectService,
	}
}

// CreateProject 创建项目
func (ctl *ProjectController) CreateProject(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	var req dto.CreateProjectRequest
	if !ctl.BindJSON(c, &req) {
		return
	}

	project := &models.Project{
		Name:            req.Name,
		Description:     req.Description,
		RefreshInterval: req.RefreshInterval,
	}

	if err := ctl.projectService.Create(project, authUser.UserID); err != nil {
		ctl.HandleCode(c, errcode.StatusServerError)
		return
	}
	ctl.Success(c, project)
}

// GetProject 获取项目详情
func (ctl *ProjectController) GetProject(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		ctl.HandleCode(c, errcode.StatusInvalidParams)
		return
	}

	project, err := ctl.projectService.GetByID(id, authUser.UserID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrForbidden):
			ctl.HandleCode(c, errcode.StatusInsufficientPermissions)
		case errors.Is(err, service.ErrProjectNotFound):
			ctl.HandleCode(c, errcode.StatusProjectNotFound)
		default:
			ctl.HandleCode(c, errcode.StatusServerError)
		}
		return
	}
	ctl.Success(c, project)
}

// ListProjects 列出当前用户的所有项目
func (ctl *ProjectController) ListProjects(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	projects, err := ctl.projectService.ListByUserID(authUser.UserID, offset, limit)
	if err != nil {
		ctl.HandleCode(c, errcode.StatusServerError)
		return
	}
	ctl.Success(c, projects)
}

// UpdateProject 更新项目
func (ctl *ProjectController) UpdateProject(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		ctl.HandleCode(c, errcode.StatusInvalidParams)
		return
	}

	var req dto.UpdateProjectRequest
	if !ctl.BindJSON(c, &req) {
		return
	}

	project := &models.Project{
		ID:              id,
		Name:            req.Name,
		Description:     req.Description,
		RefreshInterval: req.RefreshInterval,
	}

	err = ctl.projectService.Update(project, authUser.UserID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrForbidden):
			ctl.HandleCode(c, errcode.StatusInsufficientPermissions)
		case errors.Is(err, service.ErrProjectNotFound):
			ctl.HandleCode(c, errcode.StatusProjectNotFound)
		default:
			ctl.HandleCode(c, errcode.StatusServerError)
		}
		return
	}
	ctl.Success(c, nil)
}

// DeleteProject 删除项目（级联删除预算和任务预算）
func (ctl *ProjectController) DeleteProject(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		ctl.HandleCode(c, errcode.StatusInvalidParams)
		return
	}

	err = ctl.projectService.Delete(id, authUser.UserID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrForbidden):
			ctl.HandleCode(c, errcode.StatusInsufficientPermissions)
		case errors.Is(err, service.ErrProjectNotFound):
			ctl.HandleCode(c, errcode.StatusProjectNotFound)
		default:
			ctl.HandleCode(c, errcode.StatusServerError)
		}
		return
	}
	ctl.Success(c, nil)
}

// AddProjectBudget 添加项目预算
func (ctl *ProjectController) AddProjectBudget(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		ctl.HandleCode(c, errcode.StatusInvalidParams)
		return
	}

	var req dto.ProjectBudgetRequest
	if !ctl.BindJSON(c, &req) {
		return
	}

	budget := &models.ProjectBudget{
		Type:   req.Type,
		Budget: req.Budget,
		Used:   req.Used,
	}

	err = ctl.projectService.AddBudget(projectID, budget, authUser.UserID)
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

// UpdateProjectBudget 更新项目预算
func (ctl *ProjectController) UpdateProjectBudget(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	budgetID, err := strconv.ParseUint(c.Param("budgetId"), 10, 64)
	if err != nil {
		ctl.HandleCode(c, errcode.StatusInvalidParams)
		return
	}

	var req dto.ProjectBudgetRequest
	if !ctl.BindJSON(c, &req) {
		return
	}

	budget := &models.ProjectBudget{
		ID:     budgetID,
		Type:   req.Type,
		Budget: req.Budget,
		Used:   req.Used,
	}

	err = ctl.projectService.UpdateBudget(budget, authUser.UserID)
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

// DeleteProjectBudget 删除项目预算
func (ctl *ProjectController) DeleteProjectBudget(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	budgetID, err := strconv.ParseUint(c.Param("budgetId"), 10, 64)
	if err != nil {
		ctl.HandleCode(c, errcode.StatusInvalidParams)
		return
	}

	err = ctl.projectService.DeleteBudget(budgetID, authUser.UserID)
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

// GetProjectBudgetSummary 获取项目预算汇总
func (ctl *ProjectController) GetProjectBudgetSummary(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		ctl.HandleCode(c, errcode.StatusInvalidParams)
		return
	}

	budgets, totalBudget, totalUsed, err := ctl.projectService.GetBudgetSummary(projectID, authUser.UserID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrForbidden):
			ctl.HandleCode(c, errcode.StatusInsufficientPermissions)
		case errors.Is(err, service.ErrProjectNotFound):
			ctl.HandleCode(c, errcode.StatusProjectNotFound)
		default:
			ctl.HandleCode(c, errcode.StatusServerError)
		}
		return
	}
	ctl.Success(c, gin.H{
		"budgets":      budgets,
		"total_budget": totalBudget,
		"total_used":   totalUsed,
	})
}

// GetTaskBudgetSummary 获取项目下所有任务的预算汇总
func (ctl *ProjectController) GetTaskBudgetSummary(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		ctl.HandleCode(c, errcode.StatusInvalidParams)
		return
	}

	budgets, err := ctl.projectService.GetTaskBudgetSummary(projectID, authUser.UserID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrForbidden):
			ctl.HandleCode(c, errcode.StatusInsufficientPermissions)
		case errors.Is(err, service.ErrProjectNotFound):
			ctl.HandleCode(c, errcode.StatusProjectNotFound)
		default:
			ctl.HandleCode(c, errcode.StatusServerError)
		}
		return
	}
	ctl.Success(c, budgets)
}
