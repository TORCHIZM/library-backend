package helpers

import (
	"net/url"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	FailedField string `json:"Field"`
	Tag         string `json:"Tag"`
	Value       string `json:"Value"`
}

var validate *validator.Validate = validator.New()

func RegisterCustomValidations() {
	validate.RegisterValidation("regexp", func(fl validator.FieldLevel) bool {
		var regex = regexp.MustCompile(fl.Param())
		return regex.MatchString(fl.Field().String())
	})

	validate.RegisterValidation("string", func(fl validator.FieldLevel) bool {
		var regex = regexp.MustCompile("^[A-Za-z0-9._ıİöÖüÜşŞçÇğĞ]+$")
		return regex.MatchString(fl.Field().String())
	})

	validate.RegisterValidation("stringWithSpace", func(fl validator.FieldLevel) bool {
		var regex = regexp.MustCompile("^[A-Za-z0-9._ıİöÖüÜşŞçÇğĞ ]+$")
		return regex.MatchString(fl.Field().String())
	})

	validate.RegisterValidation("platform", func(fl validator.FieldLevel) bool {
		var regex = regexp.MustCompile("ios|android|web")
		return regex.MatchString(fl.Field().String())
	})

	validate.RegisterValidation("sid", func(fl validator.FieldLevel) bool {
		var regex = regexp.MustCompile("[A-Za-z0-9.-]+$")
		return regex.MatchString(fl.Field().String())
	})

	validate.RegisterValidation("ISO8601date", func(fl validator.FieldLevel) bool {
		ISO8601DateRegexString := "^(?:[1-9]\\d{3}-(?:(?:0[1-9]|1[0-2])-(?:0[1-9]|1\\d|2[0-8])|(?:0[13-9]|1[0-2])-(?:29|30)|(?:0[13578]|1[02])-31)|(?:[1-9]\\d(?:0[48]|[2468][048]|[13579][26])|(?:[2468][048]|[13579][26])00)-02-29)T(?:[01]\\d|2[0-3]):[0-5]\\d:[0-5]\\d(?:\\.\\d{1,9})?(?:Z|[+-][01]\\d:[0-5]\\d)$"
		ISO8601DateRegex := regexp.MustCompile(ISO8601DateRegexString)
		return ISO8601DateRegex.MatchString(fl.Field().String())
	})

	validate.RegisterValidation("commentType", func(fl validator.FieldLevel) bool {
		var regex = regexp.MustCompile("post|quotation")
		return regex.MatchString(fl.Field().String())
	})

	validate.RegisterValidation("userBookStatus", func(fl validator.FieldLevel) bool {
		var regex = regexp.MustCompile("reading|finished|overtime")
		return regex.MatchString(fl.Field().String())
	})

	validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		var number, upper, special bool = false, false, false

		for _, c := range fl.Field().String() {
			switch {
			case unicode.IsNumber(c):
				number = true
			case unicode.IsUpper(c):
				upper = true
			case unicode.IsPunct(c) || unicode.IsSymbol(c):
				special = true
			case unicode.IsSpace(c):
				return false
			}
		}

		return (number && upper && special)
	})

	validate.RegisterValidation("imagelink", func(fl validator.FieldLevel) bool {
		var prefix, suffix bool = false, false
		var value string = fl.Field().String()

		/**
		Uncomment and replace blow in production // TORCHIZM 13.04.2022 12:30
		var endpointPrefix string = os.Getenv("ENDPOINT") + "/api/storage"
		if strings.HasPrefix(value, "http://"+endpointPrefix) || strings.HasPrefix(value, "https://"+endpointPrefix) {
		*/

		if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
			prefix = true
		}

		if strings.HasSuffix(value, ".png") || strings.HasSuffix(value, ".jpg") || strings.HasSuffix(value, ".jpeg") {
			suffix = true
		}

		uri, err := url.Parse(value)

		if err != nil && uri.Scheme != "http" && uri.Scheme != "https" {
			return false
		}

		return (prefix && suffix)
	})
}

func ValidateStruct(object interface{}) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(object)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			element := &ErrorResponse{
				FailedField: err.StructNamespace(),
				Tag:         err.Tag(),
				Value:       err.Param(),
			}

			errors = append(errors, element)
		}
	}

	return errors
}
