package server

import (
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type Settings struct {
	Debug   bool
	Port    string
	Timeout int
}

type FuzzyHTTPServer struct {
	Settings Settings
	Metrics  Metrics
}

func MakeFuzzyHTTPServer(s Settings) FuzzyHTTPServer {
	return FuzzyHTTPServer{Settings: s, Metrics: RegisterPrometheusMetrics()}
}

func (f FuzzyHTTPServer) Status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{"status": "ok"}`)
}

func (f FuzzyHTTPServer) Api(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("api call %d", f.Settings.Timeout)
	n := rand.Intn(f.Settings.Timeout * 1000) // n will be between 0 and 10
	time.Sleep(time.Duration(n) * time.Millisecond)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{"status": "ok"}`)
	f.Metrics.ReuestCount.Inc()
	return
}
