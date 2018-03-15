package models

type User struct {
	BaseModel
	Username string `sql:"size:64" json:"username" validate:"min=1, max=64" gorm:"not null;unique"`
	Password string `sql:"size:64" json:"password" validate:"min=6, max=64" gorm:"not null"`
	Role     string `sql:"size:32" json:"role" validate:"min=0, max=32"`
	Token    string `sql:"size:64" json:"token" validate:"max=64"`
}
