package ranking

import (
	"database/sql"
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"time"
)

var insertCommonDataStmt *sql.Stmt

func UploadCommonData(pid types.PID, uniqueID types.UInt64, commonData types.Buffer) error {
	now := time.Now()
	_, err := insertCommonDataStmt.Exec(uniqueID, pid, commonData, now)
	if err != nil {
		return nex.NewError(nex.ResultCodes.Core.SystemError, err.Error())
	}

	return nil
}

func initInsertCommonDataStmt() error {
	stmt, err := Database.Prepare(`
		INSERT INTO ranking.common_data (unique_id, owner_pid, data, creation_date, update_date)
		VALUES ($1, $2, $3, $4, $4)
		ON CONFLICT (owner_pid, unique_id) DO UPDATE
			SET data = excluded.data, update_date = excluded.update_date;
	`)
	if err != nil {
		return err
	}

	insertCommonDataStmt = stmt
	return nil
}
