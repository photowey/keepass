package common

import (
	"errors"
	"fmt"
)

const (
	ExitCodeGeneric        = 1
	ExitCodeUsage          = 2
	ExitCodeNotInitialized = 3
	ExitCodeUnlockFailed   = 4
)

type CLIError struct {
	ExitCode int
	Err      error
}

func (e CLIError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("exit code %d", e.ExitCode)
	}
	return e.Err.Error()
}

func (e CLIError) Unwrap() error { return e.Err }

func WithExitCode(code int, err error) error {
	if err == nil {
		return nil
	}
	return CLIError{ExitCode: code, Err: err}
}

func UsageError(message string) error {
	return WithExitCode(ExitCodeUsage, errors.New(message))
}
