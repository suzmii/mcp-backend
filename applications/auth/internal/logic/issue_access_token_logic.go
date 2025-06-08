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

type IssueAccessTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIssueAccessTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IssueAccessTokenLogic {
	return &IssueAccessTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IssueAccessTokenLogic) IssueAccessToken(in *auth.IssueAccessTokenRequest) (*auth.IssueAccessTokenResponse, error) {
	now := time.Now()
	secret := jwtutil.GenerateSecret(in.UserId, now.Unix())

	token, err := jwtutil.GenerateToken(secret[:], in.UserId, "access", in.Permissions, l.svcCtx.Config.AccessExp)
	if err != nil {
		l.Logger.Errorf("Failed to generate access token: %v", err)
		return nil, status.Errorf(codes.Internal, "生成访问令牌失败")
	}

	err = l.svcCtx.Auth.SetTokenSecret(l.ctx, "access", in.UserId, secret[:], l.svcCtx.Config.AccessExp)
	if err != nil {
		l.Logger.Errorf("Failed to set token secret: %v", err)
		return nil, status.Errorf(codes.Internal, "保存令牌密钥失败")
	}

	return &auth.IssueAccessTokenResponse{
		AccessToken: token,
	}, nil
}
