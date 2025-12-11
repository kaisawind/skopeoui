package server

import (
	"context"
	"fmt"
	"sync"

	"github.com/kaisawind/skopeoui/internal/db"
	"github.com/kaisawind/skopeoui/internal/http"
	pkgdb "github.com/kaisawind/skopeoui/pkg/db"
	pkgsvc "github.com/kaisawind/skopeoui/pkg/service"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type Server struct {
	opts        *Options
	crons       *cron.Cron
	tasks       sync.Map
	db          pkgdb.IDB
	httpService pkgsvc.IService
	quit        chan struct{}
}

func NewServer(ops ...*Options) *Server {
	opts := NewOptions()
	for _, o := range ops {
		o.Apply(opts)
	}
	return &Server{
		opts:  opts,
		crons: cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger))),
		quit:  make(chan struct{}),
	}
}

func (s *Server) Start() error {
	go s.Serve()
	return nil
}

func (s *Server) Serve() (err error) {
	err = s.initDB()
	if err != nil {
		logrus.WithError(err).Error("init db failed")
		return
	}
	logrus.Info("db initialized")

	// do cron job
	err = s.doCron()
	if err != nil {
		logrus.WithError(err).Error("do cron error")
		return
	}
	logrus.Info("cron job done")

	var wg sync.WaitGroup
	wg.Go(func() {
		logrus.Info("http service started")
		err = s.doHttp()
		if err != nil {
			logrus.WithError(err).Error("do http error")
			return
		}
	})
	wg.Wait()
	return
}

func (s *Server) Close() {
	if s.httpService != nil {
		s.httpService.Close()
	}
	if s.db != nil {
		s.db.Close()
	}
}

func (s *Server) initDB() (err error) {
	s.db, err = db.New()
	if err != nil {
		return
	}
	return
}

func (s *Server) doCron() (err error) {
	ctx := context.Background()
	count, tasks, err := s.db.ListTask(ctx, 0, 0)
	if err != nil {
		return
	}
	logrus.Infof("list %d tasks", count)
	for _, t := range tasks {
		err = s.createTaskCallback(ctx, t)
		if err != nil {
			logrus.WithError(err).Errorf("create task callback failed, task: %v", t)
		}
	}
	return
}

func (s *Server) doHttp() (err error) {
	opts := pkgsvc.NewOptions().SetAddress(s.opts.HttpAddress)
	s.httpService = http.NewService(s.db, opts)
	if s.httpService == nil {
		err = fmt.Errorf("http service created failed")
		return
	}
	err = s.httpService.Start()
	if err != nil {
		return
	}
	httpSvc := s.httpService.(*http.Service)
	httpSvc.RegisterCreateTaskCallback(s.createTaskCallback)
	httpSvc.RegisterDeleteTaskCallback(s.deleteTaskCallback)
	httpSvc.RegisterUpdateTaskCallback(s.updateTaskCallback)
	<-s.quit
	s.httpService.Close()
	return
}
