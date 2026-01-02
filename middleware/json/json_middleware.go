package json

import (
	"bytes"
	"encoding/json"
	"strconv"
	utilsjson "tiny-admin-api-serve/utils/json"

	"github.com/gin-gonic/gin"
)

// CustomJSON 自定义JSON序列化中间件
func CustomJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建一个响应写入器包装器
		writer := &responseWriter{Body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = writer

		// 处理请求
		c.Next()

		// 检查响应状态码和内容类型
		contentType := c.Writer.Header().Get("Content-Type")
		// 调试：打印Content-Type和响应头
		// log.Printf("Content-Type: %s", contentType)
		// log.Printf("Response Headers: %v", c.Writer.Header())
		// log.Printf("Response Status: %d", c.Writer.Status())

		// 更宽松地检查JSON响应类型
		if contentType == "application/json; charset=utf-8" ||
			contentType == "application/json" ||
			contentType == "application/json; charset=UTF-8" {
			// 如果是JSON响应，使用自定义的JSON序列化重新处理
			var data interface{}
			// 调试：打印原始响应内容
			// log.Printf("Original Response: %s", writer.Body.Bytes())

			if err := json.Unmarshal(writer.Body.Bytes(), &data); err == nil {
				// 使用自定义JSON序列化
				customData, err := utilsjson.Marshal(data)
				if err == nil {
					// 更新Content-Length
					writer.ResponseWriter.Header().Set("Content-Length", strconv.Itoa(len(customData)))
					// 调试：打印自定义序列化后的内容
					// log.Printf("Custom Response: %s", customData)
					// 写入自定义序列化后的内容
					writer.ResponseWriter.Write(customData)
					return
				}
			} else {
				// 调试：打印反序列化错误
				// log.Printf("Unmarshal error: %v", err)
			}
		}
		// 如果不是JSON响应或者处理失败，直接写入原始内容
		writer.ResponseWriter.Write(writer.Body.Bytes())
	}
}

// responseWriter 响应写入器包装器
type responseWriter struct {
	gin.ResponseWriter
	Body      *bytes.Buffer
	isWritten bool
}

// Write 实现http.ResponseWriter接口
func (w *responseWriter) Write(b []byte) (int, error) {
	// 只将内容写入缓冲区，不立即写入响应
	w.isWritten = true
	return w.Body.Write(b)
}

// WriteHeader 实现http.ResponseWriter接口
func (w *responseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

// Written 实现gin.ResponseWriter接口，返回是否已经写入过响应
func (w *responseWriter) Written() bool {
	return w.isWritten
}
