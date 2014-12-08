package models

type User struct {
	Id                  int64
	Email               string
	HashedPassword      string
	IsAdmin             bool
	IsActive            bool
	SignupToken         string
	PasswordChangeToken string
	IsNew               bool
}

func generateToken() (string, error) {
	b := make([]byte, settings.GetInt("Auth", "TokenBytes"))
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "="), nil
}

func NewUser(email, pw1 string, pw2 string, isAdmin bool) (*User, error) {
	// Ensure email is lower case or we can create the user here with an
	// upper case email, it gets saved to the database and converted to
	// lower case and then the hash will never match as when logging in
	// the email is retrieved from the database.
	email = strings.ToLower(email)
	if pw1 != pw2 {
		return nil, NewUserError("Given passwords do not match.")
	}

	if err := password.ValidatePassword(pw1); err != nil {
		return nil, err
	}

	hashInput, err := password.GenerateHashInputFromStrings(email, pw1)
	if err != nil {
		return nil, err
	}

	hash, err := password.GenerateHash(hashInput)
	if err != nil {
		return nil, err
	}
	tok, err := generateToken()
	if err != nil {
		return nil, err
	}
	u := User{Email: email, HashedPassword: hash, IsAdmin: isAdmin, SignupToken: tok}

	return &u, nil
}
