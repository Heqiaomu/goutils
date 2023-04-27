package validate

import (
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type CusValidate struct {
	Tag string
	F   validator.Func
}

// CheckParam validate obj fields
func CheckParam(obj interface{}, cusvs ...CusValidate) (bool, error) {
	v := validator.New()
	for _, cv := range cusvs {
		err := v.RegisterValidation(cv.Tag, cv.F)
		if err != nil {
			return false, errors.Wrap(err, "register validation")
		}
	}
	//todo 优化参数检查方式
	if err := v.Struct(obj); err != nil {
		return false, err
	}
	return true, nil
}
