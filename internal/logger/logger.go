package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type responseData struct {
	status int
	size   int
}

type logResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

var Log *zap.Logger = zap.NewNop()

func Init(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapLogger, err := cfg.Build()
	if err != nil {
		return err
	}
	Log = zapLogger
	return nil
}

func RequestLogger(h http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method

		h.ServeHTTP(w, r)
		duration := time.Since(start)

		Log.Sugar().Debugln(
			"incoming HTTP request",
			"uri", uri,
			"method", method,
			"duration", duration,
		)
	})
}

func ResponseLogger(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respLogger := logResponseWriter{
			ResponseWriter: w,
			responseData:   &responseData{},
		}
		h.ServeHTTP(&respLogger, r)
		Log.Sugar().Debugln(
			"sending response",
			"status", respLogger.responseData.status,
			"size", respLogger.responseData.size,
		)
	})
}

func (r *logResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *logResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.Header().Set("Content-Type", "application/json")
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}
