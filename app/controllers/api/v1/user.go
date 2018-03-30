package controllers

import (
	"AiCompServer/app/db"
	"AiCompServer/app/models"
	"encoding/hex"
	"github.com/revel/revel"
	"golang.org/x/crypto/scrypt"
	"gopkg.in/validator.v2"
	"strconv"
)

type ApiUser struct {
	ApiV1Controller
}

type ReseposeUser struct {
	User *models.User `json:"user"`
}

type ReseposeUsers struct {
	Users []models.User `json:"users"`
}

/* User List
JSON
 Request
  URL : /GET     /api/v1/user
 Response
    {
     "results": [
       {
         "id": 2,
         "created_at": "2018-01-19T01:48:56.00547399+09:00",
         "username": "tester2",
         "password": "41ab4ad97e3cc6a9c52375694dc50df9ca8d99ba0ebf904864a880e586574587",
         "role": "test",
         "token": "",
         "updated_at": "0001-01-01T00:00:00Z"
       },
       {
         "id": 1,
         "created_at": "2018-01-19T01:48:53.109434454+09:00",
         "username": "tester1",
         "password": "41ab4ad97e3cc6a9c52375694dc50df9ca8d99ba0ebf904864a880e586574587",
         "role": "test",
         "token": "",
         "updated_at": "0001-01-01T00:00:00Z"
       }
     ]
    }
*/
func (c ApiUser) Index() revel.Result {
	if err := CheckRole(c.ApiV1Controller, []string{"admin"}); err != nil {
		return err
	}
	if err := CheckToken(c.ApiV1Controller); err != nil {
		return err
	}
	users := []models.User{}
	if err := db.DB.Order("id desc").Find(&users).Error; err != nil {
		return c.HandleNotFoundError("Record Find Failure")
	}
	r := Response{ReseposeUsers{users}}
	return c.RenderJSON(r)
}

/* User Detail
JSON
 Request
  URL : /GET     /api/v1/user/:id
 Response
    {
     "results": {
       "id": 1,
       "created_at": "2018-01-19T01:48:53.109434454+09:00",
       "username": "tester1",
       "password": "41ab4ad97e3cc6a9c52375694dc50df9ca8d99ba0ebf904864a880e586574587",
       "role": "test",
       "token": "",
       "updated_at": "0001-01-01T00:00:00Z"
     }
    }
*/
func (c ApiUser) Show(id int) revel.Result {
	if err := CheckRole(c.ApiV1Controller, []string{"admin", strconv.Itoa(id)}); err != nil {
		return err
	}
	if err := CheckToken(c.ApiV1Controller); err != nil {
		return err
	}
	user := &models.User{}
	if err := db.DB.First(&user, id).Error; err != nil {
		return c.HandleNotFoundError(err.Error())
	}
	r := Response{ReseposeUser{user}}
	return c.RenderJSON(r)
}

// User Create
func (c ApiUser) Create() revel.Result {
	user := &models.User{}
	if err := c.BindParams(user); err != nil {
		return c.HandleBadRequestError(err.Error())
	}
	if err := validator.Validate(user); err != nil {
		return c.HandleBadRequestError(err.Error())
	}
	if user.Role != "" {
		s, b := revel.Config.String("admin.role")
		if !b {
			s = "AdminPower"
		}
		if user.Role == s {
			user.Role = "admin"
		} else {
			user.Role = "common"
		}
	}
	salt := []byte("yatuhashi")
	converted, _ := scrypt.Key([]byte(user.Password), salt, 16384, 8, 1, 32)
	user.Password = hex.EncodeToString(converted[:])
	if err := db.DB.Create(user).Error; err != nil {
		return c.HandleBadRequestError(err.Error())
	}
	r := Response{"Success"}
	return c.RenderJSON(r)
}

func (c ApiUser) Update(id int) revel.Result {
	if err := CheckRole(c.ApiV1Controller, []string{"admin", strconv.Itoa(id)}); err != nil {
		return err
	}
	if err := CheckToken(c.ApiV1Controller); err != nil {
		return err
	}
	userOld := &models.User{}
	if err := db.DB.First(&userOld, id).Error; err != nil {
		return c.HandleNotFoundError(err.Error())
	}
	userNew := &models.User{}
	if err := c.BindParams(userNew); err != nil {
		return c.HandleBadRequestError(err.Error())
	}
	if err := validator.Validate(userNew); err != nil {
		return c.HandleBadRequestError(err.Error())
	}
	salt := []byte("yatuhashi")
	converted, _ := scrypt.Key([]byte(userNew.Password), salt, 16384, 8, 1, 32)
	userNew.Password = hex.EncodeToString(converted[:])
	if err := CheckRole(c.ApiV1Controller, []string{"admin"}); err != nil {
		userNew.Role = userOld.Role
	}
	if err := db.DB.Model(&userOld).Update(&userNew).Error; err != nil {
		return c.HandleNotFoundError(err.Error())
	}
	r := Response{ReseposeUser{userNew}}
	return c.RenderJSON(r)
}

func (c ApiUser) Delete(id int) revel.Result {
	if err := CheckRole(c.ApiV1Controller, []string{"admin", strconv.Itoa(id)}); err != nil {
		return err
	}
	if err := CheckToken(c.ApiV1Controller); err != nil {
		return err
	}
	user := &models.User{}
	if err := db.DB.First(&user, id).Error; err != nil {
		return c.HandleNotFoundError(err.Error())
	}
	if err := db.DB.Delete(&user).Error; err != nil {
		return c.HandleInternalServerError("Record Delete Failure")
	}
	r := Response{"Success Delete"}
	return c.RenderJSON(r)
}
