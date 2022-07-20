package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

const (
	// NotFoundMessage is the string written to the response
	// body when a given entity cannot be found.
	NotFoundMessage = "Not Found"
)

var (
	// ErrInvalidIntFormat is an error saying what was given was not a valid int.
	ErrInvalidIntFormat = FormatError{errors.New("Invalid int format")}

	// ErrInvalidBoolFormat is an error saying what was given was not a valid bool.
	ErrInvalidBoolFormat = FormatError{errors.New("Invalid bool format")}
)

// FormatError is an alias for error.
type FormatError struct {
	error
}

// NotFound is a convenience method to say an entity was not found in
// an http response.
func NotFound(c echo.Context) error {
	return c.String(http.StatusNotFound, NotFoundMessage)
}

// WithInt extracts an int from the route with the given name.
// It then runs a given action with that extracted parameter.
// If the value in the route is not a valid int, it will return
// an ErrInvalidIntFormat.
func WithInt(c echo.Context, action func(ints map[string]int) error, params ...string) error {
	paramMap := make(map[string]int)

	for _, param := range params {
		iStr := c.Param(param)

		i, err := strconv.Atoi(iStr)
		if err != nil {
			return ErrInvalidIntFormat
		}

		paramMap[param] = i
	}

	return action(paramMap)
}

// WithBoolQuery extracts bools from the query string with the given names.
// It then runs a given action with the extracted parameters.
// If a value in the query string is not a valid bool, it will return
// an ErrInvalidBoolFormat.
func WithBoolQuery(c echo.Context, action func(map[string]bool) error, names ...string) error {
	bools := make(map[string]bool)

	for _, name := range names {
		b, err := parseBool(c.QueryParam(name))
		if err != nil {
			return ErrInvalidBoolFormat
		}

		bools[name] = b
	}

	return action(bools)
}

func parseBool(bStr string) (bool, error) {
	var b bool

	if bStr != "" {
		bo, err := strconv.ParseBool(bStr)
		if err != nil {
			return b, ErrInvalidBoolFormat
		}

		b = bo
	}

	return b, nil
}

// WithStringQuery extracts strings from the query string with the given names.
// It then runs a given action with the extracted parameters.
func WithStringQuery(c echo.Context, action func(map[string]string) error, names ...string) error {
	strs := make(map[string]string)

	for _, name := range names {
		s := c.QueryParam(name)

		strs[name] = s
	}

	return action(strs)
}
