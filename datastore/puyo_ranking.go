package datastore

import (
	"database/sql"
	"errors"
	"fmt"
	nextypes "github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/nex-protocols-go/v2/ranking/types"
	"github.com/PretendoNetwork/puyo-puyo-tetris/globals"
	"github.com/lib/pq"
	"strconv"
	"time"
)

// https://medium.com/analytics-vidhya/leaderboards-and-rankings-with-sql-f0c7700d41d3

//todo
// filter by groups
// unique_id is something else
// orderBy

/*
	  assumes SELECT
						user_pid,
						unique_id,
						score,
						groups,
						param,
						common_data,
						update_date
*/
func parseRankingDataList(rows *sql.Rows) (*nextypes.List[*types.RankingRankData], error) {
	results := nextypes.NewList[*types.RankingRankData]()

	for rows.Next() {
		result := types.NewRankingRankData()
		var updateDate time.Time
		var userPid uint64

		err := rows.Scan(
			&userPid,
			&result.UniqueID.Value,
			&result.Score.Value,
			&result.Groups.Value,
			&result.Param.Value,
			&result.CommonData.Value,
			&result.Order.Value,
			&updateDate,
		)
		if err != nil {
			return nil, err
		}

		result.PrincipalID = nextypes.NewPID(userPid)
		result.UpdateTime.FromTimestamp(updateDate)

		results.Append(result)
	}

	return results, nil
}

func isUndefinedTable(err error) bool {
	var pqErr *pq.Error
	// 42P01 is "undefined_table"
	return err != nil && errors.As(err, &pqErr) && pqErr.SQLState() == "42P01"
}

func GetRankingsAndCountByCategoryAndRankingOrderParam(category *nextypes.PrimitiveU32, rankingOrderParam *types.RankingOrderParam) (*nextypes.List[*types.RankingRankData], uint32, error) {
	globals.Logger.Info(rankingOrderParam.FormatToString(1))
	// todo ordinal ranking
	rankingTable := `ranking.ranks_` + strconv.Itoa(int(category.Value))
	rows, err := Postgres.Query(`
		SELECT
		    rt.user_pid,
		    rt.unique_id,
		    rt.score,
		    rt.groups,
		    rt.param,
		    rt.common_data,
		    rt.rank,
		    rt.update_date
		FROM `+rankingTable+` as rt
		LIMIT $1
		OFFSET $2
	`,
		rankingOrderParam.Length.Value,
		rankingOrderParam.Offset.Value,
	)
	// undefined table is expected if rankingTable isn't existing
	if errors.Is(err, sql.ErrNoRows) || isUndefinedTable(err) {
		return nil, 0, nil
	} else if err != nil {
		return nil, 0, err
	}

	results, err := parseRankingDataList(rows)
	if err != nil {
		return nil, 0, err
	}

	// todo totalCount
	return results, uint32(results.Length()), nil
}
func GetNearbyRankingsAndCountByCategoryAndRankingOrderParam(pid *nextypes.PID, category *nextypes.PrimitiveU32, rankingOrderParam *types.RankingOrderParam) (*nextypes.List[*types.RankingRankData], uint32, error) {
	globals.Logger.Infof("pid: %d cat: %d", pid.Value(), category.Value)
	globals.Logger.Info(rankingOrderParam.FormatToString(1))

	// todo ordinal ranking
	// https://stackoverflow.com/a/9852512
	rankingTable := `ranking.ranks_` + strconv.Itoa(int(category.Value))
	rows, err := Postgres.Query(`
		WITH central_user as (
		        SELECT
		            ordinal,
		            GREATEST(ordinal - $2 / 2, 1) as min_ord,
		            LEAST(ordinal + $2 / 2, (SELECT max(ordinal) FROM `+rankingTable+`)) as max_ord
		        FROM `+rankingTable+`
		        WHERE user_pid = $1
		    )
		SELECT
		    rt.user_pid,
		    rt.unique_id,
		    rt.score,
		    rt.groups,
		    rt.param,
		    rt.common_data,
		    rt.rank,
		    rt.update_date
		FROM `+rankingTable+` AS rt, central_user
		WHERE rt.ordinal >= central_user.min_ord
		AND rt.ordinal < central_user.max_ord
	`,
		pid.Value(),
		rankingOrderParam.Length.Value,
	)
	// undefined table is expected if rankingTable isn't existing
	if errors.Is(err, sql.ErrNoRows) || isUndefinedTable(err) {
		return nil, 0, nil
	} else if err != nil {
		return nil, 0, err
	}

	results, err := parseRankingDataList(rows)
	if err != nil {
		return nil, 0, err
	}

	// todo totalCount
	return results, uint32(results.Length()), nil
}
func GetFriendsRankingsAndCountByCategoryAndRankingOrderParam(pid *nextypes.PID, category *nextypes.PrimitiveU32, rankingOrderParam *types.RankingOrderParam) (*nextypes.List[*types.RankingRankData], uint32, error) {
	globals.Logger.Info(rankingOrderParam.FormatToString(1))
	return nil, 0, nil
}
func GetNearbyFriendsRankingsAndCountByCategoryAndRankingOrderParam(pid *nextypes.PID, category *nextypes.PrimitiveU32, rankingOrderParam *types.RankingOrderParam) (*nextypes.List[*types.RankingRankData], uint32, error) {
	globals.Logger.Info(rankingOrderParam.FormatToString(1))
	return nil, 0, nil
}
func GetOwnRankingByCategoryAndRankingOrderParam(pid *nextypes.PID, category *nextypes.PrimitiveU32, rankingOrderParam *types.RankingOrderParam) (*nextypes.List[*types.RankingRankData], uint32, error) {
	// todo filter by groups
	rankingTable := `ranking.ranks_` + strconv.Itoa(int(category.Value))
	rows, err := Postgres.Query(`
		SELECT
		    user_pid,
		    unique_id,
		    score,
		    groups,
		    param,
		    common_data,
		    rank,
		    update_date
		FROM `+rankingTable+`
		WHERE user_pid = $1
		LIMIT $2
	`,
		pid.Value(),
		rankingOrderParam.Length.Value,
	)
	// undefined table is expected if rankingTable isn't existing
	if errors.Is(err, sql.ErrNoRows) || isUndefinedTable(err) {
		return nil, 0, nil
	} else if err != nil {
		return nil, 0, err
	}

	results, err := parseRankingDataList(rows)
	if err != nil {
		return nil, 0, err
	} else if results.Length() != 1 {
		return nil, 0, nil
	}

	return results, uint32(results.Length()), nil
}

func createCategory(category uint32, golfScoring bool) error {
	now := time.Now()
	_, err := Postgres.Exec(`INSERT INTO ranking.categories (category, golf_scoring, creation_date)
		VALUES ($1, $2, $3) 
	`,
		category, golfScoring, now,
	)
	if err != nil {
		return err
	}

	// can convert to materialized view when ready
	// also yikes lmao
	order := "DESC"
	if golfScoring {
		order = "ASC"
	}
	_, err = Postgres.Exec(fmt.Sprintf(`CREATE VIEW ranking.ranks_%d AS
		SELECT
			user_pid,
			unique_id,
			score,
			groups,
			param,
			common_data,
			update_date,
			RANK() over (ORDER BY score %s) AS rank,
			ROW_NUMBER() over (ORDER BY score %s) AS ordinal
		FROM ranking.scores
		WHERE category = %d
	`, category, order, order, category))
	if err != nil {
		return err
	}

	globals.Logger.Infof("Created category %d (golf: %t)", category, golfScoring)
	return nil
}

func InsertRankingByPIDAndRankingScoreData(pid *nextypes.PID, rankingScoreData *types.RankingScoreData, uniqueID *nextypes.PrimitiveU64) error {
	globals.Logger.Info(rankingScoreData.FormatToString(1))
	now := time.Now()
	res, err := Postgres.Exec(`
		UPDATE ranking.scores SET score = $1, groups = $2, param = $3, update_date = $4
		WHERE category = $5 AND user_pid = $6
	`,
		rankingScoreData.Score.Value,
		rankingScoreData.Groups.Value,
		rankingScoreData.Param.Value,
		now,
		rankingScoreData.Category.Value,
		pid.Value(),
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 0 {
		// * All happy
		return nil
	}

	var categoryExists bool
	err = Postgres.QueryRow(`SELECT EXISTS(SELECT 1 FROM ranking.categories WHERE category = $1)`, rankingScoreData.Category.Value).Scan(&categoryExists)
	if err != nil {
		return err
	}

	if !categoryExists {
		err = createCategory(
			rankingScoreData.Category.Value,
			rankingScoreData.OrderBy.Value == 1,
		)
		if err != nil {
			return err
		}
	}

	_, err = Postgres.Exec(`
			INSERT INTO ranking.scores (user_pid, category, groups, score, param, creation_date, update_date) 
			VALUES ($1, $2, $3, $4, $5, $6, $6)
	`,
		pid.Value(),
		rankingScoreData.Category.Value,
		rankingScoreData.Groups.Value,
		rankingScoreData.Score.Value,
		rankingScoreData.Param.Value,
		now,
	)

	return err
}

func GetCommonData(uniqueID *nextypes.PrimitiveU64) (*nextypes.Buffer, error) {
	globals.Logger.Infof("GetCommonData %d", uniqueID.Value)
	return nil, nil
}

func UploadCommonData(pid *nextypes.PID, uniqueID *nextypes.PrimitiveU64, commonData *nextypes.Buffer) error {
	_, err := Postgres.Exec(`
		UPDATE ranking.scores SET common_data = $1 WHERE user_pid = $2
	`,
		commonData.Value,
		pid.Value(),
	)
	return err
}

func initRanking() {

}
