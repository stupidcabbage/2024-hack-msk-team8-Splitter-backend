package user_service

import (
	"context"
	"time"

	"example.com/m/internal/api/v1/adapters/repositories"
	"example.com/m/internal/api/v1/core/application/dto"
	"example.com/m/internal/api/v1/core/application/errorz"
	"example.com/m/internal/api/v1/infrastructure/prom"
	"example.com/m/internal/api/v1/utils"
)

type UserService struct {
	ur repositories.UserRepository
}

func NewUserService(ur *repositories.UserRepository) *UserService {
	return &UserService{ur: *ur}
}

func (s *UserService) isUserExist(username string, phoneNumber string) (*bool, *errorz.Error_) {
	foundUserByEmail, err := s.ur.GetByUsername(&username)
	if err != nil {
		return nil, &errorz.ErrDatabaseError
	}
	foundUserByPhoneNumber, err := s.ur.GetByPhoneNumber(&phoneNumber)
	if err != nil {
		return nil, &errorz.ErrDatabaseError
	}

	state := foundUserByEmail != nil || foundUserByPhoneNumber != nil

	return &state, nil
}

func (s *UserService) CreateUser(ctx context.Context, u dto.CreateUserDto) (*dto.UserDto, *errorz.Error_) {
	userExists, exception := s.isUserExist(u.Username, u.PhoneNumber)
	if exception != nil {
		return nil, exception
	}
	if *userExists {
		return nil, &errorz.ErrUserAlreadyExists
	}

	hashedPassword, _ := utils.HashPassword(u.Password)
	userToCreate := dto.UserDto{
		Username:    u.Username,
		Password:    hashedPassword,
		PhoneNumber: u.PhoneNumber,
		InviteCode:  utils.GenerateInviteCode(),
		CreatedAt:   time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   time.Now().UTC().Format("2006-01-02T15:04:05Z"),
	}

	err := s.ur.Create(&userToCreate)
	if err != nil {
		return nil, &errorz.ErrDatabaseError
	}

	prom.UserCreatedCounter.WithLabelValues("method").Inc()
	return &userToCreate, nil
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*dto.UserDto, *errorz.Error_) {
	user, err := s.ur.GetByUsername(&username)
	if err != nil {
		return nil, &errorz.ErrDatabaseError
	}

	if user == nil {
		return nil, &errorz.ErrUserNotFound
	}

	return user, nil
}

func (s *UserService) GetUserByInviteCode(ctx context.Context, code string) (*dto.UserDto, *errorz.Error_) {
	user, err := s.ur.GetByInviteCode(&code)
	if err != nil {
		return nil, &errorz.ErrDatabaseError
	}

	if user == nil {
		return nil, &errorz.ErrUserNotFound
	}

	return user, nil
}

func (s *UserService) UpdateUserByUsername(ctx context.Context, username string, u dto.UpdateUserDto) (*dto.UserDto, *errorz.Error_) {
	_, exception := s.GetUserByUsername(ctx, username)
	if exception != nil {
		return nil, exception
	}

	utils.UpdateUserTimestamps(&u)

	if err := s.ur.UpdateByUsername(&username, &u); err != nil {
		return nil, &errorz.ErrDatabaseError
	}

	updatedUser, exception := s.GetUserByUsername(ctx, username)
	if exception != nil {
		return nil, exception
	}

	return updatedUser, nil
}

func (s *UserService) RegenerateInviteCode(ctx context.Context, username string) *errorz.Error_ {
	user, exception := s.GetUserByUsername(ctx, username)
	if exception != nil {
		return exception
	}

	updateData := dto.UpdateUserDto{InviteCode: utils.GenerateInviteCode(), Password: user.Password}
	if err := s.ur.UpdateByUsername(&username, &updateData); err != nil {
		return &errorz.ErrDatabaseError
	}

	return nil
}
