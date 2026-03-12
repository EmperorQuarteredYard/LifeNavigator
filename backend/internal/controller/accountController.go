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

// CreateAccount 创建账户
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
		UserID:  authUser.UserID,
		Type:    req.Type,
		Balance: req.Balance, // 创建时可指定初始余额，服务层可做校验
	}
	created, err := ctl.accountServ.CreateAccount(account)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, created)
}

// DeleteAccount 删除账户
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

	// 先获取账户（确保所有权）
	account, err := ctl.accountServ.GetByAccountID(authUser.UserID, accountID)
	if err != nil {
		ctl.Error(c, err)
		return
	}

	if err := ctl.accountServ.DeleteAccount(account); err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, gin.H{"message": "account deleted successfully"})
}

// AdjustBalance 调整账户余额
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

// GetAccount 获取单个账户详情
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
	var res = dto.AccountModel{
		ID:         account.ID,
		UserID:     account.UserID,
		Name:       account.Name,
		Type:       account.Type,
		Balance:    account.Balance,
		NetBalance: account.NetBalance,
	}
	ctl.Success(c, res)
}

// ListAccounts 获取当前用户的所有账户
func (ctl *AccountController) ListAccounts(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	accounts, err := ctl.accountServ.ListByUserID(authUser.UserID)
	var result []dto.AccountModel
	for _, account := range accounts {
		result = append(result, dto.AccountModel{
			ID:         account.ID,
			UserID:     account.UserID,
			Name:       account.Name,
			Type:       account.Type,
			Balance:    account.Balance,
			NetBalance: account.NetBalance,
		})
	}
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, result)
}

// ListLinkedTasks 获取与账户关联的任务及付款
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

	// 解析时间参数
	startTime, endTime, err := ctl.parseTimeRange(c)
	if err != nil {
		ctl.Code(c, errcode.StatusInvalidParams)
		return
	}

	tasks, payments, err := ctl.accountServ.ListLinkedTask(authUser.UserID, accountID, startTime, endTime)
	if err != nil {
		ctl.Error(c, err)
		return
	}

	ctl.Success(c, gin.H{
		"tasks":    tasks,
		"payments": payments,
	})
}

// ListLinkedPayments 获取与账户关联的付款（不含任务详情）
func (ctl *AccountController) ListLinkedPayments(c *gin.Context) {
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

	payments, err := ctl.accountServ.ListLinkedTaskPayment(authUser.UserID, accountID, startTime, endTime)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, payments)
}

// parseTimeRange 从查询参数解析 start_time 和 end_time，格式 RFC3339
// 如果参数缺失，返回默认时间范围（如 1970-01-01 到 现在）
func (ctl *AccountController) parseTimeRange(c *gin.Context) (start, end time.Time, err error) {
	startStr := c.Query("start_time")
	endStr := c.Query("end_time")

	// 默认起始时间：很久以前
	start = time.Unix(0, 0)
	if startStr != "" {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			return
		}
	}

	// 默认结束时间：当前时间
	end = time.Now()
	if endStr != "" {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			return
		}
	}
	return
}
