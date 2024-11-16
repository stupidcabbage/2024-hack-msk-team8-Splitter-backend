package dto

type AuthorizeUserDto struct {
	Username string `json:"username" db:"username" binding:"required,max=32,min=6"`
	Password string `json:"password" db:"password" binding:"required,max=64,min=6"`
}

type ChangeUserPasswordDto struct {
	OldPassword string `json:"old_password" db:"password" binding:"required,max=64,min=6"`
	NewPassword string `json:"new_password" db:"password" binding:"required,max=64,min=6"`
}
