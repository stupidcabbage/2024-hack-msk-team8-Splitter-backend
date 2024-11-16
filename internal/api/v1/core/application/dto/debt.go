package dto

type DebtDto struct {
	Id           int64  `json:"id" db:"id"`
	Name         string `json:"name" db:"name" binding:"required,max=32,min=6"`
	FromUsername string `json:"from_username" db:"from_username" binding:"required,max=32,min=6"`
	ToUsername   string `json:"to_username" db:"to_username" binding:"required,max=32,min=6"`
	// sets automatically to false
	IsClosed    bool `json:"is_closed" db:"is_closed"`
	TotalAmount int  `json:"total_amount" db:"total_amount" binding:"required"`
	// sets automatically to 0
	PayedAmount int    `json:"payed_amount" db:"payed_amount" binding:"required"`
	CreatedAt   string `json:"created_at" db:"created_at"`
	UpdatedAt   string `json:"updated_at" db:"updated_at"`
	GroupId     *int64 `json:"group_id" db:"group_id"`
}

type DebtDtoWOId struct {
	Name         string `json:"name" db:"name" binding:"required,max=32,min=6"`
	FromUsername string `json:"from_username" db:"from_username" binding:"required,max=32,min=6"`
	ToUsername   string `json:"to_username" db:"to_username" binding:"required,max=32,min=6"`
	// sets automatically to false
	IsClosed    bool `json:"is_closed" db:"is_closed"`
	TotalAmount int  `json:"total_amount" db:"total_amount" binding:"required"`
	// sets automatically to 0
	PayedAmount int    `json:"payed_amount" db:"payed_amount" binding:"required"`
	CreatedAt   string `json:"created_at" db:"created_at"`
	UpdatedAt   string `json:"updated_at" db:"updated_at"`
	GroupId     *int64 `json:"group_id" db:"group_id"`
}

type CreateDebtDto struct {
	Name string `json:"name" db:"name" binding:"required,max=32,min=1"`
	// FromUsername string `json:"from_username" db:"from_username"`
	ToUsername  string `json:"to_username" db:"to_username"`
	TotalAmount int    `json:"total_amount" db:"total_amount" binding:"required"`
	InviteCode  string `json:"invite_code" binding:"required,max=8,min=8"`
	GroupId     *int64 `json:"group_id" db:"group_id"`
}

type IncreaseDebtPayedAmountDto struct {
	Amount int `json:"amount" binding:"required,min=0"`
}

type GetDebtDto struct {
	Id           int64  `json:"id" db:"id"`
	Name         string `json:"name" db:"name" binding:"required,max=32,min=6"`
	FromUsername string `json:"from_user_email" db:"from_user_email" binding:"required,max=32,min=6"`
	ToUsername   string `json:"to_user_email" db:"to_user_email" binding:"required,max=32,min=6"`
	// sets automatically to false
	IsClosed    bool `json:"is_closed" db:"is_closed"`
	TotalAmount int  `json:"total_amount" db:"total_amount" binding:"required"`
	// sets automatically to 0
	PayedAmount int    `json:"payed_amount" db:"payed_amount" binding:"required"`
	CreatedAt   string `json:"created_at" db:"created_at"`
	UpdatedAt   string `json:"updated_at" db:"updated_at"`
}

type UpdateDebtDto struct {
	IsClosed    bool   `json:"is_closed" db:"is_closed"`
	PayedAmount int    `json:"payed_amount" db:"payed_amount"`
	UpdatedAt   string `json:"updated_at" db:"updated_at"`
}
