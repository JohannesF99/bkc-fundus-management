package member

import (
	"database/sql"
	"github.com/JohannesF99/bkc-fundus-management/pkg/models"
	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	db *sql.DB
}

func connect() (DB, error) {
	db, err := sql.Open("mysql", "root:2678@/bkc?parseTime=true")
	if err != nil {
		return DB{}, err
	}
	return DB{db}, nil
}

func (db DB) GetAllMembers() ([]models.Member, error) {
	var members []models.Member
	tx, err := db.db.Begin()
	defer tx.Commit()
	if err != nil {
		return nil, err
	}
	rows, err := tx.Query("SELECT * FROM bkc.members")
	defer rows.Close()
	if err != nil {
		return nil, err
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
			return nil, err
		}
		members = append(members, newMember)
	}
	return members, nil
}

func (db DB) GetMemberWithId(userId int) (models.Member, error) {
	tx, err := db.db.Begin()
	defer tx.Commit()
	if err != nil {
		return models.Member{}, err
	}
	rows, err := tx.Query("SELECT * FROM bkc.members WHERE id=?", userId)
	defer rows.Close()
	if err != nil {
		return models.Member{}, err
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
		return models.Member{}, err
	}
	return member, nil
}

func (db DB) CreateMember(member models.Member) (int64, error) {
	tx, err := db.db.Begin()
	defer tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return -1, err
	}
	res, err := tx.Exec(
		"INSERT INTO bkc.members(name, comment) VALUES (?,?)",
		member.Name, member.Comment)
	if err != nil {
		_ = tx.Rollback()
		return -1, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return -1, err
	}
	return id, nil
}

func (db DB) UpdateBorrowedItemCount(userId int, diff int) (int, error) {
	tx, err := db.db.Begin()
	defer tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return -1, err
	}
	_, err = tx.Exec(
		"UPDATE bkc.members SET borrowed_item_count=borrowed_item_count+? WHERE id=?",
		diff, userId)
	if err != nil {
		_ = tx.Rollback()
		return -1, err
	}
	return userId, nil
}

func (db DB) ChangeMemberStatus(userId int, status bool) error {
	tx, err := db.db.Begin()
	defer tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	res, err := tx.Exec("UPDATE bkc.members SET active=? WHERE id=?", status, userId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected != 1 {
		_ = tx.Rollback()
		return err
	}
	return nil
}
