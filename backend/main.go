package main

import (
	"log"
	"os"
	"wikipedia/backend/controllers"
	"wikipedia/backend/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func main() {
	r := gin.New()
	db := ConnectDatabase()

	handler := controllers.DBHandler(db)
	r.Use(gin.Recovery(), middlewares.Logger())

	r.POST("/login", handler.Login)
	//r.POST("/register", handler.Register(db))

	//r.GET("/articles", handler.GetArticles(db))
	//r.POST("/articles", handler.CreateArticle(db))

	if err := r.Run(":3000"); err != nil {
		log.Fatal(err)
	}
}
