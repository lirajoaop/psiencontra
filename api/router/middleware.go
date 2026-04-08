package router

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joaop/psiencontra/api/handler"
)

const (
	// AuthCookieName is the name of the cookie carrying the JWT.
	AuthCookieName = "psiencontra_auth"
	// userContextKey is the gin.Context key under which the authenticated
	// user ID is stored.
	userContextKey = "user_id"
)

// OptionalAuth populates the request context with a user ID when a valid JWT
// is present (cookie or Authorization header), but does not block the request
// when authentication is missing or invalid. Used by endpoints that work both
// for anonymous and logged-in users (e.g. creating a session).
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if userID, ok := extractUserID(c); ok {
			c.Set(userContextKey, userID)
		}
		c.Next()
	}
}

// RequireAuth aborts the request with 401 if no valid JWT is present.
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := extractUserID(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Set(userContextKey, userID)
		c.Next()
	}
}

func extractUserID(c *gin.Context) (uuid.UUID, bool) {
	token := readToken(c)
	if token == "" {
		return uuid.Nil, false
	}
	claims, err := handler.AuthSvc.ParseToken(token)
	if err != nil {
		return uuid.Nil, false
	}
	return claims.UserID, true
}

func readToken(c *gin.Context) string {
	if cookie, err := c.Cookie(AuthCookieName); err == nil && cookie != "" {
		return cookie
	}
	if h := c.GetHeader("Authorization"); strings.HasPrefix(h, "Bearer ") {
		return strings.TrimPrefix(h, "Bearer ")
	}
	return ""
}

// UserIDFromContext returns the authenticated user ID, if any.
func UserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	v, exists := c.Get(userContextKey)
	if !exists {
		return uuid.Nil, false
	}
	id, ok := v.(uuid.UUID)
	return id, ok
}
