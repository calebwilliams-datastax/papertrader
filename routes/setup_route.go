package routes

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/calebwilliams-datastax/papertrader-api/models"
	"github.com/calebwilliams-datastax/papertrader-api/util"
)

func (ec *EndpointContext) SetupTestData(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("SetupTestData\n")
	util.LogRequest(r)
	ec.RefreshAuthToken()
	res := map[string]string{}

	user := models.User{
		ID:      models.GenerateID(),
		Name:    "system",
		Created: time.Now(),
		Email:   "system@papertrader.io",
	}

	game := models.Game{
		ID:        models.GenerateID(),
		Name:      "default game",
		Cap:       "1000000.00",
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

	symbols := []string{"GME", "SPY", "AMC", "SCHW"}
	orders := []models.Order{}
	for i := 1; i < 30; i++ {
		order := models.Order{
			ID:          models.GenerateID(),
			UserID:      user.ID,
			PortfolioID: portfolio.ID,
			Created:     time.Now().AddDate(0, 0, -rand.Intn(20)),
			Symbol:      symbols[rand.Intn(len(symbols))],
			OrderAction: models.Buy,
			Amount:      rand.Intn(3000),
			Ask:         "",
		}
		orders = append(orders, order)
	}

	ucode, udata, uerr := ec.PostDB("users", user)
	res["user"] = fmt.Sprintf(`%v:%s:%s`, ucode, udata, uerr)
	gcode, gdata, gerr := ec.PostDB("games", game)
	res["games"] = fmt.Sprintf(`%v:%s:%s`, gcode, gdata, gerr)
	pcode, pdata, perr := ec.PostDB("portfolios", portfolio)
	res["portfolios"] = fmt.Sprintf(`%v:%s:%s`, pcode, pdata, perr)
	for _, o := range orders {
		ocode, odata, oerr := ec.PostDB("orders", o)
		key := fmt.Sprintf("order:%s", o.ID)
		res[key] = fmt.Sprintf(`%v:%s:%s`, ocode, odata, oerr)
	}
	w.Write([]byte(fmt.Sprintf("%v", res)))
}
