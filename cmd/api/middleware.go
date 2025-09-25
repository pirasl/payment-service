package main

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pascaldekloe/jwt"
	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
)

func (s *service) authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Vary", "Authorization")

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			s.logger.Error("no auth token.")
			s.invalidAuthenticationTokenResponse(c)
			c.Abort()
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {

			s.logger.Error("malformed authorization header", "header", authHeader)
			s.malformedAuthTokenResponse(c)
			c.Abort()
			return
		}

		token := headerParts[1]

		claims, err := jwt.HMACCheck([]byte(token), []byte(s.config.jwtConfig.secret))
		if err != nil {
			s.logger.Error("jwt signature validation failed", "error", err)
			s.invalidAuthenticationTokenResponse(c)
			c.Abort()
			return
		}

		if !claims.Valid(time.Now()) {
			s.logger.Error("jwt token is invalid or expired")
			s.expiredTokenResponse(c)
			c.Abort()
			return
		}

		if claims.Issuer != "api-gateway" {
			s.logger.Error("jwt issuer mismatch", "expected_issuer", "api-gateway", "actual_issuer", claims.Issuer)
			s.invalidAuthenticationTokenResponse(c)
			c.Abort()
			return
		}

		if !claims.AcceptAudience("quizify.leo-piras.com") {

			s.logger.Error("jwt audience mismatch", "expected_audience", "quizify.leo-piras.com", "actual_audiences", claims.Audiences)
			s.invalidAuthenticationTokenResponse(c)
			c.Abort()
			return
		}

		userID, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			s.logger.Error("failed to parse user ID from jwt subject", "subject", claims.Subject, "error", err)
			s.InternalServerErrorResponse(c, err)
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}

func (s *service) rateLimiter() gin.HandlerFunc {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := realip.FromRequest(c.Request)

		mu.Lock()

		if _, found := clients[ip]; !found {
			clients[ip] = &client{
				limiter: rate.NewLimiter(rate.Limit(s.config.rateLimiterConfig.rps), s.config.rateLimiterConfig.burst),
			}
		}

		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			s.rateLimitExceededResponse(c)
			return
		}

		mu.Unlock()

		c.Next()
	}
}

func (s *service) limitBodySize() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, int64(8<<20))
		c.Next()
	}
}
