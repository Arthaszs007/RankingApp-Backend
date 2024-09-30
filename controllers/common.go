package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// json construction
type JsonStruct struct {
	Code  int         `json:"code"`
	Msg   interface{} `json:"msg"`
	Data  interface{} `json:"data"`
	Count int64       `json:"count"`
}

// return the success json content
func ReturnSuccess(c *gin.Context, code int, msg interface{}, data interface{}, count int64) {
	json := &JsonStruct{Code: code, Msg: msg, Data: data, Count: count}
	c.JSON(200, json)
}

// return the error json content
func ReturnError(c *gin.Context, code int, msg interface{}) {
	json := &JsonStruct{Code: code, Msg: msg}
	c.JSON(200, json)
}

// hash a string by bcrypto with a salt valu,will return the error if has
func HashStr(s string) (string, error) {
	res, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// verfiy the origin string with hash string, if match return true,or false,
func VerifyHash(hash string, s string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(s))
	return err == nil
}

// create a timestamp,recieve a hour as param and return the new timestamp
func CreateTimestamp(hour int) time.Time {
	now := time.Now()

	newTime := now.Add(time.Duration(hour) * time.Hour)

	return newTime
}
