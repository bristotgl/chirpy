package auth

import (
	"testing"
)

func TestCheckPasswordHash(t *testing.T) {
	password1 := "correctPassword"
	hash1, _ := HashPassword(password1)

	cases := []struct {
		name          string
		password      string
		hash          string
		matchPassword bool
		wantErr       bool
	}{
		{
			name:          "Correct password",
			password:      password1,
			hash:          hash1,
			matchPassword: true,
			wantErr:       false,
		},
		{
			name:          "Incorrect password",
			password:      "wrong_password",
			hash:          hash1,
			matchPassword: false,
			wantErr:       false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      password1,
			hash:          "invalid_hash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			match, err := CheckPasswordHash(c.password, c.hash)
			if (err != nil) != c.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr = %v", err, c.wantErr)
				return
			}

			if !c.wantErr && match != c.matchPassword {
				t.Errorf("CheckPasswordHash() got %v but want %v", match, c.matchPassword)
			}
		})
	}
}
