package http

import (
	"bytes"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockLogger struct {
	buffer bytes.Buffer
}

func (m *mockLogger) Write(p []byte) (n int, err error) {
	return m.buffer.Write(p)
}

func TestNewRouter(t *testing.T) {
	t.Run("creates a new router with custom logger", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		mockLog := &mockLogger{}
		router := NewRouter(mockLog)

		assert.NotNil(t, router)
		assert.IsType(t, &Router{}, router)

		castRouter, ok := router.(*Router)
		assert.True(t, ok)
		assert.NotNil(t, castRouter.Engine)

		testMessage := "test message"
		_, _ = mockLog.Write([]byte(testMessage))
		assert.Contains(t, mockLog.buffer.String(), testMessage)
	})
}
