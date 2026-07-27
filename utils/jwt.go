package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"trashcan/global"
)

var (
	TokenExpired     = errors.New("token已过期")
	TokenNotValidYet = errors.New("token尚未激活")
	TokenMalformed   = errors.New("token格式错误")
	TokenInvalid     = errors.New("token无效")
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT token
func GenerateToken(userID uint, username string, isAdmin bool) (string, error) {
	secret := global.CONFIG.JWTConfig.Secret
	expireHours := global.CONFIG.JWTConfig.ExpireHours
	if secret == "" {
		secret = "default-secret-key"
	}
	if expireHours == 0 {
		expireHours = 24
	}

	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(expireHours) * time.Hour)

	claims := Claims{
		UserID:   userID,
		Username: username,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(nowTime),
			NotBefore: jwt.NewNumericDate(nowTime),
			Issuer:    "trashcan-system",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(secret))
	return token, err
}

// ParseToken 解析并验证token
func ParseToken(tokenString string) (*Claims, error) {
	secret := global.CONFIG.JWTConfig.Secret
	if secret == "" {
		secret = "default-secret-key"
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, TokenMalformed
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, TokenExpired
		} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, TokenNotValidYet
		} else {
			return nil, TokenInvalid
		}
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, TokenInvalid
}
