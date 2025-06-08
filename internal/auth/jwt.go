package auth

import (
	"time"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"errors"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	myLovelySigningKey := os.Getenv("JWT_SIGNING_KEY")
	if myLovelySigningKey == "" {
		return "", errors.New("could not get signing key from env")
	}

	now := time.Now().UTC()
	claims := &jwt.RegisteredClaims{
		Issuer: "chirpy",
		Subject: userID.String(),
		IssuedAt: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(myLovelySigningKey))
	if err != nil {
		return "", err
	}
	return ss, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	myLovelySigningKey := os.Getenv("JWT_SIGNING_KEY")
	if myLovelySigningKey == "" {
		return uuid.Nil, errors.New("jwt signing key env var either missing or is problematic")
	}
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token)(any, error){
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}
	userIdStr := claims.Subject
	if userIdStr == "" {
		return uuid.Nil, errors.New("token subj (userID) is missing")
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return uuid.Nil, err
	}
	return userId, nil
}
