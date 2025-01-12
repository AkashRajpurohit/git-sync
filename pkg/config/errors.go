package config

type InvalidConfigError struct {
	message string
}

func (e *InvalidConfigError) Error() string {
	return e.message
}

func NewInvalidConfigError(message string) *InvalidConfigError {
	return &InvalidConfigError{message: message}
}
