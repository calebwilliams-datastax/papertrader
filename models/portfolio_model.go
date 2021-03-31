package models

import "time"

type Portfolio struct {
	ID      string    `json:"id"`
	Created time.Time `json:"created"`
	GameID  string    `json:"game_id"`
	UserID  string    `json:"user_id"`
	Cash    float64   `json:"cash"`
	Spent   float64   `json:"spent"`
	Value   float64   `json:"value"`
}

func (p *Portfolio) SetDefaults() {
	p.ID = GenerateID()
	p.Created = time.Now()
}
