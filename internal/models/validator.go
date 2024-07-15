package models

import "github.com/go-playground/validator/v10"

type Item struct {
	Data string `json:"data" validate:"required"`
}

var validate = validator.New()

func (item *Item) Validate() error {
	return validate.Struct(item)
}
