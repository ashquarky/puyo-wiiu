package ranking

import (
	commonranking "github.com/PretendoNetwork/nex-protocols-common-go/v2/ranking"
)

func NewRankingProtocol(protocol *commonranking.CommonProtocol) error {
	err := initDatabase()
	if err != nil {
		return err
	}

	protocol.InsertRankingByPIDAndRankingScoreData = InsertRankingByPIDAndRankingScoreData
	protocol.GetRankingsAndCountByCategoryAndRankingOrderParam = GetRankingsAndCountByCategoryAndRankingOrderParam
	protocol.GetNearbyRankingsAndCountByCategoryAndRankingOrderParam = GetNearbyRankingsAndCountByCategoryAndRankingOrderParam
	// protocol.GetFriendsRankingsAndCountByCategoryAndRankingOrderParam
	// protocol.GetNearbyFriendsRankingsAndCountByCategoryAndRankingOrderParam
	// protocol.GetOwnRankingByCategoryAndRankingOrderParam
	protocol.UploadCommonData = UploadCommonData

	return nil
}

//func GetCommonData(uniqueID types.UInt64) (types.Buffer, error)                            {}
