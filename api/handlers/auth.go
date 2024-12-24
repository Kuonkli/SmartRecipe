package handlers

import (
	"SmartRecipe/database"
	"SmartRecipe/models"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"time"
)

var jwtKey = []byte("j2f3l2kd658vK-09S=GHksvSKGLSDGJ")

type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

type SignUpRequest struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignUp(c *gin.Context) {
	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user := &models.User{
		Name:     req.Name,
		Surname:  req.Surname,
		Email:    req.Email,
		Password: string(hashedPassword),
		RoleId:   1,
	}
	newUserId, err := database.Database.Users.AddUser(user)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println(err)
		return
	}

	accessTokenExpirationTime := time.Now().Add(15 * time.Minute)
	refreshTokenExpirationTime := time.Now().Add(24 * time.Hour)

	accessClaims := &Claims{
		UserID: strconv.Itoa(newUserId),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessTokenExpirationTime.Unix(),
		},
	}

	refreshClaims := &Claims{
		UserID: strconv.Itoa(newUserId),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshTokenExpirationTime.Unix(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}
	// Сохранение токенов в куки
	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    accessTokenString,
		Expires:  accessTokenExpirationTime,
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	}
	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshTokenString,
		Expires:  refreshTokenExpirationTime,
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	}
	http.SetCookie(c.Writer, accessCookie)
	http.SetCookie(c.Writer, refreshCookie)

	c.Redirect(http.StatusSeeOther, "/c/smart_recipe")
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	user, err := database.Database.Users.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password or email"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password or email"})
		return
	}

	accessTokenExpirationTime := time.Now().Add(15 * time.Minute)
	refreshTokenExpirationTime := time.Now().Add(24 * time.Hour)

	accessClaims := &Claims{
		UserID: strconv.Itoa(user.Id),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessTokenExpirationTime.Unix(),
		},
	}

	refreshClaims := &Claims{
		UserID: strconv.Itoa(int(user.Id)),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshTokenExpirationTime.Unix(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Сохранение access и refresh токенов в куки
	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    accessTokenString,
		Expires:  accessTokenExpirationTime,
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	}
	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshTokenString,
		Expires:  refreshTokenExpirationTime,
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	}
	http.SetCookie(c.Writer, accessCookie)
	http.SetCookie(c.Writer, refreshCookie)

	c.Redirect(http.StatusSeeOther, "/c/smart_recipe")
}

func Logout(c *gin.Context) {
	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Устанавливаем срок действия в прошлое
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	}
	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Устанавливаем срок действия в прошлое
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	}
	http.SetCookie(c.Writer, accessCookie)
	http.SetCookie(c.Writer, refreshCookie)

	c.Redirect(http.StatusSeeOther, "/o/auth")
}

func RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.Redirect(http.StatusFound, "/o/auth")
		return
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		c.Redirect(http.StatusFound, "/o/auth")
		return
	}

	newAccessTokenExpirationTime := time.Now().Add(15 * time.Minute)
	newAccessClaims := &Claims{
		UserID: claims.UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: newAccessTokenExpirationTime.Unix(),
		},
	}

	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newAccessClaims)
	newAccessTokenString, err := newAccessToken.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new access token"})
		return
	}

	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    newAccessTokenString,
		Expires:  newAccessTokenExpirationTime,
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	}
	http.SetCookie(c.Writer, accessCookie)
	redirectURI := c.DefaultQuery("redirect_uri", "/")
	c.Redirect(http.StatusFound, redirectURI)
}
