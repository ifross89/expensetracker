package auth

import (
	"code.google.com/p/go.crypto/bcrypt"

	"encoding/base64"
	"errors"
)

const (
	defaultMinLength = 6
	defaultMaxLength = 64
	defaultBcryptCost = 10
)

var (
	ErrPwTooShort = errors.New("Supplied password too short")
	ErrPwTooLong  = errors.New("Supplied password too long")
	ErrIncorrectPw = errors.New("Supplied password does not match the hash")
)

type PasswordHasher interface {
	Hash(string) (string, error)
	Compare(string, string) error
}


type bcryptHasher struct {
	minPwLen int
	maxPwLen int
	bcryptCost int
	
}

func defaultIfZero(value, def int) int {
	if value == 0 { return def }
	return value
}

func NewBcryptHasher(minPwLen, maxPwLen, bcryptCost int) PasswordHasher {
	minPwLen = defaultIfZero(minPwLen, defaultMinLength)
	maxPwLen = defaultIfZero(maxPwLen, defaultMaxLength)
	bcryptCost = defaultIfZero(bcryptCost, defaultBcryptCost)
	return bcryptHasher{minPwLen, maxPwLen, bcryptCost}
}

func (b bcryptHasher) Hash(pw string) (string, error) {
	if err := b.validatePw(pw); err != nil {
		return "", err
	}
		
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), b.bcryptCost)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(hash), nil
}

func (b bcryptHasher) validatePw(pw string) error {
	if len(pw) > b.maxPwLen {
		return ErrPwTooLong
	}

	if len(pw) < b.minPwLen {
		return ErrPwTooShort
	}

	return nil
}

func (b bcryptHasher) Compare(hash, pw string) error {
	blob, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword(blob, []byte(pw))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return ErrIncorrectPw
	}

	return err
}


