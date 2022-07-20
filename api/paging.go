package api

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

const (
	// SkipDefault is the default number of items to skip when paging.
	SkipDefault = 0

	// TakeDefault is the default number of items to take when paging.
	TakeDefault = 25

	// SkipMin is the minimum number of items to skip.
	SkipMin = 0

	// TakeMin is the minimum number of items that can be taken.
	TakeMin = 1

	// TakeMax is the maximum number of items that can be taken.
	TakeMax = 250
)

// NormalizePages adjusts the skip and take parameters to be sane if needed.
func NormalizePages(skip, take int) (int, int) {
	if skip < SkipMin {
		skip = SkipDefault
	}

	if take < TakeMin {
		take = TakeDefault
	}

	if take > TakeMax {
		take = TakeMax
	}

	return skip, take
}

// WithPaging will extract a skip and page query parameter from echo context and
// run a function with those parameters.
func WithPaging(c echo.Context, action func(skip, take int) error) error {
	skipStr := c.QueryParam("skip")
	takeStr := c.QueryParam("take")

	if skipStr == "" {
		skipStr = strconv.Itoa(SkipDefault)
	}

	if takeStr == "" {
		takeStr = strconv.Itoa(TakeDefault)
	}

	skip, err := strconv.Atoi(skipStr)
	if err != nil {
		return c.String(http.StatusBadRequest, ErrInvalidIntFormat.Error())
	}

	take, err := strconv.Atoi(takeStr)
	if err != nil {
		return c.String(http.StatusBadRequest, ErrInvalidIntFormat.Error())
	}

	return action(NormalizePages(skip, take))
}
