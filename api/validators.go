package api

import (
	"github.com/aalug/blog-go/utils"
	"github.com/go-playground/validator/v10"
)

var isValidSlug validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if slug, ok := fieldLevel.Field().Interface().(string); ok {
		return utils.IsSlug(slug)
	}
	return false
}
