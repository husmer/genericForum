package validateData

import "net/mail"

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func ValidatePassword(password1, password2 string) (bool, bool) {
	passLen := false
	passMatch := false
	if len(password1) >= 6 {
		passLen = true
	}
	if password1 == password2 {
		passMatch = true
	}
	return passLen, passMatch
}

func ValidateName(firstName string) bool {
	name := false
	if len(firstName) >= 3 {
		name = true
	}
	return name
}
