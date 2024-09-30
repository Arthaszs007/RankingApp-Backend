package models

import "server/dao"

type Event struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

// automatic to search table based on return name
func (Event) TableName() string {
	return "events"
}

// open or close the event based on a id
func SetEventState(event_id int, state bool) {}

// get the current first active is true eventInfo
func GetActiveEvent() (Event, error) {
	var e Event
	err := dao.DB.Model(&Event{}).Where("active = true").First(&e).Error
	return e, err
}
