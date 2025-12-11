package server

import "net/http"

func (s *Server) WebMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux = ServeWebMux(mux)
	mux = s.ServeLogMux(mux)
	mux = s.ServeTaskMux(mux)
	return mux
}
