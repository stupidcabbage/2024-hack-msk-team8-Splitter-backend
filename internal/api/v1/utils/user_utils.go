package utils

import (
	"math/rand"
	"time"

	"example.com/m/internal/api/v1/core/application/dto"
	"golang.org/x/crypto/bcrypt"
)

func UpdateUserTimestamps(u *dto.UpdateUserDto) {
	u.UpdatedAt = time.Now().UTC().Format("2006-01-02T15:04:05Z")
}

func GenerateInviteCode() string {
	var letterRunes = []rune("0123456789")

	b := make([]rune, 8)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func ExcludeUserCredentials(u *dto.UserDto) dto.GetUserDto {
	return dto.GetUserDto{
		Username:    u.Username,
		InviteCode:  u.InviteCode,
		PhoneNumber: u.PhoneNumber,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 7)
	return string(bytes), err
}
