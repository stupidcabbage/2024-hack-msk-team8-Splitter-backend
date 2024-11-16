package dto

type CreateUserDto struct {
	Username    string `json:"username" db:"username" binding:"required,max=32,min=6"`
	Password    string `json:"password" db:"password" binding:"required,max=64,min=6"`
	InviteCode  string `json:"invite_code" db:"invite_code"`
	PhoneNumber string `json:"phone_number" db:"phone_number" binding:"required,max=32,min=6"`
}

type UserDto struct {
	Username    string `json:"username" db:"username"`
	Password    string `json:"password" db:"password"`
	InviteCode  string `json:"invite_code" db:"invite_code"`
	CreatedAt   string `json:"created_at" db:"created_at"`
	UpdatedAt   string `json:"updated_at" db:"updated_at"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`
}

type GetUserDto struct {
	Username    string `json:"username" db:"username"`
	InviteCode  string `json:"invite_code" db:"invite_code"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`
	CreatedAt   string `json:"created_at" db:"created_at"`
	UpdatedAt   string `json:"updated_at" db:"updated_at"`
}

type UpdateUserDto struct {
	Password   string `json:"password" db:"password"`
	InviteCode string `json:"invite_code" db:"invite_code"`
	// sets in service automatically
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}
