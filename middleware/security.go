package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// SecurityHeaders sets essential HTTP security headers on every response.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// XSS protection
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")

		// Enforce HTTPS (1 year)
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Content Security Policy
		c.Header("Content-Security-Policy", "default-src 'self'")

		// Prevent referrer information leakage
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Restrict browser feature access
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")

		// Hide server information
		c.Header("Server", "")

		c.Next()
	}
}

// CORSConfig holds settings for the CORS middleware.
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	MaxAge         int // preflight cache duration in seconds
}

// DefaultCORSConfig returns a sensible default CORS configuration.
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Origin", "Content-Type", "Authorization"},
		MaxAge:         86400,
	}
}

// CORS middleware handles Cross-Origin Resource Sharing.
func CORS(cfg CORSConfig) gin.HandlerFunc {
	origins := strings.Join(cfg.AllowedOrigins, ", ")
	methods := strings.Join(cfg.AllowedMethods, ", ")
	headers := strings.Join(cfg.AllowedHeaders, ", ")

	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", origins)
		c.Header("Access-Control-Allow-Methods", methods)
		c.Header("Access-Control-Allow-Headers", headers)
		c.Header("Access-Control-Max-Age", formatInt(cfg.MaxAge))

		// Respond immediately to preflight requests
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// visitor keeps track of requests from a single IP.
type visitor struct {
	tokens    int
	lastSeen  time.Time
}

// RateLimiter middleware limits requests per IP.
// maxRequests: maximum number of requests allowed within the time window.
// window: duration of the time window.
func RateLimiter(maxRequests int, window time.Duration) gin.HandlerFunc {
	var mu sync.Mutex
	visitors := make(map[string]*visitor)

	// Periodically clean up stale entries
	go func() {
		for {
			time.Sleep(window)
			mu.Lock()
			for ip, v := range visitors {
				if time.Since(v.lastSeen) > window {
					delete(visitors, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		v, exists := visitors[ip]
		if !exists || time.Since(v.lastSeen) > window {
			visitors[ip] = &visitor{tokens: maxRequests - 1, lastSeen: time.Now()}
			mu.Unlock()
			c.Next()
			return
		}

		if v.tokens <= 0 {
			mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests, please try again later",
			})
			return
		}

		v.tokens--
		v.lastSeen = time.Now()
		mu.Unlock()

		c.Next()
	}
}

// MaxBodySize middleware limits the request body size.
// maxBytes: maximum allowed body size in bytes.
func MaxBodySize(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxBytes {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "Request body too large",
			})
			return
		}

		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}

// ContentTypeCheck middleware ensures that POST, PUT, and PATCH requests
// have a valid Content-Type header (application/json).
func ContentTypeCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
			ct := c.GetHeader("Content-Type")
			if !strings.Contains(ct, "application/json") {
				c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, gin.H{
					"error": "Content-Type must be application/json",
				})
				return
			}
		}
		c.Next()
	}
}

// ApplySecurityMiddlewares registers all security middlewares on the given engine.
func ApplySecurityMiddlewares(router *gin.Engine) {
	router.Use(SecurityHeaders())
	router.Use(CORS(DefaultCORSConfig()))
	router.Use(RateLimiter(100, 1*time.Minute)) // 100 requests per minute
	router.Use(MaxBodySize(1 << 20))             // 1 MB
	router.Use(ContentTypeCheck())
}

// formatInt converts an int to its string representation without importing strconv.
func formatInt(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	negative := n < 0
	if negative {
		n = -n
	}
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	if negative {
		result = "-" + result
	}
	return result
}
