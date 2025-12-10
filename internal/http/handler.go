package http

import "net/http"

func (s *Service) WebMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux = ServeMux(mux)
	return mux
}
