package response

import "github.com/gin-gonic/gin"

type Error struct {
	Error string `json:"error"`
}

func WriteError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Error{
		Error: message,
	})
}
