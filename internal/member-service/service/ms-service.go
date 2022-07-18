package service

import (
	"github.com/JohannesF99/bkc-fundus-management/internal/member-service/database"
	"github.com/JohannesF99/bkc-fundus-management/pkg/models"
)

func GetAllMember() ([]models.Member, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}
	members, err := db.GetAllMembers()
	if err != nil {
		return nil, err
	}
	return members, nil
}

func GetMemberWithId(userId int) (models.Member, error) {
	db, err := database.Connect()
	if err != nil {
		return models.Member{}, err
	}
	member, err := db.GetMemberWithId(userId)
	if err != nil {
		return models.Member{}, err
	}
	return member, nil
}

func InsertNewMember(newMember models.Member) (int64, error) {
	db, err := database.Connect()
	if err != nil {
		return -1, err
	}
	id, err := db.CreateMember(models.Member{
		Name:    newMember.Name,
		Comment: newMember.Comment,
	})
	if err != nil {
		return -1, err
	}
	return id, nil
}

func UpdateBorrowCount(userId int, diff int) (int, error) {
	db, _ := database.Connect()
	count, err := db.UpdateBorrowedItemCount(userId, diff)
	if err != nil {
		return -1, err
	}
	return count, nil
}
