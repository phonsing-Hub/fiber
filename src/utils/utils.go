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
func CreateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
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

func Decoded(tokenString string) (string, error) {
	// ตรวจสอบความถูกต้องของ token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// ตรวจสอบว่ามีการใช้ signing method ที่ถูกต้อง (เช่น HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		// คืนค่า secret key เพื่อให้ตรวจสอบ signature
		return secret, nil
	})

	// หาก token ไม่ถูกต้อง ส่งกลับ error
	if err != nil {
		return "", err
	}

	// ดึงข้อมูลจาก claims (ส่วน payload)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// ตรวจสอบว่ามี user_id ใน claims หรือไม่
		if userID, ok := claims["user_id"].(string); ok {
			return userID, nil
		}
		return "", errors.New("user_id not found in token")
	}

	// คืนค่า error หาก token ไม่ถูกต้อง
	return "", errors.New("invalid token")
}
