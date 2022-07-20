package api

// Validator is a type that can validate its own data.
type Validator interface {
	Validate() error
	Escape()
}
