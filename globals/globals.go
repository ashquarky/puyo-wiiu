package globals

import (
	"database/sql"
	pb "github.com/PretendoNetwork/grpc/go/account"
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/datastore"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	"github.com/PretendoNetwork/plogger-go"
	"github.com/minio/minio-go/v7"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var Postgres *sql.DB
var MatchmakingManager *common_globals.MatchmakingManager

var Logger *plogger.Logger
var KerberosPassword = "password" // * Default password

var AuthenticationServer *nex.PRUDPServer
var AuthenticationEndpoint *nex.PRUDPEndPoint

var SecureServer *nex.PRUDPServer
var SecureEndpoint *nex.PRUDPEndPoint

var GRPCAccountClientConnection *grpc.ClientConn
var GRPCAccountClient pb.AccountClient
var GRPCAccountCommonMetadata metadata.MD

var MinIOClient *minio.Client
var Presigner *datastore.S3Presigner

var TokenAESKey []byte
var LocalAuthMode bool
