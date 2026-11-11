package controllers

import (
	"errors"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"

	"wikipedia/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func getArticleID(c *gin.Context) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nieprawidlowe ID artykulu",
		})
		return 0, false
	}

	return id, true
}

// GET /articles
func (h *DBS) GetArticles(c *gin.Context) {
	var articles []models.Article

	err := h.DB.
		Preload("Images").
		Order("id DESC").
		Limit(100).
		Find(&articles).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Nie udalo sie pobrac artykulow",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"articles": articles,
	})
}

// GET /articles/:id
func (h *DBS) GetArticle(c *gin.Context) {
	id, ok := getArticleID(c)
	if !ok {
		return
	}

	var article models.Article

	err := h.DB.
		Preload("Images").
		First(&article, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
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

	c.JSON(http.StatusOK, gin.H{
		"article": article,
	})
}

// POST /articles
func (h *DBS) CreateArticle(c *gin.Context) {
	var input struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nieprawidlowe dane",
		})
		return
	}

	input.Title = strings.TrimSpace(input.Title)
	input.Content = strings.TrimSpace(input.Content)

	titleLength := utf8.RuneCountInString(input.Title)

	if titleLength < 3 || titleLength > 200 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Tytul musi miec od 3 do 200 znakow",
		})
		return
	}

	if input.Content == "" ||
		utf8.RuneCountInString(input.Content) > 100000 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nieprawidlowa dlugosc tresci",
		})
		return
	}

	userID := c.GetUint64("userID")

	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Nie jestes zalogowany",
		})
		return
	}

	article := models.Article{
		Title:    input.Title,
		Content:  input.Content,
		AuthorID: uint(userID),
		Images:   []models.ArticleImage{},
	}

	if err := h.DB.Create(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Nie udalo sie utworzyc artykulu",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"article": article,
	})
}

// DELETE /articles/:id
func (h *DBS) DeleteArticle(c *gin.Context) {
	id, ok := getArticleID(c)
	if !ok {
		return
	}

	var article models.Article

	err := h.DB.Preload("Images").First(&article, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
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

	// Usunac moze tylko autor artykulu
	if uint64(article.AuthorID) != c.GetUint64("userID") {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Nie mozesz usunac cudzego artykulu",
		})
		return
	}

	// Usuniecie rekordow w jednej transakcji
	err = h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Where("article_id = ?", article.ID).
			Delete(&models.ArticleImage{}).Error; err != nil {
			return err
		}

		return tx.Delete(&article).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Nie udalo sie usunac artykulu",
		})
		return
	}

	// Usuniecie zdjec z dysku
	for _, img := range article.Images {
		if strings.HasPrefix(img.URL, "/uploads/") {
			filePath := filepath.Join("uploads", path.Base(img.URL))
			_ = os.Remove(filePath)
		}
	}

	c.Status(http.StatusNoContent)
}
