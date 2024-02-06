package controllers

import (
	"net/http"
	"simplegoapp/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type Auth struct {
	DB *gorm.DB
}

type authForm struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type authResponse struct {
	ID string `json:"id"`
	Email string `json:"email"`
}

func (a *Auth) Signup(ctx *gin.Context) {
	var form authForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	copier.Copy(&user, &form) 
	user.Password = user.GenerateEncryptedPassword()
	if err := a.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	println(&user)

	var serializeUser authResponse
	copier.Copy(&serializeUser, &user)
	ctx.JSON(http.StatusCreated, gin.H{"user": serializeUser})
}