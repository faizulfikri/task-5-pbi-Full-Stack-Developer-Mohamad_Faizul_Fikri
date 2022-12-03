package model

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/faizulfikri/task-5-vix-btpn-Mohamad_Faizul_Fikri/app"
	"github.com/faizulfikri/task-5-vix-btpn-Mohamad_Faizul_Fikri/helper/helperhash"
	"github.com/google/uuid"
)

type User struct {
	ID       string    `gorm:"primaryKey; unique; not null" json:"id"`
	Username string    `gorm:"size:255; not null" json:"username"`
	Email    string    `gorm:"size:255; not null; unique" json:"email"`
	Password string    `gorm:"size:255; not null" json:"password"`
	Photos   Photo     `gorm:"constraint:OnUpdate:CASCADE, OnDelete:SET NULL;" json:"photos"`
	CreateAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"create_at"`
	UpdateAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"update_at"`
}

type Photo struct {
	ID       int       `gorm:"primaryKey; autoIncrement" json:"id"`
	Title    string    `gorm:"size:255; not null" json:"title"`
	Caption  string    `gorm:"size:255; not null" json:"caption"`
	PhotoUrl string    `gorm:"size:255; not null" json:"photo_url"`
	UserID   string    `gorm:"not null" json:"user_id"`
	Owner    app.Owner `gorm:"owner"`
}

// User Methods
func (u *User) Init() {
	u.ID = uuid.New().String()
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
}

func (u *User) HashPassword() error {
	hashPass, err := helperhash.GetHashPassword(u.Password)
	if err != nil {
		return err
	}

	u.Password = string(hashPass)
	return nil

}

func (u *User) CheckPassword(cek string) error {
	err := helperhash.ComparePassword(u.Password, cek)
	if err != nil {
		return err
	}

	return nil
}

// Validate data
func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {

	case "register":
		if u.ID == "" {
			return errors.New("ID is required")
		} else if u.Email == "" {
			return errors.New("Email is required")
		} else if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Email is invalid")
		} else if u.Password == "" {
			return errors.New("Password is required")
		} else if length := len([]rune(u.Password)); length < 6 {
			return errors.New("Password must be at least 6 character")
		} else if u.Username == "" {
			return errors.New("Username is required")
		}

		return nil

	case "login":
		if u.Email == "" {
			return errors.New("Email is required")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Email is invalid")
		}
		if u.Password == "" {
			return errors.New("Password is required")
		}

		return nil

	default:
		return nil

	}
}

//Photo Method

func (p *Photo) Init() {
	p.Title = html.EscapeString(strings.TrimSpace(p.Title)) //Escape string
	p.Caption = html.EscapeString(strings.TrimSpace(p.Caption))
	p.PhotoUrl = html.EscapeString(strings.TrimSpace(p.PhotoUrl))
}

// Validate data
func (p *Photo) Validate(action string) error {
	switch strings.ToLower(action) {

	case "upload":
		if p.Title == "" {
			return errors.New("Title is required")
		} else if p.Caption == "" {
			return errors.New("Caption is required")
		} else if p.UserID == "" {
			return errors.New("User ID is required")
		}

		return nil

	case "change":
		if p.Title == "" {
			return errors.New("Title is required")
		} else if p.Caption == "" {
			return errors.New("Caption is required")
		} else if p.PhotoUrl == "" {
			return errors.New("Url is required")
		}

		return nil

	default:
		return nil
	}
}
