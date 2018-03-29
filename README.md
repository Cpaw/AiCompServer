# Ai Competion Server

Thie System uses revel
A high-productivity web framework for the [Go language](http://www.golang.org/).

### Start the server:

```
git clone https://github.com/Cpaw/AiComp.git
cd AiComp
docker-compose up -d
```

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
METHOD  URL                                     Function                      Allow Role
// AiCompServer/app/controllers/api/v1/user.go
GET     /api/v1/user                            ApiUser.Index                 Admin
GET     /api/v1/user/:id                        ApiUser.Show                  Admin, Signed in User
PUT     /api/v1/user/:id                        ApiUser.Update                Admin, Signed in User
DELETE  /api/v1/user/:id                        ApiUser.Delete                Admin, Signed in User

POST    /api/v1/signup                          ApiUser.Create                Everyone
// AiCompServer/app/controllers/api/v1/auth.go
GET     /api/v1/signin                          ApiAuth.GetSessionID          Everyone
POST    /api/v1/signin                          ApiAuth.SignIn                Everyone
DELETE  /api/v1/signin                          ApiAuth.SignOut               Signed in
GET     /api/v1/role                            ApiAuth.Role                  Everyone

// AiCompServer/app/controllers/api/v1/ranking.go
GET    /api/v1/ranking                          ApiChallenge.Ranking          Signed in

// AiCompServer/app/controllers/api/v1/challenge.go
GET    /api/v1/challenges                       ApiChallenge.Index            Signed in
POST   /api/v1/challenges                       ApiChallenge.Create           Admin
GET    /api/v1/challenges/:id                   ApiChallenge.Show             Signed in
PUT    /api/v1/challenges/:id                   ApiChallenge.Update           Admin
DELETE /api/v1/challenges/:id                   ApiChallenge.Delete           Admin

// AiCompServer/app/controllers/api/v1/answer.go
GET    /api/v1/answers                          ApiAnswer.Index               Admin
POST   /api/v1/answers                          ApiAnswer.Create              Admin
GET    /api/v1/answers/:id                      ApiAnswer.Show                Admin
PUT    /api/v1/answers/:id                      ApiAnswer.Update              Admin
DELETE /api/v1/answers/:id                      ApiAnswer.Delete              Admin

GET     /api/v1/challengeanswer/:id             ApiAnswer.UserChallengeAnswer Signed in User

POST   /api/v1/submit                           ApiAnswer.Submit              Signed in
```
