package controllers

import (
	"AiCompServer/app/db"
	"AiCompServer/app/models"
	"bufio"
	"github.com/revel/revel"
	"gopkg.in/validator.v2"
	"log"
	"os"
	"strconv"
	"strings"
)

type ApiAnswer struct {
	ApiV1Controller
}

type ResponseAnswer struct {
	Answer *models.Answer `json:"answer"`
}

type ResponseAnswers struct {
	Answers []models.Answer `json:"answers"`
}

// Answer Index
func (c ApiAnswer) Index() revel.Result {
	if err := CheckToken(c.ApiV1Controller); err != nil {
		return err
	}
	answers := []models.Answer{}
	if err := db.DB.Order("id desc").Find(&answers).Error; err != nil {
		return c.HandleNotFoundError("Record Find Failure")
	}
	r := Response{ResponseAnswers{answers}}
	return c.RenderJSON(r)
}

// Answer Show
func (c ApiAnswer) Show(id int) revel.Result {
	if err := CheckToken(c.ApiV1Controller); err != nil {
		return err
	}
	answer := &models.Answer{}
	if err := db.DB.First(&answer, id).Error; err != nil {
		return c.HandleNotFoundError(err.Error())
	}
	r := Response{ResponseAnswer{answer}}
	return c.RenderJSON(r)
}

// Answer User's Answer of Some Challenge
func (c ApiAnswer) UserChallengeAnswer(id int) revel.Result {
	if err := CheckToken(c.ApiV1Controller); err != nil {
		return err
	}
	token := c.Request.Header.Get("Authorization")
	user := &models.User{}
	if err := db.DB.Find(&user, models.User{Token: token}).Error; err != nil {
		return c.HandleNotFoundError(err.Error())
	}
	answer := &models.Answer{}
	if err := db.DB.Where("user_id = ? AND challenge_id = ?", user.ID, id).First(&answer).Error; err != nil {
		return c.HandleNotFoundError(err.Error())
	}
	r := Response{ResponseAnswer{answer}}
	return c.RenderJSON(r)
}

// Answer Create
func (c ApiAnswer) Create() revel.Result {
	if err := CheckRole(c.ApiV1Controller, []string{"admin"}); err != nil {
		return err
	}
	if err := CheckToken(c.ApiV1Controller); err != nil {
		return err
	}
	answer := &models.Answer{}
	if err := c.BindParams(answer); err != nil {
		return c.HandleBadRequestError(err.Error())
	}
	token := c.Request.Header.Get("Authorization")
	user := &models.User{}
	if err := db.DB.Find(&user, models.User{Token: token}).Error; err != nil {
		return c.HandleNotFoundError(err.Error())
	}
	if err := validator.Validate(answer); err != nil {
		return c.HandleBadRequestError(err.Error())
	}
	if err := db.DB.Create(answer).Error; err != nil {
		return c.HandleBadRequestError(err.Error())
	}
	r := Response{ResponseAnswer{answer}}
	return c.RenderJSON(r)
}

// Answer Update
func (c ApiAnswer) Update(id int) revel.Result {
	if err := CheckRole(c.ApiV1Controller, []string{"admin"}); err != nil {
		return err
	}
	if err := CheckToken(c.ApiV1Controller); err != nil {
		return err
	}
	answerOld := &models.Answer{}
	if err := db.DB.First(&answerOld, id).Error; err != nil {
		return c.HandleNotFoundError(err.Error())
	}
	answerNew := &models.Answer{}
	if err := c.BindParams(answerNew); err != nil {
		return c.HandleBadRequestError(err.Error())
	}
	if err := validator.Validate(answerNew); err != nil {
		return c.HandleBadRequestError(err.Error())
	}
	if err := db.DB.Model(&answerOld).Update(&answerNew).Error; err != nil {
		return c.HandleNotFoundError(err.Error())
	}
	r := Response{ResponseAnswer{answerNew}}
	return c.RenderJSON(r)
}

// Answer Delete
func (c ApiAnswer) Delete(id int) revel.Result {
	if err := CheckRole(c.ApiV1Controller, []string{"admin"}); err != nil {
		return err
	}
	if err := CheckToken(c.ApiV1Controller); err != nil {
		return err
	}
	answer := &models.Answer{}
	if err := db.DB.First(&answer, id).Error; err != nil {
		return c.HandleNotFoundError(err.Error())
	}
	if err := db.DB.Delete(&answer).Error; err != nil {
		return c.HandleInternalServerError("Record Delete Failure")
	}
	r := Response{"Success Delete"}
	return c.RenderJSON(r)
}

func (c ApiAnswer) Submit(ChallengeID uint64, ansFP *os.File) revel.Result {
	if err := CheckToken(c.ApiV1Controller); err != nil {
		return err
	}
	var fp *os.File
	ansFile := os.Getenv("ANSWERFILE")
	ansFile = ansFile + strconv.Itoa(int(ChallengeID)) + ".csv"
	fp, err := os.Open(ansFile)
	if err != nil {
		return c.HandleBadRequestError("送信されたファイルが開けませんでした")
	}
	scanner1 := bufio.NewScanner(ansFP)
	scanner2 := bufio.NewScanner(fp)
	acc := 0
	a1 := map[string]string{}
	a2 := map[string]string{}
	for scanner1.Scan() && scanner2.Scan() {
		st1 := scanner1.Text()
		st2 := scanner2.Text()
		log.Println(st1)
		log.Println(st2)
		l1 := strings.Split(st1, ",")
		l2 := strings.Split(st2, ",")
		if len(l1) > 1 && len(l2) > 1 {
			a1[l1[0]] = strings.Replace(l1[1], " ", "", -1)
			a2[l2[0]] = strings.Replace(l2[1], " ", "", -1)
		}
		log.Println(st1, l1)
		log.Println(st2, l2)
	}
	for asKey, asVal := range a1 {
		if aaVal, err := a2[asKey]; err == true {
			if aaVal == asVal {
				acc = acc + 1
			}
		}
	}
	if err := scanner1.Err(); err != nil {
		return c.HandleBadRequestError("採点中に解答ファイルにエラーが起きました")
	}
	if err := scanner2.Err(); err != nil {
		return c.HandleBadRequestError("採点中に正解ファイルにエラーが起きました")
	}
	// Submitしたユーザーを特定する
	token := c.Request.Header.Get("Authorization")
	user := &models.User{}
	if err := db.DB.Find(&user, models.User{Token: token}).Error; err != nil {
		return c.HandleNotFoundError(err.Error())
	}
	// まずはDBを探してあったらUpdate処理、なかったらCreate
	// そのユーザーがSubmitした問題のanswerを探す
	answer := &models.Answer{}
	if err := db.DB.Where("user_id = ? AND challenge_id = ?", user.ID, ChallengeID).First(&answer).Error; err != nil {
		// なかった場合のCreate
		answer.ChallengeID = ChallengeID
		answer.UserID = user.ID
		answer.Score = acc
		if err := validator.Validate(answer); err != nil {
			return c.HandleBadRequestError(err.Error())
		}
		if err := db.DB.Create(answer).Error; err != nil {
			return c.HandleBadRequestError(err.Error())
		}
		r := Response{"Success Submit"}
		return c.RenderJSON(r)
	}
	// そのユーザーがSubmitした問題のanswerを更新する
	answerNew := &models.Answer{}
	answerNew.ChallengeID = ChallengeID
	answerNew.UserID = user.ID
	if answer.Score < acc {
		answerNew.Score = acc
	}
	if err := validator.Validate(answerNew); err != nil {
		return c.HandleBadRequestError(err.Error())
	}
	if err := db.DB.Model(&answer).Update(&answerNew).Error; err != nil {
		return c.HandleNotFoundError(err.Error())
	}
	r := Response{"Success Submit"}
	return c.RenderJSON(r)
}
