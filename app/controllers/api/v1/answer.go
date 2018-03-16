package controllers

import (
	"AiCompServer/app/db"
	"AiCompServer/app/models"
	"github.com/revel/revel"
	"gopkg.in/validator.v2"
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
	token := c.Request.Header.Get("authentication")
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
	token := c.Request.Header.Get("authentication")
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

func (c ApiAnswer) Submit() revel.Result {
	// まずはDBを探してあったらUpdate処理、なかったらCreate
	r := Response{"Success Submit"}
	return c.RenderJSON(r)
}
