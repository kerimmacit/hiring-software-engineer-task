package domain_errors

import "errors"

var (
	ErrLineItemNotFound       = errors.New("line item not found")
	ErrLineItemAlreadyUpdated = errors.New("line item already updated")
)
