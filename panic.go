package fn

import "fmt"

type ErrPanic struct {
	V any
}

func (p ErrPanic) Error() string {
	return fmt.Sprintf("recovered panic: %v", p.V)
}

func (p ErrPanic) Unwrap() error {
	if err, ok := p.V.(error); ok {
		return err
	} else {
		return nil
	}
}
