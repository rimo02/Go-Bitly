package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/patrickmn/go-cache"
	"github.com/rimo02/url-shortener/controller"
	"github.com/rimo02/url-shortener/database"
	"github.com/rimo02/url-shortener/routes"
)

func StartExpiryCheck() {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		controller.DeleteExpiredUrls()
	}
}
func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	database.Cache = cache.New(1*time.Hour, 3*time.Minute)
	database.ConnectDB()
}
func main() {
	r := gin.Default()
	go StartExpiryCheck()
	routes.RegisterRoutes(r)
	err := r.Run(":" + os.Getenv("PORT"))
	if err != nil {
		fmt.Println("Error running at port ", os.Getenv("PORT"), " error = ", err.Error())
	}
}
