package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/whiskey-back/pkg/infrastruct"
	"net/http"
	"strings"
)

type Manager struct{}

func NewManager() (*Manager, error) {
	return &Manager{}, nil
}

const (
	authorizationHeader = "Authorization"
)

type CustomClaims struct {
	UserID string `json:"id"`
	jwt.StandardClaims
}

func (m *Manager) GenerateJWT(userID string, secretKey string) (string, error) {

	claims := CustomClaims{
		userID,
		jwt.StandardClaims{},
	}

	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := tokenJWT.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (m *Manager) GetClaimsByRequest(r *http.Request, JWTkey string) (*CustomClaims, error) {

	header := r.Header.Get(authorizationHeader)
	if header == "" {
		return nil, infrastruct.ErrorJWTIsBroken
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, infrastruct.ErrorJWTIsBroken
	}

	if len(headerParts[1]) == 0 {
		return nil, infrastruct.ErrorJWTIsBroken
	}

	token, err := m.ValidateJwt(headerParts[1], JWTkey)
	if err != nil {
		return nil, infrastruct.ErrorJWTIsBroken
	}

	if claims, ok := token.Claims.(*CustomClaims); ok {
		err = claims.Valid()
		if err != nil {
			return nil, err
		}
		return claims, nil
	}

	return nil, infrastruct.ErrorJWTIsBroken
}

func (m *Manager) ValidateJwt(tokenString string, key string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unknown signature method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})
	if err != nil {
		return token, err
	}

	if !token.Valid {
		return token, fmt.Errorf("token is not valid, %v", token)
	}
	return token, nil
}

func (c CustomClaims) Valid() error {

	if c.UserID == "" {
		return infrastruct.ErrorJWTIsBroken
	}

	return nil
}
