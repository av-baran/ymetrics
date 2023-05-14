package logger

import (
	"fmt"
	"net/http"
	"time"

	"github.com/av-baran/ymetrics/internal/config"
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

var logger *zap.SugaredLogger = zap.NewNop().Sugar()

func Init(cfg config.LoggerConfig) error {
	lvl, err := zap.ParseAtomicLevel(cfg.Level)
	if err != nil {
		return fmt.Errorf("cannot parse log level: %w", err)
	}

	zapCfg := zap.NewProductionConfig()
	zapCfg.Level = lvl
	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapLogger, err := zapCfg.Build()
	if err != nil {
		return fmt.Errorf("cannot build logger config: %w", err)
	}
	logger = zapLogger.Sugar()
	return nil
}

func RequestLogMiddlware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method

		h.ServeHTTP(w, r)
		duration := time.Since(start)

		logger.Infow(
			"incoming HTTP request",
			"uri", uri,
			"method", method,
			"duration", duration,
		)
	})
}

func ResponseLogMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respLogger := logResponseWriter{
			ResponseWriter: w,
			responseData:   &responseData{},
		}
		h.ServeHTTP(&respLogger, r)
		logger.Infow(
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

func Fatal(args ...any) {
	logger.Fatal(args...)
}
func Error(args ...any) {
	logger.Error(args...)
}
func Warn(args ...any) {
	logger.Warn(args...)
}
func Info(args ...any) {
	logger.Info(args...)
}
func Debug(args ...any) {
	logger.Debug(args...)
}

func Fatalf(format string, opts ...any) {
	logger.Fatalf(format, opts...)
}
func Errorf(format string, opts ...any) {
	logger.Errorf(format, opts...)
}
func Warnf(format string, opts ...any) {
	logger.Warnf(format, opts...)
}
func Infof(format string, opts ...any) {
	logger.Infof(format, opts...)
}
func Debugf(format string, opts ...any) {
	logger.Debugf(format, opts...)
}

func Sync() {
	logger.Sync()
}
