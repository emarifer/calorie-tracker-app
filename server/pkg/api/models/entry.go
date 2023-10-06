package models

import (
	"time"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Entry struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt   time.Time          `json:"createdAt,omitempty" bson:"created_at,omitempty"`
	Dish        *string            `json:"dish,omitempty" validate:"required,min=1" bson:"dish,omitempty"`
	Fat         *float64           `json:"fat,omitempty" validate:"required,number" bson:"fat,omitempty"`
	Ingredients *string            `json:"ingredients,omitempty" validate:"required,min=1" bson:"ingredients,omitempty"`
	Calories    *string            `json:"calories,omitempty" validate:"required,min=1" bson:"calories,omitempty"`
}

type IngredientsUpdate struct {
	Ingredients *string `json:"ingredients" validate:"required,min=1"`
}

// Start Validator
var validate = validator.New()

type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

func ValidateStruct[T any](payload T) []*ErrorResponse {
	var errors []*ErrorResponse

	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}

	return errors
}

// End Validator

/* Ejemplos de uso del Validator. VER:
https://github.com/go-playground/validator/blob/master/_examples/simple/main.go
https://github.com/go-playground/validator/blob/master/_examples/struct-level/main.go
*/
