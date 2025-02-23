package main

import (
	"context"
	authpb "goshop/api/protobuf/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BearerSchema = "Bearer "
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) <= len(BearerSchema) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			c.Abort()
			return
		}

		tokenString := authHeader[len(BearerSchema):]
		ret, err := AuthClient().VerifyToken(context.Background(), &authpb.ReqVerifyToken{
			Token:    tokenString,
			IsAccess: true,
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}
		c.Set("user_id", ret.UserId)

		c.Next()
	}
}
