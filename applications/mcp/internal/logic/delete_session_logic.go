package logic

import (
	"context"
	"errors"

	"mcp/applications/mcp/internal/svc"
	"mcp/applications/mcp/pb/mcp"
	"mcp/dao/query"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type DeleteSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSessionLogic {
	return &DeleteSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteSessionLogic) DeleteSession(in *mcp.DeleteSessionRequest) (*mcp.DeleteSessionResponse, error) {
	qs := query.Q.Session
	s, err := qs.WithContext(l.ctx).Where(qs.UUID.Eq(in.SessionUuid)).Take()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "session not found")
		}
		logx.Errorf("获取session失败: %v", err)
		return nil, status.Error(codes.Internal, "dberror")
	}

	if s.UserID != in.UserId {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	_, err = qs.WithContext(l.ctx).Where(qs.ID.Eq(s.ID)).Delete()
	if err != nil {
		logx.Errorf("failed to delete session: %v", err)
		return nil, status.Error(codes.Internal, "dberror")
	}

	return &mcp.DeleteSessionResponse{}, nil
}
