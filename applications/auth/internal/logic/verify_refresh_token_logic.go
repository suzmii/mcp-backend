package logic

import (
	"context"

	"mcp/applications/auth/internal/svc"
	"mcp/applications/auth/pb/auth"
	"mcp/core/util/jwtutil"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type VerifyRefreshTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyRefreshTokenLogic {
	return &VerifyRefreshTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *VerifyRefreshTokenLogic) VerifyRefreshToken(in *auth.VerifyRefreshTokenRequest) (*auth.VerifyRefreshTokenResponse, error) {
	claims, err := jwtutil.ParseToken(in.RefreshToken, "refresh", func(userId uint64) ([]byte, error) {
		return l.svcCtx.Auth.GetTokenSecret(l.ctx, "refresh", userId)
	})
	if err != nil {
		l.Logger.Errorf("Failed to parse refresh token: %v", err)
		return nil, status.Errorf(codes.Unauthenticated, "无效的刷新令牌")
	}

	return &auth.VerifyRefreshTokenResponse{
		Claims: &auth.Claims{
			UserId:      claims.UserID,
			Permissions: claims.Perms,
		},
	}, nil
}
