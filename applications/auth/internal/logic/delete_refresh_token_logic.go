package logic

import (
	"context"

	"mcp/applications/auth/internal/svc"
	"mcp/applications/auth/pb/auth"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeleteRefreshTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRefreshTokenLogic {
	return &DeleteRefreshTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteRefreshTokenLogic) DeleteRefreshToken(in *auth.DeleteRefreshTokenRequest) (*auth.DeleteRefreshTokenResponse, error) {
	verify := NewVerifyRefreshTokenLogic(l.ctx, l.svcCtx)
	result, err := verify.VerifyRefreshToken(&auth.VerifyRefreshTokenRequest{
		RefreshToken: in.RefreshToken,
	})
	if err != nil {
		l.Logger.Errorf("Failed to verify refresh token: %v", err)
		return nil, err
	}

	err = l.svcCtx.Auth.DelTokenSecret(l.ctx, "refresh", result.Claims.UserId)
	if err != nil {
		l.Logger.Errorf("Failed to delete refresh token secret: %v", err)
		return nil, status.Errorf(codes.Internal, "删除刷新令牌失败")
	}

	return &auth.DeleteRefreshTokenResponse{}, nil
}
