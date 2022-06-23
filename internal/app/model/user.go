package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

type Role struct {
	ID int
}

type Request struct {
	Date1 string `json:"date1"`
	Date2 string `json:"date2"`
	Token string `json:"token"`
	Line  int    `json:"line"`
}

type TokenPerson struct {
	Name string
	Item int
	Role string
}

type User struct {
	ID                int    `json:"id"`
	Email             string `json:"email"`
	Password          string `json:"password,omitempty"`
	EncryptedPassword string `json:"-"`
	Role              string `json:"role,omitempty"`
}

type Token struct {
	Role  string `json:"role"`
	Email string `json:"email"`
	Token string `json:"token"`
}

type RegisterUser struct {
	Role        string `json:"role"`
	Email       string `json:"email"`
	Token       string `json:"token"`
	Regemail    string `json:"regemail"`
	Regpassword string `json:"regpassword"`
}

// Validate ...
func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.By(requiredIf(u.EncryptedPassword == "")), validation.Length(6, 100)),
	)
}

// BeforeCreate ...
func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		enc, err := encryptString(u.Password)
		if err != nil {
			return err
		}

		u.EncryptedPassword = enc
	}

	return nil
}

// Sanitize ...
func (u *User) Sanitize() {
	u.Password = ""
}

// ComparePassword ...
func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
