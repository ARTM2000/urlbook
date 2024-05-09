package pkg

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	MessageRequired              = "%s is required"
	MessageLen                   = "%s must be exactly %s characters long"
	MessageMin                   = "%s must be at least %s characters long"
	MessageMax                   = "%s must be at most %s characters long"
	MessageEqual                 = "%s must be equal to %s"
	MessageNotEqual              = "%s must not be equal to %s"
	MessageOneOf                 = "%s must be one of %s"
	MessageLessThan              = "%s must be less than %s"
	MessageLessThanEqual         = "%s must be less than or equal to %s"
	MessageGreaterThan           = "%s must be greater than %s"
	MessageGreaterThanEqual      = "%s must be greater than or equal to %s"
	MessageEqualField            = "%s must be equal to %s"
	MessageNotEqualField         = "%s must not be equal to %s"
	MessageGreaterThanField      = "%s must be greater than %s"
	MessageGreaterThanEqualField = "%s must be greater than or equal to %s"
	MessageLessThanField         = "%s must be less than %s"
	MessageLessThanEqualField    = "%s must be less than or equal to %s"
	MessageAlpha                 = "%s must contain only letters"
	MessageAlphaNum              = "%s must contain only letters and numbers"
	MessageAlphaUnicode          = "%s must contain only Unicode letters"
	MessageAlphaNumUnicode       = "%s must contain only Unicode letters and numbers"
	MessageNumeric               = "%s must be a numeric value"
	MessageNumber                = "%s must be a number"
	MessageHexadecimal           = "%s must be a hexadecimal number"
	MessageEmail                 = "%s must be a valid email address"
	MessageURL                   = "%s must be a valid URL"
	MessageURI                   = "%s must be a valid URI"
	MessageBase64                = "%s must be a valid base64-encoded string"
	MessageContains              = "%s must contain %s"
	MessageContainsAny           = "%s must contain at least one of the following characters: %s"
	MessageExcludes              = "%s may not contain %s"
	MessageExcludesAll           = "%s may not contain any of the following characters: %s"
	MessageExcludesRune          = "%s may not contain the following character: %s"
	MessageUUID                  = "%s must be a valid UUID"
	MessageUUID3                 = "%s must be avalid UUIDv3"
	MessageUUID4                 = "%s must be a valid UUIDv4"
	MessageUUID5                 = "%s must be a valid UUIDv5"
	MessageDataURI               = "%s must be a valid data URI"
	MessageIPv4                  = "%s must be a valid IPv4 address"
	MessageIP                    = "%s must be a valid IP address"
	MessageBoolean               = "%s must ba a valid Boolean"
	MessageGroupInvalid          = "%s is invalid. must be %s"
	MessageCron                  = "%s must be valid cron schedule"
	MessagePassword              = "%s is not strong enough"
	MessageDirPath               = "%s is not valid directory path"
	MessageDir                   = "%s is not a existing directory"
	MessageFilename              = "%s is not valid filename"
	MessageHTTPUrl               = "%s is not a valid http url"
)

func defaultGetValidatorErrorMessage(tag string, field string, param string, info ...string) (string, bool) {
	switch tag {
	case "required":
		return fmt.Sprintf(MessageRequired, field), true
	case "omitempty":
		return "", true
	case "len":
		return fmt.Sprintf(MessageLen, field, param), true
	case "min":
		return fmt.Sprintf(MessageMin, field, param), true
	case "max":
		return fmt.Sprintf(MessageMax, field, param), true
	case "eq":
		return fmt.Sprintf(MessageEqual, field, param), true
	case "ne":
		return fmt.Sprintf(MessageNotEqual, field, param), true
	case "oneof":
		return fmt.Sprintf(MessageOneOf, field, param), true
	case "lt":
		return fmt.Sprintf(MessageLessThan, field, param), true
	case "lte":
		return fmt.Sprintf(MessageLessThanEqual, field, param), true
	case "gt":
		return fmt.Sprintf(MessageGreaterThan, field, param), true
	case "gte":
		return fmt.Sprintf(MessageGreaterThanEqual, field, param), true
	case "eqfield":
		return fmt.Sprintf(MessageEqualField, field, param), true
	case "nefield":
		return fmt.Sprintf(MessageNotEqualField, field, param), true
	case "gtfield":
		return fmt.Sprintf(MessageGreaterThanField, field, param), true
	case "gtefield":
		return fmt.Sprintf(MessageGreaterThanEqualField, field, param), true
	case "ltfield":
		return fmt.Sprintf(MessageLessThanField, field, param), true
	case "ltefield":
		return fmt.Sprintf(MessageLessThanEqualField, field, param), true
	case "alpha":
		return fmt.Sprintf(MessageAlpha, field), true
	case "alphanum":
		return fmt.Sprintf(MessageAlphaNum, field), true
	case "alphaunicode":
		return fmt.Sprintf(MessageAlphaUnicode, field), true
	case "alphanumunicode":
		return fmt.Sprintf(MessageAlphaNumUnicode, field), true
	case "numeric":
		return fmt.Sprintf(MessageNumeric, field), true
	case "number":
		return fmt.Sprintf(MessageNumber, field), true
	case "hexadecimal":
		return fmt.Sprintf(MessageHexadecimal, field), true
	case "email":
		return fmt.Sprintf(MessageEmail, field), true
	case "url":
		return fmt.Sprintf(MessageURL, field), true
	case "http_url":
		return fmt.Sprintf(MessageHTTPUrl, field), true
	case "uri":
		return fmt.Sprintf(MessageURI, field), true
	case "base64":
		return fmt.Sprintf(MessageBase64, field), true
	case "contains":
		return fmt.Sprintf(MessageContains, field, param), true
	case "containsany":
		return fmt.Sprintf(MessageContainsAny, field, param), true
	case "excludes":
		return fmt.Sprintf(MessageExcludes, field, param), true
	case "excludesall":
		return fmt.Sprintf(MessageExcludesAll, field, param), true
	case "excludesrune":
		return fmt.Sprintf(MessageExcludesRune, field, param), true
	case "uuid":
		return fmt.Sprintf(MessageUUID, field), true
	case "uuid3":
		return fmt.Sprintf(MessageUUID3, field), true
	case "uuid4":
		return fmt.Sprintf(MessageUUID4, field), true
	case "uuid5":
		return fmt.Sprintf(MessageUUID5, field), true
	case "datauri":
		return fmt.Sprintf(MessageDataURI, field), true
	case "ipv4":
		return fmt.Sprintf(MessageIPv4, field), true
	case "ip":
		return fmt.Sprintf(MessageIP, field), true
	case "boolean":
		return fmt.Sprintf(MessageBoolean, field), true
	case "cron":
		return fmt.Sprintf(MessageCron, field), true
	case "dirpath":
		return fmt.Sprintf(MessageDirPath, field), true
	case "dir":
		return fmt.Sprintf(MessageDir, field), true

	// custom tag
	case "groupinvalid":
		return fmt.Sprintf(MessageGroupInvalid, field, info[0]), true
	case "password":
		return fmt.Sprintf(MessagePassword, field), true
	case "filename":
		return fmt.Sprintf(MessageFilename, field), true

	default:
		fmt.Printf("not defined tag for validation in messages: %s\n", tag)
		return fmt.Sprintf("Validation error for %s", field), false
	}
}

type validationError struct {
	FailedField string `json:"field"`
	Message     string `json:"message"`
}

type CustomMessageMapper = func(tag string, field string, param string, info ...string) (string, bool)

type superValidator struct {
	validate      *validator.Validate
	messageMapper CustomMessageMapper
}

type SuperValidator interface {
	RegisterMessageMapper(cmm CustomMessageMapper)
	RegisterCustomValidators(tag string, fn validator.Func)
	ValidateStruct(s interface{}) (errors []*validationError, ok bool)
	ValidateSliceParamUniqueness(s []any) (bool, any)
}

func NewSuperValidator() SuperValidator {
	var v = validator.New()
	sv := &superValidator{
		validate: v,
		messageMapper: func(tag, field, param string, info ...string) (string, bool) {
			return "", false
		},
	}

	// add some default custom validators
	sv.RegisterCustomValidators("password", validatePassword)
	sv.RegisterCustomValidators("filename", validateFilename)

	return sv
}

func (sv *superValidator) RegisterMessageMapper(cmm CustomMessageMapper) {
	sv.messageMapper = cmm
}

func (sv *superValidator) RegisterCustomValidators(tag string, fn validator.Func) {
	sv.validate.RegisterValidation(tag, fn)
}

func (sv *superValidator) ValidateStruct(s interface{}) (errors []*validationError, ok bool) {
	err := sv.validate.Struct(s)
	ok = true
	if err != nil {
		ok = false
		for _, err := range err.(validator.ValidationErrors) {
			var e validationError
			e.FailedField = err.Field()
			if strings.Contains(err.Tag(), "|") {
				tags := strings.Split(err.Tag(), "|")
				last := len(tags) - 1
				infoMsg := strings.Join(tags[:last], ",")
				if len(tags) > 1 {
					infoMsg = strings.Join(tags[:last], ",") + fmt.Sprintf(" or %s", tags[last])
				}
				cMsg, isHandled := sv.messageMapper(
					"groupinvalid",
					err.Field(),
					err.Param(),
					infoMsg,
				)
				if isHandled {
					e.Message = cMsg
					continue
				}

				msg, _ := defaultGetValidatorErrorMessage(
					"groupinvalid",
					err.Field(),
					err.Param(),
					infoMsg,
				)
				e.Message = msg
			} else {
				cMsg, isHandled := sv.messageMapper(err.Tag(), err.Field(), err.Param())
				if isHandled {
					e.Message = cMsg
					continue
				}

				msg, _ := defaultGetValidatorErrorMessage(err.Tag(), err.Field(), err.Param())
				e.Message = msg
			}
			errors = append(errors, &e)
		}
	}
	return
}

func (sv *superValidator) ValidateSliceParamUniqueness(s []any) (bool, any) {
	sOccurrence := map[any]bool{}

	for _, cn := range s {
		if sOccurrence[cn] {
			return false, &cn
		} else {
			sOccurrence[cn] = true
		}
	}

	return true, nil
}
