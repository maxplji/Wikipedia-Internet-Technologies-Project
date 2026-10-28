package controllers

import (
	"net/http"
	"wikipedia/backend/models"

	"errors"

	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type DBS struct {
	DB *gorm.DB
}

func DBHandler(db *gorm.DB) *DBS {
	return &DBS{
		DB: db,
	}
}

func (h *DBS) Login(c *gin.Context) {

	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nieprawidlowe dane",
		})
		return
	}

	var user models.User

	err := h.DB.Where("email = ?", input.Email).Take(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Nieprawidlowy email lub haslo",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Blad serwera",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash),
		[]byte(input.Password),
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Nieprawidlowy email lub haslo",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})

}

func (h *DBS) Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required,min=3,max=18"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	err := c.ShouldBindJSON(&input)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nieprawidlowe dane",
		})
		return
	}

	if utf8.RuneCountInString(input.Password) < 15 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Haslo musi miec conajmniej 15 znakow",
		})
		return
	}

	var user models.User

	err = h.DB.Where("email = ?", input.Email).Take(&user).Error

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Uzytkownik o podanym emailu juz istnieje!!!",
		})
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Blad bazy danych",
		})
		return
	}

	err = h.DB.Where("username = ?", input.Username).Take(&user).Error

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Uzytkownik o podanym nicku juz istnieje!!!",
		})
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Blad bazy danych",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Blad zabezpieczania hasla",
		})
		return
	}

	user = models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hash),
	}

	err = h.DB.Create(&user).Error

}
