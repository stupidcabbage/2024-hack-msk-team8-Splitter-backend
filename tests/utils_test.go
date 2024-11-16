package utils

import (
	"testing"

	"example.com/m/internal/api/v1/core/application/dto"
	"example.com/m/internal/api/v1/utils"
)


func TestGroupInviteCode(t *testing.T) {
	code1 := utils.GenerateGroupInviteCode()
	code2 := utils.GenerateGroupInviteCode()
	if code1 == code2 {
		t.Fatalf("code should be different")
	}
}


func TestInviteCode(t *testing.T) {
	code1 := utils.GenerateInviteCode()
	code2 := utils.GenerateInviteCode()
	if code1 == code2 {
		t.Fatalf("code should be different")
	}
}

func TestHashPassword(t *testing.T) {
	hash, err := utils.HashPassword("123456")
	if err != nil {
		t.Fatal(err)
	}
	if hash == "" {
		t.Fatal("hash should not be empty")
	}
}


func TestExcludeCredentials(t *testing.T) {
	user := dto.UserDto{
		Username:    "username",
		Password:    "password",
		PhoneNumber: "1234567890",
		InviteCode:  "inviteCode",
		CreatedAt:   "createdAt",
		UpdatedAt:   "updatedAt",
	}
	credentials := utils.ExcludeUserCredentials(&user)
	if credentials.Username != user.Username {
		t.Fatalf("username should be equal")		
	}
	if credentials.PhoneNumber != user.PhoneNumber {
		t.Fatalf("phoneNumber should be equal")		
	}
	if credentials.InviteCode != user.InviteCode {
		t.Fatalf("inviteCode should be equal")		
	}
	if credentials.CreatedAt != user.CreatedAt {
		t.Fatalf("createdAt should be equal")		
	}
	if credentials.UpdatedAt != user.UpdatedAt {
		t.Fatalf("updatedAt should be equal")		
	}
}

func TestUpdateUserTimestamps(t *testing.T) {
	u := &dto.UpdateUserDto{}
	u.UpdatedAt = "updatedAt"
	utils.UpdateUserTimestamps(u)

	if u.UpdatedAt == "" {
		t.Fatalf("updatedAt should not be empty")
	}
	if u.UpdatedAt == "updatedAt" {
		t.Fatalf("updatedAt should not be equal")
	}
}