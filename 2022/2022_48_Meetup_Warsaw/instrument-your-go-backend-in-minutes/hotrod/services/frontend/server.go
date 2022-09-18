// Copyright (c) 2019 The Jaeger Authors.
// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package frontend

import (
	"embed"
	"encoding/json"
	"net/http"
	"path"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/jaegertracing/jaeger/pkg/httpfs"

	"example.com/pkg/httperr"
	"example.com/pkg/log"
	"example.com/pkg/tracing"
)

//go:embed web_assets/*
var assetFS embed.FS

// Server implements Frontend service
type Server struct {
	hostPort  string
	tp        trace.TracerProvider
	logger    log.Factory
	bestETA   *bestETA
	assetFS   http.FileSystem
	basepath  string
	grafanaUI string
}

// ConfigOptions used to make sure service clients
// can find correct server ports
type ConfigOptions struct {
	FrontendHostPort string
	DriverHostPort   string
	CustomerHostPort string
	RouteHostPort    string
	Basepath         string
	GrafanaUI        string
}

// NewServer creates a new frontend.Server
func NewServer(options ConfigOptions, tp trace.TracerProvider, logger log.Factory) *Server {
	return &Server{
		hostPort:  options.FrontendHostPort,
		tp:        tp,
		logger:    logger,
		bestETA:   newBestETA(tp, logger, options),
		assetFS:   httpfs.PrefixedFS("web_assets", http.FS(assetFS)),
		basepath:  options.Basepath,
		grafanaUI: options.GrafanaUI,
	}
}

// Run starts the frontend server
func (s *Server) Run() error {
	mux := s.createServeMux()
	s.logger.Bg().Info("Starting", zap.String("address", "http://"+path.Join(s.hostPort, s.basepath)))
	return http.ListenAndServe(s.hostPort, mux)
}

func (s *Server) createServeMux() http.Handler {
	mux := tracing.NewServeMux(s.tp)
	p := path.Join("/", s.basepath)
	mux.Handle(p, http.StripPrefix(p, http.FileServer(s.assetFS)))
	mux.Handle(path.Join(p, "/dispatch"), http.HandlerFunc(s.dispatch))
	mux.Handle(path.Join(p, "/config"), http.HandlerFunc(s.config))
	mux.Handle("/metrics", promhttp.Handler())
	return mux
}

func (s *Server) config(w http.ResponseWriter, r *http.Request) {
	config := map[string]string{
		"grafana": s.grafanaUI,
	}
	s.writeResponse(config, "", w, r)
}

func (s *Server) dispatch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s.logger.For(ctx).Info("HTTP request received", zap.String("method", r.Method), zap.Stringer("url", r.URL))
	if err := r.ParseForm(); httperr.HandleError(w, err, http.StatusBadRequest) {
		s.logger.For(ctx).Error("bad request", zap.Error(err))
		return
	}

	customerID := r.Form.Get("customer")
	if customerID == "" {
		http.Error(w, "Missing required 'customer' parameter", http.StatusBadRequest)
		return
	}

	// TODO distinguish between user errors (such as invalid customer ID) and server failures
	response, err := s.bestETA.Get(ctx, customerID)
	if httperr.HandleError(w, err, http.StatusInternalServerError) {
		s.logger.For(ctx).Error("request failed", zap.Error(err))
		return
	}

	traceID := trace.SpanFromContext(ctx).SpanContext().TraceID().String()
	s.writeResponse(response, traceID, w, r)
}

func (s *Server) writeResponse(response interface{}, traceID string, w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(response)
	if httperr.HandleError(w, err, http.StatusInternalServerError) {
		s.logger.For(r.Context()).Error("cannot marshal response", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("trace-id", traceID)
	w.Write(data)
}
