package azuredevops

import (
	"github.com/giantswarm/microerror"
)

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var unexpectedResponseError = &microerror.Error{
	Kind: "unexpectedResponseError",
}

// IsUnexpectedResponse asserts unexpectedResponseError.
func IsUnexpectedResponse(err error) bool {
	return microerror.Cause(err) == unexpectedResponseError
}
