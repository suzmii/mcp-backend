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

type GetSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSessionLogic {
	return &GetSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSessionLogic) GetSession(in *mcp.GetSessionRequest) (*mcp.GetSessionResponse, error) {
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

	qm := query.Q.Message

	msgs, err := qm.WithContext(l.ctx).Order(qm.ID).Where(qm.SessionID.Eq(s.ID)).Find()
	if err != nil {
		logx.Errorf("获取消息失败: %v", err)
		return nil, status.Error(codes.Internal, "dberror")
	}

	rsp := &mcp.GetSessionResponse{
		Session: &mcp.SessionInfo{
			SessionUuid: s.UUID,
			SessionHint: s.Hint,
		},
	}

	for _, v := range msgs {
		rsp.Messages = append(rsp.Messages, &mcp.Message{
			Role:    v.Role,
			Content: v.Content,
		})
	}

	return rsp, nil
}
