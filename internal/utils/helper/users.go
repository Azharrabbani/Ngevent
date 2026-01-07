package helper

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		return "", err
	}

	return string(hashPassword), nil
}

func ValidePassword(hashPassoword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassoword), []byte(password))

	return err
}