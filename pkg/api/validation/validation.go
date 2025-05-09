package validation

import (
	"fmt"
	"strings"

	"github.com/distribution/reference"

	"github.com/kubesphere/s2irun/pkg/api"
)

// ValidateConfig returns a list of error from validation.
func ValidateConfig(config *api.Config) []Error {
	allErrs := []Error{}
	if len(config.BuilderImage) == 0 {
		allErrs = append(allErrs, NewFieldRequired("builderImage"))
	}
	switch config.BuilderPullPolicy {
	case api.PullNever, api.PullAlways, api.PullIfNotPresent:
	default:
		allErrs = append(allErrs, NewFieldInvalidValueObj("builderPullPolicy", config.BuilderPullPolicy))
	}
	if config.DockerConfig == nil || len(config.DockerConfig.Endpoint) == 0 {
		allErrs = append(allErrs, NewFieldRequired("dockerConfig.endpoint"))
	}
	if config.DockerNetworkMode != "" && !validateDockerNetworkMode(config.DockerNetworkMode) {
		allErrs = append(allErrs, NewFieldInvalidValueObj("dockerNetworkMode", config.DockerNetworkMode))
	}
	if config.Labels != nil {
		for k := range config.Labels {
			if len(k) == 0 {
				allErrs = append(allErrs, NewFieldInvalidValueWithReason("labels", "contains empty label"))
			}
		}
	}
	if config.Tag != "" {
		if err := validateDockerReference(config.Tag); err != nil {
			allErrs = append(allErrs, NewFieldInvalidValueWithReasonAndValue("tag", err.Error(), config.Tag))
		}
	}
	return allErrs
}

// validateDockerNetworkMode checks wether the network mode conforms to the docker remote API specification (v1.19)
// Supported values are: bridge, host, container:<name|id>, and netns:/proc/<pid>/ns/net
func validateDockerNetworkMode(mode api.DockerNetworkMode) bool {
	switch mode {
	case api.DockerNetworkModeBridge, api.DockerNetworkModeHost:
		return true
	}
	if strings.HasPrefix(string(mode), api.DockerNetworkModeContainerPrefix) {
		return true
	}
	if strings.HasPrefix(string(mode), api.DockerNetworkModeNetworkNamespacePrefix) {
		return true
	}
	return false
}

func validateDockerReference(ref string) error {
	_, err := reference.Parse(ref)
	return err
}

// NewFieldRequired returns a *ValidationError indicating "value required"
func NewFieldRequired(field string) Error {
	return Error{Type: ErrorTypeRequired, Field: field}
}

// NewFieldInvalidValue returns a ValidationError indicating "invalid value"
func NewFieldInvalidValue(field string) Error {
	return Error{Type: ErrorInvalidValue, Field: field}
}

// NewFieldInvalidValueObj returns a ValidationError indicating "invalid value"
func NewFieldInvalidValueObj(field string, value interface{}) Error {
	return Error{Type: ErrorInvalidValue, Field: field, Value: value}
}

// NewFieldInvalidValueWithReason returns a ValidationError indicating "invalid value" and a reason for the error
func NewFieldInvalidValueWithReason(field, reason string) Error {
	return Error{Type: ErrorInvalidValue, Field: field, Reason: reason}
}

// NewFieldInvalidValueWithReasonAndValue returns a ValidationError indicating the value and reason
func NewFieldInvalidValueWithReasonAndValue(field, reason string, value interface{}) Error {
	return Error{Type: ErrorInvalidValue, Field: field, Value: value, Reason: reason}
}

// ErrorType is a machine readable value providing more detail about why a field
// is invalid.
type ErrorType string

const (
	// ErrorTypeRequired is used to report required values that are not provided
	// (e.g. empty strings, null values, or empty arrays).
	ErrorTypeRequired ErrorType = "FieldValueRequired"

	// ErrorInvalidValue is used to report values that do not conform to the
	// expected schema.
	ErrorInvalidValue ErrorType = "InvalidValue"
)

// Error is an implementation of the 'error' interface, which represents an
// error of validation.
type Error struct {
	Type   ErrorType
	Field  string
	Value  interface{}
	Reason string
}

func (v Error) Error() string {
	var msg string
	switch v.Type {
	case ErrorInvalidValue:
		msg = fmt.Sprintf("Invalid value specified for %q, value: %v", v.Field, v.Value)
	case ErrorTypeRequired:
		msg = fmt.Sprintf("Required value not specified for %q", v.Field)
	default:
		msg = fmt.Sprintf("%s: %s", v.Type, v.Field)
	}
	if len(v.Reason) > 0 {
		msg = fmt.Sprintf("%s: %s", msg, v.Reason)
	}
	return msg
}
