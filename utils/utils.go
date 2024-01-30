package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"os"
	"reflect"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func GenerateResetToken() (string, error) {
	uid := make([]byte, 16)
	_, err := rand.Read(uid)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(uid), nil
}

func SendEmail(to_receiver string, msg string) {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	from := os.Getenv("MAIL_FROM")
	password := os.Getenv("MAIL_PASSWORD")
	auth := smtp.PlainAuth("", from, password, smtpHost)
	to := []string{to_receiver}
	message := []byte(msg)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Function to update fields in a struct based on map values
func UpdateFields(data interface{}, updateData map[string]interface{}) {
	valueOf := reflect.ValueOf(data).Elem()
	for key, value := range updateData {
		field := valueOf.FieldByName(key)
		if field.IsValid() {
			field.Set(reflect.ValueOf(value))
		}
	}
}
