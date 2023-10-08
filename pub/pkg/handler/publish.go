package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) publish(c *gin.Context) {
	// Read the request body into a byte slice
	channel := "my-channel"
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		http.Error(c.Writer, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Publish the message to the NATS Streaming server
	if err := h.sc.Publish(channel, body); err != nil {
		http.Error(c.Writer, fmt.Sprintf("Error publishing message: %v", err), http.StatusInternalServerError)
		return
	}

	// Send a response to the client
	c.JSON(http.StatusOK, gin.H{
		"message": "Message published successfully",
	})
}

// func (h *Handler) kek(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Kek",
// 	})
// }
