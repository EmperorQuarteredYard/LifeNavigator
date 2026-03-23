package controller

import (
	"LifeNavigator/internal/interfaces/Service"
	"LifeNavigator/internal/models"
	"LifeNavigator/pkg/dto"
	"LifeNavigator/pkg/errcode"
	"strconv"

	"github.com/gin-gonic/gin"
)

func NewBudgetController(budgetServ Service.BudgetService) *BudgetController {
	return &BudgetController{
		budgetServ: budgetServ,
	}
}

type BudgetController struct {
	budgetServ Service.BudgetService
	*BaseController
}

func (ctl *BudgetController) SetPayment(c *gin.Context) {
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
	err = ctl.budgetServ.AddPayment(authUser.UserID, taskID, payment)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, gin.H{"message": "payment added successfully", "id": payment.ID})
}
func (ctl *BudgetController) UpdatePayment(c *gin.Context) {
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
	err = ctl.budgetServ.UpdatePayment(authUser.UserID, payment)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, gin.H{"message": "payment updated successfully"})
}

func (ctl *BudgetController) DeletePayment(c *gin.Context) {
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
	err = ctl.budgetServ.DeletePayment(authUser.UserID, paymentID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, gin.H{"message": "payment deleted successfully"})
}

func (ctl *BudgetController) GetPayments(c *gin.Context) {
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
	payments, err := ctl.budgetServ.GetPaymentByTaskID(authUser.UserID, taskID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, payments)
}

func (ctl *BudgetController) AddBudget(c *gin.Context) {
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
		ProjectID: projectID,
		AccountID: req.AccountID,
		Unit:      req.Type,
		Budget:    req.Budget,
		Used:      req.Used,
	}
	err = ctl.budgetServ.AddBudget(authUser.UserID, projectID, budget)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, gin.H{"message": "budget added successfully", "id": budget.ID})
}
func (ctl *BudgetController) UpdateBudget(c *gin.Context) {
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
		Unit:      req.Type,
		Budget:    req.Budget,
		Used:      req.Used,
	}
	err = ctl.budgetServ.UpdateBudget(authUser.UserID, budget)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, gin.H{"message": "budget updated successfully"})
}
func (ctl *BudgetController) DeleteBudget(c *gin.Context) {
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

	err = ctl.budgetServ.DeleteBudget(authUser.UserID, budgetID)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Success(c, gin.H{"message": "budget deleted successfully"})
}
