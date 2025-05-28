package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vfa-khuongdv/golang-cms/internal/utils"
	"github.com/vfa-khuongdv/golang-cms/pkg/logger"
)

type LogResponse struct {
	Method     string `json:"method"`
	URL        string `json:"url"`
	Header     any    `json:"header"`
	Request    any    `json:"request,omitempty"`
	Response   any    `json:"response,omitempty"`
	Latency    string `json:"latency,omitempty"`
	StatusCode string `json:"status_code"`
}

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
			"access_token",
			"refresh_token",
			"ccv",
			"credit_card",
			"debit_card",
			"social_security_number",
			"ssn",
			"bank_account",
			"bank_account_number",
			"email",
			"phone",
			"address",
		}

		timeStart := time.Now()

		var logEntry LogResponse
		logEntry.Method = c.Request.Method
		logEntry.URL = c.Request.URL.String()
		logEntry.Header = c.Request.Header
		logEntry.Request = c.Request.URL.Query()
		logEntry.Response = nil

		// Read and log request body
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(io.LimitReader(c.Request.Body, maxBodySize))
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Restore body for the next handler

			if strings.Contains(c.Request.Header.Get("Content-Type"), "application/json") {
				if err := json.Unmarshal(bodyBytes, &requestBody); err == nil {
					// Mask sensitive data in the request body
					requestBody = utils.CensorSensitiveData(requestBody, sensitiveKeys)
					// Store the masked data directly without converting to string
					logEntry.Request = requestBody
				} else {
					logEntry.Request = string(bodyBytes)
				}
			} else {
				logEntry.Request = string(bodyBytes)
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
		// Calculate latency
		logEntry.Latency = fmt.Sprintf("%d (ms)", timeEnd.Sub(timeStart).Milliseconds())

		// Log response
		var responseBodyData any
		if strings.Contains(c.Writer.Header().Get("Content-Type"), "application/json") {
			if err := json.Unmarshal(responseBody.Bytes(), &responseBodyData); err == nil {
				// Mask sensitive data in the response body
				responseBodyData = utils.CensorSensitiveData(responseBodyData, sensitiveKeys)
				// Store the masked data directly
				logEntry.Response = responseBodyData
			} else {
				logEntry.Response = responseBody.String()
			}
		} else {
			logEntry.Response = responseBody.String()
		}

		logEntry.StatusCode = fmt.Sprintf("%d", c.Writer.Status())

		jsonData, err := json.Marshal(logEntry)
		if err != nil {
			logger.Error("Failed to marshal log entry:", err)
			return
		}
		logger.Info(string(jsonData))
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
