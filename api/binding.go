package api

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

var (
	// ErrQueryBindPtr is an error stating the a given item must be a pointer.
	ErrQueryBindPtr = errors.New("item needs to be a pointer")

	// ErrQueryBindStruct is an error stating that a given item's pointer value must be a struct.
	ErrQueryBindStruct = errors.New("pointer value needs to be a struct")

	// ErrInvalidTimeFormat is an error stating that a given item is not a valid time format.
	ErrInvalidTimeFormat = echo.NewHTTPError(http.StatusBadRequest, "invalid time format. supported format is RFC3339")

	// ErrInvalidNumberFormat is an error stating that the given item is not a valid number format.
	ErrInvalidNumberFormat = echo.NewHTTPError(http.StatusBadRequest, "invalid number format")

	// ErrInvalidBooleanFormat is an error stating that the given item is not a valid boolean format.
	ErrInvalidBooleanFormat = echo.NewHTTPError(http.StatusBadRequest, "invalid boolean format")
)

// Bind combines the echo bind function along with a model validator.
// To use this, your model must implement the models.Validator interface.
func Bind(c echo.Context, model Validator) error {
	if err := c.Bind(model); err != nil {
		return err
	}

	model.Escape()

	if err := model.Validate(); err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	return nil
}

// QueryBinder will bind query parameters to a given struct based on an echo.Context.
type QueryBinder interface {
	BindQuery(echo.Context, interface{}) error
}

// BindQuery is a convenience function to bind a query with the DefaultQueryBinder.
func BindQuery(c echo.Context, item interface{}) error {
	return new(DefaultQueryBinder).BindQuery(c, item)
}

// DefaultQueryBinder is the default query binder.
type DefaultQueryBinder struct{}

// BindQuery binds (populates a struct) based on values in a query string.
func (s *DefaultQueryBinder) BindQuery(c echo.Context, item interface{}) error {
	rv := reflect.ValueOf(item)

	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return ErrQueryBindPtr
	}

	elem := rv.Elem()

	if elem.Kind() != reflect.Struct {
		return ErrQueryBindStruct
	}

	return s.loadData(c, elem)
}

func (s *DefaultQueryBinder) loadData(c echo.Context, r reflect.Value) error {
	for i := 0; i < r.NumField(); i++ {
		tag := r.Type().Field(i).Tag.Get("query")

		field := r.Field(i)

		param := c.QueryParam(tag)
		if field.IsValid() && field.CanSet() && param != "" {
			if err := s.set(field, param); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *DefaultQueryBinder) set(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.Int:
		i, err := strconv.Atoi(value)
		if err != nil {
			return ErrInvalidNumberFormat
		}

		s.setValue(field, i)
	case reflect.Int16:
		i, err := strconv.ParseInt(value, 10, 16)
		if err != nil {
			return ErrInvalidNumberFormat
		}

		s.setValue(field, int16(i))
	case reflect.Int32:
		i, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return ErrInvalidNumberFormat
		}

		s.setValue(field, int32(i))
	case reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return ErrInvalidNumberFormat
		}

		s.setValue(field, int64(i))
	case reflect.String:
		s.setValue(field, value)
	case reflect.Float32:
		f, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return ErrInvalidNumberFormat
		}

		s.setValue(field, float32(f))
	case reflect.Float64:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return ErrInvalidNumberFormat
		}

		s.setValue(field, f)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return ErrInvalidBoolFormat
		}

		s.setValue(field, b)
	case reflect.Ptr:
		field.Set(reflect.New(field.Type().Elem()))
		return s.set(field.Elem(), value)
	case reflect.Struct:
		switch field.Type() {
		case reflect.TypeOf(time.Time{}):
			if t, err := time.Parse(time.RFC3339, value); err == nil {
				s.setValue(field, t)
			} else {
				return ErrInvalidTimeFormat
			}
		}
	}

	return nil
}

func (s *DefaultQueryBinder) setValue(field reflect.Value, value interface{}) {
	field.Set(reflect.ValueOf(value))
}
