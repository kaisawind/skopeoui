package server

import (
	"fmt"
	"sync"

	"github.com/kaisawind/skopeoui/internal/db"
	"github.com/kaisawind/skopeoui/internal/http"
	pkgdb "github.com/kaisawind/skopeoui/pkg/db"
	pkgsvc "github.com/kaisawind/skopeoui/pkg/service"
	"github.com/sirupsen/logrus"
)

type Server struct {
	opts        *Options
	quit        chan struct{}
	db          pkgdb.IDB
	httpService pkgsvc.IService
}

func NewServer(ops ...*Options) *Server {
	opts := NewOptions()
	for _, o := range ops {
		o.Apply(opts)
	}
	return &Server{
		opts: opts,
		quit: make(chan struct{}),
	}
}

func (s *Server) Start() error {
	go s.Serve()
	return nil
}

func (s *Server) Serve() (err error) {
	err = s.initdb()
	if err != nil {
		logrus.WithError(err).Error("init db failed")
		return
	}
	logrus.Info("db initialized")

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

func (s *Server) initdb() (err error) {
	s.db, err = db.New()
	if err != nil {
		return
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
	<-s.quit
	s.httpService.Close()
	return
}
