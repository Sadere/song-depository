package main

import "github.com/gin-gonic/gin"

// This is a mock music info service for testing purposes
func main() {
	r := gin.Default()
	r.GET("/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"releaseDate": "22.06.2019",
			"text":        "Ooh baby, don't you know I suffer?\\nOoh baby, can you hear me moan?\\nYou caught me under false pretenses\\nHow long before you let me go?\\n\\nOoh\\nYou set my soul alight\\nOoh\\nYou set my soul alight",
			"link":        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
		})
	})
	_ = r.Run(":8081")
}
