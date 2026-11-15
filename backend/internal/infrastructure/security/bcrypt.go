package security

import "golang.org/x/crypto/bcrypt"

// BcryptHasher implementa ports.PasswordHasher con costo 12 (sección 6).
type BcryptHasher struct{}

func NewBcryptHasher() *BcryptHasher { return &BcryptHasher{} }

func (BcryptHasher) Hash(plano string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plano), 12)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (BcryptHasher) Comparar(hash, plano string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plano)) == nil
}
