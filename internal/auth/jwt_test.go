package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
	validUuid := uuid.New()
	validSecret := "secret"
	validToken, _ := MakeJWT(validUuid, validSecret, time.Hour)
	expiredToken, _ := MakeJWT(validUuid, validSecret, time.Hour*-5)

	cases := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Happy path",
			tokenString: validToken,
			tokenSecret: validSecret,
			wantUserID:  validUuid,
			wantErr:     false,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Expired token",
			tokenString: expiredToken,
			tokenSecret: validSecret,
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid_token_string",
			tokenSecret: validSecret,
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(c.tokenString, c.tokenSecret)
			if (err != nil) != c.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr = %v", err, c.wantErr)
				return
			}

			if gotUserID != c.wantUserID {
				t.Errorf("ValidateJWT() got %v but want %v", gotUserID, c.wantUserID)
			}
		})
	}
}
