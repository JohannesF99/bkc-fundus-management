package member

import (
	"github.com/JohannesF99/bkc-fundus-management/pkg/models"
)

func getAllMembers() ([]models.Member, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}
	members, err := db.GetAllMembers()
	if err != nil {
		return nil, err
	}
	return members, nil
}

func getMemberWithUserId(userId int) (models.Member, error) {
	db, err := Connect()
	if err != nil {
		return models.Member{}, err
	}
	member, err := db.GetMemberWithId(userId)
	if err != nil {
		return models.Member{}, err
	}
	return member, nil
}

func insertNewMember(newMember models.NewMemberInfos) (int64, error) {
	db, err := Connect()
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

func updateMemberBorrowCount(userId int, diff int) (int, error) {
	db, _ := Connect()
	count, err := db.UpdateBorrowedItemCount(userId, diff)
	if err != nil {
		return -1, err
	}
	return count, nil
}

func changeMemberStatus(userId int, status bool) (models.Member, error) {
	db, err := Connect()
	if err != nil {
		return models.Member{}, err
	}
	err = db.ChangeMemberStatus(userId, status)
	if err != nil {
		return models.Member{}, err
	}
	member, err := db.GetMemberWithId(userId)
	if err != nil {
		return models.Member{}, err
	}
	return member, nil
}
