package models

import (
	"server/dao"
)

type Vote struct {
	User_id   int
	Event_id  int
	Member_id int
}

// automatic to search table based on return name
func (Vote) TableName() string {
	return "votes"
}

// receive 3 params and based on result return true or false with error if have (1 per 1 time in same event)
func CheckVotesExist(event_id int, user_id int) (bool, error) {
	var count int64
	err := dao.DB.Model(&Vote{}).Where(" user_id = ? AND event_id = ?", user_id, event_id).Count(&count).Error
	if count > 0 {
		return true, err
	}
	return false, err
}

// receive 3 params and write in table of "votes"
func RecordVotes(event_id int, user_id int, member_id int) error {
	vote := Vote{User_id: user_id, Event_id: event_id, Member_id: member_id}
	err := dao.DB.Model(&Vote{}).Create(&vote).Error
	return err
}
