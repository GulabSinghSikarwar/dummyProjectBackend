package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/go-playground/validator.v9"
)

type Watchlist struct {
	UserId primitive.ObjectID   `bson:"userId" validate:"requried"`
	ID     primitive.ObjectID   `bson:"_id"`
	Stocks []primitive.ObjectID `json:"stocks",`
}

type AddStockRequestBody struct {
	StockId primitive.ObjectID `json:"stockId" validate:"requried"`
	// WatchlistId primitive.ObjectID `json:"watchlistId" validate:"requried"`
}
type WatchlistGetReqBody struct {
	WatchlistId primitive.ObjectID `json:"watchlistId"`
}

func ValidateAddStockRequestBody[T any](payload T) []*ErrorResponse {

	var errors []*ErrorResponse
	err := validate.Struct(payload)

	if err != nil {
		validationErrors, ok := err.(*validator.ValidationErrors)
		if !ok {
			for _, err := range *validationErrors {
				var element *ErrorResponse
				element.Field = err.StructField()
				element.Value = err.Param()
				element.Tag = err.Tag()
				errors = append(errors, element)

			}

		}

	}
	return errors
}
