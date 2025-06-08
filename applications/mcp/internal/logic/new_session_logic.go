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

type NewSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNewSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NewSessionLogic {
	return &NewSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NewSessionLogic) NewSession(in *mcp.NewSessionRequest) (*mcp.NewSessionResponse, error) {
	uuid := uuid.New().String()
	session := models.Session{
		UUID:   uuid,
		UserID: in.UserId,
	}

	qu := query.Q.Session

	err := qu.WithContext(l.ctx).Create(&session)
	if err != nil {
		logx.Errorf("创建session失败: %v", err)
		return nil, status.Error(codes.Internal, "dberror")
	}
	return &mcp.NewSessionResponse{
		SessionUuid: uuid,
	}, nil
}
