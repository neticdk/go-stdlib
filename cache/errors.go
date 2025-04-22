package cache

type ErrNotFound struct{}

func (e *ErrNotFound) Error() string {
	return "cache: not found"
}

func (e *ErrNotFound) Unwrap() error {
	return nil
}

func (e *ErrNotFound) Code() int {
	return 0
}

func NewErrNotFound() error {
	return &ErrNotFound{}
}

type ErrExpired struct{}

func (e *ErrExpired) Error() string {
	return "cache: element expired"
}

func (e *ErrExpired) Unwrap() error {
	return nil
}

func (e *ErrExpired) Code() int {
	return 0
}

func NewErrExpired() error {
	return &ErrExpired{}
}
