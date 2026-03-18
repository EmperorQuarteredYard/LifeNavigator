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

type AccountController struct {
	accountServ service.AccountService
	*BaseController
}

func NewAccountController(accountServ service.AccountService) *AccountController {
	return &AccountController{
		accountServ:    accountServ,
		BaseController: &BaseController{},
	}
}

func (ctl *AccountController) CreateAccount(c *gin.Context) {
	var req dto.CreateAccountRequest
	if !ctl.BindJSON(c, &req) {
		return
	}
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	account := &models.Account{
		Name:    req.Name,
		Type:    req.Type,
		Balance: req.Balance,
	}
	created, err := ctl.accountServ.CreateAccount(authUser.UserID, account)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, created)
}

func (ctl *AccountController) DeleteAccount(c *gin.Context) {
	idStr := c.Param("id")
	accountID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctl.Code(c, errcode.StatusInvalidParams)
		return
	}
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	if err := ctl.accountServ.DeleteAccount(authUser.UserID, accountID); err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, gin.H{"message": "account deleted successfully"})
}

func (ctl *AccountController) AdjustBalance(c *gin.Context) {
	idStr := c.Param("id")
	accountID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctl.Code(c, errcode.StatusInvalidParams)
		return
	}
	var req dto.AdjustBalanceRequest
	if !ctl.BindJSON(c, &req) {
		return
	}
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	newBalance, err := ctl.accountServ.AdjustBalance(authUser.UserID, accountID, req.Amount)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, gin.H{"new_balance": newBalance})
}

func (ctl *AccountController) GetAccount(c *gin.Context) {
	idStr := c.Param("id")
	accountID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctl.Code(c, errcode.StatusInvalidParams)
		return
	}
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	account, err := ctl.accountServ.GetByAccountID(authUser.UserID, accountID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, account)
}

func (ctl *AccountController) ListAccounts(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	accounts, err := ctl.accountServ.ListByUserID(authUser.UserID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, accounts)
}

func (ctl *AccountController) ListLinkedTasks(c *gin.Context) {
	idStr := c.Param("id")
	accountID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctl.Code(c, errcode.StatusInvalidParams)
		return
	}
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	startTime, endTime, err := ctl.parseTimeRange(c)
	if err != nil {
		ctl.Code(c, errcode.StatusInvalidParams)
		return
	}

	taskList, err := ctl.accountServ.ListLinkedTask(authUser.UserID, accountID, startTime, endTime)
	if err != nil {
		ctl.Error(c, err)
		return
	}

	ctl.Success(c, taskList)
}

func (ctl *AccountController) parseTimeRange(c *gin.Context) (start, end time.Time, err error) {
	startStr := c.Query("start_time")
	endStr := c.Query("end_time")

	start = time.Unix(0, 0)
	if startStr != "" {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			return
		}
	}

	end = time.Now()
	if endStr != "" {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			return
		}
	}
	return
}
