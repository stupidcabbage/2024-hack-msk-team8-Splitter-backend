package dto

type GroupDto struct {
	Id         int64  `json:"id" db:"id"`
	Name       string `json:"name" db:"name" binding:"required,min=1,max=32"`
	InviteCode string `json:"invite_code" db:"invite_code"`
	CreatedBy  string `json:"created_by" db:"created_by"`
}

type GroupDtoIn struct {
	CreatedBy string `json:"created_by"`
}

type CreateGroupDto struct {
	Name       string `json:"name" db:"name" binding:"required,min=1,max=32"`
	InviteCode string `json:"invite_code" db:"invite_code"`
	CreatedBy  string `json:"created_by"`
}

type GroupDtoComeIn struct {
	InviteCode string `json:"invite_code"`
	Username   string `json:"username"`
}
