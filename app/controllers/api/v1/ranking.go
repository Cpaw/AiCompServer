package controllers

import (
	"AiCompServer/app/db"
	"AiCompServer/app/models"
	"github.com/revel/revel"
	"sort"
)

type Rank struct {
	rank     int
	username string
	score    int
}

type Ranks []Rank

func (c ApiChallenge) Ranking() revel.Result {
	if err := CheckToken(c.ApiV1Controller); err != nil {
		return err
	}
	users := []models.User{}
	if err := db.DB.Order("id desc").Find(&users).Error; err != nil {
		return c.HandleNotFoundError("Record Find Failure")
	}
	answer := []models.Answer{}
	var rank Ranks
	for _, user := range users {
		if user.Role == "admin" {
			continue
		}
		score := 0
		if err := db.DB.Find(&answer, "user_id = ?", user.ID).Error; err != nil {
			return c.HandleNotFoundError(err.Error())
		}
		for _, ans := range answer {
			score = score + ans.Score
		}
		rank = append(rank, Rank{rank: 0, username: user.Username, score: score})
	}
	sort.Slice(rank, func(i, j int) bool { return rank[i].score > rank[j].score })
	for index, r := range rank {
		r.rank = index
	}
	r := Response{rank}
	return c.RenderJSON(r)
}
