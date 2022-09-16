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

package log

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type spanLogger struct {
	logger     *zap.Logger
	span       trace.Span
	spanFields []zapcore.Field
}

func (sl spanLogger) Info(msg string, fields ...zapcore.Field) {
	sl.logToSpan("info", msg)
	sl.logger.Info(msg, append(sl.spanFields, fields...)...)
}

func (sl spanLogger) Error(msg string, fields ...zapcore.Field) {
	sl.logToSpan("error", msg)
	sl.logger.Error(msg, append(sl.spanFields, fields...)...)
}

func (sl spanLogger) Fatal(msg string, fields ...zapcore.Field) {
	sl.logToSpan("fatal", msg)
	sl.logger.Fatal(msg, append(sl.spanFields, fields...)...)
}

// With creates a child logger, and optionally adds some context fields to that logger.
func (sl spanLogger) With(fields ...zapcore.Field) Logger {
	return spanLogger{logger: sl.logger.With(fields...), span: sl.span, spanFields: sl.spanFields}
}

func (sl spanLogger) logToSpan(level string, msg string) {
	sl.span.SetAttributes(
		attribute.String("event", msg),
		attribute.String("level", level),
	)
}
