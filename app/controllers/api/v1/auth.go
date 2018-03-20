package controllers

import (
	"AiCompServer/app/db"
	"AiCompServer/app/models"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
	"github.com/revel/revel/cache"
	"golang.org/x/crypto/scrypt"
	"gopkg.in/validator.v2"
	"log"
	"strconv"
	"time"
)

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func TokenGenerator(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

type ApiAuth struct {
	ApiV1Controller
}

func (c ApiAuth) GetSessionID() revel.Result {
	session := TokenGenerator(32)
	go cache.Set("session_"+session, session, 2*time.Minute)
	c.Response.Out.Header().Add("Authorization", session)
	r := Response{"Get Session ID"}
	log.Print("&&", session, "&&")
	return c.RenderJSON(r)
}

func (c ApiAuth) SignIn() revel.Result {
	session := c.Request.Header.Get("Authorization")
	if session == "" {
		return c.HandleNotFoundError("Retry Session")
	}
	log.Print("&&", session, "&&")
	var res string
	if err := cache.Get("session_"+session, &res); err != nil {
		return c.HandleBadRequestError("Session Timeout")
	}
	go cache.Delete("session_" + session)

	jsonData := &Auth{}
	if err := c.BindParams(jsonData); err != nil {
		return c.HandleBadRequestError(err.Error())
	}
	log.Print(jsonData.Password)
	salt := []byte("yatuhashi")
	converted, _ := scrypt.Key([]byte(jsonData.Password), salt, 16384, 8, 1, 32)
	password := hex.EncodeToString(converted[:])

	userOld := &models.User{}
	if err := db.DB.Find(&userOld, models.User{Username: jsonData.Username}).Error; err != nil {
		return c.HandleNotFoundError(err.Error())
	}

	if userOld.Password != password {
		return c.HandleNotFoundError("Password is incorrect")
	}

	userNew := &models.User{}
	userNew = userOld
	userNew.Token = TokenGenerator(32)
	if err := validator.Validate(userNew); err != nil {
		return c.HandleBadRequestError(err.Error())
	}
	if err := db.DB.Model(&userOld).Update(&userNew).Error; err != nil {
		return c.HandleNotFoundError(err.Error())
	}
	r := Response{"Success"}
	go cache.Set("auth_"+userNew.Token, userNew.Username, 30*time.Minute)
	c.Response.Out.Header().Add("Authorization", userNew.Token)
	return c.RenderJSON(r)
}

func (c ApiAuth) SignOut() revel.Result {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		return c.HandleNotFoundError("Not SignIn")
	}
	user := &models.User{}
	if err := db.DB.Find(&user, models.User{Token: token}).Error; err != nil {
		return c.HandleNotFoundError(err.Error())
	}
	if err := db.DB.Model(&user).Update("Token", gorm.Expr("NULL")).Error; err != nil {
		return c.HandleBadRequestError(err.Error())
	}
	go cache.Delete("auth_" + token)
	r := Response{"Sign Out"}
	return c.RenderJSON(r)
}

func (c ApiAuth) Role() revel.Result {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		return c.HandleNotFoundError("Not SignIn")
	}
	// Check Token Timeout
	var res string
	if err := cache.Get("auth_"+token, &res); err != nil {
		return c.HandleBadRequestError("Session Timeout")
	}
	go cache.Set("auth_"+user.Token, user.Username, 30*time.Minute)
	user := &models.User{}
	if err := db.DB.Find(&user, models.User{Token: token}).Error; err != nil {
		return c.HandleNotFoundError(err.Error())
	}
	r := Response{user.Role}
	return c.RenderJSON(r)
}

func CheckRole(c ApiV1Controller, AllowRole []string) revel.Result {
	log.Print("CheckRole")
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		return c.HandleNotFoundError("Not SignIn")
	}
	user := &models.User{}
	if err := db.DB.Find(&user, models.User{Token: token}).Error; err != nil {
		return c.HandleNotFoundError("This token does not exist")
	}
	for _, role := range AllowRole {
		if role == user.Role || role == strconv.FormatUint(user.ID, 10) {
			return nil
		}
	}
	r := Response{"Permission Denied"}
	return c.RenderJSON(r)
}

func CheckToken(c ApiV1Controller) revel.Result {
	log.Print("CheckToken")
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		return c.HandleNotFoundError("Not SignIn")
	}
	user := &models.User{}
	if err := db.DB.Find(&user, models.User{Token: token}).Error; err != nil {
		return c.HandleNotFoundError("This token does not exist")
	}
	// Check Token Timeout
	var res string
	if err := cache.Get("auth_"+token, &res); err != nil {
		return c.HandleBadRequestError("Session Timeout")
	}
	go cache.Set("auth_"+user.Token, user.Username, 30*time.Minute)
	return nil
}
