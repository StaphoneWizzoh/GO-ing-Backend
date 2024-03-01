package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)


var (
	jwtAccessSecret  = []byte("eK_sS2AgDstNYrh0Bx5LK3nPx-z1h2l_ZdjchgQjvyA=")
	jwtRefreshSecret = []byte("bS4RqAvfuWhiAjZJ_104wBUcDAbp4cEt2ChP1IYskI8=")
)

type UserClaims struct{
	UserID 		uuid.UUID 		`json:"userId"`
	Username	string			`json:"username"`
	Email		string			`json:"email"`
	Role		string			`json:"role"`
	jwt.RegisteredClaims
}

func GenerateTokens(userID uuid.UUID, username, email, role string)(string, string, time.Time, error){
	// Generating access token
	accessToken, err := generateAccessToken(userID, username, email, role)
	if err != nil{
		return "", "", time.Time{}, err
	}

	// Generating refresh token
	refreshToken, expireTime, err := generateRefreshToken(userID, username, email, role)
	if err != nil{
		return "", "", time.Time{}, err
	}

	return accessToken, refreshToken, expireTime, nil
}

func generateAccessToken(userID uuid.UUID, username, email, role string) (string,  error){

	// Creating claims
	claims := UserClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // expires in 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID.String(),
		},
	}

	// Creating token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Signing token
	return token.SignedString(jwtAccessSecret)
}

func generateRefreshToken(userID uuid.UUID, username string, email string, role string) (string, time.Time, error) {

	// Token expires in 90 days (3 months)
	expireTime := time.Now().Add(24 * 90 * time.Hour) 

	// Creating claims
	claims := UserClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID.String(),
		},
	}

	// Creating token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Signing token
	refreshToken, err := token.SignedString(jwtRefreshSecret)
	if err != nil {
		return "", time.Time{}, err
	}

	return refreshToken, expireTime, nil
}

func ParseToken(tokenString string, isAccessToken bool)(*UserClaims, error){
	var claims UserClaims
	var jwtSecret []byte

	if isAccessToken{
		jwtSecret = jwtAccessSecret
	}else{
		jwtSecret = jwtRefreshSecret
	}

	// Parsing token
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token)(interface{}, error){
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok{
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return jwtSecret, nil
	})

	if err != nil{
		return nil, fmt.Errorf("error parsing token: %v", err)
	}

	// Checking if the token is valid
	if !token.Valid{
		return nil, jwt.ErrSignatureInvalid
	}

	return &claims, nil
}

// RefreshToken generates a new access token
func RefreshToken(refreshToken string) (string, error){
	
	// Parsing refresh token
	claims, err := ParseToken(refreshToken, false)
	if err != nil {
		return "", err
	}

	// Validation
	if claims.UserID.String() == "" || claims.Username == "" || claims.Email == "" || claims.Role == "" {
		return "", errors.New("invalid token claims")
	}

	// Generating new access token
	newAccessToken, err := generateAccessToken(claims.UserID, claims.Username, claims.Email, claims.Role)
	if err != nil{
		return "", err
	}

	return newAccessToken, nil
}

// ValidateToken validates a given token
func ValidateToken(tokenString string, isAccessToken bool) error{
	var jwtSecret []byte

	if isAccessToken {
		jwtSecret = jwtAccessSecret
	} else {
		jwtSecret = jwtRefreshSecret
	}

	// Parsing token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token)(interface{}, error){
		return jwtSecret, nil
	})
	if err != nil {
		return err
	}

	// Checking if the token is valid
	if !token.Valid{
		return jwt.ErrSignatureInvalid
	}

	return nil
}

// IsTokenExpired checks if a given token is expired
func IsTokenExpired(tokenString string, isAccessToken bool) bool{
	var jwtSecret []byte

	if isAccessToken {
		jwtSecret = jwtAccessSecret
	} else {
		jwtSecret = jwtRefreshSecret
	}

	// Parsing token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return true
	}

	// Checking if token is valid
	if !token.Valid {
		return true
	}

	// Checking if token is expired
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid{
		exp := claims["exp"].(float64)
		if time.Now().Unix() > int64(exp){
			return true
		}
	}

	return false
}

func ExtractUserRoleFromToken(tokenString string)(string, error){
	// Parse Token
	token, err := jwt.ParseWithClaims( tokenString, &UserClaims{}, func(token *jwt.Token)(interface{}, error){
		// Secret Key
		return []byte(jwtAccessSecret), nil
	})
	if err != nil{
		if err == jwt.ErrSignatureInvalid{
			return "", fmt.Errorf("invalid token signature: %s", err)
		}
		return "", fmt.Errorf("error parsing token: %s",err)
	}

	// Checking if claims are valid
	if !token.Valid{
		return "", fmt.Errorf("invalid token: %s", err)
	}

	// Extracting user role from claims
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return "", fmt.Errorf("could not extract user claims:%s", err)
	}
	return claims.Role, nil
}