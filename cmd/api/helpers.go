package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *service) readJSON(c *gin.Context, dst any) error {
	dec := json.NewDecoder(c.Request.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

// getRequiredStringEnv retrieves a required environment variable as a string.
// It will log a fatal error and exit if the variable is not set.
func getRequiredStringEnv(key string) (*string, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return nil, fmt.Errorf("fatal Error: Required environment variable '%s' not set", key)
	}
	return &value, nil
}

// getOptionalStringEnv retrieves an optional environment variable as a string,
// returning a fallback value if the variable is not set.
func getOptionalStringEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// getRequiredIntEnv retrieves a required environment variable as an integer.
// It will log a fatal error and exit if the variable is not set or is not a valid integer.
func getRequiredIntEnv(key string) (*int, error) {
	value, err := getRequiredStringEnv(key)
	if err != nil {
		return nil, err
	}

	intValue, err := strconv.Atoi(*value)
	if err != nil {
		return nil, fmt.Errorf("fatal Error: Environment variable '%s' is not a valid integer: %v", key, err)
	}
	return &intValue, nil
}

// getOptionalIntEnv retrieves an optional environment variable as an integer,
// returning a fallback value if the variable is not set or is not a valid integer.
func getOptionalIntEnv(key string, fallback int) int {
	value := getOptionalStringEnv(key, "")
	if value == "" {
		return fallback
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Warning: Environment variable '%s' is not a valid integer. Using fallback value %d. Error: %v", key, fallback, err)
		return fallback
	}
	return intValue
}

// getRequiredBoolEnv retrieves a required environment variable as a boolean.
// It will log a fatal error and exit if the variable is not set or is not a valid boolean.
func getRequiredBoolEnv(key string) (*bool, error) {
	value, err := getRequiredStringEnv(key)
	if err != nil {
		return nil, err
	}

	boolValue, err := strconv.ParseBool(*value)
	if err != nil {
		return nil, fmt.Errorf("fatal Error: Environment variable '%s' is not a valid boolean: %w", key, err)
	}
	return &boolValue, nil
}

// getOptionalBoolEnv retrieves an optional environment variable as a boolean,
// returning a fallback value if the variable is not set or is not a valid boolean.
func getOptionalBoolEnv(key string, fallback bool) bool {
	value := getOptionalStringEnv(key, "")
	if value == "" {
		return fallback
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("Warning: Environment variable '%s' is not a valid boolean. Using fallback value %t. Error: %v", key, fallback, err)
		return fallback
	}
	return boolValue
}
