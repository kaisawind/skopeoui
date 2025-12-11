package http

import (
	"encoding/json/v2"
	"io/fs"
	"net/http"
	"os"

	"github.com/kaisawind/skopeoui/web"
	"github.com/sirupsen/logrus"
)

const distPath = "web/dist"

type WebMux struct {
	*http.ServeMux
}

func ServeWebMux(mux *http.ServeMux) *http.ServeMux {
	if _, err := os.Stat(distPath); err != nil {
		logrus.Infoln("local web disabled")
	} else {
		logrus.Infoln("local web enabled")
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(distPath); err != nil {
			fsys, err := fs.Sub(web.WebUI, "dist")
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			fs := http.FS(fsys)
			http.FileServer(fs).ServeHTTP(w, r)
			return
		}
		fs := http.Dir(distPath)
		http.FileServer(fs).ServeHTTP(w, r)
	})
	return mux
}

func Cross(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			rw.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				rw.Header().Add("Access-Control-Allow-Headers", "*")
				rw.Header().Add("Access-Control-Allow-Methods", "*")
				return
			}
		}
		h.ServeHTTP(rw, r)
	})
}

func HttpResponse(rw http.ResponseWriter, code int, data any) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	out := map[string]any{
		"code":    code,
		"success": true,
		"data":    data,
	}
	_ = json.MarshalWrite(rw, out)
}

func HttpError(rw http.ResponseWriter, code int, message string) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	data := map[string]any{
		"success": false,
		"code":    code,
		"error":   message,
	}
	_ = json.MarshalWrite(rw, data)
}
