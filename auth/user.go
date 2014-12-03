package auth

import (
	"net/http"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strings"
)

var (
	ErrPaswordMismatch = errors.New("Passwords supplied do not match")
)

type User struct {
	Id int64
	Email string
	PwHash string
	Admin bool
	Active bool
	Token string
}

type Storer interface {
	UserByEmail(string) (*User, error)
	UserById(int64) (*User, error)
	UserByToken(int64) (*User, error)
	Delete(*User) error
	Insert(*User) error
	Update(*User) error
}

type Mailer interface {
	Signup(*User) error
	PasswordReset(*User) error
}

type UserManager struct {
	hasher PasswordHasher
	storer Storer
	mailer Mailer
	sess   SessionManager
}

func NewUserManager(h PasswordHasher, s Storer, m Mailer, sm SessionManager) UserManager {
	if h == nil { h = newBcryptHasher(0, 0, 0) } // Set default hasher with default values
	return UserManager{h, s, m, sm}
}

func (m UserManager) New(email, pw, confirmPw string, active, admin bool) (*User, error) {
	// Lower case all email addresses, for consistency.
	email = strings.ToLower(email)

	if pw != confirmPw {
		return ErrPasswordMismatch
	}

	hash, err := m.hasher.Hash(pw)
	if err != nil {
		return nil, err
	}
	tok := ""
	if !active {
		if tok, err := generateToken(); err != nil {
			return nil, err
		}
	}

	return &User{0, email, hash, active, admin, tok}, nil
}

func (m UserManager) SignupUser(email, pw, confirmPw string admin, active bool) (*User, error) {
	u, err := New(email, pw, confirmPw, active, admin)
	if err != nil {
		return nil, err
	}

	if !u.Active {
		if err = m.mailer.Signup(u); err != nil {
			return nil, err
		}
	}

	err = m.storer.Insert(u)
	if err != nil {
		return nil, err
	}

	// TODO: Create a session for the user
}

func (m UserManager) ById(id int64) (*User, error) {
	return m.storer.UserById(id)
}

func (m UserManager) ByEmail(email string) (*User, error) {
	return m.storer.UserByEmail(email)
}

func (m UserManager) ByToken(tok int64) (*User, error) {
	return m.storer.UserByToken(tok)
}

func (m UserManager) Insert(u *User) error {
	return m.storer.Save(u)
}

func (m UserManager) Update(u *User) error {
	return m.storer.Save(u)
}

func (m UserManager) FromSession(w http.ResponseWriter, h *http.Request) (*User, error) {
	return m.sess.User(w, h)
}

func (m UserManager) Authenticate(u *User, pw string) error {
	return hasher.Compare(u.PwHash, pw)
}

func (m UserManager) Activate(u *User) error {
	if u.Active { return nil }
	u.Active = true
	return m.Update(u)
}

func (m UserManager) Deactivate(u *User) error {
	if !u.Active {return nil}
	u.Active = false
	return m.Update(u)
}

func (m UserManager) SendSignupMail(u *User) error {
	return m.mailer.Signup(u)
}

func (m UserManager) RequestPwReset(u *User, disableCurrentPw bool) error {
	var err error
	u.Token, err = generateToken()
	if err != nil {
		return err
	}

	if disableCurrentPw {
		// Disable logins with current credentials
		u.PwHash = ""
	}

	err = m.Update(u)
	if err != nil {
		return err
	}

	return m.mailer.ResetPassword(u)
}

func (m UserManager) ResetPw(u *User, pw, confirm string) error {
	return nil
}


func generateToken() (string, error) {
	b := make([]byte, 64) // 512 bit token should be enough for anyone :)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// The "=" at the end can get lost when emailing links for some reason
	// so ensure that there are no "=" at the end for linking or storing.
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "="), nil
}

