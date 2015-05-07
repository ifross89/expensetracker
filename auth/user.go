package auth

import (
	"github.com/juju/errors"

	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"
	"fmt"
	"time"
)

var (
	ErrPwMismatch   = errors.New("Passwords supplied do not match")
	ErrNoToken      = errors.New("There is no token associated with the user")
	ErrAlreadySaved = errors.New("Cannot insert as user already saved")
)

type User struct {
	ID        int64      `db:"id" json:"id"`
	Email     string     `db:"email" json:"email"`
	PwHash    string     `db:"pw_hash" json:"-"`
	Admin     bool       `db:"admin" json:"-"`
	Active    bool       `db:"active" json:"active"`
	Token     string     `db:"token" json:"token"`
	Name      string     `db:"name" json:"name"`
	CreatedAt *time.Time `db:"created_at" json:"createdAt"`
}

func (u *User) String() string {
	return fmt.Sprintf("%s <%s>", u.Name, u.Email)
}

type Storer interface {
	UserByEmail(string) (*User, error)
	UserByID(int64) (*User, error)
	UserByToken(string) (*User, error)
	Users() ([]*User, error)
	Delete(*User) error
	Insert(*User) error
	Update(*User) error
}

type Mailer interface {
	Signup(*User) error
	PasswordReset(*User) error
}

type nopMailer struct{}

func (nopMailer) Signup(*User) error {
	return nil
}

func (nopMailer) PasswordReset(*User) error {
	return nil
}

type UserManager struct {
	hasher PasswordHasher
	store  Storer
	mailer Mailer
	sess   SessionStore
}

// NewUserManager creates an object which can be used to manipulate User objects.
func NewUserManager(h PasswordHasher, s Storer, m Mailer, sm SessionStore) *UserManager {
	// Set default hasher with default values
	if h == nil {
		h = NewBcryptHasher(0, 0, 0)
	}
	if m == nil {
		m = &nopMailer{}
	}

	return &UserManager{h, s, m, sm}
}

// New creates a new user. Note that this only creates the user, it does
// not save the user.
func (m UserManager) New(name, email, pw, confirmPw string, active, admin bool) (*User, error) {
	// Lower case all email addresses, for consistency.
	email = strings.ToLower(email)

	if pw != confirmPw {
		return nil, errors.Trace(ErrPwMismatch)
	}

	hash, err := m.hasher.Hash(pw)
	if err != nil {
		return nil, errors.Trace(err)
	}
	tok := ""
	if !active {
		if tok, err = generateToken(); err != nil {
			return nil, errors.Trace(err)
		}
	}

	return &User{0, email, hash, active, admin, tok, name, nil}, nil
}

func (m UserManager) Users() ([]*User, error) {
	return m.store.Users()
}

// SignupUser creates a new user, emails a signup email, saves the user and logs them in.
func (m UserManager) SignupUser(w http.ResponseWriter, r *http.Request, name, email, pw, confirmPw string, admin, active bool) (*User, error) {
	u, err := m.New(name, email, pw, confirmPw, active, admin)
	if err != nil {
		return nil, errors.Trace(err)
	}

	if !u.Active {
		if err = m.mailer.Signup(u); err != nil {
			return nil, errors.Trace(err)
		}
	}

	err = m.store.Insert(u)
	if err != nil {
		return nil, errors.Trace(err)
	}

	err = m.sess.LogUserIn(w, r, u)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return u, nil
}

// ById obtains the user by their id field
func (m UserManager) ById(id int64) (*User, error) {
	return m.store.UserByID(id)
}

// ByEmail obtains the user from the email supplied
func (m UserManager) ByEmail(email string) (*User, error) {
	return m.store.UserByEmail(email)
}

// ByToken obtains the user from the generated token. This can be used for password
// resets or signup links
func (m UserManager) ByToken(tok string) (*User, error) {
	return m.store.UserByToken(tok)
}

// Insert saves a new User. This cannot be called if the user has already
// been saved.
func (m UserManager) Insert(u *User) error {
	return m.store.Insert(u)
}

// Update saves any changes to the User object. This can only be called
// if the user has already been saved.
func (m UserManager) Update(u *User) error {
	return m.store.Update(u)
}

// FromSession retrieves the current user associated with the session
func (m UserManager) FromSession(w http.ResponseWriter, r *http.Request) (*User, error) {
	return m.sess.User(w, r, m.store)
}

func (m UserManager) AdminFromSession(w http.ResponseWriter, r *http.Request) (*User, error) {
	u, err := m.FromSession(w, r)

	if err != nil {
		return nil, errors.Annotate(err, "Error getting admin from session")
	}

	if !u.Active {
		return nil, errors.Errorf("Error getting admin from session: user %s not active", u)
	}

	if !u.Admin {
		return nil, errors.Errorf("Error getting admin from session: user %s not admin", u)
	}

	return u, nil
}

// Authenticate checks to see if the password supplies is the same as the
// password that was used to create the hash
func (m UserManager) Authenticate(u *User, pw string) error {
	return m.hasher.Compare(u.PwHash, pw)
}

func (m UserManager) LogOut(w http.ResponseWriter, r *http.Request) error {
	return m.sess.LogUserOut(w, r)
}

func (m UserManager) LogIn(w http.ResponseWriter, r *http.Request, u *User) error {
	return m.sess.LogUserIn(w, r, u)
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
	if err := m.mailer.Signup(u); err != nil {
		return err
	}
	return nil
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

	// Only send an email if the mailer is set
	if err = m.mailer.PasswordReset(u); err != nil {
		return err
	}
	return nil
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
	return m.store.Update(u)
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

func (m UserManager) DeleteUser(u *User) error {
	return m.store.Delete(u)
}

func (m UserManager) DeleteUserById(id int64) error {
	return m.store.Delete(&User{ID: id})
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
