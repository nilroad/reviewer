package rule

import (
	"club/internal/api/rest/request"
	"club/internal/config"
	"club/pkg/bank"
	"reflect"
	"regexp"
	"strings"
	"time"

	"git.oceantim.com/backend/packages/golang/essential/validation"
)

const (
	MinimumAge = 13
)

func getRules(ruleConfig config.RuleConfig) map[string]validation.ValidateFunc {
	return map[string]validation.ValidateFunc{
		"ir_phone":             isIrPhone,
		"username":             usernameValidator(ruleConfig),
		"require_if_not_empty": requireIfNotExists,
		"postal_code":          postalCodeValidator,
		"birth_date":           birthDateValidator,
		"card_number":          cardNumberValidator,
		"not_expired_date":     notExpiredDateValidator,
	}
}

func isIrPhone(vl validation.FieldLevel) bool {
	phoneNumber := vl.Field().String()
	rgx := regexp.MustCompile(`[+](989)\d{9}$`)

	return rgx.MatchString(phoneNumber)
}

func RegisterRules(ruleConfig config.RuleConfig) error {
	rules := getRules(ruleConfig)

	for r, f := range rules {
		if err := validation.RegisterNewRule(r, f); err != nil {
			return err
		}
	}

	return nil
}

func usernameValidator(ruleConfig config.RuleConfig) validation.ValidateFunc {
	return func(vl validation.FieldLevel) bool {

		value := vl.Field().String()
		length := len(value)

		if length < ruleConfig.UserNameMinLength || length > ruleConfig.UserNameMaxLength {
			return false
		}

		if strings.HasPrefix(value, ".") || strings.HasSuffix(value, ".") {
			return false
		}

		if strings.Contains(value, "..") {
			return false
		}

		rgx := regexp.MustCompile(`^\d*$`)
		if matched := rgx.MatchString(value); matched {
			return false
		}

		rgx = regexp.MustCompile(`^[A-Za-z._0-9]*$`)
		if matched := rgx.MatchString(value); !matched {
			return false
		}

		return true
	}
}

func requireIfNotExists(fl validation.FieldLevel) bool {
	// Get the name of the other field passed as a parameter
	otherFieldName := fl.Param()
	if otherFieldName == "" {
		return false // no field name provided
	}

	// Access the struct this field belongs to
	parent := fl.Parent()

	// Get the value of the "other" field
	otherField := parent.FieldByName(otherFieldName)
	if !otherField.IsValid() {
		return false // invalid field name
	}

	// Get the value of the current field (the one being validated)
	currentField := fl.Field()

	// If the other field is empty and current field is also empty, return false
	if isEmptyValue(otherField) && isEmptyValue(currentField) {
		return false
	}

	return true
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return strings.TrimSpace(v.String()) == ""
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Slice, reflect.Map, reflect.Array:
		return v.Len() == 0
	default:
		zero := reflect.Zero(v.Type()).Interface()

		return reflect.DeepEqual(v.Interface(), zero)
	}
}

func postalCodeValidator(vl validation.FieldLevel) bool {
	value := vl.Field().String()

	rgx := regexp.MustCompile(`^[1-9]\d{9}$`)
	if matched := rgx.MatchString(value); !matched {
		return false
	}

	return true
}

func birthDateValidator(vl validation.FieldLevel) bool {
	field := vl.Field()

	dateOnlyPtr, ok := field.Interface().(request.DateOnly)
	if !ok {
		return false
	}

	birthDate := time.Time(dateOnlyPtr)

	return checkAge(birthDate)
}

func checkAge(birthDate time.Time) bool {
	now := time.Now()
	age := now.Year() - birthDate.Year()

	if now.Month() < birthDate.Month() || (now.Month() == birthDate.Month() && now.Day() < birthDate.Day()) {
		age--
	}

	return age >= MinimumAge
}

func cardNumberValidator(vl validation.FieldLevel) bool {
	value := vl.Field().String()

	return bank.ValidateCard(value)
}

func notExpiredDateValidator(vl validation.FieldLevel) bool {
	field := vl.Field()

	// Support string (RFC3339/RFC3339Nano) and pointers to string, as well as time.Time
	switch field.Kind() {
	case reflect.String:
		s := strings.TrimSpace(field.String())
		if s == "" {
			return false
		}
		if t, ok := parseRFC3339Flexible(s); ok {
			return isDateNotExpired(t)
		}

		return false
	case reflect.Ptr:
		if field.IsNil() {
			return false
		}
		elem := field.Elem()
		if elem.Kind() == reflect.String {
			s := strings.TrimSpace(elem.String())
			if s == "" {
				return false
			}
			if t, ok := parseRFC3339Flexible(s); ok {
				return isDateNotExpired(t)
			}

			return false
		}
		if elem.Type() == reflect.TypeOf(time.Time{}) {
			t := elem.Interface().(time.Time) //nolint:all

			return isDateNotExpired(t)
		}
	default:
		if field.Type() == reflect.TypeOf(time.Time{}) {
			t := field.Interface().(time.Time) //nolint:all

			return isDateNotExpired(t)
		}
	}

	return false
}

func isDateNotExpired(date time.Time) bool {
	return !date.Before(time.Now())
}

func parseRFC3339Flexible(s string) (time.Time, bool) {
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return t, true
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, true
	}

	return time.Time{}, false
}
