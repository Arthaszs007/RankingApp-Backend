package controllers

import (
	"fmt"
	"server/models"
	"server/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

// inquire the user exist or not by username
func (u UserController) GetUserInfoByUsername(c *gin.Context) {
	username := c.Param("username")

	user, err := models.GetUserInfoByUsername(username)

	if err != nil {
		msg := fmt.Sprintf("%s is not exist", username)
		ReturnError(c, 400, msg)
		return
	}
	msg := fmt.Sprintf("%s is exist", username)
	ReturnSuccess(c, 0, msg, user, 1)
}

// user login by api,1. check the username and password 2. create the jwt and send to client
func (u UserController) Login(c *gin.Context) {

	var data map[string]interface{}
	if err := c.BindJSON(&data); err != nil {
		ReturnError(c, 400, "get params failed")
		return
	}
	username := data["username"].(string)
	password := data["password"].(string)
	//check content is empty
	if username == "" || password == "" {
		ReturnError(c, 400, "content can't be empty")
		return
	}
	// get user info by username and match password
	user, err := models.GetUserInfoByUsername(username)
	if err != nil {
		ReturnError(c, 400, "Username does not exist")
		return
	} else if !VerifyHash(user.Password, password) {

		ReturnError(c, 400, "Password is wrong")
		return
	}
	//create a jwt
	jwtTonke, err := jwt.GenerateJWT(user.Username, CreateTimestamp(1))
	if err != nil {
		ReturnError(c, 400, "jwtToken is failed to create")
		fmt.Println(err)
		return
	}

	// send the jwt to client
	ReturnSuccess(c, 200, "Log in successfully", jwtTonke, 1)

}

// user register by api 1. check the input 2. exist check 3.hash and shorage the password
func (u UserController) Register(c *gin.Context) {

	var data map[string]interface{}
	if err := c.BindJSON(&data); err != nil {
		ReturnError(c, 400, "get params failed")
		return
	}
	username := data["username"].(string)
	password := data["password"].(string)
	repeat := data["repeat"].(string)
	//empty check ,password repeat check and  length check
	if username == "" || password == "" || repeat == "" {
		ReturnError(c, 400, "content shouldn't be empty")
		return
	} else if password != repeat {
		ReturnError(c, 400, "twice password should be same")
		return
	} else if len(username) > 20 || len(password) > 20 || len(repeat) > 20 {
		ReturnError(c, 400, "username and password can't be over 20 letters")
		return
	} else if len(username) < 6 || len(password) < 6 || len(repeat) < 6 {
		ReturnError(c, 400, "username and password  can't be less 6 letters")
		return
	}
	//username exist check
	_, err := models.GetUserInfoByUsername(username)

	if err == nil {
		msg := fmt.Sprintf("%s is exist", username)
		ReturnError(c, 400, msg)
		return
	}
	// hash password
	psd, hashErr := HashStr(password)
	if hashErr != nil {
		ReturnError(c, 400, "password is failed to hash")
		return
	}

	// push to database by model
	userId, addErr := models.AddUser(username, psd)
	if addErr != nil {
		ReturnError(c, 400, "Register is failed")
		return
	}
	ReturnSuccess(c, 200, "Register is done", userId, 1)

}

// verify the token coming from the api headers
func (u UserController) Verify(c *gin.Context) {

	// get token and remove the "Bearer" in front if exist
	tokenStr := c.GetHeader("Authorization")

	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	// get username and send to client if passed
	username, err := jwt.VerifyJWT(tokenStr)
	if err != nil {
		ReturnError(c, 400, err.Error())
		return
	}

	ReturnSuccess(c, 200, "token verify passed", username, 1)

}
