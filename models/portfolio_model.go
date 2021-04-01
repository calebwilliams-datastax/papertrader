package models

import "time"

type Portfolio struct {
	ID      string    `json:"id"`
	Created time.Time `json:"created"`
	GameID  string    `json:"game_id"`
	UserID  string    `json:"user_id"`
	Cash    string   `json:"cash"`
	Spent   string   `json:"spent"`
	Value   string   `json:"value"`
}

type APIPortfolioResponse struct {
	Count int         `json:"count"`
	Data  []Portfolio `json:"data"`
}

func (p *Portfolio) SetDefaults() {
	p.ID = GenerateID()
	p.Created = time.Now()
}
