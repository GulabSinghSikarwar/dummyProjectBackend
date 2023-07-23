package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/go-playground/validator.v9"
	// "gopkg.in/go-playground/validator.v9"
)

// type User struct {
// 	ID       *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
// 	Name     string     `gorm:"type:varchar(100);not null"`
// 	Email    string     `gorm:"type:varchar(100);uniqueIndex;not null"`
// 	Password string     `gorm:"type:varchar(100);not null"`
// 	Role     *string    `gorm:"type:varchar(50);default:'user';not null"`
// 	Provider *string    `gorm:"type:varchar(50);default:'local';not null"`

// 	Verified  *bool      `gorm:"not null;default:false"`
// 	CreatedAt *time.Time `gorm:"not null;default:now()"`
// 	UpdatedAt *time.Time `gorm:"not null;default:now()"`
// }

type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	Name     string             `json:"name" validate:"required,min=2,max=100"`
	Password string             `json:"password" validate:"required", min=6"`
	Email    string             `json:"email"  validate:"email, required"`
}

//	type SignUpInput struct {
//		Name           string `json:"name" validate:"required"`
//		Email          string `json:"email" validate:"required"`
//		Password       string `json:"password" validate:"required min=8"`
//		ConfirmPasswod string `json:"confirmPassword" validate:"required min=8"`
//	}
type SignUpInput struct {
	Name           string `json:"name" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required,min=8"`
	ConfirmPasswod string `json:"confirmPassword" validate:"required,min=8"`
}

type SignInInput struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
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
