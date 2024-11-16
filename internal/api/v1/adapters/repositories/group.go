package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"example.com/m/internal/api/v1/core/application/dto"
	"github.com/doug-martin/goqu/v9"
)

type GroupRepository struct {
	db *sql.DB
}

func NewGroupRepository(db *sql.DB) *GroupRepository {
	return &GroupRepository{
		db: db,
	}
}

func (r *GroupRepository) Create(d *dto.CreateGroupDto) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return -1, err
	}

	group_stmt, _, _ := goqu.Insert("groups").Cols("invite_code", "created_by", "name").Vals(
		goqu.Vals{d.InviteCode, d.CreatedBy, d.Name},
	).Returning("id").ToSQL()

	var group_id int64
	err = tx.QueryRow(group_stmt).Scan(&group_id)
	if err != nil {
		return -1, err
	}

	group_user_stmt, _, _ := goqu.Insert("groups_users").Cols("username", "group_id").Vals(
		goqu.Vals{d.CreatedBy, group_id},
	).ToSQL()

	_, err = tx.Exec(group_user_stmt)
	if err != nil {
		return -1, err
	}

	err = tx.Commit()
	if err != nil {
		return -1, err
	}

	return group_id, nil
}

func (r *GroupRepository) GetCreatedGroups(username string) (*[]dto.GroupDto, error) {
	var groups []dto.GroupDto
	query, _, _ := goqu.From("groups").Where(goqu.Ex{
		"created_by": username,
	}).ToSQL()
	rows, err := r.db.Query(query)
	if errors.Is(err, sql.ErrNoRows) {
		return &groups, nil
	}
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for rows.Next() {
		var group dto.GroupDto

		rows.Scan(&group.Id, &group.Name, &group.InviteCode, &group.CreatedBy)
		groups = append(groups, group)
	}

	return &groups, nil
}

func (r *GroupRepository) IfInviteCodeExists(ic string) (bool, error) {
	stmt, _, _ := goqu.From("groups").
		Select(
			goqu.COUNT("invite_code"),
		).
		Where(
			goqu.C("invite_code").Eq(ic),
		).ToSQL()

	var inviteCode string
	err := r.db.QueryRow(stmt).Scan(&inviteCode)
	if err != nil {
		return false, err
	}
	if inviteCode == "" {
		return false, nil
	}
	return true, nil
}

func (r *GroupRepository) JoinToGroup(d *dto.GroupDtoComeIn) error {
	group_user_stmt, _, _ := goqu.
		Insert("groups_users").
		Prepared(true).
		Cols("username", "group_id").
		FromQuery(
			goqu.From("groups").
				Select(
					goqu.L("$1", d.Username),
					goqu.Cast(goqu.L("id"), "integer"),
				).Where(goqu.L("invite_code").Eq(goqu.L("$2"))),
		).ToSQL()
	_, err := r.db.Exec(group_user_stmt, d.Username, d.InviteCode)
	if err != nil {
		return err
	}

	return nil
}

func (r *GroupRepository) IsGroupExistsAndUserInGroup(groupId int64, username string) (bool, error) {
	stmt, _, _ := goqu.From("groups_users").
		Select(
			goqu.COUNT("group_id"),
		).
		Where(
			goqu.C("group_id").Eq(groupId),
			goqu.C("username").Eq(username),
		).ToSQL()
	var count int64
	err := r.db.QueryRow(stmt).Scan(&count)
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}
