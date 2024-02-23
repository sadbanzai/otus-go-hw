package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type CustomError struct {
	Message string
}

func (e CustomError) Error() string {
	return e.Message
}

var (
	ErrNotStruct = CustomError{"value is not struct"}
	ErrBadParam  = CustomError{"bad parameter"}
	ErrNotImpl   = CustomError{"not implemented"}
)

type ValidateError struct {
	Message string
}

func (e ValidateError) Error() string {
	return e.Message
}

var (
	ErrMin    = ValidateError{"less than min"}
	ErrMax    = ValidateError{"more than max"}
	ErrIn     = ValidateError{"not in set"}
	ErrLen    = ValidateError{"incorrect length"}
	ErrRegexp = ValidateError{"does not match regexp"}
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var str strings.Builder
	for _, valErr := range v {
		str.WriteString(fmt.Sprintf("%s: %s, ", valErr.Field, valErr.Err.Error()))
	}
	return strings.TrimSuffix(str.String(), ", ")
}

type ValidationFunc func(v interface{}, param string) error

type Validator struct {
	validationFuncs map[string]ValidationFunc
}

func NewValidator() *Validator {
	return &Validator{
		validationFuncs: map[string]ValidationFunc{
			"min":    Min,
			"max":    Max,
			"in":     In,
			"len":    Len,
			"regexp": Regexp,
		},
	}
}

func toInt(param string) (int64, error) {
	i, err := strconv.ParseInt(param, 0, 64)
	if err != nil {
		return 0, ErrBadParam
	}
	return i, nil
}

func Min(v interface{}, param string) error {
	p, err := toInt(param)
	if err != nil {
		return err
	}
	invalid := false
	st := reflect.ValueOf(v)
	switch st.Kind() {
	case reflect.Int:
		invalid = st.Int() < p
	default:
		return ErrNotImpl
	}
	if invalid {
		return ErrMin
	}
	return nil
}

func Max(v interface{}, param string) error {
	p, err := toInt(param)
	if err != nil {
		return err
	}
	invalid := false
	st := reflect.ValueOf(v)
	switch st.Kind() {
	case reflect.Int:
		invalid = st.Int() > p
	default:
		return ErrNotImpl
	}
	if invalid {
		return ErrMax
	}
	return nil
}

func In(v interface{}, param string) error {
	ins := strings.Split(param, ",")
	if len(ins) == 0 {
		return ErrBadParam
	}
	invalid := false
	st := reflect.ValueOf(v)
	switch st.Kind() {
	case reflect.Int:
		var insInt []int64
		for _, inVal := range ins {
			inInt, err := toInt(inVal)
			if err != nil {
				return err
			}
			insInt = append(insInt, inInt)
		}
		invalid = !slices.Contains(insInt, st.Int())
	case reflect.String:
		invalid = !slices.Contains(ins, st.String())
	default:
		return ErrNotImpl
	}
	if invalid {
		return ErrIn
	}
	return nil
}

func Len(v interface{}, param string) error {
	p, err := toInt(param)
	if err != nil {
		return err
	}
	invalid := false
	st := reflect.ValueOf(v)
	switch st.Kind() {
	case reflect.String:
		invalid = int64(len(st.String())) != p
	default:
		return ErrNotImpl
	}
	if invalid {
		return ErrLen
	}
	return nil
}

func Regexp(v interface{}, param string) error {
	re, err := regexp.Compile(param)
	if err != nil {
		return ErrBadParam
	}
	invalid := false
	st := reflect.ValueOf(v)
	switch st.Kind() {
	case reflect.String:
		invalid = !re.MatchString(st.String())
	default:
		return ErrNotImpl
	}
	if invalid {
		return ErrRegexp
	}
	return nil
}

func Validate(v interface{}) error {
	var validationErrors ValidationErrors
	validator := NewValidator()
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Struct {
		return ErrNotStruct
	}
	st := reflect.TypeOf(v)
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		fieldValue := value.Field(i)
		validate, ok := field.Tag.Lookup("validate")
		if !ok || validate == "" {
			continue
		}
		rules := strings.Split(validate, "|")
		for _, ruleStr := range rules {
			rule := strings.Split(ruleStr, ":")
			ruleName := rule[0]
			ruleParam := rule[1]
			ruleFunc, ruleExists := validator.validationFuncs[ruleName]
			if ruleExists {
				switch field.Type.Kind() {
				case reflect.Int, reflect.String:
					valError := ruleFunc(fieldValue.Interface(), ruleParam)
					if valError != nil {
						validationErrors = append(
							validationErrors,
							ValidationError{Field: field.Name, Err: valError})
					}
				case reflect.Slice:
					for i := 0; i < fieldValue.Len(); i++ {
						valError := ruleFunc(fieldValue.Index(i).Interface(), ruleParam)
						if valError != nil {
							sliceElem := fmt.Sprintf("%s %d", field.Name, i)
							validationErrors = append(
								validationErrors,
								ValidationError{Field: sliceElem, Err: valError})
						}
					}
				}
			}
		}

	}
	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}
