package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message,omitempty"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			if validationErr, ok := err.Err.(validator.ValidationErrors); ok {
				fields := make(map[string]string)
				for _, fieldErr := range validationErr {
					fields[fieldErr.Field()] = getValidationErrorMessage(fieldErr)
				}

				c.JSON(http.StatusBadRequest, ErrorResponse{
					Error:  "validation_error",
					Fields: fields,
				})
				return
			}

			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "internal_server_error",
				Message: err.Error(),
			})
		}
	}
}

func getValidationErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value is too short"
	case "max":
		return "Value is too long"
	case "numeric":
		return "Value must be numeric"
	default:
		return "Invalid value"
	}
}

