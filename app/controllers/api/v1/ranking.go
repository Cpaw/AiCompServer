package controllers

import (
	"AiCompServer/app/db"
	"AiCompServer/app/models"
	"github.com/revel/revel"
	"sort"
)

type Rank struct {
	Rank     int     `json:"rank"`
	Username string  `json:"username"`
	Score    float64 `json:"score"`
}

type Ranks []Rank

type ResponseRanking struct {
	Ranking Ranks `json:"ranking"`
}

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
		score := 0.0
		if err := db.DB.Find(&answer, "user_id = ?", user.ID).Error; err != nil {
			return c.HandleNotFoundError(err.Error())
		}
		for _, ans := range answer {
			score = score + ans.Score
		}
		rank = append(rank, Rank{Rank: 0, Username: user.Username, Score: score})
	}
	sort.Slice(rank, func(i, j int) bool { return rank[i].Score > rank[j].Score })
	Rindex := 1
	TmpScore := 0.0
	for Tindex, _ := range rank {
		// スコアが前の人と一緒の時は前の人の順位と同じにする
		if rank[Tindex].Score == TmpScore {
			rank[Tindex].Rank = Rindex
			TmpScore = rank[Tindex].Score
		} else {
			rank[Tindex].Rank = Tindex + 1
			Rindex = Tindex + 1
			TmpScore = rank[Tindex].Score
		}
	}
	r := Response{ResponseRanking{rank}}
	return c.RenderJSON(r)
}
