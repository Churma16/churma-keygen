package middleware

import (
	"net/http"
	"sync"
	"time"

	"churma-keygen/backend/dtos"

	"github.com/gin-gonic/gin"
)

type clientLimit struct {
	lastSeen time.Time
	count    int
}

var (
	clientsMu sync.Mutex
	clients   = make(map[string]*clientLimit)

	// Configurable parameters
	maxRequests = 10
	timeWindow  = time.Minute
)

func init() {
	// Cleanup old clients periodically in the background
	go func() {
		for {
			time.Sleep(10 * time.Minute)
			clientsMu.Lock()
			for ip, cl := range clients {
				if time.Since(cl.lastSeen) > timeWindow {
					delete(clients, ip)
				}
			}
			clientsMu.Unlock()
		}
	}()
}

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		clientsMu.Lock()
		cl, exists := clients[ip]
		now := time.Now()

		if !exists {
			clients[ip] = &clientLimit{
				lastSeen: now,
				count:    1,
			}
			clientsMu.Unlock()
			c.Next()
			return
		}

		// Reset count if time window has passed
		if now.Sub(cl.lastSeen) > timeWindow {
			cl.count = 1
			cl.lastSeen = now
			clientsMu.Unlock()
			c.Next()
			return
		}

		cl.count++
		cl.lastSeen = now

		if cl.count > maxRequests {
			clientsMu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, dtos.NewErrorResponse(http.StatusTooManyRequests, "Too many requests. Please try again later."))
			return
		}

		clientsMu.Unlock()
		c.Next()
	}
}
