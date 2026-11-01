package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"wikipedia/backend/models"

	"errors"

	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": fmt.Sprint(user.ID),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	secrets := os.Getenv("SECRET")
	tokenString, err := token.SignedString([]byte(secrets))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Nie udalo sie stworzyc tokenu",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, int((24 * time.Hour).Seconds()), "/", "", false, true)

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

	if len([]byte(input.Password)) > 72 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Haslo jest za dlugie",
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

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Nie udalo sie utworzyc konta",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Konto zostalo utworzone",
	})

}
