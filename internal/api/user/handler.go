package user

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"Dash/core/user"
	"Dash/pkg/logger"
)

type Handler struct {
	service *user.Service
}

func NewHandler(service *user.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Login:VALIDATION_FAILED", map[string]interface{}{
			"error": err.Error(),
		})

		if validationErr, ok := err.(validator.ValidationErrors); ok {
			fields := make(map[string]string)
			for _, fieldErr := range validationErr {
				fieldName := getJSONFieldName(req, fieldErr.Field())
				fields[fieldName] = getValidationErrorMessage(fieldErr)
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"error":  "validation_error",
				"fields": fields,
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	phoneNumber, err := user.NewPhoneNumber(req.Phone)
	if err != nil {
		logger.Error("Login:INVALID_PHONE_FORMAT", map[string]interface{}{
			"phone": req.Phone,
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Login(phoneNumber, req.Password)
	if err != nil {
		if err == user.ErrCredentialsNotValid || err == user.ErrInputCredentialsInvalid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid phone or password"})
			return
		}
		logger.Error("Login:SERVICE_ERROR", map[string]interface{}{
			"phone": phoneNumber.String(),
			"error": err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login"})
		return
	}

	response := LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
	c.JSON(http.StatusOK, response)
}

func getJSONFieldName(req LoginRequest, structFieldName string) string {
	t := reflect.TypeOf(req)
	field, found := t.FieldByName(structFieldName)
	if !found {
		return structFieldName
	}

	jsonTag := field.Tag.Get("json")
	if jsonTag == "" || jsonTag == "-" {
		return structFieldName
	}

	for idx := 0; idx < len(jsonTag); idx++ {
		if jsonTag[idx] == ',' {
			return jsonTag[:idx]
		}
	}
	return jsonTag
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
