package models

import (
	"server/dao"

	"gorm.io/gorm"
)

type Events_members struct {
	Event_id  int `json:"event_id"`
	Member_id int `json:"member_id"`
	Votes     int `json:"votes"`
}

type EventsInfo struct {
	Member_id int `json:"member_id"`
	Votes     int `json:"votes"`
}

// automatic to search table based on return name
func (Events_members) TableName() string {
	return "events_members"
}

// pass a event_id to get info
func GetMembersAndVotesByEventID(event_id int) ([]EventsInfo, error) {
	var eventInfos []EventsInfo
	err := dao.DB.Model(&Events_members{}).Where("event_id = ?", event_id).Select("member_id,votes").Find(&eventInfos).Error

	return eventInfos, err
}

// receive 2 params to filter and add 1 on votes
func IncreaseVote(event_id int, member_id int) error {

	err := dao.DB.Model(&Events_members{}).Where("event_id = ? AND  member_id = ?", event_id, member_id).UpdateColumn("votes", gorm.Expr("votes + ?", 1)).Error

	return err
}
