package models

import "server/dao"

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Register string `json:"register"`
	Ip       string `json:"ip"`
}

// automatic to search table based on return name
func (User) TableName() string {
	return "users"
}

//get user info based on ID
func GetUserInfoByUsername(username string) (User, error) {
	var user User
	err := dao.DB.Model(&User{}).Where("username = ?", username).First(&user).Error
	return user, err
}

// add user info
func AddUser(username string, password string) (int, error) {
	user := User{Username: username, Password: password}
	err := dao.DB.Model(&User{}).Create(&user).Error
	return user.Id, err
}
