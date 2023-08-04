package jwt

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"log"
	"net/http"
	"time"
)

var (
	jwtKey = []byte("secret")
)

type Claims struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	jwt.StandardClaims
}

func ValidateToken(tokenString string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("JWT validation failed: %v", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("JWT is not valid")
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")[(len("Bearer ")):]
		if tokenString == "" {
			log.Println("no header")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}

		claims, err := ValidateToken(tokenString)
		if err != nil {
			log.Println(err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GenerateJWT(id uuid.UUID, username, role string) (string, error) {
	expTime := time.Now().Add(60 * time.Minute)
	claims := &Claims{
		ID:       id,
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}
