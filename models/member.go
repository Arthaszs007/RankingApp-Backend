package models

import "server/dao"

type Member struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Img_url string `json:"img_url"`
	Bio     string `json:"bio"`
}

// automatic to search table based on return name
func (Member) TableName() string {
	return "members"
}

//receive a member id slice to search in db and return a slice of result
func GetMembersInfo(member_id []int) ([]Member, error) {
	var members []Member
	err := dao.DB.Model(&Member{}).Where("id IN ?", member_id).Order("id ASC").Find(&members).Error
	return members, err
}
