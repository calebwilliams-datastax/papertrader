package models

import "time"

type User struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	Created time.Time `json:"created"`
}

func (u *User) SetDefaults() {
	u.ID = GenerateID()
	u.Created = time.Now()

}
