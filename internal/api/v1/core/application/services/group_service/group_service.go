package group_service

import (
	"example.com/m/internal/api/v1/adapters/repositories"
	"example.com/m/internal/api/v1/core/application/dto"
	"example.com/m/internal/api/v1/core/application/errorz"
	"example.com/m/internal/api/v1/infrastructure/prom"
	"example.com/m/internal/api/v1/utils"
)

type GroupService struct {
	gr repositories.GroupRepository
	ur repositories.UserRepository
	dr repositories.DebtRepository
}

func NewGroupService(gr *repositories.GroupRepository, ur *repositories.UserRepository, dr *repositories.DebtRepository) *GroupService {
	return &GroupService{gr: *gr, ur: *ur, dr: *dr}
}

func (s *GroupService) CreateGroup(group dto.CreateGroupDto) (*dto.GroupDto, *errorz.Error_) {

	group.InviteCode = utils.GenerateGroupInviteCode()

	isInviteCodeExists, err := s.gr.IfInviteCodeExists(group.InviteCode)
	if err != nil {
		return nil, &errorz.ErrDatabaseError
	}

	if isInviteCodeExists {
		group.InviteCode = utils.GenerateGroupInviteCode()
	}

	id, err := s.gr.Create(
		&group,
	)

	if err != nil {
		return nil, &errorz.ErrDatabaseError
	}

	prom.GroupCreatedCount.WithLabelValues("method").Inc()
	return &dto.GroupDto{
		Id:         id,
		InviteCode: group.InviteCode,
		CreatedBy:  group.CreatedBy,
		Name:       group.Name,
	}, nil
}

func (s *GroupService) AddToGroup(d *dto.GroupDtoComeIn) *errorz.Error_ {
	err := s.gr.JoinToGroup(d)
	if err != nil {
		return &errorz.ErrDatabaseError
	}

	return nil
}

func (s *GroupService) GetUsersByGroup(groupId int64) (*[]dto.UserDto, *errorz.Error_) {
	users, err := s.ur.GetUsersByGroup(&groupId)
	if err != nil {
		return nil, &errorz.ErrDatabaseError
	}
	return users, nil
}

func (s *GroupService) GetCreatedGroups(username string) (*[]dto.GroupDto, *errorz.Error_) {
	groups, err := s.gr.GetCreatedGroups(username)
	if err != nil {
		return nil, &errorz.ErrDatabaseError
	}

	return groups, nil
}

func (s *GroupService) GetIncomingDebtsByGroup(username string, groupId int64) (*[]dto.DebtDto, *errorz.Error_) {
	debts, err := s.dr.GetIncomingDebtsByGroup(&username, &groupId)
	if err != nil {
		return nil, &errorz.ErrDatabaseError
	}
	return debts, nil
}

func (s *GroupService) GetOutcomingDebtsByGroup(username string, groupId int64) (*[]dto.DebtDto, *errorz.Error_) {
	debts, err := s.dr.GetOutcomingDebtsByGroup(&username, &groupId)
	if err != nil {
		return nil, &errorz.ErrDatabaseError
	}
	return debts, nil
}

func (s *GroupService) IsGroupExistsAndUserInGroup(groupId int64, username string) (bool, *errorz.Error_) {
	isExists, err := s.gr.IsGroupExistsAndUserInGroup(groupId, username)
	if err != nil {
		return false, &errorz.ErrDatabaseError
	}
	return isExists, nil
}
