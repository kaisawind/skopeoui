package http

import "net/http"

func (s *Service) WebMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux = ServeWebMux(mux)
	mux = s.ServeLogMux(mux)
	mux = s.ServeTaskMux(mux)
	return mux
}
