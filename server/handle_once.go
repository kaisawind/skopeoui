package server

import (
	"context"
	"encoding/json/v2"
	"net/http"
	"strings"
	"time"

	"github.com/kaisawind/skopeoui/pkg/pb"
	"github.com/rs/xid"
)

func (s *Server) ServeOnceMux(mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("POST /v1/once", s.CreateOnce)
	mux.HandleFunc("POST /v1/once/log", s.GetOnceLog)
	mux.HandleFunc("DELETE /v1/once", s.DeleteOnce)
	mux.HandleFunc("GET /v1/onces", s.ListOnce)
	return mux
}

func (s *Server) CreateOnce(rw http.ResponseWriter, r *http.Request) {
	t := &pb.Once{}
	err := json.UnmarshalRead(r.Body, t)
	if err != nil {
		HttpError(rw, http.StatusBadRequest, err.Error())
		return
	}
	go s.onceJob(context.Background(), t)
	HttpResponse(rw, http.StatusOK, t)
}

func (s *Server) GetOnceLog(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sid := r.URL.Query().Get("id")
	if sid == "" {
		HttpError(rw, http.StatusBadRequest, "id is empty")
		return
	}
	_, ok := s.onces.Load(sid)
	if !ok {
		HttpError(rw, http.StatusNotFound, "job not found")
		return
	}
	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")

	if f, ok := rw.(http.Flusher); ok {
		f.Flush()
	} else {
		http.Error(rw, "Streaming unsupported", http.StatusInternalServerError)
		return
	}
	s.rds.Store(sid+xid.New().String(), rw)
	timer := time.NewTicker(10 * time.Hour)
	defer timer.Stop()
	done := make(chan struct{})
	defer close(done)
	s.dones.Store(sid, done)
	defer s.dones.Delete(sid)
	select {
	case <-timer.C:
	case <-done:
	case <-ctx.Done():
	}
}

func (s *Server) DeleteOnce(rw http.ResponseWriter, r *http.Request) {
	sid := r.URL.Query().Get("id")
	if sid == "" {
		HttpError(rw, http.StatusBadRequest, "id is empty")
		return
	}
	s.onces.Delete(sid)
	s.rds.Range(func(key, value any) bool {
		if strings.HasPrefix(key.(string), sid) {
			s.rds.Delete(key)
		}
		return true
	})
	done, ok := s.dones.Load(sid)
	if ok {
		done.(chan struct{}) <- struct{}{}
	}
	HttpResponse(rw, http.StatusOK, nil)
}

func (s *Server) ListOnce(rw http.ResponseWriter, r *http.Request) {
	values := []map[string]string{}
	s.onces.Range(func(key, value any) bool {
		once, ok := value.(*pb.Once)
		if ok {
			values = append(values, map[string]string{
				"id":     key.(string),
				"source": once.Source,
				"dest":   once.Destination,
			})
		}
		return true
	})
	HttpResponse(rw, http.StatusOK, values)
}
