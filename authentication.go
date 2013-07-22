package mvc

import (
	"code.google.com/p/go.crypto/bcrypt"
)

func EncryptPassword(password string) (string, error) {
	enc, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(enc), err
}
