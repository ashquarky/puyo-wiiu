package datastore

import (
	"database/sql"
	"errors"
	"github.com/PretendoNetwork/nex-go/v2"
	nextypes "github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	"github.com/PretendoNetwork/puyo-puyo-tetris/globals"
	"github.com/lib/pq"
	"runtime"
	"time"
)

// get own gamedata (among other things)
func GetObjectInfosByDataStoreSearchParam(param *types.DataStoreSearchParam, pid *nextypes.PID) ([]*types.DataStoreMetaInfo, uint32, *nex.Error) {
	// No clue. A guess.
	if param.DataType.Value == 65535 && param.SearchTarget.Value == 10 {
		var dataID uint64

		err := Postgres.QueryRow(`SELECT data_id from datastore.objects WHERE name = 'gamedata' AND owner = $1`, pid.Value()).Scan(&dataID)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, 0, nil
		} else if err != nil {
			// COOL FACTS: this crashes Puyo
			return nil, 0, nex.NewError(nex.ResultCodes.DataStore.SystemFileError, err.Error())
		}

		info, errCode := GetObjectInfoByDataID(nextypes.NewPrimitiveU64(dataID))
		if errCode != nil {
			return nil, 0, nil
		}

		result := make([]*types.DataStoreMetaInfo, 1)
		result[0] = info

		return result, 1, nil
	}

	globals.Logger.Warning("Stubbed GetObjectInfosByDataStoreSearchParam!")
	return nil, 0, nil
}

// part of GetMeta, not used
func GetObjectInfoByPersistenceTargetWithPassword(persistenceTarget *types.DataStorePersistenceTarget, password *nextypes.PrimitiveU64) (*types.DataStoreMetaInfo, *nex.Error) {
	globals.Logger.Warning("Actually, GetObjectInfoByPersistenceTargetWithPassword *is* used!")
	runtime.Breakpoint() // TODO disable in prod
	return nil, nil
}

// part of GetMeta
func GetObjectInfoByDataIDWithPassword(dataID *nextypes.PrimitiveU64, _ *nextypes.PrimitiveU64) (*types.DataStoreMetaInfo, *nex.Error) {
	// TODO check password
	return GetObjectInfoByDataID(dataID)
}

// part of ChangeMeta
func GetObjectInfoByDataID(dataID *nextypes.PrimitiveU64) (*types.DataStoreMetaInfo, *nex.Error) {
	result := types.NewDataStoreMetaInfo()

	var createdTime time.Time
	var updatedTime time.Time
	var ownerID uint64
	var permissionRecipients []uint64
	var delPermissionRecipients []uint64
	var tags []string

	result.ExpireTime = nextypes.NewDateTime(0x9C3f3E0000) // * 9999-12-31T00:00:00.000Z. This is what the real server sends

	err := Postgres.QueryRow(`SELECT
    	data_id,
    	owner,
    	size,
    	name,
    	data_type,
		meta_binary,
        permission,
     	permission_recipients,
     	delete_permission,
     	delete_permission_recipients,
     	flag,
        period,
     	refer_data_id,
     	tags,
     	creation_date,
     	update_date
	FROM datastore.objects WHERE data_id = $1`, dataID.Value).Scan(
		&result.DataID.Value,
		&ownerID,
		&result.Size.Value,
		&result.Name.Value,
		&result.DataType.Value,
		&result.MetaBinary.Value,
		&result.Permission.Permission.Value,
		pq.Array(&permissionRecipients),
		&result.DelPermission.Permission.Value,
		pq.Array(&delPermissionRecipients),
		&result.Flag.Value,
		&result.Period.Value,
		&result.ReferDataID.Value,
		pq.Array(&tags),
		&createdTime,
		&updatedTime,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nex.NewError(nex.ResultCodes.DataStore.NotFound, "Object not found")
	} else if err != nil {
		return nil, nex.NewError(nex.ResultCodes.DataStore.SystemFileError, err.Error())
	}

	result.OwnerID = nextypes.NewPID(ownerID)
	result.Permission.RecipientIDs = createPIDList(&permissionRecipients)
	result.DelPermission.RecipientIDs = createPIDList(&delPermissionRecipients)
	result.Tags = createStringList(&tags)
	result.CreatedTime.FromTimestamp(createdTime)
	result.UpdatedTime.FromTimestamp(updatedTime)
	result.ReferredTime.FromTimestamp(createdTime)

	return result, nil
}

// uploading gamedata meta object
func InitializeObjectByPreparePostParam(ownerPID *nextypes.PID, param *types.DataStorePreparePostParam) (uint64, *nex.Error) {
	now := time.Now()

	var dataID uint64

	err := Postgres.QueryRow(`INSERT INTO datastore.objects 
    (
    	owner,
     	size,
     	name,
     	data_type,
     	meta_binary,
        permission,
     	permission_recipients,
     	delete_permission,
     	delete_permission_recipients,
     	flag,
        period,
     	refer_data_id,
     	tags, 
     	persistence_slot_id,
     	extra_data,
     	creation_date,
     	update_date
    )
	VALUES (
	        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
	) RETURNING data_id`,
		ownerPID.Value(),
		param.Size.Value,
		param.Name.Value,
		param.DataType.Value,
		param.MetaBinary.Value,
		param.Permission.Permission.Value,
		pq.Array(convertPIDList(param.Permission.RecipientIDs)),
		param.DelPermission.Permission.Value,
		pq.Array(convertPIDList(param.DelPermission.RecipientIDs)),
		param.Flag.Value,
		param.Period.Value,
		param.ReferDataID.Value,
		pq.Array(convertStringList(param.Tags)),
		param.PersistenceInitParam.PersistenceSlotID.Value, // todo DeleteLastObject
		pq.Array(convertStringList(param.ExtraData)),
		now,
		now,
	).Scan(&dataID)
	if err != nil {
		return 0, nex.NewError(nex.ResultCodes.DataStore.SystemFileError, err.Error())
	}

	return dataID, nil
}

// maybe unused?
func InitializeObjectRatingWithSlot(_ uint64, _ *types.DataStoreRatingInitParamWithSlot) *nex.Error {
	globals.Logger.Warning("Actually, InitializeObjectRatingWithSlot *is* used!")
	runtime.Breakpoint() // TODO disable in prod
	return nil
}

// Parts of ChangeMeta
func UpdateObjectPeriodByDataIDWithPassword(dataID *nextypes.PrimitiveU64, dataType *nextypes.PrimitiveU16, password *nextypes.PrimitiveU64) *nex.Error {
	return nil
}

// Parts of ChangeMeta
func UpdateObjectMetaBinaryByDataIDWithPassword(dataID *nextypes.PrimitiveU64, metaBinary *nextypes.QBuffer, password *nextypes.PrimitiveU64) *nex.Error {
	return nil
}

// Parts of ChangeMeta
func UpdateObjectDataTypeByDataIDWithPassword(dataID *nextypes.PrimitiveU64, period *nextypes.PrimitiveU16, password *nextypes.PrimitiveU64) *nex.Error {
	return nil
}

// Helpers for nex types
func convertPIDList(list *nextypes.List[*nextypes.PID]) []uint64 {
	result := make([]uint64, list.Length())

	for i, pid := range list.Slice() {
		result[i] = pid.Value()
	}

	return result
}

func createPIDList(list *[]uint64) *nextypes.List[*nextypes.PID] {
	result := make([]*nextypes.PID, len(*list))

	for i, u := range *list {
		result[i] = nextypes.NewPID(u)
	}

	nexlist := nextypes.NewList[*nextypes.PID]()
	nexlist.SetFromData(result)

	return nexlist
}

func convertStringList(list *nextypes.List[*nextypes.String]) []string {
	result := make([]string, list.Length())

	for i, n := range list.Slice() {
		result[i] = n.Value
	}

	return result
}

func createStringList(list *[]string) *nextypes.List[*nextypes.String] {
	result := make([]*nextypes.String, len(*list))

	for i, u := range *list {
		result[i] = nextypes.NewString(u)
	}

	nexlist := nextypes.NewList[*nextypes.String]()
	nexlist.SetFromData(result)

	return nexlist
}
