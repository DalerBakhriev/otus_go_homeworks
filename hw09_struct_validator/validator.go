package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const validationTagKey = "validate"

var (
	// Validation errors.
	errLength       = errors.New("invalid length")
	errNotInSet     = errors.New("value not from set")
	errRegexp       = errors.New("value must match regular expression")
	errMinThreshold = errors.New("value is less than minimum")
	errMaxThreshold = errors.New("value is greater than maximum")

	// Program errors.
	errWrongType              = errors.New("wrong type: expected struct")
	errUnknownValidationQuery = errors.New("unknown validation query")
	errFieldType              = errors.New("invalid field type for validation")
)

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: %s\n", v.Field, v.Err)
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	validationErrs := &strings.Builder{}
	for _, validationErr := range v {
		validationErrs.WriteString(validationErr.Error())
	}
	return validationErrs.String()
}

func ParseValidationQuery(tag reflect.StructTag) string {
	valQs := tag.Get(validationTagKey)
	valQsCutPrefix := strings.ReplaceAll(valQs, "validate:", "")
	validationQuery := strings.ReplaceAll(valQsCutPrefix, `"`, "")

	return validationQuery
}

func validateLength(condition, value string) error {
	validLen, err := strconv.Atoi(condition)
	if err != nil {
		return err
	}
	if len(value) != validLen {
		return fmt.Errorf("%w, must be %d, not %d", errLength, validLen, len(value))
	}
	return nil
}

func validateInSetString(condition, value string) error {
	validSet := strings.Split(condition, ",")
	valueInSet := false
	for _, s := range validSet {
		if value == s {
			valueInSet = true
			break
		}
	}
	if !valueInSet {
		return fmt.Errorf("%w, must be one of %v, not %s", errNotInSet, validSet, value)
	}

	return nil
}

func validateRegexp(condition, value string) error {
	re, err := regexp.Compile(condition)
	if err != nil {
		return err
	}
	if value != re.FindString(value) {
		return fmt.Errorf("%w: %s", errRegexp, condition)
	}
	return nil
}

func validateInSetInt(condition string, value int64) error {
	validSet := strings.Split(condition, ",")
	valueInSet := false
	for _, s := range validSet {
		if strconv.Itoa(int(value)) == s {
			valueInSet = true
			break
		}
	}
	if !valueInSet {
		return fmt.Errorf("%w, must be one of %v, not %d", errNotInSet, validSet, value)
	}

	return nil
}

func validateMin(condition string, value int64) error {
	minThresh, err := strconv.Atoi(condition)
	if err != nil {
		return err
	}
	if value < int64(minThresh) {
		return fmt.Errorf("%w, should be at least %d", errMinThreshold, minThresh)
	}

	return nil
}

func validateMax(condition string, value int64) error {
	maxThresh, err := strconv.Atoi(condition)
	if err != nil {
		return err
	}
	if value > int64(maxThresh) {
		return fmt.Errorf("%w, should be at most %d", errMaxThreshold, maxThresh)
	}

	return nil
}

func validateStringField(v string, fn string, query string) error {
	conditions := strings.Split(query, "|")
	validationErrors := make(ValidationErrors, 0)
	for _, condition := range conditions {
		condKeyWord := strings.Split(condition, ":")[0]
		cond := strings.Split(condition, ":")[1]
		var validateFunc func(string, string) error
		var validateErr error
		switch condKeyWord {
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
			return fmt.Errorf("%w: %s for field %s", errUnknownValidationQuery, query, fn)
		}
		err := validateFunc(cond, v)
		if err != nil {
			if !errors.Is(err, validateErr) {
				return err
			}
			validationErrors = append(validationErrors, ValidationError{Field: fn, Err: err})
		}
	}
	return validationErrors
}

func validateStringSliceField(v []string, fn string, query string) error {
	fmt.Printf("input is %v\n", v)
	for _, el := range v {
		err := validateStringField(el, fn, query)
		if err != nil {
			var valErrs ValidationErrors
			if errors.As(err, &valErrs) && len(valErrs) == 0 {
				continue
			}
			return err
		}
	}
	return nil
}

func validateIntField(v int64, fn string, query string) error {
	conditions := strings.Split(query, "|")
	validationErrors := make(ValidationErrors, 0)
	for _, condition := range conditions {
		condKeyWord := strings.Split(condition, ":")[0]
		cond := strings.Split(condition, ":")[1]
		var validateFunc func(string, int64) error
		var validateErr error
		switch condKeyWord {
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
			return fmt.Errorf("%w: %s for field %s", errUnknownValidationQuery, query, fn)
		}
		err := validateFunc(cond, v)
		if err != nil {
			if !errors.Is(err, validateErr) {
				return err
			}
			validationErrors = append(validationErrors, ValidationError{Field: fn, Err: err})
		}
	}
	return validationErrors
}

func validateIntSliceField(v []int, fn string, query string) error {
	for _, el := range v {
		err := validateIntField(int64(el), fn, query)
		if err != nil {
			var valErrs ValidationErrors
			if errors.As(err, &valErrs) && len(valErrs) == 0 {
				continue
			}
			return err
		}
	}
	return nil
}

func Validate(v interface{}) error {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Struct {
		return fmt.Errorf("%w got %T", errWrongType, v)
	}

	t := value.Type()
	validationErrors := make(ValidationErrors, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		validationQuery := ParseValidationQuery(field.Tag)
		if validationQuery == "" {
			continue
		}
		fieldVal := value.Field(i)
		var err error
		switch fieldVal.Kind() {
		case reflect.String:
			err = validateStringField(fieldVal.String(), field.Name, validationQuery)
		case reflect.Int:
			err = validateIntField(fieldVal.Int(), field.Name, validationQuery)
		case reflect.Slice:
			switch fieldVal.Type().Elem().Kind() {
			case reflect.String:
				arr, ok := fieldVal.Interface().([]string)
				if !ok {
					return fmt.Errorf("%w, field %s", errFieldType, field.Name)
				}
				err = validateStringSliceField(arr, field.Name, validationQuery)
			case reflect.Int:
				arr, ok := fieldVal.Interface().([]int)
				if !ok {
					return fmt.Errorf("%w, field %s", errFieldType, field.Name)
				}
				err = validateIntSliceField(arr, field.Name, validationQuery)
			default:
				return fmt.Errorf("%w, field %s", errFieldType, field.Name)
			}
		default:
			return fmt.Errorf("%w, field %s", errFieldType, field.Name)
		}
		if err != nil {
			var fieldValidationErrs ValidationErrors
			if !errors.As(err, &fieldValidationErrs) {
				return err
			}
			validationErrors = append(validationErrors, fieldValidationErrs...)
		}
	}

	return validationErrors
}
