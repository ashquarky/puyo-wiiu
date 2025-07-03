package nex

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	commonmatchmaking "github.com/PretendoNetwork/nex-protocols-common-go/v2/match-making"
	commonmatchmakingext "github.com/PretendoNetwork/nex-protocols-common-go/v2/match-making-ext"
	commonmatchmakeextension "github.com/PretendoNetwork/nex-protocols-common-go/v2/matchmake-extension"
	commonranking "github.com/PretendoNetwork/nex-protocols-common-go/v2/ranking"
	commonsecure "github.com/PretendoNetwork/nex-protocols-common-go/v2/secure-connection"
	matchmaking "github.com/PretendoNetwork/nex-protocols-go/v2/match-making"
	matchmakingext "github.com/PretendoNetwork/nex-protocols-go/v2/match-making-ext"
	matchmakeextension "github.com/PretendoNetwork/nex-protocols-go/v2/matchmake-extension"
	ranking "github.com/PretendoNetwork/nex-protocols-go/v2/ranking"
	secure "github.com/PretendoNetwork/nex-protocols-go/v2/secure-connection"
	puyodatastore "github.com/PretendoNetwork/puyo-puyo-tetris/datastore"
	"github.com/PretendoNetwork/puyo-puyo-tetris/globals"

	commonnattraversal "github.com/PretendoNetwork/nex-protocols-common-go/v2/nat-traversal"
	nattraversal "github.com/PretendoNetwork/nex-protocols-go/v2/nat-traversal"

	matchmakingtypes "github.com/PretendoNetwork/nex-protocols-go/v2/match-making/types"
)

//func MatchmakeExtensionCloseParticipation(err error, packet nex.PacketInterface, callID uint32, gid *types.PrimitiveU32) (*nex.RMCMessage, *nex.Error) {
//	if err != nil {
//		globals.Logger.Error(err.Error())
//		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
//	}
//
//	session, ok := commonglobals.Sessions[gid.Value]
//	if !ok {
//		return nil, nex.NewError(nex.ResultCodes.RendezVous.SessionVoid, "change_error")
//	}
//
//	connection := packet.Sender().(*nex.PRUDPConnection)
//	endpoint := connection.Endpoint().(*nex.PRUDPEndPoint)
//
//	// * PUYOPUYOTETRIS has everyone send CloseParticipation here, not just the owner of the room.
//	// * So, if a non-owner asks, just lie and claim success without actually changing anything.
//	if !session.GameMatchmakeSession.Gathering.OwnerPID.Equals(connection.PID()) {
//		session.GameMatchmakeSession.OpenParticipation = types.NewPrimitiveBool(false)
//	}
//
//	rmcResponse := nex.NewRMCSuccess(endpoint, nil)
//	rmcResponse.ProtocolID = matchmakeextension.ProtocolID
//	rmcResponse.MethodID = matchmakeextension.MethodCloseParticipation
//	rmcResponse.CallID = callID
//
//	return rmcResponse, nil
//}

func CreateReportDBRecord(_ types.PID, _ types.UInt32, _ types.QBuffer) error {
	return nil
}

// TO DO:
// How do clubs work?
// GetObjectInfoByDataID
// UpdateObjectPeriodByDataIDWithPassword
// UpdateObjectMetaBinaryByDataIDWithPassword
// UpdateObjectDataTypeByDataIDWithPassword

func registerCommonSecureServerProtocols() {
	secureProtocol := secure.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(secureProtocol)
	commonSecureProtocol := commonsecure.NewCommonProtocol(secureProtocol)
	commonSecureProtocol.EnableInsecureRegister() // Game uses TicketGranting::LoginEx

	commonSecureProtocol.CreateReportDBRecord = CreateReportDBRecord

	// Ranking - ??
	rankingProtocol := ranking.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(rankingProtocol)
	commonRankingProtocol := commonranking.NewCommonProtocol(rankingProtocol)
	commonRankingProtocol.GetRankingsAndCountByCategoryAndRankingOrderParam = puyodatastore.GetRankingsAndCountByCategoryAndRankingOrderParam
	commonRankingProtocol.GetOwnRankingByCategoryAndRankingOrderParam = puyodatastore.GetOwnRankingByCategoryAndRankingOrderParam
	commonRankingProtocol.GetNearbyFriendsRankingsAndCountByCategoryAndRankingOrderParam = puyodatastore.GetNearbyFriendsRankingsAndCountByCategoryAndRankingOrderParam
	commonRankingProtocol.GetFriendsRankingsAndCountByCategoryAndRankingOrderParam = puyodatastore.GetFriendsRankingsAndCountByCategoryAndRankingOrderParam
	commonRankingProtocol.GetNearbyRankingsAndCountByCategoryAndRankingOrderParam = puyodatastore.GetNearbyRankingsAndCountByCategoryAndRankingOrderParam
	commonRankingProtocol.InsertRankingByPIDAndRankingScoreData = puyodatastore.InsertRankingByPIDAndRankingScoreData
	commonRankingProtocol.UploadCommonData = puyodatastore.UploadCommonData
	commonRankingProtocol.GetCommonData = puyodatastore.GetCommonData

	// Matchmaking stuff - National Puzzle League
	natTraversalProtocol := nattraversal.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(natTraversalProtocol)
	commonnattraversal.NewCommonProtocol(natTraversalProtocol)

	matchMakingProtocol := matchmaking.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(matchMakingProtocol)
	commonMatchMakingProtocol := commonmatchmaking.NewCommonProtocol(matchMakingProtocol)
	commonMatchMakingProtocol.SetManager(globals.MatchmakingManager)

	matchMakingExtProtocol := matchmakingext.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(matchMakingExtProtocol)
	commonMatchMakingExtProtocol := commonmatchmakingext.NewCommonProtocol(matchMakingExtProtocol)
	commonMatchMakingExtProtocol.SetManager(globals.MatchmakingManager)

	matchmakeExtensionProtocol := matchmakeextension.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(matchmakeExtensionProtocol)
	commonMatchmakeExtensionProtocol := commonmatchmakeextension.NewCommonProtocol(matchmakeExtensionProtocol)
	// * Handle custom CloseParticipation behaviour
	//matchmakeExtensionProtocol.SetHandlerCloseParticipation(MatchmakeExtensionCloseParticipation)
	commonMatchmakeExtensionProtocol.SetManager(globals.MatchmakingManager)

	commonMatchmakeExtensionProtocol.OnAfterAutoMatchmakeWithSearchCriteriaPostpone = func(packet nex.PacketInterface, lstSearchCriteria types.List[matchmakingtypes.MatchmakeSessionSearchCriteria], anyGathering matchmakingtypes.GatheringHolder, strMessage types.String) {
		globals.Logger.Info("Matchmake search criteria:")
		for _, criteria := range lstSearchCriteria {
			globals.Logger.Info(criteria.FormatToString(1))
		}

		//globals.Logger.Info("Active matchmaking sessions:")
		//for _, session := range commonglobals.Sessions {
		//	globals.Logger.Info(session.GameMatchmakeSession.FormatToString(1))
		//}
	}

}
