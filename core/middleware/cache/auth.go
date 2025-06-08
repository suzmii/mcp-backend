package cache

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth Redis

func NewAuth(conf Config) *Auth {
	return (*Auth)(MustNewRedis(conf))
}

func tokenKey(subject string, uid uint64) string {
	return "token:" + subject + ":" + strconv.FormatUint(uid, 10)
}

func (a *Auth) SetTokenSecret(ctx context.Context, subject string, uid uint64, secret []byte, exp time.Duration) error {
	err := a.rds.Set(ctx, tokenKey(subject, uid), secret, exp).Err()
	if err != nil {
		return err
	}
	return nil
}

func (a *Auth) DelTokenSecret(ctx context.Context, subject string, uid uint64) error {
	err := a.rds.Del(ctx, tokenKey(subject, uid)).Err()
	if err != nil {
		return err
	}
	return nil
}

func (a *Auth) GetTokenSecret(ctx context.Context, subject string, uid uint64) ([]byte, error) {
	result := a.rds.Get(ctx, tokenKey(subject, uid))
	if err := result.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, status.Error(codes.Unauthenticated, "invaild token")
		}
		return nil, err
	}
	secret, err := result.Bytes()
	return secret, err
}
