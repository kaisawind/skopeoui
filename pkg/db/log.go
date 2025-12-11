package db

import (
	"context"

	"github.com/kaisawind/skopeoui/pkg/pb"
)

type ILog interface {
	CreateLog(ctx context.Context, in *pb.Log) (err error)
	DeleteLog(ctx context.Context, id int32) (err error)
	DeleteLogByTaskId(ctx context.Context, taskId int32) (err error)
	UpdateLog(ctx context.Context, in *pb.Log) (err error)
	GetLog(ctx context.Context, id int32) (out *pb.Log, err error)
	ListLog(ctx context.Context, skip, limit int32) (count int32, out []*pb.Log, err error)
	ListLogByTaskId(ctx context.Context, taskId int32, skip, limit int32) (count int32, out []*pb.Log, err error)
}
