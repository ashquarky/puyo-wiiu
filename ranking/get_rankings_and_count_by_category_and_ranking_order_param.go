package ranking

import (
	"database/sql"
	"errors"
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	rankingtypes "github.com/PretendoNetwork/nex-protocols-go/v2/ranking/types"
	"time"
)

var getGlobalRankings *sql.Stmt

func GetRankingsAndCountByCategoryAndRankingOrderParam(category types.UInt32, rankingOrderParam rankingtypes.RankingOrderParam) (types.List[rankingtypes.RankingRankData], uint32, error) {
	// We are going to assume GLOBAL ranking until nex-protocols-common-go is improved to support the other types
	rows, err := getGlobalRankings.Query(category, rankingOrderParam.OrderCalculation == 0, rankingOrderParam.GroupIndex, rankingOrderParam.GroupNum, rankingOrderParam.Offset, rankingOrderParam.Length)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, nil
	} else if err != nil {
		return nil, 0, nex.NewError(nex.ResultCodes.Core.SystemError, err.Error())
	}
	defer rows.Close()

	totalCount := uint32(0)
	results := make(types.List[rankingtypes.RankingRankData], 0, rankingOrderParam.Length)
	for rows.Next() {
		result := rankingtypes.NewRankingRankData()
		var updateDate time.Time

		err := rows.Scan(
			&result.UniqueID,
			&result.PrincipalID,
			&result.Category,
			&result.Groups,
			&result.Score,
			&result.Param,
			&updateDate,
			&result.CommonData,
			&result.Order,
			&totalCount,
		)
		if err != nil {
			return nil, 0, nex.NewError(nex.ResultCodes.Core.SystemError, err.Error())
		}

		result.UpdateTime = result.UpdateTime.FromTimestamp(updateDate)

		results = append(results, result)
	}

	return results, totalCount, nil
}

func initGetGlobalRankingsStmt() error {
	stmt, err := Database.Prepare(`
		WITH cat AS (
		    SELECT golf_scoring FROM ranking.categories WHERE category = $1
		)
		SELECT
		    ranking.scores.unique_id,
		    ranking.scores.owner_pid,
		    ranking.scores.category,
		    ranking.scores.groups,
		    ranking.scores.score,
		    ranking.scores.param,
		    ranking.scores.update_date,
		    COALESCE(ranking.common_data.data, ''::bytea),
		    CASE WHEN $2 THEN
		        RANK() OVER (ORDER BY 
		            CASE WHEN cat.golf_scoring THEN ranking.scores.score END DESC,
					CASE WHEN NOT cat.golf_scoring THEN ranking.scores.score END ASC
				)
			ELSE
			    ROW_NUMBER() OVER (ORDER BY 
		            CASE WHEN cat.golf_scoring THEN ranking.scores.score END DESC,
					CASE WHEN NOT cat.golf_scoring THEN ranking.scores.score END ASC
				)
			END AS ord,
		    /* highly unfortunate */
			COUNT(*) OVER () AS count
		FROM cat, ranking.scores
		    LEFT OUTER JOIN ranking.common_data
			ON ranking.scores.unique_id = ranking.common_data.unique_id
		    AND ranking.scores.owner_pid = ranking.common_data.owner_pid
		WHERE
		    category = $1 AND
		    /* $3: GroupIndex; $4: GroupNum */
		    CASE WHEN $3 < 2 AND length(groups) >= 2 THEN get_byte(groups, $3) = $4 ELSE TRUE END 
		ORDER BY ord
		OFFSET $5
		LIMIT $6
	`)
	if err != nil {
		return err
	}

	getGlobalRankings = stmt
	return nil
}
