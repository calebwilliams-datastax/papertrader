package models

import "time"

type Game struct {
	ID        string    `json:"id"`
	Created   time.Time `json:"created"`
	CreatedBy string    `json:"created_by"`
	Name      string    `json:"name"`
	End       time.Time `json:"end"`
	Cap       float64   `json:"cap"`
}

func (g *Game) SetDefaults() {
	g.ID = GenerateID()
	g.Created = time.Now()
}