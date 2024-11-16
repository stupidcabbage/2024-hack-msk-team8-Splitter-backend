package controllers

import (
	"fmt"
	"strconv"

	"example.com/m/internal/api/v1/core/application/dto"
	"example.com/m/internal/api/v1/core/application/services/debt_service"
	"example.com/m/internal/api/v1/utils"
	"github.com/gin-gonic/gin"
)

type DebtController struct {
	ds debt_service.DebtService
}

func NewDebtController(s *debt_service.DebtService) *DebtController {
	return &DebtController{
		ds: *s,
	}
}

// @BasePath /api/v1

// Create debt
// @Schemes
// @Description Creates debt and returns it
// @Accept json
// @Produce json
// @Tags debt
// @Param debt body dto.CreateDebtDto true "Debt data"
// @Success 201 {array} dto.DebtDto
// @Failure 500 {object} errorz.Error_
// @Failure 503 {object} errorz.Error_
// @Failure 401 {object} errorz.Error_
// @Failure 400 {object} errorz.Error_
// @Security BearerAuth
// @Param Authorization header string true "Bearer JWT token" // Добавляем параметр для JWT
// @Router /debts [post]
func (c *DebtController) CreateDebt(ctx *gin.Context) {
	token, err := utils.ExtractTokenFromHeaders(ctx)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	payload, err := utils.ExtractPayloadFromJWT(*token)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	username := payload["username"].(string)

	var debt dto.CreateDebtDto
	if err := ctx.ShouldBindBodyWithJSON(&debt); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	debtToCreate := dto.CreateDebtDto{
		Name:        debt.Name,
		ToUsername:  username,
		InviteCode:  debt.InviteCode,
		TotalAmount: debt.TotalAmount,
		GroupId:     debt.GroupId,
	}

	createdDebt, err := c.ds.CreateDebt(ctx, &debtToCreate)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	ctx.JSON(201, &createdDebt)
}

// Get outcoming debts
// @Schemes
// @Description Returns user's outcoming debts(requires JWT in "Bearer" header)
// @Tags debt
// @Produce json
// @Success 200 {array} dto.DebtDto
// @Failure 500 {object} errorz.Error_
// @Failure 503 {object} errorz.Error_
// @Failure 401 {object} errorz.Error_
// @Security BearerAuth
// @Param Authorization header string true "Bearer JWT token" // Добавляем параметр для JWT
// @Router /debts/outcoming [get]
func (c *DebtController) GetOutcomingDebts(ctx *gin.Context) {
	token, err := utils.ExtractTokenFromHeaders(ctx)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	payload, err := utils.ExtractPayloadFromJWT(*token)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	username := payload["username"].(string)

	debts, err := c.ds.GetOutcomingDebts(ctx, username)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}
	if len(*debts) == 0 {
		ctx.JSON(200, make([]string, 0))
		return
	}

	ctx.JSON(200, *debts)
}

// Get incoming debts
// @Schemes
// @Description Returns user's incoming debts(requires JWT in "Bearer" header)
// @Tags debt
// @Produce json
// @Success 200 {array} dto.DebtDto
// @Failure 500 {object} errorz.Error_
// @Failure 503 {object} errorz.Error_
// @Failure 401 {object} errorz.Error_
// @Security BearerAuth
// @Param Authorization header string true "Bearer JWT token" // Добавляем параметр для JWT
// @Router /debts/incoming [get]
func (c *DebtController) GetIncomingDebts(ctx *gin.Context) {
	token, err := utils.ExtractTokenFromHeaders(ctx)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	payload, err := utils.ExtractPayloadFromJWT(*token)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	username := payload["username"].(string)

	debts, err := c.ds.GetIncomingDebts(ctx, username)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	if len(*debts) == 0 {
		ctx.JSON(200, make([]string, 0))
		return
	}

	ctx.JSON(200, *debts)
}

// Get closed debts
// @Schemes
// @Description Returns user's closed debts (outcoming and incoming) (requires JWT in "Bearer" header)
// @Tags debt
// @Produce json
// @Success 200 {array} dto.DebtDto
// @Failure 500 {object} errorz.Error_
// @Failure 503 {object} errorz.Error_
// @Failure 401 {object} errorz.Error_
// @Security BearerAuth
// @Param Authorization header string true "Bearer JWT token" // Добавляем параметр для JWT
// @Router /debts/closed [get]
func (c *DebtController) GetClosedDebts(ctx *gin.Context) {
	token, err := utils.ExtractTokenFromHeaders(ctx)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	payload, err := utils.ExtractPayloadFromJWT(*token)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	username := payload["username"].(string)

	debts, err := c.ds.GetClosedDebts(ctx, username)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}
	fmt.Println(debts)

	ctx.JSON(200, &debts)
}

// Get debt by id
// @Schemes
// @Description Returns debt object (outcoming and incoming) (requires JWT in "Bearer" header)
// @Tags debt
// @Produce json
// @Success 200 {object} dto.DebtDto
// @Failure 500 {object} errorz.Error_
// @Failure 503 {object} errorz.Error_
// @Failure 401 {object} errorz.Error_
// @Security BearerAuth
// @Param Authorization header string true "Bearer JWT token" // Добавляем параметр для JWT
// @Router /debts/{id} [get]
func (c *DebtController) GetDebtById(ctx *gin.Context) {
	token, err := utils.ExtractTokenFromHeaders(ctx)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	payload, err := utils.ExtractPayloadFromJWT(*token)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	username := payload["username"].(string)

	debtId, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)

	debt, err := c.ds.GetDebtById(ctx, username, debtId)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	ctx.JSON(200, &debt)
}

// Close debt
// @Schemes
// @Description Closes debt (requires JWT in "Bearer" header)
// @Tags debt
// @Produce json
// @Success 200 {object} bool
// @Failure 500 {object} errorz.Error_
// @Failure 503 {object} errorz.Error_
// @Failure 401 {object} errorz.Error_
// @Failure 400 {object} errorz.Error_
// @Security BearerAuth
// @Param Authorization header string true "Bearer JWT token" // Добавляем параметр для JWT
// @Router /debts/{id}/close [patch]
func (c *DebtController) CloseDebt(ctx *gin.Context) {
	token, err := utils.ExtractTokenFromHeaders(ctx)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	payload, err := utils.ExtractPayloadFromJWT(*token)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	username := payload["username"].(string)

	debtId, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)

	err = c.ds.CloseDebt(ctx, username, debtId)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
	})
}

// Increase debt payed amount
// @Schemes
// @Description Increases debt current payed amount (requires JWT in "Bearer" header)
// @Tags debt
// @Produce json
// @Param debt body dto.IncreaseDebtPayedAmountDto true "Amount data"
// @Success 200 {object} bool
// @Failure 500 {object} errorz.Error_
// @Failure 503 {object} errorz.Error_
// @Failure 401 {object} errorz.Error_
// @Failure 400 {object} errorz.Error_
// @Security BearerAuth
// @Param Authorization header string true "Bearer JWT token" // Добавляем параметр для JWT
// @Router /debts/{id}/payedAmount/increase [put]
func (c *DebtController) IncreaseDebtPayedAmount(ctx *gin.Context) {
	token, err := utils.ExtractTokenFromHeaders(ctx)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	payload, err := utils.ExtractPayloadFromJWT(*token)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	var amount dto.IncreaseDebtPayedAmountDto
	if err := ctx.ShouldBindBodyWithJSON(&amount); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	username := payload["username"].(string)

	debtId, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)

	isClosed, err := c.ds.IncreaseDebtPayedAmount(ctx, username, debtId, amount.Amount)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	ctx.JSON(200, gin.H{
		"is_closed": isClosed,
	})
}

// Get debts stats
// @Schemes
// @Description Returns user's debts stats object (requires JWT in "Bearer" header)
// @Tags debt
// @Produce json
// @Success 200 {object} debt_service.DebtsMetrics
// @Failure 500 {object} errorz.Error_
// @Failure 503 {object} errorz.Error_
// @Failure 401 {object} errorz.Error_
// @Security BearerAuth
// @Param Authorization header string true "Bearer JWT token" // Добавляем параметр для JWT
// @Router /debts/stats [get]
func (c *DebtController) GetDebtsMetrics(ctx *gin.Context) {
	token, err := utils.ExtractTokenFromHeaders(ctx)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	payload, err := utils.ExtractPayloadFromJWT(*token)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	username := payload["username"].(string)

	metrics, err := c.ds.CountProfileMetrics(ctx, username)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	ctx.JSON(200, &metrics)
}
