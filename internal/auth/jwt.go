package auth

import (
	"time"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
	"errors"
)

var (
	errEnv		= errors.New("could not get signing key from env")
	errSignedString = errors.New("error signing string")
	errUnexpectedSigningMethod = errors.New("unexpected signing method")
	errInvalidToken			= errors.New("invalid token")
	errTokenParse			= errors.New("error parsing token")
	errSubjectMissing		= errors.New("token subj (userID) is missing")
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	now := time.Now().UTC()
	claims := &jwt.RegisteredClaims{
		Issuer: "chirpy",
		Subject: userID.String(),
		IssuedAt: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return ss, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token)(any, error){
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errUnexpectedSigningMethod
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err 
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return uuid.Nil, jwt.ErrTokenExpired
	}

	if !token.Valid {
		return uuid.Nil, errInvalidToken
	}
	userIdStr := claims.Subject
	if userIdStr == "" {
		return uuid.Nil, errSubjectMissing 
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return uuid.Nil, err
	}
	return userId, nil
}
