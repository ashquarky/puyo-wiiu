package ranking

import (
	"database/sql"
	"github.com/PretendoNetwork/puyo-puyo-tetris/globals"
)

var Database *sql.DB

func initDatabase() error {
	inits := []func() error{
		initTables,
		initInsertScoreStmt,
		initGetGlobalRankingsStmt,
		initInsertCommonDataStmt,
	}

	for _, init := range inits {
		err := init()
		if err != nil {
			return err
		}
	}

	return nil
}

func initTables() error {
	_, err := Database.Exec(`CREATE SCHEMA IF NOT EXISTS ranking`)
	if err != nil {
		return err
	}

	globals.Logger.Success("ranking Postgres schema created")

	_, err = Database.Exec(`CREATE TABLE IF NOT EXISTS ranking.scores (
		deleted boolean NOT NULL DEFAULT FALSE,
		unique_id bigint,
		owner_pid bigint,
		category int,
		groups bytea,
		score int,
		param bigint,
		creation_date timestamp,
		update_date timestamp,
		PRIMARY KEY (category, owner_pid, unique_id)
	)`)
	if err != nil {
		return err
	}

	_, err = Database.Exec(`CREATE INDEX IF NOT EXISTS score_index
		ON ranking.scores (score)
	`)

	_, err = Database.Exec(`CREATE INDEX IF NOT EXISTS category_index
		ON ranking.scores (category)
	`)

	_, err = Database.Exec(`CREATE TABLE IF NOT EXISTS ranking.categories (
		category int PRIMARY KEY,   
    	golf_scoring boolean,
		creation_date timestamp                 
	)`)
	if err != nil {
		return err
	}

	_, err = Database.Exec(`CREATE TABLE IF NOT EXISTS ranking.common_data (
    	deleted boolean NOT NULL DEFAULT FALSE,
		unique_id bigint,
		owner_pid bigint,
    	data bytea,
		creation_date timestamp,
		update_date timestamp,
		PRIMARY KEY (owner_pid, unique_id)
	)`)
	if err != nil {
		return err
	}

	globals.Logger.Success("Postgres tables created")

	return nil
}
