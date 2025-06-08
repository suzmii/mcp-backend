package logic

import (
	"context"
	"errors"

	"mcp/applications/user/internal/svc"
	"mcp/applications/user/pb/user"
	"mcp/core/util/userutil"
	"mcp/dao/models"
	"mcp/dao/query"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterRequest) (*user.RegisterResponse, error) {
	if in.Username == "" || in.Password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "用户名和密码不能为空")
	}

	if len(in.Password) < 6 {
		return nil, status.Errorf(codes.InvalidArgument, "密码长度不能少于6位")
	}

	u := &models.User{
		Name: in.Username,
	}

	err := query.Q.Transaction(func(tx *query.Query) error {
		qu := tx.User
		_, err := qu.WithContext(l.ctx).Where(qu.Name.Eq(in.Username)).Take()
		if err == nil {
			return status.Errorf(codes.AlreadyExists, "用户名已存在")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logx.Errorf("查询用户信息失败: %v", err)
			return status.Error(codes.Internal, "dberror")
		}

		hash, err := userutil.PasswordHash(in.Password)
		if err != nil {
			logx.Errorf("failed to generate password hash: %v", err)
			return status.Error(codes.Internal, "")
		}

		u.PasswordHash = hash

		err = qu.WithContext(l.ctx).Create(u)

		if err != nil {
			logx.Errorf("创建用户失败: %v", err)
			return status.Error(codes.Internal, "dberror")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &user.RegisterResponse{
		UserId: uint64(u.ID),
	}, nil
}
