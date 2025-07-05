package nex

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	commondatastore "github.com/PretendoNetwork/nex-protocols-common-go/v2/datastore"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	commonmatchmaking "github.com/PretendoNetwork/nex-protocols-common-go/v2/match-making"
	commonmatchmakingext "github.com/PretendoNetwork/nex-protocols-common-go/v2/match-making-ext"
	commonmatchmakeextension "github.com/PretendoNetwork/nex-protocols-common-go/v2/matchmake-extension"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/matchmake-extension/database"
	commonranking "github.com/PretendoNetwork/nex-protocols-common-go/v2/ranking"
	commonsecure "github.com/PretendoNetwork/nex-protocols-common-go/v2/secure-connection"
	datastore "github.com/PretendoNetwork/nex-protocols-go/v2/datastore"
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

func MatchmakeExtensionCloseParticipation(err error, packet nex.PacketInterface, callID uint32, gid types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	connection := packet.Sender().(*nex.PRUDPConnection)
	manager := globals.MatchmakingManager
	endpoint := connection.Endpoint().(*nex.PRUDPEndPoint)

	manager.Mutex.Lock()

	session, _, nexError := database.GetMatchmakeSessionByID(manager, endpoint, uint32(gid))
	if nexError != nil {
		manager.Mutex.Unlock()
		return nil, nexError
	}

	// * PUYOPUYOTETRIS has everyone send CloseParticipation here, not just the owner of the room.
	// * So, if a non-owner asks, just lie and claim success without actually changing anything.
	if session.Gathering.OwnerPID.Equals(connection.PID()) {
		nexError = database.UpdateParticipation(manager, uint32(gid), false)
		if nexError != nil {
			manager.Mutex.Unlock()
			return nil, nexError
		}
	}

	manager.Mutex.Unlock()

	rmcResponse := nex.NewRMCSuccess(endpoint, nil)
	rmcResponse.ProtocolID = matchmakeextension.ProtocolID
	rmcResponse.MethodID = matchmakeextension.MethodCloseParticipation
	rmcResponse.CallID = callID

	return rmcResponse, nil
}

func CreateReportDBRecord(_ types.PID, _ types.UInt32, _ types.QBuffer) error {
	return nil
}

// TO DO:
// Persistent gatherings for clubs

func registerCommonSecureServerProtocols() {
	secureProtocol := secure.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(secureProtocol)
	commonSecureProtocol := commonsecure.NewCommonProtocol(secureProtocol)
	commonSecureProtocol.EnableInsecureRegister() // Game uses TicketGranting::LoginEx

	commonSecureProtocol.CreateReportDBRecord = CreateReportDBRecord

	// DataStore - player stats, replays, clubs
	datastoreProtocol := datastore.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(datastoreProtocol)
	commonDatastoreProtocol := commondatastore.NewCommonProtocol(datastoreProtocol)
	commonDatastoreProtocol.SetManager(globals.DatastoreManager)

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
	matchmakeExtensionProtocol.SetHandlerCloseParticipation(MatchmakeExtensionCloseParticipation)
	commonMatchmakeExtensionProtocol.SetManager(globals.MatchmakingManager)
	commonMatchmakeExtensionProtocol.CleanupMatchmakeSessionSearchCriterias = func(searchCriterias types.List[matchmakingtypes.MatchmakeSessionSearchCriteria]) {
		// lol ok
	}

	commonMatchmakeExtensionProtocol.OnAfterAutoMatchmakeWithSearchCriteriaPostpone = func(packet nex.PacketInterface, lstSearchCriteria types.List[matchmakingtypes.MatchmakeSessionSearchCriteria], anyGathering matchmakingtypes.GatheringHolder, strMessage types.String) {
		globals.Logger.Info("Matchmake search criteria:")
		for _, criteria := range lstSearchCriteria {
			globals.Logger.Info(criteria.FormatToString(1))
		}
	}
}
