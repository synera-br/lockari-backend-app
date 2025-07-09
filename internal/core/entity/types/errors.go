package types

import "fmt"

func ErrRepositoryNotFound(repository string) error {
	return fmt.Errorf("repository '%s' not found", repository)
}

func ErrServiceNotFound(service string) error {
	return fmt.Errorf("service '%s' not found", service)
}

func ErrHandlerNotFound(handler string) error {
	return fmt.Errorf("handler '%s' not found", handler)
}

func ErrInvalidRequest(message string) error {
	return fmt.Errorf("invalid request: %s", message)
}

func ErrInvalidResponse(message string) error {
	return fmt.Errorf("invalid response: %s", message)
}

func ErrInternalServer(message string) error {
	return fmt.Errorf("internal server error: %s", message)
}

func ErrUnauthorized(message string) error {
	return fmt.Errorf("unauthorized: %s", message)
}

func ErrForbidden(message string) error {
	return fmt.Errorf("forbidden: %s", message)
}

func ErrNotFound(message string) error {
	return fmt.Errorf("not found: %s", message)
}

func ErrConflict(message string) error {
	return fmt.Errorf("conflict: %s", message)
}

func ErrMethodNotAllowed(message string) error {
	return fmt.Errorf("method not allowed: %s", message)
}

func ErrNotImplemented(message string) error {
	return fmt.Errorf("not implemented: %s", message)
}

func ErrServiceUnavailable(message string) error {
	return fmt.Errorf("service unavailable: %s", message)
}

func ErrGatewayTimeout(message string) error {
	return fmt.Errorf("gateway timeout: %s", message)
}

func ErrBadRequest(message string) error {
	return fmt.Errorf("bad request: %s", message)
}

func ErrTooManyRequests(message string) error {
	return fmt.Errorf("too many requests: %s", message)
}

func ErrRequestTimeout(message string) error {
	return fmt.Errorf("request timeout: %s", message)
}

func ErrRequestEntityTooLarge(message string) error {
	return fmt.Errorf("request entity too large: %s", message)
}

func ErrGenericError(message string) error {
	return fmt.Errorf("%s", message)
}
