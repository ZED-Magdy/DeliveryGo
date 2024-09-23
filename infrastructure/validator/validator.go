package validator

import (
	"github.com/go-playground/validator/v10"
	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/go-playground/locales/en"
)

type Service struct {
	Validator *validator.Validate
	Trans *ut.Translator
}

var validateService *Service

func New() Service {
	if(validateService != nil) {
		return *validateService
	}

	en := en.New()
	uni := ut.New(en, en)

	trans, _ := uni.GetTranslator("en")

	validate := validator.New()
	en_translations.RegisterDefaultTranslations(validate, trans)
	return Service{Validator: validate, Trans: &trans}
}