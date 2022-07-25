package member

import (
	"github.com/JohannesF99/bkc-fundus-management/pkg/models"
	"time"
)

func getAllMembers() ([]models.Member, error) {
	db, err := connect()
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
	db, err := connect()
	if err != nil {
		return models.Member{}, err
	}
	member, err := db.GetMemberWithId(userId)
	if err != nil {
		return models.Member{}, err
	}
	if !member.Active {
		return models.Member{}, models.Error{
			Details: "Member does exist, but is currently marked as inactive",
			Path:    "Member Service - getMemberWithUserId()",
			Object:  member.String(),
			Time:    time.Now(),
		}
	}
	return member, nil
}

func insertNewMember(newMember models.NewMemberInfos) (int64, error) {
	db, err := connect()
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
	db, _ := connect()
	count, err := db.UpdateBorrowedItemCount(userId, diff)
	if err != nil {
		return -1, err
	}
	return count, nil
}

func changeMemberStatus(memberId int, status bool) (models.Member, error) {
	db, err := connect()
	if err != nil {
		return models.Member{}, err
	}
	member, err := db.GetMemberWithId(memberId)
	if err != nil {
		return models.Member{}, err
	}
	if member.BorrowedItemCount > 0 {
		return models.Member{}, models.Error{
			Details: "User-Status cannot be changed until every borrowed Item has been returned",
			Path:    "Member Service - changeMemberStatus",
			Object:  member.String(),
			Time:    time.Now(),
		}
	}
	err = db.ChangeMemberStatus(memberId, status)
	if err != nil {
		return models.Member{}, err
	}
	member, err = db.GetMemberWithId(memberId)
	if err != nil {
		return models.Member{}, err
	}
	return member, nil
}
