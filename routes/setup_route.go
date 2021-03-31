package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/papertrader-api/models"
)

func (ec *EndpointContext) SetupTestData(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("SetupTestData\n")
	res := map[string]int{}
	user := models.User{
		ID:      models.GenerateID(),
		Name:    "system",
		Created: time.Now(),
		Email:   "system@papertrader.io",
	}
	game := models.Game{
		ID:        models.GenerateID(),
		Name:      "default game",
		Cap:       1000000.00,
		End:       time.Now().AddDate(0, 1, 0),
		CreatedBy: user.ID,
	}
	portfolio := models.Portfolio{
		ID:      models.GenerateID(),
		Created: time.Now(),
		GameID:  game.ID,
		UserID:  user.ID,
		Cash:    game.Cap,
		Value:   game.Cap,
		Spent:   game.Cap,
	}
	var code int
	code, _, _ = ec.PostDB("users", user)
	res["user"] = code
	code, _, _ = ec.PostDB("games", game)
	res["game"] = code
	code, _, _ = ec.PostDB("portfolios", portfolio)
	res["portfolio"] = code
	w.Write([]byte(fmt.Sprintf("%v", res)))
}
