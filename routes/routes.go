package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rimo02/url-shortener/controller"
	"github.com/rimo02/url-shortener/middleware"
)

var RegisterRoutes = func(c *gin.Engine) {
	c.POST("/shorten", controller.ShortenTheUrl)
	c.GET("/:shorturl", middleware.PerClientRateLimiter(), controller.GetTheUrl)
	c.GET("/dashboard/:shorturl", controller.Dashboard)
}
