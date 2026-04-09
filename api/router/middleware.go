package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joaop/psiencontra/api/handler"
)

const (
	// userContextKey is the gin.Context key under which the authenticated
	// user ID is stored. The reader lives in handler.UserIDFromContext.
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
	token := handler.ReadToken(c)
	if token == "" {
		return uuid.Nil, false
	}
	claims, err := handler.AuthSvc.ParseToken(token)
	if err != nil {
		return uuid.Nil, false
	}
	return claims.UserID, true
}
