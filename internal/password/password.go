package password

import "golang.org/x/crypto/bcrypt"

func HashPassword(plaintext string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plaintext), 12)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func Verify(plainText, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainText)) == nil
}
