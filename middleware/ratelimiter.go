package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

const (
	rateLimit       = 100.0 / 3600.0 // 100 requests for every client every 1 hour
	burstLimit      = 100            // maximum upto 100 requests a client can make at once
	cleanUpInterval = time.Minute    // checks for idle clients every cleanUp Interval
	idleTimeout     = time.Hour      // checks if the client has been inactive fir idle time
)

type Client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

var (
	mu      sync.Mutex
	clients = make(map[string]*Client)
)

func getClientLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	limiter, exist := clients[ip]

	if !exist {
		limiter := rate.NewLimiter(rateLimit, burstLimit)
		clients[ip] = &Client{
			limiter:  limiter,
			lastSeen: time.Now(),
		}
	}
	return limiter.limiter
}

func cleanUpOldclients() {
	for {
		time.Sleep(cleanUpInterval)
		mu.Lock()
		for ip, client := range clients {
			if time.Since(client.lastSeen) > idleTimeout {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}

func PerClientRateLimiter() gin.HandlerFunc {
	go cleanUpOldclients()
	return func(c *gin.Context) {
		ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "unable to parse IP address"})
			return
		}
		limiter := getClientLimiter(ip)
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "Maximum Rate limit exceeded. Please try again after 1 hour"})
			return
		}
		c.Next()
	}
}
