package http

import (
	"context"
	"net/http"
	"sync"

	pkgdb "github.com/kaisawind/skopeoui/pkg/db"
	"github.com/kaisawind/skopeoui/pkg/pb"
	pkgsvc "github.com/kaisawind/skopeoui/pkg/service"
	"github.com/sirupsen/logrus"
)

type CreateTaskCallback func(ctx context.Context, task *pb.Task) (err error)
type DeleteTaskCallback func(ctx context.Context, task *pb.Task) (err error)
type UpdateTaskCallback func(ctx context.Context, task *pb.Task) (err error)

type Service struct {
	opts       *pkgsvc.Options
	httpServer *http.Server
	db         pkgdb.IDB
	callbacks  sync.Map
}

func NewService(db pkgdb.IDB, opt ...pkgsvc.Option) pkgsvc.IService {
	opts := pkgsvc.NewOptions().SetAddress(":8080")
	for _, o := range opt {
		o.Apply(opts)
	}
	return &Service{
		opts: opts,
		db:   db,
	}
}

func (s *Service) Close() (err error) {
	if s.httpServer != nil {
		s.httpServer.Shutdown(context.Background())
	}
	return
}

func (s *Service) Start() (err error) {
	go s.Serve()
	return
}

func (s *Service) RegisterCreateTaskCallback(cb CreateTaskCallback) {
	s.callbacks.Store(cb, cb)
}

func (s *Service) RegisterDeleteTaskCallback(cb DeleteTaskCallback) {
	s.callbacks.Store(cb, cb)
}

func (s *Service) RegisterUpdateTaskCallback(cb UpdateTaskCallback) {
	s.callbacks.Store(cb, cb)
}

func (s *Service) Serve() (err error) {
	s.opts.Address, err = pkgsvc.GetListeningAddress(s.opts.Address)
	if err != nil {
		logrus.WithError(err).Errorln("failed to get listening address", "error", err)
		return
	}
	s.httpServer = &http.Server{
		Addr:    s.opts.Address,
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

func (s *Service) Address() string {
	return s.opts.Address
}
