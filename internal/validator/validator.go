package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	ptBRTranslation "github.com/go-playground/validator/v10/translations/pt_BR"
)

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
)

func Validate(data any) (bool, map[string][]string) {
	ptBR := pt_BR.New()
	uni = ut.New(ptBR, ptBR)

	trans, _ := uni.GetTranslator("pt_BR")
	ptBRTranslation.RegisterDefaultTranslations(validate, trans)

	err := validate.Struct(data)

	errors := make(map[string][]string)
	reflected := reflect.ValueOf(data)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			field, _ := reflected.Type().FieldByName(err.StructField())
			var name string

			if name = field.Tag.Get("json"); name == "" {
				name = strings.ToLower(err.StructField())
			}

			errors[name] = append(errors[name], err.Translate(trans))
		}
	}

	if len(errors) == 0 {
		return true, nil
	}

	return false, errors
}

func RegisterValidators() error {
	validate = validator.New()
	err := validate.RegisterValidation("notblank", validators.NotBlank)
	return err
}
