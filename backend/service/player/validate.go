package player

import (
	"net/mail"
	"unicode/utf8"

	"github.com/importandrew/MaCRa/models"
)

const (
	// maxNameLength matches the users.name VARCHAR(50) column, which counts
	// characters, not bytes.
	maxNameLength = 50
	// maxEmailLength matches the users.email VARCHAR(50) column.
	maxEmailLength = 50
	// maxPasswordBytes is bcrypt's limit: it refuses longer passwords.
	maxPasswordBytes = 72
)

// checkName reports whether name is a valid player name: non-empty and at
// most maxNameLength characters.
func checkName(name string) error {
	if name == "" {
		return models.ErrEmptyName
	}

	// RuneCountInString, not len: len counts bytes, and a Cyrillic letter is
	// two bytes in UTF-8, so a 26-letter Russian name was rejected.
	if utf8.RuneCountInString(name) > maxNameLength {
		return models.ErrLongName
	}

	return nil
}

// checkEmail reports whether email is a plausible bare address (no display
// name, no surrounding spaces) of at most maxEmailLength characters.
func checkEmail(email string) error {
	if email == "" {
		return models.ErrEmptyEmail
	}

	if utf8.RuneCountInString(email) > maxEmailLength {
		return models.ErrLongEmail
	}

	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return models.ErrInvalidEmail
	}

	return nil
}

// checkNewPassword reports whether password can be set: non-empty and at most
// maxPasswordBytes bytes. Short passwords are allowed on purpose.
func checkNewPassword(password string) error {
	if password == "" {
		return models.ErrEmptyPassword
	}

	if len(password) > maxPasswordBytes {
		return models.ErrLongPassword
	}

	return nil
}
