package services

import "github.com/go-playground/validator/v10"

type (
	// Validator provides validation mainly by validating structs within the web context.
	Validator struct {
		validator *validator.Validate
	}
)

// NewValidator creates a new Validator.
func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

// Validate validates the given interface as a struct. The struct type definition should include
// validation annotations.
func (v *Validator) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		return err
	}
	return nil
}
