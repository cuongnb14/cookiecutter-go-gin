package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"{{ cookiecutter.project_slug }}/configs"
)

func RequestLogger() gin.HandlerFunc {
	//logger := configs.GetLogger()
	return func(c *gin.Context) {
		startTime := time.Now()
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		headerMap := map[string][]string{}
		_ = copier.Copy(&headerMap, c.Request.Header)
		headerMap["Authorization"] = []string{"filtered"}

		header, _ := json.Marshal(headerMap)
		query := c.Request.URL.RawQuery
		body := strings.ReplaceAll(string(bodyBytes), "\n", "")

		c.Next()

		if c.Request.URL.Path == "/healthz" {
			return
		}

		latency := time.Since(startTime)
		timestamp := startTime.Format("2006-01-02 15:04:05")

		if configs.Env.EnableLogRequestDetail {
			fmt.Printf(
				"%s %s %s %s %d %d headers=%s query=%s, body=%s\n",
				timestamp,
				c.ClientIP(),
				c.Request.Method,
				c.Request.URL.Path,
				c.Writer.Status(),
				latency.Milliseconds(),
				header, query, body,
			)
		} else {
			fmt.Printf(
				"%s %s %s %s %d %d\n",
				timestamp,
				c.ClientIP(),
				c.Request.Method,
				c.Request.URL.Path,
				c.Writer.Status(),
				latency.Milliseconds(),
			)
		}
	}
}
