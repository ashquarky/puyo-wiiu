package datastore

import (
	"database/sql"
	"os"

	"github.com/PretendoNetwork/puyo-puyo-tetris/globals"
	_ "github.com/lib/pq"
)

func ConnectPostgres() {
	var err error

	globals.Postgres, err = sql.Open("postgres", os.Getenv("PN_PUYOPUYOTETRIS_POSTGRES_URI"))
	if err != nil {
		globals.Logger.Critical(err.Error())
	}

	globals.Logger.Success("Connected to Postgres!")

	initRankingTable()
	initRanking()
}
