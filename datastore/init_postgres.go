package datastore

import (
	"github.com/PretendoNetwork/puyo-puyo-tetris/globals"
	"os"
)

func initRankingTable() {
	_, err := globals.Postgres.Exec(`CREATE SCHEMA IF NOT EXISTS ranking`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		os.Exit(0)
	}

	globals.Logger.Success("ranking globals.Postgres schema created")

	_, err = globals.Postgres.Exec(`CREATE SEQUENCE IF NOT EXISTS ranking.unique_id_seq
		INCREMENT 1
		MINVALUE 1
		MAXVALUE 281474976710656
		START 1
		CACHE 1`, // * Honestly I don't know what CACHE does but I saw it recommended so here it is
	)
	if err != nil {
		globals.Logger.Critical(err.Error())
		os.Exit(0)
	}

	_, err = globals.Postgres.Exec(`CREATE TABLE IF NOT EXISTS ranking.scores (
		unique_id bigint NOT NULL DEFAULT nextval('ranking.unique_id_seq') PRIMARY KEY,
		deleted boolean NOT NULL DEFAULT FALSE,
		user_pid bigint,
		category int,
		groups bytea,
		score int,
		param bigint,
		common_data bytea,
		creation_date timestamp,
		update_date timestamp
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		os.Exit(0)
	}

	_, err = globals.Postgres.Exec(`CREATE INDEX IF NOT EXISTS ranking.score_index
		ON ranking.scores (score)
	`)

	_, err = globals.Postgres.Exec(`CREATE INDEX IF NOT EXISTS ranking.category_index
		ON ranking.scores (category)
	`)

	_, err = globals.Postgres.Exec(`CREATE TABLE IF NOT EXISTS ranking.categories (
		category int NOT NULL PRIMARY KEY,   
    	golf_scoring boolean,
		creation_date timestamp                 
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		os.Exit(0)
	}

	globals.Logger.Success("Postgres tables created")
}
