package common

import "errors"

func UnwrapAll(err error) []error {
	var errs []error
	for err != nil {
		errs = append(errs, err)
		err = errors.Unwrap(err)
	}
	return errs
}
