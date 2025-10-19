package handler

import (
	"catcher/pkg/logging"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
)

func ResponseWriterMW(logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer = &ResponseWriter{
			ResponseWriter: c.Writer,
			logger:         logger,
		}
		c.Next()
	}
}

func SafeBodyReaderMW(logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body != nil {
			c.Request.Body = &ReadCloser{
				ReadCloser: c.Request.Body,
				logger:     logger}
		}
		c.Next()
	}
}

type ResponseWriter struct {
	gin.ResponseWriter
	logger logging.Logger
}

type ReadCloser struct {
	io.ReadCloser
	logger logging.Logger
}

func (r *ResponseWriter) Write(data []byte) (int, error) {
	n, err := r.ResponseWriter.Write(data)
	if err != nil {
		if isConnectionError(err) {
			r.logger.Info("Client disconnected",
				r.logger.Err(err))
			return n, nil // возвращаем nil вместо ошибки
		}
	}
	return n, err
}

func (r *ResponseWriter) WriteString(s string) (int, error) {
	n, err := r.ResponseWriter.WriteString(s)
	if err != nil {
		if isConnectionError(err) {
			r.logger.Info("Client disconnected",
				r.logger.Err(err))
			return n, nil
		}
	}
	return n, err
}

func (r *ReadCloser) Read(p []byte) (n int, err error) {
	n, err = r.ReadCloser.Read(p)
	if err != nil && err != io.EOF {
		if isConnectionError(err) {
			r.logger.Info("Client disconnected",
				r.logger.Err(err))
			return n, io.EOF // преобразуем в EOF
		}
	}
	return n, err
}

func isConnectionError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := strings.ToLower(err.Error())

	// Явные индикаторы разрыва соединения клиентом
	return strings.Contains(errMsg, "wsarecv") ||
		strings.Contains(errMsg, "wsasend") ||
		strings.Contains(errMsg, "connection reset by peer") ||
		strings.Contains(errMsg, "broken pipe") ||
		strings.Contains(errMsg, "forcibly closed") ||
		strings.Contains(errMsg, "use of closed network connection")

}
