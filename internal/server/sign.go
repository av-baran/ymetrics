package server

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/av-baran/ymetrics/internal/logger"
	"github.com/av-baran/ymetrics/pkg/interrors"
)

type SignResponseWriter struct {
	http.ResponseWriter
	encodedSign   string
	signSecretKey string
}

func (s *Server) checkSignMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := s.cfg.SignSecretKey
		if key == "" || !isSignHeaderPresent(r) {
			logger.Infof("key is empty keep running without sign check")
			h.ServeHTTP(w, r)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			sendError(w, "cannot read request body", err)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		gotHash := r.Header.Get("HashSHA256")
		decodedHash, err := hex.DecodeString(gotHash)
		if err != nil {
			sendError(w, "cannot decode sign", err)
			return
		}

		sign := signBody(s.cfg.SignSecretKey, body)
		if hmac.Equal(decodedHash, sign) {
			h.ServeHTTP(w, r)
		} else {
			sendError(w, "cannot handle request", interrors.ErrInvalidSign)
			return
		}
	})
}

func (s *Server) addSignMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := s.cfg.SignSecretKey
		if key == "" || !isSignHeaderPresent(r) {
			logger.Infof("key is empty keep running without adding sign")
			h.ServeHTTP(w, r)
			return
		}

		signWriter := &SignResponseWriter{
			ResponseWriter: w,
		}

		h.ServeHTTP(signWriter, r)
	})
}

func (r *SignResponseWriter) Write(b []byte) (int, error) {
	sign := signBody(r.signSecretKey, b)
	r.encodedSign = hex.EncodeToString(sign)
	size, err := r.ResponseWriter.Write(b)
	return size, err
}

func (r *SignResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.Header().Set("HashSHA256", r.encodedSign)
	r.ResponseWriter.WriteHeader(statusCode)
}

func signBody(key string, body []byte) []byte {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(body)
	result := h.Sum(nil)

	return result
}

func isSignHeaderPresent(r *http.Request) bool {
	signHeaderName := http.CanonicalHeaderKey("HashSHA256")
	_, signOk := r.Header[signHeaderName]
	return signOk
}
