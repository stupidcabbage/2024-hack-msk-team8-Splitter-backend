package repositories

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"example.com/m/internal/api/v1/core/application/dto"
	"github.com/doug-martin/goqu/v9"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(u *dto.UserDto) error {
	query, _, _ := goqu.Insert("users").Rows(*u).ToSQL()
	_, err := r.db.Exec(query)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetByUsername(username *string) (*dto.UserDto, error) {
	query, _, _ := goqu.From("users").Where(goqu.Ex{
		"username": *username,
	}).ToSQL()

	var user dto.UserDto
	err := r.db.QueryRow(query).Scan(&user.Username, &user.Password, &user.InviteCode, &user.PhoneNumber, &user.CreatedAt, &user.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByPhoneNumber(phoneNumber *string) (*dto.UserDto, error) {
	query, _, _ := goqu.From("users").Where(goqu.Ex{
		"phone_number": *phoneNumber,
	}).ToSQL()

	var user dto.UserDto
	err := r.db.QueryRow(query).Scan(&user.Username, &user.Password, &user.InviteCode, &user.PhoneNumber, &user.CreatedAt, &user.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByInviteCode(code *string) (*dto.UserDto, error) {
	query, _, _ := goqu.From("users").Where(goqu.Ex{
		"invite_code": *code,
	}).ToSQL()

	var user dto.UserDto
	err := r.db.QueryRow(query).Scan(&user.Username, &user.Password, &user.InviteCode, &user.PhoneNumber, &user.CreatedAt, &user.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateByUsername(username *string, u *dto.UpdateUserDto) error {
	var uMap map[string]interface{}
	inrec, _ := json.Marshal(*u)
	json.Unmarshal(inrec, &uMap)

	var rec goqu.Record = uMap

	query, _, _ := goqu.From("users").Where(goqu.C("username").Eq(*username)).Update().Set(
		rec,
	).ToSQL()

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetUsersByGroup(groupId *int64) (*[]dto.UserDto, error) {
	query, _, _ := goqu.From("users").
		Join(goqu.T("groups_users"), goqu.On(goqu.Ex{"users.username": goqu.I("groups_users.username")})).
		Where(goqu.Ex{"groups_users.group_id": groupId}).ToSQL() 
	log.Println(query)

	var users []dto.UserDto
	rows, err := r.db.Query(query)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user dto.UserDto
		var groupId string

		err = rows.Scan(&user.Username, &user.Password, &user.InviteCode, &user.PhoneNumber, &user.CreatedAt, &user.UpdatedAt, &user.Username, &groupId)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		fmt.Println(user)
		users = append(users, user)
	}
	fmt.Println(&users)
	return &users, nil
}
