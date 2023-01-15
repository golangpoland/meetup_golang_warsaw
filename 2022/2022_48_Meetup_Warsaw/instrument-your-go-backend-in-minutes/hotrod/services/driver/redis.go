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

package driver

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"example.com/pkg/delay"
	"example.com/pkg/log"
	"example.com/pkg/tracing"
	"example.com/services/config"
)

// Redis is a simulator of remote Redis cache
type Redis struct {
	tracerProvider trace.TracerProvider
	logger         log.Factory
	errorSimulator
}

func newRedis(logger log.Factory) *Redis {
	return &Redis{
		tracerProvider: tracing.Init("redis"),
		logger:         logger,
	}
}

// FindDriverIDs finds IDs of drivers who are near the location.
func (r *Redis) FindDriverIDs(ctx context.Context, _ string) []string {
	ctx, span := r.tracerProvider.Tracer("").Start(ctx, "FindDriverIDs", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	// simulate RPC delay
	delay.Sleep(config.RedisFindDelay, config.RedisFindDelayStdDev)

	drivers := make([]string, 10)
	for i := range drivers {
		// #nosec
		drivers[i] = fmt.Sprintf("T7%05dC", rand.Int()%100000)
	}
	r.logger.For(ctx).Info("Found drivers", zap.Strings("drivers", drivers))
	return drivers
}

// GetDriver returns driver and the current car location
func (r *Redis) GetDriver(ctx context.Context, driverID string) (Driver, error) {
	ctx, span := r.tracerProvider.Tracer("").Start(ctx, "GetDriver")
	span.SetAttributes(attribute.String("param.driverID", driverID))
	defer span.End()

	// simulate RPC delay
	delay.Sleep(config.RedisGetDelay, config.RedisGetDelayStdDev)
	if err := r.checkError(); err != nil {
		if span := trace.SpanFromContext(ctx); span != nil {
			span.SetStatus(codes.Error, "err")
		}
		r.logger.For(ctx).Error("redis timeout", zap.String("driver_id", driverID), zap.Error(err))
		return Driver{}, err
	}

	// #nosec
	return Driver{
		DriverID: driverID,
		Location: fmt.Sprintf("%d,%d", rand.Int()%1000, rand.Int()%1000),
	}, nil
}

var errTimeout = errors.New("redis timeout")

type errorSimulator struct {
	sync.Mutex
	countTillError int
}

func (es *errorSimulator) checkError() error {
	es.Lock()
	es.countTillError--
	if es.countTillError > 0 {
		es.Unlock()
		return nil
	}
	es.countTillError = 5
	es.Unlock()
	delay.Sleep(2*config.RedisGetDelay, 0) // add more delay for "timeout"
	return errTimeout
}