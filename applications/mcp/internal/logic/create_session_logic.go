package logic

import (
	"context"

	"mcp/applications/mcp/internal/svc"
	"mcp/applications/mcp/pb/mcp"
	"mcp/dao/models"
	"mcp/dao/query"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSessionLogic {
	return &CreateSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateSessionLogic) CreateSession(in *mcp.CreateSessionRequest) (*mcp.CreateSessionResponse, error) {
	uuid := uuid.New().String()
	session := models.Session{
		UUID:   uuid,
		UserID: in.UserId,
		Hint:   in.Hint,
	}

	qu := query.Q.Session

	err := qu.WithContext(l.ctx).Create(&session)
	if err != nil {
		logx.Errorf("创建session失败: %v", err)
		return nil, status.Error(codes.Internal, "dberror")
	}
	return &mcp.CreateSessionResponse{}, nil
}
