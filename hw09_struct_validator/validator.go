package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	validationTag = regexp.MustCompile(`validate:".+"`)

	// Validation errors
	errLength       = errors.New("invalid length")
	errNotInSet     = errors.New("value not from set")
	errRegexp       = errors.New("value must match regular expression")
	errMinThreshold = errors.New("value is less than minimum")
	errMaxThreshold = errors.New("value is greater than maximum")

	// Program errors
	errWrongType              = errors.New("wrong type: expected struct")
	errUnknownValidationQuery = errors.New("unknown validation query")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	validationErrs := &strings.Builder{}
	for _, validationErr := range v {
		fieldErr := fmt.Sprintf("%s: %s\n", validationErr.Field, validationErr.Err.Error())
		validationErrs.WriteString(fieldErr)
	}
	return validationErrs.String()
}

func ParseValidationQuery(tag reflect.StructTag) string {
	valQs := validationTag.FindString(string(tag))
	valQsCutPrefix := strings.ReplaceAll(valQs, "validate:", "")
	validationQuery := strings.ReplaceAll(valQsCutPrefix, `"`, "")

	return validationQuery
}

func validateLength(condition, value string) error {
	validLenStr := strings.Split(condition, ":")[1]
	validLen, err := strconv.Atoi(validLenStr)
	if err != nil {
		return err
	}
	if len(value) != validLen {
		return fmt.Errorf("%s, must be %d, not %d", errLength, validLen, len(value))
	}

	return nil
}

func validateInSetString(condition, value string) error {
	validSetStr := strings.Split(condition, ":")[1]
	validSet := strings.Split(validSetStr, ",")
	valueInSet := false
	for _, s := range validSet {
		if value == s {
			valueInSet = true
			break
		}
	}
	if !valueInSet {
		return fmt.Errorf("%s, must be one of %v, not %s", errNotInSet, validSet, value)
	}

	return nil
}

func validateRegexp(condition, value string) error {
	reStr := strings.Split(condition, ":")[1]
	re, err := regexp.Compile(reStr)
	if err != nil {
		return err
	}
	if value != re.FindString(value) {
		return fmt.Errorf("%s: %s", errRegexp, reStr)
	}
	return nil
}

func validateInSetInt(condition string, value int64) error {
	validSetStr := strings.Split(condition, ":")[1]
	validSet := strings.Split(validSetStr, ",")
	valueInSet := false
	for _, s := range validSet {
		if strconv.Itoa(int(value)) == s {
			valueInSet = true
			break
		}
	}
	if !valueInSet {
		return fmt.Errorf("%s, must be one of %v, not %d", errNotInSet, validSet, value)
	}

	return nil
}

func validateMin(condition string, value int64) error {
	minThreshStr := strings.Split(condition, ":")[1]
	minThresh, err := strconv.Atoi(minThreshStr)
	if err != nil {
		return err
	}
	if value < int64(minThresh) {
		return fmt.Errorf("%s, should be at least %d", errMinThreshold, minThresh)
	}

	return nil
}

func validateMax(condition string, value int64) error {
	maxThreshStr := strings.Split(condition, ":")[1]
	maxThresh, err := strconv.Atoi(maxThreshStr)
	if err != nil {
		return err
	}
	if value > int64(maxThresh) {
		return fmt.Errorf("%s, should be at most %d", errMaxThreshold, maxThresh)
	}

	return nil
}

func validateFieldByQuery(fv reflect.Value, sv reflect.StructField, query string) error {
	conditions := strings.Split(query, "|")
	validationErrors := make(ValidationErrors, 0)
	for _, condition := range conditions {
		cond := strings.Split(condition, ":")[0]
		switch fv.Kind() {
		case reflect.String:
			v := fv.String()
			var validateFunc func(string, string) error
			var validateErr error
			switch cond {
			case "len":
				validateFunc = validateLength
				validateErr = errLength
			case "in":
				validateFunc = validateInSetString
				validateErr = errNotInSet
			case "regexp":
				validateFunc = validateRegexp
				validateErr = errRegexp
			default:
				return fmt.Errorf("%s: %s for field %s", errUnknownValidationQuery, query, sv.Name)
			}
			err := validateFunc(cond, v)
			if err != nil {
				if !errors.Is(err, validateErr) {
					return err
				}
				validationErrors = append(validationErrors, ValidationError{Field: sv.Name, Err: err})
			}
		case reflect.Int:
			v := fv.Int()
			var validateFunc func(string, int64) error
			var validateErr error
			switch cond {
			case "min":
				validateFunc = validateMin
				validateErr = errMinThreshold
			case "max":
				validateFunc = validateMax
				validateErr = errMaxThreshold
			case "in":
				validateFunc = validateInSetInt
				validateErr = errNotInSet
			default:
				return fmt.Errorf("%s: %s for field %s", errUnknownValidationQuery, query, sv.Name)
			}
			err := validateFunc(cond, v)
			if err != nil {
				if !errors.Is(err, validateErr) {
					return err
				}
				validationErrors = append(validationErrors, ValidationError{Field: sv.Name, Err: err})
			}
		}
	}

	return validationErrors
}

func Validate(v interface{}) error {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Struct {
		return fmt.Errorf("%s got %T", errWrongType, v)
	}

	t := value.Type()
	validationErrors := make(ValidationErrors, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		validationQuery := ParseValidationQuery(field.Tag)
		if validationQuery == "" {
			continue
		}
		err := validateFieldByQuery(value.Field(i), field, validationQuery)
		if err != nil {
			var fieldValidationErrs *ValidationErrors
			if !errors.As(err, &fieldValidationErrs) {
				return err
			}
			validationErrors = append(validationErrors, *fieldValidationErrs...)
		}
	}

	return validationErrors
}
