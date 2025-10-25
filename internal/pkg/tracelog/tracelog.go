package tracelog

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var globalLogger = log.Logger

func SetLogger(logger zerolog.Logger) {
	globalLogger = logger
}

type TraceLogger struct {
	name      string
	span      trace.Span
	logger    zerolog.Logger
	err       error
	startedAt time.Time
}

func (t *TraceLogger) StartedAt() time.Time {
	return t.startedAt
}

func (t *TraceLogger) Duration() time.Duration {
	return time.Since(t.startedAt)
}

func Begin(ctx context.Context, name string) (tl *TraceLogger, ctxWithTraceId context.Context) {
	tl = &TraceLogger{name: name, logger: globalLogger}
	ctx, tl.span = otel.Tracer("").Start(ctx, name)

	tl.AddAttributes(String(traceIDLogKey, tl.GetTraceID()), String(spanNameLogKey, tl.name))
	tl.AddEvent("Begin span", zerolog.TraceLevel)
	tl.startedAt = time.Now().UTC()

	return tl, ctx
}

func (t *TraceLogger) End() {
	timeElapsed := time.Since(t.startedAt).String()

	if t.span == nil || !t.span.IsRecording() {
		return
	}

	attrs := []KeyValue{String(spanTimeElapsedLogKey, timeElapsed)}

	if t.err != nil {
		t.span.RecordError(t.err)
		t.span.SetStatus(codes.Error, t.err.Error())
		attrs = append(attrs, String(spanErrorLogKey, t.err.Error()))
	} else {
		t.span.SetStatus(codes.Ok, "")
	}

	t.AddEvent(fmt.Sprintf("End span '%s', Time: %s", t.name, timeElapsed), zerolog.TraceLevel, attrs...)
	t.span.End()
}

func (t *TraceLogger) AddAttributes(attributes ...KeyValue) {
	for _, attr := range attributes {
		switch attr.vT {
		case vTString:
			t.span.SetAttributes(attribute.String(attr.k, attr.v.(string)))
			t.logger = t.logger.With().Str(attr.k, attr.v.(string)).Logger()
		case vTInt:
			t.span.SetAttributes(attribute.Int(attr.k, attr.v.(int)))
			t.logger = t.logger.With().Int(attr.k, attr.v.(int)).Logger()
		case vTBool:
			t.span.SetAttributes(attribute.Bool(attr.k, attr.v.(bool)))
			t.logger = t.logger.With().Bool(attr.k, attr.v.(bool)).Logger()
		default:
			continue
		}
	}
}

func (t *TraceLogger) AddEvent(eventName string, logLevel zerolog.Level, attributes ...KeyValue) {
	var eventAttrs []trace.EventOption
	var tmpLogger = t.logger

	for i := range attributes {
		switch attributes[i].vT {
		case vTString:
			eventAttrs = append(eventAttrs, trace.WithAttributes(attribute.String(attributes[i].k, attributes[i].v.(string))))
			tmpLogger = tmpLogger.With().Str(attributes[i].k, attributes[i].v.(string)).Logger()
		case vTInt:
			eventAttrs = append(eventAttrs, trace.WithAttributes(attribute.Int(attributes[i].k, attributes[i].v.(int))))
			tmpLogger = tmpLogger.With().Int(attributes[i].k, attributes[i].v.(int)).Logger()
		case vTBool:
			eventAttrs = append(eventAttrs, trace.WithAttributes(attribute.Bool(attributes[i].k, attributes[i].v.(bool))))
			tmpLogger = tmpLogger.With().Bool(attributes[i].k, attributes[i].v.(bool)).Logger()
		default:
			continue
		}
	}

	t.span.AddEvent(eventName, eventAttrs...)
	tmpLogger.WithLevel(logLevel).Msg(eventName)
}

func (t *TraceLogger) AddError(err error, logLevel ...zerolog.Level) {
	if err == nil {
		return
	}

	t.err = err

	if len(logLevel) > 0 {
		t.logger.WithLevel(logLevel[0]).Msg(err.Error())
		return
	}

	t.logger.Error().Msg(err.Error())
}

func (t *TraceLogger) GetTraceID() (id string) {
	if t.span != nil {
		return t.span.SpanContext().TraceID().String()
	}
	return
}
