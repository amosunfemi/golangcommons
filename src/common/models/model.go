package models

import "time"

//User ...
type User struct {
	BaseModel
	UserName          string `json:"user_name"`
	Salt              string `json:"salt"`
	Provider          string `json:"provider"`
	VerificationToken string `json:"verification_token"`
	Email             string `json:"email"`
	Password          string `json:"password"`
	PhoneNo           string `json:"phone_no"`
	Active            bool   `json:"active"`
	UUID              string
	UserGroup         string `json:"user_group"`
}

//TableName ...
func (user User) TableName() string {
	return "users"
}

//UserGroup ...
type UserGroup struct {
	BaseModel
	Status      string `json:"status"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// TableName ...
func (ind *UserGroup) TableName() string {
	return "user_group"
}

//BaseModel ...
type BaseModel struct {
	ID        int        `json:"id" gorm:"primary_key"`
	CreatedAt *time.Time `json:"created_at" sql:"timestamp with time zone DEFAULT ('now'::text)::date"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Status    string     `json:"status" sql:"DEFAULT:'ACTIVE'"`
	RecStatus string     `json:"rec_status" sql:"DEFAULT:'A'"`
}

// TokenAuthentication ...
type TokenAuthentication struct {
	Token string `json:"token" form:"token"`
}
