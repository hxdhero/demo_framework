package validate

import "github.com/go-playground/validator/v10"

var V *validator.Validate

func init() {
	V = validator.New()
}
