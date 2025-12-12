package user

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

var (
	ErrInvalidPhoneFormat = errors.New("phone must start with +998 and have correct length")
)

type PhoneNumber string

func NewPhoneNumber(phone string) (PhoneNumber, error) {
	phone = strings.TrimSpace(phone)
	if !strings.HasPrefix(phone, "+998") {
		return "", fmt.Errorf("%w: phone must start with +998", ErrInvalidPhoneFormat)
	}

	if len(phone) != 13 {
		return "", fmt.Errorf("%w: phone must be 13 characters (format: +998XXXXXXXXX)", ErrInvalidPhoneFormat)
	}

	for _, r := range phone[4:] {
		if !unicode.IsDigit(r) {
			return "", fmt.Errorf("%w: phone must contain only digits after +998", ErrInvalidPhoneFormat)
		}
	}

	return PhoneNumber(phone), nil
}

func (p PhoneNumber) String() string {
	return string(p)
}

type User struct {
	ID           string      `json:"id" db:"id"`
	PhoneNumber  PhoneNumber `json:"phone_number" db:"phone_number"`
	PasswordHash string      `json:"-" db:"password_hash"`
	Name         string      `json:"name" db:"name"`
	CreatedAt    time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at" db:"updated_at"`
}

func (u *User) GetId() string {
	return u.ID
}

func NewUser(phoneNumber PhoneNumber, passwordHash, name string) *User {
	return &User{
		ID:           uuid.New().String(),
		PhoneNumber:  phoneNumber,
		PasswordHash: passwordHash,
		Name:         name,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
