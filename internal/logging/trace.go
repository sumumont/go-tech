package logging

import (
	"context"
	"github.com/rs/xid"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const traceKey = "TraceId"
const spanId = "Spanid"

var traceIdKeys = []string{"x-request-id", "traceid", "x-b3-traceid"}

type traceLogger struct {
	logger zerolog.Logger
	ctx    context.Context
}

func Trace() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从header里面拿 trace-id
		var traceId string
		for _, headKey := range traceIdKeys {
			traceId = ctx.GetHeader(headKey)
			if traceId != "" {
				break
			}
		}

		if traceId == "" {
			// 没有则创建
			traceId = xid.New().String()
		}
		ctx.Set(traceKey, traceId)
		ctx.Request.Header.Set(traceKey, traceId)
		ctx.Next()
	}
}

func (l *traceLogger) Debug() *zerolog.Event {
	traceID := traceIdFromCtx(l.ctx)
	spanID := spanIdFromCtx(l.ctx)
	return logger.Debug().Str("time", Timer()).Caller().
		Str(traceKey, traceID).Str(spanId, spanID)
}

func (l *traceLogger) Info() *zerolog.Event {
	traceID := traceIdFromCtx(l.ctx)
	spanID := spanIdFromCtx(l.ctx)
	return logger.Info().Str("time", Timer()).Caller().
		Str(traceKey, traceID).Str(spanId, spanID)
}

func (l *traceLogger) Warn() *zerolog.Event {
	traceID := traceIdFromCtx(l.ctx)
	spanID := spanIdFromCtx(l.ctx)
	return logger.Warn().Str("time", Timer()).Caller().
		Str(traceKey, traceID).Str(spanId, spanID)
}

func (l *traceLogger) Error(err error) *zerolog.Event {
	traceID := traceIdFromCtx(l.ctx)
	spanID := spanIdFromCtx(l.ctx)

	return logger.Err(err).Str("time", Timer()).Caller().
		Str(traceKey, traceID).Str(spanId, spanID)
}

func (l *traceLogger) Fatal() *zerolog.Event {
	traceID := traceIdFromCtx(l.ctx)
	spanID := spanIdFromCtx(l.ctx)

	return logger.Fatal().Str("time", Timer()).Caller().
		Str(traceKey, traceID).Str(spanId, spanID)
}

func WithContext(ctx context.Context) *traceLogger {
	return &traceLogger{
		logger: logger,
		ctx:    ctx,
	}
}

func spanIdFromCtx(ctx context.Context) string {
	xTraceId, ok := ctx.Value(spanId).(string)
	if ok {
		return xTraceId
	}

	return ""
}

func traceIdFromCtx(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	xSpanId, ok := ctx.Value(traceKey).(string)
	if ok {
		return xSpanId
	}

	return ""
}
