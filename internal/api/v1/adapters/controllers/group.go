package controllers

import (
	"strconv"

	"example.com/m/internal/api/v1/core/application/dto"
	"example.com/m/internal/api/v1/core/application/errorz"
	"example.com/m/internal/api/v1/core/application/services/group_service"
	"example.com/m/internal/api/v1/utils"
	"github.com/gin-gonic/gin"
)

type GroupController struct {
	gs group_service.GroupService
}

func NewGroupController(s *group_service.GroupService) *GroupController {
	return &GroupController{
		gs: *s,
	}
}

// @BasePath /api/v1

// Create group
// @Schemes
// @Description Creates new group (requires JWT in "Bearer" header)
// @Tags group
// @Accept json
// @Produce json
// @Param group body dto.CreateGroupDto true "Group data"
// @Success 201 {object} dto.GroupDto
// @Failure 500 {object} errorz.Error_
// @Failure 503 {object} errorz.Error_
// @Failure 400 {object} errorz.Error_
// @Failure 401 {object} errorz.Error_
// @Security BearerAuth
// @Param Authorization header string true "Bearer JWT token" // Добавляем параметр для JWT
// @Router /groups [post]
func (c *GroupController) CreateGroup(ctx *gin.Context) {
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

	var group dto.CreateGroupDto
	if err := ctx.ShouldBindBodyWithJSON(&group); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	group.CreatedBy = username

	createdGroup, exception := c.gs.CreateGroup(group)
	if exception != nil {
		ctx.JSON(int(exception.StatusCode), exception)
		return
	}

	ctx.JSON(201, createdGroup)
}

// Add to group
// @Schemes
// @Description Connect user to group (requires JWT in "Bearer" header)
// @Tags group
// @Accept json
// @Produce json
// @Param invite_code path string true "Invite Code"
// @Success 201
// @Failure 500 {object} errorz.Error_
// @Failure 503 {object} errorz.Error_
// @Failure 400 {object} errorz.Error_
// @Failure 401 {object} errorz.Error_
// @Security BearerAuth
// @Param Authorization header string true "Bearer JWT token" // Добавляем параметр для JWT
// @Router /groups/{invite_code}/join [post]
func (c *GroupController) AddToGroup(ctx *gin.Context) {
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
	inviteCode := ctx.Param("invite_code")
	err = c.gs.AddToGroup(
		&dto.GroupDtoComeIn{
			InviteCode: inviteCode,
			Username:   username,
		})
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}
	ctx.JSON(201, "ok")
}

// Gets group users
// @Schemes
// @Description Get all group users (requires JWT in "Bearer" header)
// @Tags group
// @Accept json
// @Produce json
// @Param group_id path string true "Group id"
// @Success 200 {array} []dto.GetUserDto
// @Failure 500 {object} errorz.Error_
// @Failure 503 {object} errorz.Error_
// @Failure 401 {object} errorz.Error_
// @Security BearerAuth
// @Param Authorization header string true "Bearer JWT token" // Добавляем параметр для JWT
// @Router /groups/{group_id}/users [GET]
func (c *GroupController) GetUsersByGroup(ctx *gin.Context) {
	var groupId int64
	groupId, e := strconv.ParseInt(ctx.Param("group_id"), 10, 64)
	if e != nil {
		ctx.JSON(int(errorz.InvalidGroupIdError.StatusCode), errorz.InvalidGroupIdError)
		return
	}
	users, err := c.gs.GetUsersByGroup(groupId)

	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}
	var usersToReturn []dto.GetUserDto

	for _, user := range *users {
		usersToReturn = append(usersToReturn, utils.ExcludeUserCredentials(&user))
	}

	if len(usersToReturn) == 0 {
		ctx.JSON(200, make([]string, 0))
		return
	}

	ctx.JSON(200, &usersToReturn)
}

// Получить список своих долгов
// @Schemes
// @Description Получить весь список своих долгов в группе (requires JWT in "Bearer" header)
// @Tags group
// @Accept json
// @Produce json
// @Param group_id path string true "Group id"
// @Success 200 {array} []dto.DebtDto
// @Failure 500 {object} errorz.Error_
// @Failure 503 {object} errorz.Error_
// @Failure 401 {object} errorz.Error_
// @Security BearerAuth
// @Param Authorization header string true "Bearer JWT token" // Добавляем параметр для JWT
// @Router /groups/{group_id}/debts/incoming [GET]
func (c *GroupController) GetIncomingDebtsByGroup(ctx *gin.Context) {
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

	var groupId int64
	groupId, e := strconv.ParseInt(ctx.Param("group_id"), 10, 64)
	if e != nil {
		ctx.JSON(int(errorz.InvalidGroupIdError.StatusCode), errorz.InvalidGroupIdError)
		return
	}
	debts, err := c.gs.GetIncomingDebtsByGroup(username, groupId)

	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	if len(*debts) == 0 {
		ctx.JSON(200, make([]string, 0))
		return
	}
	ctx.JSON(200, debts)
}

// Получить список должников
// @Schemes
// @Description Получить весь список должников в группе (requires JWT in "Bearer" header)
// @Tags group
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer JWT token" // Добавляем параметр для JWT
// @Param group_id path string true "Group id"
// @Success 200 {array} []dto.DebtDto
// @Failure 500 {object} errorz.Error_
// @Failure 503 {object} errorz.Error_
// @Failure 401 {object} errorz.Error_
// @Failure 400 {object} errorz.Error_
// @Router /groups/{group_id}/debts/outcoming [GET]
func (c *GroupController) GetOutcomingDebtsByGroup(ctx *gin.Context) {
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

	var groupId int64
	groupId, e := strconv.ParseInt(ctx.Param("group_id"), 10, 64)
	if e != nil {
		ctx.JSON(int(errorz.InvalidGroupIdError.StatusCode), errorz.InvalidGroupIdError)
		return
	}
	debts, err := c.gs.GetOutcomingDebtsByGroup(username, groupId)

	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}

	if len(*debts) == 0 {
		ctx.JSON(200, make([]string, 0))
		return
	}
	ctx.JSON(200, debts)
}

// Получения списка созданных групп
// @Schemes
// @Description Получить список созданных групп (requires JWT in "Bearer" header)
// @Tags group
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} []dto.GroupDto
// @Failure 500 {object} errorz.Error_
// @Failure 503 {object} errorz.Error_
// @Failure 401 {object} errorz.Error_
// @Router /groups/created [GET]
func (c *GroupController) GetCreatedGroups(ctx *gin.Context) {
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

	groups, err := c.gs.GetCreatedGroups(username)
	if err != nil {
		ctx.JSON(int(err.StatusCode), err)
		return
	}
	if len(*groups) == 0 {
		ctx.JSON(200, make([]string, 0))
		return
	}

	ctx.JSON(200, groups)
}
