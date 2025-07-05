package ranking

import (
	"database/sql"
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	rankingtypes "github.com/PretendoNetwork/nex-protocols-go/v2/ranking/types"
	"time"
)

var insertCategoryStmt *sql.Stmt
var insertScoreStmt *sql.Stmt

func InsertRankingByPIDAndRankingScoreData(pid types.PID, rankingScoreData rankingtypes.RankingScoreData, uniqueID types.UInt64) error {
	now := time.Now()
	tx, err := Database.Begin()
	if err != nil {
		return nex.NewError(nex.ResultCodes.Core.SystemError, err.Error())
	}
	defer tx.Rollback()

	_, err = tx.Stmt(insertCategoryStmt).Exec(
		rankingScoreData.Category,
		rankingScoreData.OrderBy == 1,
		now,
	)
	if err != nil {
		return nex.NewError(nex.ResultCodes.Core.SystemError, err.Error())
	}

	_, err = tx.Stmt(insertScoreStmt).Exec(
		uniqueID,
		pid,
		rankingScoreData.Category,
		rankingScoreData.Groups,
		rankingScoreData.Score,
		rankingScoreData.Param,
		now,
	)
	if err != nil {
		return nex.NewError(nex.ResultCodes.Core.SystemError, err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return nex.NewError(nex.ResultCodes.Core.SystemError, err.Error())
	}

	return err
}

func initInsertScoreStmt() error {
	stmt, err := Database.Prepare(`
		INSERT INTO ranking.categories (category, golf_scoring, creation_date)
		VALUES ($1, $2, $3)
		ON CONFLICT (category) DO NOTHING
	`)
	if err != nil {
		return err
	}
	stmt2, err := Database.Prepare(`
		INSERT INTO ranking.scores (unique_id, owner_pid, category, groups, score, param, creation_date, update_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $7)
		ON CONFLICT (category, owner_pid, unique_id) DO UPDATE
			SET score = excluded.score, update_date = excluded.update_date;
	`)
	// todo: case rankingscoredata.updatemode == ifbetter where category.golf_scoring etc etc..
	if err != nil {
		return err
	}

	insertCategoryStmt = stmt
	insertScoreStmt = stmt2
	return nil
}
