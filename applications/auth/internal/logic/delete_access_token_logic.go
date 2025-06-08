package logic

import (
	"context"

	"mcp/applications/auth/internal/svc"
	"mcp/applications/auth/pb/auth"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeleteAccessTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteAccessTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAccessTokenLogic {
	return &DeleteAccessTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteAccessTokenLogic) DeleteAccessToken(in *auth.DeleteAccessTokenRequest) (*auth.DeleteAccessTokenResponse, error) {
	verify := NewVerifyAccessTokenLogic(l.ctx, l.svcCtx)
	result, err := verify.VerifyAccessToken(&auth.VerifyAccessTokenRequest{
		AccessToken: in.AccessToken,
	})
	if err != nil {
		l.Logger.Errorf("Failed to verify access token: %v", err)
		return nil, err
	}

	err = l.svcCtx.Auth.DelTokenSecret(l.ctx, "access", result.Claims.UserId)
	if err != nil {
		l.Logger.Errorf("Failed to delete access token secret: %v", err)
		return nil, status.Errorf(codes.Internal, "删除访问令牌失败")
	}

	return &auth.DeleteAccessTokenResponse{}, nil
}
