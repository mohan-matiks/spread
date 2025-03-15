package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var Validate = validator.New()

// Function to validate struct and return errors in a map format
// how to use:
// err := ValidateStruct(&App{})
//
//	if err != nil {
//		fmt.Println(err)
//	}
func ValidateStruct(s interface{}) map[string]string {
	err := Validate.Struct(s)
	if err != nil {
		errors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			errors[e.Field()] = fmt.Sprintf("Invalid value for %s", e.Field())
		}
		return errors
	}
	return nil
}

// BindAndValidate binds the request body to the struct and validates it
// how to use:
// err := BindAndValidate(c, &App{})
//
//	if err != nil {
//		fmt.Println(err)
//	}
func BindAndValidate(c *fiber.Ctx, i interface{}) []string {
	if err := c.BodyParser(i); err != nil {
		return []string{err.Error()}
	}

	// Validate request data
	if err := Validate.Struct(i); err != nil {
		var errors []string
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, e.Field()+" failed validation: "+e.Tag())
		}
		return errors
	}
	return nil
}
