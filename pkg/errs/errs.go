package errs

import "errors"

var (
	ErrInvalidStatusUpdate  = errors.New("invalid status update")
	ErrInvalidTransition    = errors.New("invalid state transition")
	ErrDriverNotAvailable   = errors.New("driver is not available")
	ErrOrderAlreadyAssigned = errors.New("order is already assigned")
	ErrDriverNotFound       = errors.New("driver not found")
	ErrOrderNotFound        = errors.New("order not found")
)
