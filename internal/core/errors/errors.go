package errors

import "fmt"

type BaseError struct {
	Code           string
	Message        string
	IsDisplayable  bool
	HTTPCode       int
	InternalReason *string
}

func (e *BaseError) Error() string {
	if e.InternalReason != nil {
		return fmt.Sprintf("%s: %s", e.Message, *e.InternalReason)
	}
	return e.Message
}

func (e *BaseError) GetCode() string {
	return e.Code
}

func (e *BaseError) GetMessage() string {
	return e.Message
}

func (e *BaseError) GetHTTPCode() int {
	return e.HTTPCode
}

func (e *BaseError) IsUserDisplayable() bool {
	return e.IsDisplayable
}

type DontWorryError struct {
	*BaseError
}

func NewDontWorryError(internalReason *string) *DontWorryError {
	return &DontWorryError{
		BaseError: &BaseError{
			Code:           "NO_ACTION_REQUIRED",
			Message:        "Não se preocupe, nenhuma ação é necessária!",
			IsDisplayable:  true,
			HTTPCode:       200,
			InternalReason: internalReason,
		},
	}
}

type ValidationError struct {
	*BaseError
	Details map[string]string
}

func NewValidationError(details map[string]string) *ValidationError {
	return &ValidationError{
		BaseError: &BaseError{
			Code:          "VALIDATION_ERROR",
			Message:       "Validation failed",
			IsDisplayable: true,
			HTTPCode:      400,
		},
		Details: details,
	}
}

type InvalidParameterValueError struct {
	*BaseError
	Parameter string
	Value     string
}

func NewInvalidParameterValueError(parameter, value string) *InvalidParameterValueError {
	return &InvalidParameterValueError{
		BaseError: &BaseError{
			Code:          "INVALID_PARAMETER_VALUE",
			Message:       fmt.Sprintf("Invalid value for parameter '%s': %s", parameter, value),
			IsDisplayable: true,
			HTTPCode:      400,
		},
		Parameter: parameter,
		Value:     value,
	}
}

type EntityNotFoundError struct {
	*BaseError
	EntityName string
}

func NewEntityNotFoundError(entityName, message string) *EntityNotFoundError {
	return &EntityNotFoundError{
		BaseError: &BaseError{
			Code:          "ENTITY_NOT_FOUND",
			Message:       message,
			IsDisplayable: true,
			HTTPCode:      404,
		},
		EntityName: entityName,
	}
}

type UnableToCreateEntityError struct {
	*BaseError
	EntityName string
}

func NewUnableToCreateEntityError(entityName, message string) *UnableToCreateEntityError {
	return &UnableToCreateEntityError{
		BaseError: &BaseError{
			Code:          "UNABLE_TO_CREATE_ENTITY",
			Message:       message,
			IsDisplayable: true,
			HTTPCode:      500,
		},
		EntityName: entityName,
	}
}

type UnableToUpdateEntityError struct {
	*BaseError
	EntityName string
}

func NewUnableToUpdateEntityError(entityName, message string) *UnableToUpdateEntityError {
	return &UnableToUpdateEntityError{
		BaseError: &BaseError{
			Code:          "UNABLE_TO_UPDATE_ENTITY",
			Message:       message,
			IsDisplayable: true,
			HTTPCode:      500,
		},
		EntityName: entityName,
	}
}

type InvalidIpAddressError struct {
	*BaseError
}

func NewInvalidIpAddressError() *InvalidIpAddressError {
	return &InvalidIpAddressError{
		BaseError: &BaseError{
			Code:          "INVALID_IP_ADDRESS",
			Message:       "Invalid IP address format",
			IsDisplayable: true,
			HTTPCode:      400,
		},
	}
}
