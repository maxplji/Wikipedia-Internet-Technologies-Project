package main

import (
	"log"
	"os"
	"wikipedia/backend/controllers"
	"wikipedia/backend/middlewares"
	"wikipedia/backend/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		log.Fatal("Brak DATABASE_URL w zmiennych srodowiskowych")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Nie udalo sie polaczyc z baza danych: ", err)
	}

	return db
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("Nie udalo sie wczytac .env, uzywam zmiennych systemowych")
	}

	if os.Getenv("SECRET") == "" {
		log.Fatal("Brak klucza SECRET")
	}

	db := ConnectDatabase()

	err := db.AutoMigrate(
		&models.User{},
		&models.Article{},
		&models.ArticleImage{},
	)

	if err != nil {
		log.Fatal("Blad migracji bazy danych: ", err)
	}

	if err := os.MkdirAll("./uploads", 0755); err != nil {
		log.Fatal("Nie udalo sie utworzyc katalogu uploads: ", err)
	}

	r := gin.New()

	r.Use(
		gin.Recovery(),
		middlewares.Logger(),
	)

	frontendOrigin := os.Getenv("FRONTEND_ORIGIN")

	if frontendOrigin != "" {
		r.Use(cors.New(cors.Config{
			AllowOrigins: []string{
				frontendOrigin,
			},
			AllowMethods: []string{
				"GET",
				"POST",
				"DELETE",
				"OPTIONS",
			},
			AllowHeaders: []string{
				"Content-Type",
			},
			AllowCredentials: true,
		}))
	}

	r.MaxMultipartMemory = 8 << 20

	// Udostepnianie zapisanych zdjec
	r.Static("/uploads", "./uploads")

	handler := controllers.DBHandler(db)

	// Logowanie i rejestracja
	r.POST("/login", handler.Login)
	r.POST("/register", handler.Register)

	// Endpointy artykulow
	articles := r.Group("/articles")
	{

		articles.GET("", handler.GetArticles)
		articles.GET("/:id", handler.GetArticle)

		articles.POST(
			"",
			middlewares.JWTAuth(),
			handler.CreateArticle,
		)

		articles.POST(
			"/:id/images",
			middlewares.JWTAuth(),
			handler.UploadArticleImages,
		)

		articles.DELETE(
			"/:id",
			middlewares.JWTAuth(),
			handler.DeleteArticle,
		)
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	log.Println("Serwer uruchamiany na porcie:", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
