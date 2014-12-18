package auth

import (
	"github.com/gorilla/sessions"
	"github.com/juju/errors"

	"net/http"
)

const (
	cookieSessionName = "cookie-session"
	emailKey          = "email"
)

var (
	ErrNoSession = errors.New("No session present for the request")
	ErrWrongType = errors.New("The value stored in the session was the wrong type")
)

type SessionStore interface {
	User(http.ResponseWriter, *http.Request, Storer) (*User, error)
	LogUserOut(http.ResponseWriter, *http.Request) error
	LogUserIn(http.ResponseWriter, *http.Request, *User) error
}

type cookieSessionStore struct {
	store *sessions.CookieStore
}

func NewCookieSessionStore(pairs ...[]byte) *cookieSessionStore {
	cs := &cookieSessionStore{store: sessions.NewCookieStore(pairs...)}
	return cs
}

func (s *cookieSessionStore) session(r *http.Request) (*sessions.Session, error) {
	return s.store.Get(r, cookieSessionName)
}

func (s *cookieSessionStore) LogUserOut(w http.ResponseWriter, r *http.Request) error {
	sess, err := s.session(r)
	if err != nil {
		// could not get the session. No need to log out
		// TODO: log this when I have figured out a logging strategy.
		return nil
	}

	delete(sess.Values, emailKey)
	return s.store.Save(r, w, sess)
}

func (s *cookieSessionStore) User(w http.ResponseWriter, r *http.Request, us Storer) (*User, error) {
	sess, err := s.session(r)
	if err != nil {
		return nil, errors.Annotate(err, "No cookie session present for request")
	}

	val, ok := sess.Values[emailKey]
	if !ok {
		return nil, errors.Trace(ErrNoSession)
	}

	email, ok := val.(string)
	if !ok {
		return nil, errors.Trace(ErrWrongType)
	}

	u, err := us.UserByEmail(email)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return u, nil
}

func (s *cookieSessionStore) LogUserIn(w http.ResponseWriter, r *http.Request, u *User) error {
	// TODO: Log any error returned. A session is always returned, so do not want
	// to return early as it could be the first time.
	sess, _ := s.session(r)

	sess.Values[emailKey] = u.Email
	err := s.store.Save(r, w, sess)
	if err != nil {
		if logoutErr := s.LogUserOut(w, r); logoutErr != nil {
			return errors.Wrap(logoutErr, err)
		}
		return errors.Trace(err)
	}

	return nil
}
