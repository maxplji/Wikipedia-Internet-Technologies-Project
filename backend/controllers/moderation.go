package controllers

import (
	"errors"
	"net/http"
	"strings"
	"unicode/utf8"

	"wikipedia/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Sprawdza w bazie, czy zalogowany użytkownik jest administratorem.
func (h *DBS) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		err := h.DB.
			Select("id", "is_admin").
			First(&user, c.GetUint64("userID")).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Uzytkownik nie istnieje",
			})
			return
		}

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Blad bazy danych",
			})
			return
		}

		if !user.IsAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Wymagane uprawnienia administratora",
			})
			return
		}

		c.Next()
	}
}

// Administrator pobiera artykuły oczekujące na zatwierdzenie.
func (h *DBS) GetPendingArticles(c *gin.Context) {
	var articles []models.Article

	err := h.DB.
		Preload("Images").
		Where("status = ?", "pending").
		Order("id ASC").
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

// Administrator zatwierdza artykuł lub jego edycję.
func (h *DBS) ApproveArticle(c *gin.Context) {
	id, ok := getArticleID(c)
	if !ok {
		return
	}

	result := h.DB.
		Model(&models.Article{}).
		Where("id = ? AND status = ?", id, "pending").
		Update("status", "published")

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Nie udalo sie zatwierdzic artykulu",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Nie znaleziono artykulu oczekujacego na akceptacje",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Artykul zostal opublikowany",
	})
}

// Autor pobiera własny artykuł do edycji, także oczekujący.
func (h *DBS) GetArticleForEdit(c *gin.Context) {
	id, ok := getArticleID(c)
	if !ok {
		return
	}

	var article models.Article

	err := h.DB.
		Preload("Images").
		Where("author_id = ?", c.GetUint64("userID")).
		First(&article, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Nie znaleziono Twojego artykulu",
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

// Autor zapisuje edycję i przesyła artykuł do ponownej akceptacji.
func (h *DBS) UpdateArticle(c *gin.Context) {
	id, ok := getArticleID(c)
	if !ok {
		return
	}

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
	contentLength := utf8.RuneCountInString(input.Content)

	if titleLength < 3 || titleLength > 200 ||
		contentLength == 0 || contentLength > 100000 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nieprawidlowa dlugosc tytulu lub tresci",
		})
		return
	}

	result := h.DB.
		Model(&models.Article{}).
		Where("id = ? AND author_id = ?", id, c.GetUint64("userID")).
		Updates(map[string]any{
			"title":   input.Title,
			"content": input.Content,
			"status":  "pending",
		})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Nie udalo sie zapisac edycji",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Nie znaleziono Twojego artykulu",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Edycja oczekuje na akceptacje",
	})
}
