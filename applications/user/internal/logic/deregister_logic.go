package logic

import (
	"context"
	"errors"

	"mcp/applications/user/internal/svc"
	"mcp/applications/user/pb/user"
	"mcp/dao/query"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"mcp/core/util/userutil"
)

type DeregisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeregisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeregisterLogic {
	return &DeregisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户注销
func (l *DeregisterLogic) Deregister(in *user.DeregisterRequest) (*user.DeregisterResponse, error) {
	err := query.Q.Transaction(func(tx *query.Query) error {
		qu := tx.User
		u, err := qu.WithContext(l.ctx).Where(qu.Name.Eq(in.Username)).First()
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return status.Error(codes.NotFound, "用户不存在")
			}
			logx.Error(err)
			return status.Error(codes.Internal, "数据库错误")
		}
		if !userutil.VerifyPassword(u.PasswordHash, in.Password) {
			return status.Error(codes.Unauthenticated, "密码错误")
		}
		_, err = qu.WithContext(l.ctx).Where(qu.ID.Eq(u.ID)).Delete()
		if err != nil {
			logx.Error(err)
			return status.Error(codes.Internal, "数据库错误")
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &user.DeregisterResponse{}, nil
}
