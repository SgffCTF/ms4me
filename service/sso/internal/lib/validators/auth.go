package validators

import "fmt"

var (
	ErrUsernameEmpty    = fmt.Errorf("username cannot be empty")
	ErrUsernameTooShort = fmt.Errorf("username must be at least 3 characters long")
	ErrUsernameTooLong  = fmt.Errorf("username cannot be longer than 64 characters")
	ErrPasswordEmpty    = fmt.Errorf("password cannot be empty")
	ErrPasswordTooShort = fmt.Errorf("password must be at least 8 characters long")
	ErrPasswordTooLong  = fmt.Errorf("password cannot be longer than 64 characters")
)

func ValidateUsername(username string) error {
	if username == "" {
		return ErrUsernameEmpty
	}
	if len(username) < 3 {
		return ErrUsernameTooShort
	}
	if len(username) > 64 {
		return ErrUsernameTooLong
	}
	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return ErrPasswordEmpty
	}
	if len(password) < 8 {
		return ErrPasswordTooShort
	}
	if len(password) > 64 {
		return ErrPasswordTooLong
	}
	return nil
}
