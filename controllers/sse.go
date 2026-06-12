package controllers

import (
	"fmt"
	"io"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SSEHub struct {
	mu      sync.RWMutex
	clients map[chan string]struct{}
	log     *zap.Logger
}

func NewSSEHub(log *zap.Logger) *SSEHub {
	return &SSEHub{
		clients: make(map[chan string]struct{}),
		log:     log,
	}
}

func (h *SSEHub) Subscribe() (chan string, func()) {
	h.mu.Lock()
	defer h.mu.Unlock()

	ch := make(chan string, 64)
	h.clients[ch] = struct{}{}

	unsubscribe := func() {
		h.mu.Lock()
		delete(h.clients, ch)
		h.mu.Unlock()
		close(ch)
	}

	return ch, unsubscribe
}

func (h *SSEHub) Broadcast(data string) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for ch := range h.clients {
		select {
		case ch <- data:
		default:
		}
	}
}

type SSEController struct {
	hub *SSEHub
}

func NewSSEController(hub *SSEHub) *SSEController {
	return &SSEController{hub: hub}
}

func (ctrl *SSEController) Stream(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	ch, unsubscribe := ctrl.hub.Subscribe()
	defer unsubscribe()

	c.Stream(func(w io.Writer) bool {
		msg, ok := <-ch
		if !ok {
			return false
		}
		fmt.Fprintf(w, "data: %s\n\n", msg)
		return true
	})
}
