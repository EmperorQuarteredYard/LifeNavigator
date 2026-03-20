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
	created, err := ctl.projectServ.Create(authUser.UserID, project)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, created)
}

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

	project, err := ctl.projectServ.GetByID(authUser.UserID, id)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, project)
}

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
	ctl.Success(c, projects)
}

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
	err = ctl.projectServ.Update(authUser.UserID, project)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, gin.H{"message": "project updated successfully"})
}

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

	err = ctl.projectServ.Delete(authUser.UserID, id)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, gin.H{"message": "project deleted successfully"})
}
