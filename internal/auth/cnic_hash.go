package auth

import "golang.org/x/crypto/bcrypt"

func HashedCnic(cnic string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(cnic), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
