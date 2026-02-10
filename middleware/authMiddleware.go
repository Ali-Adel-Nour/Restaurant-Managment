package middleware

import (
	"net/http"

	"github.com/ali-adel-nour/restaurant-management/helpers"
	"github.com/gin-gonic/gin"
)

// Authentication is a middleware that validates JWT tokens
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from header
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "No authorization header provided",
			})
			c.Abort()
			return
		}

		// Validate token
		claims, err := helpers.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err,
			})
			c.Abort()
			return
		}

		// Set claims in context for use in handlers
		c.Set("email", claims.Email)
		c.Set("first_name", claims.FirstName)
		c.Set("last_name", claims.LastName)
		c.Set("uid", claims.UID)

		c.Next()
	}
}
