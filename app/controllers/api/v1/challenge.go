package controllers

import (
	"AiCompServer/app/db"
	"AiCompServer/app/models"
	"github.com/revel/revel"
	"gopkg.in/validator.v2"
)

type ApiChallenge struct {
	ApiV1Controller
}

type ResponseChallenge struct {
	Challenge *models.Challenge `json:"challenge"`
}

type ResponseChallenges struct {
	Challenges []models.Challenge `json:"challenges"`
}

func (c ApiChallenge) Index() revel.Result {
	if err := CheckToken(c.ApiV1Controller); err != nil {
		return err
	}
	challenges := []models.Challenge{}
	if err := db.DB.Order("id desc").Find(&challenges).Error; err != nil {
		return c.HandleNotFoundError("Record Find Failure")
	}
	r := Response{ResponseChallenges{challenges}}
	return c.RenderJSON(r)
}

func (c ApiChallenge) Show(id int) revel.Result {
	if err := CheckToken(c.ApiV1Controller); err != nil {
		return err
	}
	challenge := &models.Challenge{}
	if err := db.DB.First(&challenge, id).Error; err != nil {
		return c.HandleNotFoundError(err.Error())
	}
	r := Response{ResponseChallenge{challenge}}
	return c.RenderJSON(r)
}

func (c ApiChallenge) Create() revel.Result {
	if err := CheckRole(c.ApiV1Controller, []string{"admin"}); err != nil {
		return err
	}
	if err := CheckToken(c.ApiV1Controller); err != nil {
		return err
	}
	challenge := &models.Challenge{}
	if err := c.BindParams(challenge); err != nil {
		return c.HandleBadRequestError(err.Error())
	}
	if err := validator.Validate(challenge); err != nil {
		return c.HandleBadRequestError(err.Error())
	}
	if err := db.DB.Create(challenge).Error; err != nil {
		return c.HandleBadRequestError(err.Error())
	}
	r := Response{ResponseChallenge{challenge}}
	return c.RenderJSON(r)
}

func (c ApiChallenge) Update(id int) revel.Result {
	if err := CheckRole(c.ApiV1Controller, []string{"admin"}); err != nil {
		return err
	}
	if err := CheckToken(c.ApiV1Controller); err != nil {
		return err
	}
	challengeOld := &models.Challenge{}
	if err := db.DB.First(&challengeOld, id).Error; err != nil {
		return c.HandleNotFoundError(err.Error())
	}
	challengeNew := &models.Challenge{}
	if err := c.BindParams(challengeNew); err != nil {
		return c.HandleBadRequestError(err.Error())
	}
	if err := validator.Validate(challengeNew); err != nil {
		return c.HandleBadRequestError(err.Error())
	}
	if err := db.DB.Model(&challengeOld).Update(&challengeNew).Error; err != nil {
		return c.HandleNotFoundError(err.Error())
	}
	r := Response{ResponseChallenge{challengeNew}}
	return c.RenderJSON(r)
}

func (c ApiChallenge) Delete(id int) revel.Result {
	if err := CheckRole(c.ApiV1Controller, []string{"admin"}); err != nil {
		return err
	}
	if err := CheckToken(c.ApiV1Controller); err != nil {
		return err
	}
	challenge := &models.Challenge{}
	if err := db.DB.First(&challenge, id).Error; err != nil {
		return c.HandleNotFoundError(err.Error())
	}
	if err := db.DB.Delete(&challenge).Error; err != nil {
		return c.HandleInternalServerError("Record Delete Failure")
	}
	r := Response{"Success Delete"}
	return c.RenderJSON(r)
}
