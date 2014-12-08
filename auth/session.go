package auth

import (
	"github.com/gorilla/sessions"
)

type cookieSession struct {
	store *sessions.CookieStore
}
