package models

import "time"

type User struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	Created time.Time `json:"created"`
}

type APIUserResponse struct {
	Count int    `json:"count"`
	Data  []User `json:"data"`
}

func (u *User) SetDefaults() {
	u.ID = GenerateID()
	u.Created = time.Now()
}
