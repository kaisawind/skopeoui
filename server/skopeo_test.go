package server

import (
	"context"
	"testing"

	"github.com/kaisawind/skopeoui/pkg/pb"
)

func TestSkopeoTask(t *testing.T) {
	ctx := context.Background()
	task := &pb.Task{
		Source:      "m.daocloud.io/docker.io/nginx:alpine",
		Destination: "192.168.1.118:5000/nginx:alpine",
	}
	err := skopeoTask(ctx, task.Source, task.Destination, func(txt string) {
		t.Logf("skopeoTask output: %s", txt)
	})
	if err != nil {
		t.Fatalf("skopeoTask failed: %v", err)
	}
}
