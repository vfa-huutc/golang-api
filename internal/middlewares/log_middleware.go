package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/vfa-khuongdv/golang-cms/internal/utils"
)

// Middleware for logging requests and responses in Gin
func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody any
		const maxBodySize = 1 << 20 // Limit body size to 1 MB
		// Define sensitive keys to be masked in logs
		// These keys will be masked in both request and response logs
		sensitiveKeys := []string{
			"password",
			"api-key",
			"token",
			"email",
			"phone",
			"address",
		}

		timeStart := time.Now()

		logEntry := logrus.WithFields(logrus.Fields{
			"method":  c.Request.Method,
			"url":     c.Request.URL.String(),
			"header":  c.Request.Header,
			"request": c.Request.URL.Query(),
		})

		// Read and log request body
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(io.LimitReader(c.Request.Body, maxBodySize))
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Restore body for the next handler

			if strings.Contains(c.Request.Header.Get("Content-Type"), "application/json") {
				if err := json.Unmarshal(bodyBytes, &requestBody); err == nil {
					// Mask sensitive data in the request body
					requestBody = utils.CensorSensitiveData(requestBody, sensitiveKeys)
					logEntry = logEntry.WithField("request", requestBody)
				} else {
					logEntry = logEntry.WithField("request_raw", string(bodyBytes))
				}
			} else {
				logEntry = logEntry.WithField("request_raw", string(bodyBytes))
			}
		}

		// Create a buffer to capture the response body
		responseBody := &bytes.Buffer{}
		c.Writer = &bodyWriter{
			ResponseWriter: c.Writer,
			body:           responseBody,
		}

		// Process the request
		c.Next()

		timeEnd := time.Now()

		logEntry = logEntry.WithFields(logrus.Fields{
			"latency": fmt.Sprintf("%d (ms)", timeEnd.Sub(timeStart).Milliseconds()),
		})

		// Log response
		var responseBodyData any
		if strings.Contains(c.Writer.Header().Get("Content-Type"), "application/json") {
			if err := json.Unmarshal(responseBody.Bytes(), &responseBodyData); err == nil {
				// Mask sensitive data in the response body
				responseBodyData = utils.CensorSensitiveData(responseBodyData, sensitiveKeys)
				logEntry = logEntry.WithField("response", responseBodyData)
			} else {
				logEntry = logEntry.WithField("response_raw", responseBody.String())
			}
		} else {
			logEntry = logEntry.WithField("response_raw", responseBody.String())
		}

		logEntry = logEntry.WithField("status_code", c.Writer.Status())
		logEntry.Info("Request handled")

	}
}

// bodyWriter is a custom ResponseWriter to capture the response body
type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
