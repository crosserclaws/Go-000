package error

import (
	"fmt"
)

// CustomError is an error which encapsulates internal errors that should not be depended by other layers.
type CustomError struct {
	ErrorCode *customErrorCode
	cause     error
}

func (c *CustomError) String() string {
	return c.Error()
}

func (c *CustomError) Error() string {
	return c.ErrorCode.Error()
}

// Is method is used for checking if this error has the same code as the target.
func (c *CustomError) Is(target error) bool {
	t, ok := target.(*CustomError)
	if !ok {
		code, yes := target.(*customErrorCode)
		if !yes {
			return false
		}
		return c.ErrorCode == code
	}
	return c.ErrorCode == t.ErrorCode
}

// UnWrap is used for checking the cause if it's necessary.
func (c *CustomError) UnWrap() error {
	return c.cause
}

type customErrorCode struct {
	code        string
	description string
}

func (c *customErrorCode) Error() string {
	return fmt.Sprintf("[%s] %s", c.code, c.description)
}
