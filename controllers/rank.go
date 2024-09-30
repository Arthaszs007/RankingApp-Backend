package controllers

import (
	"fmt"
	"server/cache"
	"server/models"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RankController struct{}

type VotesRank struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Img_url string `json:"img_url"`
	Bio     string `json:"bio"`
	Votes   int    `json:"votes"`
}

func (r RankController) GetRank(c *gin.Context) {
	// get the active event id
	event, err := models.GetActiveEvent()
	if err != nil {
		msg := ("get active event faild")
		ReturnError(c, 400, msg)
		return
	}
	// combin the redisKey
	redisKey := "rank" + string(event.Id)

	// try to check value in redis, if data exist , return data or visit database
	val, err := getFromRedis(redisKey)
	if err == nil {
		ReturnSuccess(c, 200, "redis success", val, 1)
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
	// combin the data slice
	var res []VotesRank
	for i := range membersInfo {
		res = append(res, VotesRank{
			Id:      membersInfo[i].Id,
			Name:    membersInfo[i].Name,
			Img_url: membersInfo[i].Img_url,
			Bio:     membersInfo[i].Bio,
			Votes:   votes[i],
		})
	}

	// based on "votes" to sort from big to small
	sort.Slice(res, func(i, j int) bool {
		return res[i].Votes > res[j].Votes
	})
	// store the data into redis
	storeToRedis(redisKey, res, 10)
	ReturnSuccess(c, 200, "success", res, 0)
}

// store the data of votesRank type to redis, need pass a key , data slice and a timer as minute
func storeToRedis(key string, VotesInfo []VotesRank, timer int) error {
	// instance a pipeline
	pipe := cache.Rdb.Pipeline()

	// loop to store data to redis
	for _, item := range VotesInfo {
		// store the votes and id as a map
		pipe.ZAdd(cache.Rctx, key, redis.Z{
			Score:  float64(item.Votes),
			Member: strconv.Itoa(item.Id),
		})

		//use a key type to store the hash info
		hashKey := fmt.Sprintf("rank:%d", item.Id)
		pipe.HSet(cache.Rctx, hashKey,
			"name", item.Name,
			"img_url", item.Img_url,
			"bio", item.Bio)

		// add a expire on hash info
		pipe.Expire(cache.Rctx, hashKey, time.Duration(timer)*time.Minute)

	}
	// add a expire on id-votes map
	pipe.Expire(cache.Rctx, key, time.Duration(timer)*time.Minute)

	// run the pipeline and check the error
	_, err := pipe.Exec(cache.Rctx)
	if err != nil {
		return err
	}
	return nil
}

// get data from the redis
func getFromRedis(key string) ([]VotesRank, error) {
	// try to get ids from the redis, if empty return error
	memberIDS, err := cache.Rdb.ZRevRange(cache.Rctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	} else if len(memberIDS) == 0 {
		return nil, fmt.Errorf("no members found for key: %s", key)
	}

	// instance a pipeline
	pipe := cache.Rdb.Pipeline()
	// a slice to hold the command results
	cmds := make([]*redis.SliceCmd, len(memberIDS))      // to hold hash info data
	votesCmds := make([]*redis.FloatCmd, len(memberIDS)) // to hold votes data
	memberIDsInt := make([]int, len(memberIDS))          // to hold id data

	// loop to get all data and store them
	for i, id := range memberIDS {
		memberid, _ := strconv.Atoi(id)
		memberIDsInt[i] = memberid                                                                   // Store the integer ID
		cmds[i] = pipe.HMGet(cache.Rctx, fmt.Sprintf("rank:%d", memberid), "name", "img_url", "bio") // reverse the hash info
		votesCmds[i] = pipe.ZScore(cache.Rctx, key, id)
	}

	//run the pipeline
	_, err = pipe.Exec(cache.Rctx)
	if err != nil {
		return nil, err
	}

	// loop to deal with the data
	var members []VotesRank
	for i, cmd := range cmds {
		data, err := cmd.Result()
		if err != nil {
			return nil, err
		}

		// get votes score
		votes, err := votesCmds[i].Result()
		if err != nil {
			return nil, err
		}

		// combine the data into a structure
		member := VotesRank{
			Id:      memberIDsInt[i], // Use stored member ID
			Name:    data[0].(string),
			Img_url: data[1].(string),
			Bio:     data[2].(string),
			Votes:   int(votes),
		}
		members = append(members, member)
	}

	// return the data if successfully
	return members, nil
}
