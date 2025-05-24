package controllers

import (
	"net/http"
	"time"
	"todo-api/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var jwtKey = []byte("secret") // Depois eu coloco um os.Getenv("JWT_KEY")

type Claims struct {
	UserID uint
	jwt.StandardClaims
}

func Register(c *gin.Context, db *gorm.DB) {
	var input struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
		return
	}

	var exinsting models.User
	if err := db.Where("email = ? OR username = ?", input.Email, input.Username).First(&exinsting).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or Username already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not ecrypt password"})
		return
	}

	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User Created"})
}

func Login(c *gin.Context, db *gorm.DB) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user models.User
	if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email"})
		return
	}

	pass := []byte(input.Password)
	hashedPassword := []byte(user.Password)

	if err := bcrypt.CompareHashAndPassword(hashedPassword, pass); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication Error"})
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generateToken"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenStr})
}

func Test(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"funciona": jwt.SigningMethodHS256.Alg(), "outro": jwt.SigningMethodES256.Alg()})
}
