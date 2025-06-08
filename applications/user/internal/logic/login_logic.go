package logic

import (
	"context"
	"errors"

	"mcp/applications/auth/pb/auth"
	"mcp/applications/user/internal/svc"
	"mcp/applications/user/pb/user"
	"mcp/core/enums"
	"mcp/core/util/userutil"
	"mcp/dao/query"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginRequest) (*user.LoginResponse, error) {
	if in.Username == "" || in.Password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "用户名和密码不能为空")
	}

	qu := query.Q.User
	u, err := qu.WithContext(l.ctx).Where(qu.Name.Eq(in.Username)).Take()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "用户不存在")
		}
		logx.Error(err)
		return nil, status.Errorf(codes.Internal, "dberror")
	}

	if !userutil.VerifyPassword(u.PasswordHash, in.Password) {
		return nil, status.Error(codes.Unauthenticated, "密码错误")
	}

	accessTokenResult, err := l.svcCtx.AuthService.IssueAccessToken(l.ctx, &auth.IssueAccessTokenRequest{
		UserId:      uint64(u.ID),
		Permissions: []enums.Permission{},
	})

	if err != nil {
		logx.Errorf("failed to issue access token: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to issue access token")
	}

	refreshTokenResult, err := l.svcCtx.AuthService.IssueRefreshToken(l.ctx, &auth.IssueRefreshTokenRequest{
		UserId:      uint64(u.ID),
		Permissions: []enums.Permission{},
	})

	if err != nil {
		logx.Errorf("failed to issue refresh token: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to issue refresh token")
	}

	return &user.LoginResponse{
		AccessToken:  accessTokenResult.AccessToken,
		RefreshToken: refreshTokenResult.RefreshToken,
		UserId:       uint64(u.ID),
	}, nil
}
