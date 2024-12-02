package types

import "fmt"

var (
	ErrInvalidState = fmt.Errorf("invalid chrome state")
	ErrInitTimeout  = fmt.Errorf("chrome initialization timeout")
	ErrNotAvailable = fmt.Errorf("chrome instance not available")
)

type ChromeError struct {
	ID  int
	Op  string
	Err error
}

func (e *ChromeError) Error() string {
	return fmt.Sprintf("chrome %d: %s: %v", e.ID, e.Op, e.Err)
}

func (e *ChromeError) Unwrap() error {
	return e.Err
}
