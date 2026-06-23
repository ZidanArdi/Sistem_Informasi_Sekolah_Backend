package helpers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"
)

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Exp    int64  `json:"exp"`
}

func GenerateToken(userID uint, email string, role string) (string, error) {
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		Exp:    time.Now().Add(24 * time.Hour).Unix(),
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	unsignedToken := base64.RawURLEncoding.EncodeToString(headerJSON) + "." + base64.RawURLEncoding.EncodeToString(claimsJSON)
	signature := signJWT(unsignedToken)

	return unsignedToken + "." + signature, nil
}

func ParseToken(token string) (*JWTClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("token tidak valid")
	}

	unsignedToken := parts[0] + "." + parts[1]
	expectedSignature := signJWT(unsignedToken)
	if !hmac.Equal([]byte(expectedSignature), []byte(parts[2])) {
		return nil, errors.New("signature token tidak valid")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var claims JWTClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, err
	}

	if claims.Exp < time.Now().Unix() {
		return nil, errors.New("token sudah kedaluwarsa")
	}

	return &claims, nil
}

func signJWT(unsignedToken string) string {
	hash := hmac.New(sha256.New, []byte(jwtSecret()))
	hash.Write([]byte(unsignedToken))
	return base64.RawURLEncoding.EncodeToString(hash.Sum(nil))
}

func jwtSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "sistem-informasi-sekolah-dev-secret"
	}
	return secret
}
