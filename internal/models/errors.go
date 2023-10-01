package models

type UserAlreadyExistsError struct{}

func (e *UserAlreadyExistsError) Error() string {
	return "user already exists"
}
