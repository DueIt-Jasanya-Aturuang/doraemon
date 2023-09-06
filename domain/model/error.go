package model

import (
	"fmt"
)

type ErrResponseHTTP struct {
	Code    int
	Message any
}

func (e *ErrResponseHTTP) Error() string {
	return fmt.Sprintf("%d: %d", e.Code, e.Message)
}
