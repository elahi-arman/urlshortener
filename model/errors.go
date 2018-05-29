package model

//NotFoundError signifies that a given link was not found
type NotFoundError struct {
	msg string // description of error
}

func (e NotFoundError) Error() string { return e.msg }
