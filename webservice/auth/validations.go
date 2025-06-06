package auth

import "regexp"

var (
	usernameRegex       = regexp.MustCompile(`^[a-zA-z0-9_-]{3,32}$`)
	emailRegex          = regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)
	hasLowerCharRegex   = regexp.MustCompile(`[a-z]`)
	hasUpperCharRegex   = regexp.MustCompile(`[A-Z]`)
	hasDigitRegex       = regexp.MustCompile(`\d`)
	hasSpecialCharRegex = regexp.MustCompile(`[!@#$%^&*_]`)
	lengthRegex         = regexp.MustCompile(`^.{8,72}$`)
)

func ValidateUsername(username string) bool {
	return usernameRegex.MatchString(username)
}

func ValidateEmail(mail string) bool {
	return emailRegex.MatchString(mail)
}

func ValidatePassword(password string) bool {
	return hasLowerCharRegex.MatchString(password) && hasUpperCharRegex.MatchString(password) &&
		hasDigitRegex.MatchString(password) && hasSpecialCharRegex.MatchString(password) &&
		lengthRegex.MatchString(password)
}
