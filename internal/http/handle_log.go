package http

import (
	"net/http"
	"strconv"
)

func (s *Service) ServeLogMux(mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("DELETE /v1/log", s.DeleteLog)
	mux.HandleFunc("DELETE /v1/logs/task", s.DeleteLogByTaskId)
	mux.HandleFunc("GET /v1/log", s.GetLog)
	mux.HandleFunc("GET /v1/logs", s.ListLog)
	mux.HandleFunc("GET /v1/logs/task", s.ListLogByTaskId)
	return mux
}

func (s *Service) DeleteLog(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sid := r.URL.Query().Get("id")
	if sid == "" {
		HttpError(rw, http.StatusBadRequest, "id is empty")
		return
	}
	iid, _ := strconv.Atoi(sid)
	err := s.db.DeleteLog(ctx, int32(iid))
	if err != nil {
		HttpError(rw, http.StatusBadRequest, err.Error())
		return
	}
	HttpResponse(rw, http.StatusOK, "ok")
}

func (s *Service) DeleteLogByTaskId(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sid := r.URL.Query().Get("id")
	if sid == "" {
		HttpError(rw, http.StatusBadRequest, "id is empty")
		return
	}
	iid, _ := strconv.Atoi(sid)
	err := s.db.DeleteLogByTaskId(ctx, int32(iid))
	if err != nil {
		HttpError(rw, http.StatusBadRequest, err.Error())
		return
	}
	HttpResponse(rw, http.StatusOK, "ok")
}

func (s *Service) GetLog(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sid := r.URL.Query().Get("id")
	if sid == "" {
		HttpError(rw, http.StatusBadRequest, "id is empty")
		return
	}
	iid, _ := strconv.Atoi(sid)
	t, err := s.db.GetLog(ctx, int32(iid))
	if err != nil {
		HttpError(rw, http.StatusBadRequest, err.Error())
		return
	}
	HttpResponse(rw, http.StatusOK, t)
}

func (s *Service) ListLog(rw http.ResponseWriter, r *http.Request) {
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
	ccount, items, err := s.db.ListLog(ctx, int32(iskip), int32(ilimit))
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

func (s *Service) ListLogByTaskId(rw http.ResponseWriter, r *http.Request) {
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
	sid := r.URL.Query().Get("taskId")
	if sid == "" {
		HttpError(rw, http.StatusBadRequest, "taskId is empty")
		return
	}
	iid, _ := strconv.Atoi(sid)
	ccount, items, err := s.db.ListLogByTaskId(ctx, int32(iid), int32(iskip), int32(ilimit))
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
