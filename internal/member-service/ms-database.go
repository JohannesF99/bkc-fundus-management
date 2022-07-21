package member

import (
	"database/sql"
	"github.com/JohannesF99/bkc-fundus-management/pkg/models"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"time"
)

type DB struct {
	db *sql.DB
}

func connect() (DB, error) {
	db, err := sql.Open("mysql", "root:2678@/bkc?parseTime=true")
	if err != nil {
		return DB{}, models.Error{
			Details: err.Error(),
			Path:    "Member Service - Database Connection",
			Object:  "",
			Time:    time.Now(),
		}
	}
	return DB{db}, nil
}

func (db DB) GetAllMembers() ([]models.Member, error) {
	var members []models.Member
	tx, err := db.db.Begin()
	defer tx.Commit()
	if err != nil {
		return nil, models.Error{
			Details: err.Error(),
			Path:    "Member Service - GetAllMembers()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	rows, err := tx.Query("SELECT * FROM bkc.members")
	defer rows.Close()
	if err != nil {
		return nil, models.Error{
			Details: err.Error(),
			Path:    "Member Service - GetAllMembers()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	for rows.Next() {
		var newMember models.Member
		err := rows.Scan(
			&newMember.Id,
			&newMember.Name,
			&newMember.BorrowedItemCount,
			&newMember.Comment,
			&newMember.Active,
			&newMember.Created,
			&newMember.Modified)
		if err != nil {
			return nil, models.Error{
				Details: err.Error(),
				Path:    "Member Service - GetAllMembers()",
				Object:  "",
				Time:    time.Now(),
			}
		}
		members = append(members, newMember)
	}
	return members, nil
}

func (db DB) GetMemberWithId(userId int) (models.Member, error) {
	tx, err := db.db.Begin()
	defer tx.Commit()
	if err != nil {
		return models.Member{}, models.Error{
			Details: err.Error(),
			Path:    "Member Service - GetMemberWithId()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	rows, err := tx.Query("SELECT * FROM bkc.members WHERE id=?", userId)
	defer rows.Close()
	if err != nil {
		return models.Member{}, models.Error{
			Details: err.Error(),
			Path:    "Member Service - GetMemberWithId()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	var member models.Member
	rows.Next()
	err = rows.Scan(
		&member.Id,
		&member.Name,
		&member.BorrowedItemCount,
		&member.Comment,
		&member.Active,
		&member.Created,
		&member.Modified)
	if err != nil {
		return models.Member{}, models.Error{
			Details: err.Error(),
			Path:    "Member Service - GetMemberWithId()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	return member, nil
}

func (db DB) CreateMember(member models.Member) (int64, error) {
	tx, err := db.db.Begin()
	defer tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return -1, models.Error{
			Details: err.Error(),
			Path:    "Member Service - CreateMember()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	res, err := tx.Exec(
		"INSERT INTO bkc.members(name, comment) VALUES (?,?)",
		member.Name, member.Comment)
	if err != nil {
		_ = tx.Rollback()
		return -1, models.Error{
			Details: err.Error(),
			Path:    "Member Service - CreateMember()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	id, err := res.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return -1, models.Error{
			Details: err.Error(),
			Path:    "Member Service - CreateMember()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	return id, nil
}

func (db DB) UpdateBorrowedItemCount(userId int, diff int) (int, error) {
	tx, err := db.db.Begin()
	defer tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return -1, models.Error{
			Details: err.Error(),
			Path:    "Member Service - UpdateBorrowedItemCount()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	res, err := tx.Exec(
		"UPDATE bkc.members SET borrowed_item_count=borrowed_item_count+? WHERE id=?",
		diff, userId)
	if err != nil {
		_ = tx.Rollback()
		return -1, models.Error{
			Details: err.Error(),
			Path:    "Member Service - UpdateBorrowedItemCount()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return -1, models.Error{
			Details: err.Error(),
			Path:    "Member Service - UpdateBorrowedItemCount()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	if rowsAffected != 1 {
		_ = tx.Rollback()
		return -1, models.Error{
			Details: "Rows Affected: " + strconv.Itoa(int(rowsAffected)),
			Path:    "Member Service - UpdateBorrowedItemCount()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	return userId, nil
}

func (db DB) ChangeMemberStatus(userId int, status bool) error {
	tx, err := db.db.Begin()
	defer tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return models.Error{
			Details: err.Error(),
			Path:    "Member Service - ChangeMemberStatus()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	res, err := tx.Exec("UPDATE bkc.members SET active=? WHERE id=?", status, userId)
	if err != nil {
		_ = tx.Rollback()
		return models.Error{
			Details: err.Error(),
			Path:    "Member Service - ChangeMemberStatus()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return models.Error{
			Details: err.Error(),
			Path:    "Member Service - ChangeMemberStatus()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	if rowsAffected != 1 {
		_ = tx.Rollback()
		return models.Error{
			Details: "Rows Affected: " + strconv.Itoa(int(rowsAffected)),
			Path:    "Member Service - ChangeMemberStatus()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	return nil
}
