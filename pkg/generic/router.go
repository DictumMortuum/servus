package generic

import (
	"github.com/DictumMortuum/servus/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"time"
)

// SetConfig gin Middlware to push some config values
func SetConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("CorsOrigin", "*")
		c.Set("Verbose", true)
		c.Next()
	}
}

// Options common response for rest options
func Options(c *gin.Context) {
	Origin := c.MustGet("CorsOrigin").(string)

	c.Writer.Header().Set("Access-Control-Allow-Origin", Origin)
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,DELETE,POST,PUT")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Next()
}

func SetupMainRouter() (*gin.Engine, error) {
	err := config.Read()
	if err != nil {
		return nil, err
	}

	r := gin.Default()
	r.Use(SetConfig())
	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type, Bearer, range",
		ExposedHeaders:  "x-total-count, Content-Range",
		MaxAge:          50 * time.Second,
		Credentials:     false,
		ValidateHeaders: false,
	}))

	return r, nil
}
