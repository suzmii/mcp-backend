package logic

import (
	"context"

	"mcp/applications/mcp/internal/svc"
	"mcp/applications/mcp/pb/mcp"
	"mcp/dao/query"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetSessionListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSessionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSessionListLogic {
	return &GetSessionListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSessionListLogic) GetSessionList(in *mcp.GetSessionListRequest) (*mcp.GetSessionListResponse, error) {
	if in.Page < 1 || in.PageSize < 1 || in.PageSize > 100 {
		return nil, status.Error(codes.InvalidArgument, "bad page or page size")
	}

	qs := query.Q.Session
	sessions, err := qs.WithContext(l.ctx).Where(qs.UserID.Eq(in.UserId)).Offset(int(in.PageSize) * int(in.Page-1)).Limit(int(in.PageSize)).Find()
	if err != nil {
		logx.Errorf("获取session列表失败: %v", err)
		return nil, status.Error(codes.Internal, "dberror")
	}

	rsp := &mcp.GetSessionListResponse{}

	for _, v := range sessions {
		rsp.Sessions = append(rsp.Sessions, &mcp.SessionInfo{
			SessionUuid: v.UUID,
			SessionHint: v.Hint,
		})
	}

	return rsp, nil
}
