package auth

import (
	"net/http"
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

func TestGetBearerToken(t *testing.T) {
	cases := []struct {
		name      string
		headers   http.Header
		wantToken string
		wantErr   bool
	}{
		{
			name:      "Valid Authorization header",
			headers:   http.Header{"Authorization": []string{"Bearer sfsdf"}},
			wantToken: "sfsdf",
			wantErr:   false,
		},
		{
			name:      "No Authorization header",
			headers:   http.Header{},
			wantToken: "",
			wantErr:   true,
		},
		{
			name:      "Invalid Authorization header",
			headers:   http.Header{"Authorization": []string{"Basic1234"}},
			wantToken: "",
			wantErr:   true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotToken, err := GetBearerToken(c.headers)
			if (err != nil) != c.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr = %v", err, c.wantErr)
				return
			}

			if gotToken != c.wantToken {
				t.Errorf("GetBearerToken() got %v, but want %v", gotToken, c.wantToken)
			}
		})
	}
}
