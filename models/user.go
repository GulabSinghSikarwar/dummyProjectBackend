package models

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/go-playground/validator.v9"
	// "gopkg.in/go-playground/validator.v9"
)

type User struct {
	ID       *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name     string     `gorm:"type:varchar(100);not null"`
	Email    string     `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password string     `gorm:"type:varchar(100);not null"`
	Role     *string    `gorm:"type:varchar(50);default:'user';not null"`
	Provider *string    `gorm:"type:varchar(50);default:'local';not null"`
	// Photo     *string    `gorm:"not null;default:'default.png'"`
	Verified  *bool      `gorm:"not null;default:false"`
	CreatedAt *time.Time `gorm:"not null;default:now()"`
	UpdatedAt *time.Time `gorm:"not null;default:now()"`
}

//	type User struct {
//		Id            primitive.ObjectID
//		Firt_Name     *string `json:"first_name" validate:"required, min=2 ,max=100"`
//		Last_Name     *string `json:"last_name" validate:"required,min=2,max=100"`
//		Password      *string `json:"password" validate:"required", min=6"`
//		Email         *string `json:"email"  validate:"email, required"`
//		Phone         *string `json:"phone"  validate:"required" `
//		Token         *string `json:"token"`
//		User_type     *string `json:"user_type" validate:"required, eq=ADMIN|eq=USER"`
//		Refresh_Token *string `json:"refresh_token"`
//	}
type SignUpInput struct {
	Name           string `json:"name" validate:"required"`
	Email          string `json:"name" validate:"required"`
	Password       string `json:"name" validate:"required mi=8"`
	ConfirmPasswod string `json:"name" validate:"required min=8"`
}
type SignIpInput struct {
	Email    string `json:"name" validate:"required"`
	Password string `json:"name" validate:"required mi=8"`
}
type ErrorResponse struct {
	Tag   string `json:tag`
	Field string `json:field`
	Value string `json :"value ,omitempty"`
}

var validate = validator.New()

func ValidateStruct[T any](payload T) []*ErrorResponse {

	var error []*ErrorResponse

	err := validate.Struct(payload)

	if err != nil {
		validationErrors, ok := err.(*validator.ValidationErrors)
		if !ok {
			for _, err := range *validationErrors {

				var element ErrorResponse
				element.Field = err.StructNamespace()
				element.Tag = err.Tag()
				element.Value = err.Param()
				error = append(error, &element)

			}
		}
	}
	return error

}
