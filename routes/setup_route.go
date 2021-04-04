package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/papertrader-api/models"
	"github.com/papertrader-api/util"
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

	ucode, udata, uerr := ec.PostDB("users", user)
	res["user"] = fmt.Sprintf(`%v:%s:%s`, ucode, udata, uerr)
	gcode, gdata, gerr := ec.PostDB("games", game)
	res["games"] = fmt.Sprintf(`%v:%s:%s`, gcode, gdata, gerr)
	pcode, pdata, perr := ec.PostDB("portfolios", portfolio)
	res["portfolios"] = fmt.Sprintf(`%v:%s:%s`, pcode, pdata, perr)
	w.Write([]byte(fmt.Sprintf("%v", res)))
}
