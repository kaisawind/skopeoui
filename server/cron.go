package server

import (
	"context"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"net/http"
	"strings"
	"time"

	"github.com/kaisawind/skopeoui/pkg/pb"
	"github.com/sirupsen/logrus"
)

func (s *Server) taskJob(ctx context.Context, task *pb.Task) {
	txts := []string{}
	e := skopeoTask(ctx, task.Source, task.Destination, func(txt string) {
		txts = append(txts, txt)
	}, func() {})
	if e != nil {
		logrus.WithError(e).Error("skopeo copy failed")
		return
	}
	e = s.db.CreateLog(ctx, &pb.Log{
		TaskId: task.Id,
		Msg:    strings.Join(txts, "\n"),
		Time:   time.Now().Unix(),
	})
	if e != nil {
		logrus.WithError(e).Error("create log failed")
		return
	}
}

func (s *Server) onceJob(ctx context.Context, t *pb.Once) {
	hash := fnv.New128()
	hash.Write([]byte(t.Source + t.Destination))
	sum := hash.Sum(nil)
	onceKey := hex.EncodeToString(sum)
	s.onces.Store(onceKey, t)
	e := skopeoTask(ctx, t.Source, t.Destination, func(txt string) {
		s.rds.Range(func(key, value any) bool {
			if strings.HasPrefix(key.(string), onceKey) {
				rw, ok := value.(http.ResponseWriter)
				if ok {
					_, err := fmt.Fprintf(rw, "data: %s\n\n", txt)
					if err != nil {
						logrus.WithError(err).Error("write response failed")
						return true
					}
					if f, ok := rw.(http.Flusher); ok {
						f.Flush()
					}
				}
			}
			return true
		})
	}, func() {
		s.onces.Delete(onceKey)
		s.rds.Delete(onceKey)
		done, ok := s.dones.Load(onceKey)
		if ok {
			done.(chan struct{}) <- struct{}{}
		}
	})
	if e != nil {
		logrus.WithError(e).Error("skopeo copy failed")
		return
	}
}
