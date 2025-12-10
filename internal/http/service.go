package http

import (
	"context"
	"net/http"

	pkgsvc "github.com/kaisawind/skopeoui/pkg/service"
	"github.com/sirupsen/logrus"
)

type Service struct {
	opts       *pkgsvc.Options
	httpServer *http.Server
}

func NewService(opt ...pkgsvc.Option) pkgsvc.IService {
	opts := pkgsvc.NewOptions().SetAddress(":8080")
	for _, o := range opt {
		o.Apply(opts)
	}
	return &Service{
		opts: opts,
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
