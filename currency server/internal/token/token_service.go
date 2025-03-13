package token

import (
	"crypto/rsa"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenService struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

type TSConfig struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func NewTokenService(c *TSConfig) *TokenService {
	return &TokenService{
		PrivateKey: c.PrivateKey,
		PublicKey:  c.PublicKey,
	}
}

type IDTokenCustomClaimns struct {
	UserId uint `json:"user_id"`
	jwt.StandardClaims
}

func (s *TokenService) ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Authorization"] == nil {
			fmt.Fprintf(w, "Authorizarion header is empty")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		jwtToken := strings.Split(r.Header["Authorization"][0], " ")[1]
		claims := &IDTokenCustomClaimns{}

		token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
			return s.PublicKey, nil
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Unauthrized "+err.Error())
			return
		}

		if token.Valid {
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Token is invalid")
			return
		}

	})
}

func GenerateToken(user_id *uint, key *rsa.PrivateKey) (string, error) {
	unixTime := time.Now().Unix()
	tokenExp := unixTime + 60*120 //120 min

	claims := IDTokenCustomClaimns{
		UserId: *user_id,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  unixTime,
			ExpiresAt: tokenExp,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, err := token.SignedString(key)

	if err != nil {
		log.Println("Failed to sign id token string")
		return "", err
	}

	return ss, nil
}

func (t *TokenService) GetUserInfo(header string) (uint, error) {
	claims := &IDTokenCustomClaimns{}

	_, err := jwt.ParseWithClaims(header, claims, func(token *jwt.Token) (interface{}, error) {
		return t.PublicKey, nil
	})
	if err != nil {
		return 0, err
	}

	return claims.UserId, nil
}
