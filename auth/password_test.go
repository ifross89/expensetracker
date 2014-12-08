package auth

import (
	"testing"
)

func TestPasswordValidation(t *testing.T) {
	hasher := NewBcryptHasher(2, 4, 4)
	tests := []struct {
		pw       string
		expected error
	}{
		{pw: "a", expected: ErrPwTooShort},
		{pw: "ab", expected: nil},
		{pw: "abcd", expected: nil},
		{pw: "abcde", expected: ErrPwTooLong},
	}
	for _, test := range tests {
		if _, err := hasher.Hash(test.pw); err != test.expected {
			t.Fatalf("Expected %v, got %v", test.expected, err)
			return
		}
	}
}

func TestHashMatching(t *testing.T) {
	hasher := NewBcryptHasher(6, 20, 4)
	tests := []struct {
		pw1, pw2 string
		expected error
	}{
		{pw1: "matching", pw2: "matching", expected: nil},
		{pw1: "notMatching", pw2: "different", expected: ErrIncorrectPw},
		{pw1: "caseSensitive", pw2: "casesensitive", expected: ErrIncorrectPw},
	}
	for _, test := range tests {
		hash, err := hasher.Hash(test.pw1)
		if err != nil {
			t.Fatalf("Error during hash creation: %v (pw=%s)", err, test.pw1)
			return
		}

		err = hasher.Compare(hash, test.pw2)
		if err != test.expected {
			t.Fatalf("Expected: %v, got: %v (pw1=%s, pw2=%s)", test.expected, err, test.pw1, test.pw2)
			return
		}

		// Ensure that this is a symmetric relationship
		hash, err = hasher.Hash(test.pw2)
		if err != nil {
			t.Fatalf("Error during hash creation: %v (pw=%s)", err, test.pw2)
			return
		}

		err = hasher.Compare(hash, test.pw1)
		if err != test.expected {
			t.Fatalf("Expected: %v, got: %v (pw1=%s, pw2=%s)", test.expected, err, test.pw2, test.pw1)
			return
		}
	}
}

func TestBcryptSalt(t *testing.T) {
	hashers := []PasswordHasher{
		NewBcryptHasher(6, 25, 5),
		NewBcryptHasher(6, 25, 5),
		NewBcryptHasher(6, 25, 6),
	}

	testPasswords := []string{"abcdefg"}
	for _, pw := range testPasswords {
		hashes := make([]string, len(hashers), len(hashers))
		for j, hasher := range hashers {
			hash, err := hasher.Hash(pw)
			if err != nil {
				t.Fatalf("Error during hash creation: %v (pw=%s)", err, pw)
				return
			}
			hashes[j] = hash
		}
		for j, hasher := range hashers {
			err := hasher.Compare(hashes[j], pw)
			if err != nil {
				t.Fatalf("Error during sanity check, hasher %d failing compare, pw=%s, err=%s", j, pw, err)
			}

		}
		if hashes[0] == hashes[1] {
			t.Fatal("Two passwords hashed produced the same hash (not salting properly)")
			return
		}
		if hashes[0] == hashes[2] {
			t.Fatalf("Two hashers with different costs produce the same hashes. 1: pw:%s, %s, 2: %s", pw, hashes[0], hashes[2])
			return
		}
	}
}
