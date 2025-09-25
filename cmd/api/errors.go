package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorMessage struct {
	reason  string
	message string
	details *string
}

func newErrorMessage(reason string, message string, details string) *ErrorMessage {
	return &ErrorMessage{
		reason:  reason,
		message: message,
		details: &details,
	}
}

func (s *service) malformedAuthTokenResponse(c *gin.Context) {
	err := newErrorMessage("INVALID_HEADERS", "incorrect auth header format", "ensure the header follows 'Bearer <token>' format")
	c.JSON(400, err)
}

func (s *service) expiredTokenResponse(c *gin.Context) {
	err := newErrorMessage("INVALID_HEADERS", "incorrect auth header format", "ensure the header follows 'Bearer <token>' format")
	c.JSON(400, err)
}

// notFoundResponse sends a 404 Not Found response.
func (s *service) notFoundResponse(c *gin.Context) {
	err := newErrorMessage("NOT_FOUND", "resource not found", "the requested resource could not be found")
	c.JSON(http.StatusNotFound, err)
}

// methodNotAllowedResponse sends a 405 Method Not Allowed response.
func (s *service) methodNotAllowedResponse(c *gin.Context) {
	message := fmt.Sprintf("the %s method is not supported for this resource", c.Request.Method)
	err := newErrorMessage("METHOD_NOT_ALLOWED", "method not supported", message)
	c.JSON(http.StatusMethodNotAllowed, err)
}

// badRequestResponse sends a 400 Bad Request response.
func (s *service) badRequestResponse(c *gin.Context, errMsg string) {
	err := newErrorMessage("BAD_REQUEST", "invalid request", errMsg)
	c.JSON(http.StatusBadRequest, err)
}

// failedValidationResponse sends a 422 Unprocessable Entity response for validation failures.
// func FailedValidationResponse(c *gin.Context, validationErrors map[string]string) {
// 	var details string
// 	for key, value := range validationErrors {
// 		details += fmt.Sprintf("%s: %s; ", key, value)
// 	}
// 	err := newErrorMessage("VALIDATION_FAILED", "invalid data submitted", details)
// 	c.JSON(http.StatusUnprocessableEntity, err)
// }

// editConflictResponse sends a 409 Conflict response.
func (s *service) editConflictResponse(c *gin.Context) {
	err := newErrorMessage("EDIT_CONFLICT", "edit conflict", "unable to update the record due to an edit conflict, please try again")
	c.JSON(http.StatusConflict, err)
}

// rateLimitExceededResponse sends a 429 Too Many Requests response.
func (s *service) rateLimitExceededResponse(c *gin.Context) {
	err := newErrorMessage("RATE_LIMIT_EXCEEDED", "rate limit exceeded", "you have exceeded the allowed number of requests")
	c.JSON(http.StatusTooManyRequests, err)
}

// invalidCredentialsResponse sends a 401 Unauthorized response for bad credentials.
func (s *service) invalidCredentialsResponse(c *gin.Context) {
	err := newErrorMessage("INVALID_CREDENTIALS", "invalid credentials", "invalid authentication credentials")
	c.JSON(http.StatusUnauthorized, err)
}

// invalidAuthenticationTokenResponse sends a 401 Unauthorized response with the WWW-Authenticate header.
func (s *service) invalidAuthenticationTokenResponse(c *gin.Context) {
	c.Header("WWW-Authenticate", "Bearer")
	err := newErrorMessage("INVALID_TOKEN", "invalid token", "invalid or missing authentication token")
	c.JSON(http.StatusUnauthorized, err)
}

// authenticationRequiredResponse sends a 401 Unauthorized response for unauthenticated requests.
func (s *service) authenticationRequiredResponse(c *gin.Context) {
	err := newErrorMessage("AUTH_REQUIRED", "authentication required", "you must be authenticated to access this resource")
	c.JSON(http.StatusUnauthorized, err)
}

// inactiveAccountResponse sends a 403 Forbidden response for inactive accounts.
func (s *service) inactiveAccountResponse(c *gin.Context) {
	err := newErrorMessage("ACCOUNT_INACTIVE", "account not active", "your user account must be activated to access this resource")
	c.JSON(http.StatusForbidden, err)
}

// notPermittedResponse sends a 403 Forbidden response for insufficient permissions.
func (s *service) notPermittedResponse(c *gin.Context) {
	err := newErrorMessage("PERMISSION_DENIED", "permission denied", "your user account doesn't have the necessary permissions to access this resource")
	c.JSON(http.StatusForbidden, err)
}

func (s *service) serviceUnavailableResponse(c *gin.Context) {
	errMsg := "the server is temporarily unable to handle the request"
	response := newErrorMessage("SERVICE_UNAVAILABLE", "service unavailable", errMsg)
	c.JSON(http.StatusServiceUnavailable, response)
}

func (s *service) InternalServerErrorResponse(c *gin.Context, err error) {
	errMsg := "the server encountered a problem and could not process your request"
	response := newErrorMessage("INTERNAL_SERVER_ERROR", "internal server error", errMsg)
	c.JSON(http.StatusInternalServerError, response)
}
