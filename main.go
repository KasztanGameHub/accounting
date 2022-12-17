package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	_ "github.com/joho/godotenv/autoload"
)

var PORT = os.Getenv("PORT")
var GOOGLE_CLIENT_ID = os.Getenv("GOOGLE_CLIENT_ID")
var SECRET = []byte(os.Getenv("SECRET"))

type accClaims struct {
	Email string `json:"email"`
	Avatar string `json:"avatar"`
	Name string `json:"name"`
	jwt.StandardClaims
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.POST("/auth/callback", handleCallback)
	r.GET("/auth/logout", handleLogout)
	r.GET("/me", handleMe)

	r.Run(":"+PORT)
}