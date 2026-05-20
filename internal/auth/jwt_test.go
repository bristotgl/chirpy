package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	type ExpectedValues struct {
		JWT string
		Err error
	}

	cases := map[string]struct {
		UserID      uuid.UUID
		TokenSecret string
		ExpiresIn   time.Duration
		Expected    ExpectedValues
	}{
		"happy path": {
			UserID:      uuid.New(),
			TokenSecret: "secret",
			ExpiresIn:   time.Minute * 5,
			Expected: ExpectedValues{
				JWT: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHktYWNjZXNzIiwic3ViIjoiNTBkYjlhN2ItM2U0Ny00NTU3LWFkZDctMjkyZDJiZWI0NzAyIiwiZXhwIjoxNzc5MjM2MDkwLCJpYXQiOjE3NzkyMzU3OTB9.0H6K4ta3ou7DmRWGa5gw-GS6th0byVPNKx3hQqup_UY",
				Err: nil,
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual, err := MakeJWT(c.UserID, c.TokenSecret, c.ExpiresIn)
			if err != nil {
				t.Errorf("Error creating tokens with given params. Got error %v but want %v", err, c.Expected.Err)
				return
			}

			if len(actual) < 1 {
				t.Errorf("Invalid token: got %v but want %v", actual, c.Expected.JWT)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	type ExpectedValues struct {
		UserID uuid.UUID
	}

	validUuid := uuid.New()
	validSecret := "secret"
	validToken, _ := MakeJWT(validUuid, validSecret, time.Minute*5)

	expiredToken, _ := MakeJWT(validUuid, validSecret, time.Minute*-5)

	cases := map[string]struct {
		TokenString string
		TokenSecret string
		Expected    ExpectedValues
	}{
		"happy path": {
			TokenString: validToken,
			TokenSecret: validSecret,
			Expected: ExpectedValues{
				UserID: validUuid,
			},
		},
		"wrong secret": {
			TokenString: validToken,
			TokenSecret: "wrong_secret",
			Expected: ExpectedValues{
				UserID: uuid.Nil,
			},
		},
		"expired token": {
			TokenString: expiredToken,
			TokenSecret: validSecret,
			Expected: ExpectedValues{
				UserID: uuid.Nil,
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual, err := ValidateJWT(c.TokenString, c.TokenSecret)
			if err != nil && c.Expected.UserID != uuid.Nil {
				t.Errorf("Error validating token: %v", err)
				return
			} 	

			if actual != c.Expected.UserID {
				t.Errorf("UserIDs are different: got %v but want %v", actual, c.Expected.UserID)
			}
		})
	}
}
