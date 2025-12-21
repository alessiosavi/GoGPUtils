package aws

import (
	"errors"
	"fmt"
)

// Common sentinel errors.
var (
	// ErrNilConfig is returned when a nil config is provided.
	ErrNilConfig = errors.New("aws: config is nil")

	// ErrNilClient is returned when a nil client is provided.
	ErrNilClient = errors.New("aws: client is nil")

	// ErrEmptyBucket is returned when a bucket name is empty.
	ErrEmptyBucket = errors.New("aws: bucket name is empty")

	// ErrEmptyKey is returned when a key/path is empty.
	ErrEmptyKey = errors.New("aws: key is empty")

	// ErrEmptyTable is returned when a table name is empty.
	ErrEmptyTable = errors.New("aws: table name is empty")

	// ErrEmptyQueue is returned when a queue name/URL is empty.
	ErrEmptyQueue = errors.New("aws: queue name is empty")

	// ErrEmptySecret is returned when a secret name is empty.
	ErrEmptySecret = errors.New("aws: secret name is empty")

	// ErrEmptyParameter is returned when a parameter name is empty.
	ErrEmptyParameter = errors.New("aws: parameter name is empty")

	// ErrEmptyFunction is returned when a function name is empty.
	ErrEmptyFunction = errors.New("aws: function name is empty")

	// ErrInvalidInput is returned when input validation fails.
	ErrInvalidInput = errors.New("aws: invalid input")
)

// ConfigError represents an error loading AWS configuration.
type ConfigError struct {
	Err error
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("aws: failed to load config: %v", e.Err)
}

func (e *ConfigError) Unwrap() error {
	return e.Err
}

// OperationError represents an error during an AWS operation.
type OperationError struct {
	Service   string
	Operation string
	Err       error
}

func (e *OperationError) Error() string {
	return fmt.Sprintf("aws %s: %s failed: %v", e.Service, e.Operation, e.Err)
}

func (e *OperationError) Unwrap() error {
	return e.Err
}

// WrapError wraps an error with service and operation context.
func WrapError(service, operation string, err error) error {
	if err == nil {
		return nil
	}

	return &OperationError{
		Service:   service,
		Operation: operation,
		Err:       err,
	}
}

// NotFoundError represents a resource not found error.
type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("aws: %s not found: %s", e.Resource, e.ID)
}

// NewNotFoundError creates a new NotFoundError.
func NewNotFoundError(resource, id string) *NotFoundError {
	return &NotFoundError{
		Resource: resource,
		ID:       id,
	}
}

// ValidationError represents an input validation error.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("aws: validation failed for %s: %s", e.Field, e.Message)
}

// NewValidationError creates a new ValidationError.
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}
