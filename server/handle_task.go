package server

import (
	"context"
	"encoding/json/v2"
	"net/http"
	"strconv"

	"github.com/kaisawind/skopeoui/pkg/pb"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func (s *Server) ServeTaskMux(mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("POST /v1/task", s.CreateTask)
	mux.HandleFunc("DELETE /v1/task", s.DeleteTask)
	mux.HandleFunc("PUT /v1/task", s.UpdateTask)
	mux.HandleFunc("GET /v1/task", s.GetTask)
	mux.HandleFunc("GET /v1/tasks", s.ListTask)
	return mux
}

func (s *Server) CreateTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	t := &pb.Task{}
	err := json.UnmarshalRead(r.Body, t)
	if err != nil {
		HttpError(rw, http.StatusBadRequest, err.Error())
		return
	}
	err = s.db.CreateTask(ctx, t)
	if err != nil {
		HttpError(rw, http.StatusBadRequest, err.Error())
		return
	}
	eid, err := s.crons.AddFunc(t.Cron, func() {
		s.taskJob(context.Background(), t)
	})
	if err != nil {
		logrus.WithError(err).Error("add cron job failed")
		HttpError(rw, http.StatusInternalServerError, err.Error())
		return
	}
	// store task id and cron job id
	s.tasks.Store(t.Id, eid)
	HttpResponse(rw, http.StatusOK, t)
}

func (s *Server) DeleteTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sid := r.URL.Query().Get("id")
	if sid == "" {
		HttpError(rw, http.StatusBadRequest, "id is empty")
		return
	}
	iid, _ := strconv.Atoi(sid)
	t, err := s.db.GetTask(ctx, int32(iid))
	if err != nil {
		HttpError(rw, http.StatusBadRequest, err.Error())
		return
	}
	err = s.db.DeleteTask(ctx, int32(iid))
	if err != nil {
		HttpError(rw, http.StatusBadRequest, err.Error())
		return
	}
	// get cron job id
	eid, ok := s.tasks.Load(t.Id)
	if ok {
		// delete cron job
		s.crons.Remove(eid.(cron.EntryID))
		// delete task id and cron job id
		s.tasks.Delete(t.Id)
	}
	HttpResponse(rw, http.StatusOK, "ok")
}

func (s *Server) UpdateTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	t := &pb.Task{}
	err := json.UnmarshalRead(r.Body, t)
	if err != nil {
		HttpError(rw, http.StatusBadRequest, err.Error())
		return
	}
	err = s.db.UpdateTask(ctx, t)
	if err != nil {
		HttpError(rw, http.StatusBadRequest, err.Error())
		return
	}
	eid, ok := s.tasks.Load(t.Id)
	if ok {
		s.crons.Remove(eid.(cron.EntryID))
	}
	eid, err = s.crons.AddFunc(t.Cron, func() {
		s.taskJob(context.Background(), t)
	})
	if err != nil {
		logrus.WithError(err).Error("add cron job failed")
		HttpError(rw, http.StatusInternalServerError, err.Error())
		return
	}
	// store task id and cron job id
	s.tasks.Store(t.Id, eid)
	HttpResponse(rw, http.StatusOK, t)
}

func (s *Server) GetTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sid := r.URL.Query().Get("id")
	if sid == "" {
		HttpError(rw, http.StatusBadRequest, "id is empty")
		return
	}
	iid, _ := strconv.Atoi(sid)
	t, err := s.db.GetTask(ctx, int32(iid))
	if err != nil {
		HttpError(rw, http.StatusBadRequest, err.Error())
		return
	}
	HttpResponse(rw, http.StatusOK, t)
}

func (s *Server) ListTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	iskip := 0
	sskip := r.URL.Query().Get("skip")
	if sskip != "" {
		iskip, _ = strconv.Atoi(sskip)
	}
	ilimit := 0
	slimit := r.URL.Query().Get("limit")
	if slimit != "" {
		ilimit, _ = strconv.Atoi(slimit)
	}
	ccount, items, err := s.db.ListTask(ctx, int32(iskip), int32(ilimit))
	if err != nil {
		HttpError(rw, http.StatusBadRequest, err.Error())
		return
	}
	out := map[string]any{
		"count": ccount,
		"items": items,
	}
	HttpResponse(rw, http.StatusOK, out)
}
