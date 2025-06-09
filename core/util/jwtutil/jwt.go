package jwtutil

import (
	"crypto/sha256"
	"mcp/core/enums"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var signMethod = jwt.SigningMethodHS256

type Claims struct {
	UserID uint64             `json:"uid"`
	Perms  []enums.Permission `json:"permissions"`
	jwt.RegisteredClaims
}

func GenerateSecret(uid uint64, timestamp int64) [32]byte {
	return sha256.Sum256([]byte(
		strconv.FormatUint(uid, 10) +
			"-salt-copyright-suzmii-" +
			strconv.FormatInt(timestamp, 10)))
}

// GenerateToken 生成带签名的 JWT 字符串
func GenerateToken(secret []byte, uid uint64, subject string, perms []enums.Permission, exp time.Duration) (string, error) {
	now := time.Now()

	claims := &Claims{
		UserID: uid,
		Perms:  perms,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "backend server", // (optional) identify auth service ID
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(exp)),
			Subject:   subject,
		},
	}

	token := jwt.NewWithClaims(signMethod, claims)

	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", errors.Wrap(err, "failed to sign token")
	}

	return signedToken, nil
}

// ParseToken 解析并验证 JWT，如果有效则返回 Claims
func ParseToken(tokenStr, subject string, secretGetter func(userId uint64) ([]byte, error)) (*Claims, error) {
	claims := &Claims{}

	_, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
		if token.Method != signMethod {
			return nil, errors.New("unexpected signing method")
		}
		if claims.Subject != subject {
			return nil, errors.New("subject mismatch")
		}
		if claims.IssuedAt == nil || claims.ExpiresAt == nil {
			return nil, errors.New("missing issued or expiry time")
		}
		if time.Now().After(claims.ExpiresAt.Time) {
			return nil, errors.New("token is expired")
		}
		return secretGetter(claims.UserID)
	})

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "failed to parse or verify jwt")
	}

	return claims, nil
}
