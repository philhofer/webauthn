package protocol

import (
	"fmt"
)

type Error struct {
	// Short name for the type of error that has occurred.
	Type string `json:"type"`

	// Additional details about the error.
	Details string `json:"error"`

	// Information to help debug the error.
	DevInfo string `json:"debug"`

	// Inner error.
	Err error `json:"-"`
}

func (e *Error) Error() string {
	if e.Err == nil {
		return e.Details
	}
	return fmt.Sprintf("%s: %s", e.Details, e.Err)
}

func (e *Error) Unwrap() error {
	return e.Err
}

// Is compares an Error for [errors.Is].
//
// Is checks that [err] has type [*Error]
// and that its [Type] field is equal to [e.Type].
func (e *Error) Is(err error) bool {
	et, ok := err.(*Error)
	return ok && et.Type == e.Type
}

// WithDetails copies the error and sets [e.Details] to [details].
func (e *Error) WithDetails(details string) *Error {
	err := *e
	err.Details = details
	return &err
}

// WithInfo copies the error and sets [e.DevInfo] to [info].
func (e *Error) WithInfo(info string) *Error {
	err := *e
	err.DevInfo = info
	return &err
}

// WithError copies the error and sets [e.Err] to [err].
func (e *Error) WithError(err error) *Error {
	errCopy := *e
	errCopy.Err = err
	return &errCopy
}

var (
	ErrBadRequest = &Error{
		Type:    "invalid_request",
		Details: "error reading the request data",
	}
	ErrChallengeMismatch = &Error{
		Type:    "challenge_mismatch",
		Details: "stored challenge and received challenge do not match",
	}
	ErrParsingData = &Error{
		Type:    "parse_error",
		Details: "error parsing the authenticator response",
	}
	ErrAuthData = &Error{
		Type:    "auth_data",
		Details: "error verifying the authenticator data",
	}
	ErrVerification = &Error{
		Type:    "verification_error",
		Details: "error validating the authenticator response",
	}
	ErrAttestation = &Error{
		Type:    "attestation_error",
		Details: "error validating the attestation data provided",
	}
	ErrInvalidAttestation = &Error{
		Type:    "invalid_attestation",
		Details: "invalid attestation data",
	}
	ErrMetadata = &Error{
		Type:    "invalid_metadata",
		Details: "",
	}
	ErrAttestationFormat = &Error{
		Type:    "invalid_attestation",
		Details: "invalid attestation format",
	}
	ErrAttestationCertificate = &Error{
		Type:    "invalid_certificate",
		Details: "invalid attestation certificate",
	}
	ErrAssertionSignature = &Error{
		Type:    "invalid_signature",
		Details: "assertion Signature against auth data and client hash is not valid",
	}
	ErrUnsupportedKey = &Error{
		Type:    "invalid_key_type",
		Details: "unsupported Public Key Type",
	}
	ErrUnsupportedAlgorithm = &Error{
		Type:    "unsupported_key_algorithm",
		Details: "unsupported public key algorithm",
	}
	ErrNotSpecImplemented = &Error{
		Type:    "spec_unimplemented",
		Details: "this field is not yet supported by the WebAuthn spec",
	}
	ErrNotImplemented = &Error{
		Type:    "not_implemented",
		Details: "this field is not yet supported by this library",
	}
)
