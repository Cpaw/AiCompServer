# Ai Competion Server

A high-productivity web framework for the [Go language](http://www.golang.org/).

### Start the web server:

   revel run 
 
### Build

   revel package AiCompServer

### This is Score Server System

```
.
├── README.md
├── app
│   ├── controllers
│   │   └── api
│   │       └── v1
│   │           ├── answer.go
│   │           ├── auth.go
│   │           ├── base.go
│   │           ├── challenge.go
│   │           ├── ranking.go
│   │           └── user.go
│   ├── db
│   │   └── gorm.go
│   ├── init.go
│   ├── models
│   │   ├── base.go
│   │   ├── challenge.go
│   │   └── user.go
│   ├── routes
│   │   └── routes.go
├── conf
│   ├── app.conf
│   └── routes
└── tests
    └── apptest.go
```



### API LIST:

```
// AiCompServer/app/controllers/api/v1/user.go
GET     /api/v1/user                            ApiUser.Index
GET     /api/v1/user/:id                        ApiUser.Show
PUT     /api/v1/user/:id                        ApiUser.Update
DELETE  /api/v1/user/:id                        ApiUser.Delete

POST    /api/v1/signup                          ApiUser.Create
// AiCompServer/app/controllers/api/v1/auth.go
GET     /api/v1/signin                          ApiAuth.GetSessionID
POST    /api/v1/signin                          ApiAuth.SignIn
DELETE  /api/v1/signin                          ApiAuth.SignOut
GET     /api/v1/role                            ApiAuth.Role

// AiCompServer/app/controllers/api/v1/ranking.go
GET    /api/v1/ranking                          ApiChallenge.Ranking

// AiCompServer/app/controllers/api/v1/challenge.go
GET    /api/v1/challenges                       ApiChallenge.Index
POST   /api/v1/challenges                       ApiChallenge.Create
GET    /api/v1/challenges/:id                   ApiChallenge.Show
PUT    /api/v1/challenges/:id                   ApiChallenge.Update
DELETE /api/v1/challenges/:id                   ApiChallenge.Delete

// AiCompServer/app/controllers/api/v1/answer.go
GET    /api/v1/answers                          ApiAnswer.Index
POST   /api/v1/answers                          ApiAnswer.Create
GET    /api/v1/answers/:id                      ApiAnswer.Show
PUT    /api/v1/answers/:id                      ApiAnswer.Update
DELETE /api/v1/answers/:id                      ApiAnswer.Delete

GET     /api/v1/challengeanswer/:id             ApiAnswer.UserChallengeAnswer

POST   /api/v1/submit                           ApiAnswer.Submit
```
