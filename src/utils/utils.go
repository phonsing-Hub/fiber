package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
	// "fmt"
	// "encoding/json"
)

var secret = []byte(os.Getenv("JWT_SECRET_KEY"))

// HashPassword hashes the given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash checks if the provided password matches the hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// CreateJWT creates a JWT token
func CreateJWT(ID uint, userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": ID,
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})

	return token.SignedString(secret)
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// ตรวจสอบว่า algorithm ที่ใช้เป็น "HS256"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}
	return token, nil
}

func Decoded(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})

	if err != nil {
		return 0, err
	}

	// ดึงข้อมูลจาก claims (ส่วน payload)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// ดึงค่า id จาก claims และแปลงเป็น uint
		if idFloat, ok := claims["id"].(float64); ok {
			return uint(idFloat), nil
		}
		return 0, errors.New("id not found in token")
	}

	// คืนค่า error หาก token ไม่ถูกต้อง
	return 0, errors.New("invalid token")
}
