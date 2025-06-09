package auth

import (
	"testing"
	"time"
	"errors"
	"net/http"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert" // Import the assert package
	"github.com/golang-jwt/jwt/v5"
)

func TestJWT(t *testing.T) {
	type testCase struct {
		name	string
		userId	uuid.UUID
		tokenSecret	string
		expiresIn	time.Duration
	}
	
	validTests := []testCase{
		{
			"validTest1",
			uuid.New(),
			"shhh this is a secret",
			time.Second,
		},
	}

	for _, test := range validTests {
		t.Run(test.name, func(t * testing.T) {
			signedString, err := MakeJWT(test.userId, test.tokenSecret, test.expiresIn)
			if err != nil {
				t.Errorf("err in MakeJWT %v", err)
			}

			userId, err := ValidateJWT(signedString, test.tokenSecret)
			if err != nil {
				t.Errorf("err in validateJWT: %v", err)
			}
			if test.userId != userId{
				t.Errorf("userId is not equal to userId")
			}

		})
	}

	//invalid tests
	test := testCase{
		"jwt expires",
		uuid.New(),
		"secretHereYo",
		(100*time.Nanosecond),
	}
	t.Run(test.name, func(t *testing.T) {
		signedString, err := MakeJWT(test.userId, test.tokenSecret, test.expiresIn)
		if err != nil {
			t.Errorf("err in MakeJWT %v", err)
		}

		time.Sleep(120*time.Nanosecond)

		_, err = ValidateJWT(signedString, test.tokenSecret)
		assert.Truef(t, errors.Is(err, jwt.ErrTokenExpired), "Expected invalid token from expiring, ERR: %v", err)
	})
}

func TestGetBearerToken(t *testing.T) {
	type testCase struct {
		name			string
		headahString 	string
		expectErr		bool
		expected		string
	}

	tests := []testCase{
		{
			"validTest1",
			"Bearer sillywurds ",
			false,
			"sillywurds",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t * testing.T) {
			headah := http.Header{} 
			headah.Add("Authorization", test.headahString)

			token, err := GetBearerToken(headah)
			if err != nil {
				if !test.expectErr {
					assert.Errorf(t, err, "error unexpected")
				}
				return
			}

			if token != test.expected {
				assert.Failf(t, "token wrong :(, token: %v, expected %v", token, test.expected)
			}
		})
	}
}
