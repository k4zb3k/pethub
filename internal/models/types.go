package models

import "time"

type Config struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Phone     string    `json:"phone"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type Token struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Token     string    `db:"token"`
	Expire    time.Time `db:"expire"`
	CreatedAt time.Time `db:"created_at"`
}

type Ads struct {
	ID          int    `json:"id"      db:"id"`
	UserID      int    `json:"user_id" db:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	PhotoPath   string `json:"photo_path"`
	TypeId      int    `json:"type_id" db:"type_id"`
	PetId       int    `json:"pet_id"  db:"pet_id"`
	CityId      int    `json:"city_id" db:"city_id"`
	Reward      int    `json:"reward" db:"reward"`
	Active      bool   `json:"active" db:"is_active"`
}
