package services

import (
	"encoding/json"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func Validate(validatorInstance *validator.Validate, trans ut.Translator, request interface{}) []byte {
	err := validatorInstance.Struct(request)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			
			return []byte(`{"error": "Error validating request"}`)
		}
		errs := err.(validator.ValidationErrors)
		errors := errs.Translate(trans)

		res, err := json.Marshal(errors)
		if err != nil {
			return []byte(`{"error": "Error marshalling errors"}`)
		}
		return res

	}
	return nil
}