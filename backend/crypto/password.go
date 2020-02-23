package crypto

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) []byte {
	b := []byte(password)

	hashedPassword, err := bcrypt.GenerateFromPassword(b, bcrypt.DefaultCost)

	if err != nil {
		panic(err)
	}

	return hashedPassword
}

func ComparePasswords(provided string, saved []byte) error {
	return bcrypt.CompareHashAndPassword(saved, []byte(provided))
}
