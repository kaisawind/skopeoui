package server

import (
	"context"
	"strings"
	"time"

	"github.com/kaisawind/skopeoui/pkg/pb"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func (s *Server) taskJob(ctx context.Context, task *pb.Task) {
	txts := []string{}
	e := skopeoTask(ctx, task.Source, task.Destination, func(txt string) {
		txts = append(txts, txt)
	})
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

func (s *Server) createTaskCallback(_ context.Context, task *pb.Task) (err error) {
	eid, err := s.crons.AddFunc(task.Cron, func() {
		s.taskJob(context.Background(), task)
	})
	if err != nil {
		logrus.WithError(err).Error("add cron job failed")
		return
	}
	// store task id and cron job id
	s.tasks.Store(task.Id, eid)
	return
}

func (s *Server) deleteTaskCallback(_ context.Context, task *pb.Task) (err error) {
	// get cron job id
	eid, ok := s.tasks.Load(task.Id)
	if !ok {
		logrus.WithField("task_id", task.Id).Warnln("task not found")
		return
	}
	// delete cron job
	s.crons.Remove(eid.(cron.EntryID))
	// delete task id and cron job id
	s.tasks.Delete(task.Id)
	return
}

func (s *Server) updateTaskCallback(_ context.Context, task *pb.Task) (err error) {
	// get cron job id
	eid, ok := s.tasks.Load(task.Id)
	if !ok {
		logrus.WithField("task_id", task.Id).Warnln("task not found")
		return
	}
	// update cron job
	s.crons.Remove(eid.(cron.EntryID))
	eid, err = s.crons.AddFunc(task.Cron, func() {
		s.taskJob(context.Background(), task)
	})
	if err != nil {
		logrus.WithError(err).Error("add cron job failed")
		return
	}
	// store task id and cron job id
	s.tasks.Store(task.Id, eid)
	return
}
