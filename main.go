package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/tnosaj/fuzzyhttpserver/server"
)

func main() {
	var wait time.Duration
	s, err := evaluateInputs()
	if err != nil {
		log.Fatalf("could not evaluate inputs: %q", err)
	}

	setupLogger(s.Debug)

	fuzzyHttpServer := server.MakeFuzzyHTTPServer(s)

	// Generate seedyness
	rand.New(rand.NewSource(time.Now().UnixNano()))

	http.HandleFunc("/status", fuzzyHttpServer.Status)
	http.HandleFunc("/api", fuzzyHttpServer.Api)
	http.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%s", s.Port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * time.Duration(s.Timeout),
		ReadTimeout:  time.Second * time.Duration(s.Timeout),
		IdleTimeout:  time.Second * time.Duration(s.Timeout),
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logrus.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

func setupLogger(debug bool) {
	//logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.Debug("Configured logger")
}
