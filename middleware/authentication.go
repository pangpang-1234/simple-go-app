package middleware

import (
	"errors"
	"log"
	"os"
	"simplegoapp/config"
	"simplegoapp/models"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type login struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

var identitykey = "sub"

func Authenticate() *jwt.GinJWTMiddleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		// secret key
		Key: []byte(os.Getenv("SECRET_KEY")),

		IdentityKey: identitykey,
		// login => user
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var form login
			var user models.User

			if err := c.ShouldBindJSON(&form); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}

			db := config.GetDB()
			if err := db.Where("email = ?", form.Email).First(&user).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound){	
				return nil, jwt.ErrFailedAuthentication
			}
			if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			return &user, nil
		},
		// user => payload(sub) => JWT
		PayloadFunc: func (data interface{}) jwt.MapClaims {
			if v, ok :=	data.(*models.User); ok {
				claims := jwt.MapClaims{
					identitykey: v.ID,
				}
				return claims
			}	
			return jwt.MapClaims{}
		},

		// handle error
		Unauthorized: func (c *gin.Context, code int, message string)  {
			c.JSON(code, gin.H{
				"error": message,
			})
		},
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	return authMiddleware
}
