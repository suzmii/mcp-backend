package logic

import (
	"context"
	"errors"

	"mcp/applications/mcp/internal/svc"
	"mcp/applications/mcp/pb/mcp"
	"mcp/dao/models"
	"mcp/dao/query"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type AppendMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAppendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppendMessageLogic {
	return &AppendMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AppendMessageLogic) AppendMessage(in *mcp.AppendMessageRequest) (*mcp.AppendMessageResponse, error) {
	qs := query.Q.Session
	s, err := qs.WithContext(l.ctx).Where(qs.UUID.Eq(in.SessionUuid)).Take()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "session not found")
		}
		logx.Errorf("获取sesson失败: %v", err)
		return nil, status.Error(codes.Internal, "dberror")
	}

	msgs := make([]*models.Message, 0, len(in.Messages))

	for _, v := range in.Messages {
		msgs = append(msgs, &models.Message{
			SessionID: s.ID,
			Role:      v.Role,
			Content:   v.Content,
		})
	}

	qm := query.Q.Message
	err = qm.WithContext(l.ctx).Create(msgs...)
	if err != nil {
		logx.Errorf("创建消息失败: %v", err)
		return nil, status.Error(codes.Internal, "dberror")
	}
	return &mcp.AppendMessageResponse{}, nil
}
