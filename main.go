package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

// Credit scoring logic
type CreditRequest struct {
	UserID string  `json:"user_id"`
	Amount float64 `json:"amount"`
	Income float64 `json:"income"`
}

type CreditResponse struct {
	UserID   string `json:"user_id"`
	Score    int    `json:"score"`
	Approved bool   `json:"approved"`
	Reason   string `json:"reason"`
}

// Simple ML-based scoring (simulated)
func scoreCreditApplication(req CreditRequest) CreditResponse {
	rand.Seed(time.Now().UnixNano())

	// Basic risk model: income/amount ratio + random noise
	baseScore := 600
	if req.Income > 0 {
		ratio := req.Income / req.Amount
		baseScore += int(ratio * 100)
	}

	// Add random variation (simulate ML model)
	score := baseScore + rand.Intn(100) - 50
	if score > 850 {
		score = 850
	}
	if score < 300 {
		score = 300
	}

	approved := score >= 650
	reason := "Approved based on income-to-loan ratio"
	if !approved {
		reason = "Insufficient creditworthiness"
	}

	return CreditResponse{
		UserID:   req.UserID,
		Score:    score,
		Approved: approved,
		Reason:   reason,
	}
}

// Kafka event publisher (audit log)
func publishAuditEvent(ctx context.Context, event CreditResponse) error {
	writer := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "credit-applications",
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	eventBytes, _ := json.Marshal(event)
	return writer.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(event.UserID),
			Value: eventBytes,
		},
	)
}

func main() {
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "ArdaCredit Backend Prototype"})
	})

	// Credit application endpoint
	r.POST("/api/v1/apply", func(c *gin.Context) {
		var req CreditRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		// Score the application
		response := scoreCreditApplication(req)

		// Async: Publish to Kafka (audit trail) - fire and forget
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := publishAuditEvent(ctx, response); err != nil {
				log.Printf("Failed to publish audit event: %v", err)
			}
		}()

		c.JSON(200, response)
	})

	fmt.Println(" ArdaCredit Backend Prototype running on :8080")
	r.Run(":8080")
}
