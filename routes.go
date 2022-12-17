package main

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// POST /auth/callback
func handleCallback(c *gin.Context) {
	// get google credential
	cred := c.PostForm("credential")
	if cred == "" {
		c.JSON(401, gin.H{
			"ok": false,
		})
		return
	}

	// validate the google credential
	claims, err := validateGoogleJWT(cred)
	if err != nil {
		c.JSON(401, gin.H{
			"ok": false,
		})
		return
	}

	// make a token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accClaims{
		Email: claims.Email,
		Avatar: claims.Avatar,
		Name: claims.Name,
	})

	tokenString, _ := token.SignedString(SECRET)

	c.SetCookie("token", tokenString, 60*60*24, "/", "", false, true)
	c.JSON(200, gin.H{
		"ok": true,
	})
}

// GET /auth/logout
func handleLogout(c *gin.Context) {
	// remove the cookie
	c.SetCookie("token", "", -1, "/", "", false, true)

	c.JSON(200, gin.H{
		"ok": true,
	})
}

// GET /me
func handleMe(c *gin.Context) {
	// get a token
	cred, err := c.Cookie("token")
	if err != nil {
		c.JSON(401, gin.H{
			"ok": false,
		})
		return
	}

	// parse claims from the token
	token, err := jwt.ParseWithClaims(cred, &accClaims{}, func (t *jwt.Token) (interface{}, error) {
		return SECRET, nil
	})
	if err != nil {
		c.JSON(401, gin.H{
			"ok": false,
		})
		return
	}

	claims := token.Claims.(*accClaims)
	c.JSON(200, gin.H{
		"ok": true,
		"email": claims.Email,
		"avatar": claims.Avatar,
		"name": claims.Name,
	})
}