package constants

import "regexp"

// Regex for email
var EmailRegex = regexp.MustCompile(`^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$`)

// Regex for password
var PasswordRegex = regexp.MustCompile(`^(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{6,20}$`)
