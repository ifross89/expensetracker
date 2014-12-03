package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
)

var (
	ErrPwMismatch = errors.New("Passwords supplied do not match")
	ErrNoToken    = errors.New("There is no token associated with the user")
)

type User struct {
	Id     int64
	Email  string
	PwHash string
	Admin  bool
	Active bool
	Token  string
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

// NewUserManager creates an object which can be used to manipulate User objects.
func NewUserManager(h PasswordHasher, s Storer, m Mailer, sm SessionManager) UserManager {
	if h == nil {
		h = newBcryptHasher(0, 0, 0)
	} // Set default hasher with default values
	return UserManager{h, s, m, sm}
}

// New creates a new user. Note that this only creates the user, it does
// not save the user.
func (m UserManager) New(email, pw, confirmPw string, active, admin bool) (*User, error) {
	// Lower case all email addresses, for consistency.
	email = strings.ToLower(email)

	if pw != confirmPw {
		return ErrPwMismatch
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

// SignupUser creates a new user, emails a signup email and saves the user.
func (m UserManager) SignupUser(email, pw, confirmPw string, admin, active bool) (*User, error) {
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

	// TODO: Create a session for the user and log them in
}

// ById obtains the user by their id field
func (m UserManager) ById(id int64) (*User, error) {
	return m.storer.UserById(id)
}

// ByEmail obtains the user from the email supplied
func (m UserManager) ByEmail(email string) (*User, error) {
	return m.storer.UserByEmail(email)
}

// ByToken obtains the user from the generated token. This can be used for password
// resets or signup links
func (m UserManager) ByToken(tok int64) (*User, error) {
	return m.storer.UserByToken(tok)
}

// Insert saves a new User. This cannot be called if the user has already
// been saved.
func (m UserManager) Insert(u *User) error {
	return m.storer.Save(u)
}

// Update saves any changes to the User object. This can only be called
// if the user has already been saved.
func (m UserManager) Update(u *User) error {
	return m.storer.Save(u)
}

// FromSession retrieves the current user associated with the session
func (m UserManager) FromSession(w http.ResponseWriter, h *http.Request) (*User, error) {
	return m.sess.User(w, h)
}

// Authenticate checks to see if the password supplies is the same as the
// password that was used to create the hash
func (m UserManager) Authenticate(u *User, pw string) error {
	return hasher.Compare(u.PwHash, pw)
}

// Activate ensures that a user is able to log on.
func (m UserManager) Activate(u *User) error {
	if u.Active {
		return nil
	}
	u.Active = true
	return m.Update(u)
}

// Deactivate deactivates the user, disabling the user from logging on.
func (m UserManager) Deactivate(u *User) error {
	if !u.Active {
		return nil
	}
	u.Active = false
	return m.Update(u)
}

// SendSignupMail sends the user a signup email
func (m UserManager) SendSignupMail(u *User) error {
	if u.Token == "" {
		return ErrNoToken
	}
	return m.mailer.Signup(u)
}

// RequestPwReset sends a password reset email to the user and optionally
// disables the user from logging on/
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

// UpdatePw forces a password change. This can be useful for situations
// when the user does not know their password or if an admin wants to
// request a password change.
func (m UserManager) UpdatePw(u *User, pw, confirm string) error {
	if pw != confirm {
		return ErrPwMismatch
	}

	hash, err := m.hasher.Hash(pw)
	if err != nil {
		return err
	}

	u.PwHash = hash
	return m.storer.Update(u)
}

// UserResetPw is used when a user requests a password update.
// The user must supply the correct password in order for the update
// to be successful
func (m UserManager) UserResetPw(u *User, old, pw, confirm string) error {
	if err := m.Authenticate(u, old); err != nil {
		return err
	}

	return m.UpdatePw(u, pw, confirm)
}

// Helper function that generates random tokens. The length of the token
// created is currenly static (512 bit). The random sequence is base-64
// encoded into a string.
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
