package errs

// CustomError is a simple error type with just a message.
type CustomError struct {
	Message string
}

// Error implements the error interface.
func (e *CustomError) Error() string {
	return e.Message
}

// New creates a new CustomError with the given message.
func New(message string) error {
	return &CustomError{Message: message}
}
