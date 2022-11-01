package exceptions

import "errors"

var (
	ErrUserNotLoggedIn = errors.New("ErrUserNotLoggedIn")
)

type ErrLoginError struct {
	Reason string
}

func (ErrLoginError) Error() string {
	return "ErrLoginError"
}

type ErrMissingRequirements struct {
	MissingRequirements []string
}

func (ErrMissingRequirements) Error() string {
	return "ErrMissingRequirements"
}

type ErrVSCodeError struct {
	Logs         string
	ErrorMessage string
}

func (ErrVSCodeError) Error() string {
	return "ErrVSCodeError"
}
