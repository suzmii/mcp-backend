package logic

import (
	"context"

	"mcp/applications/auth/pb/auth"
	"mcp/applications/user/internal/svc"
	"mcp/applications/user/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LogoutLogic) Logout(in *user.LogoutRequest) (*user.LogoutResponse, error) {
	_, err := l.svcCtx.AuthService.DeleteAccessToken(l.ctx, &auth.DeleteAccessTokenRequest{
		AccessToken: in.AccessToken,
	})
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid access token")
	}
	_, err = l.svcCtx.AuthService.DeleteRefreshToken(l.ctx, &auth.DeleteRefreshTokenRequest{
		RefreshToken: in.RefreshToken,
	})
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
	}
	return &user.LogoutResponse{}, nil
}
