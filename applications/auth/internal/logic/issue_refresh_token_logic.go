package logic

import (
	"context"
	"time"

	"mcp/applications/auth/internal/svc"
	"mcp/applications/auth/pb/auth"
	"mcp/core/util/jwtutil"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IssueRefreshTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIssueRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IssueRefreshTokenLogic {
	return &IssueRefreshTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IssueRefreshTokenLogic) IssueRefreshToken(in *auth.IssueRefreshTokenRequest) (*auth.IssueRefreshTokenResponse, error) {
	now := time.Now()
	secret := jwtutil.GenerateSecret(in.UserId, now.Unix())

	token, err := jwtutil.GenerateToken(secret[:], in.UserId, "refresh", in.Permissions, l.svcCtx.Config.RefreshExp)
	if err != nil {
		l.Logger.Errorf("Failed to generate refresh token: %v", err)
		return nil, status.Errorf(codes.Internal, "生成刷新令牌失败")
	}

	err = l.svcCtx.Auth.SetTokenSecret(l.ctx, "refresh", in.UserId, secret[:], l.svcCtx.Config.RefreshExp)
	if err != nil {
		l.Logger.Errorf("Failed to set refresh token secret: %v", err)
		return nil, status.Errorf(codes.Internal, "保存刷新令牌密钥失败")
	}

	return &auth.IssueRefreshTokenResponse{
		RefreshToken: token,
	}, nil
}
