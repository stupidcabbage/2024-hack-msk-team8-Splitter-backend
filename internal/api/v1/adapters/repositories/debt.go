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

type DebtRepository struct {
	db *sql.DB
}

func NewDebtRepository(db *sql.DB) *DebtRepository {
	return &DebtRepository{
		db: db,
	}
}

func (r *DebtRepository) Create(d *dto.DebtDtoWOId) (*int64, error) {
	var id int64
	query, _, _ := goqu.Insert("debts").Rows(*d).Returning("id").ToSQL()
	err := r.db.QueryRow(query).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (r *DebtRepository) GetById(id *int64) (*dto.DebtDto, error) {
	query, _, _ := goqu.From("debts").Where(goqu.Ex{
		"id": *id,
	}).ToSQL()
	log.Println(query)

	var debt dto.DebtDto

	err := r.db.QueryRow(query).Scan(&debt.Id, &debt.Name, &debt.FromUsername, &debt.ToUsername, &debt.IsClosed, &debt.TotalAmount, &debt.PayedAmount, &debt.CreatedAt, &debt.UpdatedAt, &debt.GroupId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &debt, nil
}

func (r *DebtRepository) CountUnpayedIncomingDebtsAmount(username *string) (int, error) {
	query := fmt.Sprintf("select SUM(total_amount) from debts where to_username = '%s';", *username)
	var amount int
	err := r.db.QueryRow(query).Scan(&amount)
	if errors.Is(sql.ErrNoRows, err) {
		return 0, nil
	}
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return amount, nil
}

func (r *DebtRepository) CountUnpayedOutcomingDebtsAmount(username *string) (int, error) {
	query := fmt.Sprintf("select SUM(total_amount) from debts where from_username = '%s';", *username)
	var amount int
	err := r.db.QueryRow(query).Scan(&amount)
	if errors.Is(sql.ErrNoRows, err) {
		return 0, nil
	}
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return amount, nil
}

func (r *DebtRepository) CountPayedIncomingDebtsAmount(username *string) (int, error) {
	query := fmt.Sprintf("select SUM(payed_amount) from debts where to_username = '%s';", *username)
	var amount int
	err := r.db.QueryRow(query).Scan(&amount)
	if errors.Is(sql.ErrNoRows, err) {
		return 0, nil
	}
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return amount, nil
}

func (r *DebtRepository) CountPayedOutcomingDebtsAmount(username *string) (int, error) {
	query := fmt.Sprintf("select SUM(payed_amount) from debts where from_username = '%s';", *username)
	var amount int
	err := r.db.QueryRow(query).Scan(&amount)
	if errors.Is(sql.ErrNoRows, err) {
		return 0, nil
	}
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return amount, nil
}

func (r *DebtRepository) GetOutcoming(username *string) (*[]dto.DebtDto, error) {
	var debts []dto.DebtDto
	query, _, _ := goqu.From("debts").Where(goqu.Ex{
		"from_username": *username,
	}).ToSQL()

	rows, err := r.db.Query(query)
	if errors.Is(err, sql.ErrNoRows) {
		return &debts, nil
	}
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for rows.Next() {
		var debt dto.DebtDto

		rows.Scan(&debt.Id, &debt.Name, &debt.FromUsername, &debt.ToUsername, &debt.IsClosed, &debt.TotalAmount, &debt.PayedAmount, &debt.CreatedAt, &debt.UpdatedAt, &debt.GroupId)
		fmt.Println(debt)
		debts = append(debts, debt)
	}

	return &debts, nil
}

func (r *DebtRepository) GetIncoming(username *string) (*[]dto.DebtDto, error) {
	var debts []dto.DebtDto
	query, _, _ := goqu.From("debts").Where(goqu.Ex{
		"to_username": *username,
	}).ToSQL()

	rows, err := r.db.Query(query)
	if errors.Is(err, sql.ErrNoRows) {
		return &debts, nil
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var debt dto.DebtDto

		rows.Scan(&debt.Id, &debt.Name, &debt.FromUsername, &debt.ToUsername, &debt.IsClosed, &debt.TotalAmount, &debt.PayedAmount, &debt.CreatedAt, &debt.UpdatedAt, &debt.GroupId)
		debts = append(debts, debt)
	}

	return &debts, nil
}

func (r *DebtRepository) GetClosed(username *string) (*[]dto.DebtDto, error) {
	var debts []dto.DebtDto
	query, _, _ := goqu.From("debts").Where(goqu.ExOr{
		"to_username":   *username,
		"from_username": *username,
	}, goqu.Ex{
		"is_closed": true,
	}).ToSQL()
	rows, err := r.db.Query(query)
	if errors.Is(err, sql.ErrNoRows) {
		return &debts, nil
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var debt dto.DebtDto

		rows.Scan(&debt.Id, &debt.Name, &debt.FromUsername, &debt.ToUsername, &debt.IsClosed, &debt.TotalAmount, &debt.PayedAmount, &debt.CreatedAt, &debt.UpdatedAt, &debt.GroupId)
		debts = append(debts, debt)
	}

	return &debts, nil
}

func (r *DebtRepository) UpdateById(id *int64, d *dto.UpdateDebtDto) error {
	var uMap map[string]interface{}
	inrec, _ := json.Marshal(*d)
	json.Unmarshal(inrec, &uMap)

	var rec goqu.Record = uMap

	query, _, _ := goqu.From("debts").Where(goqu.C("id").Eq(*id)).Update().Set(
		rec,
	).ToSQL()

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *DebtRepository) GetIncomingDebtsByGroup(username *string, groupId *int64) (*[]dto.DebtDto, error) {
	var debts []dto.DebtDto
	query, _, _ := goqu.From("debts").
		Select("debts.*").
		Join(goqu.T("users"), goqu.On(goqu.Ex{"debts.to_username": goqu.I("users.username")})).
		Join(goqu.T("groups_users"), goqu.On(goqu.Ex{"users.username": goqu.I("users.username")})).
		Where(goqu.Ex{
			"users.username": username,
			"debts.group_id": groupId,
		}).
		ToSQL()
	rows, err := r.db.Query(query)
	if errors.Is(err, sql.ErrNoRows) {
		return &debts, nil
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var debt dto.DebtDto

		rows.Scan(&debt.Id, &debt.Name, &debt.FromUsername, &debt.ToUsername, &debt.IsClosed, &debt.TotalAmount, &debt.PayedAmount, &debt.CreatedAt, &debt.UpdatedAt, &debt.GroupId)
		debts = append(debts, debt)
	}

	return &debts, nil
}

func (r *DebtRepository) GetOutcomingDebtsByGroup(username *string, groupId *int64) (*[]dto.DebtDto, error) {
	var debts []dto.DebtDto
	query, _, _ := goqu.From("debts").
		Select("debts.*").
		InnerJoin(goqu.T("users"), goqu.On(goqu.Ex{"debts.from_username": goqu.I("users.username")})).
		InnerJoin(goqu.T("groups_users"), goqu.On(goqu.Ex{"users.username": goqu.I("users.username")})).
		Where(goqu.Ex{
			"users.username": username,
			"debts.group_id": groupId,
		}).
		ToSQL()
	rows, err := r.db.Query(query)
	if errors.Is(err, sql.ErrNoRows) {
		return &debts, nil
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var debt dto.DebtDto

		rows.Scan(&debt.Id, &debt.Name, &debt.FromUsername, &debt.ToUsername, &debt.IsClosed, &debt.TotalAmount, &debt.PayedAmount, &debt.CreatedAt, &debt.UpdatedAt, &debt.GroupId)
		debts = append(debts, debt)
	}

	return &debts, nil
}
