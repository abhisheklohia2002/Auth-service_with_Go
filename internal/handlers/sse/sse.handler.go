package handlerSSE

import (
	"fmt"
	"io"

	"example.com/m/internal/common/sse"
	"github.com/gin-gonic/gin"
)

type SSEHandler struct {
	hub *sse.Hub
}

func NewSSEHandler(hub *sse.Hub) *SSEHandler {
	return &SSEHandler{
		hub: hub,
	}
}

func (h *SSEHandler) StreamNotifications(c *gin.Context) {
	client := make(sse.ClientChan)

	h.hub.AddClient(client)
	defer h.hub.RemoveClient(client)

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-client; ok {
			fmt.Fprintf(w, "event: notification\n")
			fmt.Fprintf(w, "data: %s\n\n", msg)
			return true
		}

		return false
	})
}
