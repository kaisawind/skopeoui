package server

import (
	"context"
	"net/http"
	"sync"

	"github.com/kaisawind/skopeoui/internal/db"
	pkgdb "github.com/kaisawind/skopeoui/pkg/db"
	pkgsvc "github.com/kaisawind/skopeoui/pkg/service"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type Server struct {
	opts       *Options
	crons      *cron.Cron
	tasks      sync.Map // cron <-> task id
	onces      sync.Map // hash <-> once task
	rds        sync.Map //
	dones      sync.Map // hash <-> once task done
	db         pkgdb.IDB
	httpServer *http.Server
	quit       chan struct{}
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
	if s.httpServer != nil {
		s.httpServer.Shutdown(context.Background())
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
		eid, e := s.crons.AddFunc(t.Cron, func() {
			s.taskJob(context.Background(), t)
		})
		if e != nil {
			logrus.WithError(e).Error("add cron job failed")
			continue
		}
		// store task id and cron job id
		s.tasks.Store(t.Id, eid)
	}
	s.crons.Start()
	return
}

func (s *Server) doHttp() (err error) {
	naddr, err := pkgsvc.GetListeningAddress(s.opts.HttpAddress)
	if err != nil {
		return
	}
	s.httpServer = &http.Server{
		Addr:    naddr,
		Handler: Cross(s.WebMux()), // or your custom handler
	}
	logrus.Infoln("http service is started ...", "addr", s.httpServer.Addr)
	err = s.httpServer.ListenAndServe()
	if err != nil {
		logrus.WithError(err).Errorln("failed to http serve", "error", err)
		return err
	}
	logrus.Infoln("http server existed ...")
	return
}
