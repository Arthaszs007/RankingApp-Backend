package controllers

import (
	"fmt"
	"server/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EventController struct{}

type MembersWithVotes struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Img_url string `json:"img_url"`
	Bio     string `json:"bio"`
	Votes   int    `json:"votes"`
}

// 1. get active evnet info 2. get the members info based on event info 3. push to client
func (e EventController) GetList(c *gin.Context) {
	// get active event's info
	event, err := models.GetActiveEvent()
	if err != nil {
		msg := ("get active event faild")
		ReturnError(c, 400, msg)
		return
	}

	// get the slice of [members id,votes] based on event id
	eventInfos, err := models.GetMembersAndVotesByEventID(event.Id)
	if err != nil {
		msg := ("get members and votes failed")
		ReturnError(c, 400, msg)
		return
	}
	// splite the event info and storage into 2 slice
	var memberIDs []int
	var votes []int

	for _, item := range eventInfos {
		memberIDs = append(memberIDs, item.Member_id)
		votes = append(votes, item.Votes)
	}
	// get members info based on memberIDS
	membersInfo, err := models.GetMembersInfo(memberIDs)
	if err != nil {
		msg := ("get members info faild")
		ReturnError(c, 400, msg)
		return
	}

	var res []MembersWithVotes
	for i := range membersInfo {
		res = append(res, MembersWithVotes{
			Id:      membersInfo[i].Id,
			Name:    membersInfo[i].Name,
			Img_url: membersInfo[i].Img_url,
			Bio:     membersInfo[i].Bio,
			Votes:   votes[i],
		})
	}

	// members,err := models.GetMembersInfo()
	ReturnSuccess(c, 200, "success", res, 1)
}

// 1. get the active event id 2. get user id to check in table of votes whether already exist 3. if not ,votes increase 1 by member id 4. write the vote record in table of votes
func (e EventController) VoteToMember(c *gin.Context) {
	// get active event's info
	event, err := models.GetActiveEvent()
	if err != nil {
		msg := ("get active event faild")
		ReturnError(c, 400, msg)
		return
	}

	// get the daya from json body
	var data map[string]interface{}
	if err := c.BindJSON(&data); err != nil {
		ReturnError(c, 400, "get params failed")
		return
	}
	//get user id and member id from api body and transfer string to int
	memberstr := data["member_id"].(string)
	member_id, _ := strconv.Atoi(memberstr)

	userName := data["username"].(string)

	userInfo, err := models.GetUserInfoByUsername(userName)
	if err != nil {
		ReturnError(c, 400, "Get user info failed")
		return
	}

	// check the user vote record,if exsited, return the message
	res, _ := models.CheckVotesExist(event.Id, userInfo.Id)
	fmt.Println(res)
	if res {
		ReturnError(c, 400, "User already voted ")
		return
	}

	// invoke func to increase the votes
	err1 := models.IncreaseVote(event.Id, member_id)
	if err1 != nil {
		ReturnError(c, 400, "increase votes failed")
		return
	}
	// invoke func to write record of vote in database
	err2 := models.RecordVotes(event.Id, userInfo.Id, member_id)
	if err2 != nil {
		ReturnError(c, 400, "Write record of vote in table of votes failed ")
		return
	}

	ReturnSuccess(c, 200, "success", member_id, 1)
}
