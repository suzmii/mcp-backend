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

type VerifyAccessTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyAccessTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyAccessTokenLogic {
	return &VerifyAccessTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *VerifyAccessTokenLogic) VerifyAccessToken(in *auth.VerifyAccessTokenRequest) (*auth.VerifyAccessTokenResponse, error) {
	claims, err := jwtutil.ParseToken(in.AccessToken, "access", func(userId uint64) ([]byte, error) {
		return l.svcCtx.Auth.GetTokenSecret(l.ctx, "access", userId)
	})
	if err != nil {
		l.Logger.Errorf("Failed to parse access token: %v", err)
		return nil, status.Errorf(codes.Unauthenticated, "无效的访问令牌")
	}

	return &auth.VerifyAccessTokenResponse{
		Claims: &auth.Claims{
			UserId:      claims.UserID,
			Permissions: claims.Perms,
		},
	}, nil
}
