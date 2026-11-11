package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"wikipedia/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *DBS) UploadArticleImages(c *gin.Context) {
	id, ok := getArticleID(c)
	if !ok {
		return
	}

	var article models.Article

	err := h.DB.First(&article, id).Error

	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Artykul nie istnieje",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Blad bazy danych",
		})
		return
	}

	if uint64(article.AuthorID) != c.GetUint64("userID") {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Nie mozesz dodawac zdjec do tego artykulu",
		})
		return
	}

	c.Request.Body = http.MaxBytesReader(
		c.Writer,
		c.Request.Body,
		26<<20,
	)

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nieprawidlowe dane lub zbyt duzy upload",
		})
		return
	}

	files := form.File["images"]

	if len(files) == 0 || len(files) > 5 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Wyslij od 1 do 5 zdjec",
		})
		return
	}

	var count int64

	err = h.DB.Model(&models.ArticleImage{}).
		Where("article_id = ?", article.ID).
		Count(&count).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Blad bazy danych",
		})
		return
	}

	if count+int64(len(files)) > 10 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Artykul moze miec maksymalnie 10 zdjec",
		})
		return
	}

	var images []models.ArticleImage
	var savedPaths []string

	cleanup := func() {
		for _, filePath := range savedPaths {
			_ = os.Remove(filePath)
		}
	}

	for _, file := range files {
		if file.Size == 0 || file.Size > 5<<20 {
			cleanup()
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Zdjecie przekracza limit 5 MiB lub jest puste",
			})
			return
		}

		src, err := file.Open()
		if err != nil {
			cleanup()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Nie udalo sie odczytac zdjecia",
			})
			return
		}

		var header [512]byte
		n, readErr := src.Read(header[:])
		_ = src.Close()

		if readErr != nil && readErr != io.EOF {
			cleanup()
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Nieprawidlowy plik",
			})
			return
		}

		mimeType := http.DetectContentType(header[:n])

		var extension string

		switch mimeType {
		case "image/jpeg":
			extension = ".jpg"
		case "image/png":
			extension = ".png"
		case "image/webp":
			extension = ".webp"
		default:
			cleanup()
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Dozwolone sa tylko JPG, PNG i WebP",
			})
			return
		}

		randomBytes := make([]byte, 16)

		if _, err := rand.Read(randomBytes); err != nil {
			cleanup()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Blad generowania nazwy pliku",
			})
			return
		}

		filename := hex.EncodeToString(randomBytes) + extension
		destination := filepath.Join("uploads", filename)

		if err := c.SaveUploadedFile(file, destination); err != nil {
			cleanup()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Nie udalo sie zapisac zdjecia",
			})
			return
		}

		savedPaths = append(savedPaths, destination)

		images = append(images, models.ArticleImage{
			ArticleID: article.ID,
			URL:       "/uploads/" + filename,
		})
	}

	if err := h.DB.Create(&images).Error; err != nil {
		cleanup()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Nie udalo sie zapisac zdjec w bazie",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"images": images,
	})
}
