package lib

import "errors"

var (
	// handler
	ErrFailedToParseRequest   = errors.New("failed to parse request")
	ErrUnknownError           = errors.New("unknown error")
	ErrInsufficientPermission = errors.New("insufficient permission")

	// lib.jwt_claims
	ErrFailedToParseJWTClaimsInContext = errors.New("failed to parse jwt claims in context")
	ErrJWTClaimsNotFoundInContext      = errors.New("jwt claims not found in context")
)
