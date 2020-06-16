package directive

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-playground/validator/v10"
)

// ValidateDirective .
func ValidateDirective(ctx context.Context, obj interface{}, next graphql.Resolver, field string, rules string) (res interface{}, err error) {
	if field == "" || rules == "" {
		return nil, errors.New("Missing validate field or rules")
	}

	value := obj.(map[string]interface{})[field]
	if err := validate().Var(value, rules); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return nil, fmt.Errorf(`Validation failed on field "%[1]v", condition "%[2]v"`, field, err.ActualTag())
		}
	}
	if err := requiredIfField(field, obj, rules); err != nil {
		return nil, err
	}

	return next(ctx)
}

func validate() *validator.Validate {
	// https://github.com/go-playground/validator/issues/494
	// register custom validation: rfe(Required if Field is Equal to some value).
	v := validator.New()
	v.RegisterValidation(`rfe`, func(fl validator.FieldLevel) bool { return true })
	return v
}

func requiredIfField(field interface{}, obj interface{}, rules string) error {
	const rfeTag = "rfe="
	tags := strings.Split(rules, `,`) // example: tags := "email,rfe=id:1"
	tag := tags[find(tags, rfeTag)]   // find and get string "rfe=id:1"

	// If rfe not exist, bypass validation
	if !(len(tag) > 0) {
		return nil
	}

	/*
		example: {
			status: "NEW",
			newData: "1", // @rfe=status:NEW
			oldData: "0", // @rfe=status:OLD
		}

		param 		= [rfe=status, NEW]
		paramField 	= status (remove string "rfe=")
		paramValue  = NEW

		paramFieldValue = "NEW" (value from root object "status")
		selfValue		= "1"   (value from root object "newData")
	*/

	param := strings.Split(tag, `:`)
	paramField := strings.Replace(param[0], rfeTag, "", 1)
	paramValue := fmt.Sprintf("%v", param[1])

	paramFieldValue := fmt.Sprintf("%v", obj.(map[string]interface{})[paramField])
	selfValue := fmt.Sprintf("%s", obj.(map[string]interface{})[field.(string)])

	/*
		Error example: {
			status: "NEW",
			newData: "",  // @rfe=status:NEW
			oldData: "x", // @rfe=status:OLD
		}

		1./ same ref value and data is empty
		2./ diff ref value and data is NOT empty
	*/
	if (paramValue == paramFieldValue && selfValue == "") || (paramValue != paramFieldValue && selfValue != "") {
		return fmt.Errorf(`Field: '%[1]v' is required field '%[2]v' to match value: '%[3]v'`, field, paramField, paramValue)
	}
	return nil
}

// Find returns the smallest index i at which x == a[i],
// or len(a) if there is no such index.
func find(a []string, x string) int {
	for i, n := range a {
		if strings.Contains(n, x) {
			return i
		}
	}
	return len(a)
}
