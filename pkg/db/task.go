package db

import (
	"context"

	"github.com/kaisawind/skopeoui/pkg/pb"
)

type ITask interface {
	CreateTask(ctx context.Context, in *pb.Task) (err error)
	DeleteTask(ctx context.Context, id int32) (err error)
	UpdateTask(ctx context.Context, in *pb.Task) (err error)
	GetTask(ctx context.Context, id int32) (out *pb.Task, err error)
	ListTask(ctx context.Context, skip, limit int32) (count int32, out []*pb.Task, err error)
}
