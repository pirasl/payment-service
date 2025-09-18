package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *service) routes() http.Handler {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, int64(8<<20))
		c.Next()
	})
	r.GET("/healthcheck", s.healthCheckHandler)

	rv1 := r.Group("/stripe/v1")
	rv1.POST("/create-payment-intent", s.createPaymentIntentHandler)

	rv1.POST("/webhook", s.webhookHandler)

	return r
}
