package transport

import (
	"encoding/json"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"storeservice/pkg/common/errors"
	"time"
)

var grpcServeMuxOptions = &runtime.JSONPb{
	EmitDefaults: true,
	OrigName:     true,
}

func NewServeMux() *runtime.ServeMux {
	return runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, grpcServeMuxOptions))
}

func LogMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		h.ServeHTTP(w, r)
		log.WithFields(log.Fields{
			"method":     r.Method,
			"url":        r.URL,
			"remoteAddr": r.RemoteAddr,
			"userAgent":  r.UserAgent(),
			"duration":   time.Since(startTime).String(),
			"at":         startTime,
		}).Info("got request")
	})
}

func RenderJson(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Error(err)
		ProcessError(w, errors.InternalError)
		return
	}
}

func ReadJson(r *http.Request, output interface{}) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer CloseBody(r.Body)

	err = json.Unmarshal(b, &output)
	if err != nil {
		err = fmt.Errorf("can't parse %s to json", b)
	}

	return err
}

func ProcessError(w http.ResponseWriter, e error) {
	if e == errors.InternalError {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	} else {
		http.Error(w, e.Error(), http.StatusBadRequest)
	}
}

func CloseBody(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		log.Error()
	}
}
