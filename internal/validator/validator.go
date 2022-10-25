package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
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

	for _, err := range err.(validator.ValidationErrors) {
		// Attempt to find field by name and get json tag name
		field, _ := reflected.Type().FieldByName(err.StructField())
		var name string

		//If json tag doesn't exist, use lower case of name
		if name = field.Tag.Get("json"); name == "" {
			name = strings.ToLower(err.StructField())
		}

		// Add error to map translating it
		errors[name] = append(errors[name], err.Translate(trans))
	}

	if len(errors) == 0 {
		return true, nil
	}

	return false, errors
}
