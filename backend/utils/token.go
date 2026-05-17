package utils

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/azdonald/pharmd/backend/models"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	ID             string
	OrganisationID string
}

type JwtCustomClaims struct {
	Name       string `json:"name"`
	ID         string `json:"id"`
	BusinessID string `json:"business_id"`
	jwt.RegisteredClaims
}

type JwtCustomRefreshClaims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

func CreateAccessToken(user models.User) (string, error) {
	tokenExpiry, err := strconv.Atoi(os.Getenv("ACCESSTOKEN_EXPIRY"))
	if err != nil {
		log.Println("failed to get access token expiry, defaulting to 15mins")
		tokenExpiry = 15
	}
	exp := time.Now().Add(time.Minute * time.Duration(tokenExpiry))

	claims := JwtCustomClaims{
		Name:       user.FirstName + " " + user.LastName,
		ID:         user.ID,
		BusinessID: user.OrganisationID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	key := os.Getenv("TOKEN_KEY")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(key))
}

func CreateRefreshToken(user models.User) (string, error) {
	tokenExpiry, err := strconv.Atoi(os.Getenv("REFRESHTOKEN_EXPIRY"))
	if err != nil {
		log.Println("failed to get refresh token expiry, defaulting to 24hrs")
		tokenExpiry = 24
	}
	exp := time.Now().Add(time.Hour * time.Duration(tokenExpiry))

	claims := JwtCustomRefreshClaims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	key := os.Getenv("TOKEN_KEY")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(key))
}

func ExtractClaimFromToken(reqToken string) (*JwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(reqToken, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_KEY")), nil
	})

	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JwtCustomClaims); ok {
		return claims, nil
	}

	return nil, errors.New("unknown claims type")
}

func ExtractRefreshClaimFromToken(reqToken string) (*JwtCustomRefreshClaims, error) {
	token, err := jwt.ParseWithClaims(reqToken, &JwtCustomRefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_KEY")), nil
	})

	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JwtCustomRefreshClaims); ok {
		return claims, nil
	}

	return nil, errors.New("unknown refresh claims type")
}

func ExtractClaimFromContext(ctx context.Context) *Claims {
	jwtclaims, ok := ctx.Value("userClaims").(*JwtCustomClaims)
	if !ok {
		return nil
	}

	return &Claims{
		ID:             jwtclaims.ID,
		OrganisationID: jwtclaims.BusinessID,
	}
}
