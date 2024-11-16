package debt_service

import (
	"context"
	"time"

	"example.com/m/internal/api/v1/adapters/repositories"
	"example.com/m/internal/api/v1/core/application/dto"
	"example.com/m/internal/api/v1/core/application/errorz"
	"example.com/m/internal/api/v1/core/application/services/group_service"
	"example.com/m/internal/api/v1/core/application/services/user_service"
	"example.com/m/internal/api/v1/infrastructure/prom"
)

type DebtService struct {
	dr repositories.DebtRepository
	us user_service.UserService
	gs group_service.GroupService
}

func NewDebtService(dr *repositories.DebtRepository, us *user_service.UserService, gs *group_service.GroupService) *DebtService {
	return &DebtService{
		dr: *dr,
		us: *us,
		gs: *gs,
	}
}

// TODO проверка from_username и to_username на существование
func (s *DebtService) CreateDebt(ctx context.Context, d *dto.CreateDebtDto) (*dto.DebtDto, *errorz.Error_) {
	debter, exception := s.us.GetUserByInviteCode(ctx, d.InviteCode)
	if exception != nil {
		return nil, exception
	}

	if debter.Username == d.ToUsername {
		return nil, &errorz.ErrCantDebtYourself
	}

	if d.GroupId != nil {
		isValidGroup, exception := s.gs.IsGroupExistsAndUserInGroup(*d.GroupId, debter.Username)
		if exception != nil {
			return nil, exception
		}
		if !isValidGroup {
			return nil, &errorz.ErrCantDebtNotInGroup
		}
		prom.DebtCreatedCount.WithLabelValues("method").Inc()
	}

	var debtToCreate dto.DebtDtoWOId = dto.DebtDtoWOId{
		Name:         d.Name,
		FromUsername: debter.Username,
		ToUsername:   d.ToUsername,
		TotalAmount:  d.TotalAmount,
		// ставлю нулевые значения для оплаченной части долга
		IsClosed:    false,
		PayedAmount: 0,
		CreatedAt:   time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		GroupId:     d.GroupId,
	}

	id, err := s.dr.Create(&debtToCreate)
	if err != nil {
		return nil, &errorz.ErrDatabaseError
	}

	debt, err := s.dr.GetById(id)
	if err != nil {
		return nil, &errorz.ErrDatabaseError
	}
	if debt == nil {
		return nil, &errorz.ErrDebtNotFound
	}

	return debt, nil
}

// стейтмент и для создателя/получателя долга и для должника
func canUserAccessDebt(debt *dto.DebtDto, username *string) bool {
	return debt.FromUsername == *username || debt.ToUsername == *username
}

// стейтмент для создателя/получателя долга
func canUserAffectDebt(debt *dto.DebtDto, username *string) bool {
	return debt.ToUsername == *username
}

func (s *DebtService) GetDebtById(ctx context.Context, username string, debtId int64) (*dto.DebtDto, *errorz.Error_) {
	debt, err := s.dr.GetById(&debtId)
	if err != nil {
		return nil, &errorz.ErrDatabaseError
	}
	if debt == nil {
		return nil, &errorz.ErrDebtNotFound
	}
	if !canUserAccessDebt(debt, &username) {
		return nil, &errorz.ErrDebtNotFound
	}

	return debt, nil
}

func (s *DebtService) CloseDebt(ctx context.Context, username string, debtId int64) *errorz.Error_ {
	debt, err := s.dr.GetById(&debtId)
	if err != nil {
		return &errorz.ErrDatabaseError
	}
	if debt == nil {
		return &errorz.ErrDebtNotFound
	}
	if !canUserAffectDebt(debt, &username) {
		return &errorz.ErrDebtAffectingPermissionDenied
	}

	updateData := dto.UpdateDebtDto{IsClosed: true, PayedAmount: debt.PayedAmount, UpdatedAt: time.Now().UTC().Format("2006-01-02T15:04:05Z")}
	err = s.dr.UpdateById(&debtId, &updateData)
	if err != nil {
		return &errorz.ErrDatabaseError
	}

	return nil
}

type DebtsMetrics struct {
	PayedIncomingDebtsAmount    int `json:"payed_incoming_debts_amount"`
	UnpayedIncomingDebtsAmount  int `json:"unpayed_incoming_debts_amount"`
	PayedOutcomingDebtsAmount   int `json:"payed_outcoming_debts_amount"`
	UnpayedOutcomingDebtsAmount int `json:"unpayed_outcoming_debts_amount"`
}

func (s *DebtService) CountProfileMetrics(ctx context.Context, username string) (*DebtsMetrics, *errorz.Error_) {
	// уверенным в себе разработчикам обработка ошибок не нужна
	payedIncomingDebtsAmount, _ := s.dr.CountPayedIncomingDebtsAmount(&username)
	unpayedIncomingDebtsAmount, _ := s.dr.CountUnpayedIncomingDebtsAmount(&username)
	payedOutcomingDebtsAmount, _ := s.dr.CountPayedOutcomingDebtsAmount(&username)
	unpayedOutcomingDebtsAmount, _ := s.dr.CountUnpayedOutcomingDebtsAmount(&username)

	return &DebtsMetrics{
		PayedIncomingDebtsAmount:    payedIncomingDebtsAmount,
		UnpayedIncomingDebtsAmount:  unpayedIncomingDebtsAmount,
		PayedOutcomingDebtsAmount:   payedOutcomingDebtsAmount,
		UnpayedOutcomingDebtsAmount: unpayedOutcomingDebtsAmount,
	}, nil
}

// возвращает значение payed_amount
func calculateNewPayedAmount(debt dto.DebtDto, increase int) int {
	if (debt.PayedAmount + increase) >= debt.TotalAmount {
		return debt.TotalAmount
	}
	return debt.PayedAmount + increase
}

func (s *DebtService) IncreaseDebtPayedAmount(ctx context.Context, username string, debtId int64, increase int) (bool, *errorz.Error_) {
	debt, err := s.dr.GetById(&debtId)
	if err != nil {
		return false, &errorz.ErrDatabaseError
	}
	if debt == nil {
		return false, &errorz.ErrDebtNotFound
	}

	if !canUserAffectDebt(debt, &username) {
		return false, &errorz.ErrDebtAffectingPermissionDenied
	}
	if debt.IsClosed {
		return false, &errorz.ErrDebtIsClosed
	}
	payedAmount := calculateNewPayedAmount(*debt, increase)

	updateData := dto.UpdateDebtDto{PayedAmount: payedAmount, IsClosed: false, UpdatedAt: time.Now().UTC().Format("2006-01-02T15:04:05Z")}
	if payedAmount == debt.TotalAmount {
		updateData.IsClosed = true
	}

	err = s.dr.UpdateById(&debtId, &updateData)
	if err != nil {
		return false, &errorz.ErrDatabaseError
	}

	return updateData.IsClosed, nil

}

func (s *DebtService) GetIncomingDebts(ctx context.Context, username string) (*[]dto.DebtDto, *errorz.Error_) {
	debts, err := s.dr.GetIncoming(&username)
	if err != nil {
		return nil, &errorz.ErrDatabaseError
	}

	return debts, nil
}
func (s *DebtService) GetOutcomingDebts(ctx context.Context, username string) (*[]dto.DebtDto, *errorz.Error_) {
	debts, err := s.dr.GetOutcoming(&username)
	if err != nil {
		return nil, &errorz.ErrDatabaseError
	}

	return debts, nil
}

func (s *DebtService) GetClosedDebts(ctx context.Context, username string) (*[]dto.DebtDto, *errorz.Error_) {
	debts, err := s.dr.GetClosed(&username)
	if err != nil {
		return nil, &errorz.ErrDatabaseError
	}

	return debts, nil
}
