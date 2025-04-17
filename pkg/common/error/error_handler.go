package error

import "log"

// DomainError represents a domain-specific error
type DomainError struct {
	Code    string
	Message string
}

// Error implements the error interface
func (e *DomainError) Error() string {
	return e.Message
}

// NewDomainError creates a new DomainError instance
func NewDomainError(code, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
	}
}

// HandleError logs the error and returns a user-friendly message
func HandleError(err error, userMessage string) string {
	if err != nil {
		log.Printf("Error: %v", err)
	}
	return userMessage
}
