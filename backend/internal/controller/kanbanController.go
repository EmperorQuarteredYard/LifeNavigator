package controller

import (
	"LifeNavigator/internal/service"
	"LifeNavigator/pkg/dto"
	"LifeNavigator/pkg/errcode"
	"strconv"

	"github.com/gin-gonic/gin"
)

type KanbanController struct {
	kanbanServ service.KanbanService
	*BaseController
}

func NewKanbanController(kanbanServ service.KanbanService) *KanbanController {
	return &KanbanController{
		kanbanServ:     kanbanServ,
		BaseController: &BaseController{},
	}
}

func (ctl *KanbanController) CreateKanban(c *gin.Context) {
	var req dto.CreateKanbanRequest
	if !ctl.BindJSON(c, &req) {
		return
	}
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	kanban, err := ctl.kanbanServ.Create(authUser.UserID, &req)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, kanban)
}

func (ctl *KanbanController) GetKanban(c *gin.Context) {
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

	kanban, err := ctl.kanbanServ.GetByID(authUser.UserID, id)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, kanban)
}

func (ctl *KanbanController) ListKanbans(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	list, err := ctl.kanbanServ.ListByUserID(authUser.UserID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, list)
}

func (ctl *KanbanController) UpdateKanban(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctl.Code(c, errcode.StatusInvalidParams)
		return
	}
	var req dto.UpdateKanbanRequest
	if !ctl.BindJSON(c, &req) {
		return
	}
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	err = ctl.kanbanServ.Update(authUser.UserID, id, &req)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, gin.H{"message": "kanban updated successfully"})
}

func (ctl *KanbanController) DeleteKanban(c *gin.Context) {
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

	err = ctl.kanbanServ.Delete(authUser.UserID, id)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, gin.H{"message": "kanban deleted successfully"})
}

func (ctl *KanbanController) GetKanbanTasks(c *gin.Context) {
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

	page, pageSize := ctl.parsePagination(c)

	var status *uint8
	statusStr := c.Query("status")
	if statusStr != "" {
		s, err := strconv.ParseUint(statusStr, 10, 8)
		if err != nil {
			ctl.Code(c, errcode.StatusInvalidParams)
			return
		}
		st := uint8(s)
		status = &st
	}

	tasks, err := ctl.kanbanServ.GetKanbanTasks(authUser.UserID, id, status, page, pageSize)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, tasks)
}

func (ctl *KanbanController) GetDefaultKanbanTasks(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	page, pageSize := ctl.parsePagination(c)

	var status *uint8
	statusStr := c.Query("status")
	if statusStr != "" {
		s, err := strconv.ParseUint(statusStr, 10, 8)
		if err != nil {
			ctl.Code(c, errcode.StatusInvalidParams)
			return
		}
		st := uint8(s)
		status = &st
	}

	tasks, err := ctl.kanbanServ.GetDefaultKanbanTasks(authUser.UserID, status, page, pageSize)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, tasks)
}

func (ctl *KanbanController) SetDefaultKanban(c *gin.Context) {
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

	err = ctl.kanbanServ.SetDefaultKanban(authUser.UserID, id)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, gin.H{"message": "default kanban set successfully"})
}
