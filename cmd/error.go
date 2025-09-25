package cmd

import "errors"

// JoinErrors returns an error that wraps the given errors with any nil value discarded.
// Like errors.Join, JoinErrors returns nil if every value in errs is nil.
// Unlike errors.Join, any previously joined error in errs are flattened.
func JoinErrors(errs ...error) error {
	flattened := make([]error, 0)
	for _, err := range errs {
		flattened = append(flattened, UnwrapJoinedErrors(err)...)
	}
	return errors.Join(flattened...)
}

// UnwrapJoinedErrors calls errs.Unwrap if it returns []error (i.e. errs is the result of
// errors.Join).
// If errs.Unwrap returns an error, UnwrapJoinedErrors returns []error{errs.Unwrap()}.
// Otherwise, UnwrapJoinedErrors returns []error{errs}.
func UnwrapJoinedErrors(errs error) []error {
	if joinedErrs, ok := errs.(interface{ Unwrap() []error }); ok {
		return joinedErrs.Unwrap()
	} else if unwrapable, ok := errs.(interface{ Unwrap() error }); ok {
		return []error{unwrapable.Unwrap()}
	} else {
		return []error{errs}
	}
}
