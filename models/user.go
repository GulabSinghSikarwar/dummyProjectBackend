package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id            primitive.ObjectID
	Firt_Name     *string `json:"first_name" validate:"required, min=2 ,max=100"`
	Last_Name     *string `json:"last_name" validate:"required,min=2,max=100"`
	Password      *string `json:"password" validate:"required", min=6"`
	Email         *string `json:"email"  validate:"email, required"`
	Phone         *string `json:"phone"  validate:"required" `
	Token         *string `json:"token"`
	User_type     *string `json:"user_type" validate:"required, eq=ADMIN|eq=USER"`
	Refresh_Token *string `json:"refresh_token"`
}
